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

	// Create list element
	list := widget.NewList(
		func() int {
			return 20
		}, func() fyne.CanvasObject {
			return widget.NewLabel("42\nlolo")
		}, func(id widget.ListItemID, object fyne.CanvasObject) {
			object = widget.NewLabel("Here is a list number")
		})

	// Create a container to hold the button and label
	content := container.New(
		new(customLayout),
		openFolderButton,
		folderLabel,
		inputBox1,
		inputBox2,
		searchBtn,
		list,
	)

	// Set the window content and show the window
	window.SetContent(content)
	window.ShowAndRun()
}
