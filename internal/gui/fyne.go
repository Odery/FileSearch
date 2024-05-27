package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/Odery/FileSearch/internal/docx"
	"log"
)

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
}

func DrawGUI() {
	// Create a new application
	gui := app.New()
	window := gui.NewWindow("File Search - Степан Помічник")

	// Set default settings and size
	window.Resize(fyne.NewSize(1200, 600))
	window.CenterOnScreen()
	window.SetFixedSize(true)

	// Declare URI variable
	var workingDIR fyne.ListableURI

	// Create a label to display the selected folder path
	folderLabel := widget.NewLabel("Не вибрана")

	// Create a button to open the folder selection dialog
	openFolderButton := widget.NewButton("Робоча директорія:", func() {
		// Open the folder selection dialog
		dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			// Handle error (if any)
			if err != nil {
				log.Println("[ERROR]: ", err)
				return
			}
			// Handle case where no folder was selected
			if uri == nil && workingDIR == nil {
				folderLabel.SetText("Не вибрана")
			}

			// Update the label with the selected folder path
			if uri != nil {
				workingDIR = uri
				folderLabel.SetText(workingDIR.Path())
			}
		}, window).Show()
	})

	// Create input boxes for search
	inputBox1 := widget.NewEntry()
	inputBox1.SetPlaceHolder("Пошуковий запит")
	inputBox2 := widget.NewEntry()
	inputBox2.SetPlaceHolder("Додатковий параметр")

	// Create search button
	searchBtn := widget.NewButton("Шукати", func() {
		//TODO make use of wg (wait group) to make a progressbar
		result, _, err := docx.ProcessSearchRequest(workingDIR.Path(), inputBox1.Text, inputBox2.Text)
		if err != nil {
			log.Println("[ERROR]: ", err)
		}

		//TODO
		result.Lock()
		result.Unlock()
	})

	// Create a container to hold the button and label
	content := container.New(
		new(customLayout),
		openFolderButton,
		folderLabel,
		inputBox1,
		inputBox2,
		searchBtn,
	)

	// Set the window content and show the window
	window.SetContent(content)
	window.ShowAndRun()
}
