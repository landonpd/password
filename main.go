package main

import (
	"fmt"
	"os"
	stor "password/passwordStorage"
	uI "password/userInterface"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	//variables and constants
	const FILENAME = "passwordFile.bin"
	display := uI.GetInputDisplay

	//reading file to get all the passwords, everything is still encrypted at this stage
	fileData, err := stor.ReadData(FILENAME)
	if err != nil {
		fmt.Println("file doesn't exist")
		//maybe make the file, or have it ship with a blank file, not sure wha thte process here should be
		return
	} else if len(fileData) == 0 {
		//file is empty indicating that this is a new account
		display = uI.NewAcctDisplay
	}

	//---------------------------------------------------------------BubbleTea code below----------------------------------------------------------------------------------------\\

	//running bubble tea wich gives really nice terminal interface, see extra file for details

	p := tea.NewProgram(uI.InitialModel(FILENAME, fileData, display)) //db, pswrds.Config
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
