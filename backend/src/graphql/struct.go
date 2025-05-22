package graphql

import (
	"database/sql"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

type GraphqlResolverParams struct {
	Db       *sql.DB
	Redis    *redis.Client
	Neo4j    *neo4j.SessionWithContext
	Rabbitmq *amqp091.Connection
}
