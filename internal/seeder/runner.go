package seeder

import (
	"context"
	"fmt"
	"time"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/container"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/logger"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/organizations/domain/entity"
	permissionEntity "bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/permissions/domain/entity"
	roleEntity "bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/roles/domain/entity"
	userEntity "bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/users/domain/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

// ExecuteSeeder runs seed operations for permissions, organizations, roles, users.
// Order is important: permissions -> organizations -> roles -> users
func ExecuteSeeder(appContainer *container.AppContainer, data *SeedData) error {
	ctx := context.Background()

	logger.Log.Info("🔧 Starting seeder execution")

	// Step 1: Seed Permissions
	logger.Log.Info("📋 Step 1: Seeding permissions")
	if err := seedPermissions(ctx, appContainer, data.Permissions); err != nil {
		logger.Log.Error("❌ Failed to seed permissions", zap.Error(err))
		return fmt.Errorf("failed to seed permissions: %w", err)
	}
	logger.Log.Info("✅ Step 1 completed: Permissions seeded")

	// Step 2: Seed Organizations
	logger.Log.Info("🏢 Step 2: Seeding organizations")
	if err := seedOrganizations(ctx, appContainer, data.Organizations); err != nil {
		logger.Log.Error("❌ Failed to seed organizations", zap.Error(err))
		return fmt.Errorf("failed to seed organizations: %w", err)
	}
	logger.Log.Info("✅ Step 2 completed: Organizations seeded")

	// Step 3: Seed Roles
	logger.Log.Info("👤 Step 3: Seeding roles")
	if err := seedRoles(ctx, appContainer, data.Roles); err != nil {
		logger.Log.Error("❌ Failed to seed roles", zap.Error(err))
		return fmt.Errorf("failed to seed roles: %w", err)
	}
	logger.Log.Info("✅ Step 3 completed: Roles seeded")

	// Step 4: Seed Users (with organization linking)
	logger.Log.Info("👥 Step 4: Seeding users")
	if err := seedUsers(ctx, appContainer, data.Users); err != nil {
		logger.Log.Error("❌ Failed to seed users", zap.Error(err))
		return fmt.Errorf("failed to seed users: %w", err)
	}
	logger.Log.Info("✅ Step 4 completed: Users seeded")

	logger.Log.Info("🎉 All seeding steps completed successfully")
	return nil
}

// seedPermissions handles permission seeding
func seedPermissions(ctx context.Context, appContainer *container.AppContainer, permissions []PermissionSeed) error {
	logger.Log.Info("🔐 Starting to seed permissions", zap.Int("count", len(permissions)))

	for i, p := range permissions {
		logger.Log.Debug("Processing permission",
			zap.Int("index", i+1),
			zap.Int("total", len(permissions)),
			zap.String("permission", fmt.Sprintf("%s:%s", p.Resource, p.Action)))

		existing, _, err := appContainer.Permission.Repository.List(ctx, make(map[string]interface{}), 1, 5000)
		if err != nil {
			logger.Log.Error("Error fetching existing permissions", zap.Error(err))
			return err
		}

		if permissionExists(existing, p.Resource, p.Action) {
			logger.Log.Debug("Permission exists, skipping", zap.String("permission", fmt.Sprintf("%s:%s", p.Resource, p.Action)))
			continue
		}

		perm := &permissionEntity.Permission{
			ID:          primitive.NewObjectID(),
			Resource:    p.Resource,
			Action:      permissionEntity.PermissionAction(p.Action),
			Description: p.Description,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		if err := appContainer.Permission.CreatePermissionUseCase.Execute(ctx, perm); err != nil {
			logger.Log.Error("Failed to seed permission", zap.String("permission", fmt.Sprintf("%s:%s", p.Resource, p.Action)), zap.Error(err))
			return err
		}

		logger.Log.Info("✅ Seeded permission", zap.String("permission", fmt.Sprintf("%s:%s", p.Resource, p.Action)))
	}

	logger.Log.Info("🔐 Permissions seeding completed", zap.Int("count", len(permissions)))
	return nil
}

// seedOrganizations handles organization seeding
func seedOrganizations(ctx context.Context, appContainer *container.AppContainer, organizations []OrganizationSeed) error {
	logger.Log.Info("🏢 Starting to seed organizations", zap.Int("count", len(organizations)))

	if len(organizations) == 0 {
		logger.Log.Warn("No organizations to seed")
		return nil
	}

	// Check if organization container exists
	if appContainer.Organization == nil {
		logger.Log.Error("Organization container is nil")
		return fmt.Errorf("organization container is not initialized")
	}

	if appContainer.Organization.CreateOrganizationUseCase == nil {
		logger.Log.Error("CreateOrganizationUseCase is nil")
		return fmt.Errorf("CreateOrganizationUseCase is not initialized")
	}

	orgCollection := appContainer.MongoDatabase.Collection("organizations")

	for i, o := range organizations {
		logger.Log.Debug("Processing organization",
			zap.Int("index", i+1),
			zap.Int("total", len(organizations)),
			zap.String("name", o.Name),
			zap.String("slug", o.Slug))

		// Check if organization already exists by slug
		count, err := orgCollection.CountDocuments(ctx, bson.M{"slug": o.Slug})
		if err != nil {
			logger.Log.Error("Error counting existing organizations", zap.String("slug", o.Slug), zap.Error(err))
			return err
		}

		if count > 0 {
			logger.Log.Info("Organization exists, skipping", zap.String("slug", o.Slug))
			continue
		}

		// Create organization entity
		org := &entity.Organization{
			ID:      primitive.NewObjectID(),
			Name:    o.Name,
			Slug:    o.Slug,
			Type:    o.Type,
			Email:   o.Email,
			Phone:   o.Phone,
			Website: o.Website,
			TaxIDs:  o.TaxIDs,
			Logo:    o.Logo,
			Address: entity.Address{
				Street:  o.Address.Street,
				City:    o.Address.City,
				State:   o.Address.State,
				Country: o.Address.Country,
				Pincode: o.Address.Pincode,
			},
			Status:    o.Status,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		logger.Log.Debug("Creating organization",
			zap.String("name", o.Name),
			zap.String("type", o.Type),
			zap.String("id", org.ID.Hex()))

		if err := appContainer.Organization.CreateOrganizationUseCase.Execute(ctx, org); err != nil {
			logger.Log.Error("Failed to seed organization",
				zap.String("name", o.Name),
				zap.String("slug", o.Slug),
				zap.Error(err))
			return err
		}

		logger.Log.Info("✅ Seeded organization",
			zap.String("name", o.Name),
			zap.String("slug", o.Slug),
			zap.String("id", org.ID.Hex()))
	}

	logger.Log.Info("🏢 Organizations seeding completed", zap.Int("count", len(organizations)))
	return nil
}

// seedRoles handles role seeding
func seedRoles(ctx context.Context, appContainer *container.AppContainer, roles []RoleSeed) error {
	logger.Log.Info("👤 Starting to seed roles", zap.Int("count", len(roles)))

	roleCollection := appContainer.MongoDatabase.Collection("roles")

	for i, r := range roles {
		logger.Log.Debug("Processing role",
			zap.Int("index", i+1),
			zap.Int("total", len(roles)),
			zap.String("name", r.Name),
			zap.String("scope", r.Scope))

		// Check if role already exists
		count, _ := roleCollection.CountDocuments(ctx, bson.M{"name": r.Name})
		if count > 0 {
			logger.Log.Info("Role exists, skipping", zap.String("role", r.Name))
			continue
		}

		// Fetch all permissions for role assignment
		allPerms, _, err := appContainer.Permission.Repository.List(ctx, make(map[string]interface{}), 1, 5000)
		if err != nil {
			logger.Log.Error("Failed to fetch permissions for role", zap.String("role", r.Name), zap.Error(err))
			return err
		}

		// Build permission IDs based on role requirements
		var permIDs []string
		if len(r.Permissions) == 1 && r.Permissions[0] == "*" {
			// "*" means "give me every permission"
			for _, p := range allPerms {
				permIDs = append(permIDs, p.ID.Hex())
			}
			logger.Log.Debug("Assigning all permissions to role",
				zap.String("role", r.Name),
				zap.Int("permission_count", len(permIDs)))
		} else {
			// Otherwise only the named ones
			for _, want := range r.Permissions {
				for _, p := range allPerms {
					permString := fmt.Sprintf("%s:%s", p.Resource, p.Action)
					if permString == want {
						permIDs = append(permIDs, p.ID.Hex())
						break
					}
				}
			}
			logger.Log.Debug("Assigning specific permissions to role",
				zap.String("role", r.Name),
				zap.Int("permission_count", len(permIDs)),
				zap.Strings("permissions", r.Permissions))
		}

		role := &roleEntity.Role{
			ID:          primitive.NewObjectID(),
			Name:        r.Name,
			Description: r.Description,
			Scope:       roleEntity.RoleScope(r.Scope),
			Permissions: permIDs,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			IsSystem:    true,
		}

		if err := appContainer.Role.CreateRoleUseCase.Execute(ctx, role); err != nil {
			logger.Log.Error("Failed to seed role", zap.String("role", r.Name), zap.Error(err))
			return fmt.Errorf("failed to seed role %s: %w", r.Name, err)
		}

		logger.Log.Info("✅ Seeded role",
			zap.String("role", r.Name),
			zap.String("scope", r.Scope),
			zap.String("id", role.ID.Hex()))
	}

	logger.Log.Info("👤 Roles seeding completed", zap.Int("count", len(roles)))
	return nil
}

// seedUsers handles user seeding with organization linking
func seedUsers(ctx context.Context, appContainer *container.AppContainer, users []UserSeed) error {
	logger.Log.Info("👥 Starting to seed users", zap.Int("count", len(users)))

	userCollection := appContainer.MongoDatabase.Collection("users")
	organizationCollection := appContainer.MongoDatabase.Collection("organizations")

	for i, u := range users {
		if len(u.Emails) == 0 {
			logger.Log.Error("User has no emails, skipping", zap.String("user", u.FullName))
			continue
		}

		email := u.Emails[0].Email
		logger.Log.Debug("Processing user",
			zap.Int("index", i+1),
			zap.Int("total", len(users)),
			zap.String("email", email),
			zap.String("org", u.OrganizationSlug))

		// Check if user already exists
		count, err := userCollection.CountDocuments(ctx, bson.M{"emails.email": email})
		if err != nil {
			logger.Log.Error("Error counting existing users", zap.String("email", email), zap.Error(err))
			return err
		}

		if count > 0 {
			logger.Log.Info("User exists, skipping", zap.String("email", email))
			continue
		}

		// Find organization by slug
		var orgDoc bson.M
		if err := organizationCollection.FindOne(ctx, bson.M{"slug": u.OrganizationSlug}).Decode(&orgDoc); err != nil {
			logger.Log.Error("Failed to find organization for user",
				zap.String("email", email),
				zap.String("orgSlug", u.OrganizationSlug),
				zap.Error(err))
			return err
		}

		rawOrgID, ok := orgDoc["_id"].(primitive.ObjectID)
		if !ok {
			logger.Log.Error("Invalid organization ID type", zap.Any("_id", orgDoc["_id"]))
			return fmt.Errorf("invalid organization ID type: %v", orgDoc["_id"])
		}
		organizationID := rawOrgID.Hex()

		// Find role by name
		var roleDoc bson.M
		if err := appContainer.MongoDatabase.Collection("roles").FindOne(ctx, bson.M{"name": u.Role}).Decode(&roleDoc); err != nil {
			logger.Log.Error("Failed to find role for user",
				zap.String("email", email),
				zap.String("role", u.Role),
				zap.Error(err))
			return err
		}

		rawRoleID, ok := roleDoc["_id"].(primitive.ObjectID)
		if !ok {
			logger.Log.Error("Invalid role ID type", zap.Any("_id", roleDoc["_id"]))
			return fmt.Errorf("invalid role ID type: %v", roleDoc["_id"])
		}
		roleID := rawRoleID.Hex()

		// Create user entity
		user := &userEntity.User{
			ID:              primitive.NewObjectID(),
			FullName:        u.FullName,
			RoleID:          roleID,
			Role:            u.Role,
			OrganizationID:  organizationID, // Link to organization
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
			Password:        u.Password,
			Status:          mapUserStatus(u.Status),
			Phones:          mapPhones(u.Phones),
			Emails:          mapEmails(u.Emails),
			ProfilePhotoURL: u.ProfilePhotoURL,
			AuditTrail:      userEntity.AuditTrail{},
		}

		if err := appContainer.User.CreateUserUseCase.Execute(ctx, user); err != nil {
			logger.Log.Error("Failed to seed user", zap.String("email", email), zap.Error(err))
			return err
		}

		logger.Log.Info("✅ Seeded user",
			zap.String("email", email),
			zap.String("role", u.Role),
			zap.String("organization", u.OrganizationSlug),
			zap.String("id", user.ID.Hex()))
	}

	logger.Log.Info("👥 Users seeding completed", zap.Int("count", len(users)))
	return nil
}

// permissionExists checks if permission already exists by name
func permissionExists(existing []*permissionEntity.Permission, resource, action string) bool {
	for _, e := range existing {
		if e.Resource == resource && e.Action == permissionEntity.PermissionAction(action) {
			return true
		}
	}
	return false
}
