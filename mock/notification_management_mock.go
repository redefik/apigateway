package mock

import (
	"encoding/json"
	"github.com/bitly/go-simplejson"
	"github.com/gorilla/mux"
	"github.com/redefik/sdccproject/apigateway/config"
	"github.com/redefik/sdccproject/apigateway/microservice"
	"io/ioutil"
	"log"
	"net/http"
)

// NotificationManagementMockCreateCourse simulates the behaviour of the notification management microservice
// when receives a request of course creation.
func NotificationManagementMockCreateCourse(w http.ResponseWriter, r *http.Request) {

	var course microservice.Course
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Panicln(err)
	}
	err = json.Unmarshal(body, &course)
	if err != nil {
		log.Panicln(err)
	}
	if course.Name == "courseSuccess" || course.Name == "courseFailInCourseManagement" {
		w.WriteHeader(http.StatusCreated)
		jsonBody := simplejson.New()
		jsonBody.Set("id", "idCourse")
		jsonBody.Set("name", "courseSuccess")
		jsonBody.Set("department", "department")
		jsonBody.Set("year", "2019-2020")
		body, _ := jsonBody.MarshalJSON()
		_, _ = w.Write(body)
	} else if course.Name == "courseFailInNotificationManagement" {
		w.WriteHeader(http.StatusInternalServerError)
		jsonBody := simplejson.New()
		jsonBody.Set("error", "internal server error")
		body, _ := jsonBody.MarshalJSON()
		_, _ = w.Write(body)
	}
}

// NotificationManagementMockDeleteCourse simulates the behaviour of the notification management micro-service
// when receives a request of course deletion.
func NotificationManagementMockDeleteCourse(w http.ResponseWriter, r *http.Request) {

	var course microservice.Course
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Panicln(err)
	}
	err = json.Unmarshal(body, &course)
	if err != nil {
		log.Panicln(err)
	}
	if course.Name == "courseFailInCourseManagement" {
		w.WriteHeader(http.StatusOK)
	}
}

// NotificationManagementMockAddStudentToCourse simulates the behaviour of the notification management micro-service
// when receives a request to add a student to a course
func NotificationManagementMockAddStudentToCourse(w http.ResponseWriter, r *http.Request) {

	var course microservice.Course
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Panicln(err)
	}
	err = json.Unmarshal(body, &course)
	if err != nil {
		log.Panicln(err)
	}

	if course.Name == "courseSuccess" || course.Name == "courseFailingInCourseManagement" ||
		course.Name == "courseToUnregisterFailureInCourseManagement" {
		w.WriteHeader(http.StatusOK)
	} else if course.Name == "courseFailingInNotificationManagement" {
		w.WriteHeader(http.StatusBadRequest)
	}
}

// NotificationManagementMockRemoveStudentToCourse simulates the behaviour of the notification management micro-service
// when receives a request to remove a student from a course
func NotificationManagementMockRemoveStudentToCourse(w http.ResponseWriter, r *http.Request) {

	var course microservice.Course
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Panicln(err)
	}
	err = json.Unmarshal(body, &course)
	if err != nil {
		log.Panicln(err)
	}

	if course.Name == "courseFailingInCourseManagement" || course.Name == "courseToUnregisterSuccess" ||
		course.Name == "courseToUnregisterFailureInCourseManagement" {
		w.WriteHeader(http.StatusOK)
	} else if course.Name == "courseToUnregisterFailureInNotificationManagement" {
		w.WriteHeader(http.StatusBadRequest)
	}

}

// starts a notification management micro-service mock
func LaunchNotificationManagementMock() {
	r := mux.NewRouter()
	r.HandleFunc("/notification_management/api/v1.0/course", NotificationManagementMockCreateCourse).Methods(http.MethodPost)
	r.HandleFunc("/notification_management/api/v1.0/course", NotificationManagementMockDeleteCourse).Methods(http.MethodDelete)
	r.HandleFunc("/notification_management/api/v1.0/course/student/{mail}", NotificationManagementMockAddStudentToCourse).Methods(http.MethodPut)
	r.HandleFunc("/notification_management/api/v1.0/course/student/{mail}", NotificationManagementMockRemoveStudentToCourse).Methods(http.MethodDelete)
	_ = http.ListenAndServe(config.Configuration.ApiGatewayAddress+"81", r)
}
