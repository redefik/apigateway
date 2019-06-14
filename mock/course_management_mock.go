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

var studentCreated = false

// CourseManagementMockCreateCourse simulates the behaviour of the course management microservice when receives a request of
// course creation.
func CourseManagementMockCreateCourse(w http.ResponseWriter, r *http.Request) {

	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	//WARNING: OLD MOCK FOR COURSE CREATION. VALID UNTIL INSERTING OF DISTRIBUTED TRANSACTION
	//
	//w.Header().Set("Content-Type", "application/json")
	//w.WriteHeader(http.StatusOK)
	//response := simplejson.New()
	//response.Set("mock", "response")
	//responsePayload, err := response.MarshalJSON()
	//if err != nil {
	//	log.Panicln(err)
	//}
	//_, err = w.Write(responsePayload)
	//if err != nil {
	//	log.Panic(err)
	//}
	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	var course microservice.Course
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Panicln(err)
	}
	err = json.Unmarshal(body, &course)
	if err != nil {
		log.Panicln(err)
	}
	if course.Name == "courseSuccess" || course.Name == "courseFailInNotificationManagement" {
		w.WriteHeader(http.StatusCreated)
		jsonBody := simplejson.New()
		jsonBody.Set("id", "idCourse")
		jsonBody.Set("name", "courseSuccess")
		jsonBody.Set("department", "department")
		jsonBody.Set("year", "2019-2020")
		body, _ := jsonBody.MarshalJSON()
		_, _ = w.Write(body)
	} else if course.Name == "courseFailInCourseManagement" {
		w.WriteHeader(http.StatusInternalServerError)
		jsonBody := simplejson.New()
		jsonBody.Set("error", "internal server error")
		body, _ := jsonBody.MarshalJSON()
		_, _ = w.Write(body)
	}
}

func CourseManagementMockDeleteCourse(w http.ResponseWriter, r *http.Request) {

	courseId := mux.Vars(r)["courseId"]
	if courseId == "courseId" {
		w.WriteHeader(http.StatusOK)
	}
}

// CourseManagementMockSearchCourse simulates the behaviour of the course management microservice when receives a request of
// course research. The response is positive only if the research is for name and the string sequence is "seq".
func CourseManagementMockSearchCourse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
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

	if requestBody.Username == "notExistingStudent" {
		studentCreated = true
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// CourseManagementMockAddCourseToStudent simulates the behaviour of the course management microservice when receives
// a request of course appending to a student.
func CourseManagementMockAddCourseToStudent(w http.ResponseWriter, r *http.Request) {

	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	//WARNING: OLD MOCK FOR COURSE SUBSCRIPTION. VALID UNTIL INSERTING OF DISTRIBUTED TRANSACTION.
	//
	//if mux.Vars(r)["username"] == "existent_student" && mux.Vars(r)["id"] != "not_existent_course" {
	//	w.WriteHeader(http.StatusOK)
	//} else if mux.Vars(r)["username"] == "existent_student" && mux.Vars(r)["id"] == "not_existent_course" {
	//	w.WriteHeader(http.StatusNotFound)
	//	errorResponse := simplejson.New()
	//	errorResponse.Set("error", "Course Not Found")
	//	errorResponsePayload, err := errorResponse.MarshalJSON()
	//	if err != nil {
	//		log.Panicln(err)
	//	}
	//	_, err = w.Write(errorResponsePayload)
	//	if err != nil {
	//		log.Panicln(err)
	//	}
	//} else if mux.Vars(r)["username"] == "not_existent_student1" || mux.Vars(r)["username"] == "not_existent_student2" {
	//	if !studentCreated {
	//		w.WriteHeader(http.StatusNotFound)
	//		errorResponse := simplejson.New()
	//		errorResponse.Set("error", "Student Not Found")
	//		errorResponsePayload, err := errorResponse.MarshalJSON()
	//		if err != nil {
	//			log.Panicln(err)
	//		}
	//		_, err = w.Write(errorResponsePayload)
	//		if err != nil {
	//			log.Panicln(err)
	//		}
	//	} else {
	//		w.WriteHeader(http.StatusOK)
	//	}
	//}
	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	if (mux.Vars(r)["id"] == "idCourseSuccess" && mux.Vars(r)["username"] == "existingUser") ||
		mux.Vars(r)["id"] == "idCourseFailingInNotificationManagement" ||
		mux.Vars(r)["id"] == "idCourseToUnregisterFailureInNotificationManagement" {
		w.WriteHeader(http.StatusOK)
	} else if mux.Vars(r)["username"] == "notExistingStudent" {
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
	} else if mux.Vars(r)["id"] == "idCourseFailingInCourseManagement" {
		w.WriteHeader(http.StatusBadRequest)
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
	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	//WARNING: OLD MOCK FOR COURSE UNSUBSCRIPTION. VALID UNTIL INSERTING OF DISTRIBUTED TRANSACTION
	//studentUsername := mux.Vars(r)["username"]
	//courseId := mux.Vars(r)["id"]
	//if studentUsername == "student_test" && courseId == "course_test" {
	//	w.WriteHeader(http.StatusOK)
	//}
	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	if mux.Vars(r)["id"] == "idCourseFailingInNotificationManagement" ||
		mux.Vars(r)["id"] == "idCourseToUnregisterSuccess" ||
		mux.Vars(r)["id"] == "idCourseToUnregisterFailureInNotificationManagement" {
		w.WriteHeader(http.StatusOK)
	} else if mux.Vars(r)["id"] == "idCourseToUnregisterFailureInCourseManagement" {
		w.WriteHeader(http.StatusBadRequest)
	}
}

// starts a course management microservice mock
func LaunchCourseManagementMock() {
	r := mux.NewRouter()
	r.HandleFunc("/course_management/api/v1.0/courses", CourseManagementMockCreateCourse).Methods(http.MethodPost)
	r.HandleFunc("/course_management/api/v1.0/courses/{courseId}", CourseManagementMockDeleteCourse).Methods(http.MethodDelete)
	r.HandleFunc("/course_management/api/v1.0/courses/students/{username}", CourseManagementMockFindStudentCourses).Methods(http.MethodGet)
	r.HandleFunc("/course_management/api/v1.0/courses/{by}/{string}", CourseManagementMockSearchCourse).Methods(http.MethodGet)
	r.HandleFunc("/course_management/api/v1.0/students", CourseManagementMockCreateStudent).Methods(http.MethodPost)
	r.HandleFunc("/course_management/api/v1.0/students/{username}/courses/{id}", CourseManagementMockAddCourseToStudent).Methods(http.MethodPut)
	r.HandleFunc("/course_management/api/v1.0/exams", CourseManagementMockCreateExam).Methods(http.MethodPost)
	r.HandleFunc("/course_management/api/v1.0/exams/{course}", CourseManagementMockSearchExam).Methods(http.MethodGet)
	r.HandleFunc("/course_management/api/v1.0/exams/{examId}/students/{studentUsername}", CourseManagementMockReserveExam).Methods(http.MethodPut)
	r.HandleFunc("/course_management/api/v1.0/students/{username}/courses/{id}", CourseManagementMockUnsubscribeFromCourse).Methods(http.MethodDelete)
	http.ListenAndServe(config.Configuration.ApiGatewayAddress, r)
}
