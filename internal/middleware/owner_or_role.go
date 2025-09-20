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

// 3. CUSTOM http.Handlers ********************************************************************************************

/* ROLE-OR-OWNERSHIP-BASED AUTH Middleware ---------------------------------------------------------------------------*/
/* Middleware designed to restrict access to certain HTTP endpoints based on BOTH role AND owner.
   Higher-order function that takes the name of the URL parameter that holds the resource IDb, a function that can
   look up the owner of that resource, context key for the user role and role name to look for.*/
func AllowOwnerOrRole(paramName string, loader OwnerLoader, roleKey contextKey,
	allowedRoles ...string) func(http.Handler) http.Handler {
	/* 1. Create a set (using a map) of allowed roles for fast lookup. */
	roleSet := make(map[string]struct{}, len(allowedRoles))
	for _, role := range allowedRoles {
		roleSet[role] = struct{}{}
	}
	/* 2. Wrap the original handler (next) and add role+owner-checking logic before calling it. */
	return func(next http.Handler) http.Handler {
		/* 3. Actual Handler Function that runs for every registered HTTP request. */
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			/* 4. Get IDs and User Role... */
			/* Get the User's Role from the Context of the HTTP Request. */
			role, _ := r.Context().Value(roleKey).(string)
			/* Get the User's ID from the Context of the HTTP Request. */
			userID := r.Context().Value(UserIDKey).(int)
			/* Get the Resource's ID from the URL of the HTTP Request. */
			idStr := chi.URLParam(r, paramName)
			resourceID, _ := strconv.Atoi(idStr)
			/* Get the Owner's ID via Loader Function */
			ownerID, _ := loader(r, resourceID)
			/* 5. Check Role first and Ownership second... */
			_, isAllowed := roleSet[role]
			if userID != ownerID && !isAllowed {
				utils.WriteSafeError(w, http.StatusForbidden, "Forbidden")
				return
			}
			/* 6. If all good...move on with handling the HTTP Request */
			next.ServeHTTP(w, r)
		})
	}
}
