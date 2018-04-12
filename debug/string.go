package debug

import (
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/geom"
	"image"
	"image/color"
	"image/draw"
)

type String struct {
	images *glutil.Images
	m      *glutil.Image
	datum  string
}

const (
	fontWidth  = 5
	fontHeight = 7
	length     = 15
)

func NewString(images *glutil.Images) *String {
	return &String{
		images: images,
	}
}

func (s *String) Draw(sz size.Event, buf string) {

	const imgW, imgH = length*(fontWidth+1) + 1, fontHeight + 2
	if s.m != nil {
		s.m.Release()
	}

	s.m = s.images.NewImage(imgW, imgH)

	draw.Draw(s.m.RGBA, s.m.RGBA.Bounds(), image.White, image.Point{}, draw.Src)
	for i, c := range buf {

		glyph := glyphs[c]
		if len(glyph) != fontWidth*fontHeight {
			continue
		}
		for y := 0; y < fontHeight; y++ {
			for x := 0; x < fontWidth; x++ {
				if glyph[fontWidth*y+x] == ' ' {
					continue
				}
				s.m.RGBA.SetRGBA((fontWidth+1)*i+x+1, y+1, color.RGBA{A: 0xff})
			}
		}
	}

	s.m.Upload()

	topLeft := geom.Point{X: 0, Y: sz.HeightPt - imgH}
	topRight := geom.Point{X: imgW, Y: sz.HeightPt - imgH}
	bottomLeft := geom.Point{X: 0, Y: sz.HeightPt}

	s.m.Draw(sz, topLeft, topRight, bottomLeft, s.m.RGBA.Bounds())
}

func (s *String) Release() {
	if s.m != nil {
		s.m.Release()
		s.m = nil
		s.images = nil
	}
}

// I implemented it with reference
// https://github.com/golang/mobile/blob/master/exp/app/debug/fps.go

var glyphs = [256]string{
	'0': "" +
		"  X  " +
		" X X " +
		"X   X" +
		"X   X" +
		"X   X" +
		" X X " +
		"  X  ",
	'1': "" +
		"  X  " +
		" XX  " +
		"X X  " +
		"  X  " +
		"  X  " +
		"  X  " +
		"XXXXX",
	'2': "" +
		" XXX " +
		"X   X" +
		"    X" +
		"  XX " +
		" X   " +
		"X    " +
		"XXXXX",
	'3': "" +
		"XXXXX" +
		"    X" +
		"   X " +
		"  XX " +
		"    X" +
		"X   X" +
		" XXX ",
	'4': "" +
		"   X " +
		"  XX " +
		" X X " +
		"X  X " +
		"XXXXX" +
		"   X " +
		"   X ",
	'5': "" +
		"XXXXX" +
		"X    " +
		"X XX " +
		"XX  X" +
		"    X" +
		"X   X" +
		" XXX ",
	'6': "" +
		"  XX " +
		" X   " +
		"X    " +
		"X XX " +
		"XX  X" +
		"X   X" +
		" XXX ",
	'7': "" +
		"XXXXX" +
		"    X" +
		"   X " +
		"   X " +
		"  X  " +
		" X   " +
		" X   ",
	'8': "" +
		" XXX " +
		"X   X" +
		"X   X" +
		" XXX " +
		"X   X" +
		"X   X" +
		" XXX ",
	'9': "" +
		" XXX " +
		"X   X" +
		"X  XX" +
		" XX X" +
		"    X" +
		"   X " +
		" XX  ",
	'F': "" +
		"XXXXX" +
		"X    " +
		"X    " +
		"XXXX " +
		"X    " +
		"X    " +
		"X    ",
	'P': "" +
		"XXXX " +
		"X   X" +
		"X   X" +
		"XXXX " +
		"X    " +
		"X    " +
		"X    ",
	'S': "" +
		" XXX " +
		"X   X" +
		"X    " +
		" XXX " +
		"    X" +
		"X   X" +
		" XXX ",
	' ': "" +
		"     " +
		"     " +
		"     " +
		"     " +
		"     " +
		"     " +
		"     ",
	'X': "" +
		"X   X" +
		"X   X" +
		" X X " +
		"  X  " +
		" X X " +
		"X   X" +
		"X   X",
	'Y': "" +
		"X   X" +
		"X   X" +
		" X X " +
		"  X  " +
		"  X  " +
		"  X  " +
		"  X  ",
	':': "" +
		"     " +
		"  X  " +
		"  X  " +
		"     " +
		"  X  " +
		"  X  " +
		"     ",
	',': "" +
		"     " +
		"     " +
		"     " +
		"     " +
		"  XX " +
		"  XX " +
		"   x ",
}
