package searchExam

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

// createTestGatewaySearchExam creates an http handler that handles the test requests
func createTestGatewaySearchExam() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/didattica-mobile/api/v1.0/exams/{course}", microservice.FindExamByCourse).Methods(http.MethodGet)
	return r
}

// TestSearchExamSuccess tests the following scenario: the client sends a course research request to the api gateway,
// passing the id of the course the exam to find belong to. The gateway makes an http GET request with the given
// information and sends it to the user management microservice. It is assumed that exist an exam with course field
// matching with id "idSuccess", so in this case the research has success.
func TestSearchExamSuccess(t *testing.T) {

	_ = config.SetConfigurationFromFile("../../../config/config-test.json")

	// generate a token to be appended to the course creation request
	user := microservice.User{Name: "nome", Surname: "cognome", Username: "username", Password: "password", Type: "teacher", Mail: "name@example.com"}
	token, _ := microservice.GenerateAccessToken(user, []byte(config.Configuration.TokenPrivateKey))

	// Make the get request for course searching
	request, _ := http.NewRequest(http.MethodGet, "/didattica-mobile/api/v1.0/exams/idSuccess", nil)
	request.Header.Set("Content-Type", "application/json")
	request.AddCookie(&http.Cookie{Name: "token", Value: token})

	response := httptest.NewRecorder()
	handler := createTestGatewaySearchExam()
	// a goroutine representing the microservice listens to the requests coming from the api gateway
	go mock.LaunchCourseManagementMock()
	// simulates a request-response interaction between client and api gateway
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Error("Expected 200 OK but got " + strconv.Itoa(response.Code) + " " + http.StatusText(response.Code))
	}

}

// TestSearchExamNotSuccess tests the following scenario: the client sends a course research request to the api gateway,
// passing the id of the course the exam to find belong to. The gateway makes an http GET request with the given
// information and sends it to the user management microservice. It is assumed that not exist an exam with course field
// matching with id "idFailure", so in this case the research has no success.
func TestSearchExamNotSuccess(t *testing.T) {

	_ = config.SetConfigurationFromFile("../../../config/config-test.json")

	// generate a token to be appended to the course creation request
	user := microservice.User{Name: "nome", Surname: "cognome", Username: "username", Password: "password", Type: "teacher", Mail: "name@example.com"}
	token, _ := microservice.GenerateAccessToken(user, []byte(config.Configuration.TokenPrivateKey))

	// Make the get request for course searching
	request, _ := http.NewRequest(http.MethodGet, "/didattica-mobile/api/v1.0/exams/idFailure", nil)
	request.Header.Set("Content-Type", "application/json")
	request.AddCookie(&http.Cookie{Name: "token", Value: token})

	response := httptest.NewRecorder()
	handler := createTestGatewaySearchExam()
	// a goroutine representing the microservice listens to the requests coming from the api gateway
	go mock.LaunchCourseManagementMock()
	// simulates a request-response interaction between client and api gateway
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Error("Expected 404 Not Found but got " + strconv.Itoa(response.Code) + " " + http.StatusText(response.Code))
	}

}
