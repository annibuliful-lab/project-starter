package authentication

import (
	"database/sql"

	"github.com/graph-gophers/graphql-go"
	"github.com/redis/go-redis/v9"
)

type Authentication struct {
	AccountId    graphql.ID
	ProjectIds   *[]graphql.ID
	Token        string
	RefreshToken string
}

type NewAuthenticationResolverParams struct {
	Db    *sql.DB
	Redis *redis.Client
}

type LoginInput struct {
	Username string
	Password string
}

type LoginData struct {
	Username string
	Password string
}

type Logout struct {
	Success bool
	Message string
}
