package v1

type PackageResponse struct {
	ID                *string  `json:"id"`
	TrackingCode      *string  `json:"tracking_code"`
	Product           *string  `json:"product"`
	WeightKg          *float64 `json:"weight_kg"`
	DestinationState  *string  `json:"destination_state"`
	Status            *string  `json:"status"`
	HiredCarrierID    *string  `json:"hired_carrier_id"`
	HiredPrice        *string  `json:"hired_price"`
	HiredDeliveryDays *int32   `json:"hired_delivery_days"`
	CreatedAt         *string  `json:"created_at"`
	UpdatedAt         *string  `json:"updated_at"`
}

type CreatePackageRequest struct {
	Product          string  `json:"product" validate:"required"`
	WeightKg         float64 `json:"weight_kg" validate:"required,gt=0"`
	DestinationState string  `json:"destination_state" validate:"required,len=2"`
}

type UpdatePackageStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=created awaiting_pickup picked_up shipped delivered lost"`
}

type HireCarrierRequest struct {
	CarrierID    string `json:"carrier_id" validate:"required,uuid"`
	Price        string `json:"price" validate:"required"`
	DeliveryDays int32  `json:"delivery_days" validate:"required,gt=0"`
}

type QuoteResponse struct {
	CarrierName           *string  `json:"carrier_name"`
	EstimatedPrice        *float64 `json:"estimated_price"`
	EstimatedDeliveryDays *int32   `json:"estimated_delivery_days"`
}

type GetQuotesQuery struct {
	StateCode string  `form:"state_code" validate:"required,len=2"`
	WeightKg  float64 `form:"weight_kg" validate:"required,gt=0"`
}

type CarrierResponse struct {
	ID        *string `json:"id"`
	Name      *string `json:"name"`
	CreatedAt *string `json:"created_at"`
}

type StateResponse struct {
	Code       *string `json:"code"`
	Name       *string `json:"name"`
	RegionName *string `json:"region_name"`
}
