package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aakosarev/banner-rotation/internal/model"
	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type service interface {
	AddBannerToSlot(ctx context.Context, bannerID, slotID *uuid.UUID) error
	RemoveBannerFromSlot(ctx context.Context, bannerID, slotID *uuid.UUID) error
	SelectBanner(ctx context.Context, slotID, socialGroupID *uuid.UUID) (*model.Banner, error)
	AddClick(ctx context.Context, bannerID, slotID, socialGroupID *uuid.UUID) error
}

type Handler struct {
	service service
}

func NewHandler(service service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodPost, "/banner", h.AddBannerToSlot)
	router.DELETE("/banner/:banner_id/slot/:slot_id", h.RemoveBannerFromSlot)
	router.GET("/slot/:slot_id/group/:group_id", h.SelectBanner)
	router.POST("/banner/:banner_id/slot/:slot_id/group/:group_id/click", h.AddClick)
}

func (h *Handler) AddBannerToSlot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	bannerSlot := model.BannerSlot{}
	err := json.NewDecoder(r.Body).Decode(&bannerSlot)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message":"Invalid request body"}`))
		return
	}
	err = h.service.AddBannerToSlot(r.Context(), &bannerSlot.BannerID, &bannerSlot.SlotID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err.Error())))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"success"}`))
}

func (h *Handler) RemoveBannerFromSlot(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	bannerID, err := uuid.Parse(params.ByName("banner_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message":"Invalid request body"}`))
		return
	}

	slotID, err := uuid.Parse(params.ByName("slot_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message":"Invalid request body"}`))
		return
	}

	err = h.service.RemoveBannerFromSlot(r.Context(), &bannerID, &slotID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err.Error())))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"success"}`))
}

func (h *Handler) SelectBanner(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	slotID, err := uuid.Parse(params.ByName("slot_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message":"Invalid request body"}`))
		return
	}

	socialGroupID, err := uuid.Parse(params.ByName("group_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message":"Invalid request body"}`))
		return
	}
	selectedBanner, err := h.service.SelectBanner(r.Context(), &slotID, &socialGroupID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err.Error())))
		return
	}

	selectedBannerJson, err := json.Marshal(selectedBanner)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err.Error())))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(selectedBannerJson)
}

func (h *Handler) AddClick(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	bannerID, err := uuid.Parse(params.ByName("banner_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message":"Invalid request body"}`))
		return
	}

	slotID, err := uuid.Parse(params.ByName("slot_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message":"Invalid request body"}`))
		return
	}

	socialGroupID, err := uuid.Parse(params.ByName("group_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message":"Invalid request body"}`))
		return
	}

	err = h.service.AddClick(r.Context(), &bannerID, &slotID, &socialGroupID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err.Error())))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"success"}`))
}
