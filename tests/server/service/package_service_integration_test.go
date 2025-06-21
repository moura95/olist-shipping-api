package service

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"github/moura95/olist-shipping-api/config"
	"github/moura95/olist-shipping-api/internal/repository"
	"github/moura95/olist-shipping-api/internal/service"
	"go.uber.org/zap"
)

var (
	integrationTestDB        *sql.DB
	integrationTestContainer *postgres.PostgresContainer
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	var err error
	integrationTestContainer, err = postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15-alpine"),
		postgres.WithDatabase("integration_test_db"),
		postgres.WithUsername("test_user"),
		postgres.WithPassword("test_password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second)),
	)
	if err != nil {
		log.Fatalf("Failed to start PostgreSQL container: %v", err)
	}

	connStr, err := integrationTestContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatalf("Failed to get connection string: %v", err)
	}

	integrationTestDB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := integrationTestDB.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	if err := runIntegrationMigrations(integrationTestDB); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	code := m.Run()

	if integrationTestContainer != nil {
		integrationTestContainer.Terminate(ctx)
	}

	if integrationTestDB != nil {
		integrationTestDB.Close()
	}

	os.Exit(code)
}

func runIntegrationMigrations(db *sql.DB) error {
	migrationsPath := "../../../db/migrations"

	files, err := ioutil.ReadDir(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %v", err)
	}

	var upFiles []string
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sql" && len(file.Name()) > 6 && file.Name()[len(file.Name())-6:] == "up.sql" {
			upFiles = append(upFiles, file.Name())
		}
	}

	for _, fileName := range upFiles {
		content, err := ioutil.ReadFile(filepath.Join(migrationsPath, fileName))
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %v", fileName, err)
		}

		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %v", fileName, err)
		}
	}

	return nil
}

func cleanupIntegrationTestData(t *testing.T) {
	ctx := context.Background()

	tables := []string{
		"packages",
	}

	for _, table := range tables {
		_, err := integrationTestDB.ExecContext(ctx, fmt.Sprintf("DELETE FROM %s WHERE tracking_code LIKE 'TEST%%' OR tracking_code LIKE 'BR%%'", table))
		if err != nil {
			t.Logf("Failed to cleanup table %s: %v", table, err)
		}
	}
}

func TestPackageServiceIntegration_CreateAndGet(t *testing.T) {

	defer cleanupIntegrationTestData(t)

	ctx := context.Background()
	logger := zap.NewNop().Sugar()
	store := repository.New(integrationTestDB)
	service := service.NewPackageService(store, config.Config{}, logger)

	createdPkg, err := service.Create(ctx, "Integration Test Product", 1.5, "SP")
	require.NoError(t, err)
	require.NotNil(t, createdPkg)

	assert.NotEmpty(t, createdPkg.ID)
	assert.NotEmpty(t, createdPkg.TrackingCode)
	assert.Equal(t, "Integration Test Product", createdPkg.Product)
	assert.Equal(t, 1.5, createdPkg.WeightKg)
	assert.Equal(t, "SP", createdPkg.DestinationState)
	assert.Equal(t, "criado", createdPkg.Status)

	retrievedPkg, err := service.GetByID(ctx, createdPkg.ID.String())
	require.NoError(t, err)
	require.NotNil(t, retrievedPkg)

	assert.Equal(t, createdPkg.ID, retrievedPkg.ID)
	assert.Equal(t, createdPkg.TrackingCode, retrievedPkg.TrackingCode)
	assert.Equal(t, createdPkg.Product, retrievedPkg.Product)
	assert.Equal(t, createdPkg.WeightKg, retrievedPkg.WeightKg)
	assert.Equal(t, createdPkg.DestinationState, retrievedPkg.DestinationState)
	assert.Equal(t, createdPkg.Status, retrievedPkg.Status)

	retrievedByTracking, err := service.GetByTrackingCode(ctx, createdPkg.TrackingCode)
	require.NoError(t, err)
	require.NotNil(t, retrievedByTracking)

	assert.Equal(t, createdPkg.ID, retrievedByTracking.ID)
	assert.Equal(t, createdPkg.TrackingCode, retrievedByTracking.TrackingCode)
}

func TestPackageServiceIntegration_GetAll(t *testing.T) {
	defer cleanupIntegrationTestData(t)

	ctx := context.Background()
	logger := zap.NewNop().Sugar()
	store := repository.New(integrationTestDB)
	service := service.NewPackageService(store, config.Config{}, logger)

	initialPackages, err := service.GetAll(ctx)
	require.NoError(t, err)
	initialCount := len(initialPackages)

	pkg1, err := service.Create(ctx, "Test Product 1", 1.0, "SP")
	require.NoError(t, err)

	pkg2, err := service.Create(ctx, "Test Product 2", 2.0, "RJ")
	require.NoError(t, err)

	allPackages, err := service.GetAll(ctx)
	require.NoError(t, err)

	assert.GreaterOrEqual(t, len(allPackages), initialCount+2)

	var foundPkg1, foundPkg2 bool
	for _, pkg := range allPackages {
		if pkg.ID == pkg1.ID {
			foundPkg1 = true
			assert.Equal(t, "Test Product 1", pkg.Product)
		}
		if pkg.ID == pkg2.ID {
			foundPkg2 = true
			assert.Equal(t, "Test Product 2", pkg.Product)
		}
	}

	assert.True(t, foundPkg1, "Package 1 should be found in list")
	assert.True(t, foundPkg2, "Package 2 should be found in list")
}

