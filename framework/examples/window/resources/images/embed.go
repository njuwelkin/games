package images

import (
	_ "embed"
)

var (
	//go:embed gobang.png
	Gobang_png []byte

	//go:embed piece.png
	Piece_png []byte
)
