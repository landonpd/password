package userInterface //user ui

import (
	"fmt"
	pswrd "password/passwordFunctions"
	stor "password/passwordStorage"
	"strconv"

	"strings"

	"bytes"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"golang.design/x/clipboard"
)

//	type CharmConfig struct {
//	    DatabaseName string
//	    AccountName string
//		AccountKey string
//	 }

//issues, update adds a new password instead of actually updating,
// deleting all of the passwords crashes the program, slice out of range error, need error check to catch that, it happens on 297,
//it might be from spamming delete key.
//

type DisplayType int

const (
	NewAcctDisplay DisplayType = iota
	GetInputDisplay
	passwordsDisplay
	generatePasswordDisplay
)

// type inputType int

// const (
// 	password DisplayType = iota
// 	length
// 	website

// )

type model struct {
	fileData      string
	key           []byte
	passwords     []pswrd.SavedPassword
	display       DisplayType
	numChoices    int
	cursor        int
	pswrdtoUpdate int
	master        bool
	website       bool
	generated     bool
	reenterPswd   bool
	textInput     textinput.Model
	// db            *kv.KV
	// config        pswrd.CharmConfig
}

// creates model with given data
func InitialModel(fileData string, displayType DisplayType) model { //DB *kv.KV, charmConfig pswrd.CharmConfig
	ti := textinput.New()
	ti.Placeholder = "Password"
	ti.Focus() //what is this
	ti.CharLimit = 156
	ti.Width = 20

	return model{

		fileData:      fileData,
		key:           []byte{},
		passwords:     []pswrd.SavedPassword{},
		display:       displayType,
		numChoices:    0,
		cursor:        0,
		pswrdtoUpdate: 0,
		master:        true,
		website:       false,
		generated:     false,
		reenterPswd:   false,
		textInput:     ti,
		// db:            DB,
		// config:        charmConfig,
	}

}

// what I need to do everytime I quit out
func exitTasks(m model) (tea.Model, tea.Cmd) {
	if len(m.passwords) != 0 { //move this if statement into WritePasswords if I want a new user to be able to put in a master password and no other passwords and still save the master password
		stor.WritePasswords(m.key, "passwordFile.txt", m.passwords)
	}
	return m, tea.Quit
}

// initializes bubbleTea u/i
func (m model) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return textinput.Blink
}

// updates which logic function to be using
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Make sure these keys always quit
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "esc" || k == "ctrl+c" || k == "ctrl+z" {
			// stor.WritePasswordsCharm(m.passwords, len(m.passwords[0].Website), m.db, m.config)
			return exitTasks(m)
		}
	}

	// Hand off the message and model to the appropriate update function for the
	// appropriate view based on the current state.
	switch m.display {
	case NewAcctDisplay:
		return updateNewAcct(msg, m)
	case GetInputDisplay:
		return updateInput(msg, m)
	case passwordsDisplay:
		return updatePasswords(msg, m)

	case generatePasswordDisplay:
		return updateGeneratePassword(msg, m)
	}
	return m, nil
}

// this function determines which screen/view to go to when
func (m model) View() string {
	var s string
	switch m.display {
	case NewAcctDisplay:
		s = newAcctView(m)
	case GetInputDisplay:
		s = inputView(m)
	case passwordsDisplay:
		s = passwordsDisplayView(m)
	case generatePasswordDisplay:
		s = generatePasswordView(m)

	}
	return s
}

// u/i screen to for new accounts, has user create and confirm a master password
func updateNewAcct(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			inputtedPassword := m.textInput.Value()
			//first time they are putting in the password
			if !m.reenterPswd {

				m.key = pswrd.HashKey([]byte(inputtedPassword)) //[]byte(inputtedPassword)
				m.reenterPswd = true
				//fmt.Println("hi")
			} else {
				if bytes.Equal(pswrd.HashKey([]byte(inputtedPassword)), m.key) {
					m.display = passwordsDisplay
					m.master = false
				} else {
					m.reenterPswd = false
				}
			}
			m.textInput.Reset()
		case "ctrl+c", "ctrl+z", "esc":
			// stor.WritePasswordsCharm(m.passwords, len(m.passwords[0].Website), m.db, m.config)
			return exitTasks(m)
		}
	}
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// view for new accounts, ran the first time the program is ran, might be changing still, only ran once then file is updated and it doesn't happen again
func newAcctView(m model) string {
	var s string
	//m.textInput.Placeholder="password"
	if m.reenterPswd {
		s = "reenter your password to confirm\nMaster Password:" //not removing old message before writing new one, don't know how to deal with that yet
	} else {
		s = "Welcome, please enter a master password to create a password manager.\nMaster Password:"
	}
	//s+=m.textInput.View()
	return s
}

