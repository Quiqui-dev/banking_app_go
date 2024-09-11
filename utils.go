package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	jwt "github.com/golang-jwt/jwt/v5"
)

func WriteJson(w http.ResponseWriter, status int, v any) error {

	w.Header().Set("Content-Type", "application/json")

	// headers should be set before we call WriteHeader
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}

func ReadJson(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func withJWTAth(handlerFunc http.HandlerFunc, s Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		tokenString := r.Header.Get("x-jwt-token")

		token, err := validateJWT(tokenString)

		if err != nil {
			permissionDeniedResponse(w)
			return
		}

		if !token.Valid {
			permissionDeniedResponse(w)
			return
		}

		userID, err := GetIDFromRoute(r)

		if err != nil {
			permissionDeniedResponse(w)
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		claimAccount := int64(claims["accountNumber"].(float64))

		account, err := s.GetAccountByNumber(claimAccount)

		if err != nil {
			permissionDeniedResponse(w)
			return
		}

		if account.ID != userID {
			permissionDeniedResponse(w)
			return
		}

		handlerFunc(w, r)
	}
}

func validateJWT(tokenString string) (*jwt.Token, error) {

	secret := os.Getenv("JWT_SECRET")

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})

}

func createJWT(account *Account) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"expiresAt":     15000,
		"accountNumber": account.Number,
	})

	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func permissionDeniedResponse(w http.ResponseWriter) {
	WriteJson(w, http.StatusUnauthorized, apiError{Error: "permission denied"})
}

func httpMethodNotAllowed(httpMethod string) error {
	return fmt.Errorf("method not permitted: %s", httpMethod)
}
