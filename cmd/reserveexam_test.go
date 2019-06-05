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

// reserveExam creates an http handler that handles the test requests
func createTestGatewayReserveExam() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/didattica-mobile/api/v1.0/exams/{examId}/students/{studentUsername}", microservice.ReserveExam).Methods(http.MethodPut)
	return r
}

/*TestReserveExamSuccess tests the following scenario: the client manages to make an exam reservation for the provided student.
Therefore, it is expected that the response to the client is 200 OK*/
func TestReserveExamSuccess(t *testing.T) {

	config.SetConfiguration("../config/config-test.json")

	// generate a token to be appended to the request
	user := microservice.User{Name: "nome", Surname: "cognome", Username: "existent_student", Password: "password", Type: "student"}
	token, _ := microservice.GenerateAccessToken(user, []byte(config.Configuration.TokenPrivateKey))

	// make the PUT request for the exam reservation
	request, _ := http.NewRequest(http.MethodPut, "/didattica-mobile/api/v1.0/exams/existent_exam/students/existent_student", nil)
	request.AddCookie(&http.Cookie{Name: "token", Value: token})

	response := httptest.NewRecorder()
	handler := createTestGatewayReserveExam()
	// a goroutine representing the microservice listens to the requests coming from the api gateway
	go mock.LaunchCourseManagementMock()
	// simulates a request-response interaction between client and api gateway
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Error("Expected 200 Ok but got " + strconv.Itoa(response.Code) + " " + http.StatusText(response.Code))
	}

}

// TestReserveExamNotExistentStudent tests the following scenario: the client requires to make an exam reservation for
// a not existent student, then the response should be 404 Not Found
func TestReserveExamNotExistentStudent(t *testing.T) {
	config.SetConfiguration("../config/config-test.json")

	// generate a token to be appended to the request
	user := microservice.User{Name: "nome", Surname: "cognome", Username: "existent_student", Password: "password", Type: "student"}
	token, _ := microservice.GenerateAccessToken(user, []byte(config.Configuration.TokenPrivateKey))

	// make the PUT request for the exam reservation
	request, _ := http.NewRequest(http.MethodPut, "/didattica-mobile/api/v1.0/exams/existent_exam/students/not_existent_student", nil)
	request.AddCookie(&http.Cookie{Name: "token", Value: token})

	response := httptest.NewRecorder()
	handler := createTestGatewayReserveExam()
	// a goroutine representing the microservice listens to the requests coming from the api gateway
	go mock.LaunchCourseManagementMock()
	// simulates a request-response interaction between client and api gateway
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Error("Expected 404 Not Found but got " + strconv.Itoa(response.Code) + " " + http.StatusText(response.Code))
	}
}

// TestReserveExamNotExistentExam tests the following scenario: the client requires to make an exam reservation for
// a not existent exam, then the response should be 404 Not Found
func TestReserveExamNotExistentExam(t *testing.T) {
	config.SetConfiguration("../config/config-test.json")

	// generate a token to be appended to the request
	user := microservice.User{Name: "nome", Surname: "cognome", Username: "existent_student", Password: "password", Type: "student"}
	token, _ := microservice.GenerateAccessToken(user, []byte(config.Configuration.TokenPrivateKey))

	// make the PUT request for the exam reservation
	request, _ := http.NewRequest(http.MethodPut, "/didattica-mobile/api/v1.0/exams/not_existent_exam/students/existent_student", nil)
	request.AddCookie(&http.Cookie{Name: "token", Value: token})

	response := httptest.NewRecorder()
	handler := createTestGatewayReserveExam()
	// a goroutine representing the microservice listens to the requests coming from the api gateway
	go mock.LaunchCourseManagementMock()
	// simulates a request-response interaction between client and api gateway
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Error("Expected 404 Not Found but got " + strconv.Itoa(response.Code) + " " + http.StatusText(response.Code))
	}
}
