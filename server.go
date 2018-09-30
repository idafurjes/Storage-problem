package main

import (
	"fmt"
	"net/http"
)

type Api struct {
	server *http.Server
	currentRecordsDir string
	stagingRecordsDir string
}

type Record struct {
	Id	string	`json: "id"`
	Price	string	`json: "price"`
	ExpDate string 	`json: "expiration_date"`

}

func NewApiServer(address, currentRecordsDir, stagingRecordsDir string) (*Api, error) {
	serveMux := http.NewServeMux();

	serveMux.HandleFunc("/promotions", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})

	server := http.Server{
		Addr: address,
		Handler: serveMux,
	}
	api := Api{&server, currentRecordsDir, stagingRecordsDir}
	return &api, nil
}

func main() {
	server, _ := NewApiServer(":8080", "./cur", "stg");
	server.server.ListenAndServe()
}
