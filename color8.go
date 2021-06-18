package ansi

const (
	colorCubeSize = 6

	MinRGB8   = CountColor4
	CountRGB8 = colorCubeSize * colorCubeSize * colorCubeSize

	MinGray8   = MinRGB8 + CountRGB8
	CountGray8 = 24

	_ = MinGray8 + CountGray8
)

func ColorRGB8(r, g, b uint8) Color {
	clamp := func(x uint8) uint8 {
		if x < colorCubeSize {
			return x
		}
		return colorCubeSize - 1
	}
	r = clamp(r)
	g = clamp(g)
	b = clamp(b)
	return Color(MinRGB8) + Color(r*36+g*6+b)
}

func ColorGray8(intensity uint8) Color {
	if intensity >= CountGray8 {
		intensity = CountGray8 - 1
	}
	return MinGray8 + Color(intensity)
}

func Color8Indexed(index uint8) Color {
	return MinColor4 + Color(index)
}
