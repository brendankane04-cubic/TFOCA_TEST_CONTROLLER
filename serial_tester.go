package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"sort"
	"strconv"
	"strings"
)

//This is a dictionary of every test we currently run organized by name-hex payload
//Each hex payload is a string of hex values converted to bytes later on
var cmds = map[string]string{
	"GPIO_SET_PUD":         "00080006000400120002",
	"GPIO_SET_DIR":         "00080006000400120102",
	"GPIO_SET_LVL":         "00080006000400120200",
	"GPIO_READ_PUD":        "000700050103001200",
	"GPIO_READ_DIR":        "000700050103001201",
	"GPIO_READ_LVL":        "000700050103001202",
	"MDIO_READ_8221_REG":   "00090107000500640A0000",
	"MDIO_READ_8489_REG_1": "000901070105000F7A0000",
	"MDIO_READ_8489_REG_2": "000901070105003F7A0000",
	"I2C_WRITE":            "00080206010400A26EFF",
	"I2C_READ":             "00070205000300A000",
}

func main() {
	//The IP address, the port, and the serial port of the test server
	//The Serial port which the physical connector is plugged into can be easily changed
	var HOST = "172.16.252.9"
	var PORT_SERIAL = 12                                                      //The Serial port the digikey connector is plugged into
	var PORT_BASE = "101XX"                                                   //The base port for the machine
	var PORT = strings.Replace(PORT_BASE, "XX", strconv.Itoa(PORT_SERIAL), 1) //The actual port passed into the socket init

	fmt.Println("RUN")

	//Set up the socket
	conn, err := net.Dial("tcp", HOST+":"+PORT)
	checkErr(err)
	defer conn.Close()

	//Get the names in the command map & sort them so that they run in alphabetical order
	//This way, the "collections" of tests run together & numbered ones run in order
	//If this is skipped, the commands will all be run in a random order which is worse
	names := make([]string, 0, len(cmds))
	for key, _ := range cmds {
		names = append(names, key)
	}
	sort.Strings(names)

	//Send every command in the list of commands we have & output their label
	for index, name := range names {
		command := cmds[name]
		fmt.Printf("Test %d: %s\n", index+1, name)
		fmt.Printf("Command: %s\n", command)
		reply, err := send_test_cmd(conn, command)
		checkErr(err)
		fmt.Printf("Response: %s \n", reply)
	}
}

//Send a command to the user & get the response
func send_test_cmd(conn net.Conn, command string) (string, error) {
	//Convert the string into an array of hex bytes
	data, err := hex.DecodeString(command)
	if err != nil {
		return "", err
	}

	//Send it to the device
	_, err = conn.Write(data)
	if err != nil {
		return "", err
	}

	//Receive a response from the machine
	reply := make([]byte, 256)
	_, err = conn.Read(reply)
	if err != nil {
		return "", err
	}

	//Remove any extraneous "0x00"s
	reply = bytes.Trim(reply, "\x00")

	//Return the reply to the user after converting it to a string
	return hex.EncodeToString(reply), nil
}

//Check an error status
func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
