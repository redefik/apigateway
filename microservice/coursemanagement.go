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
)

// FindCourse process the course searching request coming from the client and validate the embedded access token. It verify
// if the token is properly signed and not expired. Upon successful validation, the request is forwarded to the micro-service
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
	/* Upon successful validation, the request is forwarded to the course management micro-service and the response is
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
// only. Upon successful validation, the request is forwarded to the micro-service. The first time a course is added the student
// does not exist. In this case, the function provides its creation.
func AddCourseToStudent(w http.ResponseWriter, r *http.Request) {

	// For authentication purpose the access token is read from the Cookie header and validated, it's also verified if
	// requester is a student
	tokenString, err := GetToken(w, r)
	if err != nil {
		MakeErrorResponse(w, http.StatusUnauthorized, "Permission denied")
		log.Println("Permission denied")
		return
	}
	decodedToken, err := ValidateToken(tokenString, w)
	if err != nil {
		MakeErrorResponse(w, http.StatusUnauthorized, "Permission denied")
		log.Println("Permission denied")
		return
	}
	if decodedToken.Type != "student" {
		MakeErrorResponse(w, http.StatusUnauthorized, "Permission denied")
		log.Println("Permission denied")
		return
	}

	/* Upon successful validation a distributed transaction starts: api gateway send to course management and notification
	management micro-services a request to register the user to course in their own data-store. The request succeeds only
	if the operation is completed by both micro-services. The requests are send in parallel using goroutines. */

	// Collecting parameters for requests to course management and notification management micro-services
	studentUsername := mux.Vars(r)["username"]
	studentMail := decodedToken.Mail
	studentName := decodedToken.Name
	studentSurname := decodedToken.Surname
	var courseMinimized CourseMinimized
	var course Course
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		MakeErrorResponse(w, http.StatusInternalServerError, "API Gateway - Internal Server Error")
		log.Println("API Gateway - Internal Server Error")
		return
	}
	err = json.Unmarshal(body, &courseMinimized)
	if err != nil {
		MakeErrorResponse(w, http.StatusInternalServerError, "API Gateway - Internal Server Error")
		log.Println("API Gateway - Internal Server Error")
		return
	}
	err = json.Unmarshal(body, &course)
	if err != nil {
		MakeErrorResponse(w, http.StatusInternalServerError, "API Gateway - Internal Server Error")
		log.Println("API Gateway - Internal Server Error")
		return
	}

	//Initialize the channel to receive the exit of local transactions
	c := make(chan localTransaction, 2)

	//Launching goRoutines responsible to actuate local transaction
	go addSubscriptionInCourseManagement(studentUsername, studentName, studentSurname, courseMinimized.Id, c)
	go addSubscriptionInNotificationManagement(studentMail, course, c)

	isSentResponse := false     // Indicate if an internal error occurred and client already received a response
	var response *http.Response // The response for the client
	var localTransaction localTransaction
	var failingMicroservice []string // Contains the name of micro-service(s) that failed the execution of the request
	var i int

	for i = 0; i <= 1; i++ {
		// Waiting for the exit of local transactions
		localTransaction = <-c
		if localTransaction.Response == nil {
			// Any error occurred during forwarding of request: the client receive immediately an Internal Server Error
			failingMicroservice = append(failingMicroservice, localTransaction.Microservice)
			if isSentResponse == false {
				isSentResponse = true
				MakeErrorResponse(w, http.StatusInternalServerError, "Api Gateway - Internal Server Error")
				log.Println("Api Gateway - Internal Server Error")
			}
		} else if localTransaction.Response.StatusCode != http.StatusOK {
			// Failure: the client receive the last error response that the api-gateway obtained from micro-services
			failingMicroservice = append(failingMicroservice, localTransaction.Microservice)
			response = localTransaction.Response
		} else {
			// Success: the client receive the success response from course management micro-service.
			if localTransaction.Microservice == "courseManagement" {
				if response == nil {
					response = localTransaction.Response
				}
			}
		}
	}

	// If only a micro-service fail the other have to undo the action just completed
	if len(failingMicroservice) == 1 {
		if failingMicroservice[0] == "courseManagement" {
			removeSubscriptionInNotificationManagement(studentMail, course, nil)
		} else {
			removeSubscriptionInCourseManagement(studentUsername, courseMinimized.Id, nil)
		}
	}

	// If no Internal Server Error occurred the response from micro-services is forwarded to client
	if !isSentResponse {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(response.StatusCode)
		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			MakeErrorResponse(w, http.StatusInternalServerError, "Api Gateway - Internal Server Error")
			log.Println("Api Gateway - Internal Server Error")
			return
		}
		_, err = w.Write(responseBody)
		if err != nil {
			MakeErrorResponse(w, http.StatusInternalServerError, "Api Gateway - Internal Server Error")
			log.Println("Api Gateway - Internal Server Error")
			return
		}
	}

}

