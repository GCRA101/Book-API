package security

// security/ PACKAGE **********************************************************************************************
/* The security/ package is used to manage authentication, authorization and protection.
   It is used to generate hashes from passwords using the bcrypt algorithm, compare hashes with string passwords
   to grant access as well as generate authentication tokens to manage user sessions using the jwt library. */

/* IMPORTANT NOTES ----------------------------------------------------------------------------------------------*/
/* 1. Bcrypt Hashing Algorithm
- bcrypt is a cryptographic hashing algorithm designed for password hashing.
  It’s slow by design to resist brute-force attacks.
*/

// 1. IMPORT PACKAGES *******************************************************************************************
import (
	"golang.org/x/crypto/bcrypt"
)

// 2. HASHING METHODS *******************************************************************************************

/* Convert String Password to Hash */
func HashPassword(password string) (string, error) {
	/* 1. Convert the input string password into a Hash via bcrypt algorithm + return any error.
	DefaultCost=10...a cost factor value that is a good balance between security and performance.
	The cost factor is a measure of the computational complexity of the algorithm. */
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	/* 2. Convert byte slice hash to string and return it together with any error encountered */
	return string(hash), err
}

/* Compare Hash with String Password */
func CheckPasswordHash(password, hash string) bool {
	/* 1. Convert hash and password to byte slices and compares the two returning an error if not successful */
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	/* 2. Return True if match, False if not */
	return err == nil
}
