package services

// services/ PACKAGE **********************************************************************************************
/* The services/ package stores all the Business Logic, hence the methods that carry out operations and
   modifications to data/data structures while being completely decoupled from HTTP Requests and Methods. */

/* IMPORTANT NOTES ----------------------------------------------------------------------------------------------*/
/* 1. UserService and Register Method
	- The UserService Struct has only one method (non-static) Register that checks email and password input by
	  the client, convert the password into a hash using the bcrypt algorithm and stores it together with the
	  email in the corresponding DB Table.
   2. service/ and repository/ methods <-- IMPORTANT!!
   	- IMPORTANT!! The service/ package must contain all methods mirroring each single method that is defined in The
      / package! */

// 1. IMPORT PACKAGES *********************************************************************************************

/* Besides the external packages, we also need to import the necessary internal packages defined in the project */
import (
	/* INTERNAL Packages */
	"bookapi/internal/models"
	"bookapi/internal/repositories"
	"bookapi/internal/security"

	/* EXTERNAL Packages */
	"errors"
	"strings"
)

// 2. GO STRUCTS and UTILITY VARIABLES ****************************************************************************

/* STRUCT */
type UserService struct {
	Repo *repositories.UserRepository
}

/* STRUCT BUILDER */
func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{Repo: repo}
}

// 3. BUSINESS LOGIC METHODS **************************************************************************************

/* REGISTER User ------------------------------------------------------------------------------------------------*/
/* Method Mirroring STATIC HTTP Handler for POST /register */
func (s *UserService) Register(req models.RegisterRequest) (models.User, error) {
	/* 1. Extract email and textual password from the input RegisterRequest Go Struct */
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.TrimSpace(req.Password)

	/* 2. Check values - if empty return Empty user struct + error object */
	if req.Email == "" || req.Password == "" {
		return models.User{}, errors.New("Email and password are required")
	}
	/* 3. Get User matching email from DB Table + Error Handling */
	existing, err := s.Repo.FindByEmail(req.Email)
	/*...if error occured, return it with null user object */
	if err != nil {
		return models.User{}, err
	}
	/*...if mathing User exists, return error warning the client that email is already registered */
	if existing != nil {
		return models.User{}, errors.New("Email is already registered")
	}
	/*...in case the input email doesn't exist in the DB Table yet...*/

	/* 4. Generate Hash from Password + Error Handling */
	hashed, err := security.HashPassword(req.Password)
	if err != nil {
		return models.User{}, errors.New("Could not hash password")
	}

	/* 5. Build new User Go Struct with input email and generated HASH of corresponding password */
	user := models.User{
		Email:    req.Email,
		Password: hashed,
	}

	/* 6. Add the built user to the DB Table */
	return s.Repo.Create(user)
}

/* FIND USER BY EMAIL -----------------------------------------------------------------------------------------*/
/* Method Mirroring STATIC HTTP Handler for GET /register */
func (s *UserService) FindByEmail(email string) (*models.User, error) {
	/* 1. Call the Repo Method and get the user item + error object returned */
	user, err := s.Repo.FindByEmail(email)
	/* 2. Error Handling on both user and err obejcts */
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("User not found.")
	}
	/* 3. Return the found user object and null error */
	return user, nil

}

/* FIND ALL USERS --------------------------------------------------------------------------------------------*/
/* Method Mirroring STATIC HTTP Handler for GET /admin/users */
func (s *UserService) FindAll() ([]models.User, error) {
	/* 1. Call the Repo Method and return the list of users from the Database */
	return s.Repo.FindAll()
}
