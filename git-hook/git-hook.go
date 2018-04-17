package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

func main() {
	cmdInit()
	http.HandleFunc("/deploy", func(w http.ResponseWriter, r *http.Request) {
		cmdPull()
		fmt.Fprintf(w, "Deployed!")
		cmdStartServer()
	})
	if err := http.ListenAndServe(":1234", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func cmdInit() {
	cmd := exec.Command("pkill", "cockroach")
	_, _ = cmd.Output()
	go launchDB()
}

func launchDB() {
	cmd := exec.Command("task", "database")
	_, err := cmd.Output()
	perror(err)
}

func cmdPull() {
	fmt.Println("Deploying")
	cmd := exec.Command("git", "checkout", ".")
	cmd.Dir = "/root/suncork"
	out, err := cmd.Output()
	cmd = exec.Command("git", "pull", "origin", "master")
	cmd.Dir = "/root/suncork"
	out, err = cmd.Output()
	perror(err)
	fmt.Println(string(out))
}

func cmdStartServer() {
	c := exec.Command("pkill", "suncork")
	_, _ = c.Output()
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd := exec.Command("task", "run")

	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()

	var errStdout, errStderr error
	stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
	stderr := io.MultiWriter(os.Stderr, &stderrBuf)
	err := cmd.Start()
	if err != nil {
		log.Fatalf("cmd.Start() failed with '%s'\n", err)
	}

	go func() {
		_, errStdout = io.Copy(stdout, stdoutIn)
	}()

	go func() {
		_, errStderr = io.Copy(stderr, stderrIn)
	}()

	err = cmd.Wait()
	if err != nil {
		panic(err)
	}
	// if errStdout != nil || errStderr != nil {
	// 	log.Fatal("failed to capture stdout or stderr\n")
	// }
	outStr, errStr := string(stdoutBuf.Bytes()), string(stderrBuf.Bytes())
	fmt.Printf("\nout:\n%s\nerr:\n%s\n", outStr, errStr)
}

func perror(err error) {
	if err != nil {
		panic(err)
	}
}
