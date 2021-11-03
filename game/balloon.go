package game

import (
	"image"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

type ballonState int
const (
	StartWait ballonState = iota
	Floating
	Popped
)

type Balloon struct {
	x, y int

	popCount int

	state ballonState

	image *ebiten.Image
}

func (b *Balloon) isCollision(g *FirstGame) bool {
	bounds := b.image.Bounds().Add(image.Point{b.x, g.height}).Sub(image.Point{0, b.y})
	cursor := image.Rectangle{image.Point{g.cursorX, g.cursorY}, image.Point{g.cursorX + 1, g.cursorY + 1}}
	return bounds.Intersect(cursor) != image.Rectangle{}
}

func (b *Balloon) Update(g *FirstGame) {
	switch b.state {
	case StartWait:
		if rand.Intn(100) % 11 == 0 {
			b.state = Floating
		}
	case Floating:
		if g.shooting && b.isCollision(g) {
			b.state = Popped
		} else {
			b.y = (b.y + 1) % (g.height + b.image.Bounds().Dy())
		}
	case Popped:
		if b.popCount < 10 {
			b.popCount++
		} else {
			b.state = StartWait
			b.popCount = 0
			b.y = 0
		}
	}
}

func (b *Balloon) Draw(screen *ebiten.Image, g *FirstGame) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(b.x), float64(screen.Bounds().Dy() - b.y))
	screen.DrawImage(b.image, op)
	if b.state == Popped {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(g.cursorX - (g.poppedBaloon.Bounds().Dx() / 2)), float64(g.cursorY - (g.poppedBaloon.Bounds().Dy() / 2)))
		screen.DrawImage(&g.poppedBaloon, op)
	}
}
