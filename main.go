package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
	"os"
)

const winWidth, winHeight int = 800, 600
var speed float32 = 5

var font *ttf.Font
var solid *sdl.Surface
var surface *sdl.Surface

type Pos struct {
	X, Y float32
}

type Color struct {
	R, G, B byte
}

func clear(pixels []byte) {
	for i := range pixels {
		pixels[i] = 0
	}
}

func setPixel(x, y int, c Color, pixels []byte) {
	index := (y*winWidth + x) * 4
	if index < len(pixels) - 4 && index >= 0 {
		pixels[index] = c.R
		pixels[index+1] = c.G
		pixels[index+2] = c.B
	}
}

func main() {

	//err := sdl.Init(sdl.INIT_EVERYTHING)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//defer sdl.Quit()

	window, err := sdl.CreateWindow("!!!", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, int32(winWidth), int32(winHeight), sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer renderer.Destroy()

	// font
	if err := ttf.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize TTF: %s\n", err)
		os.Exit(0)
	}
	if font, err = ttf.OpenFont("./test.ttf", 32); err != nil {
		fmt.Fprint(os.Stderr, "Failed to open font: %s\n", err)
		os.Exit(0)
	}
	defer font.Close()
	if solid, err = font.RenderUTF8Solid("Hello, World!", sdl.Color{255, 0, 0, 255}); err != nil {
		fmt.Fprint(os.Stderr, "Failed to render text: %s\n", err)
		os.Exit(0)
	}
	defer solid.Free()
	if surface, err = window.GetSurface(); err != nil {
		fmt.Fprint(os.Stderr, "Failed to get window surface: %s\n", err)
		os.Exit(0)
	}

	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(winWidth), int32(winHeight))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer tex.Destroy()

	if err = solid.Blit(nil, surface, nil); err != nil {
		fmt.Fprint(os.Stderr, "Failed to put text on window surface: %s\n", err)
		os.Exit(0)
	}


	pixels := make([]byte, winWidth*(winHeight-50)*4)

	//for y := 0 ; y < winHeight; y++ {
	//	for x := 0 ; x < winWidth; x++ {
	//		setPixel(x, y, color{byte(x % 255), byte(y % 255), 0}, pixels)
	//	}
	//}

	player1 := Paddle{Pos{100,500}, 100, 20, Color{255, 255, 255}}

	block := Block{Pos{0,0}, 50, 20, Color{255, 255, 0}, true}
	field := InitField(block, Field{}, 2)

	ball := Ball{Pos{300,300}, 5, speed/10, speed-speed/10,Color{255,255,255}}

	keyState := sdl.GetKeyboardState()

	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}

		clear(pixels)
		player1.update(keyState)
		player1.draw(pixels)

		ball.update(&player1, keyState)
		ball.draw(pixels)

		field.draw(pixels)

		// reflect from blocks
		BlockCheck(&ball, &field)


		tex.Update(nil, pixels, winWidth*4)
		rect := sdl.Rect{0, int32(winHeight -50), int32(winWidth), 50}
		tex.Update(&rect, surface.Pixels(), winWidth*4)
		renderer.Copy(tex, nil, nil)
		renderer.Present()

		sdl.Delay(5)
	}

}