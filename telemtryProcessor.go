package main
import (
	"fmt"
	"encoding/json"
	"sync"
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
	var wg sync.WaitGroup // make waitgroup so that finished concurently
*-

	wg.Add(1)
	go handleSocket(queue, &wg) 


	wg.Add(1)
	go processorLoop(queue, &wg, config)

	wg.Wait()

}

// data shown like (0x09){"t:"t"}
func handleSocket(queue chan createQueue, wg *sync.WaitGroup){
	defer wg.Done() 

	for {
		
		length := socket.readUNIT32()
		bytes := socket.readBytes(length)
		value := json.Unmarshal(bytes) 

		newMemeber := createQueue{t: value[time], l : value[lable], u: value[units], v: value[value]}
		queue <- newMemeber;
	}

}

func processorLoop(queue chan currentMemeber, wg *sync.WaitGroup, createRules config){

	defer wg.Done() 

	for {
		popedQueue := queue.pop() // remove element from channel
		rule := config(popedQueue.lable)//what value do i need to multiply by
		queueValue = popedQueue[v] // second multiply value

		popedQueue[v] := queueValue * rule 

		print(popedQueue)
	}
}