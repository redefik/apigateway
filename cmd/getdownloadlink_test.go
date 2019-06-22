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

// NB: It is assumed the student with username "student_user" attends courses "course1" and "course2". The teacher "Mr Brown"
//     holds the courses "course3" and "course4". On the teaching material management micro-service are upload "file1"
//     for courses "course1", "course3" and "course4".

// createTestGatewayGetDownloadLink creates an http handler that handles the test requests
func createTestGatewayGetDownloadLink() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/didattica-mobile/api/v1.0/teachingMaterials/download/{username}/{courseId}/{fileName}",
		microservice.GetDownloadLinkToFile).Methods(http.MethodGet)
	return r
}

// TestGetDownloadLinkStudentSuccess tests the following scenario: the client is a student with username "student_user"
// and requires a file called "file1" from a course who him/her actually attend. The teacher of this course has
// previously uploaded the requested file.  Therefore the response code should be 200 OK and the body contains the link
// to download file.
func TestGetDownloadLinkStudentSuccess(t *testing.T) {

	_ = config.SetConfigurationFromFile("../config/config-test.json")

	// generate a token to be appended to request
	user := microservice.User{Name: "nome", Surname: "cognome", Username: "student_user", Password: "password", Type: "student", Mail: "name@example.com"}
	token, _ := microservice.GenerateAccessToken(user, []byte(config.Configuration.TokenPrivateKey))

	// Make the get request for teaching material searching
	request, _ := http.NewRequest(http.MethodGet,
		"/didattica-mobile/api/v1.0/teachingMaterials/download/student_user/course1/file1", nil)
	request.AddCookie(&http.Cookie{Name: "token", Value: token})

	response := httptest.NewRecorder()
	handler := createTestGatewayGetDownloadLink()
	// a goroutine representing the micro-service course management listens to the requests coming from the api gateway
	go mock.LaunchCourseManagementMock()
	// a goroutine representing the micro-service teaching material management listens to the requests coming from the api gateway
	go mock.LaunchTeachingMaterialManagementMock()
	// simulates a request-response interaction between client and api gateway
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Error("Expected 200 OK but got " + strconv.Itoa(response.Code) + " " + http.StatusText(response.Code))
	}

	jsonDecoder := json.NewDecoder(response.Body)
	var link string
	_ = jsonDecoder.Decode(&link)
	if link != "validDownloadLink" {
		t.Error("The teaching material management micro-service does not return the expected response")
	}
}

// TestGetDownloadLinkStudentUnauthorized tests the following scenario: the client is a student with username "student_user"
// and requires a file called "file1" from a course who him/her NOT actually attend. Therefore the response code should
// be 401 Unauthorized although the teacher of this course has previously uploaded the requested file
func TestGetDownloadLinkStudentUnauthorized(t *testing.T) {

	_ = config.SetConfigurationFromFile("../config/config-test.json")

	// generate a token to be appended to request
	user := microservice.User{Name: "nome", Surname: "cognome", Username: "student_user", Password: "password", Type: "student", Mail: "name@example.com"}
	token, _ := microservice.GenerateAccessToken(user, []byte(config.Configuration.TokenPrivateKey))

	// Make the get request for teaching material searching
	request, _ := http.NewRequest(http.MethodGet,
		"/didattica-mobile/api/v1.0/teachingMaterials/download/student_user/course3/file1", nil)
	request.AddCookie(&http.Cookie{Name: "token", Value: token})

	response := httptest.NewRecorder()
	handler := createTestGatewayGetDownloadLink()
	// a goroutine representing the micro-service course management listens to the requests coming from the api gateway
	go mock.LaunchCourseManagementMock()
	// a goroutine representing the micro-service teaching material management listens to the requests coming from the api gateway
	go mock.LaunchTeachingMaterialManagementMock()
	// simulates a request-response interaction between client and api gateway
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Error("Expected 401 OK but got " + strconv.Itoa(response.Code) + " " + http.StatusText(response.Code))
	}
}

// TestGetDownloadLinkStudentBadRequest tests the following scenario: the client is a student with username "student_user"
// and requires a file called "file2" from a course who him/her actually attend. However the teacher of this course has
// no uploaded the requested file. Therefore the response code should be 404 Not found
func TestGetDownloadLinkStudentBadRequest(t *testing.T) {

	_ = config.SetConfigurationFromFile("../config/config-test.json")

	// generate a token to be appended to request
	user := microservice.User{Name: "nome", Surname: "cognome", Username: "student_user", Password: "password", Type: "student", Mail: "name@example.com"}
	token, _ := microservice.GenerateAccessToken(user, []byte(config.Configuration.TokenPrivateKey))

	// Make the get request for teaching material searching
	request, _ := http.NewRequest(http.MethodGet,
		"/didattica-mobile/api/v1.0/teachingMaterials/download/student_user/course1/file2", nil)
	request.AddCookie(&http.Cookie{Name: "token", Value: token})

	response := httptest.NewRecorder()
	handler := createTestGatewayGetDownloadLink()
	// a goroutine representing the micro-service course management listens to the requests coming from the api gateway
	go mock.LaunchCourseManagementMock()
	// a goroutine representing the micro-service teaching material management listens to the requests coming from the api gateway
	go mock.LaunchTeachingMaterialManagementMock()
	// simulates a request-response interaction between client and api gateway
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Error("Expected 404 OK but got " + strconv.Itoa(response.Code) + " " + http.StatusText(response.Code))
	}
}

