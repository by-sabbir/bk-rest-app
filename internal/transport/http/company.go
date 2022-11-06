package http

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/by-sabbir/company-microservice-rest/internal/company"
	"github.com/gorilla/mux"
)

type CompanyRestService interface {
	GetCompany(context.Context, string) (company.Company, error)
	PostCompany(context.Context, company.Company) (company.Company, error)
	DeleteCompany(context.Context, string) error
	PartialUpdateCompany(context.Context, string, company.Company) (company.Company, error)
}

func (h *Handler) PostCompany(w http.ResponseWriter, r *http.Request) {
	var cmp company.Company

	if err := json.NewDecoder(r.Body).Decode(&cmp); err != nil {
		log.Println("error posting company: ", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	postedCmp, err := h.Service.PostCompany(r.Context(), cmp)
	if err != nil {
		log.Println("error posting company: ", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(postedCmp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("error encoding message: ", err)
	}
}

func (h *Handler) GetCompany(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	cmp, err := h.Service.GetCompany(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(cmp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) DeleteCompany(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := h.Service.DeleteCompany(r.Context(), id)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) PartialUpdateCompany(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var cmp company.Company

	if err := json.NewDecoder(r.Body).Decode(&cmp); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	updatedCmp, err := h.Service.PartialUpdateCompany(r.Context(), id, cmp)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"error": err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusPartialContent)
	if err := json.NewEncoder(w).Encode(updatedCmp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
