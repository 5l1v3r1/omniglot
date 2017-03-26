package omniglot

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

const ImageSize = 105

// Sample is a single labeled character.
type Sample struct {
	// Alphabet is the name of the alphabet directory, such
	// as "Early_Aramaic" or "Futurama".
	Alphabet string

	// CharName is the name of the character directory, such
	// as "character37".
	CharName string

	// Path is the path to the image for this sample, such as
	// "/data/Early_Aramaic/character19/0269_02.png".
	Path string
}

// Image reads the sample's image file.
func (s *Sample) Image() (image.Image, error) {
	f, err := os.Open(s.Path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return png.Decode(f)
}

// AugSample is an augmented sample.
// In addition to the normal Sample, an AugSample includes
// a rotation.
// The rotation is a number between 0 and 3, inclusive.
// It is a counter-clockwise angle, measured in multiples
// of 90 degrees.
type AugSample struct {
	Sample   *Sample
	Rotation int
}

// Image reads the image file and rotates it accordingly.
func (a *AugSample) Image() (image.Image, error) {
	raw, err := a.Sample.Image()
	if err != nil {
		return nil, err
	}
	return rotate(raw, float64(a.Rotation)*math.Pi/2), nil
}

func (a *AugSample) rotated(rot int) *AugSample {
	s := *a
	s.Rotation = rot
	return &s
}

func rotate(img image.Image, angle float64) image.Image {
	input := make([]float64, 0, ImageSize*ImageSize)
	for y := 0; y < ImageSize; y++ {
		for x := 0; x < ImageSize; x++ {
			r, _, _, _ := img.At(img.Bounds().Min.X+x, img.Bounds().Min.Y+y).RGBA()
			input = append(input, float64(r)/0xffff)
		}
	}
	sin, cos := math.Sin(angle), math.Cos(angle)
	out := image.NewGray(image.Rect(0, 0, ImageSize, ImageSize))
	for y := 0; y < ImageSize; y++ {
		for x := 0; x < ImageSize; x++ {
			newX, newY := rotateCoord(x, y, sin, cos)
			newColor := int(interpolate(input, newX, newY) * 0x100)
			if newColor == 0x100 {
				newColor = 0xff
			}
			out.Set(x, y, color.Gray{Y: uint8(newColor)})
		}
	}
	return out
}

func rotateCoord(x, y int, sin, cos float64) (float64, float64) {
	cx := float64(x) - ImageSize/2
	cy := float64(y) - ImageSize/2
	return ImageSize/2 + cos*cx - sin*cy, ImageSize/2 + sin*cx + cos*cy
}

func interpolate(buf []float64, x, y float64) float64 {
	right := x - math.Floor(x)
	left := 1 - right
	bottom := y - math.Floor(y)
	top := 1 - bottom
	return left*top*pixelAt(buf, int(x), int(y)) +
		right*top*pixelAt(buf, int(x)+1, int(y)) +
		left*bottom*pixelAt(buf, int(x), int(y)+1) +
		right*bottom*pixelAt(buf, int(x)+1, int(y)+1)
}

func pixelAt(buf []float64, x, y int) float64 {
	if x < 0 || y < 0 || x >= ImageSize || y >= ImageSize {
		return 1
	}
	return buf[x+y*ImageSize]
}
