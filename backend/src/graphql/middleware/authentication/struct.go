package authentication

import (
	"backend/src/.gen/cdr-intelligence/public/model"

	"github.com/google/uuid"
)

type AuthorizationHeader struct {
	Token     string
	ProjectId string
	AccountId string
}

type AuthorizationContext struct {
	Token     string
	ProjectId string
	AccountId string
}

type AuthorizationPermissionParams struct {
	PermissionSubject string
	PermissionAbility model.PermissionAbility
}

type AuthorizationWithPermissions struct {
	AccountId   uuid.UUID
	ProjectId   uuid.UUID
	Permissions []AuthorizationPermissionParams
}
