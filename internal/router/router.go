package router

// router/ PACKAGE ************************************************************************************************
/* The router/ package is responsible for tying everything together: routes, middleware,
   services, repositoreis and dependencies. It sets up and returns the HTTP router of our application.
   A proper initialization layer. */

/* IMPORTANT NOTES ----------------------------------------------------------------------------------------------*/
/* 1. Variables/Objects Accessibility
		- All Go Structs (including their fields!!), Data STructures, Variables....coming from the modules/ package
  		  MUST have the first letter CAPITAL to be accessible in other packages!
  		  In Go the difference between PUBLIC and PRIVATE variables is defined as follows:
			- CAPITAL first letter -> PUBLIC variable
			- LOWER CASE first letter -> PRIVATE variable
   2. Use of _"github.com/lib/pb
		- The PostgreSQL driver gets imported anonymously
		- It is needed for sql.Open to work with PostgreSQL).
*/

// 1. IMPORT PACKAGES *********************************************************************************************
import (
	bookConfig "bookapi/internal/config"
	"bookapi/internal/handlers"
	"bookapi/internal/middleware"
	"bookapi/internal/repositories"
	"bookapi/internal/services"
	"fmt"
	"time"

	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"                          /* 						    >>>>>> CHI Router <<<<< */
	chimiddleware "github.com/go-chi/chi/v5/middleware" /* 							>>>>>> CHI Router <<<<< */
	_ "github.com/lib/pq"

	_ "bookapi/docs" /* 						 					 				>>>>>> SWAGGER <<<<<<< */

	httpSwagger "github.com/swaggo/http-swagger/v2" /* 						 		>>>>>> SWAGGER <<<<<<< */
)

func NewRouter(cfg bookConfig.Config) http.Handler {
	/* 1. Open a connection to the PostgreSQL database using the URL from the config + Error Handling */
	db, err := initPostgres(cfg.DBURL)
	if err != nil {
		log.Fatal("Failed to connect to DB: ", err)
	}

	/* 2. Create Repository instances using the database connection. */
	userRepo := repositories.NewUserRepository(db)
	bookRepo := repositories.NewBookRepository(db)
	/* 3. Create Service instances using the repositories. */
	userService := services.NewUserService(userRepo)
	bookService := services.NewBookService(bookRepo)
	/* 4. Create Handler instances using the services. */
	userHandler := handlers.NewUserHandler(userService)
	adminHandler := handlers.NewAdminHandler(userService)
	authHandler := handlers.NewAuthHandler(userService, cfg.JWTSecret)
	bookHandler := handlers.NewBookHandler(bookService)

	/* 5. Create new CHI Router. */
	r := chi.NewRouter()
	/* 6. Apply Middleware */
	r.Use(middleware.Logging, chimiddleware.Recoverer) /*   >>>> Custom and CHI-Built-In Middleware <<<<< */
	r.Use(middleware.HSTS)                             /* 					  >>>> HTTPS Middleware <<<<< */
	if cfg.ServerPort == "6379" {
		r.Use(middleware.ProductionRateLimit()) /* 			 			 >>>> RATE LIMIT Middleware <<<<< */
	} else {
		r.Use(middleware.RateLimit) /* 			 						 >>>> RATE LIMIT Middleware <<<<< */
	}
	/* 7. Register all the Routes to the corresponding Handlers. */
	userHandler.RegisterRoutes(r)
	authHandler.RegisterRoutes(r)
	adminHandler.RegisterRoutes(r.With(middleware.JWTAuth(cfg.JWTSecret)))
	bookHandler.RegisterRoutes(r.With(middleware.JWTAuth(cfg.JWTSecret)))

	/* 8. Register the Swagger Route to its imported Handler */
	r.Group(func(r chi.Router) {
		//r.Use(middleware.JWTAuth(cfg.JWTSecret))
		r.Get("/swagger/*", httpSwagger.WrapHandler)
	})

	/* 9. Return the configured router so it can be used in main.go. */
	return r
}

// 2. DB UTILITY METHODS ******************************************************************************************

/* Initialize Connection to PostgreSQL Database */
func initPostgres(connStr string) (*sql.DB, error) {

	/* 1. Create the Connection to the DB Engine (PostgreSQL) + Error Handling */
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("Could not open DB: %w", err)
	}

	/* 2. Verify presence/status of the connection */
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("Could not connect to DB: %w", err)
	}

	/* 3. Configure the Connection Pooling */

	// Set maximum number of open connections
	db.SetMaxOpenConns(10)
	// Set maximum number of idle connections
	db.SetMaxIdleConns(5)
	// Set the maximum lifetime of an open connection
	db.SetConnMaxLifetime(time.Hour)
	// Set the maximum lifetime of an idle connection
	db.SetConnMaxIdleTime(30 * time.Minute)

	/* 4. Send Info Message to user via Console Window */
	log.Println("Connnected to PostgreSQL successfully.")

	/* 5. Return Pointer to Database Connection and Error object */
	return db, nil
}
