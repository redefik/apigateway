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

// createTestGatewayGetStudentCourses creates an http handler that handles the test requests
func createTestGatewayGetStudentCourses() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/didattica-mobile/api/v1.0/courses/students/{username}", microservice.FindStudentCourses).Methods(http.MethodGet)
	return r
}

// TestFindStudentCoursesSuccess tests the following scenario: the client requires the courses of a student
// and finds them. Therefore the response code should be 200 OK.
func TestFindStudentCoursesSuccess(t *testing.T) {

	_ = config.SetConfiguration("../config/config-test.json")

	// generate a token to be appended to the course creation request
	user := microservice.User{Name: "nome", Surname: "cognome", Username: "student_with_courses", Password: "password", Type: "student"}
	token, _ := microservice.GenerateAccessToken(user, []byte(config.Configuration.TokenPrivateKey))

	// Make the get request for course searching
	request, _ := http.NewRequest(http.MethodGet, "/didattica-mobile/api/v1.0/courses/students/student_with_courses", nil)
	request.AddCookie(&http.Cookie{Name: "token", Value: token})

	response := httptest.NewRecorder()
	handler := createTestGatewayGetStudentCourses()
	// a goroutine representing the microservice listens to the requests coming from the api gateway
	go mock.LaunchCourseManagementMock()
	// simulates a request-response interaction between client and api gateway
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Error("Expected 200 OK but got " + strconv.Itoa(response.Code) + " " + http.StatusText(response.Code))
	}

}

// TestFindStudentCoursesNotFound tests the following scenario: the client requires the courses of a student
// but the student is not subscribed to any course. Therefore the response code should be 404 NOT FOUND.
func TestFindStudentCoursesNotFound(t *testing.T) {

	_ = config.SetConfiguration("../config/config-test.json")

	// generate a token to be appended to the course creation request
	user := microservice.User{Name: "nome", Surname: "cognome", Username: "student_without_courses", Password: "password", Type: "student"}
	token, _ := microservice.GenerateAccessToken(user, []byte(config.Configuration.TokenPrivateKey))

	// Make the get request for course searching
	request, _ := http.NewRequest(http.MethodGet, "/didattica-mobile/api/v1.0/courses/students/student_without_courses", nil)
	request.AddCookie(&http.Cookie{Name: "token", Value: token})

	response := httptest.NewRecorder()
	handler := createTestGatewayGetStudentCourses()
	// a goroutine representing the microservice listens to the requests coming from the api gateway
	go mock.LaunchCourseManagementMock()
	// simulates a request-response interaction between client and api gateway
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Error("Expected 404 Not Found but got " + strconv.Itoa(response.Code) + " " + http.StatusText(response.Code))
	}

}
