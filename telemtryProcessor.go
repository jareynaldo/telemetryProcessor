package main
import (
	"fmt"
	"encoding/json"
	"sync"
	"net"
)

type createRules struct {
    Name 		string
    Label  		string
	Units 		string
    Expression  float32
	Units		string
}

type createQueue struct {
	t int 
	l string 
	u string 
	v int
}

func main (){

	config := createQueue{
		Name = "PSI to Bars"
		Label = "sensor:IPT"
		Units = "PSI"
		Expression = 0.0689476
		Units = "BARS"
	}	

	queue := make(chan crateQueue)  //make queue
	var wg sync.WaitGroup // make waitgroup so that finishes concurently


	// create Socket 
	listener, err := net.Listen("tcp", "127.0.0.1:9000")
    if err != nil {
        fmt.Println("Error listening:", err)
        return
    }
    defer listener.Close()

    fmt.Println("Server is listening on :9000")

    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("Error accepting connection:", err)
            continue
        }

        wg.Add(1)
		go handleSocket(queue, conn, &wg) 
    }

	wg.Add(1)
	go processorLoop(queue, &wg, config)

	wg.Wait()

}

func handleSocket(queue chan createQueue, conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	defer conn.Close()
	fmt.Println("Accepted connection from", conn.RemoteAddr())

	for {
		var data createQueue
		decoder := json.NewDecoder(conn)
		if err := decoder.Decode(&data); err != nil {
			fmt.Println("Error decoding JSON:", err)
			return
		}

	queue <- data // send entire struct to channel
	}
	
}

func processorLoop(queue chan currentMemeber, wg *sync.WaitGroup, createRules config){
	defer wg.Done() 

	for {
		popedQueue := <-queue // remove element from channel
		rule := config.Expression//what value do i need to multiply by

        poppedQueue.V = int(float32(poppedQueue.V) * rule)

        fmt.Printf("Processed: %+v\n", poppedQueue)
	}
}