func TestPackageServiceIntegration_UpdateStatus(t *testing.T) {

	defer cleanupIntegrationTestData(t)

	ctx := context.Background()
	logger := zap.NewNop().Sugar()
	store := repository.New(integrationTestDB)
	service := service.NewPackageService(store, config.Config{}, logger)

	createdPkg, err := service.Create(ctx, "Status Test Product", 1.0, "SP")
	require.NoError(t, err)

	assert.Equal(t, "criado", createdPkg.Status)

	err = service.UpdateStatus(ctx, createdPkg.ID.String(), "enviado")
	require.NoError(t, err)

	updatedPkg, err := service.GetByID(ctx, createdPkg.ID.String())
	require.NoError(t, err)

	assert.Equal(t, "enviado", updatedPkg.Status)

	err = service.UpdateStatus(ctx, createdPkg.ID.String(), "entregue")
	require.NoError(t, err)

	finalPkg, err := service.GetByID(ctx, createdPkg.ID.String())
	require.NoError(t, err)

	assert.Equal(t, "entregue", finalPkg.Status)
}

func TestPackageServiceIntegration_Delete(t *testing.T) {

	defer cleanupIntegrationTestData(t)

	ctx := context.Background()
	logger := zap.NewNop().Sugar()
	store := repository.New(integrationTestDB)
	service := service.NewPackageService(store, config.Config{}, logger)

	createdPkg, err := service.Create(ctx, "Delete Test Product", 1.0, "SP")
	require.NoError(t, err)

	retrievedPkg, err := service.GetByID(ctx, createdPkg.ID.String())
	require.NoError(t, err)
	assert.NotNil(t, retrievedPkg)

	err = service.Delete(ctx, createdPkg.ID.String())
	require.NoError(t, err)

	deletedPkg, err := service.GetByID(ctx, createdPkg.ID.String())
	assert.Error(t, err)
	assert.Nil(t, deletedPkg)
	assert.Contains(t, err.Error(), "get package by id")
}

func TestPackageServiceIntegration_GetQuotes(t *testing.T) {

	ctx := context.Background()
	logger := zap.NewNop().Sugar()
	store := repository.New(integrationTestDB)
	service := service.NewPackageService(store, config.Config{}, logger)

	quotes, err := service.GetQuotes(ctx, "SP", 2.0)
	require.NoError(t, err)
	require.Greater(t, len(quotes), 0)

	for _, quote := range quotes {
		assert.NotEmpty(t, quote.Carier)
		assert.Greater(t, quote.EstimatedPrice, 0.0)
		assert.Greater(t, quote.EstimatedDeliveryDays, int32(0))
	}

	quotesRJ, err := service.GetQuotes(ctx, "RJ", 1.5)
	require.NoError(t, err)
	require.Greater(t, len(quotesRJ), 0)

	quotesNorth, err := service.GetQuotes(ctx, "AM", 3.0)
	require.NoError(t, err)
	assert.True(t, len(quotesNorth) >= 0)
}

func TestPackageServiceIntegration_HireCarrier(t *testing.T) {

	defer cleanupIntegrationTestData(t)

	ctx := context.Background()
	logger := zap.NewNop().Sugar()
	store := repository.New(integrationTestDB)
	service := service.NewPackageService(store, config.Config{}, logger)

	createdPkg, err := service.Create(ctx, "Hire Carrier Test", 2.0, "SP")
	require.NoError(t, err)

	assert.Equal(t, "criado", createdPkg.Status)
	assert.False(t, createdPkg.HiredCarrierID.Valid)

	carrierID := "660e8400-e29b-41d4-a716-446655440001"
	price := "25.90"
	deliveryDays := int32(5)

	err = service.HireCarrier(ctx, createdPkg.ID.String(), carrierID, price, deliveryDays)
	require.NoError(t, err)

	updatedPkg, err := service.GetByID(ctx, createdPkg.ID.String())
	require.NoError(t, err)

	assert.Equal(t, "esperando_coleta", updatedPkg.Status)
	assert.True(t, updatedPkg.HiredCarrierID.Valid)
	assert.Equal(t, carrierID, updatedPkg.HiredCarrierID.UUID.String())
	assert.True(t, updatedPkg.HiredPrice.Valid)
	assert.Equal(t, price, updatedPkg.HiredPrice.String)
	assert.True(t, updatedPkg.HiredDeliveryDays.Valid)
	assert.Equal(t, deliveryDays, updatedPkg.HiredDeliveryDays.Int32)
}

