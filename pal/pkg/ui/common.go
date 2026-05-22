package ui

type Pos struct {
	X, Y int
}

type Rect struct {
	Top, Left, Height, Width int
}

func (r Rect) Cover(x, y int) bool {
	return x >= r.Left && x < r.Left+r.Width &&
		y >= r.Top && y < r.Top+r.Height
}
