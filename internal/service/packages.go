package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github/moura95/olist-shipping-api/config"
	"github/moura95/olist-shipping-api/internal/repository"
	"github/moura95/olist-shipping-api/pkg/tracking"
	"go.uber.org/zap"
)

type PackageService struct {
	repository repository.Querier
	config     config.Config
	logger     *zap.SugaredLogger
}

func NewPackageService(repo repository.Querier, cfg config.Config, log *zap.SugaredLogger) *PackageService {
	return &PackageService{
		repository: repo,
		config:     cfg,
		logger:     log,
	}
}

func (s *PackageService) Create(ctx context.Context, product string, weightKg float64, destinationState string) (*repository.Package, error) {
	trackingCode := tracking.GenerateUniqueTrackingCode(func(code string) bool {
		exists, err := s.repository.TrackingCodeExists(ctx, code)
		if err != nil {
			s.logger.Errorw("failed to check tracking code existence", "error", err, "tracking_code", code)
			return true
		}
		return exists
	})

	arg := repository.CreatePackageParams{
		TrackingCode:     trackingCode,
		Product:          product,
		WeightKg:         weightKg,
		DestinationState: destinationState,
	}

	pkg, err := s.repository.CreatePackage(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("create package: %v", err)
	}

	return &pkg, nil
}

func (s *PackageService) GetByID(ctx context.Context, id string) (*repository.Package, error) {
	packageID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("parse package id: %v", err)
	}

	pkg, err := s.repository.GetPackageById(ctx, packageID)
	if err != nil {
		return nil, fmt.Errorf("get package by id: %v", err)
	}

	return &pkg, nil
}

func (s *PackageService) GetByTrackingCode(ctx context.Context, trackingCode string) (*repository.Package, error) {
	pkg, err := s.repository.GetPackageByTrackingCode(ctx, trackingCode)
	if err != nil {
		return nil, fmt.Errorf("get package by tracking code: %v", err)
	}

	return &pkg, nil
}

func (s *PackageService) GetAll(ctx context.Context) ([]repository.Package, error) {
	packages, err := s.repository.ListPackages(ctx)
	if err != nil {
		return nil, fmt.Errorf("list packages: %v", err)
	}

	return packages, nil
}

func (s *PackageService) UpdateStatus(ctx context.Context, id, status string) error {
	packageID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("parse package id: %v", err)
	}

	arg := repository.UpdatePackageStatusParams{
		ID:     packageID,
		Status: status,
	}

	err = s.repository.UpdatePackageStatus(ctx, arg)
	if err != nil {
		return fmt.Errorf("update package status: %v", err)
	}

	return nil
}

func (s *PackageService) Delete(ctx context.Context, id string) error {
	packageID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("parse package id: %v", err)
	}

	err = s.repository.DeletePackage(ctx, packageID)
	if err != nil {
		return fmt.Errorf("delete package: %v", err)
	}

	return nil
}

func (s *PackageService) GetQuotes(ctx context.Context, stateCode string, weightKg float64) ([]repository.GetQuotesForPackageRow, error) {
	arg := repository.GetQuotesForPackageParams{
		StateCode: stateCode,
		WeightKg:  fmt.Sprintf("%.2f", weightKg),
	}

	quotes, err := s.repository.GetQuotesForPackage(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("get quotes for package: %v", err)
	}

	return quotes, nil
}

func (s *PackageService) HireCarrier(ctx context.Context, packageID, carrierID string, price string, deliveryDays int32) error {
	pkg, err := s.GetByID(ctx, packageID)
	if err != nil {
		return fmt.Errorf("package not found")
	}

	if err := s.ValidateCarrierForRegion(ctx, carrierID, pkg.DestinationState); err != nil {
		return err
	}

	pkgUUID, err := uuid.Parse(packageID)
	if err != nil {
		return fmt.Errorf("invalid package ID")
	}

	carrierUUID, err := uuid.Parse(carrierID)
	if err != nil {
		return fmt.Errorf("invalid carrier ID")
	}

	arg := repository.HireCarrierParams{
		ID: pkgUUID,
		HiredCarrierID: uuid.NullUUID{
			UUID:  carrierUUID,
			Valid: true,
		},
		HiredPrice: sql.NullString{
			String: price,
			Valid:  true,
		},
		HiredDeliveryDays: sql.NullInt32{
			Int32: deliveryDays,
			Valid: true,
		},
	}

	return s.repository.HireCarrier(ctx, arg)
}

func (s *PackageService) GetCarriers(ctx context.Context) ([]repository.Carrier, error) {
	carriers, err := s.repository.ListCarriers(ctx)
	if err != nil {
		return nil, fmt.Errorf("list carriers: %v", err)
	}

	return carriers, nil
}

func (s *PackageService) GetStates(ctx context.Context) ([]repository.ListStatesRow, error) {
	states, err := s.repository.ListStates(ctx)
	if err != nil {
		return nil, fmt.Errorf("list states: %v", err)
	}

	return states, nil
}

func (s *PackageService) ValidateCarrierForRegion(ctx context.Context, carrierID, stateCode string) error {
	carrierUUID, err := uuid.Parse(carrierID)
	if err != nil {
		return fmt.Errorf("invalid carrier ID")
	}

	region, err := s.repository.GetRegionByState(ctx, stateCode)
	if err != nil {
		return fmt.Errorf("invalid state code")
	}

	carrierRegions, err := s.repository.GetCarrierRegions(ctx, carrierUUID)
	if err != nil {
		return fmt.Errorf("carrier not found")
	}

	for _, cr := range carrierRegions {
		if cr.RegionID == region.ID {
			return nil
		}
	}

	return fmt.Errorf("carrier does not serve this region")
}
