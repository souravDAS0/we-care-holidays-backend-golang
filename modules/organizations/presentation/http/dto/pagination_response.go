// Add this to your dto package (dto/pagination_response.go)

package dto

import "bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/organizations/domain/entity"

// PaginatedOrganizationsResponse represents the paginated response for organizations
type PaginatedOrganizationsResponse struct {
	Items          []entity.Organization `json:"items"`
	Page           int                   `json:"page" example:"1"`
	Limit          int                   `json:"limit" example:"10"`
	Total          int64                 `json:"total" example:"2"`
	TotalPages     int64                 `json:"totalPages" example:"1"`
	IncludeDeleted bool                  `json:"includeDeleted" example:"false"`
}