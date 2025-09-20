package main

/* 1. IMPORT PACKAGES *********************************************************************************************
*******************************************************************************************************************/
import (
	/* INTERNAL Packages */
	"bookapi/internal/config"
	"bookapi/internal/router"
	"os"

	/* EXTERNAL Packages */
	"log"
	"net/http"
	_ "net/http/pprof" /* 												>>>>>> PROFILER <<<<<<< */
	"runtime"          /* 												>>>>>> PROFILER <<<<<<< */

	"github.com/joho/godotenv"
)

/* 2. ENTRY POINT *************************************************************************************************
*******************************************************************************************************************/
/* >>>>>> SWAGGER <<<<<<< */
// @title			BookAPI
// @version			1.0
// @description		Sample server for managing books from ancient roman history and computer science.
// @termsOfService  http://example.com/terms

// @contact.name	Giorgio Albieri
// @contact.email	giorgiocarloroberto.albieri@gmail.com

// @license.name	Sapienza
// @license.url		https://opensource.org/licenses/Sapienza

// @host			localhost:8080
// @BasePath		/bookapi/

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	/* 1. Connect OS to .env for extracting Environment Variables + Error Handling */
	/*...if envPath=="" means that the Go App is not running in Docker...hence we can use
	  the .env file that is stored in the folder cmd/api/ (for local testing) */
	envPath := os.Getenv("ENV_PATH")
	if envPath == "" {
		envPath = ".env" // default for local development
	}
	/*...if envPath!="" means that the Go App is dockerized and, therefore, the location
	  of the .env file is specified by envPath, set up in Dockerfile, and specific for Docker */
	err := godotenv.Load(envPath)
	if err != nil {
		log.Fatal(err)
	}

	// 2. Create custom configuration object loading in it relevant Environment Variables + Error Handling.
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// 3. ALLOCATE PROFILER on a SEPARATE PORT 							>>>>>> PROFILER <<<<<<< */
	go func() {
		/* Activate Tracking of Blocking Events */
		runtime.SetBlockProfileRate(1)
		/* Activate Tracking of waits for locks (mutexes) */
		runtime.SetMutexProfileFraction(1)
		/* Print Info Message in the Console Window */
		log.Println("Starting pprof server on %s", cfg.ProfilerPort)
		/* Allocate Server on Port + Error Handling */
		err := http.ListenAndServe(cfg.ProfilerPort, nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	// 4. CREATE NEW HTTP ROUTER
	/* The method router.NewRouter(..) is defined in the router/ package and uses the value of cfg.DBURL to
	   set up the connection to the PostgreSQL Database. */
	r := router.NewRouter(cfg)
	log.Printf("Starting server on %s", cfg.ServerPort)

	// 5. ALLOCATE SERVER ON PORT + ERROR HANDLING
	err = http.ListenAndServe(cfg.ServerPort, r)
	if err != nil {
		log.Fatal(err)
	}

}
