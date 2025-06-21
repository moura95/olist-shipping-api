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
	customValidator "github/moura95/olist-shipping-api/pkg/validator"
	"go.uber.org/zap"
)

type PackageHandler struct {
	packageService *service.PackageService
	config         *config.Config
	logger         *zap.SugaredLogger
	validate       *validator.Validate
}

func NewPackageHandler(packageService *service.PackageService, cfg *config.Config, logger *zap.SugaredLogger) *PackageHandler {
	validate := validator.New()
	customValidator.SetupCustomValidators(validate)
	return &PackageHandler{
		packageService: packageService,
		config:         cfg,
		logger:         logger,
		validate:       validate,
	}
}

// List godoc
// @Summary      List all packages
// @Description  Get all packages with pagination support
// @Tags         packages
// @Accept       json
// @Produce      json
// @Success      200  {object}  v1.Response{data=[]v1.PackageResponse}
// @Failure      500  {object}  v1.Response
// @Router       /packages [get]
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
			TrackingCode:      util.NullStringToPtr(pkg.TrackingCode),
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

// GetByID godoc
// @Summary      Get package by ID
// @Description  Get package details by package ID
// @Tags         packages
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Package ID"
// @Success      200  {object}  v1.Response{data=v1.PackageResponse}
// @Failure      400  {object}  v1.Response
// @Failure      404  {object}  v1.Response
// @Router       /packages/{id} [get]
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
		TrackingCode:      util.NullStringToPtr(pkg.TrackingCode),
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

// GetByTrackingCode godoc
// @Summary      Get package by tracking code
// @Description  Get package details by tracking code
// @Tags         packages
// @Accept       json
// @Produce      json
// @Param        tracking_code  path      string  true  "Tracking Code"
// @Success      200           {object}  v1.Response{data=v1.PackageResponse}
// @Failure      400           {object}  v1.Response
// @Failure      404           {object}  v1.Response
// @Router       /packages/tracking/{tracking_code} [get]
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
		TrackingCode:      util.NullStringToPtr(pkg.TrackingCode),
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

// Create godoc
// @Summary      Create a new package
// @Description  Create a new package for shipping
// @Tags         packages
// @Accept       json
// @Produce      json
// @Param        request  body      v1.CreatePackageRequest  true  "Package data"
// @Success      201      {object}  v1.Response{data=v1.PackageResponse}
// @Failure      400      {object}  v1.Response
// @Failure      500      {object}  v1.Response
// @Router       /packages [post]
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
		v1.HandleValidationError(ctx, err)
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
		TrackingCode:      util.NullStringToPtr(pkg.TrackingCode),
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

// UpdateStatus godoc
// @Summary      Update package status
// @Description  Update the status of a package
// @Tags         packages
// @Accept       json
// @Produce      json
// @Param        id       path  string                        true  "Package ID"
// @Param        request  body  v1.UpdatePackageStatusRequest  true  "Status data"
// @Success      204      "No Content"
// @Failure      400      {object}  v1.Response
// @Failure      500      {object}  v1.Response
// @Router       /packages/{id}/status [patch]
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
		v1.HandleValidationError(ctx, err)
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

// Delete godoc
// @Summary      Delete a package
// @Description  Delete a package by ID
// @Tags         packages
// @Accept       json
// @Produce      json
// @Param        id  path      string  true  "Package ID"
// @Success      200  {object}  v1.Response
// @Failure      400  {object}  v1.Response
// @Failure      500  {object}  v1.Response
// @Router       /packages/{id} [delete]
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

// HireCarrier godoc
// @Summary      Hire carrier for package
// @Description  Hire a carrier to deliver the package
// @Tags         packages
// @Accept       json
// @Produce      json
// @Param        id       path  string                 true  "Package ID"
// @Param        request  body  v1.HireCarrierRequest  true  "Carrier hire data"
// @Success      200      {object}  v1.Response
// @Failure      400      {object}  v1.Response
// @Failure      500      {object}  v1.Response
// @Router       /packages/{id}/hire [post]
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
		v1.HandleValidationError(ctx, err)
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
