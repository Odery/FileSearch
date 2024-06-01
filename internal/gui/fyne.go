package gui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/Odery/FileSearch/internal/docx"
	"log"
	"os/exec"
	"time"
)

func DrawGUI() {
	// Create a new application
	gui := app.New()
	window := gui.NewWindow("File Search - Помічник")

	// Set default settings, size and theme
	window.Resize(fyne.NewSize(1200, 500))
	window.CenterOnScreen()
	window.SetFixedSize(true)

	// Declare URI variable
	var workingDIR fyne.ListableURI

	// Create a label to display the selected folder path
	folderLabel := widget.NewLabel("Не вибрана")

	// Create a result variable to store search result
	result := docx.NewSearchResult()

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

	// Create an input boxes for search
	inputBox1 := widget.NewEntry()
	inputBox1.SetPlaceHolder("Пошуковий запит")
	inputBox2 := widget.NewEntry()
	inputBox2.SetPlaceHolder("Додатковий параметр")

	// Create a progress bar
	progressBar := widget.NewProgressBar()
	progressBar.Hide()

	// Create a table element
	// probably the most difficult part to understand
	table := widget.NewTableWithHeaders(
		func() (rows int, cols int) { return len(result.Results), 3 },
		func() fyne.CanvasObject { return NewTappableLabel() },
		func(id widget.TableCellID, object fyne.CanvasObject) {
			label := object.(*TappableLabel)
			label.Wrapping = fyne.TextTruncate
			// *! Only for Windows!
			label.OnDoubleTapped = func() {
				log.Println("[INFO]: Opening a file: ", result.Results[id.Row].Path)
				cmd := exec.Command("rundll32", "url.dll,FileProtocolHandler", result.Results[id.Row].Path)
				err := cmd.Start()
				if err != nil {
					log.Println("[ERROR]: When opening a file. ", err)
				}
			}
			if id.Col == 0 {
				label.SetText(result.Results[id.Row].Name)
			} else if id.Col == 1 {
				label.SetText(result.Results[id.Row].GetFormattedDate())
			} else if id.Col == 2 {
				label.SetText(result.Results[id.Row].Path)
			}
		})

	// Updating table headers and making them clickable
	table.CreateHeader = func() fyne.CanvasObject {
		return NewTappableLabel()
	}

	table.UpdateHeader = func(id widget.TableCellID, object fyne.CanvasObject) {
		obj := object.(*TappableLabel)
		obj.Wrapping = fyne.TextTruncate
		obj.Alignment = fyne.TextAlignCenter
		if id.Row == -1 && id.Col == 0 {
			obj.SetText("Ім'я")
			obj.OnDoubleTapped = func() {
				result.SortByNameAscending()
			}
		} else if id.Row == -1 && id.Col == 1 {
			obj.SetText("Дата змінення")
		} else if id.Row == -1 && id.Col == 2 {
			obj.SetText("Шлях до файлу")
		}

		if id.Col == -1 {
			obj.SetText(fmt.Sprint(id.Row + 1))

		}

	}

	// Adjust table default Column size
	table.SetColumnWidth(0, 300)
	table.SetColumnWidth(1, 148)
	table.SetColumnWidth(2, 300)

	// Create a search button
	searchBtn := widget.NewButton("Шукати", func() {
		if workingDIR == nil {
			dialog.ShowInformation("Робоча директорія не вибрана!", "Будьласка виберіть робочу директорію!", window)
			return
		}
		if inputBox1.Text == "" {
			dialog.ShowInformation("Увага!", "Основний пошуковий запит не може бути пустим!", window)
			return
		}
		progressBar.Show()

		go func() {
			r, progress, err := docx.ProcessSearchRequest(workingDIR.Path(), inputBox1.Text, inputBox2.Text)
			if err != nil {
				log.Println("[ERROR]: in docx part: ", err)
			}
			result = r

			lastDone := 0
			for !progress.IsDone() {
				if lastDone < progress.GetDone() {
					table.Refresh()
					progressBar.SetValue(progress.GetProgress())
				}
			}
			time.Sleep(1 * time.Second)
			progressBar.Hide()
		}()
	})

	// Create a container to hold the button and label
	content := container.New(
		new(customLayout),
		openFolderButton,
		folderLabel,
		inputBox1,
		inputBox2,
		searchBtn,
		progressBar,
		table,
	)

	// Set the window content and show the window
	window.SetContent(content)
	window.ShowAndRun()
}
