package microservice

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/redefik/sdccproject/apigateway/config"
	"log"
	"net/http"
)

// FindTeachingMaterialByCourse process a request for listing teaching material about a specific course, in a specific
// department, in a specific academic year. It verify if the token is properly signed and not expired. Upon successful
// validation, the request is forwarded to the micro-service and the response is forwarded to the client.
func FindTeachingMaterialByCourse(w http.ResponseWriter, r *http.Request) {

	/* For authentication purpose the access token is read from the Cookie header */
	tokenString, err := GetToken(w, r)
	if err != nil {
		return
	}
	/* The token is decoded and validated */
	_, err = ValidateToken(tokenString, w)
	if err != nil {
		MakeErrorResponse(w, http.StatusUnauthorized, "Wrong Credentials")
		log.Println("Wrong credentials")
		return
	}

	/* Upon successful validation, the request is forwarded to the teaching management management microservice and
	the response is returned to the client*/
	vars := mux.Vars(r) // url-encoded parameters
	courseId := vars["courseId"]
	err = ForwardAndReturnGet(config.Configuration.TeachingMaterialManagementAddress+"list"+"/"+courseId, w)
	if err != nil {
		MakeErrorResponse(w, http.StatusInternalServerError, "Api Gateway - Internal Server Error")
		log.Println("Api Gateway - Internal Server Error")
		return
	}
}

// GetDownloadLinkToFile process a request for obtain a link usable to download the file with specified file name.
// It verify if the token is properly signed and not expired. Then the method checks if the user is allowed to
// download file: if he/she is a student checks if him/her is subscribed to course the file belong to; if he/she is
// a teacher checks if him/her hold the course the file belong to. Upon successful
// validation, the request is forwarded to the micro-service and the response is forwarded to the client.
func GetDownloadLinkToFile(w http.ResponseWriter, r *http.Request) {

	/* The token is decoded and validated */
	tokenString, err := GetToken(w, r)
	if err != nil {
		return
	}
	decodedToken, err := ValidateToken(tokenString, w)
	if err != nil {
		MakeErrorResponse(w, http.StatusUnauthorized, "Wrong Credentials")
		log.Println("Wrong credentials")
		return
	}

	vars := mux.Vars(r)          // URL-encoded parameters
	courseId := vars["courseId"] // Represent the id of the course the course to which the file belongs
	var query string             // Representing the url used to access to micro-service course management

	/* If requester is a teacher the method checks if him/her hold the course the file belong to */
	if decodedToken.Type == "teacher" {
		// Asking to course management micro-service for the list of courses hold by the teacher with name-included token
		teacherName := decodedToken.Name + "-" + decodedToken.Surname
		query = config.Configuration.CourseManagementAddress + "courses/teacher/" + teacherName

		/* If requester is a student the method checks if him/her is subscribed to the course the file belong to */
	} else if decodedToken.Type == "student" {
		// Asking to course management micro-service for the list of courses attended by the student with given username
		username := vars["username"]
		query = config.Configuration.CourseManagementAddress + "courses/students/" + username
	}

	resp, err := http.Get(query)
	// If any error occurred during interaction with micro-service the client receive an error response
	if err != nil || resp.StatusCode != http.StatusOK {
		MakeErrorResponse(w, http.StatusInternalServerError, "Api Gateway - Internal server Error")
		log.Println("Internal Server Error")
		return
	}

	// Decoding courses holding by the teacher or attended by student
	jsonDecoder := json.NewDecoder(resp.Body)
	defer resp.Body.Close()
	var courses []CourseMinimized
	err = jsonDecoder.Decode(&courses)
	if err != nil {
		MakeErrorResponse(w, http.StatusInternalServerError, "Api Gateway - Internal server Error")
		log.Println("Internal Server Error")
		return
	}
	// Checking if the course of the requested file is between the ones attended by the requester student
	for _, course := range courses {
		if course.Id == courseId {
			/* Upon successful validation, the request is forwarded to the teaching management management micro-service and
			the response is returned to the client*/
			filename := vars["fileName"]
			err = ForwardAndReturnGet(config.Configuration.TeachingMaterialManagementAddress+
				"download"+"/"+courseId+"_"+filename, w)
			if err != nil {
				MakeErrorResponse(w, http.StatusInternalServerError, "Api Gateway - Internal Server Error")
				log.Println("Api Gateway - Internal Server Error")
				return
			}
			return
		}
	}

	/* Upon failure during validation the client receive a Permission Denied error*/
	MakeErrorResponse(w, http.StatusUnauthorized, "Permission denied")
	log.Println("Permission denied")
	return
}
