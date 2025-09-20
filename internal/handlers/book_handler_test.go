package handlers

// handlers/ PACKAGE **********************************************************************************************
/* The handlers/ package stores all the HTTP Method Handlers keeping the HTTP logic separate from
   the other packages. */

/* IMPORTANT NOTES ----------------------------------------------------------------------------------------------*/
/* 1. Scope of book_handler_test.go
    - This go file defines unit tests for the RESTful API.
   	  It uses the testing package along with httptest to simulate HTTP requests and responses.
    - It tests the methods POST /books, GET /books, GET /books/{id}, PUT /books/{id} and DELETE /books/{id} but
   	  instead of using the real database or real service, it uses a fake service (called a "mock") to test how
	  the API behaves.
   2. BookService Interface for mocking endpoint testing
	- In order to be able to use the book_handler_test.go file for testing, we need to be able to pass to
	  the BookHandler the mockBookService object. This will make possible to handle http requests without
	  having a server running and a database in place. The mockBookService and the BookService structs must
	  implement a same interface to be accepted as inputs by the BookHandler Struct (service field).
	  Hence the need to create a BookService interface that both the bookService struct and mockBookService
	  struct have to implement (in Go, it's just enough that the signatures of all their methods match with
	  the ones of the interface!)
   3. Registering middleware
    - Important!! Do not forget registering/assigning to the mock router the middleware that we use in the
	  actual router that we want to test !!!
   4. Decode nested JSON
    - If the JSON we want to decode has a simple structure, json.NewDecoder(..).Decode(&..) will do the job.
      If the JSON has a nested structure (e.g.: "data" and "meta" fields) and we are interested only in the json
	  that is stored in one of its fields (i.e. "data") the decoding process gets more complex.
	  - In such case, 1)decode the nested JSON using json.NewDecoder(..).Decode(&..), 2) get out of it the field
	    "data", 3) marshal the data back to JSON and then 4) unmarshal it into models.Book - see Helper function
		decodeNestedJSON(...) below.
*/

// 1. IMPORT PACKAGES *********************************************************************************************
import (
	/* INTERNAL Packages */
	"bookapi/internal/config"
	"bookapi/internal/middleware"
	"bookapi/internal/models"
	"bookapi/internal/security"

	/* EXTERNAL Packages */
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware" /* 							>>>>>> CHI Router <<<<<<< */
)

// 2. MOCK SERVICE - GO STRUCTS & UTILITY METHODS  ****************************************************************

/* STRUCT */
/* Fake version of the real book service and that lets control what happens when the API tries to handle the HTTP
   methods selected for testing */
type mockBookService struct {
	/* Function for creating a new Book [POST /books] */
	CreateFunc func(models.Book) (models.Book, error)
	/* Function for getting all Books [GET /books] */
	ListFunc func() ([]models.Book, error)
	/* Function for getting one Book by id [GET /books/{id}] */
	GetFunc func(int) (*models.Book, error)
	/* Function for transferring pages between two books [POST /books/transfer] */
	TransferFunc func(req models.TransferRequest) error
	/* Function for updating one book by id [PUT /books/{id}] */
	UpdateFunc func(id int, updated models.Book) (*models.Book, error)
	/* Function for deleting one book by id [DELETE /books/{id}] */
	DeleteFunc func(id int) error
	/* Function for returning the owner id of the input book id */
	GetOwnerFunc func(int) (int, error)
}

/* NON-STATIC METHODS of mockBookService */
/* ListBooks() - "When someone asks for books, use the fake function I gave you
   (i.e. m.ListFunc())." */
func (m *mockBookService) ListBooks() ([]models.Book, error) {
	return m.ListFunc()
}

/*
CreateBook() - "When someone asks to create a new book, use the fake function I gave you (i.e. m.CreateFunc()).
(i.e. m.CreateFunc())."
*/
func (m *mockBookService) CreateBook(book models.Book) (models.Book, error) {
	return m.CreateFunc(book)
}

