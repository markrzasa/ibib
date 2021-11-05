package game

import (
	"bytes"
	_ "embed"
	"image/color"
	"image/png"
	"log"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"

	"golang.org/x/image/font/basicfont"
)

type gameState int
const (
	Intro gameState = iota
	Running
)

//go:embed sprites/cloud.png
var cloud []byte

//go:embed sprites/crosshair.png
var crosshair []byte

//go:embed sprites/blue_balloon.png
var blueBalloon []byte

//go:embed sprites/red_balloon.png
var redBalloon []byte

//go:embed sprites/yellow_balloon.png
var yellowBalloon []byte

//go:embed sprites/balloon_popped.png
var balloonPopped []byte

//go:embed sprites/bullet.png
var bullet []byte

type FirstGame struct {
	state gameState

	width, height, bulletCount, bulletRate int

	halfWidth, halfHeight float64

	control Control

	balloonImages [3]*ebiten.Image
	bullet *ebiten.Image
	cloud *ebiten.Image
	crosshair *ebiten.Image
	poppedBaloon *ebiten.Image

	cursorX, cursorY int

	balloons []Balloon

	bullets Bullets

	clouds []Cloud

	shooting bool

	gamepadIdsBuffer []ebiten.GamepadID
	gamepadIds map[ebiten.GamepadID]bool
}

func (g *FirstGame) isPadButtonPressed() bool {
	pressed := false
	for id := range(g.gamepadIds) {
		for b := ebiten.StandardGamepadButtonRightBottom; b < ebiten.StandardGamepadButtonMax; b++ {
			if ebiten.IsStandardGamepadButtonPressed(id, b) {
				pressed = true
				break
			}
		}
	}

	return pressed
}

func (g *FirstGame) updateGamepads() {
	g.gamepadIdsBuffer = inpututil.AppendJustConnectedGamepadIDs(g.gamepadIdsBuffer[:0])
	for _, id := range g.gamepadIdsBuffer {
		g.gamepadIds[id] = true
	}
	for id := range g.gamepadIds {
		if inpututil.IsGamepadJustDisconnected(id) {
			delete(g.gamepadIds, id)
		}
	}
}

func (g *FirstGame) updateCrosshair() {
	g.shooting = g.control.isButtonPressed()

	g.cursorX, g.cursorY = g.control.getPosition()	
}

func (g *FirstGame) updateClouds() {
	for i := 0; i < len(g.clouds); i++ {
		g.clouds[i].Update(g)
	}
}

func (g *FirstGame) updateBalloons() {
	for i := 0; i < len(g.balloons); i++ {
		g.balloons[i].Update(g)
	}
}

func (g *FirstGame) Update() error {
	g.updateGamepads()

	switch g.state {
	case Intro:
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			g.state = Running
			g.control = &MouseControl{}
			g.control.update(g)
		} else if g.isPadButtonPressed() {
			g.state = Running
			g.control = &PadControl{}
			g.control.update(g)
		}
	case Running:
		g.control.update(g)
		g.updateClouds()
		g.updateBalloons()
		g.updateCrosshair()
	}

	return nil
}

func (g *FirstGame) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x87, 0xCE, 0xEB, 0xff})
	switch g.state {
	case Intro:
		padsConnected := len(g.gamepadIds) > 0
		y := 20
		text.Draw(
			screen,
			"Mouse: Move the crosshair with the mouse. Shoot with the left mouse button.",
			basicfont.Face7x13,
			10, y, color.RGBA{0x00, 0x00, 0x00, 0xff})
		y += 20
		if padsConnected {
			text.Draw(
				screen,
				"Pad: Move the crosshair with the analog stick. Shoot with the bottom right button.",
				basicfont.Face7x13,
				10, y, color.RGBA{0x00, 0x00, 0x00, 0xff})
				text.Draw(
					screen,
					"Click the mouse button to play using the mouse. Click a pad button to play using a pad.",
					basicfont.Face7x13,
					10, y + 20, color.RGBA{0x00, 0x00, 0x00, 0xff})
		} else {
			text.Draw(
				screen,
				"Click the mouse button to play.",
				basicfont.Face7x13,
				10, y, color.RGBA{0x00, 0x00, 0x00, 0xff})
		}
	case Running:
		g.cursorX, g.cursorY = g.control.getPosition()
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(g.cursorX - (g.crosshair.Bounds().Dx() / 2)), float64(g.cursorY - (g.crosshair.Bounds().Dy() / 2)))
		screen.DrawImage(g.crosshair, op)
	
		g.bullets.Draw(screen, g)
			for _, cloud := range g.clouds {
			cloud.Draw(screen, g)
		}

		for _, balloon := range g.balloons {
			balloon.Draw(screen, g)
		}

		g.bullets.Update(g)
	}
}

func (g *FirstGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	if outsideWidth != g.width || outsideHeight != g.height {
		g.initialize(outsideWidth, outsideHeight)
	}
	return outsideWidth, outsideHeight
}

func (g *FirstGame) newImage(imageBytes []byte) *ebiten.Image {
	image, err := png.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		log.Fatal(err)
	}
	return ebiten.NewImageFromImage(image)
}

func (g *FirstGame) initializeImages() {
	g.cloud = g.newImage(cloud)
	g.crosshair = g.newImage(crosshair)
	g.balloonImages[0] = g.newImage(blueBalloon)
	g.balloonImages[1] = g.newImage(redBalloon)
	g.balloonImages[2] = g.newImage(yellowBalloon)
	g.poppedBaloon = g.newImage(balloonPopped)
	g.bullet = g.newImage(bullet)
}

func (g *FirstGame) initialize(width, height int) {
	g.state = Intro
	g.width = width
	g.height = height
	g.halfWidth = math.Ceil(float64(g.width) / 2.0)
	g.halfHeight = math.Ceil(float64(g.height) / 2.0)
	g.bulletCount = 0
	g.bulletRate = 20

	if g.gamepadIds == nil {
		g.gamepadIds = map[ebiten.GamepadID]bool{}
	}

	balloonWidth := g.balloonImages[0].Bounds().Dx() / 3
	balloonCount := int(width / balloonWidth)
	g.balloons = make([]Balloon, balloonCount)
	for i := 0; i < balloonCount; i++ {
		g.balloons[i].x = i * balloonWidth
		g.balloons[i].y = 0
		g.balloons[i].image = g.balloonImages[i % len(g.balloonImages)]
		g.balloons[i].state = StartWait
	}

	cloudCount := int(height / (g.cloud.Bounds().Dy() + 30))
	g.clouds = make([]Cloud, cloudCount)
	y := 30
	for i := 0; i < cloudCount; i++ {
		g.clouds[i].x = 0 - rand.Intn(g.width)
		g.clouds[i].y = y
		y = y + g.cloud.Bounds().Dy() + 30
	}
}

func RunGame(width, height int) {
	game := &FirstGame{}
	game.initializeImages()
	game.initialize(width, height)
	ebiten.SetWindowSize(width, height)
	ebiten.SetWindowTitle("Infinite Balloons/Infinite Bullets")
	ebiten.SetWindowResizable(true)
	ebiten.SetScreenTransparent(true)
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}	
}
