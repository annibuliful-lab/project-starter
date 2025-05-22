package permission

import (
	"backend/src/.gen/cdr-intelligence/public/model"
	"backend/src/.gen/cdr-intelligence/public/table"
	error_utils "backend/src/error"
	graphql_enum "backend/src/graphql/enum"
	"database/sql"

	pg "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/graph-gophers/graphql-go"
	"github.com/samber/lo"
)

type PermissionService struct {
	Db *sql.DB
}

func NewPermissionService(db *sql.DB) PermissionService {
	return PermissionService{
		Db: db,
	}
}

func (service PermissionService) GetPermissionByProjectId(projectId uuid.UUID) ([]Permission, error) {
	permissions := []model.Permission{}

	getPermissionsStmt := table.Permission.
		SELECT(table.Permission.AllColumns).
		FROM(
			table.Permission.
				INNER_JOIN(
					table.ProjectRolePermission,
					table.ProjectRolePermission.PermissionId.EQ(table.Permission.ID),
				).INNER_JOIN(table.ProjectRoles, table.ProjectRoles.ID.EQ(table.ProjectRolePermission.ProjectRoleId)),
		).
		WHERE(
			table.ProjectRoles.ProjectId.EQ(pg.UUID(projectId)),
		)

	err := getPermissionsStmt.Query(service.Db, &permissions)

	if err != nil {
		return []Permission{}, error_utils.InternalServerError
	}

	return lo.Map(permissions, func(permission model.Permission, index int) Permission {
		return transformToGraphql(permission)
	}), nil
}

func transformToGraphql(permission model.Permission) Permission {
	return Permission{
		Id:          graphql.ID(permission.ID.String()),
		Name:        permission.Name,
		Description: permission.Description,
		Subject:     permission.Subject,
		Ability:     graphql_enum.GetPermissionAbility(permission.Ability.String()),
	}
}
