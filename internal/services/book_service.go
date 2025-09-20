package services

// services/ PACKAGE **********************************************************************************************
/* The services/ package stores all the Business Logic, hence the methods that carry out operations and
   modifications to data/data structures while being completely decoupled from HTTP Requests and Methods. */

/* IMPORTANT NOTES ----------------------------------------------------------------------------------------------*/
/* 1. Naming Variables
		- All Go Structs (including their fields!!), Data STructures, Variables....coming from the modules/ package
  		  MUST have the first letter CAPITAL to be accessible in other packages!
  		  In Go the difference between PUBLIC and PRIVATE variables is defined as follows:
			- CAPITAL first letter -> PUBLIC variable
			- LOWER CASE first letter -> PRIVATE variable
   2. Functions Return Type
		- ALWAYS remember to specify the DATA TYPE of the RETURNED object of the Function!
   3. Naming Methods
		- We might be tempted to name methods the same way of the corresponding HTTP Methods (e.g. putBook,
		  postBook...) but...WE HAVE TO REMEMBER that the services/ package is created to determine SEPARATION
		  of CONCERNS and DECOUPLE the BUSINESS LOGIC from HTTP METHODS HANDLING!!
		  Therefore, it is better practice to NAME METHODS using NATURAL LANGUAGE (e.g. updateBook, createBook..)
	4. Double Methods Return Outputs
		- It's good practice to make the methods return a couple of outputs made from the object they have to
		  return + an error object. In case no error occurs the first is NOT null while the second is NULL.
		  Viceversa otherwise.
	5. BookService Interface for mocking endpoint testing
		- In order to be able to use the book_handler_test.go file for testing, we need to be able to pass to
		  the BookHandler the mockBookService object. This will make possible to handle http requests without
		  having a server running and a database in place. The mockBookService and the BookService structs must
		  implement a same interface to be accepted as inputs of the BookHandler Struct (service field).
		  Hence the need to create a BookService interface that both the bookService struct and mockBookService
		  struct have to implement (in Go, it's just enough that the signatures of all their methods match with
		  the ones of the interface!)
*/

// 1. IMPORT PACKAGES *********************************************************************************************

/* Besides the external packages, we also need to import the necessary internal packages defined in the project */
import (
	/* INTERNAL Packages */
	"bookapi/internal/models"
	"bookapi/internal/repositories"

	/* EXTERNAL Packages */
	"errors"
)

// 2. GO STRUCTS and UTILITY VARIABLES ****************************************************************************

/* INTERFACE */
/* Important!!: In order to be able to use the book_handler_test.go file for testing, we need to be able to pass to
   the BookHandler the mockBookService object. This will make possible to handle http requests without having a
   server running and a database in place. The mockBookService and the BookService structs must implement a same
   interface to be accepted as inputs of the BookHandler Struct (service field).
   Hence the need to create a BookService interface that both the bookService struct and mockBookService struct
   have to implement (in Go, it's just enough that the signatures of all their methods match with the ones of the
   interface!) */
type BookService interface {
	ListBooks() ([]models.Book, error)
	GetBookByID(id int) (*models.Book, error)
	CreateBook(book models.Book) (models.Book, error)
	TransferPages(req models.TransferRequest) error
	UpdateBook(id int, updated models.Book) (*models.Book, error)
	DeleteBook(id int) error
	GetOwnerID(bookID int) (int, error)
}

/* STRUCT */
/* Such struct is part of the service layer, which connects business logic with the repository (database) layer. */
type bookService struct {
	Repo repositories.BookRepository
}

/* STRUCT BUILDER */
func NewBookService(repo repositories.BookRepository) BookService {
	return &bookService{Repo: repo}
}

// 3. BUSINESS LOGIC METHODS **************************************************************************************

/* GET AllBooks -------------------------------------------------------------------------------------------------*/
/* Method Mirroring STATIC HTTP Handler for GET /books */
func (s *bookService) ListBooks() ([]models.Book, error) {
	/* 1. Call the Repo Method and return the list of books from the Database */
	return s.Repo.FindAll()
}

