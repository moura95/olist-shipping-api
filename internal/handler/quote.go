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

type QuoteHandler struct {
	packageService *service.PackageService
	config         *config.Config
	logger         *zap.SugaredLogger
	validate       *validator.Validate
}

func NewQuoteHandler(packageService *service.PackageService, cfg *config.Config, logger *zap.SugaredLogger) *QuoteHandler {
	validate := validator.New()
	customValidator.SetupCustomValidators(validate)
	return &QuoteHandler{
		packageService: packageService,
		config:         cfg,
		logger:         logger,
		validate:       validate,
	}
}

func (h *QuoteHandler) GetQuotes(ctx *gin.Context) {
	logger := middleware.GetLoggerFromContext(ctx)
	logger.Info("get quotes started")

	var query v1.GetQuotesQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		logger.Errorw("bind query failed", "error", err)
		v1.HandleBadRequest(ctx, "Invalid query parameters")
		return
	}

	if err := h.validate.Struct(query); err != nil {
		logger.Errorw("validation failed", "error", err)
		v1.HandleValidationError(ctx, err)
		return
	}

	quotes, err := h.packageService.GetQuotes(ctx, query.StateCode, query.WeightKg)
	if err != nil {
		logger.Errorw("get quotes failed", "error", err)
		v1.HandleInternalError(ctx, fmt.Errorf("get quotes: %v", err).Error())
		return
	}

	var resp []v1.QuoteResponse
	for _, quote := range quotes {
		resp = append(resp, v1.QuoteResponse{
			CarrierName:           &quote.Carier,
			EstimatedPrice:        &quote.EstimatedPrice,
			EstimatedDeliveryDays: &quote.EstimatedDeliveryDays,
		})
	}

	logger.Infow("get quotes completed", "state_code", query.StateCode, "weight_kg", query.WeightKg, "quotes_count", len(resp))
	v1.HandleSuccess(ctx, resp)
}
