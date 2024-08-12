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

type DisplayType int

const (
	NewAcctDisplay DisplayType = iota
	GetInputDisplay
	passwordsDisplay
	generatePasswordDisplay
)

type model struct {
	fileName      string
	fileData      string
	key           []byte
	passwords     []pswrd.SavedPassword
	display       DisplayType
	numChoices    int
	cursor        int
	pswrdtoUpdate int
	wrongCount    int
	master        bool
	website       bool
	generated     bool
	reenterPswd   bool
	textInput     textinput.Model
}

// creates model with given data
func InitialModel(fileName, fileData string, displayType DisplayType) model {
	ti := textinput.New()
	ti.Placeholder = "Password"
	ti.Focus() //what is this?
	ti.CharLimit = 156
	ti.Width = 20

	return model{
		fileName:      fileName,
		fileData:      fileData,
		key:           []byte{},
		passwords:     []pswrd.SavedPassword{},
		display:       displayType,
		numChoices:    0,
		cursor:        0,
		pswrdtoUpdate: 0,
		wrongCount:    0,
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
		stor.WritePasswords(m.key, m.fileName, m.passwords)
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
				m.key = pswrd.HashKey([]byte(inputtedPassword))
				m.reenterPswd = true

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
			return exitTasks(m)
		}
	}
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// view for new accounts, ran the first time the program is ran, might be changing still, only ran once then file is updated and it doesn't happen again
func newAcctView(m model) string {
	var s string
	if m.reenterPswd {
		s = "reenter your password to confirm\nMaster Password:" //not removing old message before writing new one, don't know how to deal with that yet
	} else {
		s = "Welcome, please enter a master password to create a password manager.\nMaster Password:"
	}
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
				//check if the entered password is correct
				if plainTxt, correct := pswrd.CheckMasterPswrd(inputtedPassword, m.fileData); correct { //got the right password

					m.key = pswrd.HashKey([]byte(inputtedPassword)) //store the hashed key to use to encrypt and decrypt later on
					//splitting data from file into website and password pairs and storing them in m.paswords
					fileDataLst := strings.Split(plainTxt, "\n")
					fileDataLst = fileDataLst[1 : len(fileDataLst)-1] //don't need first line, it is a check, don't need last line, it is blank

					for _, passwords := range fileDataLst {
						// fmt.Println(passwords)
						savedPassword := strings.Split(passwords, ": ") //consider making deliminater a variable
						tempPsswrd := pswrd.SavedPassword{Website: savedPassword[0], EncryptedPswrd: savedPassword[1]}
						m.passwords = append(m.passwords, tempPsswrd)
					}
					m.display = passwordsDisplay //got password correct so going to next seciton
					m.cursor = 0
					m.master = false
					m.reenterPswd = false
				} else {
					m.wrongCount++
					m.reenterPswd = true
				}

				//entering a website
			} else if m.website {
				var newPassword pswrd.SavedPassword

				m.website = false

				newPassword.Website = inputtedPassword
				newPassword.EncryptedPswrd = ""
				m.passwords = append(m.passwords, newPassword)
				//entering a password, it is not generated, going to newest one which has already had a website or where update points to
			} else if !m.generated {
				if m.pswrdtoUpdate == -1 {
					m.passwords[len(m.passwords)-1].EncryptedPswrd = m.textInput.Value()
				} else {
					//updating a password
					m.passwords[m.pswrdtoUpdate].EncryptedPswrd = m.textInput.Value()
				}
				m.cursor = 0
				m.display = passwordsDisplay
				//we are generating a password for the newest pair
			} else {
				n, _ := strconv.Atoi(m.textInput.Value())
				if m.pswrdtoUpdate == -1 {
					m.passwords[len(m.passwords)-1].EncryptedPswrd = pswrd.GeneratePassword(n)
				} else {
					//updating a password
					m.passwords[m.pswrdtoUpdate].EncryptedPswrd = pswrd.GeneratePassword(n)
				}

				m.cursor = 0
				m.display = passwordsDisplay
			}
			m.textInput.Reset()
		case "ctrl+c", "ctrl+z", "esc":
			return exitTasks(m)
		}
	}
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

// u/i screen for inputting somethimg, could be password or website
func inputView(m model) string {
	var s string

	if (m.master && !m.reenterPswd) || (!m.website && !m.generated && !m.master) {
		m.textInput.Placeholder = "password"
		s = "Enter Password: "

	} else if m.master && m.reenterPswd {
		s = "Incorrect " + strconv.Itoa(m.wrongCount) + " time(s) try again: "
	} else if m.website {
		//new password
		m.textInput.Placeholder = "website"
		s = "Enter the website name:"
	} else if m.generated {
		m.textInput.Placeholder = "10"
		s = "Enter the length of the passwrod to be generated:"
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
		case "u", "U":
			//update
			m.display = generatePasswordDisplay //go to generate a password page
			m.pswrdtoUpdate = m.cursor
			m.cursor = 0
		case "a", "A": //add new passwords
			m.display = generatePasswordDisplay
			m.pswrdtoUpdate = -1

		case "c", "C", " ", "enter":
			//code to copy here
			strToCopy := m.passwords[m.cursor].EncryptedPswrd
			clipboard.Write(clipboard.FmtText, []byte(strToCopy))
		case "m", "M":
			//updating master password, this just assumes that we are a new account and has the user go to that page, this works because
			// it updates what the master password is and everything is already decrypted so can just encrypt later with new password
			m.display = NewAcctDisplay
		case "q", "Q", "ctrl+c", "ctrl+z", "esc":
			return exitTasks(m)
		}

	}
	return m, nil
}

// u/i screen to display the passwords, default screen display
func passwordsDisplayView(m model) string {
	var s string
	s += "\nPress c(copy), a(add), d(delete), u(update), m(master) or q(quit)\n\n"
	for i, savedPasswords := range m.passwords {
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}
		s += fmt.Sprintf("%s [ ] %s\n", cursor, savedPasswords.Website)
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
			if m.pswrdtoUpdate == -1 { //making a new password, need to enter the website as well
				m.website = true
			}
			if m.cursor == 0 {
				m.generated = true
			} else {
				m.generated = false
			}
			m.cursor = 0
			m.display = GetInputDisplay
			m.textInput.Reset()
		case "q", "Q", "ctrl+c", "ctrl+z", "esc":
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
