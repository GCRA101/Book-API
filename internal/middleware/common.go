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
	"bookapi/internal/config"
	"bookapi/internal/utils"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/rs/cors"
)

// 2. MIDDLEWARES COMPOSERS/WRAPPERS *******************************************************************************

/* PUBLIC "FACADE" Method ----------------------------------------------------------------------------------------*/
func Apply(h http.HandlerFunc) http.HandlerFunc {
	return applyMiddleware(h, requestLoggingMiddleware, recoveryMiddleware, corsMiddleware, userAgentLogMiddleware)
}

/* PRIVATE INNER Method -----------------------------------------------------------------------------------------*/
func applyMiddleware(h http.HandlerFunc, middlewares ...func(http.HandlerFunc) http.HandlerFunc) http.HandlerFunc {
	/* 1. Initialize all Middleware functions one inside the other looping backwards through their list */
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	/* 2. Return the Outer resulting Handler Function */
	return h
}

// 3. MIDDLEWARES *************************************************************************************************

// 3.1 CUSTOM http.HandlerFuncs ***********************************************************************************

/* REQUEST LOGGING Middleware ---------------------------------------------------------------------------------- */
func requestLoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	/* 1. Return a new http.HandlerFunc that wraps around the input core/base Handler (next) */
	return func(w http.ResponseWriter, r *http.Request) {
		/* 1. Get the current time and print HTTP Method infos in the Console */
		startTime := time.Now()
		log.Printf("Started HTTP Request %s %s", r.Method, r.URL.Path)
		/* 2. RUN THE CORE/BASE HTTP.HANDLERFUNC */
		next(w, r)
		/* 3. Get the duration time to handle the HTTP Response and print it in the Console */
		durationTime := time.Since(startTime)
		log.Printf("Completed %s %s in %v \n\n", r.Method, r.URL.Path, durationTime)
	}
}

/* PANIC RECOVERY Middleware ------------------------------------------------------------------------------------*/
func recoveryMiddleware(next http.HandlerFunc) http.HandlerFunc {
	/* 1. Return a new http.HandlerFunc that wraps around the input core/base Handler (next) */
	return func(w http.ResponseWriter, r *http.Request) {
		/* 2. Recover from any panic */
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Recovered from panic: %v", err)
				utils.WriteSafeError(w, http.StatusInternalServerError, "Internal Server Error")
			}
		}()
		/* 3. RUN THE CORE/BASE HTTP.HANDLERFUNC */
		next(w, r)
	}
}

/* CORS Middleware --------------------------------------------------------------------------------------------- */
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	/* 1. Return http.HandlerFunc object wrapping around the input one (next) */
	return func(w http.ResponseWriter, r *http.Request) {
		/* 1. Set Allowed Origins for HTTP Requests - Any of them in this case. */
		w.Header().Set("Access-Control-Allow-Origin", "*")
		/* 2. Set Allowed HTTP Methods for HTTP Requests - Any apart from PATCH in this case */
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		/* 3. Set Allowed Headers for HTTP Requests - Content-Type in this case */
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		/* 4. If HTTP Method OPTIONS is used, return empty HTTP Response */
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return /* <--- NEVER FORGET the RETURN keyword!! */
		}
		/* 5. RUN THE CORE/BASE HTTP.HANDLERFUNC */
		next(w, r)
	}
}

/* USER AGENT LOGGING Middleware ------------------------------------------------------------------------------- */
func userAgentLogMiddleware(next http.HandlerFunc) http.HandlerFunc {
	/* 1. Return a new http.HandlerFunc object wrapping around the input one (next) */
	return func(w http.ResponseWriter, r *http.Request) {
		/* 2. Print the User Agent of the HTTP Request in the Console window */
		log.Printf("User-Agent: %s", r.Header.Get("User-Agent"))
		/* 3. RUN THE CORE/BASE HTTP.HANDLERFUNC */
		next(w, r)
	}
}

/* ACCESS RESTRING TO INTERNAL IPs Middleware ----------------------------------------------------------------- */
func ipWhitelistMiddleware(next http.HandlerFunc) http.HandlerFunc {
	/* 1. Return a new http.HandlerFunc object wrapping around the input one (next) */
	return func(w http.ResponseWriter, r *http.Request) {
		/* 2. Return Error if IP Address is not Internal */
		if r.RemoteAddr != "127.0.0.1:1234" {
			utils.WriteSafeError(w, http.StatusForbidden, "Access Denied")
			return /* <--- NEVER FORGET the RETURN keyword AFTER calling the RESPONSE HELPER FUNCTIONS!! */
		}
		/* 3. RUN THE CORE/BASE HTTP.HANDLERFUNC */
		next(w, r)
	}
}

// 3.2 CUSTOM http.Handlers **************************************************************************************

/* AUTHENTICATION Middleware -----------------------------------------------------------------------------------*/
/*
This middleware acts as a gatekeeper. It checks if the incoming request has the correct Authorization header.
If not, it blocks the request with a 401 error.
If the header is valid, it lets the request proceed to the next handler.
*/
func AuthMiddleware(next http.Handler) http.Handler { /*				 		  >>>>>>>>> CHI Router <<<<<<<<*/
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer Secret" {
			utils.WriteSafeError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}
		/* Execute the next/inner http.Handler */
		next.ServeHTTP(w, r) /* Equivalent to next(w,r) with next http.HandlerFunc !! */
	})
}

/* REQUEST LOGGER Middleware ---------------------------------------------------------------------------------- */
/*
http.Handler version of the http.HandlerFunc requestLoggingMiddleware.
*/
func RequestLogger(next http.Handler) http.Handler { /*				 		  	  >>>>>>>>> CHI Router <<<<<<<<*/
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)
		/* Execute the next/inner http.Handler */
		next.ServeHTTP(w, r) /* Equivalent to next(w,r) with next http.HandlerFunc !! */
		duration := time.Since(start)
		log.Printf("Completed %s in %v", r.URL.Path, duration)
	})
}

/* CORS Middleware --------------------------------------------------------------------------------------------- */
/*
http.Handler version of the http.HandlerFunc corsMiddleware.
*/
func CorsMiddleware(cfg config.Config) func(http.Handler) http.Handler { /* >>>>  CONFIG-DRIVEN CORS SETUP <<<< */
	return func(next http.Handler) http.Handler {
		return cors.New(cors.Options{
			AllowedOrigins: strings.Split(cfg.CorsAllowedOrigins, ","),
			AllowedMethods: strings.Split(cfg.CorsAllowedMethods, ","),
		}).Handler(next)
	}
}
