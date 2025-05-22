package account

import (
	"backend/src/clients"
	error_utils "backend/src/error"
	"backend/src/graphql/middleware/authentication"
	"backend/src/graphql/modules/permission"
	"context"

	"github.com/google/uuid"
)

type AccountResolver struct {
	accountService    AccountService
	permissionService permission.PermissionService
}

func NewAccountResolver(params NewAccountParams) AccountResolver {

	return AccountResolver{
		accountService:    NewAccountService(params),
		permissionService: permission.NewPermissionService(params.Db),
	}
}

func (r AccountResolver) Profile(ctx context.Context, input ProfileInput) (Account, error) {
	authContext := authentication.GetAuthorizationContext(ctx)
	accountProfile, err := r.accountService.GetById(uuid.MustParse(authContext.AccountId))

	if err != nil {
		return Account{}, error_utils.GraphqlError{
			Message: err.Error(),
		}
	}
	return accountProfile, nil
}

func (r AccountResolver) Register(ctx context.Context, input RegisterInput) (Account, error) {
	account, err := r.accountService.Create(CreateAccountData{
		Username: input.Username,
		Password: input.Password,
	})

	if err != nil {
		return Account{}, error_utils.GraphqlError{
			Message: err.Error(),
		}
	}

	return account, nil
}

func (parent Account) Permissions(ctx context.Context) ([]permission.Permission, error) {
	if parent.ProjectId == nil {
		return []permission.Permission{}, nil
	}

	db, err := clients.NewPostgreSQLClient()
	if err != nil {
		return []permission.Permission{}, error_utils.GraphqlError{
			Message: error_utils.InternalServerError.Error(),
		}
	}

	permissionService := permission.PermissionService{
		Db: db,
	}

	permissions, err := permissionService.GetPermissionByProjectId(uuid.MustParse(string(*parent.ProjectId)))

	return permissions, nil
}
