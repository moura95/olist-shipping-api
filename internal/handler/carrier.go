package handler

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github/moura95/olist-shipping-api/api/v1"
	"github/moura95/olist-shipping-api/config"
	"github/moura95/olist-shipping-api/internal/middleware"
	"github/moura95/olist-shipping-api/internal/service"
	customValidator "github/moura95/olist-shipping-api/pkg/validator"
	"go.uber.org/zap"
)

type CarrierHandler struct {
	packageService *service.PackageService
	config         *config.Config
	logger         *zap.SugaredLogger
	validate       *validator.Validate
}

func NewCarrierHandler(packageService *service.PackageService, cfg *config.Config, logger *zap.SugaredLogger) *CarrierHandler {
	validate := validator.New()
	customValidator.SetupCustomValidators(validate)
	return &CarrierHandler{
		packageService: packageService,
		config:         cfg,
		logger:         logger,
		validate:       validate,
	}
}

func (h *CarrierHandler) List(ctx *gin.Context) {
	logger := middleware.GetLoggerFromContext(ctx)
	logger.Info("list carriers started")

	carriers, err := h.packageService.GetCarriers(ctx)
	if err != nil {
		logger.Errorw("list carriers failed", "error", err)
		v1.HandleInternalError(ctx, fmt.Errorf("list carriers: %v", err).Error())
		return
	}

	var resp []v1.CarrierResponse
	for _, carrier := range carriers {
		var createdAt *string
		if carrier.CreatedAt.Valid {
			formatted := carrier.CreatedAt.Time.Format(time.RFC3339)
			createdAt = &formatted
		}

		carrierID := carrier.ID.String()
		resp = append(resp, v1.CarrierResponse{
			ID:        &carrierID,
			Name:      &carrier.Name,
			CreatedAt: createdAt,
		})
	}

	logger.Infow("list carriers completed", "count", len(resp))
	v1.HandleSuccess(ctx, resp)
}
