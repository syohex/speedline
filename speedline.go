package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"

	"github.com/gographics/imagick/imagick"
)

var (
	output = flag.String("o", "output.gif", "output filename")
	color  = flag.String("color", "black", "background color")
)

func speedLine(mw *imagick.MagickWand, aw *imagick.MagickWand) error {
	rows, cols := mw.GetImageHeight(), mw.GetImageWidth()

	dw := imagick.NewDrawingWand()
	defer dw.Destroy()
	cw := imagick.NewPixelWand()
	cw.SetColor(*color)
	dw.SetFillColor(cw)

	center := []float64{float64(cols) / 2.0, float64(rows) / 2.0}
	const radiusCenter float64 = 0.75
	const step float64 = 0.02
	const bold float64 = 1.0
	var theeta float64

	for theeta < math.Pi*2 {
		stepNoise := rand.Float64() + 0.5
		theeta += step * stepNoise
		radiusCenterNoise := rand.Float64()*0.3 + 1.0
		boldNoise := rand.Float64() + 0.7 + 0.3

		point0 := imagick.PointInfo{
			X: math.Sin(theeta)*center[0]*radiusCenter*radiusCenterNoise + center[0],
			Y: math.Cos(theeta)*center[1]*radiusCenter*radiusCenterNoise + center[1],
		}
		point1 := imagick.PointInfo{
			X: math.Sin(theeta)*center[0]*2 + center[0],
			Y: math.Cos(theeta)*center[1]*2 + center[1],
		}
		point2 := imagick.PointInfo{
			X: math.Sin(theeta+step*bold*boldNoise)*center[0]*2 + center[0],
			Y: math.Cos(theeta+step*bold*boldNoise)*center[1]*2 + center[1],
		}

		dw.Polygon([]imagick.PointInfo{point0, point1, point2})
	}

	if err := aw.DrawImage(dw); err != nil {
		return err
	}

	return nil
}

func main() {
	flag.Parse()

	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()
	aw := imagick.NewMagickWand()
	defer aw.Destroy()

	if err := mw.ReadImage(flag.Arg(0)); err != nil {
		panic(err)
	}

	if int(mw.GetNumberImages()) == 1 {
		mw.SetIteratorIndex(0)
		first := mw.GetImage()
		mw.ResetIterator()
		for i := 0; i < 3; i++ {
			mw.AddImage(first.Clone())
		}
	}

	for i := 0; i < int(mw.GetNumberImages()); i++ {
		mw.SetIteratorIndex(i)
		tw := mw.GetImage()
		aw.AddImage(tw)
		if err := speedLine(tw, aw); err != nil {
			fmt.Println(err)
			return
		}
		tw.Destroy()
	}
	mw.ResetIterator()

	aw.SetOption("loop", "0")
	if err := aw.WriteImages(*output, true); err != nil {
		panic(err)
	}
}
