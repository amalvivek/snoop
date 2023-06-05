package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"net/http"
	"snoop-server/pkg/database"
)

type PropertyResult struct {
	Property `json:"property"`
}

type Property struct {
	ID      string `json:"id"`
	Address string `json:"address"`
}

func (p Property) Reconstruct(record *neo4j.Record) Property {
	p.ID, _, _ = neo4j.GetRecordValue[string](record, "id")
	p.Address, _, _ = neo4j.GetRecordValue[string](record, "address")
	return p
}

func SearchProperties(c *gin.Context) {
	searchTerm, _ := c.GetQuery("q")
	var query = "MATCH (n:Properties) WHERE n.address =~ '(?i)^.*" + searchTerm + ".*$' RETURN n.id as id, n.address as address;"

	res, err := database.ExecuteQuery(c, database.Query{String: query}, nil)
	if err != nil {
		return
	}

	properties := make([]PropertyResult, len(res.Records))

	for i, record := range res.Records {
		properties[i] = PropertyResult{Property{}.Reconstruct(record)}
	}

	c.IndentedJSON(http.StatusOK, properties)
}

func GetProperty(c *gin.Context) {
	id := c.Param("id")
	var query = "MATCH (n:Properties) WHERE n.id = '" + id + "' RETURN n.id as id, n.address as address;"

	res, err := database.ExecuteQuery(c, database.Query{String: query}, nil)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.IndentedJSON(http.StatusOK, Property{}.Reconstruct(res.Records[0]))
	}
}