// TestGetDownloadLinkTeacherSuccess tests the following scenario: the client is a teacher and requires a file called
// "file1" from "course4" (a course who him/her actually hold). The teacher of this course has previously uploaded this
// file. Therefore he can download it: the response code should be 200 OK body contains the link to download file.
func TestGetDownloadLinkTeacherSuccess(t *testing.T) {

	_ = config.SetConfigurationFromFile("../config/config-test.json")

	// generate a token to be appended to request
	user := microservice.User{Name: "Mr", Surname: "Brown", Username: "username", Password: "password", Type: "teacher", Mail: "name@example.com"}
	token, _ := microservice.GenerateAccessToken(user, []byte(config.Configuration.TokenPrivateKey))

	// Make the get request for teaching material searching
	request, _ := http.NewRequest(http.MethodGet,
		"/didattica-mobile/api/v1.0/teachingMaterials/download/username/course4/file1", nil)
	request.AddCookie(&http.Cookie{Name: "token", Value: token})

	response := httptest.NewRecorder()
	handler := createTestGatewayGetDownloadLink()
	// a goroutine representing the micro-service course management listens to the requests coming from the api gateway
	go mock.LaunchCourseManagementMock()
	// a goroutine representing the micro-service teaching material management listens to the requests coming from the api gateway
	go mock.LaunchTeachingMaterialManagementMock()
	// simulates a request-response interaction between client and api gateway
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Error("Expected 200 OK but got " + strconv.Itoa(response.Code) + " " + http.StatusText(response.Code))
	}

	jsonDecoder := json.NewDecoder(response.Body)
	var link string
	_ = jsonDecoder.Decode(&link)
	if link != "validDownloadLink" {
		t.Error("The teaching material management micro-service does not return the expected response")
	}
}

// TestGetDownloadLinkTeacherUnauthorized tests the following scenario: the client is a teacher and requires a file
// called "file1" from a course ("course1") who him/her NOT actually hold. Therefore the response code should
// be 401 Unauthorized although the teacher of that course has previously uploaded the requested file.
func TestGetDownloadLinkTeacherUnauthorized(t *testing.T) {

	_ = config.SetConfigurationFromFile("../config/config-test.json")

	// generate a token to be appended to request
	user := microservice.User{Name: "Mr", Surname: "Brown", Username: "username", Password: "password", Type: "teacher", Mail: "name@example.com"}
	token, _ := microservice.GenerateAccessToken(user, []byte(config.Configuration.TokenPrivateKey))

	// Make the get request for teaching material searching
	request, _ := http.NewRequest(http.MethodGet,
		"/didattica-mobile/api/v1.0/teachingMaterials/download/student_user/course1/file1", nil)
	request.AddCookie(&http.Cookie{Name: "token", Value: token})

	response := httptest.NewRecorder()
	handler := createTestGatewayGetDownloadLink()
	// a goroutine representing the micro-service course management listens to the requests coming from the api gateway
	go mock.LaunchCourseManagementMock()
	// a goroutine representing the micro-service teaching material management listens to the requests coming from the api gateway
	go mock.LaunchTeachingMaterialManagementMock()
	// simulates a request-response interaction between client and api gateway
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Error("Expected 401 OK but got " + strconv.Itoa(response.Code) + " " + http.StatusText(response.Code))
	}
}

// TestGetDownloadLinkTeacherBadRequest tests the following scenario: the client is a teacher and requires a file called
// "file2" from a course ("course3") who him/her actually hold. However the same teacher has no previously uploaded the requested file.
// Therefore the response code should be 404 Not found
func TestGetDownloadLinkTeacherBadRequest(t *testing.T) {

	_ = config.SetConfigurationFromFile("../config/config-test.json")

	// generate a token to be appended to request
	user := microservice.User{Name: "Mr", Surname: "Brown", Username: "username", Password: "password", Type: "teacher", Mail: "name@example.com"}
	token, _ := microservice.GenerateAccessToken(user, []byte(config.Configuration.TokenPrivateKey))

	// Make the get request for teaching material searching
	request, _ := http.NewRequest(http.MethodGet,
		"/didattica-mobile/api/v1.0/teachingMaterials/download/username/course3/file2", nil)
	request.AddCookie(&http.Cookie{Name: "token", Value: token})

	response := httptest.NewRecorder()
	handler := createTestGatewayGetDownloadLink()
	// a goroutine representing the micro-service course management listens to the requests coming from the api gateway
	go mock.LaunchCourseManagementMock()
	// a goroutine representing the micro-service teaching material management listens to the requests coming from the api gateway
	go mock.LaunchTeachingMaterialManagementMock()
	// simulates a request-response interaction between client and api gateway
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Error("Expected 404 OK but got " + strconv.Itoa(response.Code) + " " + http.StatusText(response.Code))
	}
}
