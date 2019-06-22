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

// createTestPushNotification creates an http handler that handles the test requests
func createTestGatewayPushNotification() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/didattica-mobile/api/v1.0/notification/course/{courseId}", microservice.PushCourseNotification).Methods(http.MethodPost)
	return r
}

// TestRegisterUserSuccess tests the following scenario: the client makes a request for pushing a notification about a
// course and the notification is correctly pushed. So the response should be 200 OK
func TestPushNotificationSuccess(t *testing.T) {

	config.SetConfigurationFromFile("../config/config-test.json")

	jsonBody := simplejson.New()
	jsonBody.Set("message", "courseMessage")

	// generate a token to be appended to the course creation request
	user := microservice.User{Name: "nome", Surname: "cognome", Username: "student", Password: "password", Type: "teacher", Mail: "name@example.com"}
	token, _ := microservice.GenerateAccessToken(user, []byte(config.Configuration.TokenPrivateKey))

	requestBody, _ := jsonBody.MarshalJSON()
	request, _ := http.NewRequest(http.MethodPost, "/didattica-mobile/api/v1.0/notification/course/courseId", bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	request.AddCookie(&http.Cookie{Name: "token", Value: token})

	response := httptest.NewRecorder()
	handler := createTestGatewayPushNotification()
	// a goroutine representing the microservice listens to the requests coming from the api gateway
	go mock.LaunchCourseManagementMock()
	// simulates a request-response interaction between client and api gateway
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Error("Expected 200 Ok but got " + strconv.Itoa(response.Code) + " " + http.StatusText(response.Code))
	}
}

// TestPushNotificationPermissionDenied tests the following scenario: the notification push cannot be done because the
// client is not a student. So the response should be 401 Unauthorized
func TestPushNotificationPermissionDenied(t *testing.T) {
	config.SetConfigurationFromFile("../config/config-test.json")

	jsonBody := simplejson.New()
	jsonBody.Set("message", "courseMessage")

	// generate a token to be appended to the course creation request
	user := microservice.User{Name: "nome", Surname: "cognome", Username: "student", Password: "password", Type: "student", Mail: "name@example.com"}
	token, _ := microservice.GenerateAccessToken(user, []byte(config.Configuration.TokenPrivateKey))

	requestBody, _ := jsonBody.MarshalJSON()
	request, _ := http.NewRequest(http.MethodPost, "/didattica-mobile/api/v1.0/notification/course/courseId", bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	request.AddCookie(&http.Cookie{Name: "token", Value: token})

	response := httptest.NewRecorder()
	handler := createTestGatewayPushNotification()
	// a goroutine representing the microservice listens to the requests coming from the api gateway
	go mock.LaunchCourseManagementMock()
	// simulates a request-response interaction between client and api gateway
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Error("Expected 401 Unauthorized but got " + strconv.Itoa(response.Code) + " " + http.StatusText(response.Code))
	}
}
