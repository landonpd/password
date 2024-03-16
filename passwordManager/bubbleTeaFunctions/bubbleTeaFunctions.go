package bubbleTeaFunctions

import(
	"github.com/charmbracelet/bubbles/textinput"
    tea "github.com/charmbracelet/bubbletea"
	pswrd "Users/landondixon/src/goCode/passwordManager/passwordFunctions"
	"github.com/charmbracelet/charm/kv"
	"strconv"
	"golang.design/x/clipboard"
	"fmt"
	charm "Users/landondixon/src/goCode/passwordManager/charmFunctions"
	
)

// type CharmConfig struct {
//     DatabaseName string
//     AccountName string
// 	AccountKey string
//  }
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

type model struct{
	passwords []pswrd.SavedPassword
	display DisplayType
	numChoices int 
	cursor int
	pswrdtoUpdate int
	master bool
	website bool
	generated bool
	reenterPswd bool
	textInput textinput.Model
	db *kv.KV
	config pswrd.CharmConfig

	

}



func InitialModel(psswrds []pswrd.SavedPassword,DB *kv.KV,charmConfig pswrd.CharmConfig,displayType DisplayType) model {
	ti := textinput.New()
	ti.Placeholder = "Password"
	ti.Focus() //what is this
	ti.CharLimit = 156
	ti.Width = 20
	
	return model{
		
		passwords: psswrds,
		display:displayType,
		numChoices: len(psswrds)-1,
		cursor:0,
		pswrdtoUpdate:0,
		master:true,
		website:false,
		generated:false,
		reenterPswd:false,
		textInput:ti,
		db:DB,
		config:charmConfig,

	}
	
}



func (m model) Init() tea.Cmd {
    // Just return `nil`, which means "no I/O right now, please."
    return textinput.Blink
}


func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Make sure these keys always quit
	if msg, ok := msg.(tea.KeyMsg); ok {
		k := msg.String()
		if k == "esc" || k == "ctrl+c" || k=="ctrl+z" {
			charm.WritePasswords(m.passwords,len(m.passwords[0].Website),m.db,m.config)
			return m, tea.Quit
		}
	}

	// Hand off the message and model to the appropriate update function for the
	// appropriate view based on the current state.
	switch m.display{
	case NewAcctDisplay:
		return updateNewAcct(msg,m)
	case GetInputDisplay:
		return updateInput(msg,m)
	case passwordsDisplay:
		return updatePasswords(msg,m)
	
	case generatePasswordDisplay:
		return updateGeneratePassword(msg,m)
	}
	return m,nil
}

func (m model) View() string {
	var s string
	switch m.display{
	case NewAcctDisplay:
		s=newAcctView(m)
	case GetInputDisplay:
		s=inputView(m)
	case passwordsDisplay:
		s=passwordsDisplayView(m)
	case generatePasswordDisplay:
		s=generatePasswordView(m)
	
	}
	return s
}

func updateNewAcct(msg tea.Msg, m model) (tea.Model, tea.Cmd){
	var cmd tea.Cmd
	
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			inputtedPassword:=m.textInput.Value()
			if !m.reenterPswd{
				m.passwords[0].EncryptedPswrd=inputtedPassword
				m.reenterPswd=true
				//fmt.Println("hi")
			}else{
				if inputtedPassword==m.passwords[0].EncryptedPswrd{
					m.display=passwordsDisplay
					m.master=false
				}else{
					m.reenterPswd=false
				}
			}
			m.textInput.Reset()
		case "ctrl+c","ctrl+z","esc":
			charm.WritePasswords(m.passwords,len(m.passwords[0].Website),m.db,m.config)
			return m,tea.Quit
		}
	}
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func  newAcctView(m model) string {
	 var s string
	//m.textInput.Placeholder="password"
	if m.reenterPswd{
		s="reenter your password to confirm\nMaster Password:" //not removing old message before writing new one, don't know how to deal with that yet
	}else{
		s="Welcome, please enter a master password to create a password manager.\nMaster Password:"
	}
	//s+=m.textInput.View()
	return s
} 

