package pazzle

import (
	"time"

	"github.com/shizuokago/mobile/pazzle"

	"golang.org/x/mobile/exp/gl/glutil"
	"golang.org/x/mobile/exp/sprite"
	"golang.org/x/mobile/gl"
	"golang.org/x/mobile/event/size"
	"golang.org/x/mobile/exp/sprite/clock"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/exp/app/debug"
	"golang.org/x/mobile/exp/sprite/glsprite"
)

var (
	startTime = time.Now()
	images    *glutil.Images
	eng       sprite.Engine
	scene     *sprite.Node
	game      *pazzle.Game
	fps       *debug.FPS
)

func onStart(ctx gl.Context) {
	images = glutil.NewImages(ctx)
	fps = debug.NewFPS(images)

	eng = glsprite.Engine(images)
	game = pazzle.NewGame()
	scene = game.Scene(eng)
}

func onStop() {
	eng.Release()
	images.Release()
	game.Release()
	fps.Release()
	game = nil
}

func onPaint(glctx gl.Context, sz size.Event) {
	glctx.ClearColor(0.7, 0.7, 0.7, 0.5)
	glctx.Clear(gl.COLOR_BUFFER_BIT)
	now := clock.Time(time.Since(startTime) * 60 / time.Second)
	fps.Draw(sz)
	eng.Render(scene, now, sz)
}

func onTouch(e touch.Event) bool {
	return game.Touch(e)
}

