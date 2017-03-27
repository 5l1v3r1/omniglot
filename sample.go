package omniglot

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"
	"sync"
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

	imageLock sync.Mutex
	image     image.Image
}

// Image reads the sample's image file.
func (s *Sample) Image() (image.Image, error) {
	s.imageLock.Lock()
	defer s.imageLock.Unlock()
	if s.image != nil {
		return s.image, nil
	}
	f, err := os.Open(s.Path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	s.image, err = png.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("decode %s: %s", s.Path, err)
	} else if s.image.Bounds().Dx() != s.image.Bounds().Dy() {
		return nil, fmt.Errorf("decode %s: not square", s.Path)
	}
	return s.image, nil
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
//
// If augment is true, then extra data augmentation is
// applied.
// Either way, the image is rotated according to
// a.Rotation.
//
// If outSize is non-zero, then the sample is scaled to
// outSize by outSize pixels.
func (a *AugSample) Image(augment bool, outSize int) (image.Image, error) {
	if outSize == 0 {
		outSize = ImageSize
	}
	raw, err := a.Sample.Image()
	if err != nil {
		return nil, err
	}
	var transX, transY float64
	angle := float64(a.Rotation) * math.Pi / 2
	if augment {
		transX, transY = randomTranslation(), randomTranslation()
		angle += (rand.Float64() - 0.5) * (math.Pi / 8)
	}
	return transform(raw, outSize, angle, transX, transY), nil
}

func (a *AugSample) rotated(rot int) *AugSample {
	s := *a
	s.Rotation = rot
	return &s
}

func transform(img image.Image, outSize int, angle, transX, transY float64) image.Image {
	input := Tensor(img)
	sin, cos := math.Sin(angle), math.Cos(angle)
	scaler := (ImageSize - 1) / float64(outSize-1)
	out := image.NewGray(image.Rect(0, 0, outSize, outSize))
	for y := 0; y < outSize; y++ {
		for x := 0; x < outSize; x++ {
			newX, newY := rotateCoord(float64(x)*scaler, float64(y)*scaler, sin, cos)
			newX -= transX
			newY -= transY
			newColor := int(interpolate(input, newX, newY) * 0x100)
			if newColor == 0x100 {
				newColor = 0xff
			}
			out.Set(x, y, color.Gray{Y: uint8(newColor)})
		}
	}
	return out
}

func rotateCoord(x, y float64, sin, cos float64) (float64, float64) {
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

func randomTranslation() float64 {
	return rand.Float64()*20 - 10
}
