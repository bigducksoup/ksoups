package test

import (
	"apps/center/ssh"
	"log"
	"testing"
)

func TestSSH(t *testing.T) {

	// Replace with your actual SSH server details
	session, err := ssh.NewSession("127.0.0.1:22", "ducksoup", "969690")
	if err != nil {
		t.Fatalf("Failed to create session: %v", err)
	}
	defer session.Close()

	session.OpenShell()

	bytes := make([]byte, 1024)
	n, err := session.StdoutPipe.Read(bytes)

	log.Println(string(bytes[:n]))

}
