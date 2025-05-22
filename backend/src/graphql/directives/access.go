package graphql_directives

import (
	"backend/src/.gen/cdr-intelligence/public/model"
	error_utils "backend/src/error"
	graphql_enum "backend/src/graphql/enum"
	"backend/src/graphql/middleware/authentication"
	"context"
)

type AccessDirective struct {
	RequiredProjectId *bool
	Subject           *string
	Ability           *graphql_enum.PermissionAbility
}

func (h *AccessDirective) ImplementsDirective() string {
	return "access"
}

func (h *AccessDirective) Validate(ctx context.Context, _ interface{}) error {
	authorization := authentication.GetAuthorizationContext(ctx)

	if authorization.AccountId == "" || authorization.Token == "" {
		return error_utils.GraphqlError{
			Message: "Account is unauthorized",
		}
	}

	if h.RequiredProjectId != nil && authorization.ProjectId == "" {
		return error_utils.GraphqlError{
			Message: "Project id is required",
		}
	}

	requiredPermission := h.Ability != nil && h.Subject != nil

	if requiredPermission {
		err := authentication.VerifyAuthorization(ctx, authorization, authentication.AuthorizationPermissionParams{
			PermissionSubject: *h.Subject,
			PermissionAbility: model.PermissionAbility(h.Ability.String()),
		})

		return err
	}

	return nil
}
