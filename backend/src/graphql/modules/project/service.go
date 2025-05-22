package project

import (
	"database/sql"
	"time"

	"backend/src/.gen/cdr-intelligence/public/model"
	"backend/src/.gen/cdr-intelligence/public/table"
	error_utils "backend/src/error"
	shared_module "backend/src/graphql/modules/shared"

	pg "github.com/go-jet/jet/v2/postgres"
	"github.com/google/uuid"
	"github.com/graph-gophers/graphql-go"
	"github.com/redis/go-redis/v9"
)

type ProjectService struct {
	Db    *sql.DB
	Redis *redis.Client
}

func NewProjectService(params NewProjectResolverParams) ProjectService {
	return ProjectService{
		Db:    params.Db,
		Redis: params.Redis,
	}
}

func (service ProjectService) Create(input CreateProjectDataInput) (Project, error) {

	projectId, err := uuid.NewV7()
	if err != nil {
		return Project{}, err
	}

	project := model.Projects{}

	insertProjectStmt := table.Projects.
		INSERT(
			table.Projects.ID,
			table.Projects.Title,
			table.Projects.Description,
			table.Projects.CreatedBy,
		).
		MODEL(model.Projects{
			ID:          projectId,
			Title:       input.Title,
			Description: input.Description,
			CreatedBy:   input.CreatedBy,
		}).
		RETURNING(table.Projects.AllColumns)

	err = insertProjectStmt.Query(service.Db, &project)

	if err != nil {
		return Project{}, nil
	}

	return transformToGraphql(project), nil
}

func (service ProjectService) Delete(input DeleteProjectDataInput) (shared_module.DeleteOperation, error) {
	now := time.Now()
	softDeleteProjectStmt := table.Projects.
		UPDATE(
			table.Projects.DeletedAt,
			table.Projects.DeletedBy,
		).
		MODEL(model.Projects{
			DeletedAt: &now,
			DeletedBy: &input.DeletedBy,
		}).WHERE(table.Projects.ID.EQ(pg.UUID(input.Id)))

	updatedResult, err := softDeleteProjectStmt.Exec(service.Db)

	if err != nil {
		return shared_module.DeleteOperation{}, error_utils.InternalServerError
	}

	if error_utils.HasNoAffectedRow(updatedResult) {
		return shared_module.DeleteOperation{
			Success: false,
			Message: "Project not found",
		}, nil
	}

	return shared_module.DeleteOperation{
		Success: true,
		Message: "Project is deleted",
	}, nil
}

func (service ProjectService) FindById(id uuid.UUID) (Project, error) {

	project := model.Projects{}
	getProjectStmt := table.Projects.
		SELECT(table.Projects.AllColumns).
		FROM(table.Projects).
		WHERE(
			table.Projects.ID.EQ(pg.UUID(id)).
				AND(table.Projects.DeletedAt.IS_NULL()),
		).LIMIT(1)

	err := getProjectStmt.Query(service.Db, &project)

	if err != nil && error_utils.HasNoRow(err) {
		return Project{}, error_utils.Notfound
	}

	return transformToGraphql(project), nil
}

func (service ProjectService) FindByMany() ([]Project, error) {
	return []Project{}, nil
}

func transformToGraphql(data model.Projects) Project {
	return Project{
		Id:          graphql.ID(data.ID.String()),
		Title:       data.Title,
		Description: data.Description,
		CreatedAt:   graphql.Time{Time: data.CreatedAt},
		UpdatedAt:   graphql.Time{Time: *data.UpdatedAt},
	}
}
