package projectrole

import "context"

type ProjectRoleResolver struct {
	projectRoleService ProjectRoleService
}

func NewProjectRoleResolver(params NewProjectRoleResolverParams) ProjectRoleResolver {

	return ProjectRoleResolver{
		projectRoleService: NewProjectRoleService(params),
	}
}

func (ProjectRoleResolver) CreateProjectRole(ctx context.Context, input CreateProjectRoleInput) (ProjectRole, error) {

	return ProjectRole{}, nil
}
