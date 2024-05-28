package gui

import "fyne.io/fyne/v2"

// customLayout struct to make a custom layout in Fyne
type customLayout struct{}

func (c *customLayout) MinSize(_ []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(1200, 600)
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
	elements[2].Move(fyne.NewPos(575, 75))

	elements[3].Resize(fyne.NewSize(250, 36))
	elements[3].Move(fyne.NewPos(830, 75))

	// Search button
	elements[4].Resize(fyne.NewSize(100, 36))
	elements[4].Move(fyne.NewPos(1090, 75))

	// List element layout
	elements[5].Resize(fyne.NewSize(1200, 450))
	elements[5].Move(fyne.NewPos(0, 150))
}
