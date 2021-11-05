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
	Falling
)

type Balloon struct {
	x, y int

	popCount int

	state ballonState

	image *ebiten.Image
}

func (b *Balloon) isCollision(g *FirstGame) bool {
	bounds := image.Rect(
		b.image.Bounds().Min.X, b.image.Bounds().Min.Y,
		b.image.Bounds().Max.X / 3, b.image.Bounds().Max.Y)
	bounds = bounds.Add(image.Point{b.x, g.height}).Sub(image.Point{0, b.y})
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
			b.state = Falling
			b.popCount = 0
			b.y = b.y - b.image.Bounds().Dy()
		}
	case Falling:
		b.y -= 1
		if b.y <= 0 - b.image.Bounds().Dy() {
			b.state = Floating
			b.y = 0
		}
	}
}

func (b *Balloon) Draw(screen *ebiten.Image, g *FirstGame) {
	switch b.state {
	case Popped:
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(g.cursorX - (g.poppedBaloon.Bounds().Dx() / 2)), float64(g.cursorY - (g.poppedBaloon.Bounds().Dy() / 2)))
		screen.DrawImage(g.poppedBaloon, op)
		fallthrough
	case Floating:
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(b.x), float64(screen.Bounds().Dy() - b.y))
		width := b.image.Bounds().Dx() / 3
		subImageRect := image.Rect(0, 0, width, b.image.Bounds().Dy())
		subImage := b.image.SubImage(subImageRect).(*ebiten.Image)
		screen.DrawImage(subImage, op)
	case Falling:
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(b.x), float64(screen.Bounds().Dy() - b.y))
		width := b.image.Bounds().Dx() / 3
		offset := 1
		if (b.y % 20) / 10 < 1 {
			offset = 0
		}
		x1 := width * (1 + (offset))
		x2 := x1 + width
		subImageRect := image.Rect(x1, 0, x2, b.image.Bounds().Dy())
		subImage := b.image.SubImage(subImageRect).(*ebiten.Image)
		screen.DrawImage(subImage, op)
	}
}
