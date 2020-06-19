package main

import (
   "fmt"
   "log"
   "net"
   "net/rpc"
)

type Listener int
type Reply struct {
   Data string
}

func (l *Listener) GetLine(line []byte, reply *Reply) error {
   rv := string(line)
   fmt.Printf("Receive: %v\n", rv)
   *reply = Reply{rv}
   return nil
}

func main() {
   fmt.Println("A")
   addy, err := net.ResolveTCPAddr("tcp", "0.0.0.0:12345")
   if err != nil {
      log.Fatal(err)
   }
   inbound, err := net.ListenTCP("tcp", addy)
   if err != nil {
      log.Fatal(err)
   }
   fmt.Println("B")
   listener := new(Listener)
   rpc.Register(listener)
   fmt.Println("C")
   rpc.Accept(inbound)
   fmt.Println("D")
   fmt.Println("Hello World")
}

/*func main(){
	c := 1


	fmt.Println("Hello World")
	fmt.Println(c)
}*/