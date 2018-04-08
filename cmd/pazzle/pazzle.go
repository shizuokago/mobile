package pazzle

import (
	"time"
	"math/rand"
	"golang.org/x/mobile/app"
	"golang.org/x/mobile/gl"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/paint"
	//"golang.org/x/mobile/event/touch"
	//"golang.org/x/mobile/event/key"

	"golang.org/x/mobile/event/touch"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	app.Main(func(a app.App) {

		var ctx gl.Context
		var sz size.Event

		for e := range a.Events() {
			switch e := a.Filter(e).(type) {
			case lifecycle.Event:
				switch e.Crosses(lifecycle.StageVisible) {
				case lifecycle.CrossOn:
					ctx, _ = e.DrawContext.(gl.Context)
					onStart(ctx)
					a.Send(paint.Event{})
				case lifecycle.CrossOff:
					onStop()
					ctx = nil
				}
			case size.Event:
				sz = e
			case paint.Event:
				if ctx == nil || e.External {
					continue
				}
				onPaint(ctx, sz)
				a.Publish()
			case touch.Event:
				if onTouch(e) {
					a.Send(paint.Event{})
				}
			}
		}

	})
}