// addSubscriptionInCourseManagement send a request of course subscription to course management micro-service.
// If the student to subscribe to the course is not present in data-store, he is created.
func addSubscriptionInCourseManagement(studentUsername string, studentName string, studentSurname, courseId string, channel chan localTransaction) {

	httpClient := &http.Client{}
	putRequest, err := http.NewRequest(http.MethodPut, config.Configuration.CourseManagementAddress+"students/"+
		studentUsername+"/courses/"+courseId, nil)
	if err != nil {
		if channel == nil {
			log.Panicln("Api Gateway - Consistency problem. Please, recover the system.")
		} else {
			channel <- localTransaction{"courseManagement", nil}
			return
		}
	}
	putResponse, err := httpClient.Do(putRequest)
	if err != nil {
		if channel == nil {
			log.Panicln("Api Gateway - Consistency problem. Please, recover the system.")
		} else {
			channel <- localTransaction{"courseManagement", nil}
			return
		}
	}

	// It the student does not exist, it is created by the api-gateway.
	// This case can not verify when the method is used to undo a previous operation
	if putResponse.StatusCode == http.StatusNotFound {
		var errorResponse ErrorResponse
		// Decode the microservice error putResponse
		jsonDecoder := json.NewDecoder(putResponse.Body)
		err = jsonDecoder.Decode(&errorResponse)
		if err != nil {
			channel <- localTransaction{"courseManagement", nil}
			return
		}
		if errorResponse.Error == "Student Not Found" {

			studentCreationRequest := simplejson.New()
			studentCreationRequest.Set("name", studentName+" "+studentSurname)
			studentCreationRequest.Set("username", studentUsername)
			studentCreationRequestPayload, err := studentCreationRequest.MarshalJSON()
			postResponse, err := http.Post(config.Configuration.CourseManagementAddress+"students",
				"application/json", bytes.NewBuffer(studentCreationRequestPayload))
			if err != nil {
				channel <- localTransaction{"courseManagement", nil}
				return
			}
			if postResponse.StatusCode == http.StatusCreated {
				// Upon successful student creation proceed with course appending
				req, err := http.NewRequest(http.MethodPut, config.Configuration.CourseManagementAddress+"students/"+
					studentUsername+"/courses/"+courseId, nil)
				resp, err := httpClient.Do(req)
				if err != nil {
					channel <- localTransaction{"courseManagement", nil}
					return
				}
				channel <- localTransaction{"courseManagement", resp}
				return
			}
			// If for some unexpected reason the student has not been created, an error is communicated to main thread
			channel <- localTransaction{"courseManagement", nil}
			return
		}
	}

	if putResponse.StatusCode != http.StatusOK && channel == nil {
		log.Panicln("Api Gateway - Consistency problem. Please, recover the system.")
	}

	if channel != nil {
		// Any other putResponse from the micro-service is simply forwarded to the client
		channel <- localTransaction{"courseManagement", putResponse}
	}
}

// removeSubscriptionInCourseManagement send a request to remove a course subscription to course management
// micro-service for the specified user.
func removeSubscriptionInCourseManagement(studentUsername string, courseId string, channel chan localTransaction) {

	httpClient := &http.Client{}
	req, err := http.NewRequest(http.MethodDelete, config.Configuration.CourseManagementAddress+"students/"+
		studentUsername+"/courses/"+courseId, nil)
	if err != nil {
		if channel == nil {
			log.Panicln("Api Gateway - Consistency problem. Please, recover the system.")
		} else {
			channel <- localTransaction{"courseManagement", nil}
			return
		}
	}
	resp, err := httpClient.Do(req)

	if (err != nil || resp.StatusCode != http.StatusOK) && channel == nil {
		log.Panicln("Api Gateway - Consistency problem. Please, recover the system.")
	} else if err != nil && channel != nil {
		channel <- localTransaction{"courseManagement", nil}
		return
	}
	if channel != nil {
		channel <- localTransaction{"courseManagement", resp}
	}

}