/*
GetBookByIDtBooks() - "When someone asks to get a book by id, use the fake function I gave you.
(i.e. m.GetFunc())."
*/
func (m *mockBookService) GetBookByID(id int) (*models.Book, error) {
	return m.GetFunc(id)
}

/*
TransferPages() - "When someone asks to transfer pages, use the fake function I gave you.
(i.e. m.TransferFunc())."
*/
func (m *mockBookService) TransferPages(req models.TransferRequest) error {
	return m.TransferFunc(req)
}

/*
UpdateBook() - "When someone asks to update a book, use the fake function I gave you.
(i.e. m.UpdateFunc())."
*/
func (m *mockBookService) UpdateBook(id int, updated models.Book) (*models.Book, error) {
	return m.UpdateFunc(id, updated)
}

/*
DeleteBook() - "When someone asks to delete a book, use the fake function I gave you.
(i.e. m.DeleteFunc())."
*/
func (m *mockBookService) DeleteBook(id int) error {
	return m.DeleteFunc(id)
}

/*
DeleteBook() - "When someone asks to delete a book, use the fake function I gave you.
(i.e. m.GetOwnerFunc())."
*/
func (m *mockBookService) GetOwnerID(bookID int) (int, error) {
	return m.GetOwnerFunc(bookID)
}

// 3. ROUTER - HANDLERS REGISTRATION  *****************************************************************************

/* Set up a test version of the router */
func setupTestRouter(service *mockBookService) http.Handler {
	/* 1. Create BookHandler passing the mockBookService via BookService Interface */
	handler := &BookHandler{Service: service}
	/* 2. Load the Configuration object containing main environment variables */
	cfg := config.Load()
	/* 3. Create the Chi Router */
	r := chi.NewRouter()
	/* 4. Register the main Middleware */
	r.Use(middleware.Logging, chimiddleware.Recoverer, middleware.JWTAuth(cfg.JWTSecret))
	/* 5. Register Handlers to Endpoints */
	r.Get("/books", handler.GetBooks)
	r.Post("/books", handler.PostBook)
	r.Post("/books/transfer", handler.TransferPages)
	r.Get("/books/{id}", handler.GetBookByID)
	r.Put("/books/{id}", handler.PutBook)
	r.Delete("/books/{id}", handler.DeleteBook)
	/* 6. Return router */
	return r
}

// 4. HTTP TEST HELPERS  ******************************************************************************************

/* TESTER for POST /books ---------------------------------------------------------------------------------------*/
func TestCreateBookEndpoint(t *testing.T) {

	/* 1. Set the test service createBook function and assign it to the mockBookService. */
	service := &mockBookService{
		/* The fake createBook method is designed to return always the input book with updated id and null error.*/
		CreateFunc: func(b models.Book) (models.Book, error) {
			b.ID = 42
			return b, nil
		},
	}

	/* 2. Set up the Test Router */
	router := setupTestRouter(service)

	/* 3. Create a fake HTTP Request to simulate sending a book to the server -- >> same as in POSTMAN! << */
	/* 3.1 Set up the Body */
	body := `{"title":"The Go Programming Language", "author": "Alan Donovan", "pages": 380}`
	/* 3.2 Set up the HTTP Method, Route and Body */
	req := httptest.NewRequest(http.MethodPost, "/books", strings.NewReader(body))
	/* 3.3 Set up the Headers - Content-Type */
	req.Header.Set("Content-Type", "application/json")
	/* 3.4 Set up the Headers - Authorization */
	token, err := security.GenerateToken(1, "user", config.Load().JWTSecret)
	if err != nil {
		t.Fatalf("Error in Generating the Authorization Token")
	}
	req.Header.Set("Authorization", "Bearer "+token)

	/* 4. Create a fake HTTP Response Recorder */
	rec := httptest.NewRecorder()

	/* 5. Send the Fake HTTP Request and Record the Fake HTTP Response */
	router.ServeHTTP(rec, req)

	/* 6. Check the Headers of the fake HTTP Response*/
	validateHeaders(t, rec)

	/* 7. Check the HTTP Response Status Code */
	if rec.Code != http.StatusCreated {
		/* ...if not 201, return Error message */
		t.Fatalf("expected status 201, got %d", rec.Code)
	}
	/* 8. Check the JSON Body of the HTTP Response */
	var result models.Book
	/* 8.1 Check the Decoding Process via Helper Function */
	result = decodeNestedJSON[models.Book](t, rec.Body)
	/* 8.2 Check the Content */
	if result.ID != 42 {
		/* ...if content is not as expected, return Error message */
		t.Errorf("Expected ID 42, got %d", result.ID)
	}
}

