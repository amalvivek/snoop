package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
	"net/http"
	"snoop-server/pkg/database"
)

type PropertyResult struct {
	Property `json:"property"`
}

type Property struct {
	ID      int64  `json:"id"`
	Address string `json:"address"`
}

func (p Property) Reconstruct(record *neo4j.Record) Property {
	p.ID, _, _ = neo4j.GetRecordValue[int64](record, "id")
	p.Address, _, _ = neo4j.GetRecordValue[string](record, "address")
	return p
}

func SearchProperties(c *gin.Context) {
	query, _ := c.GetQuery("q")
	var dbQuery = "MATCH (n:Properties) WHERE n.address =~ '(?i)^.*" + query + ".*$' RETURN n.id as id, n.address as address;"

	res, err := database.ExecuteQuery(c, dbQuery, nil)
	if err != nil {
		return
	}

	properties := make([]PropertyResult, len(res.Records))

	for i, record := range res.Records {
		log.Println(record)
		properties[i] = PropertyResult{Property{}.Reconstruct(record)}
	}

	c.IndentedJSON(http.StatusOK, properties)
}
