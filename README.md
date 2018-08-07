gomosaic
=====
Simple utility to build photo mosaics from a set of local images. This is mainly an exercise to learn Go. 

# Overview
This utility has two main components: indexer & mosaicmaker.
## indexer
This module will analyze all the images in a set of directories. For each image, an entry is added to an index file that contains the path to the file as well as the average RGB pixel values.
If the index file already exists, entries will be preserved (they will not be re-analyzed).

## mosaicmaker
This module uses the index file created by the indexer and a source image to generate a photo mosaic with a configurable grid/tile size. It will divide the source image into a grid of square segments of a (configurable) uniform size. For each grid segment, it will select the best matching tile (baring duplicates) and use that in the mosaic. The selected mosaic tiles will be resized to a (configurable) square size (regardless of source aspect ratio) as they are written to the destination image. 

# Dependencies
Google photos api: go get google.golang.org/api/photoslibrary/v1
Google OAUTH2: go get golang.org/x/oauth2
Google Cloud Go: go get -u cloud.google.com/go/...
oauth2-noserver (to get token): go get github.com/nmrshll/oauth2-noserver


# Usage
This application can be used via the CLI (using the cmd packages) or as a library by directly importing gomosaic/indexer and gomosaic/mosaicmaker


## Testing
To test, run 
'go test ./...'


## CLI
### Indexer
Indexer takes 2 command line arguments: a comma-delimited list of image directories and the fully-qualified path to the inedex file
#### Example
'go run cmd/indexer/main.go /home/user/images,/home/user/downloads /home/myindex.dat'
This will (recursively) scan */home/user/images* and */home/user/downloads* for image files and index them. Index will be written to /home/myindex.dat.
 
### Mosaicmaker
Mosaicmaker takes 5 arguments: sourceImage, indexFile, gridSize, tileSize and output file.
#### Example
'go run cmd/mosaicmaker.go myimg.jpg myindex.dat 10 50 mymosaic.jpg'
This will divide myimg.jpg up into a grid of 10x10 squares and will generate a mosaic using the tiles in the myindex.dat file. The resulting mosaic will be saved as mymosaic.jpg and each tile will be 50x50. So, if the input image was 100x100, the resultant mosaic would be 500x500.

#### TODO:
* unit tests
* better error handling
* Stat existing index entries on re-index & remove any files that are gone
* better data structure for index searches so we don't have to do so many comparisons
* options regarding how we want to handle duplicates (allow/disallow, min separation, etc)
* consider cropping on resize to prevent skew when making square tiles


#### Potential Enhancements:
* integration with google photos
* use 3x3 value matrix for pixel values and find best match of that
* optionally resize tiles? - will yield much faster runs