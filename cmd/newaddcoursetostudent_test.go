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

// createTestGatewayAddCourseToStudent creates an http handler that handles the test requests
func createTestGatewayAddCourseToStudent() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/didattica-mobile/api/v1.0/students/{username}",
		microservice.AddCourseToStudent).Methods(http.MethodPut)
	return r
}

// TestAddCourseToStudentSuccessNew tests the following scenario: the client requests to add a course to an existing
// student. The operation succeeds in both course management and notification management micro-service.
// Therefore, it is expected that the response to the client is 200 OK*/
func TestAddCourseToStudentSuccess(t *testing.T) {

	_ = config.SetConfigurationFromFile("../config/config-test.json")

	// generate a token to be appended to the request
	user := microservice.User{Name: "name", Surname: "surname", Username: "existingUser",
		Password: "pass", Type: "student", Mail: "name@example.com"}
	token, _ := microservice.GenerateAccessToken(user, []byte(config.Configuration.TokenPrivateKey))

	// make the body of request containing the course to add to student
	jsonBody := simplejson.New()
	jsonBody.Set("id", "idCourseSuccess")
	jsonBody.Set("name", "courseSuccess")
	jsonBody.Set("department", "department")
	jsonBody.Set("year", "2019-2020")
	requestBody, _ := json.Marshal(jsonBody)

	// make the PUT request for the course append
	request, _ := http.NewRequest(http.MethodPut,
		"/didattica-mobile/api/v1.0/students/existingUser", bytes.NewBuffer(requestBody))
	request.AddCookie(&http.Cookie{Name: "token", Value: token})
	response := httptest.NewRecorder()
	handler := createTestGatewayAddCourseToStudent()

	// Goroutines represent the micro-services listens to the requests coming from the api gateway
	go mock.LaunchCourseManagementMock()
	go mock.LaunchNotificationManagementMock()

	// simulates a request-response interaction between client and api gateway
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Error("Expected 200 Ok but got " + strconv.Itoa(response.Code) + " " + http.StatusText(response.Code))
	}
}

/*TestAddCourseToStudentSuccessNotExistentStudent tests the following scenario: the client requests to add a course to
a not existing student. It is expected that the Api gateway provides student creation and the final response is 200 OK*/
func TestAddCourseToStudentSuccessNotExistentStudent(t *testing.T) {

	_ = config.SetConfigurationFromFile("../config/config-test.json")

	// generate a token to be appended to the request
	user := microservice.User{Name: "name", Surname: "surname", Username: "notExistingStudent",
		Password: "pass", Type: "student", Mail: "name@example.com"}
	token, _ := microservice.GenerateAccessToken(user, []byte(config.Configuration.TokenPrivateKey))

	// make the body of request containing the course to add to student
	jsonBody := simplejson.New()
	jsonBody.Set("id", "idCourseSuccess")
	jsonBody.Set("name", "courseSuccess")
	jsonBody.Set("department", "department")
	jsonBody.Set("year", "2019-2020")
	requestBody, _ := json.Marshal(jsonBody)

	// make the PUT request for the course append
	request, _ := http.NewRequest(http.MethodPut,
		"/didattica-mobile/api/v1.0/students/notExistingStudent", bytes.NewBuffer(requestBody))
	request.AddCookie(&http.Cookie{Name: "token", Value: token})

	response := httptest.NewRecorder()
	handler := createTestGatewayAddCourseToStudent()

	// Goroutines represent the micro-services listens to the requests coming from the api gateway
	go mock.LaunchCourseManagementMock()
	go mock.LaunchNotificationManagementMock()

	// simulates a request-response interaction between client and api gateway
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Error("Expected 200 Ok but got " + strconv.Itoa(response.Code) + " " + http.StatusText(response.Code))
	}
}

/*
TestAddCourseToStudentFailCourseManagement tests the following scenario: a student requests to be subscribed for a
course. The subscribing fails in course management micro-service while succeeds in notification management micro-service.
The operation is not completed by both micro-services so the distributed transaction fails and the subscription is
removed from notification management micro-service to maintain consistency.
*/
func TestAddCourseToStudentFailCourseManagement(t *testing.T) {
	_ = config.SetConfigurationFromFile("../config/config-test.json")

	// generate a token to be appended to the request
	user := microservice.User{Name: "name", Surname: "surname", Username: "user",
		Password: "pass", Type: "student", Mail: "name@example.com"}
	token, _ := microservice.GenerateAccessToken(user, []byte(config.Configuration.TokenPrivateKey))

	// make the body of request containing the course to add to student
	jsonBody := simplejson.New()
	jsonBody.Set("id", "idCourseFailingInCourseManagement")
	jsonBody.Set("name", "courseFailingInCourseManagement")
	jsonBody.Set("department", "department")
	jsonBody.Set("year", "2019-2020")
	requestBody, _ := json.Marshal(jsonBody)

	// make the PUT request for the course append
	request, _ := http.NewRequest(http.MethodPut,
		"/didattica-mobile/api/v1.0/students/user", bytes.NewBuffer(requestBody))
	request.AddCookie(&http.Cookie{Name: "token", Value: token})

	response := httptest.NewRecorder()
	handler := createTestGatewayAddCourseToStudent()

	// Goroutines represent the micro-services listens to the requests coming from the api gateway
	go mock.LaunchCourseManagementMock()
	go mock.LaunchNotificationManagementMock()

	// simulates a request-response interaction between client and api gateway
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Error("Expected 400 Bad request but got " +
			strconv.Itoa(response.Code) + " " + http.StatusText(response.Code))
	}
}

/*
TestAddCourseToStudentFailNotificationManagement tests the following scenario: a student requests to be subscribed for a
course. The subscribing fails in notification management micro-service while succeeds in course management micro-service.
The operation is not completed by both micro-services so the distributed transaction fails and the subscription is
removed from course management micro-service to maintain consistency.
*/
func TestAddCourseToStudentFailNotificationManagement(t *testing.T) {
	_ = config.SetConfigurationFromFile("../config/config-test.json")

	// generate a token to be appended to the request
	user := microservice.User{Name: "name", Surname: "surname", Username: "user",
		Password: "pass", Type: "student", Mail: "name@example.com"}
	token, _ := microservice.GenerateAccessToken(user, []byte(config.Configuration.TokenPrivateKey))

	// make the body of request containing the course to add to student
	jsonBody := simplejson.New()
	jsonBody.Set("id", "idCourseFailingInNotificationManagement")
	jsonBody.Set("name", "courseFailingInNotificationManagement")
	jsonBody.Set("department", "department")
	jsonBody.Set("year", "2019-2020")
	requestBody, _ := json.Marshal(jsonBody)

	// make the PUT request for the course append
	request, _ := http.NewRequest(http.MethodPut,
		"/didattica-mobile/api/v1.0/students/user", bytes.NewBuffer(requestBody))
	request.AddCookie(&http.Cookie{Name: "token", Value: token})

	response := httptest.NewRecorder()
	handler := createTestGatewayAddCourseToStudent()

	// Goroutines represent the micro-services listens to the requests coming from the api gateway
	go mock.LaunchCourseManagementMock()
	go mock.LaunchNotificationManagementMock()

	// simulates a request-response interaction between client and api gateway
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Error("Expected 400 Bad request but got " +
			strconv.Itoa(response.Code) + " " + http.StatusText(response.Code))
	}
}
