package main

import (
	"bytes"
	"github.com/bitly/go-simplejson"
	"github.com/gorilla/mux"
	"github.com/redefik/sdccproject/apigateway/config"
	"github.com/redefik/sdccproject/apigateway/microservice"
	"github.com/redefik/sdccproject/apigateway/mock"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

// createTestGateway creates an http handler that handles the test requests
func createTestGatewayLoginUser() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/didattica-mobile/api/v1.0/token", microservice.LoginUser).Methods(http.MethodPost)
	return r
}

// TestUserLoginSuccess tests the following scenario: the client sends an access token request to the api gateway,
// passing username "admin" and password "admin_pass" in the body.
// The gateway makes an http GET request with the given information and sends it to the user management microservice.
// It is assumed that the user "admin" exists and has password "admin_pass", so the authentication should succeed
// and the gateway should create the token responding to the client with a 201 http status code.
func TestUserLoginSuccess(t *testing.T) {

	config.SetConfigurationFromFile("../config/config-test.json")

	jsonBody := simplejson.New()
	jsonBody.Set("username", "admin")
	jsonBody.Set("password", "admin_pass")

	requestBody, _ := jsonBody.MarshalJSON()
	request, _ := http.NewRequest(http.MethodPost, "/didattica-mobile/api/v1.0/token", bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	handler := createTestGatewayLoginUser()
	// a goroutine representing the microservice listens to the requests coming from the api gateway
	go mock.LaunchUserManagementMock()
	// simulates a request-response interaction between client and api gateway
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Error("Expected 201 Created but got " + strconv.Itoa(response.Code) + " " + http.StatusText(response.Code))
	}

}

// TestUserLoginFailure tests the following scenario: the client sends an access token request to the api gateway,
// passing username "admin" and password "admin_wrong_pass" in the body.
// The gateway makes an http GET request with the given information and sends it to the user management microservice.
// It is assumed that the user "admin" exists but has password "admin_pass", so the authentication should not succeed
// and the gateway should respond to the client with a 401 http status code.
func TestUserLoginFailure(t *testing.T) {

	config.SetConfigurationFromFile("../config/config-test.json")

	jsonBody := simplejson.New()
	jsonBody.Set("username", "admin")
	jsonBody.Set("password", "admin_wrong_pass")

	requestBody, _ := jsonBody.MarshalJSON()
	request, _ := http.NewRequest(http.MethodPost, "/didattica-mobile/api/v1.0/token", bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	handler := createTestGatewayLoginUser()
	// a goroutine representing the microservice listens to the requests coming from the api gateway
	go mock.LaunchUserManagementMock()
	// simulates a request-response interaction between client and api gateway
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Error("Expected 401 Not Found but got " + strconv.Itoa(response.Code) + " " + http.StatusText(response.Code))
	}
}
