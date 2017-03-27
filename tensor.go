package omniglot

import "image"

// Tensor converts an image to a gray-scale tensor.
func Tensor(img image.Image) []float64 {
	if gray, ok := img.(*image.Gray); ok {
		return grayTensor(gray)
	}
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

func grayTensor(img *image.Gray) []float64 {
	res := make([]float64, 0, img.Bounds().Dx()*img.Bounds().Dy())
	for y := 0; y < img.Bounds().Dy(); y++ {
		rowStart := y * img.Stride
		for x := 0; x < img.Bounds().Dx(); x++ {
			res = append(res, float64(img.Pix[rowStart+x])/0xff)
		}
	}
	return res
}
