package ansi

const (
	MinRGB24   = MinGray8 + CountGray8
	CountRGB24 = 1 << 24

	_ = MinRGB24 + (CountRGB24 - 1)
)

func ColorRGB24(r, g, b uint8) Color {
	rgb := (Color(r) << 16) | (Color(g) << 8) | Color(b)
	return MinRGB24 + rgb
}
