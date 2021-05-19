package main

import (
	"encoding/binary"
	"log"
)

func getClipperSerial() uint32 {
	// Read out five bytes from file 0x08 to get the UID.
	file, err := sendDesfireCommand(0xBD, []byte{0x08, 0x01, 0x00, 0x00, 0x05, 0x00, 0x00})
	check(err)

	uid := binary.BigEndian.Uint32(file)
	log.Println("Clipper Serial:", uid)
	return uid
}
