package handlers

// handlers/ PACKAGE **********************************************************************************************
/* The handlers/ package stores all the HTTP Method Handlers keeping the HTTP logic separate from
   the other packages. */

/* IMPORTANT NOTES ----------------------------------------------------------------------------------------------*/
/* 1. Naming Variables
	- All Go Structs (including their fields!!), Data STructures, Variables....coming from the modules/ package
  	  MUST have the first letter CAPITAL to be accessible in other packages!
  	  In Go the difference between PUBLIC and PRIVATE variables is defined as follows:
		- CAPITAL first letter -> PUBLIC variable
		- LOWER CASE first letter -> PRIVATE variable
   2. RETURN Keyword after Response Helper Functions
	- Whenever we use a Response Helper Function in our code (i.e. WriteJSON(..), WriteError(..), WriteSafeError(..))
	  it has always to be followed by the RETURN keyword!!....otherwise Golang will move on executing the rest of
	  the code!!...and that is not what we want!!
   3. PUT METHOD Particular Features
	- The PUT method handler requires both w and r inputs on top of the id input. This is because, in the handler,
	  we need to work both on the Request (decoding the Body JSON from the HTTP Request) as well as on the Response
	  (WriteJSON Helper Function). GET /books/{id} and DELETE /books/{id}, instead, don't need to work on the HTTP
	  Request since IT HAS NO OBJECT IN ITS BODY! (Only the {id} is the information passed.)
   4. Accessibility of RESPONSE HELPER FUNCTIONS
	- If we want to allow the Response Helper Functions to get used in whatever package of our project (i.e. not
	  only in the handlers/ package where they are defined), we need to name them with the first letter to be a
	  CAPITAL letter: i.e. - writeJSON(..) -> WriteJSON(..)
*/

/* 1. IMPORT PACKAGES *********************************************************************************************
******************************************************************************************************************/
import (
	/* INTERNAL Packages */

	"bookapi/internal/middleware"
	"bookapi/internal/models"
	"bookapi/internal/services"
	"bookapi/internal/utils"

	/* EXTERNAL Packages */
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5" /*													>>>>>>>>> CHI Router <<<<<<<<*/
)

/* 2. GO STRUCTS and UTILITY METHODS  ******************************************************************************
******************************************************************************************************************/

/* Main Struct */
type BookHandler struct {
	Service services.BookService
}

/* Constructor */
func NewBookHandler(service services.BookService) *BookHandler {
	return &BookHandler{Service: service}
}

/* Register All Routes */
func (h *BookHandler) RegisterRoutes(r chi.Router) {
	r.Route("/books", func(r chi.Router) {
		/* STATIC Routes */
		r.Get("/", h.GetBooks)
		r.Post("/", h.PostBook)
		r.With(middleware.AllowRoles("admin")).Post("/transfer", h.TransferPages) /*>>>>>> ROLE-BASED AUTH <<<<<<*/
		/* DYNAMIC Routes */
		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.GetBookByID)
			r.Group(func(r chi.Router) {
				r.Use(middleware.EnforceOwnership("id", /*					   >>>>>> OWNERSHIP-BASED AUTH <<<<<<*/
					func(r *http.Request, id int) (int, error) { return h.Service.GetOwnerID(id) }))
				r.Put("/", h.PutBook)
				r.With(middleware.AllowRoles("admin")).Delete("/", h.DeleteBook) /*>> ROLE+OWNERSHIP-BASED AUTH <<*/
			})
		})
	})
}

/* 3. HTTP REQUEST HANDLERS  ***************************************************************************************
*******************************************************************************************************************/

/* STATIC HTTP Request Handlers ------------------------------------------------------------------------------------
------------------------------------------------------------------------------------------------------------------*/

/* GET /books Handler --------------------------------------------------------------------------------------------*/
/* >>>>>> SWAGGER <<<<<<< */
// @Summary Get all books
// @Description Returns all books stored in the database
// @Tags books
// @Produce json
// @Success 200 {array} models.Book
// @Failure 500 {object} models.ErrorResponse
// @Router /books [get]
func (h *BookHandler) GetBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.Service.ListBooks()
	if err != nil {
		utils.WriteSafeError(w, http.StatusInternalServerError, "Could Not Fetch Books.")
		return
	}
	utils.WriteJSON(w, http.StatusOK, books, nil)
}

