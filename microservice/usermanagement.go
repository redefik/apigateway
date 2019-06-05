/* The package microservice implements the functions that handle the interaction between the api gateway and the microservices*/
package microservice

import (
	"encoding/json"
	"github.com/redefik/sdccproject/apigateway/config"
	"log"
	"net/http"
	"strconv"
)

// LoginUser makes an http Get request to the user management microservice in order to receive all the information needed
// to build the access token. Upon successful response, it generates the token and sends it to the client. In case of
// not found user the gateway communicates to the client that it has not the required authorization to login. Otherwise,
// it responds with a generic InternalServerError code.
func LoginUser(w http.ResponseWriter, r *http.Request) {

	var requestBody LoginRequestBody

	// Read username and password from the login request
	jsonDecoder := json.NewDecoder(r.Body)
	err := jsonDecoder.Decode(&requestBody)
	if err != nil {
		MakeErrorResponse(w, http.StatusBadRequest, "Bad Request")
		log.Println("Bad request")
		return
	}

	// Makes the get request to the microservice
	query := config.Configuration.UserManagementAddress + "users/" + requestBody.Username + "/" + requestBody.Password
	resp, err := http.Get(query)
	if err != nil {
		MakeErrorResponse(w, http.StatusInternalServerError, "Api Gateway - Internal server Error")
		log.Panicln(err)
		return
	}
	log.Println("Response status Code from Microservice: " + strconv.Itoa(resp.StatusCode))
	defer resp.Body.Close()

	// Checks the response from the user management microservice
	if resp.StatusCode == http.StatusOK {
		var responseBody LoginResponseBody

		// Decode the microservice response
		jsonDecoder = json.NewDecoder(resp.Body)
		err = jsonDecoder.Decode(&responseBody)
		if err != nil {
			MakeErrorResponse(w, http.StatusInternalServerError, "Api Gateway - Internal server Error")
			log.Println("Internal Server Error")
			return
		}

		// Generate the access token to be sent to the client
		token, err := GenerateAccessToken(responseBody.User, []byte(config.Configuration.TokenPrivateKey))
		if err != nil {
			MakeErrorResponse(w, http.StatusInternalServerError, "Api Gateway - Internal Server Error")
			log.Println("Internal Server Error")
			return
		}
		// The token is written in the 'Set-Cookie' field of the HTTP answer for the client
		http.SetCookie(w, &http.Cookie{
			Name:  "token",
			Value: token,
		})
		w.WriteHeader(http.StatusCreated)

	} else if resp.StatusCode == http.StatusNotFound {
		MakeErrorResponse(w, http.StatusUnauthorized, "Authentication failed - Wrong username or password")
		log.Println("Authentication failed - Wrong username or password")

	} else {
		MakeErrorResponse(w, http.StatusInternalServerError, "Api Gateway - Internal server Error")
		log.Println("Internal Server Error")
	}

}

// RegisterUser.md forwards the post request for registration to the user-management microservice,
// collect the response and sends it to the client
func RegisterUser(w http.ResponseWriter, r *http.Request) {

	err := ForwardAndReturnPost(config.Configuration.UserManagementAddress+"users", "application/json", w, r)
	if err != nil {
		MakeErrorResponse(w, http.StatusInternalServerError, "Api Gateway - Internal Server Error")
		log.Println("Api Gateway - Internal Server Error")
		return
	}
}
