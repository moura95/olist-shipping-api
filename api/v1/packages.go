package v1

type PackageResponse struct {
	ID                *string  `json:"id"`
	TrackingCode      *string  `json:"codigo_rastreio"`
	Product           *string  `json:"produto"`
	WeightKg          *float64 `json:"peso_kg"`
	DestinationState  *string  `json:"estado_destino"`
	Status            *string  `json:"status"`
	HiredCarrierID    *string  `json:"transportadora_id"`
	HiredPrice        *string  `json:"preco_contratado"`
	HiredDeliveryDays *int32   `json:"prazo_contratado_dias"`
	CreatedAt         *string  `json:"criado_em"`
	UpdatedAt         *string  `json:"atualizado_em"`
}

type CreatePackageRequest struct {
	Product          string  `json:"produto" validate:"required"`
	WeightKg         float64 `json:"peso_kg" validate:"required,gt=0"`
	DestinationState string  `json:"estado_destino" validate:"required,len=2,brazilian_state"`
}

type UpdatePackageStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=criado esperando_coleta coletado enviado entregue extraviado"`
}

type HireCarrierRequest struct {
	CarrierID    string `json:"transportadora_id" validate:"required,uuid"`
	Price        string `json:"preco" validate:"required"`
	DeliveryDays int32  `json:"prazo_dias" validate:"required,gt=0"`
}

type QuoteResponse struct {
	CarrierName           *string  `json:"transportadora"`
	EstimatedPrice        *float64 `json:"preco_estimado"`
	EstimatedDeliveryDays *int32   `json:"prazo_estimado_dias"`
}

type GetQuotesQuery struct {
	StateCode string  `form:"estado_destino" validate:"required,len=2,brazilian_state"`
	WeightKg  float64 `form:"peso_kg" validate:"required,gt=0"`
}

type CarrierResponse struct {
	ID        *string `json:"id"`
	Name      *string `json:"nome"`
	CreatedAt *string `json:"criado_em"`
}

type StateResponse struct {
	Code       *string `json:"codigo"`
	Name       *string `json:"nome"`
	RegionName *string `json:"nome_regiao"`
}
