package security

// security/ PACKAGE **********************************************************************************************
/* The security/ package is used to manage authentication, authorization and protection.
   It is used to generate hashes from passwords using the bcrypt algorithm, compare hashes with string passwords
   to grant access as well as generate authentication tokens to manage user sessions using the jwt library. */

/* IMPORTANT NOTES ----------------------------------------------------------------------------------------------*/
/* 1. Scope of jwt.go
	- Provide the methods that generate new tokens (GenerateToken(..)) and that check/decode existing tokens making
	  sure they are correct and they haven't expired yet (ParseToken(..)).
   2. JWT Token
	- A secure string used to identify a user (like a digital ID card) which can be used for login sessions
  	  or API authentication
*/

// 1. IMPORT PACKAGES *********************************************************************************************
import (
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5" /* 												>>>>>> JWT <<<<<<< */
)

/* Method allowing to create a secure token for a user */
func GenerateToken(userID int, userRole string, secret string) (string, error) {
	/* 1. Define the "claims" (i.e. - the inside part) of the Token */
	claims := jwt.MapClaims{
		"user_id":   userID,                                /* Embed the user's id in the token */
		"user_role": userRole,                              /* Embed the user's role in the token */
		"exp":       time.Now().Add(24 * time.Hour).Unix(), /* Set the expiration time to 24 hours from now.*/
		"iat":       time.Now().Unix(),                     /* Set the issued-at time to the current time.*/
	}
	/* 2. Create the token using the secure method HS256 including in it user info and time settings */
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	/* 3. Lock/Sign the Token using the secret key and return it as a string*/
	return token.SignedString([]byte(secret))
}

/* Method allowing to check that whether the token is valid and read the info inside it */
func ParseToken(tokenStr, secret string) (jwt.MapClaims, error) {
	/* 1. Remove empty spaces within the Token string if present */
	tokenStr = strings.ReplaceAll(tokenStr, " ", "")
	/* 2. Try to decode the input Token with the input Key */
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	/* 3. If the Token is broken (err!=nil) or expired (!token.Valid), return an error */
	if err != nil || !token.Valid {
		return nil, err
	}
	/* 4. Try to extract the Claims of the token (the part that holds user info and timestamps)
	   also checking whether they are in the expected format (jwt.MapClaims) */
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, jwt.ErrInvalidKey
	}
	/* 5. If all goes well, return the claims extracted from the Token and a null error */
	return claims, nil

}
