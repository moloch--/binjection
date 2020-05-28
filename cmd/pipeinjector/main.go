package main

import (
	"flag"
	"io/ioutil"

	"github.com/moloch--/binjection/bj"
)

func main() {
	pipeName := ""
	flag.StringVar(&pipeName, "p", `\\.\pipe\bdf`, "Pipe base name string")
	flag.Parse()

	go ListenPipeDry(pipeName + "dry")
	ListenPipeWet(pipeName + "wet")
}

func Inject(dry []byte) (wet []byte, err error) {
	config := &bj.BinjectConfig{CodeCaveMode: false}

	// *** Testing purposes
	shellcodeBytes, err := ioutil.ReadFile("./test.bin")
	if err != nil {
		return nil, err
	}
	return bj.Binject(dry, shellcodeBytes, config)
	// *** Testing purposes
	//return bj.Binject(dry, []byte{0, 0, 0, 0}, config)
}
