package sshsim

import (
	"github.com/pkg/sftp"
	"time"
	"math/rand"
	"encoding/binary"
	"io"
)

type SFTPSimulator struct {
	SSHSim *SSHSimulation
	Client *sftp.Client
}

func NewSFTPSim(sim *SSHSimulation) (*SFTPSimulator, error) {
	c, err := sftp.NewClient(sim.Client)
	if err != nil {
		return nil, err
	}
	return &SFTPSimulator{sim, c}, nil
}

func (sim *SFTPSimulator) CloseAll() {
	sim.Client.Close()
	sim.SSHSim.Client.Close()
}

func (sim *SFTPSimulator) WalkDir(dirPath string) error {
	startClock := time.Now()
	walker := sim.Client.Walk(dirPath)
	for walker.Step() {
		err := walker.Err()
		if err != nil {
			return err
		}
	}
	sim.SSHSim.RecSession.LogSample("WALKDIR", time.Now().Sub(startClock).Nanoseconds())
	return nil
}

func (sim *SFTPSimulator) PutRandomFile( fqdn string, size int64, buf int64) error {
	startClock := time.Now()
	f, err := sim.Client.Create(fqdn)
	defer f.Close()
	if err != nil {
		return err
	}
	b, err := io.Copy(f, NewRandomBufferedData(size))
	if err != nil {
		return err
	}
	stopClock := time.Now()
	sim.SSHSim.RecSession.LogSample("UPLOAD", stopClock.Sub(startClock).Nanoseconds())
	sim.SSHSim.RecSession.LogSample("UPRATE", b / (stopClock.Sub(startClock).Nanoseconds() / 1000 / 1000))
	return nil
}

type RandomBufferedData struct {
	maxsize int64
	seekptr int64
	prng *rand.Rand
	buf []byte
}

func NewRandomBufferedData(maxsize int64) *RandomBufferedData {
	prng := rand.New(rand.NewSource(time.Now().UnixNano()))
	buf := make([]byte, 4)
	return &RandomBufferedData{maxsize, 0, prng, buf}
}

//Read always fills buf p with random data up to the max size, then returns EOF
func ( r *RandomBufferedData ) Read(p []byte) (int, error) {
	bi := 0
	for i := 0 ; i < len(p) ; i++ {
		if r.seekptr == r.maxsize {
			return i, io.EOF
		} else {
			if bi > 3 {
				bi = 0
			}
			if bi == 0 {
				binary.BigEndian.PutUint32(r.buf,r.prng.Uint32())
			}
			p[i] = r.buf[bi]
			bi++
			r.seekptr++
		}
	}
	return len(p), nil
}

