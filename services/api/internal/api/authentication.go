package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Credentials user credentials
type Credentials struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

// AuthenticateResponse responds to a valid authentication request
type AuthenticateResponse struct {
	Token string `json:"token"`
}

func authenticateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}

	credentialData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	r.Body.Close()

	creds := Credentials{}
	if err := json.Unmarshal(credentialData, &creds); err != nil {
		writeError(w, err, http.StatusUnauthorized)
		return
	}

	token, err := getToken(creds)
	if err != nil {
		writeError(w, fmt.Errorf("Unauthorized"), http.StatusUnauthorized)
		return
	}

	writeResponse(w, r, getAuthResponse(token))
}

func getAuthResponse(token string) func() ([]byte, error) {
	return func() ([]byte, error) {
		authResponse := AuthenticateResponse{Token: token}
		return json.Marshal(authResponse)
	}
}

func getToken(creds Credentials) (string, error) {
	if creds.User == "bobross" && creds.Password == "painter" {
		return "notSoSecureToken!", nil
	}
	return "", fmt.Errorf("Invalid credentials")
}

func authenticate(r *http.Request) bool {
	token := r.Header.Get("Authorization")
	return "notSoSecureToken!" == token
}
