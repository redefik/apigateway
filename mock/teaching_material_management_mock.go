package mock

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/redefik/sdccproject/apigateway/config"
	"log"
	"net/http"
)

func LaunchTeachingMaterialManagementMock() {
	r := mux.NewRouter()
	r.HandleFunc("/teaching_material_management/api/v1.0/list/{courseId}",
		TeachingMaterialManagementMockFindTeachingMaterialByCourse).Methods(http.MethodGet)
	r.HandleFunc("/teaching_material_management/api/v1.0/download/{fileName}",
		TeachingMaterialManagementMockGetDownloadLink).Methods(http.MethodGet)
	_ = http.ListenAndServe(config.Configuration.ApiGatewayAddress+"80", r)
}

// TeachingMaterialManagementMockGetDownloadLink simulates the behaviour of teaching management
// micro-service upon receiving a request getting download link. In the data store there are file1 for course1 and course3
// but not file2 for course1
func TeachingMaterialManagementMockGetDownloadLink(w http.ResponseWriter, r *http.Request) {

	var response []byte
	var err error
	vars := mux.Vars(r)

	if vars["fileName"] == "course1_file1" || vars["fileName"] == "course3_file1" || vars["fileName"] == "course4_file1"{
		w.WriteHeader(http.StatusOK)
		response, err = json.Marshal("validDownloadLink")
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

// TeachingMaterialManagementMockFindTeachingMaterialByCourse simulates the behaviour of teaching management
// micro-service upon receiving a request of listing teaching material. If provided idCourse is
// "courseIdWithTeachingMaterial" two file are found, otherwise no file are found.
func TeachingMaterialManagementMockFindTeachingMaterialByCourse(w http.ResponseWriter, r *http.Request) {

	var response []byte
	var err error

	if mux.Vars(r)["courseId"] == "courseIdWithTeachingMaterial" {
		response, err = json.Marshal(&([]string{"file1", "file2"}))
	} else {
		response, err = json.Marshal(&([]string{}))
	}
	if err != nil {
		log.Panicln(err)
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		log.Panicln(err)
	}

}
