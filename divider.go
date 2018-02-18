package micro

import "image"

type LineType int

func (l LineType) Rune() rune {
	switch l {
	case HorizontalLine:
		return '─'
	case VerticalLine:
		return '│'
	default:
		log.Panic("bug")
		panic("bug")
	}
}

type LineTermination int

func (p LineTermination) Rune(defaultRune rune) rune {
	switch p {
	case LineTerminationNone:
		return defaultRune
	case LineTerminationNormal:
		return '○'
	case LineTerminationHighlight:
		return '◎'
	default:
		log.Panic("bug")
		panic("bug")
	}
}

const (
	HorizontalLine LineType = iota
	VerticalLine
	LineTerminationNormal LineTermination = iota
	LineTerminationHighlight
	LineTerminationNone
)

type Line struct {
	Start     image.Point
	Length    int
	Type      LineType
	StartDeco LineTermination
	EndDeco   LineTermination
	StyleName string
}
