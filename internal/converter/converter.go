package converter

import (
	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Task struct for a jobs channel
type Task struct {
	path string
}

// ModernizeDOCsToDOCXs converts all DOC files in the specified directory to DOCX format
func ModernizeDOCsToDOCXs(path string) error {
	// Declare number of concurrent workers
	const numberOfWorkers = 10
	var wg sync.WaitGroup
	jobs := make(chan Task, numberOfWorkers*5)

	// Spin up all workers
	for _ = range numberOfWorkers {
		wg.Add(1)
		go worker(jobs, &wg)
	}

	// Give them tasks
	err := filepath.WalkDir(path, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(d.Name(), ".doc") {
			jobs <- Task{path: path}
		}

		return nil
	})

	// Close the tasks chanel
	close(jobs)

	// Wait for them to finish
	wg.Wait()

	return err
}

func worker(jobs <-chan Task, wg *sync.WaitGroup) {
	defer wg.Done()

	// COM initialization with concurrency model
	err := ole.CoInitializeEx(0, ole.COINIT_MULTITHREADED)
	if err != nil {
		log.Printf("[ERROR]: failed to initialize COM module: %v\n", err)
		return
	}
	defer ole.CoUninitialize()

	// Create COM object to store Word Application
	comInit, err := oleutil.CreateObject("Word.Application")
	if err != nil {
		log.Printf("[ERROR]: failed to create COM object: %v\n", err)
		return
	}
	defer comInit.Release()

	// Create Word application instance
	word, err := comInit.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		log.Printf("[ERROR]: failed to get Word applicaion interface: %v\n", err)
		return
	}
	defer word.Release()

	// Set the Word application visibility
	_, err = oleutil.PutProperty(word, "Visible", false)
	if err != nil {
		log.Printf("[ERROR]: could not set Word app visibility to false: %v\n", err)
		return
	}

	// Run the jobs
	for task := range jobs {
		convertDocToDocx(word, task.path)
	}
}

func convertDocToDocx(word *ole.IDispatch, path string) {
	// Open the doc file
	doc, err := oleutil.CallMethod(word, "Document.Open", path)
	if err != nil {
		log.Printf("[ERROR]: could not open a doc file: %v\n", err)
		return
	}
	// Convert doc object to dispatch obj
	docDispatch := doc.ToIDispatch()
	defer docDispatch.Release()

	// Create path for a new docx file
	docxPath := strings.TrimSuffix(path, ".doc") + ".docx"

	// Save the doc as docx (format 16)
	_, err = oleutil.CallMethod(docDispatch, "SaveAs2", docxPath, 16)
	if err != nil {
		log.Printf("[ERROR]: could not save a document as DOCX [%s]:, %v\n", path, err)
		return
	}

	//Clean up
	//Close the DOC
	_, err = oleutil.CallMethod(docDispatch, "Close")
	if err != nil {
		log.Printf("[ERROR]: could not close the doc %v\n", err)
		return
	}

	// Delete the original doc file
	err = os.Remove(path)
	if err != nil {
		log.Printf("[WARNING]: could not delete original DOC file %s: %v\n", path, err)
		return
	}
}
