package main

import (
    "fmt"
    "time"
    "encoding/binary"
    "log"
    "sync/atomic"
    "os"
    "net"
)

const bufflen = 16

var count uint32 = 0

func main() {
    f, err := os.OpenFile("data", os.O_RDWR|os.O_CREATE, 0755)
    if err != nil {
        log.Fatal(err)
    }

    pc, err := net.ListenPacket("udp", ":8000")
    fmt.Printf("Starting\n")
    if err != nil {
        log.Fatal(err)
    }

    go status()

    buff := make([]byte, bufflen*4)

    for {
        n, _, err := pc.ReadFrom(buff)
        if err != nil {
            log.Fatal(err)
        }

        // Convert to number of uint32
        n = n / 4

        for i := 0; i < n; i++ {
            value := binary.LittleEndian.Uint32(buff[i * 4:(i+1)*4])
            fmt.Fprintf(f, "%d\n", value)
        }

        atomic.AddUint32(&count, uint32(n))
    }
}

func status() {
    var last uint32 = 0
    for {
        now := atomic.LoadUint32(&count)
        fmt.Printf("Received %d, since last %d\n", now, now - last)
        last = now
        time.Sleep(1 * time.Second)
    }
}