// u/i screen where user types something in, could be master password, a website, or a password
func updateInput(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			inputtedPassword := m.textInput.Value()
			if m.master {
				fmt.Println(inputtedPassword)
				//check if the entered password is correct
				if pswrd.CheckMasterPswrd(inputtedPassword, m.fileData) {
					fmt.Println("hi")
					m.key = pswrd.HashKey([]byte(inputtedPassword)) //store the hashed key to use to encrypt and decrypt later on
					//splitting data from file into website and password pairs and storing them in m.paswords
					plainTxt := string(pswrd.DecryptAes(m.key, []byte(m.fileData))) //something going wrong here I think
					// fmt.Println(plainTxt)
					fileDataLst := strings.Split(plainTxt, "\n")
					fileDataLst = fileDataLst[1 : len(fileDataLst)-1] //don't need first line, it is a check, don't need last line, it is blank
					fmt.Println(fileDataLst)
					for _, passwords := range fileDataLst {
						// fmt.Println(passwords)
						savedPassword := strings.Split(passwords, ": ") //consider making this a variable
						tempPsswrd := pswrd.SavedPassword{Website: savedPassword[0], EncryptedPswrd: savedPassword[1]}
						m.passwords = append(m.passwords, tempPsswrd)
					}
					// fmt.Println("passwords")
					// pswrd.DisplayPasswords(m.passwords)

					m.display = passwordsDisplay //got password correct so going to next seciton
					m.cursor = 0
					m.master = false
				} //maybe else statement to change to the reenter password thing, different thing to

				//entering a website
			} else if m.website {
				var newPassword pswrd.SavedPassword

				m.website = false

				newPassword.Website = inputtedPassword
				newPassword.EncryptedPswrd = ""
				m.passwords = append(m.passwords, newPassword)
				//entering a password, it is not generated, going to newest one which has already had a website
			} else if !m.generated {
				m.passwords[len(m.passwords)-1].EncryptedPswrd = m.textInput.Value()
				m.cursor = 0
				m.display = passwordsDisplay
				//we are generating a password for the newest pair
			} else {
				n, _ := strconv.Atoi(m.textInput.Value())
				m.passwords[len(m.passwords)-1].EncryptedPswrd = pswrd.GeneratePassword(n)
				m.cursor = 0
				m.display = passwordsDisplay
			}
			m.textInput.Reset()
		case "ctrl+c", "ctrl+z", "esc":
			// stor.WritePasswordsCharm(m.passwords, len(m.passwords[0].Website), m.db, m.config)
			return exitTasks(m)
		}
	}
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// u/i screen for inputting somethimg, could be password or website
func inputView(m model) string {
	var s string
	if m.master || (!m.website && !m.generated) {
		m.textInput.Placeholder = "password"
		s = "Enter Password: "
	} else if m.website {
		//new password
		m.textInput.Placeholder = "website"
		s = "Enter the website name:"
		//m.website=false
	} else if m.generated {
		m.textInput.Placeholder = "10"
		s = "Enter the length of the passwrod to be generated:"
		//m.generated=false
	}
	if !m.master {
		s += m.textInput.View()
	}
	return s
}

// u/i screen default screen that displays the options
func updatePasswords(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	m.numChoices = len(m.passwords)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			if m.cursor < m.numChoices-1 {
				m.cursor++
			}
		case "k", "up": //maybe add allowing the key presses for numbers to do the same thing
			if m.cursor > 0 {
				m.cursor--
			}
		case "d", "D": //delete
			//here will be delete, get cursor location, use that to delete that location from the list
			if len(m.passwords) != 0 {
				toDelete := m.cursor
				m.passwords = append(m.passwords[:toDelete], m.passwords[toDelete+1:]...)
				m.cursor = 0
			}
			//stor.WritePasswords(m.passwords,len(m.passwords[0].Website),m.db,m.config) //need to write after every delete?
		case "u", "U":
			//update
			m.display = generatePasswordDisplay //go to generate a password page
			m.pswrdtoUpdate = m.cursor
			m.cursor = 0
			// m.numChoices=2
			//stor.WritePasswords(m.passwords,len(m.passwords[0].Website),m.db,m.config)

		case "a", "A": //add new passwords
			m.display = generatePasswordDisplay
			m.pswrdtoUpdate = -1 //does this work, is the password at 0 not one of the options?, -1 is not one of the options, does work
			m.cursor = 0
			// m.website=true
			//m.numChoices=2
			//stor.WritßePasswords(m.passwords,len(m.passwords[0].Website),m.db,m.config)

		case "c", "C", " ", "enter":
			//code to copy here
			strToCopy := m.passwords[m.cursor].EncryptedPswrd
			clipboard.Write(clipboard.FmtText, []byte(strToCopy))
			//m.cursor=0
		case "q", "Q", "ctrl+c", "ctrl+z", "esc":
			// stor.WritePasswordsCharm(m.passwords, len(m.passwords[0].Website), m.db, m.config)
			return exitTasks(m)
		}

	}
	return m, nil
}

// u/i screen to display the passwords, default screen display
func passwordsDisplayView(m model) string {
	var s string
	s += "\nPress c(copy), a(add), d(delete), u(update) or q(quit)\n\n"
	for i, savedPasswords := range m.passwords {
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}
		// if i > 0 {
		s += fmt.Sprintf("%s [ ] %s\n", cursor, savedPasswords.Website)
		// }
	}
	return s
}

// u/i logic for generating a password
func updateGeneratePassword(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	m.numChoices = 2
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			if m.cursor < m.numChoices-1 {
				m.cursor++
			}
		case "k", "up":
			if m.cursor > 0 {
				m.cursor--
			}

		case "enter", " ":
			//var newPassword pswrd.SavedPassword
			if m.pswrdtoUpdate == -1 { //making a new password p$vX6WNMXZ test new firstNew

				m.website = true

				if m.cursor == 0 {
					m.generated = true
				} else {
					m.generated = false
				}
			} else {
				m.pswrdtoUpdate = m.cursor
				if m.cursor == 0 {
					m.generated = true
				} else {
					m.generated = false
				}
			}
			m.cursor = 0
			m.display = GetInputDisplay
			m.textInput.Reset()
		case "q", "Q", "ctrl+c", "ctrl+z", "esc":
			// stor.WritePasswordsCharm(m.passwords, len(m.passwords[0].Website), m.db, m.config)
			return exitTasks(m)
		}
	}
	return m, nil
}

// u/i screen for generating a password
func generatePasswordView(m model) string {
	var s string
	choices := []string{"1. Generate a password", "2. make own password"}
	for i, str := range choices { //need to skip the first one, redo how passwords are displayed
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		s += fmt.Sprintf("%s [ ] %s\n", cursor, str)

	}

	return s
}