/* GET Book -----------------------------------------------------------------------------------------------------*/
/* Method Mirroring DYNAMIC HTTP Handler for GET /books/{id} */
func (s *bookService) GetBookByID(id int) (*models.Book, error) {
	/* 1. Call the Repo Method and get the book item + error object returned */
	book, err := s.Repo.FindByID(id)
	/* 2. Error Handling on both book and err obejcts */
	if err != nil {
		return nil, err
	}
	if book == nil {
		return nil, errors.New("Book not found.")
	}
	/* 3. Return the found book object and null error */
	return book, nil
}

/* CREATE Book ---------------------------------------------------------------------------------------------------*/
/* Method Mirroring STATIC HTTP Handler for POST /books */
func (s *bookService) CreateBook(book models.Book) (models.Book, error) {
	/* 1. Check JSON Fields' values are not empty/not acceptable + Error Handling */
	err := s.validateBook(book)
	if err != nil {
		return models.Book{}, err
	}
	/* 2. Call the Repo Method and return the created book from the database + any error */
	return s.Repo.Create(book)
}

/* TRANSFER pages ------------------------------------------------------------------------------------------------*/
/* Method Mirroring STATIC HTTP Handler for POST /transfer */
func (s *bookService) TransferPages(req models.TransferRequest) error {
	/* 1. Check JSON Fields' values are not empty/not acceptable + Error Handling */
	err := s.validateTransferRequest(req)
	if err != nil {
		return err
	}
	/* 2. Call the Repo Method and return the created book from the database + any error */
	err = s.Repo.TransferPages(req)
	if err != nil {
		return err
	}
	return nil
}

/* UPDATE Book --------------------------------------------------------------------------------------------------*/
/* Method Mirroring DYNAMIC HTTP Handler for PUT /books/{id} */
func (s *bookService) UpdateBook(id int, updated models.Book) (*models.Book, error) {
	/* 1. Check JSON Fields' values are not empty/not acceptable + Error Handling */
	err := s.validateBook(updated)
	if err != nil {
		return nil, err
	}
	/* 2. Call the Repo Method and return the updated book from the database + any error */
	return s.Repo.Update(id, updated)
}

/* DELETE Book --------------------------------------------------------------------------------------------------*/
/* Method Mirroring DYNAMIC HTTP Handler for DELETE /books/{id} */
func (s *bookService) DeleteBook(id int) error {
	/* 1. Call the Repo Method and return any error */
	return s.Repo.Delete(id)
}

/* GET OwnerID --------------------------------------------------------------------------------------------------*/
/* Method Encapsulating Utility method for getting ID of book's owner */
func (s *bookService) GetOwnerID(bookID int) (int, error) {
	/* 1. Call the Repo Method and get the owner id + error object returned */
	ownerID, err := s.Repo.GetOwnerID(bookID)
	/* 2. Error Handling on both owner id and error objects */
	if err != nil {
		return 0, err
	}
	if ownerID == 0 {
		return 0, errors.New("Book not found.")
	}
	/* 3. Return the found owner id and null error */
	return ownerID, nil
}

/* Utility Method validateBook ----------------------------------------------------------------------------------*/
/* Method keeping the checks on the Body JSON Field's values out of the handlers and database code */
func (s *bookService) validateBook(book models.Book) error {
	/* If Book objects has empty title/author or negative pages, return an error...*/
	if book.Title == "" {
		return errors.New("Title is required")
	}
	if book.Author == "" {
		return errors.New("Author is required")
	}
	if book.Pages <= 0 {
		return errors.New("Pages must be greater than 0")
	}
	/*...otherwise return null */
	return nil
}

/* Utility Method transferRequest ------------------------------------------------------------------------------*/
/* Method keeping the checks on the Body JSON Field's values out of the handlers and database code */
func (s *bookService) validateTransferRequest(req models.TransferRequest) error {
	/* If Book objects has empty title/author or negative pages, return an error...*/
	if req.FromID <= 0 {
		return errors.New("Sender Book ID is invalid")
	}
	if req.ToID <= 0 {
		return errors.New("Receiver Book ID is invalid")
	}
	if req.Pages < 0 {
		return errors.New("Pages must be greater or equal to 0")
	}
	/*...otherwise return null */
	return nil
}
