package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"time"
)

// customLayout struct to make a custom layout in Fyne
type customLayout struct{}

func (c *customLayout) MinSize(_ []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(1170, 500)
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
	elements[2].Move(fyne.NewPos(10, 180))

	elements[3].Resize(fyne.NewSize(250, 36))
	elements[3].Move(fyne.NewPos(10, 225))

	// Search button
	elements[4].Resize(fyne.NewSize(100, 36))
	elements[4].Move(fyne.NewPos(75, 280))

	// Progressbar
	elements[5].Resize(fyne.NewSize(250, 36))
	elements[5].Move(fyne.NewPos(10, 350))

	// tableHeaders element layout
	elements[6].Resize(fyne.NewSize(900, 450))
	elements[6].Move(fyne.NewPos(270, 50))

	// Info dialog
	elements[7].Resize(fyne.NewSize(100, 36))
	elements[7].Move(fyne.NewPos(75, 400))

	// TODO Multiline entry element layout
	//elements[8].Resize(fyne.NewSize(895, 190))
	//elements[8].Move(fyne.NewPos(270, 505))
}

// TappableLabel is a custom Label that can be tapped
type TappableLabel struct {
	widget.Label
	OnDoubleTapped func()
	OnTapped       func()
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
	// Handle one tap
	if t.OnTapped != nil {
		t.OnTapped()
	}

	// Handle double tap
	now := time.Now()
	if now.Sub(t.lastTapped) < doubleTapTimeout {
		if t.OnDoubleTapped != nil {
			t.OnDoubleTapped()
		}
	}
	t.lastTapped = now
}

func (t *TappableLabel) TappedSecondary(*fyne.PointEvent) {}

func (t *TappableLabel) MinSize() fyne.Size {
	return fyne.NewSize(34, 34)
}

type sortStatus struct {
	nameAsc bool
	dateAsc bool
	pathAsc bool
}
