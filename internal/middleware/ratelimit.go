package middleware

// middleware/ PACKAGE ************************************************************************************************
/* The middleware/ package stores all the MIDDLEWARE functions that allow to add functionalities to the HTTP Handlers
   that are defined in the handlers/ package.
   This is achieved using the DECORATOR PATTERN. */

// 1. IMPORT PACKAGES *************************************************************************************************
import (
	/* INTERNAL Packages */
	"bookapi/internal/utils"
	/* EXTERNAL Packages */
	"net/http"
	"sync"
	"time"

	/* Allows to connect to a Redis Database */
	"github.com/redis/go-redis/v9"
	/* Allows to implement HTTP request RATE LIMIT */
	"github.com/ulule/limiter/v3"
	/* Adapter for standard HTTP */
	chimiddleware "github.com/ulule/limiter/v3/drivers/middleware/stdlib"
	/* Allows to store Rate Limit data in Redis DB */
	redisstore "github.com/ulule/limiter/v3/drivers/store/redis"
)

// 2. GO STRUCTS and UTILITY VARIABLES  *******************************************************************************

/* Requests Tracker - Go Struct */
type rateLimitEntry struct {
	LastSeen time.Time
	Count    int
}

/* Global Variable */
var (
	/* Map storing rate limit info for each IP address */
	visitors = make(map[string]*rateLimitEntry)
	/* Mutex (lock) making sure only one goroutine accesses the map at a time */
	mu sync.Mutex
)

/* Constants */
const (
	/* Time Window to limit rate */
	limitWindow = 1 * time.Minute
	/* Max number of requests allowed per IP within the limit Window */
	requestCap = 60
)

// 3. CUSTOM http.Handlers ********************************************************************************************

/* TESTING RATE-LIMIT Middleware ----------------------------------------------------------------------------------*/
/*
Middleware designed to limit the Rate of HTTP Requests to all Endpoints assigned with it.
Function returning another function — a middleware — that wraps around HTTP handlers to control
how often they can be called by a user based on their IP Address.
*/
func RateLimit(next http.Handler) http.Handler {
	/* 1. Actual Handler Function that runs for every registered HTTP request. */
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		/* 2. Get the IP address of the client sending the HTTP request */
		ip := r.RemoteAddr
		/* 3. Lock the visitors map to access it safely */
		mu.Lock()
		/* 4. Check if the IP address already has an entry in the map */
		entry, exists := visitors[ip]
		/* 5A. ...if the IP address isn't recorded in the map yet or the last request
		   has been done a while ago (beyond the limit window)... */
		if !exists || time.Since(entry.LastSeen) > limitWindow {
			/* ...create a new entry with count=1...*/
			visitors[ip] = &rateLimitEntry{LastSeen: time.Now(), Count: 1}
			/*...unlock the visitors map...*/
			mu.Unlock()
			/*...move on handling the HTTP request...*/
			next.ServeHTTP(w, r)
			return
		}
		/* 5B. ...if teh IP address has already been recorded in the map...*/
		/*...increase the requests' counter...*/
		entry.Count++
		/*...update the last seen time...*/
		entry.LastSeen = time.Now()
		/*...unlock the map...*/
		mu.Unlock()

		/* 6. If the requests count exceeds the cap/limit...*/
		if entry.Count > requestCap {
			/*...send back 429 Error via Helper Function */
			utils.WriteSafeError(w, http.StatusTooManyRequests, "Rate limit exceeded")
			return /* <--- NEVER FORGET the RETURN keyword AFTER calling the RESPONSE HELPER FUNCTIONS!! */
		}
		/* 7. If the request is within the limit, pass it to the next handler. */
		next.ServeHTTP(w, r)
	})
}

/* PRODUCTION RATE-LIMIT Middleware ----------------------------------------------------------------------------------*/
/*
Middleware designed to limit the Rate of HTTP Requests to all Endpoints assigned with it.
Function returning another function — a middleware — that wraps around HTTP handlers to control
how often they can be called.
*/
func ProductionRateLimit() func(http.Handler) http.Handler {
	/* 1. Create a Redis Client (i.e. Connection) that connects to Redis running at port 6379 */
	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379"})
	/* 2. Set up Storage System */
	store, err := redisstore.NewStoreWithOptions(rdb, limiter.StoreOptions{})
	if err != nil {
		panic(err)
	}
	/* 3. Set up Rate Limits */
	rate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  100,
	}
	/* 4. Create the limiter object that enforces the rate limit */
	limiterInstance := limiter.New(store, rate)
	/* 5. Wrap the limiter in a middleware that can be used with standard HTTP handlers */
	middleware := chimiddleware.NewMiddleware(limiterInstance)
	/* 6. Return the middleware function to protect routes */
	return middleware.Handler
}
