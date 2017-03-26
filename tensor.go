package omniglot

import "image"

// Tensor converts an image to a gray-scale tensor.
func Tensor(img image.Image) []float64 {
	size := img.Bounds().Dx()
	input := make([]float64, 0, size*size)
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			r, _, _, _ := img.At(img.Bounds().Min.X+x, img.Bounds().Min.Y+y).RGBA()
			input = append(input, float64(r)/0xffff)
		}
	}
	return input
}
