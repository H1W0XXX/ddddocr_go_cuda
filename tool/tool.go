// tool/tool.go
package tool

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
)

// TypeError represents a custom error type
type TypeError struct {
	Msg string
}

func (e *TypeError) Error() string {
	return e.Msg
}

// NewTypeError creates a new TypeError
func NewTypeError(msg string) error {
	return &TypeError{Msg: msg}
}

// Base64ToImage decodes a base64 string into an image.Image
func Base64ToImage(imgBase64 string) (image.Image, error) {
	data, err := base64.StdEncoding.DecodeString(imgBase64)
	if err != nil {
		return nil, err
	}
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return img, nil
}

// GetImgBase64 reads an image file and returns its base64 encoding
func GetImgBase64(imagePath string) (string, error) {
	data, err := ioutil.ReadFile(imagePath)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

// PngRgbaBlackPreprocess pastes an RGBA image onto a white background
func PngRgbaBlackPreprocess(img image.Image) (image.Image, error) {
	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	// Create white background

	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	white := color.RGBA{255, 255, 255, 255}
	draw.Draw(dst, dst.Bounds(), &image.Uniform{white}, image.Point{}, draw.Src)
	// Overlay source image
	draw.Draw(dst, bounds, img, bounds.Min, draw.Over)
	return dst, nil
}
