// Copyright 2018, Special Brands Holding BV
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/SpecialBrands/modbus"
)

var (
	device   string
	baudrate int
	databits int
	stopbits int
	parity   string

	enabled            bool
	rtsHighDuringSend  bool
	rtsHighAfterSend   bool
	delayRtsAfterSend  int
	delayRtsBeforeSend int

	clientid        int
	function        string
	registeraddress int
	registerbits    int
	sourceid        int
	destinationid   int

	responsewait int
	verbose      bool
	hideunknown  bool

	readerclose chan bool
)

func main() {
	flag.BoolVar(&verbose, "V", false, "")

	flag.BoolVar(&enabled, "enable_rs485", false, "")
	flag.BoolVar(&rtsHighDuringSend, "rtshighduring", false, "")
	flag.BoolVar(&rtsHighAfterSend, "rtshighafter", false, "")
	flag.IntVar(&delayRtsAfterSend, "delayaftersend", 0, "")
	flag.IntVar(&delayRtsBeforeSend, "delayduringsend", 0, "")

	flag.StringVar(&device, "device", "/dev/ttyUSB0", "device")
	flag.IntVar(&baudrate, "baud", 9600, "baud rate")
	flag.IntVar(&databits, "data", 8, "data bits")
	flag.IntVar(&stopbits, "stop", 1, "stop bits")
	flag.StringVar(&parity, "parity", "N", "parity (N/E/O)")

	flag.IntVar(&clientid, "c", 1, "Client address")
	flag.StringVar(&function, "f", "ri", "function: rh - Read Holding Register (3), ri - Read Input Register (4), w - Write Single Register (6)")
	flag.IntVar(&registeraddress, "a", 7, "Integer register address")
	flag.IntVar(&registerbits, "v", 0, "Value of the register to be written")

	flag.Parse()

	connectionhandler := modbus.NewRTUClientHandler(device)
	connectionhandler.BaudRate = baudrate
	connectionhandler.DataBits = databits
	connectionhandler.StopBits = stopbits
	connectionhandler.Parity = parity
	//connectionhandler.IdleTimeout = 2 * time.Second //For now hardcoded
	connectionhandler.RS485.Enabled = enabled
	connectionhandler.RS485.RtsHighDuringSend = rtsHighDuringSend
	connectionhandler.RS485.RtsHighAfterSend = rtsHighAfterSend
	connectionhandler.RS485.DelayRtsAfterSend = time.Duration(delayRtsAfterSend * int(time.Millisecond))
	connectionhandler.RS485.DelayRtsBeforeSend = time.Duration(delayRtsBeforeSend * int(time.Millisecond))

	connectionhandler.ClientId = byte(clientid) //Apparently this library handles only one slave per clienthandler.
	defer connectionhandler.Close()
	err := connectionhandler.Connect()
	if err != nil {
		fmt.Println("Connection error: ", err)
		os.Exit(1)
	}
	connectionclient := modbus.NewClient(connectionhandler)

	switch function {
	case "w":
		if verbose {
			fmt.Printf("Executing writing function %d. address %d value %d\n", 6, uint16(registeraddress), uint16(registerbits))
		}

		results, err := connectionclient.WriteSingleRegister(uint16(registeraddress), uint16(registerbits)) //6

		if verbose {
			fmt.Printf("Result: %v (error: %v)\n", results, err)
		}
		if err != nil {
			os.Exit(6)
		}
	case "ri":
		if verbose {
			fmt.Printf("Executing reading function %d. address %d\n", 4, uint16(registeraddress))
		}

		results, err := connectionclient.ReadInputRegisters(uint16(registeraddress), 1) //4

		if verbose {
			fmt.Printf("Result: %v (error: %v)\n", results, err)
		}
		if err != nil {
			os.Exit(4)
		}
	case "rh":
		if verbose {
			fmt.Printf("Executing reading function %d. address %d\n", 3, uint16(registeraddress))
		}

		results, err := connectionclient.ReadHoldingRegisters(uint16(registeraddress), 1) //3

		if verbose {
			fmt.Printf("Result: %v (error: %v)\n", results, err)
		}

		if err != nil {
			os.Exit(3)
		}
	}
}
