package courseSubscription

//WARNING: THESE TESTS ARE BEEN DESIGNED FOR THE OLD VERSION OF ADD COURSE TO STUDENT. THEIR RESULTS ARE NO MORE VALID.
//SEE newaddcoursetostudent_test.go
//
//import (
//	"github.com/gorilla/mux"
//	"github.com/redefik/sdccproject/apigateway/config"
//	"github.com/redefik/sdccproject/apigateway/microservice"
//	"github.com/redefik/sdccproject/apigateway/mock"
//	"net/http"
//	"net/http/httptest"
//	"strconv"
//	"testing"
//)
//
//// createTestGatewayAddCourseToStudent creates an http handler that handles the test requests
//func createTestGatewayAddCourseToStudent() http.Handler {
//	r := mux.NewRouter()
//	r.HandleFunc("/didattica-mobile/api/v1.0/students/{username}/courses/{id}", microservice.AddCourseToStudent).Methods(http.MethodPut)
//	return r
//}
//
///*TestAddCourseToStudentSuccess tests the following scenario: the client requests to add a course to an existing student.
//Therefore, it is expected that the response to the client is 200 OK*/
//func TestAddCourseToStudentSuccess(t *testing.T) {
//
//	config.SetConfiguration("../config/config-test.json")
//
//	// generate a token to be appended to the request
//	user := microservice.User{Name: "nome", Surname: "cognome", Username: "existent_student", Password: "password", Type: "student", Mail: "name@example.com"}
//	token, _ := microservice.GenerateAccessToken(user, []byte(config.Configuration.TokenPrivateKey))
//
//	// make the PUT request for the course append
//	request, _ := http.NewRequest(http.MethodPut, "/didattica-mobile/api/v1.0/students/existent_student/courses/course_id", nil)
//	request.AddCookie(&http.Cookie{Name: "token", Value: token})
//
//	response := httptest.NewRecorder()
//	handler := createTestGatewayAddCourseToStudent()
//	// a goroutine representing the microservice listens to the requests coming from the api gateway
//	go mock.LaunchCourseManagementMock()
//	// simulates a request-response interaction between client and api gateway
//	handler.ServeHTTP(response, request)
//
//	if response.Code != http.StatusOK {
//		t.Error("Expected 200 Ok but got " + strconv.Itoa(response.Code) + " " + http.StatusText(response.Code))
//	}
//
//}
//
///*TestAddCourseToStudentSuccessNotExistentStudent tests the following scenario: the client requests to add a course to a not existing student.
//It is expected that the Api gateway provides student creation and the final response is 200 OK*/
//func TestAddCourseToStudentSuccessNotExistentStudent(t *testing.T) {
//
//	config.SetConfiguration("../config/config-test.json")
//
//	// generate a token to be appended to the request
//	user := microservice.User{Name: "nome", Surname: "cognome", Username: "not_existent_student1", Password: "password", Type: "student", Mail: "name@example.com"}
//	token, _ := microservice.GenerateAccessToken(user, []byte(config.Configuration.TokenPrivateKey))
//
//	// make the PUT request for the course append
//	request, _ := http.NewRequest(http.MethodPut, "/didattica-mobile/api/v1.0/students/not_existent_student1/courses/course_id", nil)
//	request.AddCookie(&http.Cookie{Name: "token", Value: token})
//
//	response := httptest.NewRecorder()
//	handler := createTestGatewayAddCourseToStudent()
//	// a goroutine representing the microservice listens to the requests coming from the api gateway
//	go mock.LaunchCourseManagementMock()
//	// simulates a request-response interaction between client and api gateway
//	handler.ServeHTTP(response, request)
//
//	if response.Code != http.StatusOK {
//		t.Error("Expected 200 Ok but got " + strconv.Itoa(response.Code) + " " + http.StatusText(response.Code))
//	}
//
//}
//
///*TestAddCourseToStudentFailNotExistentCourse tests the following scenario: the client requests to add a not-existing course to a student
//Therefore, the response should be 404 Not Found*/
//func TestAddCourseToStudentFailNotExistentCourse(t *testing.T) {
//	config.SetConfiguration("../config/config-test.json")
//
//	// generate a token to be appended to the request
//	user := microservice.User{Name: "nome", Surname: "cognome", Username: "existent_student", Password: "password", Type: "student", Mail: "name@example.com"}
//	token, _ := microservice.GenerateAccessToken(user, []byte(config.Configuration.TokenPrivateKey))
//
//	// make the PUT request for the course append
//	request, _ := http.NewRequest(http.MethodPut, "/didattica-mobile/api/v1.0/students/existent_student/courses/not_existent_course", nil)
//	request.AddCookie(&http.Cookie{Name: "token", Value: token})
//
//	response := httptest.NewRecorder()
//	handler := createTestGatewayAddCourseToStudent()
//	// a goroutine representing the microservice listens to the requests coming from the api gateway
//	go mock.LaunchCourseManagementMock()
//	// simulates a request-response interaction between client and api gateway
//	handler.ServeHTTP(response, request)
//
//	if response.Code != http.StatusNotFound {
//		t.Error("Expected 404 Not Found but got " + strconv.Itoa(response.Code) + " " + http.StatusText(response.Code))
//	}
//}
//
///*TestAddCourseToStudentFailNotExistentStudent tests the following scenario: the client requests to add a course to a not-existing student.
//The Api Gateway attempts to create the student but something unexpected goes wrong. So the response should be 500*/
//func TestAddCourseToStudentFailNotExistentStudent(t *testing.T) {
//	config.SetConfiguration("../config/config-test.json")
//
//	// generate a token to be appended to the request
//	user := microservice.User{Name: "nome", Surname: "cognome", Username: "not_existent_student2", Password: "password", Type: "student", Mail: "name@example.com"}
//	token, _ := microservice.GenerateAccessToken(user, []byte(config.Configuration.TokenPrivateKey))
//
//	// make the PUT request for the course append
//	request, _ := http.NewRequest(http.MethodPut, "/didattica-mobile/api/v1.0/students/not_existent_student2/courses/not_existent_course", nil)
//	request.AddCookie(&http.Cookie{Name: "token", Value: token})
//
//	response := httptest.NewRecorder()
//	handler := createTestGatewayAddCourseToStudent()
//	// a goroutine representing the microservice listens to the requests coming from the api gateway
//	go mock.LaunchCourseManagementMock()
//	// simulates a request-response interaction between client and api gateway
//	handler.ServeHTTP(response, request)
//
//	if response.Code != http.StatusInternalServerError {
//		t.Error("Expected 500 Internal Server Error but got " + strconv.Itoa(response.Code) + " " + http.StatusText(response.Code))
//	}
//}
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