// removeSubscriptionInNotificationManagement send a request to remove a course subscription to notification management
// micro-service for the specified user. Channel is the chan through communicate with main thread. If channel is null
// it means the function is used as undo method because transaction fail. If an error occurred during undoing operation
// a message is show to allow system administrator to recover the system
func removeSubscriptionInNotificationManagement(studentMail string, course Course, channel chan localTransaction) {

	body, err := json.Marshal(course)
	if err != nil {
		if channel == nil {
			log.Panicln("Api Gateway - Consistency problem. Please, recover the system.")
		} else {
			channel <- localTransaction{"notificationManagement", nil}
			return
		}
	}
	httpClient := &http.Client{}
	req, err := http.NewRequest(http.MethodDelete, config.Configuration.NotificationManagementAddress+
		"course/student/"+studentMail, bytes.NewBuffer(body))
	resp, err := httpClient.Do(req)
	if (err != nil || resp.StatusCode != http.StatusOK) && channel == nil {
		log.Panicln("Api Gateway - Consistency problem. Please, recover the system.")
	} else if err != nil && channel != nil {
		channel <- localTransaction{"notificationManagement", nil}
		return
	}
	if channel != nil {
		channel <- localTransaction{"notificationManagement", resp}
	}
}

// addSubscriptionInNotificationManagement send a request of course subscription to notification management micro-service.
//Channel is the chan through communicate with main thread. If channel is null it means the function is used as undo
// method because transaction fail. If an error occurred during undoing operation a message is show to allow system
// administrator to recover the system
func addSubscriptionInNotificationManagement(studentMail string, course Course, channel chan localTransaction) {

	body, err := json.Marshal(course)
	if err != nil {
		if channel == nil {
			log.Panicln("Api Gateway - Consistency problem. Please, recover the system.")
		} else {
			channel <- localTransaction{"notificationManagement", nil}
			return
		}
	}
	httpClient := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, config.Configuration.NotificationManagementAddress+
		"course/student/"+studentMail, bytes.NewBuffer(body))
	resp, err := httpClient.Do(req)
	if (err != nil || resp.StatusCode != http.StatusOK) && channel == nil {
		log.Panicln("Api Gateway - Consistency problem. Please, recover the system.")
	} else if err != nil && channel != nil {
		channel <- localTransaction{"notificationManagement", nil}
		return
	}
	if channel != nil {
		channel <- localTransaction{"notificationManagement", resp}
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
// Upon successful validation, the request is forwarded to the micro-services of course management and notification management.
func UnsubscribeStudentFromCourse(w http.ResponseWriter, r *http.Request) {

	// For authentication purpose the access token is read from the Cookie header and validated, it's also verified if
	// requester is a student
	tokenString, err := GetToken(w, r)
	if err != nil {
		MakeErrorResponse(w, http.StatusUnauthorized, "Permission denied")
		log.Println("Permission denied")
		return
	}
	decodedToken, err := ValidateToken(tokenString, w)
	if err != nil {
		MakeErrorResponse(w, http.StatusUnauthorized, "Permission denied")
		log.Println("Permission denied")
		return
	}
	if decodedToken.Type != "student" {
		MakeErrorResponse(w, http.StatusUnauthorized, "Permission denied")
		log.Println("Permission denied")
		return
	}

	/* Upon successful validation a distributed transaction starts: api gateway send to course management and notification
	management micro-services a request to deregister the user to course in their own data-store. The request succeeds only
	if the operation is completed by both micro-services. The requests are send in parallel using goroutines. */

	// Collecting parameters for requests to course management and notification management micro-services
	studentUsername := mux.Vars(r)["username"]
	studentMail := decodedToken.Mail
	var courseMinimized CourseMinimized
	var course Course
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		MakeErrorResponse(w, http.StatusInternalServerError, "API Gateway - Internal Server Error")
		log.Println("API Gateway - Internal Server Error")
		return
	}
	err = json.Unmarshal(body, &courseMinimized)
	if err != nil {
		MakeErrorResponse(w, http.StatusInternalServerError, "API Gateway - Internal Server Error")
		log.Println("API Gateway - Internal Server Error")
		return
	}
	err = json.Unmarshal(body, &course)
	if err != nil {
		MakeErrorResponse(w, http.StatusInternalServerError, "API Gateway - Internal Server Error")
		log.Println("API Gateway - Internal Server Error")
		return
	}

	//Initialize the channel to receive the exit of local transactions
	c := make(chan localTransaction, 2)

	//Launching goRoutines responsible to actuate local transaction
	go removeSubscriptionInCourseManagement(studentUsername, courseMinimized.Id, c)
	go removeSubscriptionInNotificationManagement(studentMail, course, c)

	isSentResponse := false     // Indicate if an internal error occurred and client already received a response
	var response *http.Response // The response for the client
	var localTransaction localTransaction
	var failingMicroservice []string // Contains the name of micro-service(s) that failed the execution of the request
	var i int

	for i = 0; i <= 1; i++ {
		// Waiting for the exit of local transactions
		localTransaction = <-c
		if localTransaction.Response == nil {
			// Any error occurred during forwarding of request: the client receive immediately an Internal Server Error
			failingMicroservice = append(failingMicroservice, localTransaction.Microservice)
			if isSentResponse == false {
				isSentResponse = true
				MakeErrorResponse(w, http.StatusInternalServerError, "Api Gateway - Internal Server Error")
				log.Println("Api Gateway - Internal Server Error")
			}
		} else if localTransaction.Response.StatusCode != http.StatusOK {
			// Failure: the client receive the last error response that the api-gateway obtained from micro-services
			failingMicroservice = append(failingMicroservice, localTransaction.Microservice)
			response = localTransaction.Response
		} else {
			// Success: the client receive the success response from course management micro-service.
			if localTransaction.Microservice == "courseManagement" {
				if response == nil {
					response = localTransaction.Response
				}
			}
		}
	}

	// If only a micro-service fail the other have to undo the action just completed
	if len(failingMicroservice) == 1 {
		if failingMicroservice[0] == "courseManagement" {
			addSubscriptionInNotificationManagement(studentMail, course, nil)
		} else {
			addSubscriptionInCourseManagement(studentUsername, "", "", courseMinimized.Id, nil)
		}
	}

	// If no Internal Server Error occurred the response from micro-services is forwarded to client
	if !isSentResponse {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(response.StatusCode)
		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			MakeErrorResponse(w, http.StatusInternalServerError, "Api Gateway - Internal Server Error")
			log.Println("Api Gateway - Internal Server Error")
			return
		}
		_, err = w.Write(responseBody)
		if err != nil {
			MakeErrorResponse(w, http.StatusInternalServerError, "Api Gateway - Internal Server Error")
			log.Println("Api Gateway - Internal Server Error")
			return
		}
	}
}

