package services

import (
	"multiaura/internal/repositories"
)

type ConservationService interface {
	CreateConservation(targetUserID, userID string) error
}

type conservationService struct {
	repo repositories.ConservationRepository
}

func NewConservationService(repo *repositories.ConservationRepository) ConservationService {
	return &conservationService{*repo}
}

// CreateConservation implements ConservationService.
func (c *conservationService) CreateConservation(targetUserID string, userID string) error {
	panic("unimplemented")
}
