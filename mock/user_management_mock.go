/* Package mock implements the functions used to simulate the microservice behaviour for testing purpose.*/
package mock

import (
	"github.com/bitly/go-simplejson"
	"github.com/gorilla/mux"
	"github.com/redefik/sdccproject/apigateway/config"
	"github.com/redefik/sdccproject/apigateway/microservice"
	"log"
	"net/http"
)

// UserManagementMockRegisterUser simulates the behaviour of the user-management microservice when receives a request of
// user registration.
func UserManagementMockRegisterUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := simplejson.New()
	response.Set("mock", "response")
	responsePayload, err := response.MarshalJSON()
	if err != nil {
		log.Panicln(err)
	}
	w.Write(responsePayload)
}

// UserManagementMockLoginUser simulates the behaviour of the user-management microservice when receives a get request
// for retrieving information about an user given username and password. It is based on the assumption that only an user
// exists with username "admin" and password "admin_pass".
func UserManagementMockLoginUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	username := params["username"]
	password := params["password"]
	if username != "admin" || password != "admin_pass" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	user := microservice.User{Username: "admin",
		Password: "admin_pass",
		Name:     "name",
		Surname:  "surname",
		Type:     "teacher"}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	response := simplejson.New()
	response.Set("user", user)
	responsePayload, err := response.MarshalJSON()
	if err != nil {
		log.Panicln(err)
	}
	w.Write(responsePayload)
}

// starts a user management microservice mock
func LaunchUserManagementMock() {
	r := mux.NewRouter()
	r.HandleFunc("/user_management/api/v1.0/users", UserManagementMockRegisterUser).Methods(http.MethodPost)
	r.HandleFunc("/user_management/api/v1.0/users/{username}/{password}", UserManagementMockLoginUser).Methods(http.MethodGet)
	http.ListenAndServe(config.Configuration.ApiGatewayAddress, r)
}
