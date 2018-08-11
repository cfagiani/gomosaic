package util

import (
	"encoding/json"
	"github.com/cfagiani/gomosaic"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/photoslibrary/v1"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"errors"
)

//utility to get the path of a file by concatenating the directory, the pathSeparator, and the file name.
func GetPath(dir string, file string) string {
	return dir + string(os.PathSeparator) + file
}

//Checks the error argument and, if it is not nil, it will log the msg passed in. If isFatal is true, the log will be
//written as Fatal which will cause exit(1) to be called.
func CheckError(err error, msg string, isFatal bool) bool {
	if err != nil {
		if isFatal {
			log.Fatal(msg, err)
		} else {
			log.Println(msg)
		}
		return true
	}
	return false
}

// Converts a string to a 32-bit unsigned integer, eating any errors
func GetInt32(s string) uint32 {
	// TODO: actually handle the error?
	i, err := strconv.ParseUint(s, 10, 32)
	if err == nil {
		return uint32(i)
	}
	return 0
}

//GetPhotosService will use the client information and token file passed in to initialize a photoslibrary.Service instance
//that can be used to interact with the Google Photos API.
func GetPhotosService(clientId string, clientSecret string, tokenFile string) (*photoslibrary.Service, error) {
	conf := &oauth2.Config{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		Scopes:       []string{photoslibrary.PhotoslibraryReadonlyScope},
		Endpoint:     google.Endpoint,
	}

	client := conf.Client(context.Background(), readOauthToken(tokenFile))

	return photoslibrary.New(client)
}

//readOauthToken reads a json file containing an oauth token and unmarshals it into a Token struct.
func readOauthToken(file string) *oauth2.Token {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		log.Fatalf("Could not load token: %v", err)
		os.Exit(1)
	}
	token := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(token)
	return token
}

//ReadConfig reads in a configuration json file and unmarshals it into a Config struct.
func ReadConfig(fileName string) (gomosaic.Config, error) {
	file, e := ioutil.ReadFile(fileName)
	if e != nil {
		return gomosaic.Config{}, e
	}
	var config gomosaic.Config
	e = json.Unmarshal(file, &config)
	if e == nil {
		if len(config.Sources) == 0 {
			return config, errors.New("config file did not contain any valid sources")
		}
	}
	return config, e
}
