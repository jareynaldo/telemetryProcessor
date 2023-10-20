package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

type createRules struct {
    Name 		string
    Label  		string
	Units 		string
    Expression  float32
}


type createQueue struct {
	T int 		`json:"t"`
	L string 	`json:"l"`
	U string    `json:"u"`
	V int		`json:"v"`
}

func main (){


	config := createRules {
		Name:		 "PSI to Bars",
		Label:		 "sensor:IPT",
		Units:		 "PSI",
		Expression:	 0.0689476,
	}	

	queue := make(chan createQueue)  //make queue


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

		go handleSocket(queue, conn) 
    }

	go processorLoop(queue, config)


}

func handleSocket(queue chan createQueue, conn net.Conn) {
	defer conn.Close()
	hasHandshaked := false
	fmt.Println("Accepted connection from", conn.RemoteAddr())

	for {
		integerBuf := make([]byte, 4)
		_, err := conn.Read(integerBuf)
		if err != nil {
			return
		}
		length := binary.BigEndian.Uint32(integerBuf)

		dataBuf := make([]byte, length)
		_, err = conn.Read(dataBuf)
		if err != nil {
			return
		}

		if string(dataBuf) == "START" {
			integerBuf = make([]byte, 4)
			integerBuf = binary.BigEndian.AppendUint32(integerBuf, length)
			dataBuf = []byte("OK")
			conn.Write(integerBuf)
			conn.Write(dataBuf)
			hasHandshaked = true

		} else if (hasHandshaked){
			var sending createQueue
			err  = json.Unmarshal(dataBuf, &sending)
			if err != nil {
				return
			}
			queue <- sending
		}
	}
	
}

func processorLoop(queue chan createQueue, config createRules){

	for {
		poppedQueue := <-queue // remove element from channel
		rule := config.Expression //what value do i need to multiply by

        poppedQueue.V = int(float32(poppedQueue.V) * rule)

        fmt.Printf("Processed: %+v\n", poppedQueue)
	}
}