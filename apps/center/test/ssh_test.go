package test

import (
	"bufio"
	"config-manager/center/ssh"
	"fmt"
	"io"
	"testing"
	"time"
)

func TestSSH(t *testing.T) {

	// Replace with your actual SSH server details
	session, err := ssh.NewSession("127.0.0.1:22", "ducksoup", "969690")
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}
	defer session.Close()

	or, ow := io.Pipe()
	er, ew := io.Pipe()

	pipe, err := session.OpenShell(ow, ew)

	if err != nil {
		t.Fatal(err)
	}

	go func() {
		reader := bufio.NewReader(or)
		for {
			b, e := reader.ReadByte()
			if e != nil {
				break
			}
			fmt.Print(string(b))
		}
	}()

	go func() {
		reader := bufio.NewReader(er)
		for {
			b, e := reader.ReadByte()
			if e != nil {
				break
			}
			fmt.Print(string(b))
		}
	}()

	(*pipe).Write([]byte("ls -l\n"))

	timer := time.NewTimer(5 * time.Second)

	go func() {
		select {
		case <-timer.C:
			session.Close()
		}
	}()

	session.Wait()
}