// CreateCourse process the course creation request coming from the client validating the embedded access token. If the token
// is properly signed, it checks if the request comes from a student, because course creation is allowed for teachers
// only. Upon successful validation, the request is forwarded to the micro-service and the response is forwarded to the
// client.
func CreateCourse(w http.ResponseWriter, r *http.Request) {

	// For authentication purpose the access token is read from the Cookie header
	tokenString, err := GetToken(w, r)
	if err != nil {
		return
	}
	// The token is decoded and the claims are obtained for further checks
	decodedToken, err := ValidateToken(tokenString, w)
	if err != nil {
		return
	}
	// Upon successful authentication check if the request comes from a teacher.
	// In case of student request an Unauthorized code is returned
	if decodedToken.Type != "teacher" {
		MakeErrorResponse(w, http.StatusUnauthorized, "Permission denied")
		log.Println("Permission denied")
		return
	}

	/* Upon successful validation a distributed transaction starts: api gateway send to course management and notification
	management micro-services a request to create course in their own data-store. The creation of course succeed only if the
	operation is completed by both micro-services. The requests are send in parallel using goroutines.  */

	//Initialize the channel to receive the exit of local transactions
	c := make(chan localTransaction, 2)
	//Launching goRoutines responsible to actuate local transactions
	requestBody, _ := ioutil.ReadAll(r.Body)
	go createCourseInCourseManagement(requestBody, c)
	go createCourseInNotificationManagement(requestBody, c)

	isSentResponse := false     // Indicate if an internal error occurred and client already received a response
	var response *http.Response // The response for the client
	var localTransaction localTransaction
	var failingMicroservice []string             // Contains the name of micro-service(s) that failed the execution of the request
	var courseInCourseManagement CourseMinimized // Contains the id of course inserted in course management micro-service
	var courseInNotificationManagement Course    // Contains the name of course inserted in notification micro-service
	var i int

	for i = 0; i <= 1; i++ {
		// Waiting for the exit of local transactions
		localTransaction = <-c
		if localTransaction.Response == nil {
			// Any error occurred during forwarding of request: the client receive immediately an Internal Server Error
			failingMicroservice = append(failingMicroservice, localTransaction.Microservice)
			if isSentResponse == false {
				isSentResponse = true
				MakeErrorResponse(w, http.StatusInternalServerError, "Api Gateway - Internal Server Error")
				log.Println("Api Gateway - Internal Server Error")
			}
		} else if localTransaction.Response.StatusCode != http.StatusCreated {
			// Failure: the client receive the last error response that the api-gateway obtained from micro-services
			failingMicroservice = append(failingMicroservice, localTransaction.Microservice)
			response = localTransaction.Response
		} else {
			// Success: the client receive the success response from course management micro-service.
			// The id of inserted course in micro-services are stored for eventually future deletion
			if localTransaction.Microservice == "courseManagement" {
				if response == nil {
					response = localTransaction.Response
				}
				body, _ := ioutil.ReadAll(localTransaction.Response.Body)
				_ = json.Unmarshal(body, &courseInCourseManagement)
			} else if localTransaction.Microservice == "notificationManagement" {
				body, _ := ioutil.ReadAll(localTransaction.Response.Body)
				_ = json.Unmarshal(body, &courseInNotificationManagement)
			}
		}
	}

	// If only a micro-service fail the other have to undo the action just completed
	if len(failingMicroservice) == 1 {
		if failingMicroservice[0] == "courseManagement" {
			deleteCourseInNotificationManagement(courseInNotificationManagement)
		} else {
			deleteCourseInCourseManagement(courseInCourseManagement.Id)
		}
	}

	// If no Internal Server Error occurred the response from micro-services is forwarded to client
	if !isSentResponse {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(response.StatusCode)
		responseBody, err := ioutil.ReadAll(response.Body)
		if err != nil {
			MakeErrorResponse(w, http.StatusInternalServerError, "Api Gateway - Internal Server Error")
			log.Println("Api Gateway - Internal Server Error")
			return
		}
		_, err = w.Write(responseBody)
		if err != nil {
			MakeErrorResponse(w, http.StatusInternalServerError, "Api Gateway - Internal Server Error")
			log.Println("Api Gateway - Internal Server Error")
			return
		}
	}
}

