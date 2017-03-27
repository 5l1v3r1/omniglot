package omniglot

import (
	"image"
	"testing"
)

func BenchmarkTransform(b *testing.B) {
	img := image.NewGray(image.Rect(0, 0, ImageSize, ImageSize))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		transform(img, 20, 0.5, -2, 2)
	}
}
