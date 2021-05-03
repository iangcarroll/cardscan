package main

import (
	"log"

	"github.com/sf1/go-card/smartcard"
)

var (
	cardCtx *smartcard.Context
	reader  *smartcard.Reader
	card    *smartcard.Card
	err     error
)

func connectToCard() error {
	cardCtx, err = smartcard.EstablishContext()
	if err != nil {
		return err
	}

	log.Println("Awaiting a card to be presented.")

	reader, err = cardCtx.WaitForCardPresent()
	if err != nil {
		return err
	}

	card, err = reader.Connect()
	if err != nil {
		return err
	}

	log.Println("Card ATR:", card.ATR())

	// Select the given application.
	selectApplication()

	return nil
}

func selectApplication() error {
	log.Println("Finding application to select...")

	// Get a list of applications.
	getApps, err := sendDesfireCommand(0x6a, []byte{})
	if err != nil {
		return err
	}

	// Select the first application.
	_, err = sendDesfireCommand(0x5a, getApps.Data()[:3])
	return err
}

func safelyCloseCard() {
	if cardCtx != nil {
		cardCtx.Release()
	}

	if card != nil {
		card.Disconnect()
	}
}

func sendDesfireCommand(ins byte, data []byte) (smartcard.ResponseAPDU, error) {
	command := smartcard.CommandAPDU{0x90, ins, 0x00, 0x00} // CLA, INS, P1, P2 bytes
	if len(data) > 0 {
		command = append(command, uint8(len(data)))
		command = append(command, data...) // Data
	}
	command = append(command, 0x00) // Le byte

	response, err := card.TransmitAPDU(command)

	if debugLogging {
		if err != nil {
			log.Printf("[ Tx Error ] %s", err)
			log.Printf("[ ..cont'd ] Sent APDU Command: %x\n", []byte(command))
		} else if response.SW1() != 0x91 || response.SW2() != 0x00 {
			log.Printf("[APDU Error] SW1=%x, SW2=%x\n", response.SW1(), response.SW2())
			log.Printf("[ ..cont'd ] Sent APDU Command: %x\n", []byte(command))
		} else {
			log.Printf("[ APDU Res ] SW1=%x, SW2=%x\n", response.SW1(), response.SW2())
			log.Printf("[ ..cont'd ] Received APDU Res: %x\n", response.Data())
		}
	}

	return response, err
}
