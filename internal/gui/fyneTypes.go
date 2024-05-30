package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"time"
)

// customLayout struct to make a custom layout in Fyne
type customLayout struct{}

func (c *customLayout) MinSize(_ []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(1200, 500)
}

func (c *customLayout) Layout(elements []fyne.CanvasObject, _ fyne.Size) {
	// openFolderButton
	elements[0].Resize(fyne.NewSize(170, 36))
	elements[0].Move(fyne.NewPos(10, 10))

	// folder label
	elements[1].Resize(fyne.NewSize(0, 0))
	elements[1].Move(fyne.NewPos(185, 10))

	// Input boxes
	elements[2].Resize(fyne.NewSize(250, 36))
	elements[2].Move(fyne.NewPos(10, 80))

	elements[3].Resize(fyne.NewSize(250, 36))
	elements[3].Move(fyne.NewPos(10, 125))

	// Search button
	elements[4].Resize(fyne.NewSize(100, 36))
	elements[4].Move(fyne.NewPos(75, 180))

	// tableHeaders element layout
	elements[5].Resize(fyne.NewSize(900, 450))
	elements[5].Move(fyne.NewPos(300, 50))
}

// TappableLabel is a custom Label that can be tapped
type TappableLabel struct {
	widget.Label
	OnDoubleTapped func()
	lastTapped     time.Time
}

// NewTappableLabel is a custom label that makes a label tappable
func NewTappableLabel() *TappableLabel {
	tapLabel := new(TappableLabel)
	tapLabel.ExtendBaseWidget(tapLabel)
	return tapLabel
}

const doubleTapTimeout = 500 * time.Millisecond

// Tapped triggers only on doubly click
func (t *TappableLabel) Tapped(_ *fyne.PointEvent) {
	now := time.Now()
	if now.Sub(t.lastTapped) < doubleTapTimeout {
		if t.OnDoubleTapped != nil {
			t.OnDoubleTapped()
		}
	}
	t.lastTapped = now
}

func (t *TappableLabel) TappedSecondary(*fyne.PointEvent) {}
