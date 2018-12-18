package sliderule

import (
	"bytes"
	"fmt"
	"io"
	"math"
)

const dashWidth float64 = 0.1

type Element interface {
	Render(io.Writer)
	String() string
}

type Text struct {
	text string
	size float64
	x float64
	y float64
}

type subDash struct {
	val, len float64
}

type Dash struct {
	x, y, l float64
}

type Scale struct {
	name Text
	elems []Element
}

type Sliderule struct {
	scales []Scale
	width  float64
	height float64
}

func (s Scale) Render(w io.Writer) {
	for _, e := range s.elems {
		e.Render(w)
	}
	s.name.Render(w)
}

// Compute the log of v in the specified base
func log(v, base float64) float64 {
	return math.Log(v) / math.Log(base)
}

// Returns a Text element with the specified text, centred on and
// offset from the x, y coordinates provided.
func MakeNumber(t string, x, y float64, size float64) Text {
	xScale := float64(len(t)) / 6.0
	return Text{
		text: t,
		x: x - (xScale * size),
		y: y - size / 2.0,
		size: size}
}

// Render a Sliderule as a string
func (s Sliderule) String() string {
	buf := bytes.NewBuffer([]byte{})

	s.Render(buf)
	return buf.String()
}

// Render a Sliderule to the specified io.Writer
func (s Sliderule) Render(w io.Writer) {
	fmt.Fprintln(w, `<?xml version="1.0" standalone="no"?>`)
	fmt.Fprintln(w, "\n<!DOCTYPE svg PUBLIC \"-//W3C//DTD SVG 1.1//EN\" \"http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd\">")
	fmt.Fprintf(w, `<svg width="%fmm" heigth="%fmm" version="1.1" xmlns="http://www.w3.org/2000/svg">`, s.width, s.height)

	for _, scale := range s.scales {
		scale.Render(w)
	}
	fmt.Fprintf(w, "\n</svg>\n")
}

// Render a Text element to an output stream
func (t Text) Render(w io.Writer) {
	fmt.Fprintf(w, `<text x="%fmm" y="%fmm">%s</text>`, t.x, t.y, t.text)
}

// Return the string representation of a Text
func (t Text) String() string {
	buf := bytes.NewBuffer([]byte{})

	t.Render(buf)
	return buf.String()
}

// Render a Dash
func (d Dash) Render (w io.Writer) {
	fmt.Fprintf(w, `<line x1="%fmm" y1="%fmm" x2="%fmm" y2="%fmm" style="stroke:rgb(0,0,0);stroke-width:%f" />`, d.x, d.y, d.x, d.y + d.l, dashWidth)
}

// Render a Dash as a string
func (d Dash) String() string {
	buf := bytes.NewBuffer([]byte{})

	d.Render(buf)
	return buf.String()
}

// Compute the offset of a scale decoration. This is essentially the
// log of the decoration, with the maximum value of the scale as the
// base, multiplied by the length of the scale.
func positionOffset(number, base float64, length float64) float64 {
	factor := log(number, base)
//	fmt.Printf("DEBUG:\n  number = %f\n  base = %f\n  factor = %f\n  rv = %f\n\n", number, base, factor, length * factor)
	return length * factor
}

// Return a slice of subDashes, basically just a series of dashes with
// a specified length. Ensure that sbuDashes do not run together.
func makeSubScale(start, step, base, length float64) []subDash {
	halfway := start + 0.5 * step
	end := start + step
	tenth := end - 0.1 * step
	twentieth := end - 0.05 * step
	var rv []subDash
	endPosition := positionOffset(end, base, length)

	if math.Abs(endPosition - positionOffset(twentieth, base, length)) >= 1.0 {
		step := 0.1 * step
		for x := start + step / 2; x < end; x = x + step {
			pos := positionOffset(x, base, length)
			rv = append(rv, subDash{pos, 1.0})
		}
	}
	if math.Abs(endPosition - positionOffset(tenth, base, length)) >= 1.0 {
		step := 0.1 * step
		n := 1
		for x := start + step; x < end; x = x + step {
			if n != 5 {
				pos := positionOffset(x, base, length)
				rv = append(rv, subDash{pos, 2.0})
			}
			n++
		}
	}
	if math.Abs(endPosition - positionOffset(halfway, base, length)) >= 1.0 {
		pos := positionOffset(halfway, base, length)
		rv = append(rv, subDash{pos, 3.0})
	}

	return rv
}

func sign(f float64) float64 {
	if f == 0.0 {
		return 0.0
	}
	if f > 0.0 {
		return 1.0
	}
	return -1.0
}

// Return a scale, complete with numbers and dashes, for the given
// number of "consecutive powers of ten" passed in
func MakeLogScale(tens int, l, xOffset, yOffset float64, name string) Scale {
	//base := math.Pow(10.0, float64(tens))
	elems := []Element{}
	max := math.Pow10(tens)
	nameXOffset := xOffset - 10.0
	if l < 0.0 {
		// Reverse scale, stick the name on the opposite side
		nameXOffset = (l+xOffset) - 10.0
	}
	nameYOffset := yOffset + 4.0

	for i := 1; i <= 10; i++ {
		k := 1
		for j := 1; j <= tens; j++ {
			offset := xOffset + positionOffset(float64(i * k), max, l)
			t := fmt.Sprintf("%d", i * k)
			elems = append(elems, Element(MakeNumber(t, offset, yOffset + 2, 8.0)))
			elems = append(elems, Dash{x: offset, y: yOffset, l: 4.0})
			if i < 10 {
				for _, dash := range makeSubScale(float64(i * k), float64(k), max, l) {
					elems = append(elems, Dash{x: xOffset + dash.val, y: yOffset, l: dash.len})
				}
			}
			k = k * 10
		}
	}

	return Scale{name: MakeNumber(name, nameXOffset, nameYOffset, 8.0), elems: elems}
}

func MakeLinScale(tens int, l, xOffset, yOffset float64, name string) Scale {
	elems := []Element{}
	max := math.Pow10(tens)
	_ = max
	nameXOffset := xOffset - 10.0
	if l < 0.0 {
		// Reverse scale, stick the name on the opposite side
		nameXOffset = (l+xOffset) - 10.0
	}
	nameYOffset := yOffset + 4.0

	return Scale{name: MakeNumber(name, nameXOffset, nameYOffset, 8.0), elems: elems}
}

// Create a sliderule, with some random scales
func MakeSlideRule(width, height float64) Sliderule {
	s := Sliderule{width: width, height: height}
	l := width - 30.0

	s.scales = append(s.scales, MakeLogScale(1, l, 15.0, height/6.0 , "A"))
	s.scales = append(s.scales, MakeLogScale(2, l, 15.0, height/3, "B"))
	s.scales = append(s.scales, MakeLogScale(1, -l, l+15.0, 2*height/3, "C"))

	return s
}
