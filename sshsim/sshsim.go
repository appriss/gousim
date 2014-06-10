package sshsim

import (
	// "fmt"
	"github.com/appriss/gousim"
	"code.google.com/p/go.crypto/ssh"
	// "net"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type SSHSimulation struct {
	Client *ssh.Client
	RecSession *gousim.RecordingSession
}

func Connect(host string, config *ssh.ClientConfig, rec *gousim.RecordingSession) (*SSHSimulation, error) {
	conn, err := rec.Dial("tcp", host + ":22")
	if err != nil {
		return nil, err
	}
	startClock := time.Now()
	c, chans, reqs, err := ssh.NewClientConn(conn, host, config)
    if err != nil {
    	return nil, err
    }
    client := ssh.NewClient(c, chans, reqs)
    rec.LogSample("AUTH", time.Now().Sub(startClock).Nanoseconds() )
    return &SSHSimulation{client, rec}, nil
}









