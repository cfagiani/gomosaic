need to install via go:

go get github.com/nfnt/resize



Thoughts:
indexer - scan set of directories, resize and compute avg color values; write tiles and index to another dir; if flag is set also allow it to leave a .progress file in a dir with list of files it has processed so we don't re-scan
thumbnails should be square use 3x3 grid (so 9 values) for each image

mosaic - divide source into grid, compute avg color, search index for best match. allow params to limit reuse and tolerance
for matcH: divide target into 3x3 grid and find best match across all 9 positions



distance:
sum sq differences sqrt((sourceR - targetR)^2 + (sourceG - targetG)^2 + (sourceB - targetB)^2)
