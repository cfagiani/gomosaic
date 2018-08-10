package processor

import (
	"github.com/cfagiani/gomosaic"
	"github.com/cfagiani/gomosaic/mosaicimages"
	"github.com/cfagiani/gomosaic/util"
	"google.golang.org/api/photoslibrary/v1"
	"log"
)

type GooglePhotosProcessor struct {
	Config gomosaic.Config
	Source gomosaic.ImageSource
}

const (
	desiredPageSize    = 500
	indexTileDimension = "=w100"
	GoogleKind         = "google"
)

//Process will populate the index of MosaicTiles by querying the Google Photos api to get a list of mediaItems and then
//analyzing each to calculate average pixel values.
func (p GooglePhotosProcessor) Process(oldIndex gomosaic.MosaicTiles, newIndex gomosaic.MosaicTiles) gomosaic.MosaicTiles {

	photoService, err := util.GetPhotosService(p.Config.GoogleClientId, p.Config.GoogleClientSecret, p.Source.Options)

	count := 0

	if err == nil {
		//TODO: handle album restriction
		var nextPage = ""
		for {
			pageResp := getPage(photoService, "", nextPage)
			//TODO: refactor this as most of the logic is the same as the localdir indexer.
			for _, item := range pageResp.MediaItems {
				if item.MediaMetadata.Photo != nil { // don't index videos
					existingTile := find(item.Id, oldIndex)
					if existingTile == nil {
						imageSegment, err := mosaicimages.AnalyzeImage(item.BaseUrl + indexTileDimension)
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

//getPage will fetch a page of MediaItems from the Google Photos api.
func getPage(photoService *photoslibrary.Service, albumId string, nextPageToken string) *photoslibrary.SearchMediaItemsResponse {
	resp, apiErr := photoService.MediaItems.Search(&photoslibrary.SearchMediaItemsRequest{AlbumId: albumId,
		PageSize:  desiredPageSize,
		PageToken: nextPageToken}).Do()
	if apiErr != nil {
		log.Printf("Could not fetch results from service: %v", apiErr)
	}
	return resp
}
