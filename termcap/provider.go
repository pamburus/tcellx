package termcap

// Stub returns a provider that always returns zero values.
func Stub() Provider {
	return stubInstance
}

// Snapshot returns a snapshot of the provider.
func Snapshot(p Provider) Provider {
	return snapshot(p)
}

// Provider is an interface for querying terminal capabilities.
type Provider interface {
	ColorProfile() (ColorProfile, bool)
	BackgroundColor() (RGB, bool)
	ForegroundColor() (RGB, bool)
	LightBackgroundMode() bool

	sealed()
}

// ---

var stubInstance = &stubProvider{}

// ---

type stubProvider struct{}

func (p *stubProvider) ColorProfile() (ColorProfile, bool) {
	return Monochrome, false
}

func (p *stubProvider) BackgroundColor() (RGB, bool) {
	return RGB{}, false
}

func (p *stubProvider) ForegroundColor() (RGB, bool) {
	return RGB{}, false
}

func (p *stubProvider) LightBackgroundMode() bool {
	return false
}

func (p *stubProvider) sealed() {}

// ---

func snapshot(p Provider) *providerSnapshot {
	if s, ok := p.(*providerSnapshot); ok {
		return s
	}

	s := &providerSnapshot{}
	if cp, ok := p.ColorProfile(); ok {
		s.colorProfile = cp
		s.colorProfileValid = true
	}
	if bg, ok := p.BackgroundColor(); ok {
		s.backgroundColor = bg
		s.backgroundColorValid = true
	}
	if fg, ok := p.ForegroundColor(); ok {
		s.foregroundColor = fg
		s.foregroundColorValid = true
	}
	s.lightBackgroundMode = p.LightBackgroundMode()

	return s
}

// ---

type providerSnapshot struct {
	colorProfile         ColorProfile
	colorProfileValid    bool
	backgroundColor      RGB
	backgroundColorValid bool
	foregroundColor      RGB
	foregroundColorValid bool
	lightBackgroundMode  bool
}

func (p *providerSnapshot) ColorProfile() (ColorProfile, bool) {
	return p.colorProfile, p.colorProfileValid
}

func (p *providerSnapshot) BackgroundColor() (RGB, bool) {
	return p.backgroundColor, p.backgroundColorValid
}

func (p *providerSnapshot) ForegroundColor() (RGB, bool) {
	return p.foregroundColor, p.foregroundColorValid
}

func (p *providerSnapshot) LightBackgroundMode() bool {
	return p.lightBackgroundMode
}

func (p *providerSnapshot) sealed() {}
