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

func InitField(block Block, field Field) Field {

	row := [fieldWidth]Block{}

	for i := 0; i < len(row); i++ {
		row[i] = block
	}

	for i := 0; i < len(field.Arr); i++ {
		field.Arr[i] = row
	}

	return field
}

func (field *Field) draw(pixels []byte) {

	for y := 0; y < len(field.Arr); y ++ {
		for x := 0; x < len(field.Arr[y]); x ++ {

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
