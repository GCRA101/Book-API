package middleware

// middleware/ PACKAGE ************************************************************************************************
/* The middleware/ package stores all the MIDDLEWARE functions that allow to add functionalities to the HTTP Handlers
   that are defined in the handlers/ package.
   This is achieved using the DECORATOR PATTERN. */

// 1. IMPORT PACKAGES *************************************************************************************************
import (
	"bookapi/internal/utils"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// 2. GO STRUCTS and UTILITY METHODS  *********************************************************************************

/* Function type OwnerLoader ----------------------------------------------------------------------------------------*/
/* Function taking a request and a resource ID as inputs, and returning the owner's user ID of that resource as output.
   A function matching this type will be passed to the middleware below. */
type OwnerLoader func(r *http.Request, resourceID int) (int, error)

// 3. CUSTOM http.Handlers ********************************************************************************************

/* OWNERSHIP-BASED AUTH Middleware ----------------------------------------------------------------------------------*/
/* Middleware designed to restrict access to certain HTTP endpoints based on owner.
   Higher-order function that takes the name of the URL parameter that holds the resource ID and a function that can
   look up the owner of that resource.*/
func EnforceOwnership(paramName string, loader OwnerLoader) func(http.Handler) http.Handler {
	/* 1. Wrap the original handler (next) with ownership-checking logic. */
	return func(next http.Handler) http.Handler {
		/* 2. Actual Handler Function that runs for every registered HTTP request. */
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			/* 1. Try to get the User ID out of the Context of the HTTP Request + Error Handling via Helper Function
			- Note: The ID has been set before by the Authentication Middleware. */
			userID, ok := r.Context().Value(UserIDKey).(int)
			if !ok {
				utils.WriteSafeError(w, http.StatusUnauthorized, "Unauthorized")
				return /* <--- NEVER FORGET the RETURN keyword AFTER calling the RESPONSE HELPER FUNCTIONS!! */
			}
			/* 2. Try to extract the resource ID from the URL and convert it to an integer +
			Error Handling via Helper Function */
			idStr := chi.URLParam(r, paramName)
			resourceID, err := strconv.Atoi(idStr)
			if err != nil {
				utils.WriteSafeError(w, http.StatusBadRequest, "Invalid ID")
				return /* <--- NEVER FORGET the RETURN keyword AFTER calling the RESPONSE HELPER FUNCTIONS!! */
			}
			/* 3. Call the OwnerLoader function to find out who owns the resource + Error Handling
			via Helper Function */
			ownerID, err := loader(r, resourceID)
			if err != nil {
				utils.WriteSafeError(w, http.StatusInternalServerError, "Could not verify ownership")
				return /* <--- NEVER FORGET the RETURN keyword AFTER calling the RESPONSE HELPER FUNCTIONS!! */
			}
			/* 4. If user id and owner id don't match, that means that the user doesn't own the
			   resource...hence, an error gets returned using the Helper Function*/
			if userID != ownerID {
				utils.WriteSafeError(w, http.StatusForbidden, "Forbidden: not owner")
				return /* <--- NEVER FORGET the RETURN keyword AFTER calling the RESPONSE HELPER FUNCTIONS!! */
			}
			/* 5. If the user is also the owner of the resource, let the request continue */
			next.ServeHTTP(w, r)
		})
	}
}
