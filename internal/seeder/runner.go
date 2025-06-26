package seeder

import (
	"context"
	"fmt"
	"time"

	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/container"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/internal/logger"
	"bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/permissions/domain/entity"
	roleEntity "bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/roles/domain/entity"
	userEntity "bitbucket.org/abhishek_fordel/we-care-holidays-backend-golang/modules/users/domain/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

// ExecuteSeeder runs seed operations for permissions, roles, users.
func ExecuteSeeder(appContainer *container.AppContainer, data *SeedData) error {
	ctx := context.Background()

	// --- Permissions Section ---
	logger.Log.Debug("Starting to seed permissions", zap.Int("permissionsToSeed", len(data.Permissions)))
	for _, p := range data.Permissions {
		logger.Log.Debug("Checking permission exists", zap.String("permission", fmt.Sprintf("%s:%s", p.Resource, p.Action)))
		existing, _, err := appContainer.Permission.Repository.List(ctx, make(map[string]interface{}), 1, 5000)
		if err != nil {
			logger.Log.Error("Error fetching existing permissions", zap.Error(err))
			return err
		}
		if permissionExists(existing, p.Resource, p.Action) {
			logger.Log.Info("Permission exists, skipping", zap.String("permission", fmt.Sprintf("%s:%s", p.Resource, p.Action)))
			continue
		}

		// FIXED: Ensure description is properly set
		perm := &entity.Permission{
			ID:          primitive.NewObjectID(),
			Resource:    p.Resource,
			Action:      entity.PermissionAction(p.Action),
			Description: p.Description, // This should now work properly
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		logger.Log.Debug("Creating permission",
			zap.String("permission", fmt.Sprintf("%s:%s", p.Resource, p.Action)),
			zap.String("description", p.Description), // Log description
			zap.String("id", perm.ID.Hex()))

		if err := appContainer.Permission.CreatePermissionUseCase.Execute(ctx, perm); err != nil {
			logger.Log.Error("Failed to seed permission", zap.String("permission", fmt.Sprintf("%s:%s", p.Resource, p.Action)), zap.Error(err))
			return err
		}
		logger.Log.Info("Seeded permission",
			zap.String("permission", fmt.Sprintf("%s:%s", p.Resource, p.Action)),
			zap.String("description", p.Description))
	}

	// --- Roles Section ---
	logger.Log.Debug("Starting to seed roles", zap.Int("rolesToSeed", len(data.Roles)))
	roleCollection := appContainer.MongoDatabase.Collection("roles")

	for _, r := range data.Roles {
		logger.Log.Debug("Processing role", zap.String("roleName", r.Name), zap.String("description", r.Description))

		// skip if already exists
		count, _ := roleCollection.CountDocuments(ctx, bson.M{"name": r.Name})
		if count > 0 {
			logger.Log.Info("Role exists, skipping", zap.String("role", r.Name))
			continue
		}

		// fetch _all_ permissions once
		allPerms, _, err := appContainer.Permission.Repository.List(ctx, make(map[string]interface{}), 1, 5000)
		if err != nil {
			return err
		}

		// build the slice of IDs we actually want
		var permIDs []string
		if len(r.Permissions) == 1 && r.Permissions[0] == "*" {
			// "*" means "give me every permission"
			for _, p := range allPerms {
				permIDs = append(permIDs, p.ID.Hex())
			}
		} else {
			// otherwise only the named ones
			for _, want := range r.Permissions {
				for _, p := range allPerms {
					permString := fmt.Sprintf("%s:%s", p.Resource, p.Action)
					if permString == want {
						permIDs = append(permIDs, p.ID.Hex())
						break
					}
				}
			}
		}

		// FIXED: Ensure description is properly set from seed data
		role := &roleEntity.Role{
			ID:          primitive.NewObjectID(),
			Name:        r.Name,
			Description: r.Description, // This should now be populated from JSON
			Scope:       roleEntity.RoleScope(r.Scope),
			Permissions: permIDs,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			IsSystem:    true,
		}

		logger.Log.Debug("Creating role",
			zap.String("role", r.Name),
			zap.String("description", r.Description),
			zap.Int("permissionCount", len(permIDs)))

		if err := appContainer.Role.CreateRoleUseCase.Execute(ctx, role); err != nil {
			return fmt.Errorf("failed to seed role %s: %w", r.Name, err)
		}
		logger.Log.Info("Seeded role with permissions",
			zap.String("role", r.Name),
			zap.String("description", r.Description))
	}

	// --- Users Section (FIXED) ---
	logger.Log.Debug("Starting to seed users", zap.Int("usersToSeed", len(data.Users)))
	userCollection := appContainer.MongoDatabase.Collection("users")

	for _, u := range data.Users {
		// FIXED: Check if user has emails
		if len(u.Emails) == 0 {
			logger.Log.Error("User has no emails, skipping", zap.String("user", u.FullName))
			continue
		}

		email := u.Emails[0].Email
		logger.Log.Debug("Checking user exists", zap.String("user", email))

		count, err := userCollection.CountDocuments(ctx, bson.M{"emails.email": email})
		if err != nil {
			logger.Log.Error("Error counting existing users", zap.String("user", email), zap.Error(err))
			return err
		}

		logger.Log.Debug("Existing users count", zap.String("user", email), zap.Int64("count", count))
		if count > 0 {
			logger.Log.Info("User exists, skipping", zap.String("user", email))
			continue
		}

		// Find role by name
		var roleDoc bson.M
		if err := appContainer.MongoDatabase.Collection("roles").FindOne(ctx, bson.M{"name": u.Role}).Decode(&roleDoc); err != nil {
			logger.Log.Error("Failed to find role for user", zap.String("user", email), zap.String("role", u.Role), zap.Error(err))
			return err
		}

		rawID, ok := roleDoc["_id"].(primitive.ObjectID)
		if !ok {
			logger.Log.Error("Invalid role ID type", zap.Any("_id", roleDoc["_id"]))
			return fmt.Errorf("the provided hex string is not a valid ObjectID: %v", roleDoc["_id"])
		}
		roleID := rawID.Hex()
		logger.Log.Debug("Resolved role ID for user", zap.String("user", email), zap.String("roleID", roleID))

		// FIXED: Proper user creation with all correct mappings
		user := &userEntity.User{
			ID:              primitive.NewObjectID(), // Generate new ID
			FullName:        u.FullName,
			RoleID:          roleID,
			Role:            u.Role, // FIXED: Set role name
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
			Password:        u.Password,
			Status:          mapUserStatus(u.Status), // FIXED: Use actual status from seed
			Phones:          mapPhones(u.Phones),
			Emails:          mapEmails(u.Emails),
			ProfilePhotoURL: u.ProfilePhotoURL,
			OrganizationID:  u.OrganizationID,
			AuditTrail:      userEntity.AuditTrail{}, // Initialize empty audit trail
		}

		logger.Log.Debug("Creating user in DB",
			zap.String("user", email),
			zap.String("fullName", u.FullName),
			zap.String("status", u.Status),
			zap.String("role", u.Role))

		if err := appContainer.User.CreateUserUseCase.Execute(ctx, user); err != nil {
			logger.Log.Error("Failed to seed user", zap.String("user", email), zap.Error(err))
			return err
		}
		logger.Log.Info("Seeded user",
			zap.String("user", email),
			zap.String("fullName", u.FullName),
			zap.String("status", string(user.Status)))
	}

	return nil
}

// permissionExists checks if permission already exists by name
func permissionExists(existing []*entity.Permission, resource, action string) bool {
	for _, e := range existing {
		if e.Resource == resource && e.Action == entity.PermissionAction(action) {
			return true
		}
	}
	return false
}
