package main

import (
	"net"
	"fmt"
	"encoding/gob"
	"time"
	"sort"
)

type Process struct {
	Id           int
	Value        uint64
	KeepRunning  bool
}

func (p *Process) start() {
	for {
		if p.KeepRunning == false {
			break
		}
		p.Value += 1
		fmt.Printf("\nid %d: %d", p.Id, p.Value)
		time.Sleep(time.Millisecond * 500)
	}
}

func generateNProcess(n int) []*Process {
	var processes []*Process
	for i := 0; i < n ; i++ {
		p := Process{Id:i, Value: 0, KeepRunning: true}
		go p.start()
		processes = append(processes, &p)
	}
	return processes
}

func getProcess(processes []*Process) *Process {
	sort.Slice(processes[:], func(i, j int) bool {
		return processes[i].Id < processes[j].Id
	})
	p := processes[0]
	copy(processes[0:], processes[1:])
	processes = processes[:len(processes)-1]
	processes = processes[0:]
	p.KeepRunning = false
	return p
}

func server(processes []*Process) {
	s, err := net.Listen("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		conn, err := s.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		if len(processes) > 0 {
			go handleClient(conn, processes)
		}
	}
}

func handleClient(conn net.Conn, processes []*Process) {
	// Write
	sort.Slice(processes[:], func(i, j int) bool {
		return processes[i].Id < processes[j].Id
	})
	p := processes[0]
	copy(processes[0:], processes[1:])
	processes = processes[:len(processes)-1]
	processes = processes[0:]
	p.KeepRunning = false
	err := gob.NewEncoder(conn).Encode(p)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Read
	var decodedProcess Process
	err = gob.NewDecoder(conn).Decode(&decodedProcess)
	if err != nil {
		fmt.Println(err)
		return 
	}
	if &decodedProcess != nil {
		decodedProcess.KeepRunning = true
		go decodedProcess.start()
		processes = append(processes, &decodedProcess)
		sort.Slice(processes[:], func(i, j int) bool {
			return processes[i].Id < processes[j].Id
		})
	}
}

func main() {
	go server(generateNProcess(5))
	var input string 
	fmt.Scanln(&input)
}


