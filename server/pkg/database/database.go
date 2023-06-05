package database

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
	"net/http"
	"os"
	"strings"
)

type Result interface {
	Reconstruct(result *neo4j.EagerResult) any
}

type Query struct {
	String string
	Type   string
}

type Neo4jConfiguration struct {
	Url      string
	Username string
	Password string
	Database string
}

var driver neo4j.DriverWithContext
var configuration *Neo4jConfiguration

func init() {
	configuration = parseConfiguration()
	var err error
	driver, err = configuration.newDriver()
	if err != nil {
		log.Fatal(err)
	}
}

func (nc *Neo4jConfiguration) newDriver() (neo4j.DriverWithContext, error) {
	return neo4j.NewDriverWithContext(nc.Url, neo4j.BasicAuth(nc.Username, nc.Password, ""))
}

func parseConfiguration() *Neo4jConfiguration {
	database := lookupEnvOrGetDefault("NEO4J_DATABASE", "movies")
	if !strings.HasPrefix(lookupEnvOrGetDefault("NEO4J_VERSION", "5"), "4") {
		database = ""
	}

	return &Neo4jConfiguration{
		Url:      lookupEnvOrGetDefault("NEO4J_URI", "neo4j+s://87d3d17a.databases.neo4j.io"),
		Username: lookupEnvOrGetDefault("NEO4J_USER", "neo4j"),
		Password: lookupEnvOrGetDefault("NEO4J_PASSWORD", "password"),
		Database: database,
	}
}

func lookupEnvOrGetDefault(key string, defaultValue string) string {
	if env, found := os.LookupEnv(key); !found {
		return defaultValue
	} else {
		return env
	}
}

func UnsafeClose(ctx context.Context) {
	if err := driver.Close(ctx); err != nil {
		log.Fatal(fmt.Errorf("could not close resource: %w", err))
	}
}

func ErrorHandler(c *gin.Context) {
	c.Next()

	for _, err := range c.Errors {
		// log, handle, etc.
		log.Println(err)
	}

	c.JSON(http.StatusInternalServerError, "")
}

func ExecuteQuery(ctx context.Context, query Query, parameters map[string]any) (res *neo4j.EagerResult, err error) {

	var routing neo4j.ExecuteQueryConfigurationOption
	switch query.Type {
	case "READ":
		routing = neo4j.ExecuteQueryWithReadersRouting()
	case "WRITE":
		routing = neo4j.ExecuteQueryWithWritersRouting()
	default:
		routing = neo4j.ExecuteQueryWithReadersRouting()
	}

	return neo4j.ExecuteQuery(ctx, driver, query.String, parameters,
		neo4j.EagerResultTransformer,
		routing,
		neo4j.ExecuteQueryWithDatabase(configuration.Database))
}