/* POST /books Handler ------------------------------------------------------------------------------------------*/
/* >>>>>> SWAGGER <<<<<<< */
// @Summary Create a new book
// @Description Adds a new book to the database
// @Tags books
// @Accept json
// @Produce json
// @Param book body models.Book true "Book to create"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /books [post]
func (h *BookHandler) PostBook(w http.ResponseWriter, r *http.Request) {
	/* 1. Extract the user ID from the JWT token  + Error Handling via Helper Function */
	userID, ok := r.Context().Value(middleware.UserIDKey).(int) /*						>>>>>> JWT <<<<<<< */
	if !ok {
		utils.WriteSafeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	/* 2. Declare Go Struct to convert JSON from HTTP Request into. */
	var book models.Book

	/* 3. Create Decoder Object */
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	/* 4. Handle Error in Decoding the JSON from the HTTP Request into corresponding Go Struct */
	err := decoder.Decode(&book)
	if err != nil {
		/* Error handled using the Error Response Helper Function */
		utils.WriteError(w, http.StatusBadRequest, err, "Invalid Inputs.")
		return /* <--- NEVER FORGET the RETURN keyword AFTER calling the RESPONSE HELPER FUNCTIONS!! */
	}

	/* IMPORTANT !!
	   Handling Error due to missing values of the JSON Fields from the HTTP Request
	   is carried out by the VALIDATEBOOK Method in the services/ package and that gets executed
	   inside all the methods of the BookService object !! */

	/* 5. Assign the user_id to the book's owner_id field */
	book.OwnerID = userID

	/* 4. Add new Book record in the Database via services/ method. */
	newBook, err := h.Service.CreateBook(book)
	if err != nil {
		/* 5. If an error is returned by the service method,
		warn the client about an Internal Server Error via Helper Function. */
		utils.WriteError(w, http.StatusInternalServerError, err, "Server Error.")
	} else {
		/* 6. Convert Go Struct back to JSON, write it to the Body of the HTTP Response
		and send it to Client. */
		utils.WriteJSON(w, http.StatusCreated, newBook, nil)
	}
}

/* POST /transfer Handler ---------------------------------------------------------------------------------------*/
/* >>>>>> SWAGGER <<<<<<< */
// @Summary Transfer pages between two books
// @Description Move a number of pages from book having id=from_id to book having id=to_id
// @Tags books
// @Accept json
// @Produce json
// @Param transferpages body models.TransferRequest true "Pages transfer data"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 405 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /books/transfer [post]
func (h *BookHandler) TransferPages(w http.ResponseWriter, r *http.Request) {
	/* 1. Allow only POST HTTP Method for /transfer End Point. */
	if r.Method != http.MethodPost {
		/* If the Http Method is different than POST, send back an error message using the Helper Function */
		utils.WriteSafeError(w, http.StatusMethodNotAllowed, "Method not allowed.")
		return
	}

	/* 2. Convert the JSON Body of the HTTP Request into a TransferRequest Go Struct + Error Handling */
	var req models.TransferRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err, "Invalid Inputs.")
		return
	}

	/* 3. Check Values of JSON fields from the Body of the HTTP Request + Error Handling */
	if req.FromID <= 0 || req.ToID <= 0 || req.Pages <= 0 {
		utils.WriteSafeError(w, http.StatusBadRequest, "Missing/Invalid JSON Field values.")
		return
	}

	/* 4. EXECUTE the TRANSACTION  - Executes multiple SQL Queries in one single unit of work/function  */
	err = h.Service.TransferPages(req)

	/* 5. Check any error due to failure of Transaction and handle it with helper function */
	if err != nil {
		utils.WriteSafeError(w, http.StatusInternalServerError, "Transfer failed: "+err.Error())
		return
	}

	/* 6. Return the HTTP Response with HTTP Status Code 200 and
	the Transfer Request object via helper function*/
	utils.WriteJSON(w, http.StatusOK, req, nil)
}

/* DYNAMIC HTTP Request Handlers -----------------------------------------------------------------------------------
------------------------------------------------------------------------------------------------------------------*/

