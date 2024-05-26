// Package termenvcap provides an implementation of termcap interfaces on top of termenv module.
package termenvcap

import (
	"github.com/muesli/termenv"

	"github.com/pamburus/tcellx/termcap"
)

// New returns a new termcap provider implemented over termenv package.
func New() termcap.Provider {
	return &provider{termcap.Stub()}
}

// ---

type provider struct {
	termcap.Provider
}

func (p *provider) ColorProfile() (termcap.ColorProfile, bool) {
	profile := termenv.ColorProfile()
	switch profile {
	case termenv.ANSI:
		return termcap.ANSI, true
	case termenv.ANSI256:
		return termcap.Xterm256, true
	case termenv.TrueColor:
		return termcap.TrueColor, true
	}

	return termcap.Monochrome, false
}

func (p *provider) BackgroundColor() (termcap.RGB, bool) {
	r, g, b := termenv.ConvertToRGB(termenv.BackgroundColor()).RGB255()

	return termcap.RGB{R: r, G: g, B: b}, true
}

func (p *provider) ForegroundColor() (termcap.RGB, bool) {
	r, g, b := termenv.ConvertToRGB(termenv.ForegroundColor()).RGB255()

	return termcap.RGB{R: r, G: g, B: b}, true
}

func (p *provider) LightBackgroundMode() bool {
	return !termenv.HasDarkBackground()
}
