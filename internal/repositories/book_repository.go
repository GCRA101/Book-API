package repositories

// repositories/ PACKAGE **********************************************************************************************
/* The repositories/ package is used to store all the objects definitions and all the methods that are used to execute
   SQL Queries on the connected Database for all CRUD Operations (Create, Read, Update, Delete)
   This package is responsible for DATABASE ACCESS LOGIC. */

/* IMPORTANT NOTES ----------------------------------------------------------------------------------------------*/
/* 1. Variables/Objects Accessibility
		- All Go Structs (including their fields!!), Data STructures, Variables....coming from the modules/ package
  		  MUST have the first letter CAPITAL to be accessible in other packages!
  		  In Go the difference between PUBLIC and PRIVATE variables is defined as follows:
			- CAPITAL first letter -> PUBLIC variable
			- LOWER CASE first letter -> PRIVATE variable
   2. BookRepository Interface
		- IMPORTANT!!! All the CRUD Methods that we need to use with the PgBookRepository Go Struct are stored within
		  an interface that is called BookRepository. We have to make sure that PgBookRepository implements the
		  Interface accordingly.
		  Like all programming languages, to allow an object/class to implement an interface, we must make sure
		  that it implements all the methods declared in the interface...and with the correct signature!!
   		  LET'S REMEMBER THAT, in GO, NON-STATIC METHODS (I.E. CLASSES METHODS) ARE DEFINED SPECIFYING A POINTER
		  TO THE CORRESPONDING GO STRUCT / CLASS BEFORE THE NAME OF THE METHOD! SEE THE QUERY CRUD METHODS BELOW !
*/

// 1. IMPORT PACKAGES **********************************************************************************************
import (
	"bookapi/internal/models"
	"database/sql"
	"errors"
)

// 2. GO STRUCTS and UTILITY VARIABLES ********************************************************************************

/* Interface */
type BookRepository interface {
	Create(book models.Book) (models.Book, error)
	FindAll() ([]models.Book, error)
	FindByID(id int) (*models.Book, error)
	Update(id int, book models.Book) (*models.Book, error)
	Delete(id int) error
	TransferPages(req models.TransferRequest) error
	GetOwnerID(bookID int) (int, error)
}

/* Struct */
type PgBookRepository struct {
	DB *sql.DB
}

/* Struct Builder */
func NewBookRepository(db *sql.DB) BookRepository {
	return &PgBookRepository{DB: db}
}

// 3. QUERY CRUD METHODS **********************************************************************************************

/* CREATE - [POST /books HTTP Method] ---------------------------------------------------------------------------*/
func (r *PgBookRepository) Create(book models.Book) (models.Book, error) {
	/* 1. Build the SQL Query */
	query := `INSERT INTO books (title, author, pages, owner_id) VALUES ($1, $2, $3, $4) RETURNING id`
	/* 3. Execute the SQL Query expecting one single row from the DB Table, fill the placeholders
	      in the SQL query with the listed input values and finally read the returned id and
		  store it in book.ID */
	err := r.DB.QueryRow(query, book.Title, book.Author, book.Pages, book.OwnerID).Scan(&book.ID)
	/* 4. Return the udpated book object and any error that might occur. */
	return book, err
}

/* READ ALL - [GET /books HTTP Method] -------------------------------------------------------------------------*/
func (r *PgBookRepository) FindAll() ([]models.Book, error) {
	/* 1. Execute the SQL Query expecting a list of DB Table Rows */
	rows, err := r.DB.Query("SELECT id, title, author, pages FROM books ORDER BY id ASC")
	/* 2. If an error occurs, return null list together with encountered error */
	if err != nil {
		return nil, err
	}
	/* 3. Make sure that the DB Table Rows get CLOSED when the current function
	   finishes in order to avoid locked memory */
	defer rows.Close()
	/* 4. Create an empty list to store the book objects extracted from the DB Table */
	var books []models.Book
	/* 5. Looping through the rows of the DB Table, extract the field values and store
	      them in the corresponding attributes of each new book object that gets then
		  addedd to the books list. */
	for rows.Next() {
		/* Create a new book struct instance */
		var b models.Book
		/* Get data from the DB Table row and assign it to the book object */
		err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Pages)
		/* Return an error if an error occurs in the process. */
		if err != nil {
			return nil, err
		}
		/* Add the built book object to the list */
		books = append(books, b)
	}
	/* 6. Checks if there were any errors while reading the rows. */
	if err := rows.Err(); err != nil {
		return nil, err
	}
	/* 7. Return the list of books and a null error. */
	return books, nil
}

