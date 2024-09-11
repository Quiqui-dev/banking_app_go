package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

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

		defer r.Body.Close()
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
	router.HandleFunc("/login", makeHttpHandleFunc(s.handleLogin))
	router.Handle("/account", makeHttpHandleFunc(s.handleAccount))
	router.Handle("/account/{id}", withJWTAth(makeHttpHandleFunc(s.handleAccountById), s.store))
	router.Handle("/transfer", makeHttpHandleFunc(s.handleTransfer))

	log.Printf("JSON API server running on port: %s", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {

	switch r.Method {
	case http.MethodPost:
		break
	default:
		return httpMethodNotAllowed(r.Method)
	}

	loginRequest := new(LoginRequest)

	if err := ReadJson(r, loginRequest); err != nil {
		return err
	}

	account, err := s.store.GetAccountByNumber(loginRequest.AccountNumber)

	if err != nil {
		return err
	}

	tokenString, err := createJWT(account)

	if err != nil {
		return err
	}

	resp := &LoginResponse{
		AccountNumber: account.Number,
		Token:         tokenString,
	}

	return WriteJson(w, http.StatusOK, resp)
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {

	switch r.Method {
	case http.MethodPost:
		return s.handleCreateAccount(w, r)
	case http.MethodGet:
		return s.handleGetAccount(w, r)
	default:
		return httpMethodNotAllowed(r.Method)
	}
}

func (s *APIServer) handleAccountById(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case http.MethodGet:
		return s.handleGetAccountById(w, r)
	case http.MethodDelete:
		return s.handleDeleteAccount(w, r)
	default:
		return httpMethodNotAllowed(r.Method)
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

	id, err := GetIDFromRoute(r)

	if err != nil {
		return fmt.Errorf("invalid id provided")
	}

	account, err := s.store.GetAccountById(id)

	if err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, account)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {

	accReq := new(CreateAccountRequest)

	// read the json values if payload is corrrect
	if err := ReadJson(r, accReq); err != nil {
		return err
	}

	// create account
	account, err := NewAccount(accReq.FirstName, accReq.LastName, accReq.Password)

	if err != nil {
		return err
	}

	account, err = s.store.CreateAccount(account)

	if err != nil {
		return err
	}

	return WriteJson(w, http.StatusCreated, account)
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {

	id, err := GetIDFromRoute(r)

	if err != nil {
		return fmt.Errorf("invalid id provided")
	}

	if err := s.store.DeleteAccount(id); err != nil {
		return err
	}

	return WriteJson(w, http.StatusOK, "")
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {

	transferReq := new(TransferRequest)

	// read the json values if payload is corrrect
	if err := ReadJson(r, transferReq); err != nil {
		return err
	}

	return nil
}

func GetIDFromRoute(r *http.Request) (int, error) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))

	if err != nil {
		return 0, err
	}

	return id, nil
}
