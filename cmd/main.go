package main

import (
	"github.com/gorilla/mux"
	"github.com/redefik/sdccproject/apigateway/config"
	"github.com/redefik/sdccproject/apigateway/microservice"
	"log"
	"net/http"
)

// healthCheck handles the requests coming from an external component responsible for verifying the status of the api
// gateway
func healthCheck(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func main() {

	// Read the listening address of the gateway and the address of the other microservices
	err := config.SetConfigurationFromEnvironment()
	if err != nil {
		log.Panicln(err)
	}
	r := mux.NewRouter()
	// Register the handlers for the various HTTP requests
	r.HandleFunc("/didattica-mobile/api/v1.0/users", microservice.RegisterUser).Methods(http.MethodPost)
	r.HandleFunc("/didattica-mobile/api/v1.0/token", microservice.LoginUser).Methods(http.MethodPost)
	r.HandleFunc("/didattica-mobile/api/v1.0/courses", microservice.CreateCourse).Methods(http.MethodPost)
	r.HandleFunc("/didattica-mobile/api/v1.0/courses/students/{username}", microservice.FindStudentCourses).Methods(http.MethodGet)
	r.HandleFunc("/didattica-mobile/api/v1.0/courses/{by}/{string}", microservice.FindCourse).Methods(http.MethodGet)
	r.HandleFunc("/didattica-mobile/api/v1.0/students/{username}", microservice.UnsubscribeStudentFromCourse).Methods(http.MethodDelete)
	r.HandleFunc("/didattica-mobile/api/v1.0/students/{username}", microservice.AddCourseToStudent).Methods(http.MethodPut)
	r.HandleFunc("/didattica-mobile/api/v1.0/exams", microservice.CreateExam).Methods(http.MethodPost)
	r.HandleFunc("/didattica-mobile/api/v1.0/exams/{examId}/students/{studentUsername}", microservice.ReserveExam).Methods(http.MethodPut)
	r.HandleFunc("/didattica-mobile/api/v1.0/exams/{course}", microservice.FindExamByCourse).Methods(http.MethodGet)
	r.HandleFunc("/didattica-mobile/api/v1.0/teachingMaterials/{courseId}", microservice.FindTeachingMaterialByCourse).Methods(http.MethodGet)
	r.HandleFunc("/didattica-mobile/api/v1.0/teachingMaterials/download/{username}/{courseId}/{fileName}", microservice.GetDownloadLinkToFile).Methods(http.MethodGet)
	r.HandleFunc("/didattica-mobile/api/v1.0/notification/course/{courseId}", microservice.PushCourseNotification).Methods(http.MethodPost)
	r.HandleFunc("/", healthCheck).Methods(http.MethodGet)
	// Wait for incoming requests. A new goroutine is created to serve each request
	log.Fatal(http.ListenAndServe(config.Configuration.ApiGatewayAddress, r))
}
