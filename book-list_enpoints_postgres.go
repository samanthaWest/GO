package main

import {
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "database/sql"
	_ "github.com/lib/pq"
	_ "github.com/subosito/gotenv"

	"github.com/gorilla/mux"
}

type Book struct {
	ID int `json:"id"`
	Title string `json:"title"`
	Author string `json:"author"`
	Year string `json:"year"`
}

var books []Book

func init() {
	gotenv.Load()
}

func logFatal(err error) {
	if err != null {
		log.Fatal(err)
	}
}

func main() {

	pgUrl, err := pq.ParseURL(os.Getenv("ELEPHANTSQL_URL"))
	logFatal(err)

	log.Println(pgUrl)

	err = db.Ping()
	logFatal(err)

	router := mux.NewRouter()

	// Endpoints
	router.HandleFunc("/book", getBooks).Methods("GET")
	router.HandleFunc("/book/{id}", getBook).Methods("GET")
	router.HandleFunc("/books", addBook).Methods("POST")
	router.HandleFunc("/book", updateBook).Methods("PUT")
	router.HandleFunc("/books/{id}", removeBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))
}

func getBooks(w http.ResponseWriter, r *http.Request) {
	var book Book
	books = []Book{}

	rows, error := dbQuery("select * from books")
	logFatal(err)

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
		logFatal(err)

		books = append(book, book)
	}

	json.NewEncoder(w).Encode(books)
}

func getBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	params := mux.vars(r)

	rows := dbQueryRow("select * from books where id=$1", params["id"])
	
	err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Year)
	logFatal(err)

	json.NewEncoder(w).Encode(book)
}

func addBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	var bookId int

	json.NewEncoder(r.Body).Decode(&book)

	db.QueryRow("insert into books (title,author,year) values($1, $2, $3) RETURNING id;",
			book.Title, book.Author, book.Year).Scan(&bookID)

	logFatal(err)

	json.newEncoder(w).Encode(bookID)
}

func updateBook(w http.ResponseWriter, r *http.Request) {
	var book Book
	
	json.Decoder(r.Body).Decode(&book)

	result, err := db.Exec("update books set title=$1, author=$2, year=$3 where id=$4 RETURNING id",
		&book.Title, &book.Author, &bookYear, &book.ID)

	rowsUpdated, err := result.RowsAffected()
	logFatal(err)

	json.NewEncoder(w).Encode(rowsUpdated)
}

func removeBook(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	result, err := db.Exec("delete from books where id = $1", params["id"])
	logFatal(err)

	rowsDeleted, err := result.RowsAffected()
	logFatal(err)

	json.NewEncoder(w).Encode(rowsDeleted)
}
