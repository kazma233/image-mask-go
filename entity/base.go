package entity

import (
	"image"
	"image/color"

	"golang.org/x/image/math/fixed"
)

type (
	Shape struct {
		Width int
		Hight int
	}

	Box struct {
		Shape
		image.Point
	}

	ColorPoint struct {
		C  color.Color
		Pt fixed.Point26_6
	}

	ColorBox struct {
		Box
		color.Color
	}
)