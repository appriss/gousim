package sshsim

import (
	// "fmt"
	"github.com/appriss/gousim"
	"code.google.com/p/go.crypto/ssh"
	"net"
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
	ip := net.ParseIP(host)
	if ip == nil {
		startClock := time.Now()
		ips, err := net.LookupIP(host)
		if err != nil {
			return nil, err
		}
		rec.LogSample("DNS", time.Now().Sub(startClock).Nanoseconds() )
		ip = ips[rand.Intn(len(ips))]
	}
	startClock := time.Now()
	conn, err := net.Dial("tcp", ip.String() + ":22")
	if err != nil {
		return nil, err
	}
	rec.LogSample("TCP", time.Now().Sub(startClock).Nanoseconds() )
	startClock = time.Now()
	c, chans, reqs, err := ssh.NewClientConn(conn, host, config)
    if err != nil {
    	return nil, err
    }
    client := ssh.NewClient(c, chans, reqs)
    rec.LogSample("AUTH", time.Now().Sub(startClock).Nanoseconds() )
    return &SSHSimulation{client, rec}, nil
}






