package utils

import (
	/* INTERNAL Packages */
	"bookapi/internal/models"
	/* EXTERNAL Packages */
	"encoding/json"
	"net/http"
)

// 1. RESPONSE HELPER FUNCTIONS  **********************************************************************************

/* Success Response ---------------------------------------------------------------------------------------------*/

func WriteJSON(w http.ResponseWriter, statusCode int, data interface{}, meta interface{}) {
	/* 1. Build the Go Struct to be turned into JSON in the Body of the HTTP Response. */
	response := models.SuccessResponse{
		Data: data,
		Meta: meta,
	}
	/* 2. Set the Content-Type of the Body of the HTTP Response. */
	w.Header().Set("Content-Type", "application/json")
	/* 3. Set the Status Code of the HTTP Response. */
	w.WriteHeader(statusCode)
	/* 4. Convert the Go Struct into JSON, write it to the Body of the HTTP Response and send it to the Client */
	json.NewEncoder(w).Encode(response)
}

/* Error Response -----------------------------------------------------------------------------------------------*/

func WriteError(w http.ResponseWriter, statusCode int, err error, message string) {
	/* 1. Build up the Go Struct instance to be turned into JSON */
	response := models.ErrorResponse{
		Error:   err.Error(),
		Message: message,
	}
	/* 2. Set up the Content-Type of the Body of the HTTP Response */
	w.Header().Set("Content-Type", "application/json")
	/* 3. Set the HTTP Status Code of the HTTP Response. */
	w.WriteHeader(statusCode)
	/* 4. Convert the Go Struct into JSON, write it to the Body of the HTTP Response and send it to the Client */
	json.NewEncoder(w).Encode(response)
}

/* Error Safe Response ------------------------------------------------------------------------------------------*/

func WriteSafeError(w http.ResponseWriter, statusCode int, message string) {
	/* 1. Build up the Go Struct that gets turned into JSON */
	response := models.ErrorResponse{
		Error:   http.StatusText(statusCode),
		Message: message,
	}
	/* 2. Set the Contety-Type of the Body of the HTTP Response */
	w.Header().Set("Content-Type", "application/json")
	/* 3. Set the HTTP Status Code of the HTTP Response */
	w.WriteHeader(statusCode)
	/* 4. Convert the Go Struct into JSON, write it to the Body of the HTTP Response and send it to the Client */
	json.NewEncoder(w).Encode(response)
}
