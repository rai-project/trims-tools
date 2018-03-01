package ui

import (
	"image"

	"github.com/marcusolsson/tui-go"
)

// TailBox is a container Widget that may not show all its Widgets.
// While tui.Box attempts to show every contained Widget - sometimes shrinking
// those Widgets to do so- TailBox prioritizes completely displaying its last
// Widget, then the next-to-last widget, etc.
// It is vertically-aligned, i.e. all the contained Widgets have the same width.
type TailBox struct {
	tui.WidgetBase
	sz       image.Point
	contents []tui.Widget
}

var _ tui.Widget = &TailBox{}

// NewTailBox returns a new TailBox widget.
func NewTailBox(w ...tui.Widget) *TailBox {
	return &TailBox{
		contents: w,
	}
}

// SetContents sets the component widgets of the TailBox.
func (t *TailBox) SetContents(w ...tui.Widget) {
	t.contents = w
	t.doLayout(t.Size())
}

// Draw renders the TailBox.
func (t *TailBox) Draw(p *tui.Painter) {
	p.WithMask(image.Rect(0, 0, t.sz.X, t.sz.Y), func(p *tui.Painter) {
		// Draw background
		p.FillRect(0, 0, t.sz.X, t.sz.Y)

		// Draw from the bottom up.
		space := t.sz.Y
		p.Translate(0, space)
		defer p.Restore()
		for i := len(t.contents) - 1; i >= 0 && space > 0; i-- {
			w := t.contents[i]
			space -= w.Size().Y
			p.Translate(0, -w.Size().Y)
			defer p.Restore()
			w.Draw(p)
		}
	})
}

// Resize recalculates the layout of the box's contents.
func (t *TailBox) Resize(size image.Point) {
	t.WidgetBase.Resize(size)
	defer func() {
		t.sz = size
	}()

	// If it's just a height change, Draw should do the right thing already.
	if size.X != t.sz.X {
		t.doLayout(size)
	}
}

func (t *TailBox) doLayout(size image.Point) {
	for _, w := range t.contents {
		hint := w.SizeHint()
		// Set the width to the container width, and height to the requested height
		w.Resize(image.Pt(size.X, hint.Y))
		// ...and then resize again, now that the Y-hint has been refreshed by the X-value.
		hint = w.SizeHint()
		w.Resize(image.Pt(size.X, hint.Y))
	}
}

/*
MIT License

Copyright (c) 2017 Charles C. Eckman

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
