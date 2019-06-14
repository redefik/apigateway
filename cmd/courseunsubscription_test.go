package main

//WARNING: THESE TESTS ARE BEEN DESIGNED FOR THE OLD VERSION OF ADD COURSE TO STUDENT. THEIR RESULTS ARE NO MORE VALID.
//SEE newcourseunsubscription_test.go

//// createTestGatewayCourseUnsubscription creates an http handler that handles the test requests
//func createTestGatewayCourseUnsubscription() http.Handler {
//	r := mux.NewRouter()
//	r.HandleFunc("/didattica-mobile/api/v1.0/students/{username}/courses/{id}", microservice.UnsubscribeStudentFromCourse).Methods(http.MethodDelete)
//	return r
//}
//
///*TestUnsubscribeStudentSuccess tests the following scenario: the client successfully removes a student from a course.
//Therefore, it is expected that the response to the client is 200 OK*/
//func TestCourseUnsubscribeStudentSuccess(t *testing.T) {
//
//	config.SetConfiguration("../config/config-test.json")
//
//	// generate a token to be appended to the request
//	user := microservice.User{Name: "nome", Surname: "cognome", Username: "student_test", Password: "password", Type: "student", Mail: "name@example.com"}
//	token, _ := microservice.GenerateAccessToken(user, []byte(config.Configuration.TokenPrivateKey))
//
//	// make the PUT request for the course append
//	request, _ := http.NewRequest(http.MethodDelete, "/didattica-mobile/api/v1.0/students/student_test/courses/course_test", nil)
//	request.AddCookie(&http.Cookie{Name: "token", Value: token})
//
//	response := httptest.NewRecorder()
//	handler := createTestGatewayCourseUnsubscription()
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
