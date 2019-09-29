package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	authpb "github.com/benjohns1/basic-blog-ms/services/authentication/proto"
)

// LoginCommand user credentials
type LoginCommand struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

// LoginResponse responds to a valid login request
type LoginResponse struct {
	Token string `json:"token"`
}

// AuthenticationService authentication handlers
type AuthenticationService struct {
	Authentication authpb.AuthenticationClient
}

// LoginHandler handles http login request
func (s *AuthenticationService) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}

	// Parse API request data
	loginData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeError(w, err, http.StatusInternalServerError)
		return
	}
	r.Body.Close()

	loginCmd := LoginCommand{}
	if err := json.Unmarshal(loginData, &loginCmd); err != nil {
		writeError(w, err, http.StatusUnauthorized)
		return
	}

	// Make service RPC call
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resp, err := s.Authentication.Login(ctx, &authpb.LoginCommand{Username: loginCmd.User, Password: loginCmd.Password})
	if err != nil {
		writeError(w, fmt.Errorf("error logging in %v", err), http.StatusInternalServerError)
		return
	}

	if !resp.GetSuccess() {
		writeError(w, fmt.Errorf("unauthorized"), http.StatusUnauthorized)
		return
	}

	// Build API response
	loginResp, err := json.Marshal(&LoginResponse{Token: resp.GetToken()})
	if err != nil {
		writeError(w, fmt.Errorf("error processing login %v", err), http.StatusInternalServerError)
		return
	}
	writeResponse(w, r, loginResp, 200)
}

// Authenticate determines if an http request is authenticated
func (s *AuthenticationService) Authenticate(r *http.Request) (bool, error) {
	// Parse API request
	token := r.Header.Get("Authorization")

	// Make service RPC call
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	resp, err := s.Authentication.Authenticate(ctx, &authpb.AuthenticateQuery{Token: token})
	if err != nil {
		return false, err
	}

	return resp.GetSuccess(), nil
}
