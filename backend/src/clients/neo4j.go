package clients

import (
	"backend/src/config"
	"sync"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

var (
	neo4jDriver neo4j.DriverWithContext // Now a pointer to an interface
	neo4jOnce   sync.Once
)

// NewNeo4jClient ensures thread-safe initialization of the Neo4j driver.
func NewNeo4jClient() (neo4j.DriverWithContext, error) {

	neo4jOnce.Do(func() {
		driver, err := neo4j.NewDriverWithContext(
			config.GetEnv("NEO4J_URI", "bolt://localhost:7687"),
			neo4j.BasicAuth(
				config.GetEnv("NEO4J_USERNAME", "neo4j"),
				config.GetEnv("NEO4J_PASSWORD", "test"),
				"",
			),
		)
		if err == nil {
			neo4jDriver = driver // Store the address of the driver
		}
	})

	return neo4jDriver, nil
}
