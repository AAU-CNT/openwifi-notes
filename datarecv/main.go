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

type Measurement struct {
	Timestamp uint64
	Value     uint32
}

func NewMeasurementFromBuffer(buffer []byte) (m Measurement) {
	m.Timestamp = binary.LittleEndian.Uint64(buffer[0:])
	m.Value = binary.LittleEndian.Uint32(buffer[8:])

	return
}

func (m *Measurement) WriteToFile(f io.Writer) error {
	_, err := fmt.Fprintf(f, "%d,%d\n", m.Timestamp, m.Value)
	return err
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
