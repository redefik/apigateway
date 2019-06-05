package main

import (
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
func createTestGatewaySearchCourse() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/didattica-mobile/api/v1.0/courses/{by}/{string}", microservice.FindCourse).Methods(http.MethodGet)
	return r
}

// TestSearchCourseSuccess tests the following scenario: the client sends a course research request to the api gateway,
// passing the type of research (name or teacher) and the string representing the sequence to find.
// The gateway makes an http GET request with the given information and sends it to the user management microservice.
// It is assumed that exist a course with name matching with sequence "seq", so in this case the research has success.
func TestSearchCourseSuccess(t *testing.T) {

	_ = config.SetConfiguration("../config/config-test.json")

	// generate a token to be appended to the request
	user := microservice.User{Name: "nome", Surname: "cognome", Username: "username", Password: "password", Type: "teacher"}
	token, _ := microservice.GenerateAccessToken(user, []byte(config.Configuration.TokenPrivateKey))

	// Make the get request for course searching
	request, _ := http.NewRequest(http.MethodGet, "/didattica-mobile/api/v1.0/courses/name/seq", nil)
	request.Header.Set("Content-Type", "application/json")
	request.AddCookie(&http.Cookie{Name: "token", Value: token})

	response := httptest.NewRecorder()
	handler := createTestGatewaySearchCourse()
	// a goroutine representing the microservice listens to the requests coming from the api gateway
	go mock.LaunchCourseManagementMock()
	// simulates a request-response interaction between client and api gateway
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Error("Expected 200 OK but got " + strconv.Itoa(response.Code) + " " + http.StatusText(response.Code))
	}

}

// TestSearchCourseNotSuccess tests the following scenario: the client sends a course research request to the api gateway,
// passing the type of research (name or teacher) and the string representing the sequence to find.
// The gateway makes an http GET request with the given information and sends it to the user management microservice.
// It is assumed that not exist a course with teacher's name matching with sequence "seq", so in this case the research failed.
func TestSearchCourseNotSuccess(t *testing.T) {

	_ = config.SetConfiguration("../config/config-test.json")

	// generate a token to be appended to the request
	user := microservice.User{Name: "nome", Surname: "cognome", Username: "username", Password: "password", Type: "teacher"}
	token, _ := microservice.GenerateAccessToken(user, []byte(config.Configuration.TokenPrivateKey))

	// Make the get request for course searching
	request, _ := http.NewRequest(http.MethodGet, "/didattica-mobile/api/v1.0/courses/teacher/seq", nil)
	request.Header.Set("Content-Type", "application/json")
	request.AddCookie(&http.Cookie{Name: "token", Value: token})

	response := httptest.NewRecorder()
	handler := createTestGatewaySearchCourse()
	// a goroutine representing the microservice listens to the requests coming from the api gateway
	go mock.LaunchCourseManagementMock()
	// simulates a request-response interaction between client and api gateway
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Error("Expected 404 Not Found but got " + strconv.Itoa(response.Code) + " " + http.StatusText(response.Code))
	}

}

// TestSearchCourseFailure tests the following scenario: the client sends a course research request to the api gateway,
// passing the type of research (not name or teacher) and the string representing the sequence to find.
// The gateway makes an http GET request with the given information and sends it to the user management microservice.
// Because the type does not match with "name" or "teacher" the research failed.
func TestSearchCourseFailure(t *testing.T) {

	_ = config.SetConfiguration("../config/config-test.json")

	// generate a token to be appended to the request
	user := microservice.User{Name: "nome", Surname: "cognome", Username: "username", Password: "password", Type: "teacher"}
	token, _ := microservice.GenerateAccessToken(user, []byte(config.Configuration.TokenPrivateKey))

	// Make the get request for course searching
	request, _ := http.NewRequest(http.MethodGet, "/didattica-mobile/api/v1.0/courses/notvalid/seq", nil)
	request.Header.Set("Content-Type", "application/json")
	request.AddCookie(&http.Cookie{Name: "token", Value: token})

	response := httptest.NewRecorder()
	handler := createTestGatewaySearchCourse()
	// a goroutine representing the microservice listens to the requests coming from the api gateway
	go mock.LaunchCourseManagementMock()
	// simulates a request-response interaction between client and api gateway
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Error("Expected 400 Bad Request but got " + strconv.Itoa(response.Code) + " " + http.StatusText(response.Code))
	}

}
