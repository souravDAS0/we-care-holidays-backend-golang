package dto

import "errors"

type OrgStatusUpdateDto struct {
	Status string `json:"status" binding:"required" validate:"required" example:"Approved"`
}
// Validate performs validation on OrgStatusUpdateDto
func (dto *OrgStatusUpdateDto) Validate() error {
	if dto.Status == "" {
		return errors.New("status is required")
	}
	
	validStatuses := map[string]bool{
		"Pending":   true,
		"Approved":  true,
		"Suspended": true,
		"Archived":  true,
	}

	if !validStatuses[dto.Status] {
		return ErrInvalidStatus
	}

	return nil
}
