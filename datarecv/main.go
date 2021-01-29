package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync/atomic"
	"time"
)

type FPGAValue uint32

type Measurement struct {
	Timestamp uint64
	Value     FPGAValue
}

func (v FPGAValue) split() (uint16, uint16) {
    return uint16((v >> 16) & 0xFFFF), uint16(v & 0xFFFF)
}

func (v FPGAValue) FormatCsv() string {
    random, backoff := v.split()
    return fmt.Sprintf("%d,%d", random, backoff)
}

func (v FPGAValue) String() string {
    random, backoff := v.split()
    return fmt.Sprintf("{random: %d, backoff: %d}", random, backoff)
}

func NewMeasurementFromBuffer(buffer []byte) (m Measurement) {
	m.Timestamp = binary.LittleEndian.Uint64(buffer[0:])
	m.Value = FPGAValue(binary.LittleEndian.Uint32(buffer[8:]))

	return
}

func WriteHeader(f io.Writer) (err error) {
    _, err = fmt.Fprintln(f, "time, random, backoff")
    return
}

func (m *Measurement) WriteToFile(f io.Writer) (err error) {
	_, err = fmt.Fprintf(f, "%d,%s\n", m.Timestamp, m.Value.FormatCsv())
	return
}

func ParseBuffer(buffer []byte) []Measurement {
	log.Printf("Parsing %v", buffer)
	// First read out number of measurements
	num := buffer[0]
	log.Printf("%d measurements", num)
	buffer = buffer[1:]

	// Allocate array
	res := make([]Measurement, num)

	// Populate array
	for i := 0; i < int(num); i++ {
		res[i] = NewMeasurementFromBuffer(buffer[i*12:])
	}

	return res
}

const bufflen = 128

var count uint32 = 0

func main() {
    outname := flag.String("out", "data.csv", "File to use for outputting data")
    flag.Parse()

	f, err := os.OpenFile(*outname, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}

    err = WriteHeader(f)
    if err != nil {
        log.Fatal(err)
    }

	pc, err := net.ListenPacket("udp", ":8000")
	fmt.Printf("Starting\n")
	if err != nil {
		log.Fatal(err)
	}

	go status()

	buff := make([]byte, bufflen*12)

	for {
		n, _, err := pc.ReadFrom(buff)
		if err != nil {
			log.Fatal(err)
		}

		// Parse buffer
		results := ParseBuffer(buff[0:n])

		for i, res := range results {
			log.Printf("%d, Got %v", i, res)
			res.WriteToFile(f)
		}

		atomic.AddUint32(&count, uint32(len(results)))
	}
}

func status() {
	var last uint32 = 0
	for {
		now := atomic.LoadUint32(&count)
		fmt.Printf("Received %d, since last %d\n", now, now-last)
		last = now
		time.Sleep(1 * time.Second)
	}
}
