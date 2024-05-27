package main

import (
	"github.com/Odery/FileSearch/internal/gui"
	"log"
	"os"
)

func main() {
	logFile, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("[FATAL]: ", err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

	log.Println("[INFO]: APP STARTED")
	gui.DrawGUI()
	log.Println("[INFO]: APP EXITED")
}
