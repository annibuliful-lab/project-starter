package account

import (
	"database/sql"

	"github.com/graph-gophers/graphql-go"
	"github.com/redis/go-redis/v9"
)

type NewAccountParams struct {
	Db    *sql.DB
	Redis *redis.Client
}

type Account struct {
	Id        graphql.ID
	Username  string
	CreatedAt graphql.Time
	ProjectId *graphql.ID
}

type RegisterInput struct {
	Username string
	Password string
}

type CreateAccountData struct {
	Username string
	Password string
}

type ProfileInput struct {
	ProjectId string
}
