package models

// models/ PACKAGE ************************************************************************************************
/* The models/ package is used to store all the definitions of all objects that are used in the application.
   These includes Go Structs and Utility Variables. */

/* IMPORTANT NOTES ----------------------------------------------------------------------------------------------*/
/* 1. Variables/Objects Accessibility
		- All Go Structs (including their fields!!), Data STructures, Variables....coming from the modules/ package
  		  MUST have the first letter CAPITAL to be accessible in other packages!
  		  In Go the difference between PUBLIC and PRIVATE variables is defined as follows:
			- CAPITAL first letter -> PUBLIC variable
			- LOWER CASE first letter -> PRIVATE variable
   2. Databases vs Data Structures
		- Since, in this case, we're using PostgreSQL Databases to store the data, there's no need to declare any
		  Data Structure here (e.g. Books array) to store the Go Struct Instances. All is handled by the db/ and
		  repositories/ packages.



// 1. IMPORT PACKAGES *********************************************************************************************
/* No need to import any package in this case */

// 2. GO STRUCTS **************************************************************************************************

/* Book */
type Book struct { /* 				>>>>> SWAGGER <<<<< */
	ID      int    `json:"id" example:"1"`
	Title   string `json:"title" example:"The Go Programming Language"` /* 	Title of the book. */
	Author  string `json:"author" example:"Alan Donovan"`               /* 	Name of the author. */
	Pages   int    `json:"pages" example:"380"`                         /* 	Number of pages. */
	OwnerID int    `json:"-" example:"1"`                               // omit from JSON Responses and SWAGGER !
}

/* Transfer Request */
type TransferRequest struct { /* 	>>>>> SWAGGER <<<<< */
	FromID int `json:"from_id" example:"1"` /*Unique ID of the book that provides pages.*/
	ToID   int `json:"to_id" example:"2"`   /*Unique ID of the book that receives pages */
	Pages  int `json:"pages" example:"50"`  /*Number of pages transferred*/
}
