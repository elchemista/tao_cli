package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
)

func executeCommand(cmd string) {
	command := exec.Command("bash", "-c", cmd)
	var out bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &out
	command.Stderr = &stderr
	err := command.Run()
	if err != nil {
		log.Println(fmt.Sprint(err) + ": " + stderr.String())
		return
	}
	log.Println(out.String())
}
