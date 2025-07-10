package constants

import "net/http"

// Health Check
const (
	HealthCheckRoute = "/health"
)

const (
	HTTPOk                  = http.StatusOK                  // 200
	HTTPCreated             = http.StatusCreated             // 201
	HTTPAccepted            = http.StatusAccepted            // 202
	HTTPNoContent           = http.StatusNoContent           // 204
	HTTPBadRequest          = http.StatusBadRequest          // 400
	HTTPUnauthorized        = http.StatusUnauthorized        // 401
	HTTPForbidden           = http.StatusForbidden           // 403
	HTTPNotFound            = http.StatusNotFound            // 404
	HTTPConflict            = http.StatusConflict            // 409
	HTTPUnprocessableEntity = http.StatusUnprocessableEntity // 422
	HTTPInternalServerError = http.StatusInternalServerError // 500
)

const (
	PermissionBasePath         = "/permissions"
	ListPermissionsPath        = ""
	CreatePermissionPath       = ""
	BulkDeletePermissionsPath  = "/bulk-delete"
	BulkRestorePermissionsPath = "/bulk-restore"
	GetPermissionPath          = "/:id"
	UpdatePermissionPath       = "/:id"
	DeletePermissionPath       = "/:id"
	RestorePermissionPath      = "/:id/restore"
	HardDeletePermissionPath   = "/:id/hard-delete"
)

const (
	RoleBasePath         = "/roles"
	ListRolesPath        = ""
	CreateRolePath       = ""
	BulkDeleteRolesPath  = "/bulk-delete"
	BulkRestoreRolesPath = "/bulk-restore"
	GetRolePath          = "/:id"
	UpdateRolePath       = "/:id"
	DeleteRolePath       = "/:id"
	RestoreRolePath      = "/:id/restore"
	HardDeleteRolePath   = "/:id/hard-delete"
)

const (
	UserBasePath         = "/users"
	ListUsersPath        = ""
	CreateUserPath       = "/invite"
	BulkDeleteUsersPath  = "/bulk-delete"
	BulkRestoreUsersPath = "/bulk-restore"
	GetUserPath          = "/:id"
	UpdateUserPath       = "/:id"
	DeleteUserPath       = "/:id"
	UpdateUserStatusPath = "/:id/status"
	UploadUserAvatarPath = "/:id/profile-photo"
	RestoreUserPath      = "/:id/restore"
	HardDeleteUserPath   = "/:id/hard-delete"
)

const (
	OrganizationBasePath   = "/organizations"
	ListOrganizationsPath  = ""
	CreateOrganizationPath = ""

	BulkDeleteOrganizationsPath  = "/bulk-delete"
	BulkRestoreOrganizationsPath = "/bulk-restore"

	GetOrganizationPath        = "/:id"
	UpdateOrganizationPath     = "/:id"
	DeleteOrganizationPath     = "/:id"
	UpdateStatusPath           = "/:id/status"
	UploadOrgLogoPath          = "/:id/logo"
	RestoreOrganizationPath    = "/:id/restore"
	HardDeleteOrganizationPath = "/:id/hard-delete"
)

const (
	LocationBasePath   = "/locations"
	ListLocationsPath  = ""
	CreateLocationPath = ""

	BulkDeleteLocationsPath  = "/bulk-delete"
	BulkRestoreLocationsPath = "/bulk-restore"

	GetLocationPath         = "/:id"
	UpdateLocationPath      = "/:id"
	DeleteLocationPath      = "/:id"
	UploadLocationMediaPath = "/:id/media"
	RestoreLocationPath     = "/:id/restore"
	HardDeleteLocationPath  = "/:id/hard-delete"
)
