// +build windows

package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"net"

	"github.com/moloch--/binjection/bj"
	npipe "gopkg.in/natefinch/npipe.v2"
)

func MakePipe(pipename string) string {
	return `\\.\pipe\` + pipename

}

func ListenPipeDry(pipename string, config *bj.BinjectConfig) {
	ln, err := npipe.Listen(pipename)
	if err != nil {
		log.Fatalf("Listen(%s) failed: %v", pipename, err)
	}

	for {
		conn, err := ln.Accept()
		if err == npipe.ErrClosed {
			return
		}
		if err != nil {
			log.Fatalf("Error accepting connection: %v", err)
		}
		go handleDryConnection(conn, config)
	}
}

func ListenPipeWet(pipename string) {
	ln, err := npipe.Listen(pipename)
	if err != nil {
		log.Fatalf("Listen(%s) failed: %v", pipename, err)
	}

	for {
		conn, err := ln.Accept()
		if err == npipe.ErrClosed {
			return
		}
		if err != nil {
			log.Fatalf("Error accepting connection: %v", err)
		}
		go handleWetConnection(conn)
	}
}

var lastBytes []byte

func handleDryConnection(conn net.Conn, config *bj.BinjectConfig) {
	r := bufio.NewReader(conn)
	b, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatalf("Error reading from connection: %v", err)
	}

	i, err := Inject(b, config)
	if err != nil {
		log.Fatalf("Error injecting: %v", err)
	}
	log.Println("Set lastBytes: ", len(lastBytes))
	lastBytes = i
	if err := conn.Close(); err != nil {
		log.Fatalf("Error closing server side of connection: %v", err)
	}
}

func handleWetConnection(conn net.Conn) {
	w := bufio.NewWriter(conn)
	_, err := w.Write(lastBytes)

	log.Println("Wrote wet bytes: ", len(lastBytes))

	if err != nil {
		log.Fatalf("Error on writing to pipe: %v", err)
	}

	if err := conn.Close(); err != nil {
		log.Fatalf("Error closing server side of connection: %v", err)
	}
	lastBytes = nil
}
