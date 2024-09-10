package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
)

type APIServer struct {
	listenAddr string
	store      Storage
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error

type apiError struct {
	Error string
}

func makeHttpHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			// handle err

			WriteJson(w, http.StatusBadRequest, apiError{Error: err.Error()})
		}
	}
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {

	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTION"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	//v1Router := chi.NewRouter()

	router.Handle("/account", makeHttpHandleFunc(s.handleAccount))
	router.Handle("/account/{id}", makeHttpHandleFunc(s.handleGetAccountById))

	log.Printf("JSON API server running on port: %s", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {

	switch r.Method {
	case http.MethodPost:
		return s.handleCreateAccount(w, r)
	case http.MethodGet:
		return s.handleGetAccount(w, r)
	case http.MethodDelete:
		return s.handleDeleteAccount(w, r)
	default:
		return fmt.Errorf("method not permitted: %s", r.Method)
	}
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {

	accounts, err := s.store.GetAccounts()

	if err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, accounts)
}

func (s *APIServer) handleGetAccountById(w http.ResponseWriter, r *http.Request) error {

	vars := chi.URLParam(r, "id")

	return WriteJson(w, http.StatusOK, map[string]string{"id": vars})
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {

	accReq := new(CreateAccountRequest)

	// read the json values if payload is corrrect
	if err := ReadJson(r, accReq); err != nil {
		return err
	}

	// create account
	account := NewAccount(accReq.FirstName, accReq.LastName)

	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	return WriteJson(w, http.StatusCreated, account)
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {

	return nil
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {

	return nil
}
