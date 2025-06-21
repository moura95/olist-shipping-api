package service

import (
	"context"
	"fmt"

	"github/moura95/olist-shipping-api/internal/repository"
)

func (s *PackageService) GetStates(ctx context.Context) ([]repository.ListStatesRow, error) {
	states, err := s.repository.ListStates(ctx)
	if err != nil {
		return nil, fmt.Errorf("list states: %v", err)
	}

	return states, nil
}
