package middleware

// middleware/ PACKAGE **********************************************************************************************
/* The middleware/ package stores all the MIDDLEWARE functions that allow to add functionalities to the HTTP Handlers
   that are defined in the handlers/ package.
   This is achieved using the DECORATOR PATTERN. */

/* IMPORTANT NOTES ----------------------------------------------------------------------------------------------*/
/* 1. Use of Public and Private Composer Methods
	- It can be good practice to use a public composer method (i.e. Apply) that involves no need from the Client
  	  to know which middlewares to wrap around the core http requests handler, while using a private composer
  	  method (i.e. applyMiddleware) that requires and allows the user to specify the list of middlewares to be
  	  used.
   2. http.HandlerFunc & http.Handler	<<<<< IMPORTANT !!!!
   - It is possible to define both/either custom Handler Functions (http.HandlerFunc) and Handlers (http.Handler).
	 The two are almost equivalent but bear the following tips in mind to be able to choose the best one.const
	 	> Use http.HandlerFunc for simple, functional handlers.
		> Use http.Handler when you need more structure, like maintaining state or using methods on a custom type.
   - The CHI Router can register middlewares using the following two methods ONLY IF THEY ARE http.Handlers!...
	 ...NOT if they are http.HandlerFuncs!!
		> Register GLOBALLY -> r.Use(requestLogger)
		> Register LOCALLY 	-> r.With(requestLogger).Get/Post/Put/Patch/Delete(...)
*/

// 1. IMPORT PACKAGES *********************************************************************************************
import (
	"log"
	"net/http"
	"time"
)

// 2. CUSTOM http.Handlers ****************************************************************************************

/* REQUEST LOGGING Middleware ---------------------------------------------------------------------------------- */
func Logging(next http.Handler) http.Handler { /*				 		  	  	    >>>>>>>>> CHI Router <<<<<<<<*/
	/* 1. Return a new http.Handler that wraps around the input core/base Handler (next) */
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		/* 1. Get the current time and print HTTP Method infos in the Console */
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)
		/* 2. Execute the next/inner http.Handler */
		next.ServeHTTP(w, r)
		/* 3. Get the duration time to handle the HTTP Response and print it in the Console */
		log.Printf("Completed in %v", time.Since(start))
	})
}