func updateInput(msg tea.Msg, m model) (tea.Model, tea.Cmd){
	var cmd tea.Cmd
	
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			inputtedPassword:=m.textInput.Value()
			if m.master{
				//fmt.Println(inputtedPassword)
				if inputtedPassword==m.passwords[0].EncryptedPswrd {
					m.display=passwordsDisplay //got password correct so going to next seciton
					m.cursor=0
					m.master=false
				}
			}else if m.website{
				var newPassword pswrd.SavedPassword
				
				m.website=false
				
				
				newPassword.Website=inputtedPassword
				newPassword.EncryptedPswrd=""
				m.passwords=append(m.passwords,newPassword)
				
			}else if !m.generated{
				m.passwords[len(m.passwords)-1].EncryptedPswrd=m.textInput.Value()
				m.cursor=0
				m.display=passwordsDisplay
			}else{
				n,_:=strconv.Atoi(m.textInput.Value())
				m.passwords[len(m.passwords)-1].EncryptedPswrd=pswrd.GeneratePassword(n)
				m.cursor=0
				m.display=passwordsDisplay
			}
			m.textInput.Reset()
		case "ctrl+c","ctrl+z","esc":
			charm.WritePasswords(m.passwords,len(m.passwords[0].Website),m.db,m.config)
			return m,tea.Quit
		}
	}
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func inputView(m model) string{
	var s string
	if m.master || (!m.website&&!m.generated){
		m.textInput.Placeholder="password"
		s="Enter Password: "
	}else if m.website{
		//new password
		m.textInput.Placeholder="website"
		s="Enter the website name:"
		//m.website=false
	}else if m.generated{ 
		m.textInput.Placeholder="10"
		s="Enter the length of the passwrod to be generated:"
		//m.generated=false
	}
	if !m.master {
		s+=m.textInput.View()
	}
	return s
}

func updatePasswords(msg tea.Msg, m model) (tea.Model, tea.Cmd) {
	m.numChoices=len(m.passwords)-1
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
		case "d","D": //delete
			//here will be delete, get cursor location, use that to delete that location from the list
			toDelete:=m.cursor+1
			m.passwords = append(m.passwords[:toDelete], m.passwords[toDelete+1:]...)
			m.cursor=0
			//charm.WritePasswords(m.passwords,len(m.passwords[0].Website),m.db,m.config)
		case "u","U":
			//update
			m.display=generatePasswordDisplay //go to generate a password page
			m.pswrdtoUpdate=m.cursor
			m.cursor=0
			// m.numChoices=2
			//charm.WritePasswords(m.passwords,len(m.passwords[0].Website),m.db,m.config)

		case "a","A": //add new passwords
			m.display=generatePasswordDisplay
			m.pswrdtoUpdate=0
			m.cursor=0
			// m.website=true
			//m.numChoices=2
			//charm.WritßePasswords(m.passwords,len(m.passwords[0].Website),m.db,m.config)

		case "c","C"," ","enter":
			//code to copy here
			strToCopy:=m.passwords[m.cursor+1].EncryptedPswrd
			clipboard.Write(clipboard.FmtText, []byte(strToCopy))
			//m.cursor=0
		case "q","Q","ctrl+c","ctrl+z","esc":
			charm.WritePasswords(m.passwords,len(m.passwords[0].Website),m.db,m.config)
			return m,tea.Quit
		}
	
	}
	return m,nil
}

func passwordsDisplayView(m model ) string{
	var s string
	s+="\nPress c(copy), a(add), d(delete) or q(quit)\n\n"
	for i,savedPasswords:=range m.passwords{
		cursor := " " // no cursor
        if m.cursor == i-1 {
            cursor = ">" // cursor!
        }
		if i>0 {
		s += fmt.Sprintf("%s [ ] %s\n", cursor, savedPasswords.Website)
		}
	}
	return s
}

func updateGeneratePassword(msg tea.Msg, m model) (tea.Model, tea.Cmd) { 
	m.numChoices=2
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

		case "enter"," ":
			//var newPassword pswrd.SavedPassword
			if m.pswrdtoUpdate==0{ //making a new password
				
				m.website=true
				
			//tempPassword:=""
				if m.cursor==0{
					m.generated=true
				}else{
					m.generated=false
				}
			}else {
				m.pswrdtoUpdate=m.cursor+1
				if m.cursor==0{
					m.generated=true
				}else{
					m.generated=false
				}
			}
			m.cursor=0
			m.display=GetInputDisplay
			m.textInput.Reset()
		case "q","Q","ctrl+c","ctrl+z","esc":
			charm.WritePasswords(m.passwords,len(m.passwords[0].Website),m.db,m.config)
			return m,tea.Quit
		}
	}
	return m,nil
}



func  generatePasswordView(m model) string {
	var s string
	choices:=[]string{"1. Generate a password","2. make own password"}
	for i,str:=range choices{ //need to skip the first one, redo how passwords are displayed
		cursor := " " // no cursor
        if m.cursor == i {
            cursor = ">" // cursor!
        }

		
		s += fmt.Sprintf("%s [ ] %s\n", cursor, str)
		
	}

	return s 
}

