package main

import "os"
import "log"
import "testing"
import "net/http"
import "net/http/httptest"
import ( 
	_ "github.com/lib/pq"
	"database/sql"
)

func TestLocationHandlerResponseOK(t *testing.T) {
	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/location", nil)
	LocationHandler(response, request)
	log.Println("/location response:", response.Code)
	if response.Code != http.StatusNoContent {
		t.Fatalf("Bad Response")
	}
}

func TestLocationHandlerResponseQuery(t *testing.T) {
	response := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "http://localhost:8000/location?lat=44.09491559960329&lon=-123.0965916720434&acc=5&label=foo", nil)
	LocationHandler(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("Bad Response")
	}
	log.Println(response)
}

func TestLocationHandlerrecordlocations(t *testing.T) {
	query := "select label, acc, st_y(geom) lat, st_x(geom) lon from locations where label='foo'"
	result, err := db.Query(query)
	if err != nil {
		t.Fatalf("could not talk to database")
	}
	if result.label != 'foo' { t.Fatalf("label inserted does not match") }
	if result.acc != 5 { t.Fatalf("acc inserted does not match") }
	log.Println(result)
}

func init() {
	log.Println("smalld testing")
	db_connection := os.Getenv("SMALLD_DB_CONNECTION")
	url_base := os.Getenv("SMALLD_URL_BASE")
	options := os.Getenv("SMALLD_OPTIONS") //override command line flags 
	log.Println("SMALLD_DB_CONNECTION:", db_connection)
	log.Println("SMALLD_URL_BASE:", url_base)
	log.Println("SMALLD_OPTIONS:", options)
	var err error
	db, err = sql.Open("postgres", db_connection)
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	} 
	log.Println("connected to database")
}
