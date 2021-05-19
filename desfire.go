package main

import "log"

func getKeySettings() {
	settings, err := sendDesfireCommand(0x45, []byte{})
	check(err)

	log.Printf("Key settings: %08b\n", settings[0])
	log.Printf("Max # of keys: %02x\n", settings[1])
}

func getKeyVersion(keyNo uint8) {
	version, err := sendDesfireCommand(0x64, []byte{keyNo})
	check(err)

	if len(version.Data()) == 1 {
		log.Printf("Key %02x has version %02x.\n", keyNo, version.Data()[0])
	}
}

func getKeyVersions() {
	for i := 0; i < 0xFF; i++ {
		getKeyVersion(uint8(i))
	}
}
