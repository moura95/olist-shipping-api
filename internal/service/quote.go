package service

import (
	"context"
	"fmt"

	"github/moura95/olist-shipping-api/internal/repository"
)

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
