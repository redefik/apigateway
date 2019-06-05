package microservice

import (
	"github.com/gorilla/mux"
	"github.com/redefik/sdccproject/apigateway/config"
	"log"
	"net/http"
)

// CreateExam process the exam creation request coming from the client validating the embedded access token. If the token
// is properly signed, it checks if the request comes from a student, because exam creation is allowed for teachers
// only. Upon successful validation, the request is forwarded to the microservice and the response is forwarded to the
// client.
func CreateExam(w http.ResponseWriter, r *http.Request) {

	/* For authentication purpose the access token is read from the Cookie header */
	tokenString, err := GetToken(w, r)
	if err != nil {
		return
	}
	/* The token is decoded and the claims are obtained for further checks */
	decodedToken, err := ValidateToken(tokenString, w)
	if err != nil {
		return
	}
	/* Upon successful authentication check if the request comes from a teacher.
	In case of student request an Unauthorized code is returned */
	if decodedToken.Type != "teacher" {
		MakeErrorResponse(w, http.StatusUnauthorized, "Permission denied")
		log.Println("Permission denied")
		return
	}
	/* Upon successful validation, the request is forwarded to the course management microservice and the response is
	returned to the client */
	err = ForwardAndReturnPost(config.Configuration.CourseManagementAddress+"exams", "application/json", w, r)
	if err != nil {
		MakeErrorResponse(w, http.StatusInternalServerError, "Api Gateway - Internal Server Error")
		log.Println("Api Gateway - Internal Server Error")
		return
	}
}

// FindExamByCourse process the exam searching request coming from the client and validate the embedded access token. It verify
// if the token is properly signed and not expired. Upon successful validation, the request is forwarded to the microservice
// and the response is forwarded to the client.
func FindExamByCourse(w http.ResponseWriter, r *http.Request) {

	/* For authentication purpose the access token is read from the Cookie header */
	tokenString, err := GetToken(w, r)
	if err != nil {
		return
	}
	/* The token is decoded and the claims are obtained for further checks */
	_, err = ValidateToken(tokenString, w)
	if err != nil {
		return
	}
	/* Upon successful validation, the request is forwarded to the course management microservice and the response is
	returned to the client*/
	vars := mux.Vars(r)
	course := vars["course"]
	err = ForwardAndReturnGet(config.Configuration.CourseManagementAddress+"exams"+"/"+course, w)
	if err != nil {
		MakeErrorResponse(w, http.StatusInternalServerError, "Api Gateway - Internal Server Error")
		log.Println("Api Gateway - Internal Server Error")
		return
	}
}

// ReserveExam process the exam reservation request provided by the client and validate the embedded access token. On
// successful validatio, it forwards the request to the course management microservice and returns the response to the client
func ReserveExam(w http.ResponseWriter, r *http.Request) {
	/* For authentication purpose the access token is read from the Cookie header */
	tokenString, err := GetToken(w, r)
	if err != nil {
		return
	}
	/* The token is decoded and the claims are obtained for further checks */
	_, err = ValidateToken(tokenString, w)
	if err != nil {
		return
	}
	/* Upon successful validation, the request is forwarded to the course management microservice and the response is
	returned to the client*/
	vars := mux.Vars(r)
	examId := vars["examId"]
	studentUsername := vars["studentUsername"]
	err = ForwardAndReturnPut(config.Configuration.CourseManagementAddress+"exams"+"/"+examId+"/students/"+studentUsername, w, r)
	if err != nil {
		MakeErrorResponse(w, http.StatusInternalServerError, "Api Gateway - Internal Server Error")
		log.Println("Api Gateway - Internal Server Error")
		return
	}
}
