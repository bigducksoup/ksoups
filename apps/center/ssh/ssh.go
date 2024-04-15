package ssh

import (
	"golang.org/x/crypto/ssh"
	"io"
	"time"
)

type SSHSession struct {
	Config     *ssh.ClientConfig
	Session    *ssh.Session
	StdinPipe  *io.WriteCloser
	StdoutPipe *io.PipeReader
	StderrPipe *io.PipeReader
}

func (s *SSHSession) OpenShell() (err error) {

	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	err = s.Session.RequestPty("xterm", 40, 80, modes)
	if err != nil {
		return err
	}
	err = s.Session.Shell()
	if err != nil {
		return err
	}
	return nil
}

func (s *SSHSession) Wait() error {
	return s.Session.Wait()
}

func (s *SSHSession) Close() error {
	(*s.StdinPipe).Close()
	return s.Session.Close()
}

// NewSession creates a new SSHSession session
// address: the address of the SSHSession server
// username: the username to use
// password: the password to use
// returns: a new SSHSession session
// returns: an error if one occurs
func NewSession(address string, username string, password string) (*SSHSession, error) {

	config := ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	client, err := ssh.Dial("tcp", address, &config)

	if err != nil {
		return nil, err
	}

	session, err := client.NewSession()

	if err != nil {
		return nil, err
	}

	pipe, err := session.StdinPipe()

	if err != nil {
		return nil, err
	}

	stdoutR, stdoutW := io.Pipe()
	session.Stdout = stdoutW

	stderrR, stderrW := io.Pipe()
	session.Stderr = stderrW

	s := &SSHSession{
		Config:     &config,
		Session:    session,
		StdinPipe:  &pipe,
		StdoutPipe: stdoutR,
		StderrPipe: stderrR,
	}

	return s, nil
}
