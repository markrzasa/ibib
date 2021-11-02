package game

import (
	"bytes"
	_ "embed"
	"image/color"
	"image/png"
	"log"
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

//go:embed sprites/balloon.png
var balloon []byte

//go:embed sprites/balloon_popped.png
var balloonPopped []byte

//go:embed sprites/bullet.png
var bullet []byte

type FirstGame struct {
	state gameState

	width, height, bulletCount, bulletRate int

	balloon ebiten.Image
	bullet ebiten.Image
	cloud ebiten.Image
	crosshair ebiten.Image
	poppedBaloon ebiten.Image

	cursorX, cursorY int

	balloons []Balloon

	bullets Bullets

	clouds []Cloud

	shooting bool

	gamepadIdsBuffer []ebiten.GamepadID
	gamepadIds map[ebiten.GamepadID]bool
}

func (g *FirstGame) manageGamepads() {
	g.gamepadIdsBuffer = inpututil.AppendJustConnectedGamepadIDs(g.gamepadIdsBuffer[:0])
	for _, id := range g.gamepadIdsBuffer {
		g.gamepadIds[id] = true
	}
	for id, _ := range g.gamepadIds {
		if inpututil.IsGamepadJustDisconnected(id) {
			delete(g.gamepadIds, id)
		}
	}
}

func (g *FirstGame) isButtonPressed(id ebiten.GamepadID) bool {
	pressed := false
	maxButtons := ebiten.GamepadButton(ebiten.GamepadButtonNum(id))
	for b := ebiten.GamepadButton(id); b < maxButtons; b++ {
		if ebiten.IsGamepadButtonPressed(id, b) {
			pressed = true
			break
		}
	}

	return pressed
}

func (g *FirstGame) isShooting() bool {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		return true
	}

	shooting := false
	for id, _ := range(g.gamepadIds) {
		if g.isButtonPressed(id) {
			shooting = true
			break
		}
	}	

	return shooting
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
	g.manageGamepads()
	g.cursorX, g.cursorY = ebiten.CursorPosition()
	g.shooting = g.isShooting()

	switch g.state {
	case Intro:
		if g.shooting {
			g.state = Running
		}
	case Running:
		g.updateClouds()
		g.updateBalloons()
		g.bullets.Update(g)
	}

	return nil
}

func (g *FirstGame) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x87, 0xCE, 0xEB, 0xff})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(g.cursorX - (g.crosshair.Bounds().Dx() / 2)), float64(g.cursorY - (g.crosshair.Bounds().Dy() / 2)))
	screen.DrawImage(&g.crosshair, op)

	g.bullets.Draw(screen, g)
	switch g.state {
	case Intro:
		text.Draw(
			screen,
			"Move the crosshair with the mouse. Shoot with the left mouse button.",
			basicfont.Face7x13,
			10, 20, color.RGBA{0x00, 0x00, 0x00, 0xff})
	case Running:
		for _, cloud := range g.clouds {
			cloud.Draw(screen, g)
		}

		for _, balloon := range g.balloons {
			balloon.Draw(screen, g)
		}
	}
}

func (g *FirstGame) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	if outsideWidth != g.width || outsideHeight != g.height {
		log.Println("layout changed - reinitialize")
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
	g.cloud = *g.newImage(cloud)
	g.crosshair = *g.newImage(crosshair)
	g.balloon = *g.newImage(balloon)
	g.poppedBaloon = *g.newImage(balloonPopped)
	g.bullet = *g.newImage(bullet)
}

func (g *FirstGame) initialize(width, height int) {
	g.state = Intro
	g.width = width
	g.height = height
	g.bulletCount = 0
	g.bulletRate = 20

	if g.gamepadIds == nil {
		g.gamepadIds = map[ebiten.GamepadID]bool{}
	}

	balloonCount := int(width / g.balloon.Bounds().Dx())
	g.balloons = make([]Balloon, balloonCount)
	for i := 0; i < balloonCount; i++ {
		g.balloons[i].x = i * g.balloon.Bounds().Dx()
		g.balloons[i].y = 0
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
