package handlers

import (
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/middleware"
	permUsecases "bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/permissions/domain/usecases"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/roles/domain/usecases"
)

type RoleHandler struct {
	GetRoleUseCase             *usecases.GetRoleUseCase
	CreateRoleUseCase          *usecases.CreateRoleUseCase
	ListRolesUseCase           *usecases.ListRolesUseCase
	UpdateRoleUseCase          *usecases.UpdateRoleUseCase
	SoftDeleteRoleUseCase      *usecases.SoftDeleteRoleUseCase
	RestoreRoleUseCase         *usecases.RestoreRoleUseCase
	BulkSoftDeleteRolesUseCase *usecases.BulkSoftDeleteRolesUseCase
	HardDeleteRoleUseCase      *usecases.HardDeleteRoleUseCase
	BulkRestoreRolesUseCase    *usecases.BulkRestoreRolesUseCase
	rbacService                middleware.RBACService
	permissionValidator        *middleware.PermissionValidator
	ListPermissionsUseCase     *permUsecases.ListPermissionsUseCase
}

func NewRoleHandler(
	GetRoleUseCase *usecases.GetRoleUseCase,
	CreateRoleUseCase *usecases.CreateRoleUseCase,
	ListRolesUseCase *usecases.ListRolesUseCase,
	UpdateRoleUseCase *usecases.UpdateRoleUseCase,
	SoftDeleteRoleUseCase *usecases.SoftDeleteRoleUseCase,
	RestoreRoleUseCase *usecases.RestoreRoleUseCase,
	BulkSoftDeleteRolesUseCase *usecases.BulkSoftDeleteRolesUseCase,
	HardDeleteRoleUseCase *usecases.HardDeleteRoleUseCase,
	BulkRestoreRolesUseCase *usecases.BulkRestoreRolesUseCase,
	ListPermissionsUseCase *permUsecases.ListPermissionsUseCase,
	// rbacService middleware.RBACService,
	// permissionValidator *middleware.PermissionValidator,

) *RoleHandler {
	return &RoleHandler{
		GetRoleUseCase:             GetRoleUseCase,
		CreateRoleUseCase:          CreateRoleUseCase,
		ListRolesUseCase:           ListRolesUseCase,
		UpdateRoleUseCase:          UpdateRoleUseCase,
		SoftDeleteRoleUseCase:      SoftDeleteRoleUseCase,
		RestoreRoleUseCase:         RestoreRoleUseCase,
		BulkSoftDeleteRolesUseCase: BulkSoftDeleteRolesUseCase,
		HardDeleteRoleUseCase:      HardDeleteRoleUseCase,
		BulkRestoreRolesUseCase:    BulkRestoreRolesUseCase,
		ListPermissionsUseCase:     ListPermissionsUseCase,
		// rbacService:                rbacService,
		// permissionValidator:        permissionValidator,
	}
}
