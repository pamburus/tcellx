package tcellx

import (
	"strings"

	"github.com/gdamore/tcell/v2"
)

// ---

// Style represents a style in the terminal.
type Style struct {
	fg       Color
	bg       Color
	attrsOn  tcell.AttrMask
	attrsOff tcell.AttrMask
	url      string
	urlID    string
}

// Apply applies the style to the base style.
func (s Style) Apply(base tcell.Style) tcell.Style {
	fg, bg, attrs := base.Decompose()

	if s.fg != tcell.ColorDefault {
		fg = s.fg
	}

	if s.bg != tcell.ColorDefault {
		bg = s.bg
	}

	attrs |= s.attrsOn
	attrs &^= s.attrsOff

	result := base.
		Attributes(attrs).
		Foreground(fg).
		Background(bg)

	if s.url != "" {
		result = result.Url(s.url)
	}

	if s.urlID != "" {
		result = result.UrlId(s.urlID)
	}

	return result
}

// WithForeground sets the foreground color.
func (s Style) WithForeground(c Color) Style {
	s.fg = c

	return s
}

// WithBackground sets the background color.
func (s Style) WithBackground(c Color) Style {
	s.bg = c

	return s
}

// WithAttrs sets the attributes.
func (s Style) WithAttrs(attrs tcell.AttrMask) Style {
	s.attrsOn |= (attrs & ^s.attrsOff)
	s.attrsOff &= ^attrs

	return s
}

// WithoutAttrs removes the attributes.
func (s Style) WithoutAttrs(attrs tcell.AttrMask) Style {
	s.attrsOff |= (attrs & ^s.attrsOn)
	s.attrsOn &= ^attrs

	return s
}

// WithBold sets or removes the bold attribute.
func (s Style) WithBold(on bool) Style {
	return s.withAttrsPatch(tcell.AttrBold, on)
}

// WithDim sets or removes the dim attribute.
func (s Style) WithDim(on bool) Style {
	return s.withAttrsPatch(tcell.AttrDim, on)
}

// WithBlink sets or removes the blink attribute.
func (s Style) WithBlink(on bool) Style {
	return s.withAttrsPatch(tcell.AttrBlink, on)
}

// WithItalic sets or removes the italic attribute.
func (s Style) WithItalic(on bool) Style {
	return s.withAttrsPatch(tcell.AttrItalic, on)
}

// WithReverse sets or removes the reverse attribute.
func (s Style) WithReverse(on bool) Style {
	return s.withAttrsPatch(tcell.AttrReverse, on)
}

// WithUnderline sets or removes the underline attribute.
func (s Style) WithUnderline(on bool) Style {
	return s.withAttrsPatch(tcell.AttrUnderline, on)
}

// WithStrikeThrough sets or removes the strike-through attribute.
func (s Style) WithStrikeThrough(on bool) Style {
	return s.withAttrsPatch(tcell.AttrStrikeThrough, on)
}

// WithURL sets the URL.
func (s Style) WithURL(url string) Style {
	s.url = url

	return s
}

// WithURLID sets the URL ID.
func (s Style) WithURLID(urlID string) Style {
	s.urlID = urlID

	return s
}

// String returns the style as a string that can be used as a style tag.
func (s Style) String() string {
	if s == (Style{}) {
		return ""
	}

	var sb strings.Builder

	sb.WriteByte('[')

	color := func(c Color) {
		if c == tcell.ColorReset {
			sb.WriteByte('-')
		} else {
			sb.WriteString(ColorID(c))
		}
	}

	if c := s.fg; c != tcell.ColorDefault {
		color(c)
	}

	sb.WriteByte(':')
	if c := s.bg; c != tcell.ColorDefault {
		color(c)
	}

	sb.WriteByte(':')
	if s.attrsOff != 0 {
		sb.WriteByte('-')
	}
	if s.attrsOn != 0 {
		if s.attrsOn&tcell.AttrBold != 0 {
			sb.WriteString("b")
		}
		if s.attrsOn&tcell.AttrBlink != 0 {
			sb.WriteString("l")
		}
		if s.attrsOn&tcell.AttrDim != 0 {
			sb.WriteString("d")
		}
		if s.attrsOn&tcell.AttrItalic != 0 {
			sb.WriteString("i")
		}
		if s.attrsOn&tcell.AttrReverse != 0 {
			sb.WriteString("r")
		}
		if s.attrsOn&tcell.AttrUnderline != 0 {
			sb.WriteString("u")
		}
		if s.attrsOn&tcell.AttrStrikeThrough != 0 {
			sb.WriteString("s")
		}
	}

	if s.url != "" {
		sb.WriteByte(':')
		sb.WriteString(s.url)
	}

	sb.WriteByte(']')

	return sb.String()
}

func (s Style) withAttrsPatch(attrs tcell.AttrMask, on bool) Style {
	if on {
		return s.WithAttrs(attrs)
	}

	return s.WithoutAttrs(attrs)
}
