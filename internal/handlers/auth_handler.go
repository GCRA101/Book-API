package handlers

// handlers/ PACKAGE **********************************************************************************************
/* The handlers/ package stores all the HTTP Method Handlers keeping the HTTP logic separate from
   the other packages. */

/* IMPORTANT NOTES ----------------------------------------------------------------------------------------------*/
/* 1. Scope of auth_handler.go
   - This go file contain the method Login() that wraps around the services/ method FindByEmail that wraps around the
   	 repositories/ method FindByEmail talking directly to the Database.
     In addition to that it also carries out the creation of the Token than can be used by the client to keep getting
     access to the API endpoints during the entire user's session.
*/

// 1. IMPORT PACKAGES *********************************************************************************************
import (
	/* INTERNAL Packages */

	"bookapi/internal/security"
	"bookapi/internal/services"
	"bookapi/internal/utils"

	/* EXTERNAL Packages */
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// 2. GO STRUCTS and UTILITY METHODS  ******************************************************************************

/* STRUCT for Login */
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

/* STRUCT for Authentication via Token */
type AuthHandler struct {
	UserService *services.UserService
	JWTSecret   string
}

/* STRUCT BUILDER */
/* Creates and returns a new UserHandler instance */
func NewAuthHandler(service *services.UserService, secret string) *AuthHandler {
	return &AuthHandler{UserService: service, JWTSecret: secret}
}

/* Register All Routes */
func (h *AuthHandler) RegisterRoutes(r chi.Router) {
	/* STATIC Routes */
	r.Post("/login", h.Login)
}

// 3. HTTP REQUEST HANDLERS  ***************************************************************************************

/* STATIC HTTP Request Handlers ---------------------------------------------------------------------------------*/

/* POST /login Handler */
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	/* 1. Create blank LoginRequest Struct to hold data from the HTTP Request Body's JSON */
	var req LoginRequest
	/* 2. Convert JSON from Body of HTTP Request into LoginRequest Struct + Error Handling via Helper Function */
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.WriteSafeError(w, http.StatusBadRequest, "Invalid input")
		return
	}
	/* 3. Look into Database for User object matching input email + Error Handling via Helper Function */
	user, err := h.UserService.FindByEmail(req.Email)
	if err != nil || user == nil {
		utils.WriteSafeError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}
	/* 4. If User exists..compare input textual Password with stored Hash. + Error Handling via Helper Function */
	if !security.CheckPasswordHash(req.Password, user.Password) {
		utils.WriteSafeError(w, http.StatusUnauthorized, "Invalid email or password")
		return
	}
	/* 5. If user exists and password is correct....generate Token via JWT + Error Handling via Helper Function */
	token, err := security.GenerateToken(user.ID, user.Role, h.JWTSecret)
	if err != nil {
		utils.WriteSafeError(w, http.StatusInternalServerError, "Failed to generate token.")
		return
	}
	/* 6. Return HTTP Response with 200 Status Code + Token as JSON in the Body via Helper Function */
	utils.WriteJSON(w, http.StatusOK, token, nil)
}
