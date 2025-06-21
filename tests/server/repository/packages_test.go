package repository_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github/moura95/olist-shipping-api/internal/repository"
)

func TestCreatePackage(t *testing.T) {
	defer cleanupTestData(t)

	ctx := context.Background()

	arg := repository.CreatePackageParams{
		TrackingCode:     "BR12345678",
		Product:          "Test Product",
		WeightKg:         2.5,
		DestinationState: "SP",
	}

	pkg, err := testQueries.CreatePackage(ctx, arg)

	require.NoError(t, err)
	assert.NotEmpty(t, pkg.ID)
	assert.Equal(t, arg.TrackingCode, pkg.TrackingCode)
	assert.Equal(t, arg.Product, pkg.Product)
	assert.Equal(t, arg.WeightKg, pkg.WeightKg)
	assert.Equal(t, arg.DestinationState, pkg.DestinationState)
	assert.Equal(t, "criado", pkg.Status)
	assert.NotEmpty(t, pkg.CreatedAt)
}

func TestGetPackageById(t *testing.T) {
	defer cleanupTestData(t)

	ctx := context.Background()

	createArg := repository.CreatePackageParams{
		TrackingCode:     "BR87654321",
		Product:          "Another Product",
		WeightKg:         1.2,
		DestinationState: "RJ",
	}

	createdPkg, err := testQueries.CreatePackage(ctx, createArg)
	require.NoError(t, err)

	pkg, err := testQueries.GetPackageById(ctx, createdPkg.ID)

	require.NoError(t, err)
	assert.Equal(t, createdPkg.ID, pkg.ID)
	assert.Equal(t, createdPkg.TrackingCode, pkg.TrackingCode)
	assert.Equal(t, createdPkg.Product, pkg.Product)
	assert.Equal(t, createdPkg.WeightKg, pkg.WeightKg)
	assert.Equal(t, createdPkg.DestinationState, pkg.DestinationState)
}

func TestGetPackageByTrackingCode(t *testing.T) {
	defer cleanupTestData(t)

	ctx := context.Background()

	trackingCode := "BR11223344"
	createArg := repository.CreatePackageParams{
		TrackingCode:     trackingCode,
		Product:          "Tracked Product",
		WeightKg:         3.0,
		DestinationState: "PR",
	}

	createdPkg, err := testQueries.CreatePackage(ctx, createArg)
	require.NoError(t, err)

	pkg, err := testQueries.GetPackageByTrackingCode(ctx, trackingCode)

	require.NoError(t, err)
	assert.Equal(t, createdPkg.ID, pkg.ID)
	assert.Equal(t, trackingCode, pkg.TrackingCode)
}

func TestListPackages(t *testing.T) {
	defer cleanupTestData(t)

	ctx := context.Background()

	packages := []repository.CreatePackageParams{
		{
			TrackingCode:     "BR11111111",
			Product:          "Product 1",
			WeightKg:         1.0,
			DestinationState: "SP",
		},
		{
			TrackingCode:     "BR22222222",
			Product:          "Product 2",
			WeightKg:         2.0,
			DestinationState: "RJ",
		},
	}

	for _, pkg := range packages {
		_, err := testQueries.CreatePackage(ctx, pkg)
		require.NoError(t, err)
	}

	result, err := testQueries.ListPackages(ctx)

	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(result), 2)
}

func TestUpdatePackageStatus(t *testing.T) {
	defer cleanupTestData(t)

	ctx := context.Background()

	createArg := repository.CreatePackageParams{
		TrackingCode:     "BR99999999",
		Product:          "Status Test Product",
		WeightKg:         1.5,
		DestinationState: "SP",
	}

	createdPkg, err := testQueries.CreatePackage(ctx, createArg)
	require.NoError(t, err)

	updateArg := repository.UpdatePackageStatusParams{
		ID:     createdPkg.ID,
		Status: "enviado",
	}

	err = testQueries.UpdatePackageStatus(ctx, updateArg)
	require.NoError(t, err)

	updatedPkg, err := testQueries.GetPackageById(ctx, createdPkg.ID)
	require.NoError(t, err)
	assert.Equal(t, "enviado", updatedPkg.Status)
}

