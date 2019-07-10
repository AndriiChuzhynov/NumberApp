package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net"
	"strconv"
)

func main() {
	var hostname = flag.String("target", "localhost:4000", "Target hostname:port")
	var maxConnections = flag.Int("connections", 5, "Number of connections")
	var termination = flag.Bool("term", false, "Termination will sent randomly")
	var corruptedData = flag.Bool("corr", false, "corrupted data will be sent randomly")
	var reporting = flag.Bool("rep", false, "will report to stdout")
	flag.Parse()

	guard := make(chan struct{}, *maxConnections)
	i := 0
	for {
		i += 1
		routineName := i
		guard <- struct{}{}
		go func() {
			fmt.Println("NewRoutineStarted")
			conn, err := net.Dial("tcp", *hostname)
			if err != nil {
				// handle error
				fmt.Printf("Error: %s\n", err.Error())
				return
			}

			for {
				var data []byte

				for i := 0; i < 9; i++ {
					r := rand.Int31n(9) + 48
					data = append(data, byte(r))
				}
				data = append(data, byte(10))

				if *termination {
					if rand.Intn(100000000) == 9999 {
						data = []byte("terminate\n")
					}
				}

				if *corruptedData {
					random := rand.Intn(10000000)
					if random == 9999 {
						data = []byte("asdasda\n")
					}
					if random == 1000 {
						data = []byte("12345678K\n")
					}
				}

				if *reporting {
					fmt.Print("Routine " + strconv.Itoa(routineName) + " " + string(data))
				}

				_, err = conn.Write(data)
				if err != nil {
					fmt.Printf("Error: %s\n", err.Error())
					break
				}
			}
			<-guard
		}()
	}
}
