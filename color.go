// Package tcellx provides extensions to github.com/gdamore/tcell/v2 package.
package tcellx

import (
	"github.com/gdamore/tcell/v2"
	"github.com/lucasb-eyer/go-colorful"

	"github.com/pamburus/tcellx/termcap"
)

// Color represents a color in the terminal.
type Color = tcell.Color

// ColorID returns the ID string of the color that can be used to build style tags.
func ColorID(c Color) string {
	if !c.Valid() {
		switch c {
		case tcell.ColorNone:
			return "none"
		case tcell.ColorDefault:
			return "default"
		case tcell.ColorReset:
			return "reset"
		}

		return ""
	}

	if name := colorNames[c]; name != "" {
		return name
	}

	return c.CSS()
}

// AdaptiveColor represents a color that adapts to the terminal color profile.
type AdaptiveColor interface {
	Resolve(termcap.Provider) Color
}

// ---

// NewAdaptiveColor creates a new adaptive color builder.
func NewAdaptiveColor(c Color) AdaptiveColorBuilder {
	pc := portableColor{ANSIColor(c), ANSIExtendedColor(c), c}

	return AdaptiveColorBuilder{adaptiveColor{pc, pc}}
}

// ---

// AdaptiveColorBuilder is a builder for adaptive colors.
type AdaptiveColorBuilder struct {
	adaptiveColor
}

// WithANSI sets the ANSI color for both dark and light modes.
func (b AdaptiveColorBuilder) WithANSI(c Color) AdaptiveColorBuilder {
	c = ANSIColor(c)
	b.dm.ANSI = c
	b.lm.ANSI = c

	return b
}

// WithANSI256 sets the ANSI 256 color for both dark and light modes.
func (b AdaptiveColorBuilder) WithANSI256(c Color) AdaptiveColorBuilder {
	c = ANSIExtendedColor(c)
	b.dm.ANSI256 = c
	b.lm.ANSI256 = c

	return b
}

// WithRGB sets the RGB color for both dark and light modes.
func (b AdaptiveColorBuilder) WithRGB(c Color) AdaptiveColorBuilder {
	b.dm.RGB = c
	b.lm.RGB = c

	return b
}

// WithLightANSI sets the ANSI color for light mode.
func (b AdaptiveColorBuilder) WithLightANSI(c Color) AdaptiveColorBuilder {
	c = ANSIColor(c)
	b.lm.ANSI = c

	return b
}

// WithLightANSI256 sets the ANSI 256 color for light mode.
func (b AdaptiveColorBuilder) WithLightANSI256(c Color) AdaptiveColorBuilder {
	c = ANSIExtendedColor(c)
	b.lm.ANSI256 = c

	return b
}

// WithLightRGB sets the RGB color for light mode.
func (b AdaptiveColorBuilder) WithLightRGB(c Color) AdaptiveColorBuilder {
	b.lm.RGB = c

	return b
}

// Result returns the resulting adaptive color.
func (b AdaptiveColorBuilder) Result() AdaptiveColor {
	return b.adaptiveColor
}

// ---

// ANSIColor returns the nearest ANSI color for the given color.
func ANSIColor(c Color) Color {
	if c == tcell.ColorReset || c == tcell.ColorDefault || c >= tcell.ColorBlack && c <= tcell.ColorWhite {
		return c
	}

	return nearestColor(c, tcell.ColorBlack, tcell.ColorWhite)
}

// ANSIExtendedColor returns the nearest ANSI 256 color for the given color excluding basic 16 ANSI colors.
func ANSIExtendedColor(c Color) Color {
	if c == tcell.ColorReset || c == tcell.ColorDefault || c >= tcell.ColorBlack && c <= tcell.Color255 {
		return c
	}

	return nearestColor(c, tcell.Color16, tcell.Color255)
}

// ---

type adaptiveColor struct {
	dm portableColor // For dark mode
	lm portableColor // For light mode
}

func (c adaptiveColor) Resolve(provider termcap.Provider) Color {
	if provider.LightBackgroundMode() {
		return c.lm.Resolve(provider)
	}

	return c.dm.Resolve(provider)
}

// ---

type portableColor struct {
	ANSI    Color
	ANSI256 Color
	RGB     Color
}

func (c portableColor) Resolve(provider termcap.Provider) Color {
	if profile, ok := provider.ColorProfile(); ok {
		switch profile {
		case termcap.ANSI:
			return c.ANSI
		case termcap.Xterm256:
			return c.ANSI256
		case termcap.TrueColor:
			return c.RGB
		}
	}

	return tcell.ColorDefault
}

// ---

func nearestColor(source, from, to Color) Color {
	sc, ok := toColorful(source)
	if !ok {
		return tcell.ColorDefault
	}

	result := tcell.ColorDefault
	minDistance := 100.0

	for c := from; c <= to; c++ {
		cc, ok := toColorful(c)
		if !ok {
			continue
		}

		distance := cc.DistanceCIEDE2000(sc)
		if distance < minDistance {
			result = c
			minDistance = distance
		}
	}

	return result
}

func toColorful(c Color) (colorful.Color, bool) {
	r, g, b := c.RGB()
	if r < 0 || g < 0 || b < 0 {
		return colorful.Color{}, false
	}

	convert := func(v int32) float64 {
		return float64(v) / 255.0
	}

	return colorful.Color{R: convert(r), G: convert(g), B: convert(b)}, true
}

func collectColorNames() map[Color]string {
	colorNames := make(map[Color]string, len(tcell.ColorNames))

	for name, c := range tcell.ColorNames {
		colorNames[c] = name
	}

	return colorNames
}

// ---

var colorNames = collectColorNames()

// ---

var (
	_ AdaptiveColor = adaptiveColor{}
	_ AdaptiveColor = portableColor{}
)
