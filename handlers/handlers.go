package handlers

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

func CreateClientHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
		return
	}
	adminToken := strings.TrimPrefix(authHeader, "Bearer ")

	var clientData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&clientData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	keycloakURL := "http://localhost:8080/admin/realms/master/clients"
	clientJSON, _ := json.Marshal(clientData)
	req, err := http.NewRequest("POST", keycloakURL, bytes.NewBuffer(clientJSON))
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Authorization", "Bearer "+adminToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "Failed to communicate with Keycloak", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}

func GeneratePKCEHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var requestData struct {
		RedirectURI string `json:"redirect_uri"`
		ClientID    string `json:"client_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestData); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if requestData.RedirectURI == "" || requestData.ClientID == "" {
		http.Error(w, "Missing required parameters", http.StatusBadRequest)
		return
	}

	codeVerifier := generateCodeVerifier()
	codeChallenge := generateCodeChallenge(codeVerifier)

	authURL := "http://localhost:8080/realms/master/protocol/openid-connect/auth" +
		"?response_type=code" +
		"&client_id=" + url.QueryEscape(requestData.ClientID) +
		"&redirect_uri=" + url.QueryEscape(requestData.RedirectURI) +
		"&code_challenge=" + url.QueryEscape(codeChallenge) +
		"&code_challenge_method=S256" +
		"&prompt=login"

	response := map[string]string{
		"code_verifier":     codeVerifier,
		"code_challenge":    codeChallenge,
		"authorization_url": authURL,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Extract all query parameters
	queryParams := r.URL.Query()

	// Convert query parameters to a JSON object
	response := make(map[string]string)
	for key, values := range queryParams {
		if len(values) > 0 {
			response[key] = values[0]
		}
	}

	// Return the response as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func generateCodeVerifier() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

func generateCodeChallenge(verifier string) string {
	hash := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}
