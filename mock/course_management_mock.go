package mock

import (
	"encoding/json"
	"github.com/bitly/go-simplejson"
	"github.com/gorilla/mux"
	"github.com/redefik/sdccproject/apigateway/config"
	"github.com/redefik/sdccproject/apigateway/microservice"
	"log"
	"net/http"
)

var studentCreated = false

// CourseManagementMockCreateCourse simulates the behaviour of the course management microservice when receives a request of
// course creation.
func CourseManagementMockCreateCourse(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := simplejson.New()
	response.Set("mock", "response")
	responsePayload, err := response.MarshalJSON()
	if err != nil {
		log.Panicln(err)
	}
	_, err = w.Write(responsePayload)
	if err != nil {
		log.Panic(err)
	}
}

// CourseManagementMockSearchCourse simulates the behaviour of the course management microservice when receives a request of
// course research. The response is positive only if the research is for name and the string sequence is "seq".
func CourseManagementMockSearchCourse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//log.Panic("by = " + mux.Vars(r)["by"])
	//log.Panic("string = "+ mux.Vars(r)["string"])
	if mux.Vars(r)["by"] == "name" && mux.Vars(r)["string"] == "seq" {
		w.WriteHeader(http.StatusOK)
	} else if mux.Vars(r)["by"] == "notvalid" {
		w.WriteHeader(http.StatusBadRequest)
	} else if mux.Vars(r)["by"] == "teacher" && mux.Vars(r)["string"] == "Mr-Brown" {
		//For test TestGetDownloadLinkTeacherSuccess in getDownloadLink_test
		w.WriteHeader(http.StatusOK)
		var courses []microservice.CourseMinimized
		courses = append(courses, microservice.CourseMinimized{"course3"})
		courses = append(courses, microservice.CourseMinimized{"course4"})
		w.Header().Set("Content-Type", "application/json")
		response, err := json.Marshal(&courses)
		if err != nil {
			log.Panicln(err)
		}
		_, err = w.Write(response)
		if err != nil {
			log.Panicln(err)
		}
		return
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
	response := simplejson.New()
	response.Set("mock", "response")
	responsePayload, err := response.MarshalJSON()
	if err != nil {
		log.Panicln(err)
	}
	_, err = w.Write(responsePayload)
	if err != nil {
		log.Panic(err)
	}
}

