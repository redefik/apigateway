package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

// Contains the configurable options of the api gateway
var Configuration Config

// Encapsulates the fields of the configuration file
type Config struct {
	ApiGatewayAddress                 string
	UserManagementAddress             string
	CourseManagementAddress           string
	TeachingMaterialManagementAddress string
	NotificationManagementAddress     string
	TokenPrivateKey                   string
}

func SetConfigurationFromFile(configFile string) error {
	jsonFile, err := os.Open(configFile)
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	bytes, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, &Configuration)
	if err != nil {
		return err
	}
	return nil
}

// Reads the configuration parameters from environment variables and stores them into the Config struct
func SetConfigurationFromEnvironment() error {
	gatewayAddress, present := os.LookupEnv("APIGATEWAY_ADDR")
	if !present {
		return errors.New("couldn't load configuration parameters")
	}
	Configuration.ApiGatewayAddress = gatewayAddress
	userManagementAddress, present := os.LookupEnv("USER_ADDR")
	if !present {
		return errors.New("couldn't load configuration parameters")
	}
	Configuration.UserManagementAddress = userManagementAddress
	courseManagementAddress, present := os.LookupEnv("COURSE_ADDR")
	if !present {
		return errors.New("couldn't load configuration parameters")
	}
	Configuration.CourseManagementAddress = courseManagementAddress
	teachingMaterialManagementAddress, present := os.LookupEnv("TEACHING_ADDR")
	if !present {
		return errors.New("couldn't load configuration parameters")
	}
	Configuration.TeachingMaterialManagementAddress = teachingMaterialManagementAddress
	notificationManagementAddress, present := os.LookupEnv("NOTIFICATION_ADDR")
	if !present {
		return errors.New("couldn't load configuration parameters")
	}
	Configuration.NotificationManagementAddress = notificationManagementAddress
	tokenPrivateKey, present := os.LookupEnv("TOKEN_PRIVATE_KEY")
	if !present {
		return errors.New("couldn't load configuration parameters")
	}
	Configuration.TokenPrivateKey = tokenPrivateKey
	return nil
}
