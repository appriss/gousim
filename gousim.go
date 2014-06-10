package gousim

import (
  "errors"
  "time"
  "encoding/json"
  "io"
  "fmt"
  "net"
  "regexp"
  "math/rand"
)

type sample struct {
	key string
	value int64
}

type RecordingSession struct {
	logchan chan sample
	samplesize int
	dataMap map[string] []int64
	Recording bool `json:"-"`
	StartTime time.Time
	EndTime time.Time
	signalChan chan bool
	Prefix string `json:"-"`
}

func NewRecordingSession(samplesize int) *RecordingSession {
	s := &RecordingSession{}
	s.samplesize = samplesize
	s.logchan = make(chan sample, s.samplesize/2)
	s.dataMap = make(map[string] []int64)
	s.Recording = false
	return s
} 

func (s *RecordingSession) Stop() error {
	if !s.Recording {
		return errors.New("Recroding already stopped.")
	}
	sm := sample{"STOPRECORDING",0}
	s.logchan <- sm
	<- s.signalChan	
	return nil
}

func (s *RecordingSession) Start() error {
	s.signalChan = make(chan bool)
	if s.Recording {
		return errors.New("Recording already started.")
	}
	fmt.Println("Starting Recording...")
	go s.processSamples()
	<- s.signalChan
	return nil
}

func (s *RecordingSession) LogSample(metric string, value int64) error {
	if !s.Recording {
		return errors.New("Session is not currently recording.")
	}
	if s.Prefix != "" {
		metric = s.Prefix + "_" + metric
	}
	sm := sample{metric, value}
	s.logchan <- sm
	return nil
}

func (s *RecordingSession) processSamples() {
	s.StartTime = time.Now()
	s.Recording = true
	s.signalChan <- true
	for {
		sample := <- s.logchan
		if sample.key == "STOPRECORDING" {
			s.EndTime = time.Now()
			s.Recording = false
			s.signalChan <- true
			return
		}
		if s.dataMap[sample.key] == nil {
			s.dataMap[sample.key] = make([]int64, 0, s.samplesize)
		}
		s.dataMap[sample.key] = append(s.dataMap[sample.key], sample.value)
	}
}

type persistentSession struct {
	StartTime time.Time
	EndTime time.Time
	Samples map[string] []int64
}

func (s *RecordingSession) Export() []byte {
	e := &persistentSession{s.StartTime, s.EndTime, s.dataMap}
	j,_ := json.Marshal(e)
	return j
}

func (s *RecordingSession) ExportStream(w io.Writer) error {
	e := &persistentSession{s.StartTime, s.EndTime, s.dataMap}
	enc := json.NewEncoder(w)
	err := enc.Encode(e)
	return err
}

// Implements an instrumented dial. Use in the place of Dial in many applications
func (s *RecordingSession) Dial(network, addr string) (net.Conn, error) {
	reg, _ := regexp.Compile(`(.*):(\d+)$`)
	res := reg.FindStringSubmatch(addr)
	if res == nil {
		return nil, errors.New("Unable to parse address.")
	}
	host := res[1]
	port := res[2]
	ip := net.ParseIP(host)
	if ip == nil {
		startClock := time.Now()
		ips, err := net.LookupIP(host)
		if err != nil {
			return nil, err
		}
		s.LogSample("DNS", time.Now().Sub(startClock).Nanoseconds() )
		ip = ips[rand.Intn(len(ips))]
	}
	startClock := time.Now()
	conn, err := net.Dial("tcp", ip.String() + ":" + port)
	if err != nil {
		return nil, err
	}
	s.LogSample("TCP", time.Now().Sub(startClock).Nanoseconds() )
	return conn, nil
}

func LoadSession( data []byte, samplesize int ) (*RecordingSession,error) {
	sess := NewRecordingSession(samplesize)
	deserial := &persistentSession{}
	err := json.Unmarshal(data, deserial)
	if err != nil {
		return nil, err
	}
	sess.StartTime = deserial.StartTime
	sess.EndTime = deserial.EndTime
	sess.dataMap = deserial.Samples
	return sess, nil
}

func LoadStream(r io.Reader, samplesize int) (*RecordingSession, error) {
	sess := NewRecordingSession(samplesize)
	deserial := &persistentSession{}
	dec := json.NewDecoder(r)
	err := dec.Decode(deserial)
	if err != nil {
		return nil, err
	}
	sess.StartTime = deserial.StartTime
	sess.EndTime = deserial.EndTime
	sess.dataMap = deserial.Samples
	return sess, nil
}






