package ansi

// https://en.wikipedia.org/wiki/ANSI_escape_code#SGR

const (
	maxCode = 128

	setForegroundColorEx = 38
	setBackgroundColorEx = 48
)

var sgrParamToAction = [maxCode]Action{
	0:  Reset{},
	1:  SetBold(true),
	2:  SetFaint(true),
	3:  SetItalic(true),
	4:  SetUnderline(true),
	5:  SetBlink(true),
	7:  SetInverted(true),
	20: SetFraktur(true),

	30: SetForeground(Black),
	31: SetForeground(Red),
	32: SetForeground(Green),
	33: SetForeground(Yellow),
	34: SetForeground(Blue),
	35: SetForeground(Magenta),
	36: SetForeground(Cyan),
	37: SetForeground(White),
	39: SetForeground(DefaultColor),

	40: SetBackground(Black),
	41: SetBackground(Red),
	42: SetBackground(Green),
	43: SetBackground(Yellow),
	44: SetBackground(Blue),
	45: SetBackground(Magenta),
	46: SetBackground(Cyan),
	47: SetBackground(White),
	49: SetBackground(DefaultColor),

	90: SetForeground(BrightBlack),
	91: SetForeground(BrightRed),
	92: SetForeground(BrightGreen),
	93: SetForeground(BrightYellow),
	94: SetForeground(BrightBlue),
	95: SetForeground(BrightMagenta),
	96: SetForeground(BrightCyan),
	97: SetForeground(BrightWhite),

	100: SetBackground(BrightBlack),
	101: SetBackground(BrightRed),
	102: SetBackground(BrightGreen),
	103: SetBackground(BrightYellow),
	104: SetBackground(BrightBlue),
	105: SetBackground(BrightMagenta),
	106: SetBackground(BrightCyan),
	107: SetBackground(BrightWhite),
}

func sgrColor8(code int) Color {
	var index byte
	if code >= 0 && code <= 0xFF {
		index = byte(code)
	} else if code > 0xFF {
		index = 0xFF
	}
	return Color8Indexed(index)
}

func sgrColor24(cr, cg, cb int) Color {
	clamp := func(n int) uint8 {
		if n < 0 {
			return 0
		}
		if n > 0xFF {
			return 0xFF
		}
		return uint8(n)
	}
	return ColorRGB24(clamp(cr), clamp(cg), clamp(cb))
}

func sgrColorExtended(codes []maybeInt) (color Color, rem []maybeInt) {
	const (
		mode8bits  = 5
		mode24bits = 2
	)
	if len(codes) < 2 {
		//8-bit color requires 2 params: [5;n]
		return
	}

	switch codes[0].withDefault(-1) {
	case mode8bits:
		color = sgrColor8(codes[1].withDefault(0))
		rem = codes[2:]

	case mode24bits:
		if len(codes) < 4 {
			//24-bit color requires 4 params: [2;r;g;b]
			return
		}
		color = sgrColor24(
			codes[1].withDefault(0),
			codes[2].withDefault(0),
			codes[3].withDefault(0),
		)
		rem = codes[4:]

	default:
		//invalid/unknown, dropped
		rem = codes[1:]
	}
	return
}

func sgrLookup(codes []maybeInt) (action Action, rem []maybeInt) {
	if len(codes) == 0 {
		return
	}

	c0 := codes[0].withDefault(0)
	if c0 >= maxCode || c0 < 0 {
		rem = codes[1:]
		return
	}

	var foreground bool
	if action = sgrParamToAction[c0]; action != nil {
		rem = codes[1:]
		return
	}

	switch c0 {
	case setForegroundColorEx:
		foreground = true
		fallthrough

	case setBackgroundColorEx:
		var color Color
		color, rem = sgrColorExtended(codes[1:])

		if foreground {
			action = SetForeground(color)
		} else {
			action = SetBackground(color)
		}
		return

	default:
		//invalid/unknown, dropped
		rem = codes[1:]
		return
	}
}
