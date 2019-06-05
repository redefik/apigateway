package microservice

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

/* ForwardAndReturnPost forwards a http post request from client to microservice */
func ForwardAndReturnPost(url string, contentType string, w http.ResponseWriter, r *http.Request) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	resp, err := http.Post(url, contentType, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	_, err = w.Write(responseBody)
	if err != nil {
		return err
	}
	return nil
}

/* ForwardAndReturnGet fowards a http get request from client to microservice.
All the parameters passed from client to api-gateway are url encoded in the get request from api-gateway to microservice */
func ForwardAndReturnGet(url string, w http.ResponseWriter) error {

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	log.Println("Response status Code from Microservice: " + strconv.Itoa(resp.StatusCode))
	defer resp.Body.Close()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	_, err = w.Write(responseBody)
	if err != nil {
		return err
	}
	return nil
}

func ForwardAndReturnPut(url string, w http.ResponseWriter, r *http.Request) error {
	httpClient := &http.Client{}
	putRequest, err := http.NewRequest(http.MethodPut, url, r.Body)
	putResponse, err := httpClient.Do(putRequest)
	if err != nil {
		return err
	}
	log.Println("Response status Code from Microservice: " + strconv.Itoa(putResponse.StatusCode))
	defer putResponse.Body.Close()
	w.WriteHeader(putResponse.StatusCode)
	responseBody, err := ioutil.ReadAll(putResponse.Body)
	if err != nil {
		return err
	}
	_, err = w.Write(responseBody)
	if err != nil {
		return err
	}
	return nil
}

func ForwardAndReturnDelete(url string, w http.ResponseWriter) error {
	httpClient := &http.Client{}
	putRequest, err := http.NewRequest(http.MethodDelete, url, nil)
	putResponse, err := httpClient.Do(putRequest)
	if err != nil {
		return err
	}
	log.Println("Response status Code from Microservice: " + strconv.Itoa(putResponse.StatusCode))
	defer putResponse.Body.Close()
	w.WriteHeader(putResponse.StatusCode)
	responseBody, err := ioutil.ReadAll(putResponse.Body)
	if err != nil {
		return err
	}
	_, err = w.Write(responseBody)
	if err != nil {
		return err
	}
	return nil
}
