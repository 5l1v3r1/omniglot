package omniglot

import (
	"image"
	"testing"
)

func BenchmarkTensor(b *testing.B) {
	img := image.NewGray(image.Rect(0, 0, ImageSize, ImageSize))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Tensor(img)
	}
}
