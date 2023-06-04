package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"net/http"
	"snoop-server/pkg/database"
)

type LandlordResult struct {
	Landlord `json:"landlord"`
}

type Landlord struct {
	ID      int64  `json:"id"`
	Name    string `json:"name,omitempty"`
	Email   string `json:"email,omitempty"`
	Contact int64  `json:"contact,omitempty"`
}

func (l Landlord) Reconstruct(record *neo4j.Record) Landlord {
	l.ID, _, _ = neo4j.GetRecordValue[int64](record, "id")
	l.Name, _, _ = neo4j.GetRecordValue[string](record, "name")
	l.Email, _, _ = neo4j.GetRecordValue[string](record, "email")
	l.Contact, _, _ = neo4j.GetRecordValue[int64](record, "contact")
	return l
}

func SearchLandlords(c *gin.Context) {
	query, _ := c.GetQuery("q")
	var dbQuery = "MATCH (n:Landlords) WHERE n.name =~ '(?i)^.*" + query + ".*$' RETURN n.id as id, n.name as name, n.email as email, n.contact as contact;"

	res, err := database.ExecuteQuery(c, dbQuery, nil)
	if err != nil {
		return
	}

	landlords := make([]LandlordResult, len(res.Records))

	for i, record := range res.Records {
		landlords[i] = LandlordResult{Landlord{}.Reconstruct(record)}
	}

	c.IndentedJSON(http.StatusOK, landlords)
}
