package config

// config/ PACKAGE **********************************************************************************************
/* The config/ package is used to load configuration values from environment variables and provide default values
   for them in case they are not set */

/* IMPORTANT NOTES ----------------------------------------------------------------------------------------------*/
/* 1. JWT_SECRET
   - ALWAYS STORE the JWT_SECRET as an ENVIRONMENT VARIABLE!! Never expose it to the CLIENT!!!
*/

// 1. IMPORT PACKAGES *******************************************************************************************

/* The os package from the Go standard library allows to access environment variables via os.LookupEnv! */
import (
	"errors"
	"fmt"
	"os"
)

// 2. GO STRUCTS and CONSTANTS **********************************************************************************

/* Config Struct holding key environment variables' values extracted using the os package method LookupEnv */
type Config struct {
	ServerPort         string // The port the server will listen on (e.g. :8080)
	ProfilerPort       string // The port the pprof server will listen on (e.g. 6060) 		>>>> PROFILER <<<<
	DBURL              string // The connection string for the database.
	JWTSecret          string // The Secret used to generate Authentication Tokens			>>>>>> JWT <<<<<<<
	CorsAllowedOrigins string // The List of allowed origins for CORS
	CorsAllowedMethods string // The List of allowed methods for CORS
}

// 3. UTILITY METHODS *******************************************************************************************

/* Load Method - Gets values from environment variables and assigns them to Config Go struct object */
func Load() (Config, error) {

	/* 1. Get the Server Port + Error Handling */
	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		return Config{}, errors.New("SERVER PORT missing in .env file")
	}

	/* 2. Get the DB Connection String + Error Handling */
	dbUrl, err := buildDBConnString()
	if err != nil {
		return Config{}, err
	}

	/* 3. Get the JWT Secret + Error Handling */
	jwtSecret := os.Getenv("JWT_SECRET") /* 				>>>>>> JWT <<<<<<< */
	if jwtSecret == "" {
		return Config{}, errors.New("JWT_SECRET missing in .env file")
	}

	/* 4. Get the CORS Allowed Origins + Error Handling */
	allowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		return Config{}, errors.New("CORS_ALLOWED_ORIGINS missing in .env file")
	}

	return Config{
		/* Get the value of the SERVER_PORT environment variable, or use :8080 as a default.*/
		ServerPort: serverPort,
		/* Set the value of the Profiler Port */
		ProfilerPort: ":6060",
		/* Set the value of the Database URL */
		DBURL: dbUrl,
		/* Get the value of the JWT_SECRET environment variable, or use the default value */
		JWTSecret: jwtSecret, /* 							>>>>>> JWT <<<<<<< */
		/* Get the value of the CORS_ALLOWED_ORIGINS environment variable, or use the default value */
		CorsAllowedOrigins: allowedOrigins,
		/* Get the value of the CORS_ALLOWED_METHODS environment variable, or use the default value */
		CorsAllowedMethods: getEnv("CORS_ALLOWED_METHODS", "GET, POST, PUT, DELETE"),
	}, nil
}

/* getEnv Method - Returns values from environment variables if available, otherwise returns default values */
func getEnv(key, fallback string) string {
	/* If the variable exists (ok == true), it returns the value... */
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	/* ...otherwise, it returns the fallback value provided. */
	return fallback
}

/*
buildDBConnString Method - Returns DB connection String getting env variables from .env file.
If something goes wrong, it returns an error.
*/
func buildDBConnString() (string, error) {
	/* 1. Try first getting the DB URL directly from the Environment Variables */
	url := os.Getenv("DB_URL")
	if url != "" {
		return url, nil
	}
	/* 2. If the DB_URL env variable is missing, try building the db url from discretized db env variables */
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")
	/* 3. If also discretized db env variables are missing, return an error... */
	if username == "" || password == "" || host == "" || port == "" || dbname == "" {
		return "", errors.New("Missing/Invalid DB Environment Variables")
	}
	/* 4. If they are present in the .env file, build the URL manually combining their values */
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		username, password, host, port, dbname)
	/* 5. Retur Connection String and null error object */
	return connStr, nil
}
