package microservice

import (
	"bytes"
	"encoding/json"
	"github.com/bitly/go-simplejson"
	"github.com/gorilla/mux"
	"github.com/redefik/sdccproject/apigateway/config"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

// CreateCourse process the course creation request coming from the client validating the embedded access token. If the token
// is properly signed, it checks if the request comes from a student, because course creation is allowed for teachers
// only. Upon successful validation, the request is forwarded to the microservice and the response is forwarded to the
// client.
func CreateCourse(w http.ResponseWriter, r *http.Request) {

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
	err = ForwardAndReturnPost(config.Configuration.CourseManagementAddress+"courses", "application/json", w, r)
	if err != nil {
		MakeErrorResponse(w, http.StatusInternalServerError, "Api Gateway - Internal Server Error")
		log.Println("Api Gateway - Internal Server Error")
		return
	}

}

// FindCourse process the course searching request coming from the client and validate the embedded access token. It verify
// if the token is properly signed and not expired. Upon successful validation, the request is forwarded to the microservice
// and the response is forwarded to the client.
func FindCourse(w http.ResponseWriter, r *http.Request) {

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
	vars := mux.Vars(r) // url-encoded parameters
	by := vars["by"]
	searchString := vars["string"]
	err = ForwardAndReturnGet(config.Configuration.CourseManagementAddress+"courses"+"/"+by+"/"+searchString, w)
	if err != nil {
		MakeErrorResponse(w, http.StatusInternalServerError, "Api Gateway - Internal Server Error")
		log.Println("Api Gateway - Internal Server Error")
		return
	}

}

// AddCourseToStudent process the course creation request coming from the client validating the embedded access token. If the token
// is properly signed, it checks if the request comes from a student, because course creation is allowed for teachers
// only. Upon successful validation, the request is forwarded to the microservice. The first time a course is added the student
// does not exist. In this case, the function provides its creation.
func AddCourseToStudent(w http.ResponseWriter, r *http.Request) {

	/* For authentication purpose the access token is read from the Cookie header and validated*/
	tokenString, err := GetToken(w, r)
	if err != nil {
		return
	}
	decodedToken, err := ValidateToken(tokenString, w)
	if err != nil {
		return
	}
	/* Upon successful validation, the request is forwarded to the course management microservice*/
	vars := mux.Vars(r)
	studentUsername := vars["username"]
	courseId := vars["id"]
	httpClient := &http.Client{}
	putRequest, err := http.NewRequest(http.MethodPut, config.Configuration.CourseManagementAddress+"students/"+studentUsername+"/courses/"+courseId, r.Body)
	putResponse, err := httpClient.Do(putRequest)
	if err != nil {
		MakeErrorResponse(w, http.StatusInternalServerError, "Api Gateway - Internal Server Error")
		log.Println("Api Gateway - Internal Server Error")
		return
	}
	log.Println("Response status Code from Microservice: " + strconv.Itoa(putResponse.StatusCode))
	defer putResponse.Body.Close()

	/* It the student does not exist, it is created by the api-gateway.*/
	if putResponse.StatusCode == http.StatusNotFound {
		var errorResponse ErrorResponse
		// Decode the microservice error putResponse
		jsonDecoder := json.NewDecoder(putResponse.Body)
		err = jsonDecoder.Decode(&errorResponse)
		if err != nil {
			MakeErrorResponse(w, http.StatusInternalServerError, "Api Gateway - Internal Server Error")
			log.Println("Internal Server Error")
			return
		}
		if errorResponse.Error == "Student Not Found" {

			studentCreationRequest := simplejson.New()
			studentCreationRequest.Set("name", decodedToken.Name+" "+decodedToken.Surname)
			studentCreationRequest.Set("username", studentUsername)
			studentCreationRequestPayload, err := studentCreationRequest.MarshalJSON()
			postResponse, err := http.Post(config.Configuration.CourseManagementAddress+"students", "application/json", bytes.NewBuffer(studentCreationRequestPayload))
			if err != nil {
				MakeErrorResponse(w, http.StatusInternalServerError, "Api Gateway - Internal Server Error")
				log.Println("Api Gateway - Internal Server Error")
				return
			}
			defer postResponse.Body.Close()
			if postResponse.StatusCode == http.StatusCreated {
				// Upon successful student creation proceed with course append
				// That is, the PUT request is re-forwarded to the microservice and the response is returned to the client
				err = ForwardAndReturnPut(config.Configuration.CourseManagementAddress+"students/"+studentUsername+"/courses/"+courseId, w, r)
				if err != nil {
					MakeErrorResponse(w, http.StatusInternalServerError, "Api Gateway - Internal Server Error")
					log.Println("Api Gateway - Internal Server Error")
					return
				}
				return
			}
			// If for some unexpected reason the student has not been created, an Internal Server Error is raised
			MakeErrorResponse(w, http.StatusInternalServerError, "Api Gateway - Internal Server Error")
			log.Println("Api Gateway - Internal Server Error")
			return
		}
	}
	/*Any other putResponse from the microservice is simply forwarded to the client*/
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(putResponse.StatusCode)
	putResponseBody, err := ioutil.ReadAll(putResponse.Body)
	if err != nil {
		MakeErrorResponse(w, http.StatusInternalServerError, "Api Gateway - Internal Server Error")
		log.Println("Api Gateway - Internal Server Error")
		return
	}
	_, err = w.Write(putResponseBody)
	if err != nil {
		MakeErrorResponse(w, http.StatusInternalServerError, "Api Gateway - Internal Server Error")
		log.Println("Api Gateway - Internal Server Error")
		return
	}
}

// FindStudentCourses process the request coming from the client and validate the embedded access token. It verify
// if the token is properly signed and not expired. Upon successful validation, the request is forwarded to the microservice
// and the response is returned, as-is, to the client.
func FindStudentCourses(w http.ResponseWriter, r *http.Request) {
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
	vars := mux.Vars(r) // url-encoded parameters
	studentUsername := vars["username"]
	err = ForwardAndReturnGet(config.Configuration.CourseManagementAddress+"courses/students/"+studentUsername, w)
	if err != nil {
		MakeErrorResponse(w, http.StatusInternalServerError, "Api Gateway - Internal Server Error")
		log.Println("Api Gateway - Internal Server Error")
		return
	}
}

// UnsubscribeStudentFromCourse process the request of canceling the subscription to a course validating the embedded token.
// Upon successfully validation, the request is forwarded to the course management microservice and the response is returned
// to the client
func UnsubscribeStudentFromCourse(w http.ResponseWriter, r *http.Request) {
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
	vars := mux.Vars(r) // url-encoded parameters
	studentUsername := vars["username"]
	courseId := vars["id"]
	err = ForwardAndReturnDelete(config.Configuration.CourseManagementAddress+"students/"+studentUsername+"/courses/"+courseId, w)
	if err != nil {
		MakeErrorResponse(w, http.StatusInternalServerError, "Api Gateway - Internal Server Error")
		log.Println("Api Gateway - Internal Server Error")
		return
	}
}
