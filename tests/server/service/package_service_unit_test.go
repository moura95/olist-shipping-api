package service

import (
	"context"
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github/moura95/olist-shipping-api/config"
	"github/moura95/olist-shipping-api/internal/repository"
	"github/moura95/olist-shipping-api/internal/service"
	"go.uber.org/zap"
)

func TestPackageService_Create(t *testing.T) {
	tests := []struct {
		name             string
		product          string
		weightKg         float64
		destinationState string
		setupMocked      func(repo *repository.QuerierMocked)
		expectedError    string
	}{
		{
			name:             "Create package successfully",
			product:          "Test Product",
			weightKg:         2.5,
			destinationState: "SP",
			setupMocked: func(repo *repository.QuerierMocked) {
				expectedPackage := repository.Package{
					ID: uuid.New(),
					TrackingCode: sql.NullString{
						Valid: false,
					},
					Product:          "Test Product",
					WeightKg:         2.5,
					DestinationState: "SP",
					Status:           "criado",
					CreatedAt:        sql.NullTime{Valid: true},
				}

				repo.On("CreatePackage", mock.Anything, mock.MatchedBy(func(arg repository.CreatePackageParams) bool {
					return arg.Product == "Test Product" &&
						arg.WeightKg == 2.5 &&
						arg.DestinationState == "SP"
				})).Return(expectedPackage, nil)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMocked := repository.NewQuerierMocked(t)
			tt.setupMocked(repoMocked)

			logger := zap.NewNop().Sugar()
			packageService := service.NewPackageService(repoMocked, config.Config{}, logger)

			result, err := packageService.Create(context.Background(), tt.product, tt.weightKg, tt.destinationState)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.product, result.Product)
				assert.Equal(t, tt.weightKg, result.WeightKg)
				assert.Equal(t, tt.destinationState, result.DestinationState)
				assert.False(t, result.TrackingCode.Valid)
				assert.Equal(t, "criado", result.Status)
			}
		})
	}
}

func TestPackageService_GetByID(t *testing.T) {
	tests := []struct {
		name          string
		packageID     string
		setupMocked   func(repo *repository.QuerierMocked)
		expectedError string
	}{
		{
			name:      "Get package by valid ID",
			packageID: "550e8400-e29b-41d4-a716-446655440000",
			setupMocked: func(repo *repository.QuerierMocked) {
				expectedUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
				expectedPackage := repository.Package{
					ID:               expectedUUID,
					TrackingCode:     sql.NullString{String: "BR12345678", Valid: true},
					Product:          "Test Product",
					WeightKg:         2.5,
					DestinationState: "SP",
					Status:           "criado",
				}

				repo.On("GetPackageById", mock.Anything, expectedUUID).Return(expectedPackage, nil)
			},
		},
		{
			name:          "Get package with invalid UUID",
			packageID:     "invalid-uuid",
			setupMocked:   func(repo *repository.QuerierMocked) {},
			expectedError: "parse package id",
		},
		{
			name:      "Get package not found",
			packageID: "550e8400-e29b-41d4-a716-446655440000",
			setupMocked: func(repo *repository.QuerierMocked) {
				expectedUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
				repo.On("GetPackageById", mock.Anything, expectedUUID).Return(repository.Package{}, sql.ErrNoRows)
			},
			expectedError: "get package by id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMocked := repository.NewQuerierMocked(t)
			tt.setupMocked(repoMocked)

			logger := zap.NewNop().Sugar()
			packageService := service.NewPackageService(repoMocked, config.Config{}, logger)

			result, err := packageService.GetByID(context.Background(), tt.packageID)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if result.TrackingCode.Valid {
					assert.Equal(t, "BR12345678", result.TrackingCode.String)
				}
				assert.Equal(t, "Test Product", result.Product)
				assert.Equal(t, 2.5, result.WeightKg)
			}
		})
	}
}

