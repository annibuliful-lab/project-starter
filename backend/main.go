package main

import (
	"backend/src/clients"
	"backend/src/config"
	"backend/src/graphql"
	graphql_directives "backend/src/graphql/directives"
	"backend/src/graphql/middleware/authentication"
	uploadmiddleware "backend/src/graphql/middleware/upload"
	"backend/src/graphql/subscription/graphqlws"
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"

	merger "github.com/annibuliful-lab/merge-graphql-schema"
	gql "github.com/graph-gophers/graphql-go"
	relay "github.com/graph-gophers/graphql-go/relay"
	"github.com/joho/godotenv"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Panic().Msg("Error loading .env" + err.Error())
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	merger.MergeSchemas("./src", ".graphql", "generated.graphql")
	mergedSchema, err := os.ReadFile("generated.graphql")

	if err != nil {
		log.Panic().Msg("Error loading graphql file")
	}

	isDevelopment := config.GetEnv("ENV", "development") == "development"

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		// Debug:            isDevelopment,
	})

	opts := []gql.SchemaOpt{
		gql.UseFieldResolvers(),
		gql.MaxParallelism(20),
		gql.UseStringDescriptions(),
		gql.RestrictIntrospection(func(context.Context) bool {
			return isDevelopment
		}),
		gql.Directives(&graphql_directives.AccessDirective{}),
	}

	db, err := clients.NewPostgreSQLClient()
	if err != nil {
		log.Panic().Msg(err.Error())
	}
	redis, err := clients.NewRedisClient()
	if err != nil {
		log.Panic().Msg(err.Error())
	}

	rabbitmq, err := clients.NewRabbitMQClient()
	if err != nil {
		log.Panic().Msg(err.Error())
	}
	ctx := context.Background()
	neo4jDriver, err := clients.NewNeo4jClient()
	neo4jSession := neo4jDriver.NewSession(ctx, neo4j.SessionConfig{})
	if err != nil {
		log.Panic().Msg(err.Error())
	}

	schema, err := gql.ParseSchema(string(mergedSchema[:]), graphql.GraphqlResolver(graphql.GraphqlResolverParams{
		Db:       db,
		Redis:    redis,
		Rabbitmq: rabbitmq,
		Neo4j:    &neo4jSession,
	}), opts...)

	if err != nil {
		log.Panic().Msg(err.Error())
	}

	// graphQL handler
	graphQLHandler := corsMiddleware.Handler(
		graphqlws.NewHandlerFunc(
			schema,
			authentication.GraphqlContext(uploadmiddleware.Handler(&relay.Handler{Schema: schema})),
		),
	)

	http.Handle("/graphql", graphQLHandler)

	var listenAddress = flag.String("listen", config.GetEnv("BACKEND_PORT", ":3000"), "Listen address.")

	log.Info().Msg("Listening at http://" + *listenAddress)

	httpServer := http.Server{
		Addr: *listenAddress,
	}

	idleConnectionsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint
		if err := httpServer.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP Server Shutdown Error: %v", err)
		}
		db.Close()
		redis.Close()
		rabbitmq.Close()
		neo4jSession.Close(ctx)
		neo4jDriver.Close(ctx)

		close(idleConnectionsClosed)
	}()

	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal().Msg("HTTP server ListenAndServe Error:" + err.Error())
	}

	<-idleConnectionsClosed
}
