package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/redefik/sdccproject/apigateway/config"
	"github.com/redefik/sdccproject/apigateway/microservice"
	"github.com/redefik/sdccproject/apigateway/mock"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

// createTestGatewayFindTeachingMaterialByCourse creates an http handler that handles the test requests
func createTestGatewayFindTeachingMaterialByCourse() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/didattica-mobile/api/v1.0/teachingMaterials/{courseId}",
		microservice.FindTeachingMaterialByCourse).Methods(http.MethodGet)
	return r
}

// TestFindStudentCoursesSuccess tests the following scenario: the client requires the teaching material for an existing
// course. Therefore the response code should be 200 OK and the body contains the names of found files.
// It is assumed existing two file (named file1 and file2) for course with id "courseIdWithTeachingMaterial".
func TestFindTeachingMaterialByCourseSuccess(t *testing.T) {

	_ = config.SetConfiguration("../config/config-test.json")

	// generate a token to be appended to request
	user := microservice.User{Name: "nome", Surname: "cognome", Username: "username", Password: "password", Type: "student", Mail: "name@example.com"}
	token, _ := microservice.GenerateAccessToken(user, []byte(config.Configuration.TokenPrivateKey))

	// Make the get request for teaching material searching
	request, _ := http.NewRequest(http.MethodGet, "/didattica-mobile/api/v1.0/teachingMaterials/courseIdWithTeachingMaterial", nil)
	request.AddCookie(&http.Cookie{Name: "token", Value: token})

	response := httptest.NewRecorder()
	handler := createTestGatewayFindTeachingMaterialByCourse()
	// a goroutine representing the micro-service listens to the requests coming from the api gateway
	go mock.LaunchTeachingMaterialManagementMock()
	// simulates a request-response interaction between client and api gateway
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Error("Expected 200 OK but got " + strconv.Itoa(response.Code) + " " + http.StatusText(response.Code))
	}

	jsonDecoder := json.NewDecoder(response.Body)
	var files []string
	_ = jsonDecoder.Decode(&files)
	if files[0] != "file1" && files[1] != "file2" {
		t.Error("The teaching material management micro-service does not return the expected files")
	}
}

// TestFindStudentCoursesNotSuccess tests the following scenario: the client requires the teaching material for a not
// existing course (or a course with no associated teaching materials). Therefore the response code should be 200 OK but
// the body is an empty json array.
func TestFindTeachingMaterialByCourseNotSuccess(t *testing.T) {

	_ = config.SetConfiguration("../config/config-test.json")

	// generate a token to be appended to request
	user := microservice.User{Name: "nome", Surname: "cognome", Username: "username", Password: "password", Type: "student", Mail: "name@example.com"}
	token, _ := microservice.GenerateAccessToken(user, []byte(config.Configuration.TokenPrivateKey))

	// Make the get request for teaching material searching
	request, _ := http.NewRequest(http.MethodGet,
		"/didattica-mobile/api/v1.0/teachingMaterials/courseIdWithoutTeachingMaterial", nil)
	request.AddCookie(&http.Cookie{Name: "token", Value: token})

	response := httptest.NewRecorder()
	handler := createTestGatewayFindTeachingMaterialByCourse()
	// a goroutine representing the micro-service listens to the requests coming from the api gateway
	go mock.LaunchTeachingMaterialManagementMock()
	// simulates a request-response interaction between client and api gateway
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Error("Expected 200 OK but got " + strconv.Itoa(response.Code) + " " + http.StatusText(response.Code))
	}

	jsonDecoder := json.NewDecoder(response.Body)
	var files []string
	_ = jsonDecoder.Decode(&files)
	if len(files) != 0 {
		t.Error("The teaching material management micro-service does not return the expected files")
	}
}
