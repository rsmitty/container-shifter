package utils

import (
	b64 "encoding/base64"
	"encoding/json"
	"io/ioutil"
	"strings"

	log "github.com/Sirupsen/logrus"
)

//ErrorCheck provides error checking for the cmd package
func ErrorCheck(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func decodeEncodeUserPass(b64String string) string {
	b64Decode, err := b64.StdEncoding.DecodeString(b64String)
	ErrorCheck(err)
	authString := string(b64Decode)
	authArray := strings.SplitN(authString, ":", 2)
	log.Info(authArray)
	returnMap := map[string]string{"username": authArray[0], "password": authArray[1]}
	returnJSON, err := json.Marshal(returnMap)
	ErrorCheck(err)
	returnb64 := b64.StdEncoding.EncodeToString(returnJSON)

	return returnb64
}

//RegistryAuth parses the docker config.json path
func RegistryAuth(registry string, configPath string) string {
	configFile, _ := ioutil.ReadFile(configPath)
	var configData map[string]interface{}
	err := json.Unmarshal(configFile, &configData)
	ErrorCheck(err)
	auths := configData["auths"].(map[string]interface{})
	if value, ok := auths[registry]; ok {
		log.Info("Registry auth found, using existing credentials.")
		auth := value.(map[string]interface{})
		log.Info(auth["auth"].(string))
		return decodeEncodeUserPass(auth["auth"].(string))
	}

	log.Info("Registry auth info not found. Assuming no auth required.")
	return "noauth"

}
