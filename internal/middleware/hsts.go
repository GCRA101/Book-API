package middleware

// middleware/ PACKAGE ************************************************************************************************
/* The middleware/ package stores all the MIDDLEWARE functions that allow to add functionalities to the HTTP Handlers
   that are defined in the handlers/ package.
   This is achieved using the DECORATOR PATTERN. */

// 1. IMPORT PACKAGES *************************************************************************************************
import (
	"net/http"
)

// 2. CUSTOM http.Handlers ********************************************************************************************

/* HSTS Middleware ----------------------------------------------------------------------------------*/
/*
Middleware sending an HSTS header to enforce HTTPS at te browser level.
- It prevents downgrade attacks (where someone tries to force a user to use HTTP).
- It improves security by ensuring encrypted communication.
- Once the browser sees this header, it will automatically redirect future HTTP requests to HTTPS
	— even before contacting the server.
*/
func HSTS(next http.Handler) http.Handler {
	/* 1. Actual Handler Function that runs for every registered HTTP request. */
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		/* 2. Set up the HSTS Header */
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		/* 3. Continue handling the HTTP Requests with the next registered middleware */
		next.ServeHTTP(w, r)
	})
}
