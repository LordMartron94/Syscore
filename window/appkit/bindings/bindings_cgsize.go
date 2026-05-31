//go:build darwin

package bindings

/*
CGSize is the CoreGraphics size layout (two CGFloat values) returned by CAMetalLayer.drawableSize.
*/
type CGSize struct {
	Width  float64
	Height float64
}