/* GET /books/{id} Handler ---------------------------------------------------------------------------------------*/
/* >>>>>> SWAGGER <<<<<<< */
// @Summary Get book by ID
// @Description Retrieves a book by its ID
// @Tags books
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /books/{id} [get]
func (h *BookHandler) GetBookByID(w http.ResponseWriter, r *http.Request) {
	/* 1. Extract the id using the CHI Router directly from the HTTP Request r 		>>>>>>>>> CHI Router <<<<<<<<*/
	idStr := chi.URLParam(r, "id")
	/* 2. Convert id from string to int + Error Handling */
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteSafeError(w, http.StatusBadRequest, "Invalid id input.")
		return /* <--- NEVER FORGET the RETURN keyword AFTER calling the RESPONSE HELPER FUNCTIONS!! */
	}
	/* 3. Get Book Go Struct and corresponding Error Object based on input ID using the services/ method */
	book, err := h.Service.GetBookByID(id)
	/* 4. Handle possible returned error using the Error Response Helper Function */
	if err != nil {
		utils.WriteError(w, http.StatusNotFound, err, "Book Not Found.")
		return /* <--- NEVER FORGET the RETURN keyword AFTER calling the RESPONSE HELPER FUNCTIONS!! */
	}
	if book == nil {
		utils.WriteSafeError(w, http.StatusNotFound, "Book Not Found.")
		return /* <--- NEVER FORGET the RETURN keyword AFTER calling the RESPONSE HELPER FUNCTIONS!! */
	}
	/* 5. Convert the found Book Go Struct into JSON, write it to the Body of the HTTP Response and send it to
	Client. */
	utils.WriteJSON(w, http.StatusOK, book, nil)
}

/* PUT /books/{id} Handler ---------------------------------------------------------------------------------------*/
/* >>>>>> SWAGGER <<<<<<< */
// @Summary Update a book
// @Description Replace an existing book with a new instance
// @Tags books
// @Accept json
// @Produce json
// @Param book body models.Book true "Updated Book"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /books/{id} [put]
func (h *BookHandler) PutBook(w http.ResponseWriter, r *http.Request) {
	/* 1. Extract the id using the CHI Router directly from the HTTP Request r 		>>>>>>>>> CHI Router <<<<<<<<*/
	idStr := chi.URLParam(r, "id")
	/* 2. Convert id from string to int + Error Handling */
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteSafeError(w, http.StatusBadRequest, "Invalid id input.")
	}
	/* 3. Declare Go Struct to store the JSON passed in the Body of the HTTP Request */
	var book models.Book
	/* 4. Create the decoder object to convert the JSON into the corresponding Go Struct */
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	/* 5. Convert JSON to Go Struct and handle possible errors via Error Response Helper Function */
	err = decoder.Decode(&book)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err, "Invalid inputs.")
		return /* <--- NEVER FORGET the RETURN keyword AFTER calling the RESPONSE HELPER FUNCTIONS!! */
	}

	/* 6. Check values of JSON Fields and handle possible errors via Error Safe Response Helper Function
	   Carried out inside the services/ method UpdateBook(..) via the private method validateBook(..) */

	/* 7. Look for the book having id matching the input one and, if found, replace it with input book
	   and return the updated book object via the services/ method UpdateBook() . */
	updatedBook, err := h.Service.UpdateBook(id, book)
	/* 8. If error is returned, handle it using the Error Safe Response Helper Function */
	if err != nil {
		utils.WriteSafeError(w, http.StatusNotFound, "Book Not Found.")
		return /* <--- NEVER FORGET the RETURN keyword AFTER calling the RESPONSE HELPER FUNCTIONS!! */
	}

	/* 9. If everything has gone well, return an HTTP Response with HTTP Status 200 and a Body containing the
	   JSON of the updated object using the Success Response Helper Function */
	utils.WriteJSON(w, http.StatusOK, updatedBook, nil)

}

/* DELETE /books/{id} Handler ---------------------------------------------------------------------------------------*/
/* >>>>>> SWAGGER <<<<<<< */
// @Summary Delete book by ID
// @Description Delete a book from the database based on the input ID
// @Tags books
// @Produce json
// @Param id path int true "Book ID"
// @Success 204 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /books/{id} [delete]
func (h *BookHandler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	/* 1. Extract the id using the CHI Router directly from the HTTP Request r 		>>>>>>>>> CHI Router <<<<<<<<*/
	idStr := chi.URLParam(r, "id")
	/* 2. Convert id from string to int + Error Handling */
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.WriteSafeError(w, http.StatusBadRequest, "Invalid id input.")
	}
	/* 3. Delete book by id directly in the database via the services/ method DeleteBook() */
	err = h.Service.DeleteBook(id)
	/* 4. If an error gets returned by the services/ method, that means that the provided id doesn't
	exist in the database. The error gets handled using a Error Safe Response Helper Function */
	if err != nil {
		utils.WriteSafeError(w, http.StatusNotFound, "Book Not Found.")
		return /* <--- NEVER FORGET the RETURN keyword AFTER calling the RESPONSE HELPER FUNCTIONS!! */
	}
	/* 5. If no error has been returned, return an HTTP Status Code 204 (No Content) within an HTTP Response
	having null/empty Body */
	utils.WriteJSON(w, http.StatusNoContent, nil, nil)
}
