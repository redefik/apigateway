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

// newCreateTestGatewayCreateCourse creates an http handler that handles the test requests
func createTestGatewayCreateCourse() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/didattica-mobile/api/v1.0/courses", microservice.CreateCourse).Methods(http.MethodPost)
	return r
}

/* TestCreateCourseSuccessNew tests the following scenario: the client makes a course creation request providing a correct
token to the api gateway: it has not expired and it has been signed with the correct key. Furthermore, the user type is
teacher, so the creation operation is allowed. Both course management that notification management micro-services succeed
in registration of course so the distributed transaction has success. The Api Gateway should then return a positive
response to the client*/
func TestCreateCourseSuccessNew(t *testing.T) {

	_ = config.SetConfigurationFromFile("../config/config-test.json")

	// build the information of the course to be created (in a simplified way)
	jsonBody := simplejson.New()
	jsonBody.Set("name", "courseSuccess")

	// generate a token to be appended to the course creation request
	user := microservice.User{Name: "nome", Surname: "cognome", Username: "username", Password: "password", Type: "teacher", Mail: "name@example.com"}
	token, _ := microservice.GenerateAccessToken(user, []byte(config.Configuration.TokenPrivateKey))

	// make the POST request for the course creation
	requestBody, _ := jsonBody.MarshalJSON()
	request, _ := http.NewRequest(http.MethodPost, "/didattica-mobile/api/v1.0/courses", bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	request.AddCookie(&http.Cookie{Name: "token", Value: token})

	response := httptest.NewRecorder()
	handler := createTestGatewayCreateCourse()
	// Goroutines represent the micro-services listening to the requests coming from the api gateway
	go mock.LaunchCourseManagementMock()
	go mock.LaunchNotificationManagementMock()
	// simulates a request-response interaction between client and api gateway
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Error("Expected 201 Ok but got " + strconv.Itoa(response.Code) + " " + http.StatusText(response.Code))
	}

}

/*TestCreateCourseNotAllowedNew tests the following scenario: a student makes a course creation request. This operation is allowed
for the teachers only, therefore the Api Gateway should respond with Unauthorized.*/
func TestCreateCourseNotAllowedNew(t *testing.T) {

	config.SetConfigurationFromFile("../config/config-test.json")

	// build the information of the course to be created (in a simplified way)
	jsonBody := simplejson.New()
	jsonBody.Set("name", "corso")

	// generate a token to be appended to the course creation request
	user := microservice.User{Name: "nome", Surname: "cognome", Username: "username", Password: "password", Type: "student", Mail: "name@example.com"}
	token, _ := microservice.GenerateAccessToken(user, []byte(config.Configuration.TokenPrivateKey))

	// make the POST request for the course creation
	requestBody, _ := jsonBody.MarshalJSON()
	request, _ := http.NewRequest(http.MethodPost, "/didattica-mobile/api/v1.0/courses", bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	request.AddCookie(&http.Cookie{Name: "token", Value: token})

	response := httptest.NewRecorder()
	handler := createTestGatewayCreateCourse()
	// Goroutines represent the micro-services listening to the requests coming from the api gateway
	go mock.LaunchCourseManagementMock()
	go mock.LaunchNotificationManagementMock()
	// simulates a request-response interaction between client and api gateway
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Error("Expected 401 Unauthorized but got " + strconv.Itoa(response.Code) + " " + http.StatusText(response.Code))
	}

}

/*TestProcessTokenWithBadSignNew tests the following scenario: the client makes a course creation request containing an
access token that is not valid according to the signing key used by the Api Gateway. The gateway should respond with 401*/
func TestProcessTokenWithBadSignNew(t *testing.T) {
	config.SetConfigurationFromFile("../config/config-test.json")

	// build the information of the course to be created (in a simplified way)
	jsonBody := simplejson.New()
	jsonBody.Set("name", "corso")

	// generate a token to be appended to the course creation request
	user := microservice.User{Name: "nome", Surname: "cognome", Username: "username", Password: "password", Type: "teacher", Mail: "name@example.com"}
	token, _ := microservice.GenerateAccessToken(user, []byte("wrong-signing-key"))

	// make the POST request for the course creation
	requestBody, _ := jsonBody.MarshalJSON()
	request, _ := http.NewRequest(http.MethodPost, "/didattica-mobile/api/v1.0/courses", bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	request.AddCookie(&http.Cookie{Name: "token", Value: token})

	response := httptest.NewRecorder()
	handler := createTestGatewayCreateCourse()
	// Goroutines represent the micro-services listening to the requests coming from the api gateway
	go mock.LaunchCourseManagementMock()
	go mock.LaunchNotificationManagementMock()
	// simulates a request-response interaction between client and api gateway
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Error("Expected 401 Unauthorized but got " + strconv.Itoa(response.Code) + " " + http.StatusText(response.Code))
	}
}

