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

// createTestGatewayCreateExam creates an http handler that handles the test requests
func createTestGatewayCreateExam() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/didattica-mobile/api/v1.0/exams", microservice.CreateExam).Methods(http.MethodPost)
	return r
}

/*TestCreateExamSuccess tests the following scenario: the client makes a exam creation request providing a correct
token to the api gateway: it has not expired and it has been signed with the correct key. Furthermore, the user type is
teacher, so the creation operation is allowed. The Api Gateway should then forward the request to the microservice
and return the response to the client*/
func TestCreateExamSuccess(t *testing.T) {

	config.SetConfigurationFromFile("../config/config-test.json")

	// build the information of the exam to be created
	jsonBody := simplejson.New()
	jsonBody.Set("course", "IdCorso")
	jsonBody.Set("call", 1)
	jsonBody.Set("date", "21-03-2019")
	jsonBody.Set("startTime", "10:30")
	jsonBody.Set("room", "A2")
	jsonBody.Set("expirationDate", "20-03-2019")

	// generate a token to be appended to the exam creation request
	user := microservice.User{Name: "nome", Surname: "cognome", Username: "username", Password: "password", Type: "teacher", Mail: "name@example.com"}
	token, _ := microservice.GenerateAccessToken(user, []byte(config.Configuration.TokenPrivateKey))

	// make the POST request for the exam creation
	requestBody, _ := jsonBody.MarshalJSON()
	request, _ := http.NewRequest(http.MethodPost, "/didattica-mobile/api/v1.0/exams", bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	request.AddCookie(&http.Cookie{Name: "token", Value: token})

	response := httptest.NewRecorder()
	handler := createTestGatewayCreateExam()
	// a goroutine representing the microservice listens to the requests coming from the api gateway
	go mock.LaunchCourseManagementMock()
	// simulates a request-response interaction between client and api gateway
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Error("Expected 200 Ok but got " + strconv.Itoa(response.Code) + " " + http.StatusText(response.Code))
	}

}

/*TestCreateExamFail tests the following scenario: a student makes an exam creation request. This operation is allowed
for the teachers only, therefore the Api Gateway should respond with Unauthorized.*/
func TestCreateExamNotAllowed(t *testing.T) {

	config.SetConfigurationFromFile("../config/config-test.json")

	// build the information of the exam to be created
	jsonBody := simplejson.New()
	jsonBody.Set("course", "IdCorso")
	jsonBody.Set("call", 1)
	jsonBody.Set("date", "21-03-2019")
	jsonBody.Set("startTime", "10:30")
	jsonBody.Set("room", "A2")
	jsonBody.Set("expirationDate", "20-03-2019")

	// generate a token to be appended to the exam creation request
	user := microservice.User{Name: "nome", Surname: "cognome", Username: "username", Password: "password", Type: "student", Mail: "name@example.com"}
	token, _ := microservice.GenerateAccessToken(user, []byte(config.Configuration.TokenPrivateKey))

	// make the POST request for the exam creation
	requestBody, _ := jsonBody.MarshalJSON()
	request, _ := http.NewRequest(http.MethodPost, "/didattica-mobile/api/v1.0/exams", bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	request.AddCookie(&http.Cookie{Name: "token", Value: token})

	response := httptest.NewRecorder()
	handler := createTestGatewayCreateExam()
	// a goroutine representing the microservice listens to the requests coming from the api gateway
	go mock.LaunchCourseManagementMock()
	// simulates a request-response interaction between client and api gateway
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Error("Expected 401 Unauthorized but got " + strconv.Itoa(response.Code) + " " + http.StatusText(response.Code))
	}

}

/*TestCreateExamWithTokenBadSigned tests the following scenario: the client makes a course creation request containing an
access token that is not valid according to the signing key used by the Api Gateway. The gateway should respond with 401*/

func TestCreateExamWithTokenBadSigned(t *testing.T) {
	config.SetConfigurationFromFile("../config/config-test.json")

	// build the information of the exam to be created
	jsonBody := simplejson.New()
	jsonBody.Set("course", "IdCorso")
	jsonBody.Set("call", 1)
	jsonBody.Set("date", "21-03-2019")
	jsonBody.Set("startTime", "10:30")
	jsonBody.Set("room", "A2")
	jsonBody.Set("expirationDate", "20-03-2019")

	// generate a token to be appended to the exam creation request (with wrong secret)
	user := microservice.User{Name: "nome", Surname: "cognome", Username: "username", Password: "password", Type: "teacher", Mail: "name@example.com"}
	token, _ := microservice.GenerateAccessToken(user, []byte("wrong-signing-key"))

	// make the POST request for the exam creation
	requestBody, _ := jsonBody.MarshalJSON()
	request, _ := http.NewRequest(http.MethodPost, "/didattica-mobile/api/v1.0/exams", bytes.NewBuffer(requestBody))
	request.Header.Set("Content-Type", "application/json")
	request.AddCookie(&http.Cookie{Name: "token", Value: token})

	response := httptest.NewRecorder()
	handler := createTestGatewayCreateExam()
	// a goroutine representing the microservice listens to the requests coming from the api gateway
	go mock.LaunchCourseManagementMock()
	// simulates a request-response interaction between client and api gateway
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Error("Expected 401 Unauthorized but got " + strconv.Itoa(response.Code) + " " + http.StatusText(response.Code))
	}
}
