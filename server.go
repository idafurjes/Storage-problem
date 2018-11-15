package main

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
)

type Database struct {
	db *sql.DB
}

type Api struct {
	server   *http.Server
	database *Database
}

type Record struct {
	Id      string `json: "id"`
	Price   string `json: "price"`
	ExpDate string `json: "expiration_date"`
}

func csvReader(input io.Reader, records chan Record) {
	// Copy things and process them
	r := csv.NewReader(input)

	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		records <- Record{row[0], row[1], row[2]}
	}
	close(records)
}

func (d *Database) insertRecord(records chan Record) {
	for record := range records {
		stmt, err := d.db.Prepare("INSERT INTO promotions (id, price, exp_date) VALUES ($1,$2,$3)")
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()
		stmt.Exec(record.Id, record.Price, record.ExpDate)
	}
}

func (d *Database) processNewFile(input io.Reader) error {
	var wg sync.WaitGroup
	c := make(chan Record)

	for i := 0; i < 90; i++ {
		go d.insertRecord(c)
		wg.Add(1)
		defer wg.Done()
	}
	//defer wg.Done()

	csvReader(input, c)
	wg.Wait()

	return nil
}

func (d *Database) readPromotion(id string) ([]byte, error) {
	row := d.db.QueryRow("SELECT price, exp_date FROM promotions WHERE id = $1;", id)
	r := Record{Id: id}
	err := row.Scan(&r.Price, &r.ExpDate)
	if err == sql.ErrNoRows {
		return []byte{}, nil
	}
	if err != nil {
		return []byte{}, err
	}
	value, err := json.Marshal(r)

	return value, err
}

func NewDatabase() (*Database, error) {
	// In production, we'd use something like a DATABASE_URL environment
	// variable to get this value
	db, err := sql.Open("postgres", "postgres://postgres@localhost/promodb?sslmode=disable")
	if err != nil {
		return nil, err
	}
	return &Database{db}, err
}

func NewApiServer(address, rootDir string) (*Api, error) {
	database, err := NewDatabase()
	if err != nil {
		return nil, err
	}

	serveMux := http.NewServeMux()

	serveMux.HandleFunc("/", http.NotFound)

	serveMux.HandleFunc("/promotions", func(w http.ResponseWriter, req *http.Request) {
		if req.Method == "POST" {
			err := database.processNewFile(req.Body)
			if err != nil {
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
			_id := strings.Split(req.URL.Path, "/")
			id := _id[len(_id)-1]
			value, err := database.readPromotion(id)
			//fmt.Printf("%+v %+v\n", value, err)
			if err != nil {
				fmt.Printf("ERROR: %s\n", err)
			}
			if len(value) == 0 {
				http.NotFound(w, req)
				return
			}
			_, err = w.Write(value)
			if err != nil {
				fmt.Printf("ERROR: %s\n", err)
			}
		} else {
			http.NotFound(w, req)
			return
		}
	})

	server := http.Server{
		Addr:    address,
		Handler: serveMux,
	}
	api := Api{&server, database}
	return &api, nil
}

func main() {
	server, err := NewApiServer(":1321", "./root")
	if err != nil {
		fmt.Printf("ERROR: %s", err)
	}
	server.server.ListenAndServe()
}
