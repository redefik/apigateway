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

//TODO cambiare eventualmente il nome del modulo

// createTestGatewayCreateCourse creates an http handler that handles the test requests
func createTestGatewayCreateCourse() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/didattica-mobile/api/v1.0/courses", microservice.CreateCourse).Methods(http.MethodPost)
	return r
}

/*TestCreateCourseSuccess tests the following scenario: the client makes a course creation request providing a correct
token to the api gateway: it has not expired and it has been signed with the correct key. Furthermore, the user type is
teacher, so the creation operation is allowed. The Api Gateway should then forward the request to the microservice
and return the response to the client*/
func TestCreateCourseSuccess(t *testing.T) {

	config.SetConfiguration("../config/config-test.json")

	// build the information of the course to be created
	jsonBody := simplejson.New()
	jsonBody.Set("name", "corso")
	jsonBody.Set("department", "dipartimento")
	jsonBody.Set("teacher", "docente")
	jsonBody.Set("year", "2019-2020")
	jsonBody.Set("semester", 2)
	jsonBody.Set("description", "descrizione")
	// The course schedule is omitted for sake of simplicity

	// generate a token to be appended to the course creation request
	user := microservice.User{Name: "nome", Surname: "cognome", Username: "username", Password: "password", Type: "teacher"}
	token, _ := microservice.GenerateAccessToken(user, []byte(config.Configuration.TokenPrivateKey))

	// make the POST request for the course creation
	requestBody, _ := jsonBody.MarshalJSON()
	request, _ := http.NewRequest(http.MethodPost, "/didattica-mobile/api/v1.0/courses", bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	request.AddCookie(&http.Cookie{Name: "token", Value: token})

	response := httptest.NewRecorder()
	handler := createTestGatewayCreateCourse()
	// a goroutine representing the microservice listens to the requests coming from the api gateway
	go mock.LaunchCourseManagementMock()
	// simulates a request-response interaction between client and api gateway
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Error("Expected 200 Ok but got " + strconv.Itoa(response.Code) + " " + http.StatusText(response.Code))
	}

}

/*TestCreateCourseFail tests the following scenario: a student makes a course creation request. This operation is allowed
for the teachers only, therefore the Api Gateway should respond with Unauthorized.*/
func TestCreateCourseNotAllowed(t *testing.T) {

	config.SetConfiguration("../config/config-test.json")

	// build the information of the course to be created
	jsonBody := simplejson.New()
	jsonBody.Set("name", "corso")
	jsonBody.Set("department", "dipartimento")
	jsonBody.Set("teacher", "docente")
	jsonBody.Set("year", "2019-2020")
	jsonBody.Set("semester", 2)
	jsonBody.Set("description", "descrizione")
	// The course schedule is omitted for sake of simplicity

	// generate a token to be appended to the course creation request
	user := microservice.User{Name: "nome", Surname: "cognome", Username: "username", Password: "password", Type: "student"}
	token, _ := microservice.GenerateAccessToken(user, []byte(config.Configuration.TokenPrivateKey))

	// make the POST request for the course creation
	requestBody, _ := jsonBody.MarshalJSON()
	request, _ := http.NewRequest(http.MethodPost, "/didattica-mobile/api/v1.0/courses", bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	request.AddCookie(&http.Cookie{Name: "token", Value: token})

	response := httptest.NewRecorder()
	handler := createTestGatewayCreateCourse()
	// a goroutine representing the microservice listens to the requests coming from the api gateway
	go mock.LaunchCourseManagementMock()
	// simulates a request-response interaction between client and api gateway
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Error("Expected 401 Unauthorized but got " + strconv.Itoa(response.Code) + " " + http.StatusText(response.Code))
	}

}

/*TestProcessTokenWithBadSign tests the following scenario: the client makes a course creation request containing an
access token that is not valid according to the signing key used by the Api Gateway. The gateway should respond with 401*/
func TestProcessTokenWithBadSign(t *testing.T) {
	config.SetConfiguration("../config/config-test.json")

	// build the information of the course to be created
	jsonBody := simplejson.New()
	jsonBody.Set("name", "corso")
	jsonBody.Set("department", "dipartimento")
	jsonBody.Set("teacher", "docente")
	jsonBody.Set("year", "2019-2020")
	jsonBody.Set("semester", 2)
	jsonBody.Set("description", "descrizione")
	// The course schedule is omitted for sake of simplicity

	// generate a token to be appended to the course creation request
	user := microservice.User{Name: "nome", Surname: "cognome", Username: "username", Password: "password", Type: "teacher"}
	token, _ := microservice.GenerateAccessToken(user, []byte("wrong-signing-key"))

	// make the POST request for the course creation
	requestBody, _ := jsonBody.MarshalJSON()
	request, _ := http.NewRequest(http.MethodPost, "/didattica-mobile/api/v1.0/courses", bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	request.AddCookie(&http.Cookie{Name: "token", Value: token})

	response := httptest.NewRecorder()
	handler := createTestGatewayCreateCourse()
	// a goroutine representing the microservice listens to the requests coming from the api gateway
	go mock.LaunchCourseManagementMock()
	// simulates a request-response interaction between client and api gateway
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Error("Expected 401 Unauthorized but got " + strconv.Itoa(response.Code) + " " + http.StatusText(response.Code))
	}
}