// createCourseInNotificationManagement send a request of course creation to notification management micro-service.
func createCourseInNotificationManagement(body []byte, channel chan localTransaction) {
	// Retrieving the name of course from the body of request
	var course Course
	err := json.Unmarshal(body, &course)
	if err != nil {
		// Communicating to main thread the failure of local transaction
		channel <- localTransaction{"notificationManagement", nil}
		return
	}
	requestBody, err := json.Marshal(course)
	if err != nil {
		// Communicating to main thread the failure of local transaction
		channel <- localTransaction{"notificationManagement", nil}
		return
	}
	// Send the post request to notification management micro-service
	resp, err := http.Post(config.Configuration.NotificationManagementAddress+"course",
		"application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		// Communicating to main thread the failure of local transaction
		channel <- localTransaction{"notificationManagement", nil}
		return
	}
	// Communicating to main thread the exit of local transaction
	channel <- localTransaction{"notificationManagement", resp}
}

// createCourseInCourseManagement send a request of course creation to course management micro-service.
func createCourseInCourseManagement(requestBody []byte, channel chan localTransaction) {
	// Send the post request to course management micro-service
	resp, err := http.Post(config.Configuration.CourseManagementAddress+"courses",
		"application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		// Communicating to main thread the failure of local transaction
		channel <- localTransaction{"courseManagement", nil}
		return
	}
	// Communicating to main thread the exit of local transaction
	channel <- localTransaction{"courseManagement", resp}
}

// deleteCourseInCourseManagement send a request of course deletion to course management micro-service.
func deleteCourseInCourseManagement(courseId string) {
	client := &http.Client{}
	request, err := http.NewRequest(http.MethodDelete, config.Configuration.CourseManagementAddress+"courses/"+courseId, nil)
	if err != nil {
		log.Panicln("Api Gateway - Consistency problem. Please, recover the system.")
	}
	response, err := client.Do(request)
	if err != nil || response.StatusCode != 200 {
		log.Panicln("Api Gateway - Consistency problem. Please, recover the system.")
	}

}

// deleteCourseInNotificationManagement send a request of course deletion to notification management micro-service.
func deleteCourseInNotificationManagement(course Course) {
	client := &http.Client{}
	body, err := json.Marshal(course)
	if err != nil {
		log.Panicln("Api Gateway - Consistency problem. Please, recover the system.")
	}
	request, _ := http.NewRequest(http.MethodDelete, config.Configuration.NotificationManagementAddress+"course", bytes.NewBuffer(body))
	response, err := client.Do(request)
	if err != nil || response.StatusCode != 200 {
		log.Println(err)
		log.Panicln("Api Gateway - Consistency problem. Please, recover the system.")
	}
}

// This struct encapsulates the exit of single local transaction. Its fields are used by main thread to orchestrate
// the other micro-service involved in distributed transaction. It is assumed an internal error if response is nil
type localTransaction struct {
	Microservice string
	Response     *http.Response
}
