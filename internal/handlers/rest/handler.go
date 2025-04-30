package rest

import (
	"contact-app/internal/domain"
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type Contacts interface {
	Create(ctx context.Context, contact domain.Contact) error
	GetByID(ctx context.Context, id int64) (domain.Contact, error)
	GetAll(ctx context.Context) ([]domain.Contact, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, id int64, inp domain.UpdateContact) error
}

type Handler struct {
	contactsService Contacts
}

func NewHandler(contact Contacts) *Handler {
	return &Handler{
		contactsService: contact,
	}
}

func (h *Handler) InitRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(loggingMiddleWare)

	contacts := r.PathPrefix("/contacts").Subrouter()
	{
		contacts.HandleFunc("", h.createContact).Methods(http.MethodPost)
		contacts.HandleFunc("", h.getAllContacts).Methods(http.MethodGet)
		contacts.HandleFunc("/{id:[0-9]+}", h.GetContactById).Methods(http.MethodGet)
		contacts.HandleFunc("/{id:[0-9]+}", h.deleteContact).Methods(http.MethodDelete)
		contacts.HandleFunc("/{id:[0-9]+}", h.updateContact).Methods(http.MethodPut)
	}

	return r
}

func (h *Handler) GetContactById(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromRequest(r)
	if err != nil {
		log.Println("getContactByID() error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	contact, err := h.contactsService.GetByID(context.TODO(), id)
	if err != nil {
		if errors.Is(err, domain.ErrContactNotFound) {
			w.WriteHeader(http.StatusBadRequest)
		}

		log.Println("getBookByID() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(contact)
	if err != nil {
		log.Println("getBookByID() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}

func (h *Handler) createContact(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var contact domain.Contact
	if err = json.Unmarshal(reqBytes, &contact); err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	err = h.contactsService.Create(context.TODO(), contact)
	if err != nil {
		log.Println("createBook() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) deleteContact(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromRequest(r)
	if err != nil {
		log.Println("deleteContact() error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.contactsService.Delete(context.TODO(), id)
	if err != nil {
		log.Println("deleteContact() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) getAllContacts(w http.ResponseWriter, r *http.Request) {
	contact, err := h.contactsService.GetAll(context.TODO())
	if err != nil {
		log.Println("getAllContacts() error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response, err := json.Marshal(contact)
	if err != nil {
		log.Println("getAllContacts() error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}

func (h *Handler) updateContact(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromRequest(r)
	if err != nil {
		log.Println("error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	reqBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var inp domain.UpdateContact
	if err = json.Unmarshal(reqBytes, &inp); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.contactsService.Update(context.TODO(), id, inp)
	if err != nil {
		log.Println("error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func getIdFromRequest(r *http.Request) (int64, error) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		return 0, err
	}

	if id == 0 {
		return 0, err
	}

	return id, nil
}
