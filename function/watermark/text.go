package watermark

import (
	"github.com/disintegration/imaging"
	"github.com/fogleman/gg"
	"image"
)

func TextWatermarkImage(
	src image.Image,
	text, fontPath string,
	fontSize float64,
	angle float64,
	r, g, b, a float64,
	xGap, yGap float64) (image.Image, error) {

	w := src.Bounds().Dx()
	h := src.Bounds().Dy()
	dc := gg.NewContext(w, h)
	dc.DrawImage(src, 0, 0)

	if err := dc.LoadFontFace(fontPath, fontSize); err != nil {
		return nil, err
	}
	dc.SetRGBA(r, g, b, a)

	cx, cy := float64(w)/2, float64(h)/2
	dc.RotateAbout(angle, cx, cy)

	tw, _ := dc.MeasureString(text)
	xStep := tw + xGap
	yStep := fontSize + yGap

	startX := -float64(h)
	endX := float64(w) + float64(h)
	startY := -float64(h)
	endY := float64(h) + float64(w)

	for x := startX; x < endX; x += xStep {
		for y := startY; y < endY; y += yStep {
			dc.DrawString(text, x, y)
		}
	}

	return dc.Image(), nil
}

func TextWatermark(
	input string,
	text, fontPath string,
	fontSize float64,
	angle float64,
	r, g, b, a float64,
	xGap, yGap float64,
) (image.Image, error) {
	src, err := imaging.Open(input)
	if err != nil {
		return nil, err
	}
	return TextWatermarkImage(src, text, fontPath, fontSize, angle, r, g, b, a, xGap, yGap)
}
