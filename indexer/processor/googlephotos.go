package processor

import (
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/photoslibrary/v1"
	"os"
	"github.com/cfagiani/gomosaic"
	"encoding/json"
)

type GooglePhotosProcess struct {
	Config gomosaic.Config
	Source gomosaic.ImageSource
}

func (p GooglePhotosProcess) Process(oldIndex gomosaic.MosaicTiles, newIndex gomosaic.MosaicTiles) gomosaic.MosaicTiles {

	conf := &oauth2.Config{
		ClientID:     p.Config.GoogleClientId,
		ClientSecret: p.Config.GoogleClientSecret,
		Scopes:       []string{photoslibrary.PhotoslibraryReadonlyScope},
		Endpoint:     google.Endpoint,
	}

	client := conf.Client(context.Background(), readToken(p.Source.Options))

	//TODO replace this with the code to actually do the index
	photoService, err := photoslibrary.New(client)
	if err == nil {
		resp, apiErr := photoService.MediaItems.Search(&photoslibrary.SearchMediaItemsRequest{PageSize: 10}).Do()
		if apiErr == nil {
			for i := 0; i < len(resp.MediaItems); i++ {
				fmt.Printf("Item: %s", resp.MediaItems[i].Id)
			}
		}
	}
	return newIndex
}

func readToken(file string) (*oauth2.Token) {
	f, err := os.Open(file)
	defer f.Close()
	if err != nil {
		fmt.Printf("Could not load token: %v", err)
		os.Exit(1)
	}
	token := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(token)
	return token
}
