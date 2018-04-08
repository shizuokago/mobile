package pazzle

import (
	"image"
	_ "image/png"
	"log"
	"math/rand"

	"fmt"
	"golang.org/x/mobile/asset"
	"golang.org/x/mobile/event/touch"
	"golang.org/x/mobile/exp/f32"
	"golang.org/x/mobile/exp/sprite"
	"golang.org/x/mobile/exp/sprite/clock"
)

const (
	texRed = iota
	texBlue
	texGreen
	texYellow
	texPurple
	texPink
	texBlack
)

const (
	stateNone = iota
	statePick
	stateMove
	stateComplate
)

func randomPiece() int {
	return rand.Intn(7)
}

type arrangerFunc func(e sprite.Engine, n *sprite.Node, t clock.Time)

func (a arrangerFunc) Arrange(e sprite.Engine, n *sprite.Node, t clock.Time) { a(e, n, t) }

type Game struct {
	zoom float32

	imageX float32
	imageY float32

	boardX float32
	boardY float32

	state int
	board Board

	pick   bool
	pickX  float32
	pickY  float32
	pieceX int
	pieceY int

	startTime   clock.Time
	releaseTime clock.Time
}

type Board struct {
	column int
	row    int
	data   [][]Piece
}

type Piece struct {
	datum int
}

func NewGame() *Game {
	var g Game
	g.reset()
	return &g
}

func (g *Game) reset() {

	g.imageX = 36
	g.imageY = 36

	g.boardX = 0
	g.boardY = 0
	g.zoom = 1
	g.state = stateNone

	board := Board{
		column: 6,
		row:    5,
	}

	board.data = make([][]Piece, board.column)
	for x, _ := range board.data {
		board.data[x] = make([]Piece, board.row)
		for y, _ := range board.data[x] {
			board.data[x][y] = Piece{
				datum: randomPiece(),
			}
		}
	}
	g.board = board
}

func (g *Game) Release() {
}

func (g *Game) Scene(eng sprite.Engine) *sprite.Node {

	texs := loadTextures(eng)
	if texs != nil {
	}

	scene := &sprite.Node{}
	eng.Register(scene)
	eng.SetTransform(scene, f32.Affine{
		{g.zoom, 0, g.boardX},
		{0, g.zoom, g.boardY},
	})

	newNode := func(fn arrangerFunc) {
		n := &sprite.Node{Arranger: arrangerFunc(fn)}
		eng.Register(n)
		scene.AppendChild(n)
	}

	//盤面の状態に応じたノードを作成
	for x, _ := range g.board.data {
		for y, _ := range g.board.data[x] {

			x := x
			y := y

			newNode(func(eng sprite.Engine, n *sprite.Node, t clock.Time) {

				startX := float32(x) * g.imageX
				startY := float32(y) * g.imageY

				if g.state == statePick || g.state == stateMove {
					if x == g.pieceX && y == g.pieceY {

						//座標はオフセットや倍率がかかっているので算出
						startX = g.pickX
						startY = g.pickY

						//startX = (g.pickX - g.boardX) / g.zoom
						//startY = (g.pickY - g.boardY) / g.zoom
						//中央分引き込む
						//startX -= g.imageX / 2
						//startY -= g.imageY / 2
					}
				}

				a := f32.Affine{
					{g.imageX, 0, startX},
					{0, g.imageY, startY},
				}
				eng.SetSubTex(n, texs[g.board.data[x][y].datum])
				eng.SetTransform(n, a)
			})
		}
	}

	return scene
}

/**
 *  assetからのtextureの読み込み
 */
func loadTextures(eng sprite.Engine) []sprite.SubTex {

	a, err := asset.Open("pazzle.png")
	if err != nil {
		log.Fatal(err)
	}
	defer a.Close()
	m, _, err := image.Decode(a)
	if err != nil {
		log.Fatal(err)
	}
	t, err := eng.LoadTexture(m)
	if err != nil {
		log.Fatal(err)
	}

	const n = 36

	rtn := make([]sprite.SubTex, texBlack+1)
	//パズルのピースを切り取る
	for idx := texRed; idx <= texBlack; idx++ {
		rect := image.Rect(n*idx+1, 0, n*(idx+1)-1, n)
		rtn[idx] = sprite.SubTex{
			T: t,
			R: rect,
		}
	}

	return rtn
}

func (g *Game) Touch(e touch.Event) bool {
	touchType := e.Type
	if touchType == touch.TypeBegin {
		//座標からピックしているピースを特定
		//return g.pickPiece(e.X, e.Y)
		g.state = statePick
		g.pieceX = 0
		g.pieceY = 0

		return true
	} else if touchType == touch.TypeEnd {
		//動かしていた場合
		if g.state == stateMove {
			//終了状態
			g.state = stateComplate
		} else {
			//通常状態
			g.state = stateNone
		}
		return true
	} else if g.state == statePick || g.state == stateMove {
		//移動処理
		return g.move(e.X, e.Y)
	}
	return false
}

func (g *Game) pickPiece(x, y float32) bool {

	//盤面の領域内にいるかを判定
	bx, by, err := g.getPiece(x, y)
	if err != nil {
		return false
	}

	//状態をpick状態にする
	g.state = statePick

	g.pieceX = bx
	g.pieceY = by
	g.pickX = x
	g.pickY = y

	return true
}

func (g *Game) move(x, y float32) bool {

	g.pickX = x
	g.pickY = y

	//盤面の領域内にいるかを判定
	newPieceX, newPieceY, err := g.getPiece(x, y)
	if err != nil {
		return true
	}

	if newPieceX == g.pieceX && newPieceY == g.pieceY {
		return true
	}

	g.state = stateMove

	//盤面の位置を変更
	g.board.data[g.pieceX][g.pieceY], g.board.data[newPieceX][newPieceY] =
		g.board.data[newPieceX][newPieceY], g.board.data[g.pieceX][g.pieceY]

	//ピック対象も変更
	g.pieceX = newPieceX
	g.pieceY = newPieceY

	return true
}

func (g *Game) getPiece(x, y float32) (int, int, error) {
	//盤面の領域内にいるかを判定
	minX := g.boardX
	minY := g.boardY
	maxX := minX + (g.imageX * float32(g.board.column) * g.zoom)
	maxY := minY + (g.imageY * float32(g.board.row) * g.zoom)

	if x < minX || x > maxX &&
		y < minY || y > maxY {
		return -1, -1, fmt.Errorf("out of board")
	}

	//ピースの座標X,Yを算出
	for bx, _ := range g.board.data {
		startX := minX + (g.imageX * float32(bx) * g.zoom)
		endX := startX + (g.imageX * g.zoom)
		if x < startX || x > endX {
			continue
		}

		for by, _ := range g.board.data[bx] {
			startY := minY + (g.imageY * float32(by) * g.zoom)
			endY := startY + (g.imageY * g.zoom)
			if y < startY || y > endY {
				continue
			}

			return bx, by, nil
		}
	}
	return -1, -1, fmt.Errorf("out of board")
}