func TestPackageServiceIntegration_GetCarriers(t *testing.T) {

	ctx := context.Background()
	logger := zap.NewNop().Sugar()
	store := repository.New(integrationTestDB)
	service := service.NewPackageService(store, config.Config{}, logger)

	carriers, err := service.GetCarriers(ctx)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(carriers), 3)

	expectedCarriers := map[string]bool{
		"Nebulix Logística":     false,
		"RotaFácil Transportes": false,
		"Moventra Express":      false,
	}

	for _, carrier := range carriers {
		assert.NotEmpty(t, carrier.ID)
		assert.NotEmpty(t, carrier.Name)

		if _, exists := expectedCarriers[carrier.Name]; exists {
			expectedCarriers[carrier.Name] = true
		}
	}

	for name, found := range expectedCarriers {
		assert.True(t, found, "Carrier %s should be found", name)
	}
}

func TestPackageServiceIntegration_GetStates(t *testing.T) {

	ctx := context.Background()
	logger := zap.NewNop().Sugar()
	store := repository.New(integrationTestDB)
	service := service.NewPackageService(store, config.Config{}, logger)

	states, err := service.GetStates(ctx)
	require.NoError(t, err)
	require.GreaterOrEqual(t, len(states), 27)

	expectedStates := map[string]string{
		"SP": "Sudeste",
		"RJ": "Sudeste",
		"MG": "Sudeste",
		"RS": "Sul",
		"PR": "Sul",
		"BA": "Nordeste",
		"AM": "Norte",
		"GO": "Centro-Oeste",
	}

	foundStates := make(map[string]string)
	for _, state := range states {
		assert.NotEmpty(t, state.Code)
		assert.NotEmpty(t, state.Name)
		assert.NotEmpty(t, state.RegionName)

		foundStates[state.Code] = state.RegionName
	}

	for code, expectedRegion := range expectedStates {
		foundRegion, exists := foundStates[code]
		assert.True(t, exists, "State %s should exist", code)
		assert.Equal(t, expectedRegion, foundRegion, "State %s should be in region %s", code, expectedRegion)
	}
}

func TestPackageServiceIntegration_UniqueTrackingCodes(t *testing.T) {

	defer cleanupIntegrationTestData(t)

	ctx := context.Background()
	logger := zap.NewNop().Sugar()
	store := repository.New(integrationTestDB)
	service := service.NewPackageService(store, config.Config{}, logger)

	trackingCodes := make(map[string]bool)

	for i := 0; i < 10; i++ {
		pkg, err := service.Create(ctx, "Unique Test Product", 1.0, "SP")
		require.NoError(t, err)

		assert.False(t, trackingCodes[pkg.TrackingCode], "Tracking code %s should be unique", pkg.TrackingCode)
		trackingCodes[pkg.TrackingCode] = true

		assert.Regexp(t, "^BR[0-9]{8}$", pkg.TrackingCode, "Tracking code should match pattern")
	}
}

func TestPackageServiceIntegration_ErrorCases(t *testing.T) {

	defer cleanupIntegrationTestData(t)

	ctx := context.Background()
	logger := zap.NewNop().Sugar()
	store := repository.New(integrationTestDB)
	service := service.NewPackageService(store, config.Config{}, logger)

	t.Run("Get non-existent package", func(t *testing.T) {
		pkg, err := service.GetByID(ctx, "550e8400-e29b-41d4-a716-446655440999")
		assert.Error(t, err)
		assert.Nil(t, pkg)
		assert.Contains(t, err.Error(), "get package by id")
	})

	t.Run("Get package with invalid UUID", func(t *testing.T) {
		pkg, err := service.GetByID(ctx, "invalid-uuid")
		assert.Error(t, err)
		assert.Nil(t, pkg)
		assert.Contains(t, err.Error(), "parse package id")
	})

	t.Run("Get package by non-existent tracking code", func(t *testing.T) {
		pkg, err := service.GetByTrackingCode(ctx, "BR99999999")
		assert.Error(t, err)
		assert.Nil(t, pkg)
		assert.Contains(t, err.Error(), "get package by tracking code")
	})

	t.Run("Update status with invalid UUID", func(t *testing.T) {
		err := service.UpdateStatus(ctx, "invalid-uuid", "enviado")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "parse package id")
	})

	t.Run("Delete with invalid UUID", func(t *testing.T) {
		err := service.Delete(ctx, "invalid-uuid")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "parse package id")
	})

	t.Run("Hire carrier with invalid package UUID", func(t *testing.T) {
		err := service.HireCarrier(ctx, "invalid-uuid", "660e8400-e29b-41d4-a716-446655440001", "25.90", 5)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "parse package id")
	})

	t.Run("Hire carrier with invalid carrier UUID", func(t *testing.T) {
		pkg, err := service.Create(ctx, "Error Test Product", 1.0, "SP")
		require.NoError(t, err)

		err = service.HireCarrier(ctx, pkg.ID.String(), "invalid-uuid", "25.90", 5)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "parse carrier id")
	})

	t.Run("Get quotes for invalid state", func(t *testing.T) {
		quotes, err := service.GetQuotes(ctx, "XX", 1.0)
		require.NoError(t, err)
		assert.Empty(t, quotes)
	})
}
