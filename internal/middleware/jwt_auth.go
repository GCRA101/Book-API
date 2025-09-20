package middleware

// middleware/ PACKAGE *************************************************************************************************
/* The middleware/ package stores all the MIDDLEWARE functions that allow to add functionalities to the HTTP Handlers
   that are defined in the handlers/ package.
   This is achieved using the DECORATOR PATTERN. */

// 1. IMPORT PACKAGES **************************************************************************************************
import (
	/* INTERNAL Packages */
	"bookapi/internal/security"
	"bookapi/internal/utils"

	/* EXTERNAL Packages */
	"context"
	"net/http"
	"strings"
)

type contextKey string

const UserIDKey contextKey = "user_id"
const UserRoleKey contextKey = "user_role"

// 2. CUSTOM http.Handlers *********************************************************************************************

/* JWT TOKEN AUTHENTICATION Middleware ------------------------------------------------------------------------------ */
/*
The following middleware method carries out the following tasks:
 1. Extract the token from the header of the HTTP Request.
 2. Verify the token's signature and expiration date
 3. Inject the user ID into the request context
*/
func JWTAuth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			/* 1. Get the value of the Authorization Header of the HTTP Request + Error Handling via Helper Function*/
			auth := r.Header.Get("Authorization")
			/*..if it’s missing or doesn’t start with "Bearer", it means the user didn’t send a proper token..*/
			if auth == "" || !strings.HasPrefix(auth, "Bearer") {
				utils.WriteSafeError(w, http.StatusUnauthorized, "Unauthorized")
				return
			}
			/* 2. Extract the Token + Check its validity + Error Handling via Helper Function */
			tokenStr := strings.TrimPrefix(auth, "Bearer")
			claims, err := security.ParseToken(tokenStr, secret)
			if err != nil {
				utils.WriteSafeError(w, http.StatusUnauthorized, "Invalid or expired token.")
				return
			}
			/* 3. Try to get the user_id from the token's data + Error Handling via Helper Function */
			userIDRaw, ok := claims["user_id"]
			if !ok {
				utils.WriteSafeError(w, http.StatusUnauthorized, "Missing user_id in token.")
				return
			}
			/* 4. Try to get the user_role from the token's data + Error Handling via Helper Function */
			userRoleRaw, ok := claims["user_role"]
			if !ok {
				utils.WriteSafeError(w, http.StatusUnauthorized, "Missing user_role in token.")
				return
			}
			/* 5. Convert the user ID into an integer and user ROLE into a string*/
			userID := int(userIDRaw.(float64))
			userRole := userRoleRaw.(string)
			/* 6. Add the user ID and user ROLE to the request's context */
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			ctx = context.WithValue(ctx, UserRoleKey, userRole)
			/* 7. Passes the request (enriched with the userID info) to the next handler */
			next.ServeHTTP(w, r.WithContext(ctx))
			/*...Now the handler can access the user ID and know who made the request...*/
		})
	}
}
