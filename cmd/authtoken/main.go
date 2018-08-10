package main

import (
	"encoding/json"
	"fmt"
	"github.com/cfagiani/gomosaic/util"
	"github.com/nmrshll/oauth2-noserver"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/photoslibrary/v1"
	"os"
)

func usage() {
	fmt.Println("Too few command line arguments.\n\nUsage:\n")
	fmt.Println("authtoken <configFile> <tokenFile>\n")
}

// Utility program to validate that we can get an authentication token using the user-supplied clientId and secret.
// It will update the config file with the token so it can be used by the indexer later.
func main() {
	if len(os.Args) != 3 {
		usage()
		os.Exit(1)
	}
	config, e := util.ReadConfig(os.Args[1])
	if e != nil {
		fmt.Printf("Could not read configuration file: %v\n", e)
		os.Exit(1)
	}

	conf := &oauth2.Config{
		ClientID:     config.GoogleClientId,
		ClientSecret: config.GoogleClientSecret,
		Scopes:       []string{photoslibrary.PhotoslibraryReadonlyScope},
		Endpoint:     google.Endpoint,
	}
	client := oauth2ns.Authorize(conf)
	fmt.Printf("Writing token to %s", os.Args[2])
	f, err := os.OpenFile(os.Args[2], os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		fmt.Printf("Could not write token file: %v", err)
		os.Exit(1)
	}
	json.NewEncoder(f).Encode(client.Token)
}
