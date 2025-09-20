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

/* Success Response */
type SuccessResponse struct { /* 	>>>>> SWAGGER <<<<< */
	Data interface{} `json:"data" example"{id:1, title:"The Fractal Brain Theory", author:"Tsang", pages:"500}"`
	Meta interface{} `json:"meta"`
}

/* Error Response */
type ErrorResponse struct { /* 	>>>>> SWAGGER <<<<< */
	Error   string `json:"error"`                             /* Stringified Error Object */
	Message string `json:"message" example:"Book not found."` /* Customized Error Message */
}