func TestHireCarrier(t *testing.T) {
	defer cleanupTestData(t)

	ctx := context.Background()

	createArg := repository.CreatePackageParams{
		TrackingCode:     "BR55555555",
		Product:          "Hire Test Product",
		WeightKg:         2.0,
		DestinationState: "SP",
	}

	createdPkg, err := testQueries.CreatePackage(ctx, createArg)
	require.NoError(t, err)

	carrierID := uuid.MustParse("660e8400-e29b-41d4-a716-446655440001")
	hireArg := repository.HireCarrierParams{
		ID: createdPkg.ID,
		HiredCarrierID: uuid.NullUUID{
			UUID:  carrierID,
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

	err = testQueries.HireCarrier(ctx, hireArg)
	require.NoError(t, err)

	updatedPkg, err := testQueries.GetPackageById(ctx, createdPkg.ID)
	require.NoError(t, err)
	assert.Equal(t, "esperando_coleta", updatedPkg.Status)
	assert.True(t, updatedPkg.HiredCarrierID.Valid)
	assert.Equal(t, carrierID, updatedPkg.HiredCarrierID.UUID)
	assert.True(t, updatedPkg.HiredPrice.Valid)
	assert.Equal(t, "25.90", updatedPkg.HiredPrice.String)
	assert.True(t, updatedPkg.HiredDeliveryDays.Valid)
	assert.Equal(t, int32(5), updatedPkg.HiredDeliveryDays.Int32)
}

func TestDeletePackage(t *testing.T) {
	defer cleanupTestData(t)

	ctx := context.Background()

	createArg := repository.CreatePackageParams{
		TrackingCode:     "BR77777777",
		Product:          "Delete Test Product",
		WeightKg:         1.0,
		DestinationState: "RJ",
	}

	createdPkg, err := testQueries.CreatePackage(ctx, createArg)
	require.NoError(t, err)

	err = testQueries.DeletePackage(ctx, createdPkg.ID)
	require.NoError(t, err)

	_, err = testQueries.GetPackageById(ctx, createdPkg.ID)
	assert.Error(t, err)
	assert.Equal(t, sql.ErrNoRows, err)
}

func TestTrackingCodeExists(t *testing.T) {
	defer cleanupTestData(t)

	ctx := context.Background()

	trackingCode := "BR88888888"
	exists, err := testQueries.TrackingCodeExists(ctx, trackingCode)
	require.NoError(t, err)
	assert.False(t, exists)

	createArg := repository.CreatePackageParams{
		TrackingCode:     trackingCode,
		Product:          "Exists Test Product",
		WeightKg:         1.0,
		DestinationState: "SP",
	}

	_, err = testQueries.CreatePackage(ctx, createArg)
	require.NoError(t, err)

	exists, err = testQueries.TrackingCodeExists(ctx, trackingCode)
	require.NoError(t, err)
	assert.True(t, exists)
}

func TestGetQuotesForPackage(t *testing.T) {
	defer cleanupTestData(t)

	ctx := context.Background()

	arg := repository.GetQuotesForPackageParams{
		StateCode: "SP",
		WeightKg:  "2.0",
	}

	quotes, err := testQueries.GetQuotesForPackage(ctx, arg)

	require.NoError(t, err)
	assert.Greater(t, len(quotes), 0)

	for _, quote := range quotes {
		assert.NotEmpty(t, quote.Carier)
		assert.Greater(t, quote.EstimatedPrice, 0.0)
		assert.Greater(t, quote.EstimatedDeliveryDays, int32(0))
	}
}

func TestListCarriers(t *testing.T) {
	defer cleanupTestData(t)

	ctx := context.Background()

	carriers, err := testQueries.ListCarriers(ctx)

	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(carriers), 3)

	for _, carrier := range carriers {
		assert.NotEmpty(t, carrier.ID)
		assert.NotEmpty(t, carrier.Name)
	}
}

func TestGetCarrierById(t *testing.T) {
	defer cleanupTestData(t)

	ctx := context.Background()

	carrierID := uuid.MustParse("660e8400-e29b-41d4-a716-446655440001")

	carrier, err := testQueries.GetCarrierById(ctx, carrierID)

	require.NoError(t, err)
	assert.Equal(t, carrierID, carrier.ID)
	assert.Equal(t, "Nebulix Logística", carrier.Name)
}

func TestGetCarrierRegions(t *testing.T) {
	defer cleanupTestData(t)

	ctx := context.Background()

	carrierID := uuid.MustParse("660e8400-e29b-41d4-a716-446655440001")

	regions, err := testQueries.GetCarrierRegions(ctx, carrierID)

	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(regions), 2)

	for _, region := range regions {
		assert.Equal(t, carrierID, region.CarrierID)
		assert.NotEmpty(t, region.RegionName)
		assert.Greater(t, region.EstimatedDeliveryDays, int32(0))
	}
}

func TestListStates(t *testing.T) {
	defer cleanupTestData(t)

	ctx := context.Background()

	states, err := testQueries.ListStates(ctx)

	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(states), 5)

	for _, state := range states {
		assert.NotEmpty(t, state.Code)
		assert.NotEmpty(t, state.Name)
		assert.NotEmpty(t, state.RegionName)
	}
}

func TestGetStateByCode(t *testing.T) {
	defer cleanupTestData(t)

	ctx := context.Background()

	state, err := testQueries.GetStateByCode(ctx, "SP")

	require.NoError(t, err)
	assert.Equal(t, "SP", state.Code)
	assert.Equal(t, "São Paulo", state.Name)
	assert.Equal(t, "Sudeste", state.RegionName)
}

func TestGetRegionByState(t *testing.T) {
	defer cleanupTestData(t)

	ctx := context.Background()

	region, err := testQueries.GetRegionByState(ctx, "SP")

	require.NoError(t, err)
	assert.Equal(t, uuid.MustParse("550e8400-e29b-41d4-a716-446655440002"), region.ID)
	assert.Equal(t, "Sudeste", region.Name)
}

func TestListRegions(t *testing.T) {
	defer cleanupTestData(t)

	ctx := context.Background()

	regions, err := testQueries.ListRegions(ctx)

	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(regions), 5)

	for _, region := range regions {
		assert.NotEmpty(t, region.ID)
		assert.NotEmpty(t, region.Name)
	}
}
