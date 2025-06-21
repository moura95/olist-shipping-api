package service

import (
	"context"
	"fmt"

	"github/moura95/olist-shipping-api/internal/repository"
)

func (s *PackageService) GetCarriers(ctx context.Context) ([]repository.Carrier, error) {
	carriers, err := s.repository.ListCarriers(ctx)
	if err != nil {
		return nil, fmt.Errorf("list carriers: %v", err)
	}

	return carriers, nil
}
