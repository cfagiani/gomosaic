package processor

import (
	"fmt"
	"github.com/cfagiani/gomosaic"
	"github.com/cfagiani/gomosaic/mosaicimages"
	"github.com/cfagiani/gomosaic/util"
	"google.golang.org/api/photoslibrary/v1"
	"log"
)

type GooglePhotosProcess struct {
	Config gomosaic.Config
	Source gomosaic.ImageSource
}

const DesiredPageSize = 500
const IndexTileDimension = "=w100"

func (p GooglePhotosProcess) Process(oldIndex gomosaic.MosaicTiles, newIndex gomosaic.MosaicTiles) gomosaic.MosaicTiles {

	photoService, err := util.GetPhotosService(p.Config.GoogleClientId, p.Config.GoogleClientSecret, p.Source.Options)

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
