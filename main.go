package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"math"
)

const winWidth, winHeight int = 800, 600
var speed float32 = 5

type Block struct {
	Pos
	W int
	H int
	Color Color
}

type Pos struct {
	X, Y float32
}

type Ball struct {
	Pos
	Radius int
	Xv float32
	Yv float32
	Color Color
}

type Paddle struct {
	Pos
	W int
	H int
	Color Color
}

type Color struct {
	R, G, B byte
}

func (ball *Ball) draw(pixels []byte) {
	for y:= -ball.Radius ; y < ball.Radius ; y++ {
		for x := -ball.Radius ; x < ball.Radius ; x++ {
			if x*x+y*y < ball.Radius*ball.Radius {
				setPixel(int(ball.X)+x, int(ball.Y)+y, ball.Color, pixels)
			}
		}
	}
}

func (ball *Ball) update(paddle1 *Paddle) {
	ball.X += ball.Xv
	ball.Y += ball.Yv

	//handle collisions
	if int(ball.Y)-ball.Radius < 0 {
		ball.Yv = -ball.Yv
	}

	if int(ball.X)+ball.Radius > winWidth || int(ball.X) - ball.Radius <0 {
		ball.Xv = -ball.Xv
	}


	// reflect logic
	if int(ball.Y) + ball.Radius > int(paddle1.Y) - paddle1.H/2 &&
		ball.Yv > 0 &&
		ball.Y + float32(ball.Radius) - paddle1.Y - float32(paddle1.H/2) <= float32(ball.Yv) {
		if int(ball.X) > int(paddle1.X)-paddle1.W/2 && int(ball.X) < int(paddle1.X)+paddle1.W/2 {

			step := (ball.X - paddle1.X) / float32(paddle1.W/2)

			// ball reflecting from paddle
			ball.Yv = -(ball.Yv)
			if ball.Xv > 0 && step > 0 {
				ball.Xv += step * ball.Xv
			} else if ball.Xv > 0 && step < 0 {
				ball.Xv += step * ball.Xv
			} else if ball.Xv < 0 && step > 0 {
				ball.Xv -= step * ball.Xv
			} else if ball.Xv < 0 && step < 0 {
				ball.Xv -= step * ball.Xv
			}

			// change angle if the trajectory is too vertical
			if ball.Xv >= 0 && ball.Xv < 0.1 && float32(math.Abs(float64(step))) > 0.5 && step > 0 {
				ball.Xv += 3
			} else if ball.Xv <= 0 && ball.Xv > -0.1 && float32(math.Abs(float64(step))) > 0.5 && step < 0 {
				ball.Xv -= 3
			}

			// correct speed
			modSpeed := math.Abs(float64(ball.Xv)) + math.Abs(float64(ball.Yv))
			if modSpeed > float64(speed) || modSpeed < float64(speed) {
				ball.Yv = -(speed - float32(math.Abs(float64(ball.Xv))))
				fmt.Println("!+!")
			}

			fmt.Println("+++", step, "-", ball.Xv, ":" , ball.Yv,  modSpeed)

		}
	}

	if int(ball.Y) > winHeight {
		ball.X = 400
		ball.Y = 200
	}

}

func (block *Block) draw(pixels []byte) {

	blocColor := block.Color
	startX := int(block.X)
	startY := int(block.Y)

	for y := 0; y < block.H; y++ {
		for x := 0; x < block.W; x++ {
			if (x == 0 || y == 0) {
				block.Color = Color{blocColor.R +50, blocColor.G + 50, blocColor.B + 50}
			} else if (x == block.W || y == block.H) {
				block.Color = Color{blocColor.R - 50, blocColor.G - 50, blocColor.B - 50}
			} else {
				block.Color = blocColor
			}

			setPixel(startX+x, startY+y, block.Color, pixels)
		}
	}
}

func (paddle *Paddle) draw(pixels []byte) {
	startX := int(paddle.X) - paddle.W/2
	startY := int(paddle.Y) - paddle.H/2

	for y := 0; y < paddle.H; y++ {
		for x := 0; x < paddle.W; x++ {
			setPixel(startX+x, startY+y, paddle.Color, pixels)
		}
	}
}

func (paddle *Paddle) update(keyState []uint8) {
	if keyState[sdl.SCANCODE_RIGHT] != 0 {
		if int(paddle.X) + paddle.W/2 < winWidth {
			paddle.X += 10
		}

	}
	if keyState[sdl.SCANCODE_LEFT] != 0 {
		if int(paddle.X) - paddle.W/2 > 0 {
			paddle.X -= 10
		}
	}
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

	tex, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(winWidth), int32(winHeight))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer tex.Destroy()

	pixels := make([]byte, winWidth*winHeight*4)

	//for y := 0 ; y < winHeight; y++ {
	//	for x := 0 ; x < winWidth; x++ {
	//		setPixel(x, y, color{byte(x % 255), byte(y % 255), 0}, pixels)
	//	}
	//}

	player1 := Paddle{Pos{100,500}, 100, 20, Color{255, 255, 255}}

	block1 := Block{Pos{100,100}, 50, 20, Color{255, 255, 0}}
	block2 := Block{Pos{150,100}, 50, 20, Color{255, 255, 0}}
	block3 := Block{Pos{200,100}, 50, 20, Color{255, 255, 0}}

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

		ball.draw(pixels)
		ball.update(&player1)

		block1.draw(pixels)
		block2.draw(pixels)
		block3.draw(pixels)

		tex.Update(nil, pixels, winWidth*4)
		renderer.Copy(tex, nil, nil)
		renderer.Present()

		sdl.Delay(5)
	}

}