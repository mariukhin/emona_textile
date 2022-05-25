package api

import (
	"github.com/h2non/filetype"
	"github.com/h2non/filetype/matchers"
	"github.com/h2non/filetype/types"
)

var image = matchers.Map{
	matchers.TypeJpeg:     matchers.Jpeg,
	matchers.TypeJpeg2000: matchers.Jpeg2000,
	matchers.TypePng:      matchers.Png,
	matchers.TypeGif:      matchers.Gif,
}

func IsImage(buf []byte) bool {
	kind, _ := doMatchMap(buf, image)
	return kind != types.Unknown
}

func doMatchMap(buf []byte, machers matchers.Map) (types.Type, error) {
	kind := filetype.MatchMap(buf, machers)
	if kind != types.Unknown {
		return kind, nil
	}
	return kind, filetype.ErrUnknownBuffer
}
