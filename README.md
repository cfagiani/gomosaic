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
* Google photos api: go get google.golang.org/api/photoslibrary/v1
* Google OAUTH2: go get golang.org/x/oauth2
* Google Cloud Go: go get -u cloud.google.com/go/...
* oauth2-noserver (to get token): go get github.com/nmrshll/oauth2-noserver
* OpenURI (to open files/urls via same interface): go get -u github.com/utahta/go-openuri

# Installation
Run `make all` to build and install the binaries into your environment. This target will download the dependencies then build & install the commands into your GOPATH/bin directory. Once installed you will be able to run them without using 'go run' (i.e. as long as your PATH contains your GOPATH/bin directory you can run 'mosaicmaker' directly from the shell prompt)
# Usage
This application can be used via the CLI (using the cmd packages) or as a library by directly importing gomosaic/indexer and gomosaic/mosaicmaker


## Testing
To test, run 
`go test ./...`


## CLI

### Authtoken
This command can be run to obtain an OAuth access token from Google after the user authenticates and grants the app access. The utility will open a browser window (using the system default browser) through which the user can authenticate with Google. It will then start up a local server to accept the redirect from the browser in order to capture the access/refresh token.
#### Configuration
The command uses a configuration file that must contain __googleClientId__ and __googleClientSecret__ as top-level strings.
See the sample config  [here](https://github.com/cfagiani/gomosaic/blob/master/config-sample.json). You may obtain these keys from the [Google API Console](https://console.developers.google.com/).
#### Example
`go run cmd/authtoken/main.go config.json token.json`
This invocation will use the configuration stored in config.json to built an OAuth request and then open a browser window that allows the user to authenticate with Google. The resulting token is stored in token.json.

### Indexer
Indexer takes 2 command line arguments: the path to a configuration json file and the fully-qualified path to the index file
#### Configuration
The indexer is configured via a json file. A sample of that file can be seen [here](https://github.com/cfagiani/gomosaic/blob/master/config-sample.json).
The configuration file must contain an array of source objects:
    
    "sources": [
        {
          "kind": "",
          "path": "",
          "options": ""
        }
In the snippet above, the meaning of each field is as follows:
* __kind__ - either __google__ (to index a Google Photos account) or __local__ to index a local directory
* __path__ - either an Google Photos album name (or blank to index all photos in the account) or a local directory
* __options__ - for local this can be a __recurse__ which tells the indexer to recursively search the path location for images or, for the google indexer this is the path where the access token is stored.   

#### Example
`go run cmd/indexer/main.go config.json /home/myindex.dat`
This will (recursively) scan the locations in the config file for image files and index them. Index will be written to /home/myindex.dat.
 

 
### Mosaicmaker
Mosaicmaker takes 5 or 6 arguments: sourceImage, indexFile, gridSize, tileSize, output file, config file.
Config file is only used if the index contains tiles on google images.
#### Example
`go run cmd/mosaicmaker.go myimg.jpg myindex.dat 10 50 mymosaic.jpg`
This will divide myimg.jpg up into a grid of 10x10 squares and will generate a mosaic using the tiles in the myindex.dat file. The resulting mosaic will be saved as mymosaic.jpg and each tile will be 50x50. So, if the input image was 100x100, the resultant mosaic would be 500x500.

#### TODO:
* unit tests
* better error handling
* Stat existing index entries on re-index & remove any files that are gone
* Store index summary information including last index date so we can make indexers only look at things modiified since last index run
* better data structure for index searches so we don't have to do so many comparisons
* options regarding how we want to handle duplicates (allow/disallow, min separation, etc)
* consider cropping on resize to prevent skew when making square tiles
* refactor indexers to remove duplicate code
* refactor photo api client


#### Potential Enhancements:
* use 3x3 value matrix for pixel values and find best match of that
* optionally resize tiles? - will yield much faster runs
* center-crop after resize?
* parallelize indexing 