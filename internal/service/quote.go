package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github/moura95/olist-shipping-api/internal/repository"
)

func (s *PackageService) GetQuotes(ctx context.Context, stateCode string, weightKg float64) ([]repository.GetQuotesForPackageRow, error) {
	_, err := s.repository.GetStateByCode(ctx, stateCode)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("invalid state code: %s", stateCode)
		}
		return nil, fmt.Errorf("error validating state: %v", err)
	}

	// Busca as cotações
	arg := repository.GetQuotesForPackageParams{
		StateCode: stateCode,
		WeightKg:  fmt.Sprintf("%.2f", weightKg),
	}

	quotes, err := s.repository.GetQuotesForPackage(ctx, arg)
	if err != nil {
		return nil, fmt.Errorf("error busca cotações: %v", err)
	}

	if len(quotes) == 0 {
		return nil, fmt.Errorf("nehuma transportadora encontrada para o estado %s", stateCode)
	}

	return quotes, nil
}
