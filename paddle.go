package main

import "github.com/veandco/go-sdl2/sdl"

type Paddle struct {
	Pos
	W int
	H int
	Color Color
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

func (paddle *Paddle) draw(pixels []byte) {
	startX := int(paddle.X) - paddle.W/2
	startY := int(paddle.Y) - paddle.H/2

	for y := 0; y < paddle.H; y++ {
		for x := 0; x < paddle.W; x++ {
			setPixel(startX+x, startY+y, paddle.Color, pixels)
		}
	}
}
