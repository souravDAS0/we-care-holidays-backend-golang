package models

import "errors"

// BulkDeleteDto represents a request to delete multiple organizations
type BulkDeleteDto struct {
	IDs []string `json:"ids" binding:"required"`
}

// Validate performs validation on BulkDeleteDto
func (dto *BulkDeleteDto) Validate() error {
	if len(dto.IDs) == 0 {
		return errors.New("at least one organization ID is required")
	}
	return nil
}