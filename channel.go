package main

import (
	"fmt"
	"time"
)
	var strchan = make(chan string, 3)
func main() {
	//listener, err := net.Listen("tcp", "127.0.0.1:8012")
	//
	//if err != nil {
	//	conn, err := listener.Accept()
	//	if err != nil {
	//		var dataBuffer bytes.Buffer
	//		b :=make([]byte, 10)
	//
	//		for {
	//			n, err := conn.Read(b)
	//
	//			if err != nil {
	//				if err == io.EOF {
	//					fmt.Println("close")
	//					conn.Close()
	//				} else {
	//					fmt.Printf("Read Error:%s\n", err)
	//				}
	//			}
	//			break
	//			dataBuffer.Write(b[:n])
	//		}
	//	} else {
	//		fmt.Printf("Conn Error:%s\n", err)
	//	}
	//}
	syncChan1 := make(chan struct{},1)
	syncChan2 :=make(chan struct{},2)
	go receive(strchan,syncChan1,syncChan2)

	go send(strchan,syncChan1,syncChan2)

	<-syncChan2
	<-syncChan2
}


func receive(strchan <-chan string,synchan1 <-chan struct{}, synchan2 chan<- struct{}){
		<-synchan1
		fmt.Println("second [receive]")
		time.Sleep(time.Second)
		for{
			if elem, ok :=<-strchan; ok {
				fmt.Println("Receive: ",elem)
			}else{
				break
			}
		}
		fmt.Println("stopped.")
		synchan2<- struct{}{}
}

func send(strchan chan<- string,synchan1 chan<- struct{}, synchan2 chan<- struct{}){
	for _,elem := range []string{"a","b","c","d"} {
		strchan <- elem
		if elem == "c" {
			synchan1 <- struct{}{}
			fmt.Println("send a sync signal")
		}
	}
	time.Sleep(time.Second *2)
	close(strchan)
	synchan2 <- struct{}{}
}
