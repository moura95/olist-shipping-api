package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github/moura95/olist-shipping-api/api/v1"
	"github/moura95/olist-shipping-api/config"
	"github/moura95/olist-shipping-api/internal/middleware"
	"github/moura95/olist-shipping-api/internal/service"
	"github/moura95/olist-shipping-api/internal/util"
	"go.uber.org/zap"
)

type PackageHandler struct {
	packageService *service.PackageService
	config         *config.Config
	logger         *zap.SugaredLogger
	validate       *validator.Validate
}

func NewPackageHandler(packageService *service.PackageService, cfg *config.Config, logger *zap.SugaredLogger) *PackageHandler {
	return &PackageHandler{
		packageService: packageService,
		config:         cfg,
		logger:         logger,
		validate:       validator.New(),
	}
}

func (h *PackageHandler) List(ctx *gin.Context) {
	logger := middleware.GetLoggerFromContext(ctx)
	logger.Info("list packages started")

	packages, err := h.packageService.GetAll(ctx)
	if err != nil {
		logger.Errorw("list packages failed", "error", err)
		v1.HandleInternalError(ctx, fmt.Errorf("list packages: %v", err).Error())
		return
	}

	var resp []v1.PackageResponse
	for _, pkg := range packages {
		var createdAt, updatedAt *string
		if pkg.CreatedAt.Valid {
			formatted := pkg.CreatedAt.Time.Format(time.RFC3339)
			createdAt = &formatted
		}
		if pkg.UpdatedAt.Valid {
			formatted := pkg.UpdatedAt.Time.Format(time.RFC3339)
			updatedAt = &formatted
		}

		pkgID := pkg.ID.String()
		var hiredCarrierID *string
		if pkg.HiredCarrierID.Valid {
			carrierID := pkg.HiredCarrierID.UUID.String()
			hiredCarrierID = &carrierID
		}

		resp = append(resp, v1.PackageResponse{
			ID:                &pkgID,
			TrackingCode:      &pkg.TrackingCode,
			Product:           &pkg.Product,
			WeightKg:          &pkg.WeightKg,
			DestinationState:  &pkg.DestinationState,
			Status:            &pkg.Status,
			HiredCarrierID:    hiredCarrierID,
			HiredPrice:        util.NullStringToPtr(pkg.HiredPrice),
			HiredDeliveryDays: util.NullInt32ToPtr(pkg.HiredDeliveryDays),
			CreatedAt:         createdAt,
			UpdatedAt:         updatedAt,
		})
	}

	logger.Infow("list packages completed", "count", len(resp))
	v1.HandleSuccess(ctx, resp)
}

func (h *PackageHandler) GetByID(ctx *gin.Context) {
	logger := middleware.GetLoggerFromContext(ctx)
	logger.Info("get package by id started")

	id := ctx.Param("id")
	if id == "" {
		logger.Errorw("package id is required")
		v1.HandleBadRequest(ctx, "Package ID is required")
		return
	}

	pkg, err := h.packageService.GetByID(ctx, id)
	if err != nil {
		logger.Errorw("get package by id failed", "error", err, "id", id)
		v1.HandleNotFound(ctx, fmt.Errorf("get package by id: %v", err).Error())
		return
	}

	var createdAt, updatedAt *string
	if pkg.CreatedAt.Valid {
		formatted := pkg.CreatedAt.Time.Format(time.RFC3339)
		createdAt = &formatted
	}
	if pkg.UpdatedAt.Valid {
		formatted := pkg.UpdatedAt.Time.Format(time.RFC3339)
		updatedAt = &formatted
	}

	pkgID := pkg.ID.String()
	var hiredCarrierID *string
	if pkg.HiredCarrierID.Valid {
		carrierID := pkg.HiredCarrierID.UUID.String()
		hiredCarrierID = &carrierID
	}

	response := v1.PackageResponse{
		ID:                &pkgID,
		TrackingCode:      &pkg.TrackingCode,
		Product:           &pkg.Product,
		WeightKg:          &pkg.WeightKg,
		DestinationState:  &pkg.DestinationState,
		Status:            &pkg.Status,
		HiredCarrierID:    hiredCarrierID,
		HiredPrice:        util.NullStringToPtr(pkg.HiredPrice),
		HiredDeliveryDays: util.NullInt32ToPtr(pkg.HiredDeliveryDays),
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
	}

	logger.Infow("get package by id completed", "id", id)
	v1.HandleSuccess(ctx, response)
}

func (h *PackageHandler) GetByTrackingCode(ctx *gin.Context) {
	logger := middleware.GetLoggerFromContext(ctx)
	logger.Info("get package by tracking code started")

	trackingCode := ctx.Param("tracking_code")
	if trackingCode == "" {
		logger.Errorw("tracking code is required")
		v1.HandleBadRequest(ctx, "Tracking code is required")
		return
	}

	pkg, err := h.packageService.GetByTrackingCode(ctx, trackingCode)
	if err != nil {
		logger.Errorw("get package by tracking code failed", "error", err, "tracking_code", trackingCode)
		v1.HandleNotFound(ctx, fmt.Errorf("get package by tracking code: %v", err).Error())
		return
	}

	var createdAt, updatedAt *string
	if pkg.CreatedAt.Valid {
		formatted := pkg.CreatedAt.Time.Format(time.RFC3339)
		createdAt = &formatted
	}
	if pkg.UpdatedAt.Valid {
		formatted := pkg.UpdatedAt.Time.Format(time.RFC3339)
		updatedAt = &formatted
	}

	pkgID := pkg.ID.String()
	var hiredCarrierID *string
	if pkg.HiredCarrierID.Valid {
		carrierID := pkg.HiredCarrierID.UUID.String()
		hiredCarrierID = &carrierID
	}

	response := v1.PackageResponse{
		ID:                &pkgID,
		TrackingCode:      &pkg.TrackingCode,
		Product:           &pkg.Product,
		WeightKg:          &pkg.WeightKg,
		DestinationState:  &pkg.DestinationState,
		Status:            &pkg.Status,
		HiredCarrierID:    hiredCarrierID,
		HiredPrice:        util.NullStringToPtr(pkg.HiredPrice),
		HiredDeliveryDays: util.NullInt32ToPtr(pkg.HiredDeliveryDays),
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
	}

	logger.Infow("get package by tracking code completed", "tracking_code", trackingCode)
	v1.HandleSuccess(ctx, response)
}

