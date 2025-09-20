package handlers

// handlers/ PACKAGE **********************************************************************************************
/* The handlers/ package stores all the HTTP Method Handlers keeping the HTTP logic separate from
   the other packages. */

/* IMPORTANT NOTES ----------------------------------------------------------------------------------------------*/
/* 1. Scope of user_handler.go
- This go file contain the method Register() that wraps around the services/ method Register() that wraps
around the repositories/ method Create() talking directly to the Database.
*/

// 1. IMPORT PACKAGES *********************************************************************************************

/* Besides the external packages, we also need to import the necessary internal packages defined in the project */
import (
	/* INTERNAL Packages */

	"bookapi/internal/models"
	"bookapi/internal/services"
	"bookapi/internal/utils"

	/* EXTERNAL Packages */
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// 2. GO STRUCTS and UTILITY METHODS  *****************************************************************************

/* STRUCT */
/* Holds a reference to UserService, which contains the logic for registering users. */
type UserHandler struct {
	Service *services.UserService
}

/* STRUCT BUILDER */
/* Creates and returns a new UserHandler instance */
func NewUserHandler(service *services.UserService) *UserHandler {
	return &UserHandler{Service: service}
}

/* Register All Routes */
func (h *UserHandler) RegisterRoutes(r chi.Router) {
	r.Route("/register", func(r chi.Router) {
		/* STATIC Routes */
		r.Post("/", h.Register)
	})
}

// 3. HTTP REQUEST HANDLERS  ***************************************************************************************

/* STATIC HTTP Request Handlers ---------------------------------------------------------------------------------*/

/* POST /register Handler ---------------------------------------------------------------------------------------*/
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	/* 1. Decode JSON Body of HTTP Request + Error Handling */
	var req models.RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		utils.WriteSafeError(w, http.StatusBadRequest, "Invalid Request")
		return
	}
	/* 2. Add record in the Database via the service/ layer + Error Handling */
	user, err := h.Service.Register(req)
	if err != nil {
		utils.WriteSafeError(w, http.StatusBadRequest, err.Error())
		return
	}
	/* 3. Build Go Struct holding id and email of registered user */
	resp := struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
	}{user.ID, user.Email}

	/* 4. Return HTTP Response with 201 Status Code, registered user object and no error */
	utils.WriteJSON(w, http.StatusCreated, resp, nil)

}
