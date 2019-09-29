package main

import (
	"math"
	"fmt"
)

type Ball struct {
	Pos
	Radius int
	Xv float32
	Yv float32
	Color Color
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