func (h *PackageHandler) Create(ctx *gin.Context) {
	logger := middleware.GetLoggerFromContext(ctx)
	logger.Info("create package started")

	var req v1.CreatePackageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Errorw("bind json failed", "error", err)
		v1.HandleBadRequest(ctx, "Invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		logger.Errorw("validation failed", "error", err)
		v1.HandleValidationError(ctx, err.Error())
		return
	}

	pkg, err := h.packageService.Create(ctx, req.Product, req.WeightKg, req.DestinationState)
	if err != nil {
		logger.Errorw("create package failed", "error", err)
		v1.HandleInternalError(ctx, fmt.Errorf("create package: %v", err).Error())
		return
	}

	var createdAt, updatedAt *string
	if pkg.CreatedAt.Valid {
		formatted := pkg.CreatedAt.Time.Format(time.RFC3339)
		createdAt = &formatted
	}
	if pkg.UpdatedAt.Valid {
		formatted := pkg.UpdatedAt.Time.Format(time.RFC3339)
		updatedAt = &formatted
	}

	pkgID := pkg.ID.String()
	response := v1.PackageResponse{
		ID:                &pkgID,
		TrackingCode:      &pkg.TrackingCode,
		Product:           &pkg.Product,
		WeightKg:          &pkg.WeightKg,
		DestinationState:  &pkg.DestinationState,
		Status:            &pkg.Status,
		HiredCarrierID:    nil,
		HiredPrice:        util.NullStringToPtr(pkg.HiredPrice),
		HiredDeliveryDays: util.NullInt32ToPtr(pkg.HiredDeliveryDays),
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
	}

	logger.Infow("create package completed", "id", pkg.ID, "tracking_code", pkg.TrackingCode)
	v1.HandleCreated(ctx, response)
}

func (h *PackageHandler) UpdateStatus(ctx *gin.Context) {
	logger := middleware.GetLoggerFromContext(ctx)
	logger.Info("update package status started")

	var req v1.UpdatePackageStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Errorw("bind json failed", "error", err)
		v1.HandleBadRequest(ctx, "Invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		logger.Errorw("validation failed", "error", err)
		v1.HandleValidationError(ctx, err.Error())
		return
	}

	id := ctx.Param("id")
	if id == "" {
		logger.Errorw("package id is required")
		v1.HandleBadRequest(ctx, "Package ID is required")
		return
	}

	err := h.packageService.UpdateStatus(ctx, id, req.Status)
	if err != nil {
		logger.Errorw("update package status failed", "error", err, "id", id)
		v1.HandleInternalError(ctx, fmt.Errorf("update package status: %v", err).Error())
		return
	}

	logger.Infow("update package status completed", "id", id, "status", req.Status)
	ctx.Status(http.StatusNoContent)
}

func (h *PackageHandler) Delete(ctx *gin.Context) {
	logger := middleware.GetLoggerFromContext(ctx)
	logger.Info("delete package started")

	id := ctx.Param("id")
	if id == "" {
		logger.Errorw("package id is required")
		v1.HandleBadRequest(ctx, "Package ID is required")
		return
	}

	err := h.packageService.Delete(ctx, id)
	if err != nil {
		logger.Errorw("delete package failed", "error", err, "id", id)
		v1.HandleInternalError(ctx, fmt.Errorf("delete package: %v", err).Error())
		return
	}

	logger.Infow("delete package completed", "id", id)
	v1.HandleSuccess(ctx, "Package deleted successfully")
}

func (h *PackageHandler) GetQuotes(ctx *gin.Context) {
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
		v1.HandleValidationError(ctx, err.Error())
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

func (h *PackageHandler) HireCarrier(ctx *gin.Context) {
	logger := middleware.GetLoggerFromContext(ctx)
	logger.Info("hire carrier started")

	var req v1.HireCarrierRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		logger.Errorw("bind json failed", "error", err)
		v1.HandleBadRequest(ctx, "Invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		logger.Errorw("validation failed", "error", err)
		v1.HandleValidationError(ctx, err.Error())
		return
	}

	id := ctx.Param("id")
	if id == "" {
		logger.Errorw("package id is required")
		v1.HandleBadRequest(ctx, "Package ID is required")
		return
	}

	err := h.packageService.HireCarrier(ctx, id, req.CarrierID, req.Price, req.DeliveryDays)
	if err != nil {
		logger.Errorw("hire carrier failed", "error", err, "id", id, "carrier_id", req.CarrierID)
		v1.HandleInternalError(ctx, fmt.Errorf("hire carrier: %v", err).Error())
		return
	}

	logger.Infow("hire carrier completed", "id", id, "carrier_id", req.CarrierID)
	v1.HandleSuccess(ctx, "Carrier hired successfully")
}

func (h *PackageHandler) ListCarriers(ctx *gin.Context) {
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

func (h *PackageHandler) ListStates(ctx *gin.Context) {
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
