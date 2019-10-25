package main

type Block struct {
	Pos
	W int
	H int
	Color Color
	Exist bool
}

const fieldWidth = 16
const fieldHeight = 10
const blockWidth = 50
const blockHeight = 20

type Field struct {
	Arr [fieldHeight][fieldWidth]Block
}

func InitField(block Block, field Field, level int) Field {

	fieldMap := initField(level)

	for y := 0; y < fieldHeight; y++ {
		row := [fieldWidth]Block{}
		for x := 0; x < fieldWidth; x++ {

			blockOne := block

			if fieldMap[y][x] == "Y" {
				blockOne.Color = Color{255, 255, 0}
			} else if fieldMap[y][x] == "R" {
				blockOne.Color = Color{255, 0, 0}
			} else if fieldMap[y][x] == "0" {
				blockOne.Exist = false
			}

			row[x] = blockOne
		}
		field.Arr[y] = row
	}

	return field
}

func (field *Field) draw(pixels []byte) {

	for y := 0; y < len(field.Arr); y ++ {
		for x := 0; x < len(field.Arr[y]); x ++ {

			if !field.Arr[y][x].Exist {continue}

			field.Arr[y][x].Y = float32(y*blockHeight + 50)
			field.Arr[y][x].X = float32(x*blockWidth)
			field.Arr[y][x].draw(pixels)
		}
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

func BlockCheck(ball *Ball, field *Field) {

	for y := 0; y < len(field.Arr); y ++ {
		for x := 0; x < len(field.Arr[y]); x ++ {

			field.Arr[y][x].Y = float32(y*blockHeight + 50)
			field.Arr[y][x].X = float32(x*blockWidth)

			if !field.Arr[y][x].Exist {continue}

			if
				ball.X + float32(ball.Radius) > field.Arr[y][x].X &&
				ball.X - float32(ball.Radius) < field.Arr[y][x].X + blockWidth &&
				ball.Y + float32(ball.Radius) < field.Arr[y][x].Y + blockHeight &&
				ball.Y - float32(ball.Radius) > field.Arr[y][x].Y {

					//check from which side...
					if ball.X > field.Arr[y][x].X && ball.X < field.Arr[y][x].X + blockWidth {
						field.Arr[y][x].Exist = false
						ball.Yv = -ball.Yv
					} else if ball.Y < field.Arr[y][x].Y + blockHeight && ball.Y > field.Arr[y][x].Y {
						field.Arr[y][x].Exist = false
						ball.Xv = -ball.Xv
					}

			}
		}
	}

}
