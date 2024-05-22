package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"log"
)

func DrawGUI() {
	// Create a new application
	gui := app.New()
	window := gui.NewWindow("File Search - BANDERA")

	// Set default settings and size
	window.Resize(fyne.NewSize(1200, 600))
	window.CenterOnScreen()
	window.SetFixedSize(true)

	// Create a label to display the selected folder path
	folderLabel := widget.NewLabel("!No folder selected!")

	// Create a button to open the folder selection dialog
	openFolderButton := widget.NewButton("Робоча директорія", func() {
		// Open the folder selection dialog
		dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
			// Handle error (if any)
			if err != nil {
				log.Println("[ERROR]: ", err)
				return
			}
			// Handle case where no folder was selected
			if uri == nil {
				folderLabel.SetText("!No folder selected!")
			}

			// Update the label with the selected folder path
			if uri != nil {
				folderLabel.SetText("Робоча директорія: " + uri.Path())
			}
		}, window).Show()
	})

	// Create a container to hold the button and label
	content := container.NewVBox(
		openFolderButton,
		folderLabel,
	)

	// Set the window content and show the window
	window.SetContent(content)
	window.ShowAndRun()
}