/* TESTER for GET /books  ---------------------------------------------------------------------------------------*/
func TestListBooksEndpoint(t *testing.T) {

	/* 1. Set the test service ListBooks function and assign it to the mockBookService. */
	service := &mockBookService{
		ListFunc: func() ([]models.Book, error) {
			/* The fake ListBooks method is designed to return a list of books made by one single book only */
			return []models.Book{
				{ID: 1, Title: "Go in Action", Author: "William Kennedy", Pages: 320},
			}, nil
		},
	}

	/* 2. Set up the Test Router */
	router := setupTestRouter(service)

	/* 3. Create a fake HTTP Request to simulate requesting books from the server -- >> same as in POSTMAN! << */
	req := httptest.NewRequest(http.MethodGet, "/books", nil)
	/* Set up the Headers - Authorization */
	token, err := security.GenerateToken(1, "user", config.Load().JWTSecret)
	if err != nil {
		t.Fatalf("Error in Generating the Authorization Token")
	}
	req.Header.Set("Authorization", "Bearer "+token)

	/* 4. Create a fake HTTP Response Recorder */
	rec := httptest.NewRecorder()

	/* 5. Send the Fake HTTP Request and Record the Fake HTTP Response */
	router.ServeHTTP(rec, req)

	/* 6. Check the Headers of the fake HTTP Response*/
	validateHeaders(t, rec)

	/* 7. Check the HTTP Response Status Code */
	if rec.Code != http.StatusOK {
		/* ...if not 200, return Error message */
		t.Fatalf("Expected Status 200, got %d", rec.Code)
	}

	/* 8. Check JSON Body of HTTP Response */
	var books []models.Book
	/* 8.1 Check the Decoding Process via Helper Function */
	books = decodeNestedJSON[[]models.Book](t, rec.Body)
	/* 8.2 Check the Content */
	if len(books) != 1 || books[0].Title != "Go in Action" {
		/* ...if content is not as expected, return Error message */
		t.Errorf("Unexpected book list: %+v", books)
	}
}

/* TESTER for POST /transfer  -----------------------------------------------------------------------------------*/
func TestTransferPagesEndPoint(t *testing.T) {
	/* 1. Set the test service TransferPages function and assign it to the mockBookService. */
	service := &mockBookService{
		/* The fake TransferPages method is designed to return a null error. */
		TransferFunc: func(req models.TransferRequest) error {
			return nil
		},
	}

	/* 2. Set up the Test Router */
	router := setupTestRouter(service)

	/* 3. Create a fake HTTP Request to simulate transfering pages on the server -- >> same as in POSTMAN! << */
	/* 3.1 Set up the Body */
	body := `{"from_id": 1, "to_id": 2, "pages": 120}`
	/* 3.2 Set up the HTTP Method, Route and Body */
	req := httptest.NewRequest(http.MethodPost, "/books/transfer", strings.NewReader(body))
	/* 3.3 Set up the Headers - Content-Type */
	req.Header.Set("Content-Type", "application/json")
	/* 3.4 Set up the Headers - Authorization */
	token, err := security.GenerateToken(1, "user", config.Load().JWTSecret)
	if err != nil {
		t.Fatalf("Error in Generating the Authorization Token")
	}
	req.Header.Set("Authorization", "Bearer "+token)

	/* 4. Create a fake HTTP Response Recorder */
	rec := httptest.NewRecorder()

	/* 5. Send the Fake HTTP Request and Record the Fake HTTP Response */
	router.ServeHTTP(rec, req)

	/* 6. Check the Headers of the fake HTTP Response*/
	validateHeaders(t, rec)

	/* 7. Check HTTP Response Status Code */
	if rec.Code != http.StatusOK {
		/* ...if not 404, return Error message */
		t.Fatalf("Expected 200 Not Found, got %d", rec.Code)
	}
	/* 8. Check the JSON Body of the HTTP Response */
	var result models.TransferRequest
	/* 8.1 Check the Decoding Process via Helper Function */
	result = decodeNestedJSON[models.TransferRequest](t, rec.Body)
	/* 8.2 Check the Content */
	if result.FromID != 1 {
		/* ...if content is not as expected, return Error message */
		t.Errorf("Expected from_ID 1, got %d", result.FromID)
	}
}

