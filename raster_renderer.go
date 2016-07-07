package chart

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"

	"golang.org/x/image/font"

	"github.com/golang/freetype/truetype"
	drawing "github.com/llgcode/draw2d/draw2dimg"
)

// PNG returns a new png/raster renderer.
func PNG(width, height int) Renderer {
	i := image.NewRGBA(image.Rect(0, 0, width, height))
	return &rasterRenderer{
		i:  i,
		gc: drawing.NewGraphicContext(i),
	}
}

// RasterRenderer renders chart commands to a bitmap.
type rasterRenderer struct {
	i  *image.RGBA
	gc *drawing.GraphicContext
	fc *font.Drawer

	font      *truetype.Font
	fontColor color.RGBA
	fontSize  float64
}

// SetStrokeColor implements the interface method.
func (rr *rasterRenderer) SetStrokeColor(c color.RGBA) {
	println("RasterRenderer :: SetStrokeColor", ColorAsString(c))
	rr.gc.SetStrokeColor(c)
}

// SetFillColor implements the interface method.
func (rr *rasterRenderer) SetFillColor(c color.RGBA) {
	println("RasterRenderer :: SetFillColor", ColorAsString(c))
	rr.gc.SetFillColor(c)
}

// SetLineWidth implements the interface method.
func (rr *rasterRenderer) SetLineWidth(width float64) {
	println("RasterRenderer :: SetLineWidth", width)
	rr.gc.SetLineWidth(width)
}

// MoveTo implements the interface method.
func (rr *rasterRenderer) MoveTo(x, y int) {
	println("RasterRenderer :: MoveTo", x, y)
	rr.gc.MoveTo(float64(x), float64(y))
}

// LineTo implements the interface method.
func (rr *rasterRenderer) LineTo(x, y int) {
	println("RasterRenderer :: LineTo", x, y)
	rr.gc.LineTo(float64(x), float64(y))
}

// Close implements the interface method.
func (rr *rasterRenderer) Close() {
	println("RasterRenderer :: Close")
	rr.gc.Close()
}

// Stroke implements the interface method.
func (rr *rasterRenderer) Stroke() {
	println("RasterRenderer :: Stroke")
	rr.gc.Stroke()
}

// Fill implements the interface method.
func (rr *rasterRenderer) Fill() {
	println("RasterRenderer :: Fill")
	rr.gc.Fill()
}

// FillStroke implements the interface method.
func (rr *rasterRenderer) FillStroke() {
	println("RasterRenderer :: FillStroke")
	rr.gc.FillStroke()
}

// Circle implements the interface method.
func (rr *rasterRenderer) Circle(radius float64, x, y int) {
	xf := float64(x)
	yf := float64(y)
	rr.gc.MoveTo(xf-radius, yf)              //9
	rr.gc.QuadCurveTo(xf, yf, xf, yf-radius) //12
	rr.gc.QuadCurveTo(xf, yf, xf+radius, yf) //3
	rr.gc.QuadCurveTo(xf, yf, xf, yf+radius) //6
	rr.gc.QuadCurveTo(xf, yf, xf-radius, yf) //9
	rr.gc.Close()
	rr.gc.FillStroke()
}

// SetFont implements the interface method.
func (rr *rasterRenderer) SetFont(f *truetype.Font) {
	rr.font = f
	rr.gc.SetFont(f)
}

// SetFontSize implements the interface method.
func (rr *rasterRenderer) SetFontSize(size float64) {
	println("RasterRenderer :: SetFontSize", fmt.Sprintf("%.2f", size))
	rr.fontSize = size
	rr.gc.SetFontSize(size)
}

// SetFontColor implements the interface method.
func (rr *rasterRenderer) SetFontColor(c color.RGBA) {
	println("RasterRenderer :: SetFontColor", ColorAsString(c))
	rr.fontColor = c
	rr.gc.SetFillColor(c)
	rr.gc.SetStrokeColor(c)
}

// Text implements the interface method.
func (rr *rasterRenderer) Text(body string, x, y int) {
	println("RasterRenderer :: Text", body, x, y)
	rr.gc.CreateStringPath(body, float64(x), float64(y))
	rr.gc.Fill()
}

// MeasureText implements the interface method.
func (rr *rasterRenderer) MeasureText(body string) int {
	if rr.fc == nil && rr.font != nil {
		rr.fc = &font.Drawer{
			Face: truetype.NewFace(rr.font, &truetype.Options{
				DPI:  DefaultDPI,
				Size: rr.fontSize,
			}),
		}
	}
	if rr.fc != nil {
		dimensions := rr.fc.MeasureString(body)
		return dimensions.Floor()
	}
	return 0
}

// Save implements the interface method.
func (rr *rasterRenderer) Save(w io.Writer) error {
	println("RasterRenderer :: Save")
	return png.Encode(w, rr.i)
}