package ansi

import (
	"encoding/json"
	"fmt"
)

type Color uint32

const (
	DefaultColor Color = iota
	Black
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
	BrightBlack
	BrightRed
	BrightGreen
	BrightYellow
	BrightBlue
	BrightMagenta
	BrightCyan
	BrightWhite

	CountColor4 = iota
	MinColor4   = Black
)

const (
	gray8Low  = 0x0808
	gray8High = 0xeeee
)

var Color4Names = [...]string{
	"black",
	"red",
	"green",
	"yellow",
	"blue",
	"magenta",
	"cyan",
	"white",

	"bright-black",
	"bright-red",
	"bright-green",
	"bright-yellow",
	"bright-blue",
	"bright-magenta",
	"bright-cyan",
	"bright-white",
}

type ColorPalette [CountColor4]uint32

var XTermPalette = ColorPalette{
	0xeeeeee,

	0x000000,
	0xcd0000,
	0x00cd00,
	0xcdcd00,
	0x0000ee,
	0xcd00cd,
	0x00cdcd,
	0x5e5e5e,

	0x7f7f7f,
	0xff0000,
	0x00ff00,
	0xffff00,
	0x5c5cff,
	0xff00ff,
	0x00ffff,
	0xffffff,
}

func (c Color) String() string {
	switch {
	case c == DefaultColor:
		return ""

	case c < CountColor4:
		return Color4Names[c-1]

	case c < MinRGB8+CountRGB8:
		fallthrough
	case c < MinGray8+CountGray8:
		fallthrough
	case c < MinRGB24+CountRGB24:
		r, g, b, _ := c.RGBA()
		r >>= 8
		g >>= 8
		b >>= 8
		pack := (r << 16) | (g << 8) | b
		return fmt.Sprintf("#%06X", pack)

	default: //out of range
		return ""
	}
}

func (c Color) RGBA() (r, g, b, a uint32) {
	return c.PaletteRGBA(XTermPalette)
}

func (c Color) PaletteRGBA(palette ColorPalette) (r, g, b, a uint32) {
	switch {
	case c < CountColor4:
		tuple := palette[c]
		r = (tuple >> 16) & 0xFF
		g = (tuple >> 8) & 0xFF
		b = tuple & 0xFF
		r |= r << 8
		g |= g << 8
		b |= b << 8
		a = 0xffff
		return

	case c < MinRGB8+CountRGB8:
		index := uint32(c - MinRGB8)
		b = ((index % 6) * 0xffff) / 5
		g = (((index / 36) % 6) * 0xffff) / 5
		r = (((index / 36) % 6) * 0xffff) / 5
		a = 0xffff
		return

	case c < MinGray8+CountGray8:
		index := byte(c - MinGray8)
		r = gray8Low + uint32(index)*(gray8High-gray8Low)/(CountGray8-1)
		g = r
		b = r
		a = 0xffff
		return

	case c < MinRGB24+CountRGB24:
		packed := uint32(c - MinRGB24)
		r = (packed >> 16) & 0xFF
		g = (packed >> 8) & 0xFF
		b = packed & 0xFF
		r |= r << 8
		g |= g << 8
		b |= b << 8
		a = 0xffff
		return

	default: //out of range
		return
	}
}

func (c Color) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}
