package authentication

import (
	"backend/src/.gen/cdr-intelligence/public/model"
	"backend/src/.gen/cdr-intelligence/public/table"
	"backend/src/clients"
	error_utils "backend/src/error"
	"context"
	"database/sql"
	"encoding/json"

	"strings"
	"time"

	pg "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"
)

func GetAuthToken(authorization string) string {
	return strings.Replace(authorization, "Bearer ", "", 1)
}

func GetAuthorizationContext(ctx context.Context) AuthorizationContext {
	authContext := AuthorizationContext{}

	if ctx.Value("accountId") != nil {
		authContext.AccountId = ctx.Value("accountId").(string)
	}

	if ctx.Value("projectId") != nil {
		authContext.ProjectId = ctx.Value("projectId").(string)
	}

	if ctx.Value("token") != nil {
		authContext.Token = ctx.Value("token").(string)
	}
	return authContext
}

func VerifyAuthorization(ctx context.Context, headers AuthorizationContext, permission AuthorizationPermissionParams) error {
	if headers.Token == "" || headers.AccountId == "" || headers.ProjectId == "" {
		return error_utils.GraphqlError{
			Message: error_utils.TokenIsInvalid.Error(),
		}
	}

	key := "accountId:" + headers.AccountId + "," + "projectId:" + headers.ProjectId

	// Check if the authorization data is present in the cache
	redis, err := clients.NewRedisClient()
	if err != nil {
		return err
	}
	result, err := redis.Get(ctx, key).Result()
	if err == nil && result != "" {
		cacheError := handleCachedAuthorization(result, permission)
		if cacheError != nil {
			return error_utils.GraphqlError{
				Message: cacheError.Error(),
			}
		}

		return nil
	}

	projectId, err := uuid.Parse(headers.ProjectId)
	if err != nil {
		return error_utils.GraphqlError{
			Message: "Project id is required",
		}
	}

	db, err := clients.NewPostgreSQLClient()
	if err != nil {
		return err
	}

	projectAccount, err := getProjectAccount(db, uuid.MustParse(headers.AccountId), projectId)
	if err != nil {
		log.Printf("get-project-account-error: %v", err.Error())

		return error_utils.GraphqlError{
			Message: error_utils.ForbiddenOperation.Error(),
		}
	}

	accountPermissions, err := getAccountPermissions(db, projectAccount.ProjectRoleId)
	if err != nil {
		log.Printf("get-account-permission-error: %v", err.Error())
		return error_utils.GraphqlError{
			Message: error_utils.InternalServerError.Error(),
		}
	}

	if !hasPermission(accountPermissions, permission) {

		return error_utils.GraphqlError{
			Message: error_utils.ForbiddenOperation.Error(),
		}
	}

	cacheAuthorization(cacheAuthorizationParams{
		ctx:                ctx,
		key:                key,
		accountID:          uuid.MustParse(headers.AccountId),
		projectID:          projectId,
		accountPermissions: accountPermissions,
		redis:              *redis,
	})

	return nil
}

func handleCachedAuthorization(result string, permissionData AuthorizationPermissionParams) error {
	var data AuthorizationWithPermissions

	if err := json.Unmarshal([]byte(result), &data); err != nil {
		log.Printf("cache-authorization-error: %v", err.Error())
		return error_utils.GraphqlError{
			Message: err.Error(),
		}
	}

	_, match := lo.Find(data.Permissions, func(el AuthorizationPermissionParams) bool {
		return el.PermissionAbility == permissionData.PermissionAbility && el.PermissionSubject == permissionData.PermissionSubject
	})

	if !match {

		return error_utils.GraphqlError{
			Message: error_utils.ForbiddenOperation.Error(),
		}
	}

	return nil
}

func getProjectAccount(dbClient *sql.DB, accountID uuid.UUID, projectID uuid.UUID) (struct {
	model.ProjectAccounts
	model.ProjectRoles
}, error) {

	var projectAccount struct {
		model.ProjectAccounts
		model.ProjectRoles
	}

	selectProjectAccountStmt := pg.
		SELECT(
			table.ProjectAccounts.AccountId,
			table.ProjectRoles.ID,
			table.ProjectRoles.Title,
		).
		FROM(
			table.ProjectAccounts.
				INNER_JOIN(table.ProjectRoles, table.ProjectRoles.ID.EQ(table.ProjectAccounts.ProjectRoleId)),
		).
		WHERE(
			table.ProjectAccounts.AccountId.EQ(pg.UUID(accountID)).
				AND(table.ProjectAccounts.ProjectId.EQ(pg.UUID(projectID))))

	err := selectProjectAccountStmt.Query(dbClient, &projectAccount)
	return projectAccount, err
}

func getAccountPermissions(dbClient *sql.DB, roleID uuid.UUID) ([]struct{ model.Permission }, error) {

	var accountPermissions []struct{ model.Permission }

	selectProjectAccountPermissionsStmt := pg.
		SELECT(table.Permission.Ability, table.Permission.Subject, table.Permission.ID).
		FROM(
			table.ProjectRolePermission.
				INNER_JOIN(table.Permission, table.Permission.ID.EQ(table.ProjectRolePermission.PermissionId)),
		).WHERE(table.ProjectRolePermission.ProjectRoleId.EQ(pg.UUID(roleID)))

	err := selectProjectAccountPermissionsStmt.Query(dbClient, &accountPermissions)

	return accountPermissions, err
}

func hasPermission(accountPermissions []struct{ model.Permission }, permissionData AuthorizationPermissionParams) bool {
	match := lo.ContainsBy(accountPermissions, func(el struct{ model.Permission }) bool {
		return el.Ability == permissionData.PermissionAbility && el.Subject == permissionData.PermissionSubject
	})

	return match
}

type cacheAuthorizationParams struct {
	ctx                context.Context
	key                string
	accountID          uuid.UUID
	projectID          uuid.UUID
	accountPermissions []struct{ model.Permission }
	redis              redis.Client
}

func cacheAuthorization(params cacheAuthorizationParams) {
	data := AuthorizationWithPermissions{
		AccountId: params.accountID,
		ProjectId: params.projectID,
		Permissions: lo.Map(params.accountPermissions, func(item struct{ model.Permission }, index int) AuthorizationPermissionParams {
			return AuthorizationPermissionParams{
				PermissionSubject: item.Subject,
				PermissionAbility: item.Ability,
			}
		}),
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("cache-authorization-json-error: %v", err.Error())
		return
	}

	err = params.redis.Set(params.ctx, params.key, jsonData, 15*time.Minute).Err()
	if err != nil {
		log.Printf("cache-authorization-error: %v", err.Error())
	}
}
