package services

import "example/apies/models"

type ColorService interface {
	CreateColor(*models.Color) error
	// GetColor(*string) (*models.Color, error)
	GetAll() ([]*models.Color, error)
}
