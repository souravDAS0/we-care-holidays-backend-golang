package models

// BulkDeleteResponse represents the response for bulk delete operations
type BulkDeleteResponse struct {
	Message      string   `json:"message" example:"Roles deletion processed"`
	RequestedIDs []string `json:"requestedIds" example:"6835bf49c62fee1db6585e9f, 683467a32bf5a05aefe43cb, 68344ada06017a47db237f66"`
	DeletedCount int      `json:"deletedCount" example:"1"`
	DeletedIDs   []string `json:"deletedIds" example:"6835bf49c62fee1db6585e9f"`
	InvalidIDs   []string `json:"invalidIds" example:"683467a32bf5a05aefe43cb"`
	NotFoundIDs  []string `json:"notFoundIds" example:"68344ada06017a47db237f66"`
}

// BulkRestoreResponse represents the response for bulk restore operations
type BulkRestoreResponse struct {
	Message       string   `json:"message" example:"Roles restoration processed"`
	RequestedIDs  []string `json:"requestedIds" example:"6835bf49c62fee1db6585e9f, 683467a32bf5a05aefe43cb, 68344ada06017a47db237f66"`
	RestoredCount int      `json:"restoredCount" example:"1"`
	RestoredIDs   []string `json:"restoredIds" example:"6835bf49c62fee1db6585e9f"`
	InvalidIDs    []string `json:"invalidIds" example:"683467a32bf5a05aefe43cb"`
	NotFoundIDs   []string `json:"notFoundIds" example:"68344ada06017a47db237f66"`
}