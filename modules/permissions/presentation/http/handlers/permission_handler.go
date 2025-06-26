package handlers

import (
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/permissions/domain/usecases"
)

// PermissionHandler handles HTTP requests for Permissions
type PermissionHandler struct {
	GetPermissionUseCase        *usecases.GetPermissionUseCase
	CreatePermissionUseCase     *usecases.CreatePermissionUseCase
	ListPermissionsUseCase      *usecases.ListPermissionsUseCase
	UpdatePermissionUseCase     *usecases.UpdatePermissionUseCase
	HardDeletePermissionUseCase *usecases.HardDeletePermissionUseCase
}

func NewPermissionHandler(GetPermissionUseCase *usecases.GetPermissionUseCase,
	CreatePermissionUseCase *usecases.CreatePermissionUseCase,
	ListPermissionsUseCase *usecases.ListPermissionsUseCase,
	UpdatePermissionUseCase *usecases.UpdatePermissionUseCase,
	HardDeletePermissionUseCase *usecases.HardDeletePermissionUseCase,
) *PermissionHandler {
	return &PermissionHandler{
		GetPermissionUseCase:        GetPermissionUseCase,
		CreatePermissionUseCase:     CreatePermissionUseCase,
		ListPermissionsUseCase:      ListPermissionsUseCase,
		UpdatePermissionUseCase:     UpdatePermissionUseCase,
		HardDeletePermissionUseCase: HardDeletePermissionUseCase,
	}
}