/* TestCreateCourseFailureCourseManagement tests the following scenario: the client makes a course creation request
providing a correct token to the api gateway: it has not expired and it has been signed with the correct key.
Furthermore, the user type is teacher, so the creation operation is allowed. The creation fail in course management
micro-service and succeed in notification management micro-service. As a consequence the distributed transaction fail
and the just created course is deleted from the data-store of notification management. Client receive a http 500 internal
server error */
func TestCreateCourseFailureCourseManagement(t *testing.T) {

	_ = config.SetConfigurationFromFile("../config/config-test.json")

	// build the information of the course to be created (in a simplified way)
	jsonBody := simplejson.New()
	jsonBody.Set("name", "courseFailInCourseManagement")

	// generate a token to be appended to the course creation request
	user := microservice.User{Name: "nome", Surname: "cognome", Username: "username", Password: "password", Type: "teacher", Mail: "name@example.com"}
	token, _ := microservice.GenerateAccessToken(user, []byte(config.Configuration.TokenPrivateKey))

	// make the POST request for the course creation
	requestBody, _ := jsonBody.MarshalJSON()
	request, _ := http.NewRequest(http.MethodPost, "/didattica-mobile/api/v1.0/courses", bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	request.AddCookie(&http.Cookie{Name: "token", Value: token})

	response := httptest.NewRecorder()
	handler := createTestGatewayCreateCourse()
	// Goroutines represent the micro-services listening to the requests coming from the api gateway
	go mock.LaunchCourseManagementMock()
	go mock.LaunchNotificationManagementMock()
	// simulates a request-response interaction between client and api gateway
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Error("Expected 500 Ok but got " + strconv.Itoa(response.Code) + " " + http.StatusText(response.Code))
	}

}

/* TestCreateCourseFailureNotificationManagement tests the following scenario: the client makes a course creation request
providing a correct token to the api gateway: it has not expired and it has been signed with the correct key.
Furthermore, the user type is teacher, so the creation operation is allowed. The creation fail in notification management
micro-service and succeed in course management micro-service. As a consequence the distributed transaction fail
and the just created course is deleted from the data-store of course management. Client receive a http 500 internal
server error */
func TestCreateCourseFailureNotificationManagement(t *testing.T) {

	_ = config.SetConfigurationFromFile("../config/config-test.json")

	// build the information of the course to be created (in a simplified way)
	jsonBody := simplejson.New()
	jsonBody.Set("name", "courseFailInNotificationManagement")

	// generate a token to be appended to the course creation request
	user := microservice.User{Name: "nome", Surname: "cognome", Username: "username", Password: "password", Type: "teacher", Mail: "name@example.com"}
	token, _ := microservice.GenerateAccessToken(user, []byte(config.Configuration.TokenPrivateKey))

	// make the POST request for the course creation
	requestBody, _ := jsonBody.MarshalJSON()
	request, _ := http.NewRequest(http.MethodPost, "/didattica-mobile/api/v1.0/courses", bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	request.AddCookie(&http.Cookie{Name: "token", Value: token})

	response := httptest.NewRecorder()
	handler := createTestGatewayCreateCourse()
	// Goroutines represent the micro-services listening to the requests coming from the api gateway
	go mock.LaunchCourseManagementMock()
	go mock.LaunchNotificationManagementMock()
	// simulates a request-response interaction between client and api gateway
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Error("Expected 500 Ok but got " + strconv.Itoa(response.Code) + " " + http.StatusText(response.Code))
	}

}
