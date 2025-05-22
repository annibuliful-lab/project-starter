package projectrole

import (
	"database/sql"

	"backend/src/.gen/cdr-intelligence/public/model"
	"backend/src/.gen/cdr-intelligence/public/table"

	"github.com/google/uuid"
	"github.com/graph-gophers/graphql-go"
	"github.com/redis/go-redis/v9"
)

type ProjectRoleService struct {
	Db    *sql.DB
	redis *redis.Client
}

func NewProjectRoleService(params NewProjectRoleResolverParams) ProjectRoleService {
	return ProjectRoleService{
		Db:    params.Db,
		redis: params.Redis,
	}
}

func (service ProjectRoleService) Create(input CreateProjectDataInput) (ProjectRole, error) {
	roleId, err := uuid.NewV7()
	if err != nil {
		return ProjectRole{}, err
	}

	projectRole := model.ProjectRoles{}
	insertProjectRoleStmt := table.ProjectRoles.
		INSERT(table.ProjectRoles.ID,
			table.ProjectRoles.Title,
			table.ProjectRoles.CreatedBy,
		).
		MODEL(model.ProjectRoles{
			ID:        roleId,
			ProjectId: input.ProjectId,
			Title:     input.Title,
			CreatedBy: input.CreatedBy,
		}).
		RETURNING(table.ProjectRoles.AllColumns)

	err = insertProjectRoleStmt.Query(service.Db, &projectRole)

	if err != nil {
		return ProjectRole{}, nil
	}

	return transformToGraphql(projectRole), nil
}

func transformToGraphql(data model.ProjectRoles) ProjectRole {
	return ProjectRole{
		Id:        graphql.ID(data.ID.String()),
		title:     data.Title,
		projectId: graphql.ID(data.ProjectId.String()),
	}
}
