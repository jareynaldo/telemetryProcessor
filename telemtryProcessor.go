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
		Label:		 "sensor:MPT",
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
	go processorLoop(queue, config)

    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("Error accepting connection:", err)
            continue
        }

		go handleSocket(queue, conn) 		
    }



}

func handleSocket(queue chan createQueue, conn net.Conn) {
	hasHandshaked := false
	fmt.Println("Accepted connection from", conn.RemoteAddr())

	for {
		integerBuf := make([]byte, 4)
		_, err := conn.Read(integerBuf)

		if err != nil{
			//fmt.Printf("Error happened: %v \n", err)
			return
		}

		length := binary.BigEndian.Uint32(integerBuf)
		fmt.Print("Continues through loop \n")

		dataBuf := make([]byte, length)
		_, err = conn.Read(dataBuf)
		if err != nil {
			fmt.Printf("Equals Nill PT2: %v\n", err)
			return
		}

		if string(dataBuf) == "START" {
			integerBuf = make([]byte, 4)
			integerBuf = binary.BigEndian.AppendUint32(integerBuf, length)
			dataBuf = []byte("OK")
			conn.Write(integerBuf)
			conn.Write(dataBuf)
		}
		if string(dataBuf) == "OK" {
			// Handle the "OK" response
			fmt.Println("Received OK response.")
			hasHandshaked = true
			fmt.Printf("Has the handshake happened: %v \n", hasHandshaked)
		}  
		if hasHandshaked {
			var sending createQueue			
			err := json.Unmarshal(dataBuf, &sending)
			if err == nil {
				fmt.Printf("Error unmarshalling JSON data: %v\n", err)
			} 
			fmt.Printf("Data is being added to queue: %+v\n", sending)
			queue <- sending
		}

		
		// 	hasHandshaked = true
		// 	fmt.Printf("Has the handshake happened: %v \n", hasHandshaked)
		
		// if hasHandshaked {
		// 	var sending createQueue
		// 	fmt.Print("Starting queue upending \n")

		// 	if string(dataBuf) == "OK" {
		// 		err = json.Unmarshal(dataBuf, &sending)
		// 		if err != nil {
		// 			if err.Error() != "invalid character 'O' looking for beginning of value" {
		// 				// Handle the specific error case where 'O' is not valid JSON data.
		// 				fmt.Println("Had another issue")
		// 				return
		// 			} 
		// 		}
		// 		fmt.Printf("data is being added to queue:  %v\n", sending.V)
		// 		queue <- sending
		// 	}
		// }

	}
	
}

func processorLoop(queue chan createQueue, config createRules){
	fmt.Print("Starting processing \n")

	for {
		poppedQueue := <-queue // remove element from channel
		rule := config.Expression //what value do i need to multiply by

        poppedQueue.V = int(float32(poppedQueue.V) * rule)

        fmt.Printf("Processed: %+v\n", poppedQueue)
	}
}