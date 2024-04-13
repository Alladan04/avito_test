package http

import (
	"net/http"

	"github.com/Alladan04/avito_test/internal/models"
	"github.com/Alladan04/avito_test/internal/pkg/auth"
	"github.com/Alladan04/avito_test/internal/pkg/utils"
)

type AuthHandler struct {
	uc auth.AuthUsecase
}

func NewAuthHandler(uc auth.AuthUsecase) *AuthHandler {
	return &AuthHandler{
		uc: uc,
	}
}
func (h *AuthHandler) SignUp(w http.ResponseWriter, r *http.Request) {

	userData := models.UserForm{}
	if err := utils.GetRequestData(r, &userData); err != nil {
		utils.WriteErrorMessage(w, http.StatusBadRequest, auth.ErrIncorrectPayload.Error())
		return
	}

	if err := userData.Validate(); err != nil {
		utils.WriteErrorMessage(w, http.StatusBadRequest, err.Error())
		return
	}

	newUser, token, _, err := h.uc.SignUp(r.Context(), userData)
	if err != nil {
		utils.WriteErrorMessage(w, http.StatusBadRequest, auth.ErrCreatingUser.Error())
		return
	}

	w.Header().Set("token", "Bearer "+token)

	if err := utils.WriteResponseData(w, newUser, http.StatusCreated); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (h *AuthHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	userData := models.UserForm{}
	if err := utils.GetRequestData(r, &userData); err != nil {
		utils.WriteErrorMessage(w, http.StatusBadRequest, auth.ErrIncorrectPayload.Error())
		return
	}

	if err := userData.Validate(); err != nil {
		utils.WriteErrorMessage(w, http.StatusBadRequest, err.Error())
		return
	}

	user, token, _, err := h.uc.SignIn(r.Context(), userData)
	if err != nil {
		utils.WriteErrorMessage(w, http.StatusBadRequest, auth.ErrUserNotFound.Error())
		return
	}

	w.Header().Set("token", "Bearer "+token)

	if err := utils.WriteResponseData(w, user, http.StatusOK); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *AuthHandler) LogOut(w http.ResponseWriter, r *http.Request) {

	w.Header().Del("token")
	w.WriteHeader(http.StatusNoContent)

}
