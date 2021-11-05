package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Cloud struct {
	x, y int
}

func (c *Cloud) Update(g *FirstGame) {
	c.x = c.x + 1
	if c.x == g.width {
		c.x = 0 - g.cloud.Bounds().Dx()
	}
	c.x = c.x % g.width
}

func (c *Cloud) Draw(screen *ebiten.Image, g *FirstGame) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(c.x), float64(c.y))
	screen.DrawImage(g.cloud, op)
}
