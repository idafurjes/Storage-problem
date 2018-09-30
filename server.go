package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

type Promotions struct {
	rootDir string
	switchingMutex	  *sync.Mutex
}

type Api struct {
	server *http.Server
	promotions *Promotions
}

type Record struct {
	Id	string	`json: "id"`
	Price	string	`json: "price"`
	ExpDate string	`json: "expiration_date"`

}

func (p *Promotions) determineFileName(uuid string) string {
	shards := strings.Split(uuid[0:8], "")
	shardPath := path.Join(shards...)
	return path.Join(p.rootDir, "staging", shardPath, uuid[9:])
}

func (p *Promotions) processNewFile(input io.Reader) error {
	// Copy things and process them
	r := csv.NewReader(input)

	for {
		row, err := r.Read();
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// In production, I would validate that all these values
		// are in the format we expect them to be, including that
		// the UUID is structurally valid
		record := Record{
			Id: row[0],
			Price: row[1],
			ExpDate: row[2],
		}

		jsonRecord, err := json.Marshal(record)
		if err != nil {
			return err
		}
		filename := p.determineFileName(record.Id)

		err = os.MkdirAll(filepath.Dir(filename), 0700)
		if err != nil {
			return err
		}

		file, err := os.Create(filename)
		if err != nil {
			return err
		}

		_, err = file.Write(jsonRecord)
		if err != nil {
			return err
		}

		err = file.Close()
		if err != nil {
			return err
		}
	}

	// Once we've build the staging file tree, we'll rename it to be production
	// and delete the old production


	return nil
}

func (p *Promotions) readPromotion(id string) (string, error) {
	// Look up the file and return it
	return "", nil
}

func NewApiServer(address, rootDir string) (*Api, error) {
	promotions := Promotions{rootDir, &sync.Mutex{}}

	serveMux := http.NewServeMux()

	serveMux.HandleFunc("/", http.NotFound)

	serveMux.HandleFunc("/promotions", func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "POST" {
			err := promotions.processNewFile(req.Body)
			if err != nil{
				fmt.Printf("ERROR: %s\n", err)
			}
			defer req.Body.Close()
		} else {
			http.NotFound(w, req)
			return
		}
	})

	serveMux.HandleFunc("/promotions/", func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "GET" {
			fmt.Fprintf(w, "Looking up: %s", req.URL.Path)
		} else {
			http.NotFound(w, req)
			return
		}
	})

	server := http.Server{
		Addr: address,
		Handler: serveMux,
	}
	api := Api{&server, &promotions}
	return &api, nil
}

func main() {
	server, _ := NewApiServer(":8080", "./root");
	server.server.ListenAndServe()
}
