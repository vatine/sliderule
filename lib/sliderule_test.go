package sliderule

import (
	"bytes"
	"math"
	"reflect"
	"testing"
)

func cmpFloats(a, b, delta float64) bool {
	start := a - b
	if start < 0 {
		start = -start
	}
	return start <= delta
}

func TestPositionalOffset(t *testing.T) {
	length := 50.0
	cases := []struct {
		i, k, tens int
		expected   float64
	}{
		{1, 1, 1, 0.0},
		{3, 1, 1, 23.856},
		{1, 1, 2, 0.0},
		{1, 10, 2, 25.0},
	}

	for ix, c := range cases {
		max := math.Pow10(c.tens)
		seen := positionOffset(float64(c.i * c.k), max, length)
		if !cmpFloats(seen, c.expected, 0.001) {
			t.Errorf("Float %f is sufficiently far from expected %f (case %d)", seen, c.expected, ix)
		}
	}
}

func TestText(t *testing.T) {
	tElem := MakeNumber("1", 0.0, 2.0, 8.0)
	expect := Text{
		text: "1",
		x: -1.33333333333333333,
		y: -2.0,
		size: 8.0,
	}

	if !reflect.DeepEqual(expect, tElem) {
		t.Errorf("Expected %s, saw %s", expect, tElem)
	}

	buf := bytes.NewBufferString("")

	tElem.Render(buf)

	if string(buf.String()) != "<text x=\"-1.333333mm\" y=\"-2.000000mm\">1</text>" {
		t.Errorf("Saw %s", buf.String())
	}
	if tElem.String() != buf.String() {
		t.Errorf("Expected rendered data to be stringified data")
	}
}

func TestScale(t *testing.T) {
	s := MakeLogScale(1, 100.0, 0.0, 10.0, "T")
	
	buf := bytes.NewBufferString("")
	s.Render(buf)
//	expected := ""
//	seen := buf.String()
	
//	if seen != expected {
//		t.Errorf("Expected:\n%s\n\nSaw:\n%s", expected, buf.String())
//	}
}
