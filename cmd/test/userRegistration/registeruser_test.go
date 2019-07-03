package userRegistration

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
func createTestGatewayRegisterUser() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/didattica-mobile/api/v1.0/users", microservice.RegisterUser).Methods(http.MethodPost)
	return r
}

// TestRegisterUserSuccess tests the following scenario: the client sends a well-formed request to the api gateway, that
// forwards it to the user-management microservice. If the send fails the test does not pass.
func TestRegisterUserSuccess(t *testing.T) {

	config.SetConfigurationFromFile("../../../config/config-test.json")

	jsonBody := simplejson.New()
	jsonBody.Set("username", "user")
	jsonBody.Set("password", "pass")
	jsonBody.Set("name", "name")
	jsonBody.Set("surname", "surname")
	jsonBody.Set("fiscalCode", "code")

	requestBody, _ := jsonBody.MarshalJSON()
	request, _ := http.NewRequest(http.MethodPost, "/didattica-mobile/api/v1.0/users", bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()
	handler := createTestGatewayRegisterUser()
	// a goroutine representing the microservice listens to the requests coming from the api gateway
	go mock.LaunchUserManagementMock()
	// simulates a request-response interaction between client and api gateway
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Error("Expected 200 Ok but got " + strconv.Itoa(response.Code) + " " + http.StatusText(response.Code))
	}

}
