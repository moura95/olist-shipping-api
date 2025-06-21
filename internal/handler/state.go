package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github/moura95/olist-shipping-api/api/v1"
	"github/moura95/olist-shipping-api/config"
	"github/moura95/olist-shipping-api/internal/middleware"
	"github/moura95/olist-shipping-api/internal/service"
	customValidator "github/moura95/olist-shipping-api/pkg/validator"
	"go.uber.org/zap"
)

type StateHandler struct {
	packageService *service.PackageService
	config         *config.Config
	logger         *zap.SugaredLogger
	validate       *validator.Validate
}

func NewStateHandler(packageService *service.PackageService, cfg *config.Config, logger *zap.SugaredLogger) *StateHandler {
	validate := validator.New()
	customValidator.SetupCustomValidators(validate)
	return &StateHandler{
		packageService: packageService,
		config:         cfg,
		logger:         logger,
		validate:       validate,
	}
}

func (h *StateHandler) List(ctx *gin.Context) {
	logger := middleware.GetLoggerFromContext(ctx)
	logger.Info("list states started")

	states, err := h.packageService.GetStates(ctx)
	if err != nil {
		logger.Errorw("list states failed", "error", err)
		v1.HandleInternalError(ctx, fmt.Errorf("list states: %v", err).Error())
		return
	}

	var resp []v1.StateResponse
	for _, state := range states {
		resp = append(resp, v1.StateResponse{
			Code:       &state.Code,
			Name:       &state.Name,
			RegionName: &state.RegionName,
		})
	}

	logger.Infow("list states completed", "count", len(resp))
	v1.HandleSuccess(ctx, resp)
}