func TestPackageService_GetByTrackingCode(t *testing.T) {
	tests := []struct {
		name          string
		trackingCode  string
		setupMocked   func(repo *repository.QuerierMocked)
		expectedError string
	}{
		{
			name:         "Get package by valid tracking code",
			trackingCode: "BR12345678",
			setupMocked: func(repo *repository.QuerierMocked) {
				expectedPackage := repository.Package{
					ID:               uuid.New(),
					TrackingCode:     sql.NullString{String: "BR12345678", Valid: true},
					Product:          "Test Product",
					WeightKg:         2.5,
					DestinationState: "SP",
					Status:           "criado",
				}

				repo.On("GetPackageByTrackingCode", mock.Anything, sql.NullString{String: "BR12345678", Valid: true}).Return(expectedPackage, nil)

			},
		},
		{
			name:         "Get package by tracking code not found",
			trackingCode: "BR99999999",
			setupMocked: func(repo *repository.QuerierMocked) {
				repo.On("GetPackageByTrackingCode", mock.Anything, sql.NullString{String: "BR99999999", Valid: true}).Return(repository.Package{}, sql.ErrNoRows)

			},
			expectedError: "get package by tracking code",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMocked := repository.NewQuerierMocked(t)
			tt.setupMocked(repoMocked)

			logger := zap.NewNop().Sugar()
			packageService := service.NewPackageService(repoMocked, config.Config{}, logger)

			result, err := packageService.GetByTrackingCode(context.Background(), tt.trackingCode)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				if result.TrackingCode.Valid {
					assert.Equal(t, "BR12345678", result.TrackingCode.String)
				}
				assert.Equal(t, "Test Product", result.Product)
			}
		})
	}
}

func TestPackageService_GetAll(t *testing.T) {
	tests := []struct {
		name          string
		setupMocked   func(repo *repository.QuerierMocked)
		expectedCount int
		expectedError string
	}{
		{
			name: "Get all packages successfully",
			setupMocked: func(repo *repository.QuerierMocked) {
				expectedPackages := []repository.Package{
					{
						ID:               uuid.New(),
						TrackingCode:     sql.NullString{Valid: false},
						Product:          "Product 1",
						WeightKg:         1.0,
						DestinationState: "SP",
						Status:           "criado",
					},
					{
						ID: uuid.New(),
						TrackingCode: sql.NullString{
							String: "BR87654321",
							Valid:  true,
						},
						Product:          "Product 2",
						WeightKg:         2.0,
						DestinationState: "RJ",
						Status:           "enviado",
					},
				}

				repo.On("ListPackages", mock.Anything).Return(expectedPackages, nil)
			},
			expectedCount: 2,
		},
		{
			name: "Get all packages with empty result",
			setupMocked: func(repo *repository.QuerierMocked) {
				repo.On("ListPackages", mock.Anything).Return([]repository.Package{}, nil)
			},
			expectedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMocked := repository.NewQuerierMocked(t)
			tt.setupMocked(repoMocked)

			logger := zap.NewNop().Sugar()
			packageService := service.NewPackageService(repoMocked, config.Config{}, logger)

			result, err := packageService.GetAll(context.Background())

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCount, len(result))
			}
		})
	}
}

