package debug

import (
	"golang.org/x/mobile/exp/gl/glutil"
	"fmt"
	"golang.org/x/mobile/event/size"
)

type Point struct {
	str *String
}

func NewPoint(images *glutil.Images) *Point {
	return &Point{
		str : NewString(images),
	}
}

func (p *Point) Draw(x,y float32,sz size.Event) {
	buf := fmt.Sprintf("X:%d,Y:%d",int(x),int(y))
	p.str.Draw(sz,buf)
}

func (p *Point) Release() {
	if p.str != nil {
		p.str.Release()
	}
}

