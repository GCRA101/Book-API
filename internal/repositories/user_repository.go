package repositories

// repositories/ PACKAGE **********************************************************************************************
/* The repositories/ package is used to store all the objects definitions and all the methods that are used to execute
   SQL Queries on the connected Database for all CRUD Operations (Create, Read, Update, Delete)
   This package is responsible for DATABASE ACCESS LOGIC. */

/* IMPORTANT NOTES ----------------------------------------------------------------------------------------------*/
/* 1. UserRepository
		- Repository class/go_struct populated with methods that allow to 1) store, in the connected DB Table, an input
		  instance of User struct; and 2) find a user in the DB Table based on input email.
   2. Static vs Non-Static Methods
		- func (r *UserRepository) Create(user models.User) (models.User, error)
			-> NON-STATIC Method. It belongs to and gets executed by instances of UserRepository Struct
		- func Create(user models.User) (models.User, error)
			-> STATIC Method. It can be executed without any instance of UserRepository.

*/

// 1. IMPORT PACKAGES *************************************************************************************************
import (
	"bookapi/internal/models"
	"database/sql"
)

// 2. GO STRUCTS and UTILITY VARIABLES ********************************************************************************

/* STRUCT */
type UserRepository struct {
	DB *sql.DB
}

/* STRUCT BUILDER */
func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// 3. QUERY CRUD METHODS **********************************************************************************************

/* CREATE - [POST /register HTTP Method] ---------------------------------------------------------------------------*/
func (r *UserRepository) Create(user models.User) (models.User, error) {
	/* 1. Build SQL Query string adding user object in DB Table */
	query := `INSERT INTO users (email, password) VALUES ($1, $2) RETURNING id`
	/* 2. Execute Query passing user email and password in the placeholders and assigning id of db table row to the
	the input user object. If any error occurs, the error gets returned in err */
	err := r.DB.QueryRow(query, user.Email, user.Password).Scan(&user.ID)
	/* 3. Return input user object with updated id based on assignment in DB table + any error */
	return user, err
}

/* FIND BY EMAIL - [GET /register HTTP Method] ---------------------------------------------------------------------*/
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	/* 1. Declare a new User Go Struct to hold values extracted from the DB Table*/
	var user models.User
	/* 2. Execute SQL Query looking for user matching input email, return any encoutered error and populate the
	   fields of the Go Struct with the corresponding table row values. */
	err := r.DB.QueryRow(`SELECT id, role, email, password FROM users WHERE email = $1`, email).
		Scan(&user.ID, &user.Role, &user.Email, &user.Password)
	/* 3. If the encountered error is due to no rows returned by the query....that's not an error but just an
	      indication that there's no user in the database associated with the input email....so return null
		  user object and null error...*/
	if err == sql.ErrNoRows {
		return nil, nil
	}
	/* 4. If the encountered error is different, return the error as it is...*/
	if err != nil {
		return nil, err
	}
	/* 5. If no error has been encountered, return pointer to found user object + null error */
	return &user, nil
}

/* FIND ALL - [GET /admin/users HTTP Method] ---------------------------------------------------------------------*/
func (r *UserRepository) FindAll() ([]models.User, error) {
	/* 1. Execute the SQL Query expecting a list of DB Table Rows */
	rows, err := r.DB.Query("SELECT id, role, email, password FROM users ORDER BY id ASC")
	/* 2. If an error occurs, return null list together with encountered error */
	if err != nil {
		return nil, err
	}
	/* 3. Make sure that the DB Table Rows get CLOSED when the current function
	   finishes in order to avoid locked memory */
	defer rows.Close()
	/* 4. Create an empty list to store the user objects extracted from the DB Table */
	var users []models.User
	/* 5. Looping through the rows of the DB Table, extract the field values and store
	      them in the corresponding attributes of each new user object that gets then
		  addedd to the users list. */
	for rows.Next() {
		/* Create a new book struct instance */
		var user models.User
		/* Get data from the DB Table row and assign it to the book object */
		err := rows.Scan(&user.ID, &user.Role, &user.Email, &user.Password)
		/* Return an error if an error occurs in the process. */
		if err != nil {
			return nil, err
		}
		/* Add the built user object to the list */
		users = append(users, user)
	}
	/* 6. Checks if there were any errors while reading the rows. */
	if err := rows.Err(); err != nil {
		return nil, err
	}
	/* 7. Return the list of books and a null error. */
	return users, nil
}
