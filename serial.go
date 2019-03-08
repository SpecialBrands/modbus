// Copyright 2014 Quoc-Viet Nguyen. All rights reserved.
// This software may be modified and distributed under the terms
// of the BSD license. See the LICENSE file for details.

package modbus

import (
	"io"
	"log"
	"sync"
	"time"

	"fmt"

	"github.com/goburrow/serial"
)

const (
	// Default timeout
	serialTimeout     = 5 * time.Second
)

// serialPort has configuration and I/O controller.
type serialPort struct {
	// Serial port configuration.
	serial.Config

	Logger      *log.Logger

	mu sync.Mutex
	// port is platform-dependent data structure for serial port.
	port         io.ReadWriteCloser
}

func (mb *serialPort) Connect() (err error) {
		fmt.Println("SERIAL Opening",mb.Config)
	mb.mu.Lock()
	defer mb.mu.Unlock()

	return mb.connect()
}

var (
	portlist     = make(map[string]io.ReadWriteCloser)
	portlistlock sync.Mutex
)

// connect connects to the serial port if it is not connected. Caller must hold the mutex.
func (mb *serialPort) connect() error {
	portlistlock.Lock()
	defer portlistlock.Unlock()
	if mb.port == nil {

		if existing, found := portlist[mb.Config.Address]; found {
			mb.port = existing
		} else {
		fmt.Println("Opening",mb.Config)
			port, err := serial.Open(&mb.Config)
			if err != nil {
			fmt.Println("Got error ",err)
				return err
			}
			portlist[mb.Config.Address] = port
			mb.port = port
		}
	}
	return nil
}

func (mb *serialPort) Close() (err error) {
	mb.mu.Lock()
	defer mb.mu.Unlock()

	return mb.close()
}

// close closes the serial port if it is connected. Caller must hold the mutex.
func (mb *serialPort) close() (err error) {
	if mb.port != nil {
		err = mb.port.Close()
		mb.port = nil
	}
	return
}

func (mb *serialPort) logf(format string, v ...interface{}) {
	if mb.Logger != nil {
		mb.Logger.Printf(format, v...)
	}
}

