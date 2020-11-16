package main

import (
	"fmt"
	"net"
	"encoding/gob"
	"time"
)

type Process struct {
	Id           int
	Value        uint64
	KeepRunning  bool
}

func (p *Process) start() {
	for {
		p.Value += 1
		fmt.Printf("\nid %d: %d", p.Id, p.Value)
		time.Sleep(time.Millisecond * 500)
	}
}

func main() {
	c, err := net.Dial("tcp", ":9999")
	if err != nil {
		panic(err)
	}
	var p Process
	err = gob.NewDecoder(c).Decode(&p)
	if err != nil {
		fmt.Println(err)
	}
	go p.start()
	var input string 
	fmt.Scanln(&input)
	err = gob.NewEncoder(c).Encode(p)
	if err != nil {
		fmt.Println(err)
	}
	c.Close()
}