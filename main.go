package main

import (
	"fmt"
	"log"
	"time"
)

var (
	debugLogging = true

	clipperApplication    = []byte{0x90, 0x11, 0xf2}
	clipperiOSApplication = []byte{0x91, 0x11, 0xf2}
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type DiscoveredCommand struct {
	Ins byte
	SW1 byte
	SW2 byte
}

func tryAllCommands() (commands []DiscoveredCommand, err error) {
	for i := uint8(0); i < 255; i++ {
		res, err := sendDesfireCommand(i, []byte{})

		if err != nil {
			return commands, err
		}

		if res.SW1() != 0x91 || res.SW2() != 0x1c {
			commands = append(commands, DiscoveredCommand{
				Ins: i,
				SW1: res.SW1(),
				SW2: res.SW2(),
			})
		}

		selectApplication()
	}

	return commands, nil
}

// Returns commands that appear in the `secondCard`, but not the `firstCard`.
func diffCommands(firstCard []DiscoveredCommand, secondCard []DiscoveredCommand) (results []DiscoveredCommand) {
	for _, secondCardCommand := range secondCard {
		existsInFirstCard := false

		// Iterate through all commands discovered in the first card, and look
		// for similar commands.
		for _, firstCardCommand := range firstCard {
			if firstCardCommand.Ins == secondCardCommand.Ins {
				existsInFirstCard = true
			}
		}

		// If we found it in both cards, we don't return the result.
		if existsInFirstCard {
			continue
		}

		results = append(results, secondCardCommand)
	}

	return results
}

// Wrapper around `diffCommands` to identify all commands that differ between two cards.
func compareTwoCards() {
	firstCardCommands, err := tryAllCommands()
	check(err)

	// Disconnect from this card.
	safelyCloseCard()

	log.Println("Please remove first card...")
	time.Sleep(time.Second * 3)

	// Connect to the second card.
	if err := connectToCard(); err != nil {
		panic(err)
	}

	secondCardCommands, err := tryAllCommands()
	check(err)

	log.Println(len(firstCardCommands), len(secondCardCommands))

	onlyInFirstCard := diffCommands(secondCardCommands, firstCardCommands)
	onlyInSecondCard := diffCommands(firstCardCommands, secondCardCommands)

	fmt.Printf("\n")
	for _, command := range onlyInFirstCard {
		fmt.Printf("0x%02X is only in the first card, but not the second card.\n", command.Ins)
	}

	fmt.Printf("\n")
	for _, command := range onlyInSecondCard {
		fmt.Printf("0x%02X is only in the second card, but not the first card.\n", command.Ins)
	}
}

// Parses the DESFire file type byte.
func getFileType(kind byte) string {
	switch kind {
	case 0x00:
		return "Standard"
	case 0x01:
		return "Backup"
	case 0x02:
		return "Value"
	case 0x03:
		return "Linear"
	case 0x04:
		return "Cyclic"
	default:
		return "Unknown"
	}
}

// Lists files from a card.
func listFiles() {
	files, err := sendDesfireCommand(0x6f, []byte{})
	check(err)

	for _, file := range files.Data() {
		fileInfo, err := sendDesfireCommand(0xf5, []byte{file})
		check(err)

		data := fileInfo.Data()
		log.Printf("File 0x%02x, type %s, security level 0x%02x", file, getFileType(data[0]), data[1])
	}
}

func main() {
	// Ensure we properly close out with any issues.
	defer safelyCloseCard()

	// Connect to the card.
	if err := connectToCard(); err != nil {
		panic(err)
	}

	// List files.
	listFiles()
}
