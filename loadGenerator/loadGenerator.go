package main

import (
	"flag"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"sync"
)

func main() {
	var hostname = flag.String("target", "localhost:8888", "Target hostname:port")
	var maxConnections = flag.Int("conn", 5, "Number of connections")
	var termination = flag.Bool("term", false, "Termination will sent randomly")
	var corruptedData = flag.Bool("corr", false, "corrupted data will be sent randomly")
	var logging = flag.Bool("logg", false, "will report to stdout")
	flag.Parse()
	var wg sync.WaitGroup

	guard := make(chan struct{}, *maxConnections)
	counter := make(chan uint64, *maxConnections)
	routine := 0
	for {
		routine += 1
		routineName := routine
		guard <- struct{}{}
		fmt.Println("NewRoutineStarted")

		conn, err := net.Dial("tcp", *hostname)

		if err != nil {
			fmt.Printf("Connection error: %s\n", err.Error())
			break
		}

		wg.Add(1)
		go func() {
			var sent uint64 = 0
			for {
				var data []byte

				for i := 0; i < 9; i++ {
					r := rand.Int31n(9) + 48
					data = append(data, byte(r))
				}
				data = append(data, byte(10))

				if *termination {
					if rand.Intn(7050000) == 9999 {
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

				if *logging {
					fmt.Print("Routine " + strconv.Itoa(routineName) + " " + string(data))
				}

				_, err = conn.Write(data)
				if err != nil {
					fmt.Printf("Send error: %s\n", err.Error())
					break
				}
				sent++
			}
			fmt.Println("RoutineStopped")
			counter <- sent
			wg.Done()
			<-guard
		}()
	}
	wg.Wait()
	close(counter)
	var sum uint64
	for sent := range counter {
		sum += sent
	}
	fmt.Printf("TotalSent: %d\n", sum)
}
