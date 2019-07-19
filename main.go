package main

import (
	"NumberApp/fileWriter"
	"NumberApp/processor"
	"NumberApp/reporting"
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
)

const maxConnections = 5
const dataLength = 9
const terminationSignalMessage = "terminate"

//const terminationSignalLen = len(terminationSignalMessage)

var wg sync.WaitGroup
var gracefulStop = make(chan bool, maxConnections)

func main() {
	ln, err := net.Listen("tcp", ":4000")
	fmt.Println("Listening localhost 4000")
	if err != nil {
		panic(err)
	}

	go processor.MessagesProcessor()
	fileWriter.InitWriter()

	run(ln)
	wg.Wait()

	handleTermination(ln)
}

func run(ln net.Listener) {
	limit := make(chan struct{}, maxConnections)

	for {
		limit <- struct{}{}
		select {
		case <-gracefulStop:
			gracefulStop <- true
			return
		default:
		}

		conn, err := ln.Accept()

		if err != nil {
			fmt.Println(err)
			return
		}

		go func() {
			wg.Add(1)
			handleConnection(conn, gracefulStop)
			<-limit
			wg.Done()
		}()
	}
}

func handleConnection(connection net.Conn, gracefulStop chan bool) {
	fmt.Println("Started a new routine")

	reader := bufio.NewReader(connection)
	for {
		message, err := reader.ReadBytes('\n')

		if err != nil {
			fmt.Printf("Network event happend: %s\n", err.Error())
			break
		}

		err = checkFormat(message)
		if err != nil {
			fmt.Printf("Invalid format: %s\n", err.Error())
			break
		}

		i, err := convertToInt(message)
		if err != nil {
			if isTerminationSignal(message) {
				gracefulStop <- true
				fmt.Println("Termination signal received")
				_ = connection.Close()
				return
			}
			fmt.Printf("Invalid format: %s\n", err.Error())
			break
		}
		processor.AddMessageToQueue(i, message)

		select {
		case <-gracefulStop:
			gracefulStop <- true
			fmt.Println("Graceful stop routine")
			_ = connection.Close()
			return
		default:
		}
	}
	_ = connection.Close()
	fmt.Println("Routine closed")
}

func convertToInt(message []byte) (int, error) {
	s := string(message[0 : len(message)-1])
	return strconv.Atoi(s)
}

func checkFormat(message []byte) error {
	if len(message) != dataLength+1 {
		return fmt.Errorf("length should be %d, line %s", dataLength, message)
	}
	return nil
}

func isTerminationSignal(message []byte) bool {
	return string(message[0:len(message)-1]) == terminationSignalMessage
}

func handleTermination(listener net.Listener) {
	fmt.Println("Terminating")
	_ = listener.Close()
	processor.CloseMessagesProcessor()
	fileWriter.ShutDown()
	reporting.PrintReport()
	os.Exit(0)
}
