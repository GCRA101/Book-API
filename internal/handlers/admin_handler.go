package handlers

// handlers/ PACKAGE **********************************************************************************************
/* The handlers/ package stores all the HTTP Method Handlers keeping the HTTP logic separate from
   the other packages. */

/* IMPORTANT NOTES ----------------------------------------------------------------------------------------------*/
/* 1. Scope of admin_handler.go
- This go file contain the method GetUsers() that wraps around the services/ method FindAll() that wraps
around the repositories/ method FindAll() talking directly to the Database.
*/

// 1. IMPORT PACKAGES *********************************************************************************************

/* Besides the external packages, we also need to import the necessary internal packages defined in the project */
import (
	/* INTERNAL Packages */
	"bookapi/internal/middleware"
	"bookapi/internal/services"
	"bookapi/internal/utils"
	"fmt"

	/* EXTERNAL Packages */

	"net/http"

	"github.com/go-chi/chi/v5"
)

// 2. GO STRUCTS and UTILITY METHODS  *****************************************************************************

/* STRUCT */
/* Holds a reference to UserService, which contains the logic for registering users. */
type AdminHandler struct {
	Service *services.UserService
}

/* STRUCT BUILDER */
/* Creates and returns a new UserHandler instance */
func NewAdminHandler(service *services.UserService) *AdminHandler {
	return &AdminHandler{Service: service}
}

/* Register All Routes */
func (h *AdminHandler) RegisterRoutes(r chi.Router) {
	r.Route("/admin", func(r chi.Router) {
		r.With(middleware.AllowRoles("admin")).Get("/users", h.GetUsers)     /*		>>>>>> ROLE-BASED AUTH <<<<<<*/
		r.With(middleware.AllowRoles("admin")).Get("/profile", h.GetProfile) /*		>>>>>> ROLE-BASED AUTH <<<<<<*/
	})

}

// 3. HTTP REQUEST HANDLERS  ***************************************************************************************

/* STATIC HTTP Request Handlers ---------------------------------------------------------------------------------*/

/* GET /users Handler */
func (h *AdminHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.Service.FindAll()
	if err != nil {
		utils.WriteSafeError(w, http.StatusInternalServerError, "Could Not Fetch Books.")
		return
	}
	utils.WriteJSON(w, http.StatusOK, users, nil)
}

/* GET /profile Handler */
func (h *AdminHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey).(int)
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprintf(w, "Welcome user %d", userID)
}
