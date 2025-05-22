package project

import (
	error_utils "backend/src/error"
	"backend/src/graphql/middleware/authentication"
	"context"

	"github.com/google/uuid"
)

type ProjectResolver struct {
	projectService ProjectService
}

func NewProjectResolver(params NewProjectResolverParams) ProjectResolver {
	return ProjectResolver{
		projectService: NewProjectService(params),
	}
}

func (r ProjectResolver) CreateProject(ctx context.Context, input CreateProjectInput) (Project, error) {
	authContext := authentication.GetAuthorizationContext(ctx)

	project, err := r.projectService.Create(CreateProjectDataInput{
		Title:       input.Title,
		Description: input.Description,
		CreatedBy:   authContext.AccountId,
	})

	if err != nil {
		return Project{}, error_utils.GraphqlError{
			Message: err.Error(),
		}
	}

	return project, nil
}

func (r ProjectResolver) GetProjectById(ctx context.Context, input GetProjectByIdInput) (Project, error) {
	project, err := r.projectService.FindById(uuid.MustParse(string(input.Id)))

	if err != nil {
		return Project{}, error_utils.GraphqlError{
			Message: err.Error(),
		}
	}

	return project, nil
}