/* TRANSFER - [POST /transfer HTTP Method] -------------------------------------------------------------------------*/
func (r *PgBookRepository) TransferPages(req models.TransferRequest) error {
	/* 1. Start a new DB Transaction using the Go's standard library database/sql  + Error Handling */
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	/* 2. Define anonymous function to run after the function TransferPages finishes */
	defer func() {
		/* If errors/panic occur, ROLLBACK the Transaction */
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			/* If no errors/panic occurs, COMMIT the Transaction */
			err = tx.Commit()
		}
	}()

	/* 3. Execute an SQL Query that subtracts the input fields' value from the book record having id=fromID */
	_, err = tx.Exec(`UPDATE books SET pages = pages - $1 WHERE id = $2`, req.Pages, req.FromID)
	if err != nil {
		/* If an error occurs, stop and send out the error. */
		return err
	}

	/* 4. Execute an SQL Query that adds the input fields' value to the book record having id=toID */
	_, err = tx.Exec(`UPDATE books SET pages = pages + $1 WHERE id = $2`, req.Pages, req.ToID)
	if err != nil {
		/* If an error occurs, stop and send out the error. */
		return err
	}

	/* 5. If everything has worked out well, return null output */
	return nil
}

/* READ BY ID - [GET /books/{id} HTTP Method] ------------------------------------------------------------------*/
func (r *PgBookRepository) FindByID(id int) (*models.Book, error) {
	/* 1. Create a new instance of the Go Struct "Book" */
	var book models.Book
	/* 2. Execute the SQL Query returning one DB Table Row from which we extract the
	   fields values and assign them to the attributes of the Book object. */
	err := r.DB.QueryRow(`SELECT id, title, author, pages FROM books WHERE id = $1`, id).
		Scan(&book.ID, &book.Title, &book.Author, &book.Pages)

	/* 3. If an error has occured but this error is due to the fact that no DB table row
	   satisfies the SQL Query...that's not actually an error, so just return null. */
	if err == sql.ErrNoRows {
		return nil, errors.New("Book Not Found")
	}
	/* 4. If the error is due to some other reason, that's definitely an error so return
	it in the error output of the function. */
	if err != nil {
		return nil, err
	}
	/* 5. Return the found book object and a null error */
	return &book, nil
}

/* UPDATE - [PUT /books/{id} HTTP Method] ---------------------------------------------------------------------*/
func (r *PgBookRepository) Update(id int, book models.Book) (*models.Book, error) {
	/* 1. Build the SQL Query */
	query := `UPDATE books SET title=$1, author=$2, pages=$3 WHERE id=$4`
	/* 2. Execute the SQL Query filling in the placeholders using the DB.Exec method
	      that DOESN'T return ANY ROW as output but rather a RESULT Object that stores
		  information about how many rows were affected by the updated (RowsAffected()). */
	res, err := r.DB.Exec(query, book.Title, book.Author, book.Pages, id)
	/* 3. If the query fails, return nil and an error. */
	if err != nil {
		return nil, err
	}
	/* 4. Get the number of rows affected and whether any error occurred */
	rowsAffected, err := res.RowsAffected()
	/*...if an error occured, return it together with a null book object */
	if err != nil {
		return nil, err
	}
	/*...if no rows were affected, warn the Client that no book has been found. */
	if rowsAffected == 0 {
		return nil, errors.New("Book Not Found.")
	}
	/* 5. Update the id of the input book with the input id */
	book.ID = id
	/* 6. Return updated book object and null error */
	return &book, nil
}

/* DELETE - [DELETE /books/{id} HTTP Method] ------------------------------------------------------------------*/
func (r *PgBookRepository) Delete(id int) error {
	/* 1. Execute SQL Query deleting the record which id matches the input one.
	      The DB.Exec method DOESN'T return ANY ROW as output but rather a RESULT Object that stores
		  information about how many rows were affected by the delete operation (RowsAffected()) */
	res, err := r.DB.Exec(`DELETE FROM books WHERE id = $1`, id)
	/* 2. If an error has occured, return it as output */
	if err != nil {
		return err
	}
	/* 3. Get the number of affected rows from the res object. If we got an error, return it,
	   if no rows have been affected return error "Book Not Found", otherwise just return null */
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return errors.New("Book Not Found.")
	}
	return nil
}

/* GET OWNER ID - [GET /books/{id} HTTP Method] ------------------------------------------------------------------*/
/* This method is specifically created to encapsulate the extraction of the input book's owner id from the Database.
   This method is called exclusively within the OWNERSHIP-BASED Authorization Middleware EnforceOwnership(..) in the
   file middleware/ownership.go. to carry out authorization checks on HTTP Requests */
func (r *PgBookRepository) GetOwnerID(bookID int) (int, error) {
	/* 1. Create int variable to hold the ID of the book's owner */
	var ownerID int
	/* 2. Execute SQL Query extracting the ID of the owner of the book matching the input book ID */
	err := r.DB.QueryRow("SELECT owner_id FROM books WHERE id = $1", bookID).Scan(&ownerID)
	/* 3. Return owner ID and any error */
	return ownerID, err
}
