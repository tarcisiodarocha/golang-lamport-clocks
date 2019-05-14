package main

import (
	"fmt"
	"time"
)

type Message struct {
	Body      string
	Timestamp int
}

func event(pid, counter int) int {
	counter += 1
	fmt.Printf("Event in process pid=%v. Counter=%v\n", pid, counter)
	return counter
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func calcTimestamp(recvTimestamp, counter int) int {
	return max(recvTimestamp, counter) + 1
}

func sendMessage(ch chan Message, pid int, counter int) int {
	counter += 1
	ch <- Message{"Test msg!!!", counter}
	fmt.Printf("Message sent from pid=%v. Counter=%v\n", pid, counter)
	return counter

}

func receiveMessage(ch chan Message, pid int, counter int) int {
	message := <-ch
	counter = calcTimestamp(message.Timestamp, counter)
	fmt.Printf("Message received at pid=%v. Counter=%v\n", pid, counter)
	return counter
}

func processOne(ch12, ch21 chan Message) {
	pid := 1
	counter := 0
	counter = event(pid, counter)
	counter = sendMessage(ch12, pid, counter)
	counter = event(pid, counter)
	counter = receiveMessage(ch21, pid, counter)
	counter = event(pid, counter)

}

func processTwo(ch12, ch21, ch23, ch32 chan Message) {
	pid := 2
	counter := 0
	counter = receiveMessage(ch12, pid, counter)
	counter = sendMessage(ch21, pid, counter)
	counter = sendMessage(ch23, pid, counter)
	counter = receiveMessage(ch32, pid, counter)

}

func processThree(ch23, ch32 chan Message) {
	pid := 3
	counter := 0
	counter = receiveMessage(ch23, pid, counter)
	counter = sendMessage(ch32, pid, counter)

}

func main() {
	oneTwo := make(chan Message, 100)
	twoOne := make(chan Message, 100)
	twoThree := make(chan Message, 100)
	threeTwo := make(chan Message, 100)

	go processOne(oneTwo, twoOne)
	go processTwo(oneTwo, twoOne, twoThree, threeTwo)
	go processThree(twoThree, threeTwo)

	time.Sleep(5 * time.Second)
}
