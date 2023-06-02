package main

import (
	"context"
	"fmt"
	"go_gin_neo4j/pkg/model"
	"log"
	"net/http"
	"os"
)

func main() {
	ctx := context.Background()
	configuration := model.ParseConfiguration()
	driver, err := configuration.NewDriver()
	if err != nil {
		log.Fatal(err)
	}
	defer model.UnsafeClose(ctx, driver)
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/", model.DefaultHandler)
	serveMux.HandleFunc("/search", model.SearchHandlerFunc(ctx, driver, configuration.Database))
	serveMux.HandleFunc("/movie/vote/", model.VoteInMovieHandlerFunc(ctx, driver, configuration.Database))
	serveMux.HandleFunc("/movie/", model.MovieHandlerFunc(ctx, driver, configuration.Database))
	serveMux.HandleFunc("/graph", model.GraphHandler(ctx, driver, configuration.Database))

	var port string
	var found bool
	if port, found = os.LookupEnv("PORT"); !found {
		port = "8080"
	}
	fmt.Printf("Running on port %s, database is at %s\n", port, configuration.Url)
	panic(http.ListenAndServe(":"+port, serveMux))
}
