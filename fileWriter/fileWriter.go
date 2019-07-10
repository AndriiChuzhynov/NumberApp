package fileWriter

import (
	"fmt"
	"os"
	"sync"
	"time"
)

const fileName = "numbers.log"

var wg sync.WaitGroup
var file *os.File
var inputChan chan []byte

func WriteData(data []byte) {
	inputChan <- data
}

func ShutDown() {
	fmt.Println("Flushing data..")

	for len(inputChan) > 0 {
	}

	close(inputChan)

	wg.Wait()
	fmt.Println("Flushed")
}

func fileWriter() {
	var buffer []byte
	ticker := time.NewTicker(100 * time.Millisecond)

	for {
		select {
		case <-ticker.C:
			flush(&buffer)
		case data, ok := <-inputChan:
			buffer = append(buffer, []byte(data)...)
			if !ok {
				gracefulShutdown(&buffer)
				return
			}
		}
	}
}

func InitWriter() {
	var err error
	file, err = os.Create(fileName)
	checkCriticalError(err)

	inputChan = make(chan []byte, 100000)
	go fileWriter()
	go queueWarning()
	wg.Add(1)

}

func gracefulShutdown(buffer *[]byte) {
	flush(buffer)

	err := file.Sync()
	if err != nil {
		fmt.Printf("Sync error: %s\n", err)
	}

	err = file.Close()
	if err != nil {
		fmt.Printf("Close file error: %s\n", err)
	}

	wg.Done()
}

func flush(buffer *[]byte) {
	_, err := file.Write(*buffer)
	if err != nil {
		fmt.Printf("Write error: %s\n", err)
	}
	*buffer = nil
}

func queueWarning() {
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				if len(inputChan) == cap(inputChan) {
					fmt.Printf("Writer queue is full: %d\n", len(inputChan))
				}
			}
		}
	}()
}

func checkCriticalError(e error) {
	if e != nil {
		panic(e)
	}
}
