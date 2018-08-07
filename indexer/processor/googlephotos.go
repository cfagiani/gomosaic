package processor

import (
	"encoding/json"
	"fmt"
	"github.com/cfagiani/gomosaic"
	"github.com/cfagiani/gomosaic/mosaicimages"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/photoslibrary/v1"
	"log"
	"os"
)

type GooglePhotosProcess struct {
	Config gomosaic.Config
	Source gomosaic.ImageSource
}

const DesiredPageSize = 500
const IndexTileDimension = "=w100"

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
	count := 0

	if err == nil {
		//TODO: handle album restriction
		var nextPage = ""
		for {
			pageResp := getPage(photoService, "", nextPage)
			//TODO: refactor this as most of the logic is the same as the localdir indexer.
			for _, item := range pageResp.MediaItems {
				existingTile := find(item.Id, oldIndex)
				if existingTile == nil {
					imageSegment, err := mosaicimages.AnalyzeImage(item.BaseUrl + IndexTileDimension)
					if err == nil {
						//now add to index
						newIndex = append(newIndex,
							gomosaic.MosaicTile{Loc: "G", Filename: item.Id, AvgR: imageSegment.RVal, AvgG: imageSegment.GVal, AvgB: imageSegment.BVal})
						count++
					}
				} else {
					newIndex = append(newIndex, *existingTile)
				}
			}
			nextPage = pageResp.NextPageToken
			if len(nextPage) == 0 {
				break
			}
		}
	}
	log.Printf("Added %d new files to index\n", count)
	return newIndex
}

func getPage(photoService *photoslibrary.Service, albumId string, nextPageToken string) *photoslibrary.SearchMediaItemsResponse {
	resp, apiErr := photoService.MediaItems.Search(&photoslibrary.SearchMediaItemsRequest{AlbumId: albumId,
		PageSize:  DesiredPageSize,
		PageToken: nextPageToken}).Do()
	if apiErr != nil {
		fmt.Printf("Could not fetch results from service: %v", apiErr)
	}
	return resp
}

func readToken(file string) *oauth2.Token {
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
