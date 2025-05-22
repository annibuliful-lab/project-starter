package permission

import (
	graphql_enum "backend/src/graphql/enum"

	"github.com/graph-gophers/graphql-go"
)

type Permission struct {
	Id          graphql.ID
	Name        string
	Description *string
	Subject     string
	Ability     graphql_enum.PermissionAbility
}