func TestPackageService_UpdateStatus(t *testing.T) {
	tests := []struct {
		name          string
		packageID     string
		status        string
		setupMocked   func(repo *repository.QuerierMocked)
		expectedError string
	}{
		{
			name:      "Update package status to coletado",
			packageID: "550e8400-e29b-41d4-a716-446655440000",
			status:    "coletado",
			setupMocked: func(repo *repository.QuerierMocked) {
				expectedUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
				expectedParams := repository.UpdatePackageStatusParams{
					ID:     expectedUUID,
					Status: "coletado",
				}

				repo.On("UpdatePackageStatus", mock.Anything, expectedParams).Return(nil)
			},
		},
		{
			name:      "Update package status to enviado with tracking code",
			packageID: "550e8400-e29b-41d4-a716-446655440000",
			status:    "enviado",
			setupMocked: func(repo *repository.QuerierMocked) {
				expectedUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")

				repo.On("TrackingCodeExists", mock.Anything, mock.AnythingOfType("sql.NullString")).Return(false, nil)

				repo.On("UpdatePackageStatusWithTracking", mock.Anything, mock.MatchedBy(func(arg repository.UpdatePackageStatusWithTrackingParams) bool {
					return arg.ID == expectedUUID &&
						arg.Status == "enviado" &&
						arg.TrackingCode.Valid &&
						len(arg.TrackingCode.String) > 0
				})).Return(nil)
			},
		},
		{
			name:          "Update package with invalid UUID",
			packageID:     "invalid-uuid",
			status:        "enviado",
			setupMocked:   func(repo *repository.QuerierMocked) {},
			expectedError: "parse package id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMocked := repository.NewQuerierMocked(t)
			tt.setupMocked(repoMocked)

			logger := zap.NewNop().Sugar()
			packageService := service.NewPackageService(repoMocked, config.Config{}, logger)

			err := packageService.UpdateStatus(context.Background(), tt.packageID, tt.status)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPackageService_Delete(t *testing.T) {
	tests := []struct {
		name          string
		packageID     string
		setupMocked   func(repo *repository.QuerierMocked)
		expectedError string
	}{
		{
			name:      "Delete package successfully",
			packageID: "550e8400-e29b-41d4-a716-446655440000",
			setupMocked: func(repo *repository.QuerierMocked) {
				expectedUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
				repo.On("DeletePackage", mock.Anything, expectedUUID).Return(nil)
			},
		},
		{
			name:          "Delete package with invalid UUID",
			packageID:     "invalid-uuid",
			setupMocked:   func(repo *repository.QuerierMocked) {},
			expectedError: "parse package id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMocked := repository.NewQuerierMocked(t)
			tt.setupMocked(repoMocked)

			logger := zap.NewNop().Sugar()
			packageService := service.NewPackageService(repoMocked, config.Config{}, logger)

			err := packageService.Delete(context.Background(), tt.packageID)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPackageService_GetQuotes(t *testing.T) {
	tests := []struct {
		name          string
		stateCode     string
		weightKg      float64
		setupMocked   func(repo *repository.QuerierMocked)
		expectedCount int
		expectedError string
	}{
		{
			name:      "Get quotes successfully",
			stateCode: "SP",
			weightKg:  2.5,
			setupMocked: func(repo *repository.QuerierMocked) {
				mockState := repository.GetStateByCodeRow{
					Code:       "SP",
					Name:       "São Paulo",
					RegionName: "Sudeste",
				}
				repo.On("GetStateByCode", mock.Anything, "SP").Return(mockState, nil)

				expectedParams := repository.GetQuotesForPackageParams{
					StateCode: "SP",
					WeightKg:  "2.50",
				}

				expectedQuotes := []repository.GetQuotesForPackageRow{
					{
						Carier:                "Nebulix Logística",
						EstimatedPrice:        14.75,
						EstimatedDeliveryDays: 4,
					},
					{
						Carier:                "RotaFácil Transportes",
						EstimatedPrice:        10.88,
						EstimatedDeliveryDays: 7,
					},
				}

				repo.On("GetQuotesForPackage", mock.Anything, expectedParams).Return(expectedQuotes, nil)
			},
			expectedCount: 2,
		},
		{
			name:      "Get quotes for invalid state",
			stateCode: "XX",
			weightKg:  2.5,
			setupMocked: func(repo *repository.QuerierMocked) {
				repo.On("GetStateByCode", mock.Anything, "XX").Return(repository.GetStateByCodeRow{}, sql.ErrNoRows)
			},
			expectedError: "invalid state code",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMocked := repository.NewQuerierMocked(t)
			tt.setupMocked(repoMocked)

			logger := zap.NewNop().Sugar()
			packageService := service.NewPackageService(repoMocked, config.Config{}, logger)

			result, err := packageService.GetQuotes(context.Background(), tt.stateCode, tt.weightKg)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedCount, len(result))

				if len(result) > 0 {
					assert.NotEmpty(t, result[0].Carier)
					assert.Greater(t, result[0].EstimatedPrice, 0.0)
					assert.Greater(t, result[0].EstimatedDeliveryDays, int32(0))
				}
			}
		})
	}
}

func TestPackageService_HireCarrier(t *testing.T) {
	tests := []struct {
		name          string
		packageID     string
		carrierID     string
		price         string
		deliveryDays  int32
		setupMocked   func(repo *repository.QuerierMocked)
		expectedError string
	}{
		{
			name:         "Hire carrier successfully",
			packageID:    "550e8400-e29b-41d4-a716-446655440000",
			carrierID:    "660e8400-e29b-41d4-a716-446655440001",
			price:        "25.90",
			deliveryDays: 5,
			setupMocked: func(repo *repository.QuerierMocked) {
				expectedPkgUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
				expectedCarrierUUID := uuid.MustParse("660e8400-e29b-41d4-a716-446655440001")
				regionUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440002")

				mockPackage := repository.Package{
					ID:               expectedPkgUUID,
					TrackingCode:     sql.NullString{Valid: false},
					Product:          "Test Product",
					WeightKg:         2.5,
					DestinationState: "SP",
					Status:           "criado",
				}
				repo.On("GetPackageById", mock.Anything, expectedPkgUUID).Return(mockPackage, nil)

				mockRegion := repository.GetRegionByStateRow{
					ID:   regionUUID,
					Name: "Sudeste",
				}
				repo.On("GetRegionByState", mock.Anything, "SP").Return(mockRegion, nil)

				mockCarrierRegions := []repository.GetCarrierRegionsRow{
					{
						ID:                    uuid.New(),
						CarrierID:             expectedCarrierUUID,
						RegionID:              regionUUID,
						EstimatedDeliveryDays: 4,
						PricePerKg:            "5.90",
						RegionName:            "Sudeste",
					},
				}
				repo.On("GetCarrierRegions", mock.Anything, expectedCarrierUUID).Return(mockCarrierRegions, nil)

				expectedParams := repository.HireCarrierParams{
					ID: expectedPkgUUID,
					HiredCarrierID: uuid.NullUUID{
						UUID:  expectedCarrierUUID,
						Valid: true,
					},
					HiredPrice: sql.NullString{
						String: "25.90",
						Valid:  true,
					},
					HiredDeliveryDays: sql.NullInt32{
						Int32: 5,
						Valid: true,
					},
				}
				repo.On("HireCarrier", mock.Anything, expectedParams).Return(nil)
			},
		},
		{
			name:          "Hire carrier with invalid package UUID",
			packageID:     "invalid-uuid",
			carrierID:     "660e8400-e29b-41d4-a716-446655440001",
			price:         "25.90",
			deliveryDays:  5,
			setupMocked:   func(repo *repository.QuerierMocked) {},
			expectedError: "package not found",
		},
		{
			name:         "Hire carrier with invalid carrier UUID",
			packageID:    "550e8400-e29b-41d4-a716-446655440000",
			carrierID:    "invalid-uuid",
			price:        "25.90",
			deliveryDays: 5,
			setupMocked: func(repo *repository.QuerierMocked) {
				expectedPkgUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")

				mockPackage := repository.Package{
					ID:               expectedPkgUUID,
					TrackingCode:     sql.NullString{Valid: false},
					Product:          "Test Product",
					WeightKg:         2.5,
					DestinationState: "SP",
					Status:           "criado",
				}
				repo.On("GetPackageById", mock.Anything, expectedPkgUUID).Return(mockPackage, nil)
			},
			expectedError: "invalid carrier ID",
		},
		{
			name:         "Hire carrier that does not serve region",
			packageID:    "550e8400-e29b-41d4-a716-446655440000",
			carrierID:    "660e8400-e29b-41d4-a716-446655440001",
			price:        "25.90",
			deliveryDays: 5,
			setupMocked: func(repo *repository.QuerierMocked) {
				expectedPkgUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
				expectedCarrierUUID := uuid.MustParse("660e8400-e29b-41d4-a716-446655440001")
				regionUUID := uuid.MustParse("550e8400-e29b-41d4-a716-446655440005")

				mockPackage := repository.Package{
					ID:               expectedPkgUUID,
					TrackingCode:     sql.NullString{Valid: false},
					Product:          "Test Product",
					WeightKg:         2.5,
					DestinationState: "AM",
					Status:           "criado",
				}
				repo.On("GetPackageById", mock.Anything, expectedPkgUUID).Return(mockPackage, nil)

				mockRegion := repository.GetRegionByStateRow{
					ID:   regionUUID,
					Name: "Norte",
				}
				repo.On("GetRegionByState", mock.Anything, "AM").Return(mockRegion, nil)

				mockCarrierRegions := []repository.GetCarrierRegionsRow{}
				repo.On("GetCarrierRegions", mock.Anything, expectedCarrierUUID).Return(mockCarrierRegions, nil)
			},
			expectedError: "carrier does not serve this region",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMocked := repository.NewQuerierMocked(t)
			tt.setupMocked(repoMocked)

			logger := zap.NewNop().Sugar()
			packageService := service.NewPackageService(repoMocked, config.Config{}, logger)

			err := packageService.HireCarrier(context.Background(), tt.packageID, tt.carrierID, tt.price, tt.deliveryDays)

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestPackageService_GetCarriers(t *testing.T) {
	tests := []struct {
		name          string
		setupMocked   func(repo *repository.QuerierMocked)
		expectedCount int
		expectedError string
	}{
		{
			name: "Get carriers successfully",
			setupMocked: func(repo *repository.QuerierMocked) {
				expectedCarriers := []repository.Carrier{
					{
						ID:   uuid.MustParse("660e8400-e29b-41d4-a716-446655440001"),
						Name: "Nebulix Logística",
					},
					{
						ID:   uuid.MustParse("660e8400-e29b-41d4-a716-446655440002"),
						Name: "RotaFácil Transportes",
					},
				}

				repo.On("ListCarriers", mock.Anything).Return(expectedCarriers, nil)
			},
			expectedCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMocked := repository.NewQuerierMocked(t)
			tt.setupMocked(repoMocked)

			logger := zap.NewNop().Sugar()
			packageService := service.NewPackageService(repoMocked, config.Config{}, logger)

			result, err := packageService.GetCarriers(context.Background())

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				require.Equal(t, tt.expectedCount, len(result))

				if len(result) > 0 {
					assert.NotEmpty(t, result[0].Name)
					assert.NotEmpty(t, result[0].ID)
				}
			}
		})
	}
}

func TestPackageService_GetStates(t *testing.T) {
	tests := []struct {
		name          string
		setupMocked   func(repo *repository.QuerierMocked)
		expectedCount int
		expectedError string
	}{
		{
			name: "Get states successfully",
			setupMocked: func(repo *repository.QuerierMocked) {
				expectedStates := []repository.ListStatesRow{
					{
						Code:       "SP",
						Name:       "São Paulo",
						RegionName: "Sudeste",
					},
					{
						Code:       "RJ",
						Name:       "Rio de Janeiro",
						RegionName: "Sudeste",
					},
				}

				repo.On("ListStates", mock.Anything).Return(expectedStates, nil)
			},
			expectedCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMocked := repository.NewQuerierMocked(t)
			tt.setupMocked(repoMocked)

			logger := zap.NewNop().Sugar()
			packageService := service.NewPackageService(repoMocked, config.Config{}, logger)

			result, err := packageService.GetStates(context.Background())

			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
				require.Equal(t, tt.expectedCount, len(result))

				if len(result) > 0 {
					assert.NotEmpty(t, result[0].Code)
					assert.NotEmpty(t, result[0].Name)
					assert.NotEmpty(t, result[0].RegionName)
				}
			}
		})
	}
}
