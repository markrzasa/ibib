package game

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Bullet struct {
	startX int

	m, b float32

	end image.Point
}

type Bullets struct {
	leftBullet Bullet
	rightBullet Bullet
}

const (
	bulletCount = 50
)


func (bullets *Bullets) Update(g *FirstGame) {
	if g.shooting {
		bullets.leftBullet.startX = (bullets.leftBullet.startX + 1) % 10
		bullets.leftBullet.end = image.Point{g.cursorX, g.cursorY}
		bullets.leftBullet.m = float32(bullets.leftBullet.end.Y - g.height) / float32(bullets.leftBullet.end.X - 0)
		bullets.leftBullet.b = float32(g.height) - float32(bullets.leftBullet.m * float32(0))
		bullets.rightBullet.startX = (bullets.rightBullet.startX + 1) % 10
		bullets.rightBullet.end = image.Point{g.cursorX, g.cursorY}
		bullets.rightBullet.m = float32(bullets.rightBullet.end.Y - g.height) / float32(bullets.rightBullet.end.X - g.width)
		bullets.rightBullet.b = float32(g.height) - float32(bullets.rightBullet.m * float32(g.width))
	} else {
		bullets.leftBullet.startX = 0
		bullets.rightBullet.startX = 0
	}
}

func (bullets *Bullets) Draw(screen *ebiten.Image, g *FirstGame) {
	if g.shooting {
		offset := image.Point{g.bullet.Bounds().Dx() / 2, g.bullet.Bounds().Dx() / 2}
		xInc := int(math.Max(float64(bullets.leftBullet.end.X / bulletCount), 1))
		for x := bullets.leftBullet.startX; x < bullets.leftBullet.end.X; x += xInc {
			op := &ebiten.DrawImageOptions{}
			bulletY := int((bullets.leftBullet.m * float32(x)) + float32(bullets.leftBullet.b))
	
			op.GeoM.Translate(float64(x - offset.X), float64(bulletY - offset.Y))
			screen.DrawImage(&g.bullet, op)
		}
	
		xInc = int(math.Max(float64((g.width - bullets.rightBullet.end.X) / bulletCount), 1))
		for x := g.width - bullets.rightBullet.startX; x > bullets.rightBullet.end.X; x -= xInc {
			op := &ebiten.DrawImageOptions{}
			bulletY := int((bullets.rightBullet.m * float32(x)) + float32(bullets.rightBullet.b))
	
			op.GeoM.Translate(float64(x - offset.X), float64(bulletY - offset.Y))
			screen.DrawImage(&g.bullet, op)
		}	
	}
}
