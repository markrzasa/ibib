package game

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Control interface {
	update(g *FirstGame)
	getPosition() (int, int)
	isButtonPressed() bool
}

type MouseControl struct {
	x, y int
	pressed bool
}

func (c *MouseControl) update(g *FirstGame) {
	c.x, c.y = ebiten.CursorPosition()
	c.pressed = ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)
}

func (c *MouseControl) getPosition() (int, int) {
	return c.x, c.y
}

func (c MouseControl) isButtonPressed() bool {
	return c.pressed
}

type PadControl struct {
	x, y int
	pressed bool
}

func (c *PadControl) update(g *FirstGame) {
	for id := range g.gamepadIds {
		horizontal := ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisLeftStickHorizontal)
		vertical := ebiten.StandardGamepadAxisValue(id, ebiten.StandardGamepadAxisLeftStickVertical)
		c.x = int(g.halfWidth + (g.halfWidth * horizontal))
		c.y = int(g.halfHeight + (g.halfHeight * vertical))
		log.Println(fmt.Sprintf("x: %d, y: %d", c.x, c.y))
		log.Println(fmt.Sprintf("h: %3.2f, v: %3.2f", horizontal, vertical))
		c.pressed = g.isPadButtonPressed()
	}
}

func (c *PadControl) getPosition() (int, int) {
	return c.x, c.y
}

func (c *PadControl) isButtonPressed() bool {
	return c.pressed
}
