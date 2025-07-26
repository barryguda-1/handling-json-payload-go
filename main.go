package main

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type Book struct {
	ID     int    `json:"id,omitempty"` //exclude from response if empty
	Title  string `json:"title"`
	Author string `json:"author"`
	Pages  int    `json:"pages"`
}

type BookListResponse struct {
	Data  []Book `json:"data"`
	Total int    `json:"total"`
}

var books = []Book{
	{ID: 1, Title: "Go in Action", Author: "William Kennedy", Pages: 300},
	{ID: 2, Title: "Go Web Programming", Author: "SauSheong Chang", Pages: 350},
	{ID: 3, Title: "Go Programming Language", Author: "Alln Donovan", Pages: 400},
}

// Sample POST request from client using BOOK struct
func createBookHandler(w http.ResponseWriter, r *http.Request) {
	var book Book
	//Decode json body post
	//err := json.NewDecoder(r.Body).Decode(&book)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields() //disallow extra fields
	err := decoder.Decode(&book)
	if err != nil {
		//handle error incase of an invalid payload
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	//Handling malformed/incomplete json
	if book.Title == "" || book.Author == "" || book.Pages <= 0 {
		http.Error(w, "Missing or invalid fields", http.StatusBadRequest)
		return
	}

}

func getBookHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Sample Json Responses /books?all=true/
	all := r.URL.Query().Get("all")
	if all == "true" {
		response := BookListResponse{
			Data:  books,
			Total: len(books),
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	// Sample Json Responses /books?id=1,2,3/
	// Best method to use /books/{id}
	bookIDStr := r.URL.Query().Get("id")

	if bookIDStr == "" {
		http.Error(w, "Missing book ID", http.StatusBadRequest)
		return
	}

	bookID, _ := strconv.Atoi(bookIDStr)

	for _, book := range books {
		if book.ID == bookID {
			json.NewEncoder(w).Encode(book)
			return
		}
	}

}

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getBookHandler(w, r)
	case http.MethodPost:
		createBookHandler(w, r)
	}
}

func main() {
	http.HandleFunc("/books", handler)
	http.ListenAndServe(":8080", nil)
}
