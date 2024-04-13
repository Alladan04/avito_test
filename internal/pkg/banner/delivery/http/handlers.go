package http

import (
	"net/http"
	"strconv"

	"github.com/Alladan04/avito_test/internal/models"
	"github.com/Alladan04/avito_test/internal/pkg/banner"
	"github.com/Alladan04/avito_test/internal/pkg/utils"
	"github.com/gorilla/mux"
)

type BannerHandler struct {
	uc banner.BannerUsecase
}

func NewBannerHandler(uc banner.BannerUsecase) *BannerHandler {
	return &BannerHandler{
		uc: uc,
	}
}

// AddItem to create new banner
// for admins only
func (h *BannerHandler) AddItem(w http.ResponseWriter, r *http.Request) {

	jwtPayload, ok := r.Context().Value(models.PayloadContextKey).(models.JwtPayload)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if !jwtPayload.IsAdmin {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	item := models.BannerForm{}
	err := utils.GetRequestData(r, &item)
	if err != nil {
		utils.WriteErrorMessage(w, http.StatusBadRequest, "error unmarshalling")
		return
	}

	res, err := h.uc.AddItem(r.Context(), item)
	if err != nil {
		utils.WriteErrorMessage(w, http.StatusBadRequest, "error Adding data")
		return

	}

	err = utils.WriteResponseData(w, res, http.StatusCreated)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

}

// GetAll handler returns all banners (or those which were selected via query pramas feature_id and tag_id)
// for admins only
func (h *BannerHandler) GetAll(w http.ResponseWriter, r *http.Request) {

	//read query params
	countParam := r.URL.Query().Get("limit")
	offsetParam := r.URL.Query().Get("offset")
	featureParam := r.URL.Query().Get("feature_id")
	tagParam := r.URL.Query().Get("tag_id")

	//process offset and count params
	count, err := strconv.ParseInt(countParam, 10, 64)
	if err != nil && countParam != "" {
		utils.WriteErrorMessage(w, http.StatusBadRequest, "wrong count param")
		return
	}
	offset, err := strconv.ParseInt(offsetParam, 10, 64)
	if err != nil && offsetParam != "" {
		utils.WriteErrorMessage(w, http.StatusBadRequest, "wrong offset param")
		return
	}

	//process feature_id and tag_id
	featureId, err := strconv.ParseInt(featureParam, 10, 64)
	if err != nil && featureParam != "" {
		utils.WriteErrorMessage(w, http.StatusBadRequest, "wrong feature_id param")
		return
	}
	tagId, err := strconv.ParseInt(tagParam, 10, 64)
	if err != nil && tagParam != "" {
		utils.WriteErrorMessage(w, http.StatusBadRequest, "wrong tag_d param")
		return
	}
	//get result from usecase
	result, err := h.uc.GetAll(r.Context(), count, offset, featureId, tagId)
	if err != nil {
		utils.WriteErrorMessage(w, http.StatusBadRequest, err.Error())

		return

	}
	err = utils.WriteResponseData(w, result, http.StatusOK)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func (h *BannerHandler) GetOne(w http.ResponseWriter, r *http.Request) {
	//process_params
	featureParam := r.URL.Query().Get("feature_id")
	tagParam := r.URL.Query().Get("tag_id")
	useLastRevisionParam := r.URL.Query().Get("use_last_revision")
	featureId, err := strconv.ParseInt(featureParam, 10, 64)
	if err != nil {
		utils.WriteErrorMessage(w, http.StatusBadRequest, "wrong feature_id param")
		return
	}
	tagId, err := strconv.ParseInt(tagParam, 10, 64)
	if err != nil {
		utils.WriteErrorMessage(w, http.StatusBadRequest, "wrong tag_d param")
		return
	}
	useLastRevision, err := strconv.ParseBool(useLastRevisionParam)
	if err != nil {
		useLastRevision = false
	}

	result, err := h.uc.GetOne(r.Context(), featureId, tagId, useLastRevision)
	if err != nil {
		utils.WriteErrorMessage(w, http.StatusNotFound, "not found")
		return
	}
	err = utils.WriteResponseData(w, result, http.StatusOK)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

}

func (h *BannerHandler) UpdateBanner(w http.ResponseWriter, r *http.Request) {
	bannerIdString := mux.Vars(r)["id"]
	bannerId, err := strconv.ParseInt(bannerIdString, 10, 64)
	if err != nil {
		utils.WriteErrorMessage(w, http.StatusBadRequest, "cant parse id")
		return
	}
	item := models.BannerUpdateForm{}
	err = utils.GetRequestData(r, &item)
	if err != nil {
		utils.WriteErrorMessage(w, http.StatusBadRequest, "error unmarshalling")
		return
	}
	err = h.uc.UpdateBanner(r.Context(), item, bannerId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)

}
func (h *BannerHandler) DeleteBanner(w http.ResponseWriter, r *http.Request) {
	bannerIdString := mux.Vars(r)["id"]
	bannerId, err := strconv.ParseInt(bannerIdString, 10, 64)
	if err != nil {
		utils.WriteErrorMessage(w, http.StatusBadRequest, "cant parse id")
		return
	}

	err = h.uc.DeleteBanner(r.Context(), bannerId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
