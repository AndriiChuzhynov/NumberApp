package processor

import (
	"NumberApp/fileWriter"
	"NumberApp/reporting"
	"fmt"
	"sync"
	"time"
)

var recordedValues [1000000000]bool
var numbersChannel = make(chan messageToWrite, 1000000)
var wgProcessor sync.WaitGroup

type messageToWrite struct {
	integer int
	bytes   []byte
}

func CloseMessagesProcessor() {
	close(numbersChannel)
	wgProcessor.Wait()
}

func MessagesProcessor() {
	go queueWarning()
	wgProcessor.Add(1)
	ticker := time.NewTicker(10 * time.Second)

	for message := range numbersChannel {
		select {
		case <-ticker.C:
			reporting.PrintReport()
		default:
			if recordedValues[message.integer] {
				reporting.Duplicated()
			} else {
				fileWriter.WriteData(message.bytes)
				reporting.Uniq()
				recordedValues[message.integer] = true
			}
		}
	}
	wgProcessor.Done()
}

func AddMessageToQueue(integer int, bytes []byte) {
	numbersChannel <- messageToWrite{integer, bytes}
}

func queueWarning() {
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				if len(numbersChannel) >= cap(numbersChannel)-1000 {
					fmt.Printf("Messages queue is almost full %d\n", len(numbersChannel))
				}
			}
		}
	}()
}
