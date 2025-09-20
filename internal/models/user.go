package models

// models/ PACKAGE ************************************************************************************************
/* The models/ package is used to store all the definitions of all objects that are used in the application.
   These includes Go Structs and Utility Variables. */

/* IMPORTANT NOTES ----------------------------------------------------------------------------------------------*/
/* 1. Omitting Go Struct Fields from JSON
		- When a field/property of a Go Struct has to be kept secret and, hence, not included in the encoded JSON
		  object returned to the client via HTTP Response (e.g. Password) the json tag we need to use is as follows
		  	-> `json:"-"`


// 1. IMPORT PACKAGES *********************************************************************************************
/* No need to import any package in this case */

// 2. GO STRUCTS **************************************************************************************************

/* User */
type User struct { /* 				>>>>> SWAGGER <<<<< */
	ID       int    `json:"id" example:"1"`                       /* User's unique id */
	Role     string `json:"role" example:"user"`                  /* User's role for authorization */
	Email    string `json:"email" example:"john.golan@gmail.com"` /* User's email address */
	Password string `json:"-" example:"secretwordXXX`             // omit from JSON Responses!!
}

/* Register Request */
type RegisterRequest struct { /* 	>>>>> SWAGGER <<<<< */
	Email    string `json:"email" example:"john.golan@gmail.com"` /* User's email address */
	Password string `json:"password" example:"secretwordXXX`      /* User's login password */
}
