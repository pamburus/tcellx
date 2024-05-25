// Package termcap provides an interface for querying terminal capabilities.
package termcap

// ---

// ColorProfile represents the color profile of the terminal.
type ColorProfile int

// Color profiles.
const (
	Monochrome ColorProfile = iota
	ANSI                    // 16 colors
	Xterm256                // 256 colors
	TrueColor               // 16 million colors
)

// ---

// RGB represents a 24-bit color in RGB space.
type RGB struct {
	R, G, B uint8
}
