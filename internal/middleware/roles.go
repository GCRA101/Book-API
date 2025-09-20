package middleware

// middleware/ PACKAGE **********************************************************************************************
/* The middleware/ package stores all the MIDDLEWARE functions that allow to add functionalities to the HTTP Handlers
   that are defined in the handlers/ package.
   This is achieved using the DECORATOR PATTERN. */

/* IMPORTANT NOTES ----------------------------------------------------------------------------------------------*/
/* 1. AllowRoles Middleware
	- The middleware function below is used to decorate HTTP Request Handlers with the ability of checking whether
	  the user is allowed to access the requested endpoint/resource based on their role.
   2. Use of Hash Tables (Sets) instead of Lists/Arrays
	- In the AllowRoles function below, the lookup of allowed roles gets done on a Set rather than a List/Array.
	  This is because, as we know, Hash Tables perform way better for Search Algorithm with a computational cost of
	  θ(1) rather than θ(n) ! Worth the effort of building the Set from scratch in the function below instead of
	  relying on a simple list/array
*/

// 1. IMPORT PACKAGES *********************************************************************************************
import (
	"bookapi/internal/utils"
	"net/http"
)

// 2. CUSTOM http.Handlers ****************************************************************************************

/* ROLE-BASED AUTH Middleware ---------------------------------------------------------------------------------- */
/* Middleware designed to restrict access to certain HTTP endpoints based on the user's role.
   Higher-order function that takes a list of allowed roles and returns a middleware function.*/
func AllowRoles(allowed ...string) func(http.Handler) http.Handler {
	/* 1. Create a set (using a map) of allowed roles for fast lookup.
	Essentially create a Hash Table that has, as keys, all the different allowed roles provided in the
	input list and, as corresponding values, empty lists....These lists are useless but using a Hash
	Table rather than a list/array to lookup for role values allow to get a Search Computational Cost
	of θ(1) rather than θ(n) ! The importance of knowing Algorithms Theory ;) */
	roleSet := make(map[string]struct{}, len(allowed))
	for _, role := range allowed {
		roleSet[role] = struct{}{}
	}
	/* 2. Wrap the original handler (next) and add role-checking logic before calling it. */
	return func(next http.Handler) http.Handler {
		/* 3. Actual Handler Function that runs for every registered HTTP request. */
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			/* 4. Try to get the User's Role from the Context of the HTTP Request. */
			role, ok := r.Context().Value(UserRoleKey).(string)
			/* 5. If the role of the user is empty or cannot be extracted, return error via Helper Function. */
			if !ok || role == "" {
				utils.WriteSafeError(w, http.StatusForbidden, "Forbidden: no role provided")
				return /* <--- NEVER FORGET the RETURN keyword AFTER calling the RESPONSE HELPER FUNCTIONS!! */
			}
			/* 6. If the Role is not in the Set (Hash Table containing allowed roles),
			return error via Helper Function. */
			if _, ok := roleSet[role]; !ok {
				utils.WriteSafeError(w, http.StatusForbidden, "Forbidden: insufficient role")
				return /* <--- NEVER FORGET the RETURN keyword AFTER calling the RESPONSE HELPER FUNCTIONS!! */
			}
			/* 7. If the role is valid proceed to call the original handler. */
			next.ServeHTTP(w, r)
		})
	}
}