/* TESTER for GET /books/{id} -----------------------------------------------------------------------------------*/
func TestGetBookByIDEndPoint_NotFound(t *testing.T) {

	/* 1. Set the test service GetBookByID function and assign it to the mockBookService. */
	service := &mockBookService{
		/* The fake GetBookByID method is designed to return null book object and null error
		   whatever is the input book ID we're looking for. */
		GetFunc: func(id int) (*models.Book, error) {
			return nil, nil
		},
	}

	/* 2. Set up the Test Router */
	router := setupTestRouter(service)

	/* 3. Create a fake HTTP Request to simulate sending a book to the server -- >> same as in POSTMAN! << */
	req := httptest.NewRequest(http.MethodGet, "/books/999", nil)
	/* Set up the Headers - Authorization */
	token, err := security.GenerateToken(1, "user", config.Load().JWTSecret)
	if err != nil {
		t.Fatalf("Error in Generating the Authorization Token")
	}
	req.Header.Set("Authorization", "Bearer "+token)

	/* 4. Create a fake HTTP Response Recorder */
	rec := httptest.NewRecorder()

	/* 5. Send the Fake HTTP Request and Record the Fake HTTP Response */
	router.ServeHTTP(rec, req)

	/* 6. Check the Headers of the fake HTTP Response*/
	validateHeaders(t, rec)

	/* 7. Check HTTP Response Status Code */
	if rec.Code != http.StatusNotFound {
		/* ...if not 404, return Error message */
		t.Fatalf("Expected 404 Not Found, got %d", rec.Code)
	}
}

