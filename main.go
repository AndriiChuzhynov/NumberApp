package main

import (
	"NRApplication/btree"
	"bufio"
	"errors"
	"fmt"
	"net"
)

const maxConnections = 2

var tree = btree.New(5)

func main() {

	ln, err := net.Listen("tcp", ":4000")
	if err != nil {
		// handle error
	}
	fmt.Println("Listening localhost 4000")

	guard := make(chan struct{}, maxConnections)
	for {

		guard <- struct{}{}
		conn, err := ln.Accept()
		if err != nil {
			// handle error
		}
		go func() {
			//todo limit to max 5, add config
			handleConnection(conn)
			<-guard
		}()
	}
}

func handleConnection(connection net.Conn) {
	fmt.Println("Started new routine")
	for {
		//todo check if possible use win and unix term sequence
		//todo limit input max 10 symbols, add config
		message, err := bufio.NewReader(connection).ReadString('\n')
		if err != nil {
			fmt.Printf("Happend: %s\n", err.Error())
			break
		}

		//remove line spearator
		message = message[0 : len(message)-1]

		err = checkFormat(&message)
		if err != nil {
			fmt.Printf("Invalid format: %s\n", err.Error())
			break
		}

		fmt.Print("Message Received:", message)
	}
	_ = connection.Close()
	fmt.Println("Routine closed")
}

func checkFormat(str *string) error {

	if len(*str) != 9 {
		return errors.New("length should be 9")
	}

	//value, err := strconv.Atoi(*str)
	//if err != nil {
	//	return err
	//}

	//tree.ReplaceOrInsert(value)
	return nil
}
