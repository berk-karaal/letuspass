package controllers_test

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/berk-karaal/letuspass/backend/internal/models"
	authservice "github.com/berk-karaal/letuspass/backend/internal/services/auth"
	"github.com/berk-karaal/letuspass/backend/internal/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleAuthRegister(t *testing.T) {
	router, _, postgresDb := tests.SetupTestRouter()
	defer tests.CleanDatabase(postgresDb)

	reqBody, err := json.Marshal(map[string]interface{}{
		"email":               "test@example.com",
		"key_derivation_salt": "randomSalt",
		"name":                "Mr Test",
		"password":            "testPass",
		"public_key":          "randomPublicKey",
	})
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}
	req, _ := http.NewRequest("POST", "/api/v1/auth/register", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status %d; got %d. Response: %v", http.StatusCreated, w.Code, w.Body.String())
	}

	var users []models.User
	err = postgresDb.Find(&users).Error
	if err != nil {
		t.Fatalf("failed to query users: %v", err)
	}

	if len(users) != 1 {
		t.Fatalf("expected 1 user; got %d", len(users))
	}

	user := users[0]
	assert.Equal(t,
		map[string]interface{}{"id": 1, "email": "test@example.com", "name": "Mr Test", "public_key": "randomPublicKey",
			"key_derivation_salt": "randomSalt"},
		map[string]interface{}{"id": int(user.ID), "email": user.Email, "name": user.Name, "public_key": user.PublicKey,
			"key_derivation_salt": user.KeyDerivationSalt},
	)
}

func TestHandleAuthLogin(t *testing.T) {
	router, apiConfig, postgresDb := tests.SetupTestRouter()
	defer tests.CleanDatabase(postgresDb)

	hashedPassword, err := authservice.HashPassword("testPass")
	if err != nil {
		log.Fatalf("failed to hash password: %v", err)
	}
	user := models.User{
		Email:             "test@example.com",
		Password:          hashedPassword,
		Name:              "Mr Test",
		KeyDerivationSalt: "randomSalt",
		PublicKey:         "randomPublicKey",
		IsActive:          true,
	}
	err = postgresDb.Create(&user).Error
	if err != nil {
		log.Fatalf("failed to create test user: %v", err)
	}

	reqBody, err := json.Marshal(map[string]interface{}{
		"email":    "test@example.com",
		"password": "testPass",
	})
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status %d; got %d. Response: %v", http.StatusCreated, w.Code, w.Body.String())
	}

	expectedResponse, err := json.Marshal(map[string]interface{}{
		"email":               "test@example.com",
		"key_derivation_salt": "randomSalt",
		"name":                "Mr Test",
	})
	if err != nil {
		t.Fatalf("failed to marshal expected response: %v", err)
	}
	require.JSONEq(t, string(expectedResponse), w.Body.String())

	// Check if session token cookie is set
	setCookieHeader := w.Header().Get("Set-Cookie")
	if !strings.Contains(setCookieHeader, apiConfig.SessionTokenCookieName+"=") {
		t.Errorf("expected cookie header to contain %s; got %s", apiConfig.SessionTokenCookieName, setCookieHeader)
	}

	// Check if user session is created
	var userSession models.UserSession
	err = postgresDb.First(&userSession).Error
	if err != nil {
		t.Fatalf("failed to query user session: %v", err)
	}
	if userSession.UserID != 1 {
		t.Fatalf("expected user id %d; got %d", 1, userSession.UserID)
	}
}