/* TESTER for PUT /books/{id} -----------------------------------------------------------------------------------*/
func TestPutBookByIDEndPoint(t *testing.T) {

	/* 1. Set the test service PutBook function and assign it to the mockBookService. */
	service := &mockBookService{
		/* The fake PutBook method is designed to return a book object and null error
		   whatever is the input book ID we're looking for. */
		UpdateFunc: func(id int, updated models.Book) (*models.Book, error) {
			updated.ID = id
			return &updated, nil
		},
	}

	/* 2. Set up the Test Router */
	router := setupTestRouter(service)

	/* 3. Create a fake HTTP Request to simulate updating a book on the server -- >> same as in POSTMAN! << */
	/* 3.1 Set up the Body */
	body := `{"title":"The Go Programming Language", "author": "Alan Donovan", "pages": 380}`
	/* 3.2 Set up the HTTP Method, Route and Body */
	req := httptest.NewRequest(http.MethodPut, "/books/15", strings.NewReader(body))
	/* 3.3 Set up the Headers - Content-Type */
	req.Header.Set("Content-Type", "application/json")
	/* 3.4 Set up the Headers - Authorization */
	token, err := security.GenerateToken(1, "user", config.Load().JWTSecret)
	if err != nil {
		t.Fatalf("Error in Generating the Authorization Token")
	}
	req.Header.Set("Authorization", "Bearer "+token)

	/* 4. Create a fake HTTP Response Recorder */
	rec := httptest.NewRecorder()

	/* 5. Send the Fake HTTP Request and Record the Fake HTTP Response */
	router.ServeHTTP(rec, req)

	/* 6. Check the Headers of the fake HTTP Response*/
	validateHeaders(t, rec)

	/* 7. Check the HTTP Response Status Code */
	if rec.Code != http.StatusOK {
		/* ...if not 200, return Error message */
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	/* 8. Check the JSON Body of the HTTP Response */
	var result models.Book
	/* 8.1 Check the Decoding Process via Helper Function */
	result = decodeNestedJSON[models.Book](t, rec.Body)
	/* 8.2 Check the Content */
	if result.ID != 15 {
		/* ...if content is not as expected, return Error message */
		t.Errorf("Expected ID 15, got %d", result.ID)
	}

}

/* TESTER for DELETE /books/{id} --------------------------------------------------------------------------------*/
func TestDeleteBookEndpoint(t *testing.T) {

	/* 1. Set the test service deleteBook function and assign it to the mockBookService. */
	service := &mockBookService{
		/* The fake deleteBook method is designed to return always a null error.*/
		DeleteFunc: func(id int) error {
			return nil
		},
	}

	/* 2. Set up the Test Router */
	router := setupTestRouter(service)

	/* 3. Create a fake HTTP Request to simulate deleting a book from the server -- >> same as in POSTMAN! << */
	/* 3.1 Set up the HTTP Method, Route and Body */
	req := httptest.NewRequest(http.MethodDelete, "/books/13", nil)
	/* 3.2 Set up the Headers - Authorization */
	token, err := security.GenerateToken(1, "user", config.Load().JWTSecret)
	if err != nil {
		t.Fatalf("Error in Generating the Authorization Token")
	}
	req.Header.Set("Authorization", "Bearer "+token)

	/* 4. Create a fake HTTP Response Recorder */
	rec := httptest.NewRecorder()

	/* 5. Send the Fake HTTP Request and Record the Fake HTTP Response */
	router.ServeHTTP(rec, req)

	/* 6. Check the Headers of the fake HTTP Response*/
	validateHeaders(t, rec)

	/* 7. Check the HTTP Response Status Code */
	if rec.Code != http.StatusNoContent {
		/* ...if not 204, return Error message */
		t.Fatalf("expected status 204, got %d", rec.Code)
	}
}

// 5. TEST HELPER FUNCTIONS ***************************************************************************************

/* Decoding JSON ------------------------------------------------------------------------------------------------*/
/* Helper function encapsulating conversion of JSON into a Go object */
func decodeJSON[T any](t *testing.T, body *bytes.Buffer) T {
	var v T
	err := json.NewDecoder(body).Decode(&v)
	if err != nil {
		/* ...if error occurs, return Error message */
		t.Fatalf("Failed to decode JSON: %v", err)
	}
	return v
}

/* Decoding JSON [data+meta] ------------------------------------------------------------------------------------*/
/* Helper function encapsulating conversion of JSON into a Go object */
func decodeNestedJSON[T any](t *testing.T, body *bytes.Buffer) T {
	var resp map[string]interface{}
	if err := json.NewDecoder(body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode JSON: %v", err)
	}
	data, ok := resp["data"]
	if !ok {
		t.Fatalf("Expected 'data' field in response")
	}
	// Marshal the data back to JSON, then unmarshal into models.Book
	dataBytes, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Failed to marshal 'data': %v", err)
	}
	var v T
	if err := json.Unmarshal(dataBytes, &v); err != nil {
		t.Fatalf("Failed to unmarshal 'data' into input dataType: %v", err)
	}
	return any(v).(T)
}

/* Validating HEADERS and CONTENT-TYPE --------------------------------------------------------------------------*/
/* Helper function checking if the response has the correct Content-Type header. */
func validateHeaders(t *testing.T, rec *httptest.ResponseRecorder) {
	/* 1. Get the value of the Content-Type header of the recorded HTTP Response */
	ct := rec.Header().Get("Content-Type")
	/* 2. Check value + send error message */
	if ct != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", ct)
		return
	}
}
