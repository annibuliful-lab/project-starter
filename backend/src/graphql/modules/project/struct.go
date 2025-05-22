package project

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/graph-gophers/graphql-go"
	"github.com/redis/go-redis/v9"
)

type NewProjectResolverParams struct {
	Db    *sql.DB
	Redis *redis.Client
}

type CreateProjectInput struct {
	Title       string
	Description *string
}

type CreateProjectDataInput struct {
	Title       string
	Description *string
	CreatedBy   string
}

type DeleteProjectInput struct {
	Id graphql.ID
}

type DeleteProjectDataInput struct {
	Id        uuid.UUID
	DeletedBy string
}

type Project struct {
	Id          graphql.ID
	Title       string
	Description *string
	CreatedAt   graphql.Time
	UpdatedAt   graphql.Time
}

type GetProjectByIdInput struct {
	Id graphql.ID
}