// CourseManagementMockCreateStudent simulates the behaviour of the course management microservice when receives a request
// of student creation. This request is sent from the Api Gateway the first time that a course has to be added to a student
func CourseManagementMockCreateStudent(w http.ResponseWriter, r *http.Request) {

	type StudentCreationRequest struct {
		Name     string `json:"name"`
		Username string `json:"username"`
	}

	var requestBody StudentCreationRequest
	jsonDecoder := json.NewDecoder(r.Body)
	err := jsonDecoder.Decode(&requestBody)
	if err != nil {
		log.Panicln(err)
	}

	if requestBody.Username == "not_existent_student1" {
		studentCreated = true
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// CourseManagementMockAddCourseToStudent simulates the behaviour of the course management microservice when receives
// a request of course appending to a student.
func CourseManagementMockAddCourseToStudent(w http.ResponseWriter, r *http.Request) {
	if mux.Vars(r)["username"] == "existent_student" && mux.Vars(r)["id"] != "not_existent_course" {
		w.WriteHeader(http.StatusOK)
	} else if mux.Vars(r)["username"] == "existent_student" && mux.Vars(r)["id"] == "not_existent_course" {
		w.WriteHeader(http.StatusNotFound)
		errorResponse := simplejson.New()
		errorResponse.Set("error", "Course Not Found")
		errorResponsePayload, err := errorResponse.MarshalJSON()
		if err != nil {
			log.Panicln(err)
		}
		_, err = w.Write(errorResponsePayload)
		if err != nil {
			log.Panicln(err)
		}
	} else if mux.Vars(r)["username"] == "not_existent_student1" || mux.Vars(r)["username"] == "not_existent_student2" {
		if !studentCreated {
			w.WriteHeader(http.StatusNotFound)
			errorResponse := simplejson.New()
			errorResponse.Set("error", "Student Not Found")
			errorResponsePayload, err := errorResponse.MarshalJSON()
			if err != nil {
				log.Panicln(err)
			}
			_, err = w.Write(errorResponsePayload)
			if err != nil {
				log.Panicln(err)
			}
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}

}

// CourseManagementMockFindStudentCourses simulates the behaviour of the course management microservice when receives
// a request of course searching for a given student.
func CourseManagementMockFindStudentCourses(w http.ResponseWriter, r *http.Request) {
	if mux.Vars(r)["username"] == "student_with_courses" {
		w.WriteHeader(http.StatusOK)
	} else if mux.Vars(r)["username"] == "student_user" {
		//For test TestGetDownloadLinkStudentSuccess in getDownloadLink_test
		w.WriteHeader(http.StatusOK)
		var courses []microservice.CourseMinimized
		courses = append(courses, microservice.CourseMinimized{"course1"})
		courses = append(courses, microservice.CourseMinimized{"course2"})
		w.Header().Set("Content-Type", "application/json")
		response, err := json.Marshal(&courses)
		if err != nil {
			log.Panicln(err)
		}
		_, err = w.Write(response)
		if err != nil {
			log.Panicln(err)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func CourseManagementMockCreateExam(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := simplejson.New()
	response.Set("mock", "response")
	responsePayload, err := response.MarshalJSON()
	if err != nil {
		log.Panicln(err)
	}
	_, err = w.Write(responsePayload)
	if err != nil {
		log.Panic(err)
	}
}

// CourseManagementMockSearchExam simulates the behaviour of the course management microservice when receives a request of
// exam research. The response is positive (some exams are found) only if the id of course is idSuccess.
func CourseManagementMockSearchExam(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if mux.Vars(r)["course"] == "idSuccess" {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
	response := simplejson.New()
	response.Set("mock", "response")
	responsePayload, err := response.MarshalJSON()
	if err != nil {
		log.Panicln(err)
	}
	_, err = w.Write(responsePayload)
	if err != nil {
		log.Panic(err)
	}
}

// CourseManagementMockReserveExam simulates the behaviour of the course management microservice when receives a request
// of exam reservation creation.
func CourseManagementMockReserveExam(w http.ResponseWriter, r *http.Request) {
	examId := mux.Vars(r)["examId"]
	studentUsername := mux.Vars(r)["studentUsername"]
	if examId == "existent_exam" && studentUsername == "existent_student" {
		w.WriteHeader(http.StatusOK)
		return
	}
	if studentUsername == "not_existent_student" || examId == "not_existent_exam" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

// CourseManagementMockUnsubscribeFromCourse simulates the behaviour of the course management microservice when it is
// asked to remove a student from a course
func CourseManagementMockUnsubscribeFromCourse(w http.ResponseWriter, r *http.Request) {
	studentUsername := mux.Vars(r)["username"]
	courseId := mux.Vars(r)["id"]
	if studentUsername == "student_test" && courseId == "course_test" {
		w.WriteHeader(http.StatusOK)
	}
}

// CourseManagementMockPushNotification simulates the behaviour of the course management microservice when it receives
// a notification push request
func CourseManagementMockPushNotification(w http.ResponseWriter, r *http.Request) {
	courseId := mux.Vars(r)["courseId"]
	if courseId == "courseId" {
		w.WriteHeader(http.StatusOK)
	}
}

// starts a course management microservice mock
func LaunchCourseManagementMock() {
	r := mux.NewRouter()
	r.HandleFunc("/course_management/api/v1.0/courses", CourseManagementMockCreateCourse).Methods(http.MethodPost)
	r.HandleFunc("/course_management/api/v1.0/courses/students/{username}", CourseManagementMockFindStudentCourses).Methods(http.MethodGet)
	r.HandleFunc("/course_management/api/v1.0/courses/{by}/{string}", CourseManagementMockSearchCourse).Methods(http.MethodGet)
	r.HandleFunc("/course_management/api/v1.0/students", CourseManagementMockCreateStudent).Methods(http.MethodPost)
	r.HandleFunc("/course_management/api/v1.0/students/{username}/courses/{id}", CourseManagementMockAddCourseToStudent).Methods(http.MethodPut)
	r.HandleFunc("/course_management/api/v1.0/exams", CourseManagementMockCreateExam).Methods(http.MethodPost)
	r.HandleFunc("/course_management/api/v1.0/exams/{course}", CourseManagementMockSearchExam).Methods(http.MethodGet)
	r.HandleFunc("/course_management/api/v1.0/exams/{examId}/students/{studentUsername}", CourseManagementMockReserveExam).Methods(http.MethodPut)
	r.HandleFunc("/course_management/api/v1.0/students/{username}/courses/{id}", CourseManagementMockUnsubscribeFromCourse).Methods(http.MethodDelete)
	r.HandleFunc("/course_management/api/v1.0/courses/{courseId}/notification", CourseManagementMockPushNotification).Methods(http.MethodPost)
	http.ListenAndServe(config.Configuration.ApiGatewayAddress, r)
}
