package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"math"
)

const winWidth, winHeight int = 800, 600
var speed float32 = 5

type color struct {
	r, g, b byte
}

type pos struct {
	x, y float32
}

type ball struct {
	pos
	radius int
	xv float32
	yv float32
	color color
}

func (ball *ball) draw(pixels []byte) {
	for y:= -ball.radius ; y < ball.radius ; y++ {
		for x := -ball.radius ; x < ball.radius ; x++ {
			if x*x+y*y < ball.radius*ball.radius {
				setPixel(int(ball.x)+x, int(ball.y)+y, ball.color, pixels)
			}
		}
	}
}

func (ball *ball) update(paddle1 *paddle) {
	ball.x += ball.xv
	ball.y += ball.yv

	//handle collisions
	if int(ball.y)-ball.radius < 0 {
		ball.yv = -ball.yv
	}

	if int(ball.x)+ball.radius > winWidth || int(ball.x) - ball.radius <0 {
		ball.xv = -ball.xv
	}


	if int(ball.y) + ball.radius > int(paddle1.y) - paddle1.h/2 &&
		ball.y + float32(ball.radius) - paddle1.y - float32(paddle1.h/2) <= float32(ball.yv) {
		if int(ball.x) > int(paddle1.x)-paddle1.w/2 && int(ball.x) < int(paddle1.x)+paddle1.w/2 {

			step := (ball.x - paddle1.x) / float32(paddle1.w/2)

			ball.yv = -(ball.yv)
			if ball.xv > 0 && step > 0 {
				ball.xv += step * ball.xv
			} else if ball.xv > 0 && step < 0 {
				ball.xv += step * ball.xv
			} else if ball.xv < 0 && step > 0 {
				ball.xv -= step * ball.xv
			} else if ball.xv < 0 && step < 0 {
				ball.xv -= step * ball.xv
			}

			if ball.xv >= 0 && ball.xv < 0.1 && float32(math.Abs(float64(step))) > 0.5 && step > 0 {
				ball.xv += 3
			} else if ball.xv <= 0 && ball.xv > -0.1 && float32(math.Abs(float64(step))) > 0.5 && step < 0 {
				ball.xv -= 3
			}

			modSpeed := math.Abs(float64(ball.xv)) + math.Abs(float64(ball.yv))
			if modSpeed > float64(speed) {ball.yv = speed - ball.xv }
			//if modSpeed > float64(speed) && ball.xv < 0 {ball.xv = speed + ball.yv }

			fmt.Println("+++", step, "-", ball.xv, ":" , ball.yv,  modSpeed)

		}
	}

	if int(ball.y) > winHeight {
		ball.x = 400
		ball.y = 200
	}

}


type paddle struct {
	pos
	w int
	h int
	color color
}

func (paddle *paddle) draw(pixels []byte) {
	startX := int(paddle.x) - paddle.w/2
	startY := int(paddle.y) - paddle.h/2

	for y := 0; y < paddle.h; y++ {
		for x := 0; x < paddle.w; x++ {
			setPixel(startX+x, startY+y, paddle.color, pixels)
		}
	}
}

func (paddle *paddle) update(keyState []uint8) {
	if keyState[sdl.SCANCODE_RIGHT] != 0 {
		if int(paddle.x) + paddle.w/2 < winWidth {
			paddle.x += 10
		}

	}
	if keyState[sdl.SCANCODE_LEFT] != 0 {
		if int(paddle.x) - paddle.w/2 > 0 {
			paddle.x -= 10
		}
	}
}

func clear(pixels []byte) {
	for i := range pixels {
		pixels[i] = 0
	}
}

func setPixel(x, y int, c color, pixels []byte) {
	index := (y*winWidth + x) * 4
	if index < len(pixels) - 4 && index >= 0 {
		pixels[index] = c.r
		pixels[index+1] = c.g
		pixels[index+2] = c.b
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

	player1 := paddle{pos{100,500}, 100, 20, color{255, 255, 255}}
	ball := ball{pos{300,300}, 5, speed/10, speed-speed/10,color{255,255,255}}

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

		tex.Update(nil, pixels, winWidth*4)
		renderer.Copy(tex, nil, nil)
		renderer.Present()

		sdl.Delay(16)
	}

}