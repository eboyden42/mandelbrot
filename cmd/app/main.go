package main

import (
	"log"
	"math"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 500
	screenHeight = 500
	iterations   = 140
)

var (
	palette [iterations]byte
)

type Game struct {
	screen                 *ebiten.Image
	pixels                 []byte
	xCenter, yCenter, size float64
}

type PointCalculation struct {
	i, j                   int
	centerX, centerY, size float64
	pixels                 *[]byte
}

func init() {
	for i := range palette {
		palette[i] = byte(math.Sqrt(float64(i)/float64(len(palette))) * 0x80)
	}
}

func color(it int) (r, g, b byte) {
	if it == iterations {
		return 0xff, 0xff, 0xff
	}
	c := palette[it]
	return c, c, c
}

func NewGame() *Game {
	img := ebiten.NewImage(screenWidth, screenHeight)
	game := Game{screen: img, pixels: make([]byte, screenWidth*screenHeight*4), xCenter: 0, yCenter: 0, size: 4}
	return &game
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
		g.yCenter += 0.05 * g.size
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
		g.yCenter -= 0.05 * g.size
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		g.xCenter += 0.05 * g.size
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		g.xCenter -= 0.05 * g.size
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		g.size *= 9.0 / 10
	} else if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		g.size *= 10.0 / 9.0
	}

	g.calculateFractal(g.xCenter, g.yCenter, g.size)
	return nil
}

func (gm *Game) calculateFractal(centerX, centerY, size float64) {
	var wg sync.WaitGroup
	for j := range screenHeight {
		for i := range screenHeight {
			wg.Add(1)
			go calculatePoint(i, j, centerX, centerY, size, &gm.pixels, &wg)
		}
	}
	wg.Wait()
	gm.screen.WritePixels(gm.pixels)
}

func calculatePoint(i, j int, centerX, centerY, size float64, pixels *[]byte, wg *sync.WaitGroup) {
	defer wg.Done()
	x := float64(i)*size/screenWidth - size/2 + centerX
	y := (screenHeight-float64(j))*size/screenHeight - size/2 + centerY
	c := complex(x, y)
	z := complex(0, 0)
	it := 0
	for ; it < iterations; it++ {
		z = z*z + c
		if real(z)*real(z)+imag(z)*imag(z) > 4 {
			break
		}
	}
	r, g, b := color(it)
	p := 4 * (i + j*screenWidth)
	pixelSlice := *pixels
	pixelSlice[p] = r
	pixelSlice[p+1] = g
	pixelSlice[p+2] = b
	pixelSlice[p+3] = 0xff
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(g.screen, nil)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeigh int) {
	return 400, 400
}

func main() {
	game := NewGame()
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("mandelbrot")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
