package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
	"net/http"
	"snoop-server/pkg/database"
)

type LandlordResult struct {
	Landlord `json:"landlord"`
}

type Landlord struct {
	ID      string `json:"id"`
	Name    string `json:"name,omitempty"`
	Email   string `json:"email,omitempty"`
	Contact string `json:"contact,omitempty"`
}

func (l Landlord) Reconstruct(record *neo4j.Record) Landlord {
	l.ID, _, _ = neo4j.GetRecordValue[string](record, "id")
	l.Name, _, _ = neo4j.GetRecordValue[string](record, "name")
	l.Email, _, _ = neo4j.GetRecordValue[string](record, "email")
	l.Contact, _, _ = neo4j.GetRecordValue[string](record, "contact")
	return l
}

func SearchLandlords(c *gin.Context) {
	searchTerm, _ := c.GetQuery("q")
	var query = "MATCH (n:Landlords) WHERE n.name =~ '(?i)^.*" + searchTerm + ".*$' RETURN n.id as id, n.name as name, n.email as email, n.contact as contact;"

	res, err := database.ExecuteQuery(c, database.Query{String: query}, nil)
	if err != nil {
		return
	}

	landlords := make([]LandlordResult, len(res.Records))

	for i, record := range res.Records {
		landlords[i] = LandlordResult{Landlord{}.Reconstruct(record)}
	}

	c.IndentedJSON(http.StatusOK, landlords)
}

func GetLandlord(c *gin.Context) {
	id := c.Param("id")
	var query = "MATCH (n:Landlords) WHERE n.id = '" + id + "' RETURN n.id as id, n.name as name, n.email as email, n.contact as contact;"

	res, err := database.ExecuteQuery(c, database.Query{String: query}, nil)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.IndentedJSON(http.StatusOK, Landlord{}.Reconstruct(res.Records[0]))
	}
}

func AddLandlord(c *gin.Context) {

	l := new(Landlord)
	err := c.Bind(l)

	if err != nil {
		return
	}

	var query = "CREATE (l:Landlords {id: randomUUID(), name: '" + l.Name + "', email: '" + l.Email + "', contact: '" + l.Contact + "'}) RETURN l;"
	log.Println(query)

	res, err := database.ExecuteQuery(c, database.Query{String: query, Type: "WRITE"}, nil)
	log.Println(res)
	if err != nil {
		log.Println(err)
		return
	}

	c.IndentedJSON(http.StatusCreated, res)
}
