package main

import (
	"bytes"
	"encoding/json"
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

// createTestGatewayCourseUnsubscription creates an http handler that handles the test requests
func createTestGatewayCourseUnsubscription() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/didattica-mobile/api/v1.0/students/{username}",
		microservice.UnsubscribeStudentFromCourse).Methods(http.MethodDelete)
	return r
}

/*
TestCourseUnsubscribeSuccess tests the following scenario: the client send a request to unsubscribe a student from
a course. The api gateway forwards the request both to course management and notification management. The operation
succeeds on both micro-services so the client obtain an http 200 ok as response.
*/
func TestCourseUnsubscribeSuccess(t *testing.T) {

	_ = config.SetConfigurationFromFile("../config/config-test.json")

	// generate a token to be appended to the request
	user := microservice.User{Name: "name", Surname: "surname", Username: "user",
		Password: "pass", Type: "student", Mail: "name@example.com"}
	token, _ := microservice.GenerateAccessToken(user, []byte(config.Configuration.TokenPrivateKey))

	// make the body of request containing the course to add to student
	jsonBody := simplejson.New()
	jsonBody.Set("id", "idCourseToUnregisterSuccess")
	jsonBody.Set("name", "courseToUnregisterSuccess")
	jsonBody.Set("department", "department")
	jsonBody.Set("year", "2019-2020")
	requestBody, _ := json.Marshal(jsonBody)

	// make the Delete request
	request, _ := http.NewRequest(http.MethodDelete,
		"/didattica-mobile/api/v1.0/students/user", bytes.NewBuffer(requestBody))
	request.AddCookie(&http.Cookie{Name: "token", Value: token})

	response := httptest.NewRecorder()
	handler := createTestGatewayCourseUnsubscription()

	// Goroutines represent the micro-services listening to the requests coming from the api gateway
	go mock.LaunchCourseManagementMock()
	go mock.LaunchNotificationManagementMock()

	// simulates a request-response interaction between client and api gateway
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Error("Expected 200 Ok but got " + strconv.Itoa(response.Code) + " " + http.StatusText(response.Code))
	}
}

/*
TestCourseUnsubscribeFailureInCourseManagement tests the following scenario: the client send a request to unsubscribe a
student from a course. The api gateway forwards the request both to course management and notification management.
The operation succeeds in notification management micro-service and fails in the other one. So the distributed transaction
fails and the client get a 400 bad request.
*/
func TestCourseUnsubscribeFailureInCourseManagement(t *testing.T) {

	_ = config.SetConfigurationFromFile("../config/config-test.json")

	// generate a token to be appended to the request
	user := microservice.User{Name: "name", Surname: "surname", Username: "user",
		Password: "pass", Type: "student", Mail: "name@example.com"}
	token, _ := microservice.GenerateAccessToken(user, []byte(config.Configuration.TokenPrivateKey))

	// make the body of request containing the course to add to student
	jsonBody := simplejson.New()
	jsonBody.Set("id", "idCourseToUnregisterFailureInCourseManagement")
	jsonBody.Set("name", "courseToUnregisterFailureInCourseManagement")
	jsonBody.Set("department", "department")
	jsonBody.Set("year", "2019-2020")
	requestBody, _ := json.Marshal(jsonBody)

	// make the Delete request
	request, _ := http.NewRequest(http.MethodDelete,
		"/didattica-mobile/api/v1.0/students/user", bytes.NewBuffer(requestBody))
	request.AddCookie(&http.Cookie{Name: "token", Value: token})

	response := httptest.NewRecorder()
	handler := createTestGatewayCourseUnsubscription()

	// Goroutines represent the micro-services listening to the requests coming from the api gateway
	go mock.LaunchCourseManagementMock()
	go mock.LaunchNotificationManagementMock()

	// simulates a request-response interaction between client and api gateway
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Error("Expected 400 Bad request but got " + strconv.Itoa(response.Code) + " " + http.StatusText(response.Code))
	}
}

/*
TestCourseUnsubscribeFailureInNotificationManagement tests the following scenario: the client send a request to unsubscribe a
student from a course. The api gateway forwards the request both to course management and notification management.
The operation succeeds in course management micro-service and fails in the other one. So the distributed transaction
fails and the client get a 400 bad request.
*/
func TestCourseUnsubscribeFailureInNotificationManagement(t *testing.T) {

	_ = config.SetConfigurationFromFile("../config/config-test.json")

	// generate a token to be appended to the request
	user := microservice.User{Name: "name", Surname: "surname", Username: "user",
		Password: "pass", Type: "student", Mail: "name@example.com"}
	token, _ := microservice.GenerateAccessToken(user, []byte(config.Configuration.TokenPrivateKey))

	// make the body of request containing the course to add to student
	jsonBody := simplejson.New()
	jsonBody.Set("id", "idCourseToUnregisterFailureInNotificationManagement")
	jsonBody.Set("name", "courseToUnregisterFailureInNotificationManagement")
	jsonBody.Set("department", "department")
	jsonBody.Set("year", "2019-2020")
	requestBody, _ := json.Marshal(jsonBody)

	// make the Delete request
	request, _ := http.NewRequest(http.MethodDelete,
		"/didattica-mobile/api/v1.0/students/user", bytes.NewBuffer(requestBody))
	request.AddCookie(&http.Cookie{Name: "token", Value: token})

	response := httptest.NewRecorder()
	handler := createTestGatewayCourseUnsubscription()

	// Goroutines represent the micro-services listening to the requests coming from the api gateway
	go mock.LaunchCourseManagementMock()
	go mock.LaunchNotificationManagementMock()

	// simulates a request-response interaction between client and api gateway
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Error("Expected 400 Bad request but got " + strconv.Itoa(response.Code) + " " + http.StatusText(response.Code))
	}
}
