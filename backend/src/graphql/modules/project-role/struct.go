package projectrole

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/graph-gophers/graphql-go"
	"github.com/redis/go-redis/v9"
)

type ProjectRole struct {
	Id        graphql.ID
	projectId graphql.ID
	title     string
}

type NewProjectRoleResolverParams struct {
	Db    *sql.DB
	Redis *redis.Client
}

type CreateProjectRoleInput struct {
	Title string
}

type CreateProjectDataInput struct {
	CreatedBy string
	Title     string
	ProjectId uuid.UUID
}
