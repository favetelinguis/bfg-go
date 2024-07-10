package main

// import (
// 	"fmt"
// 	"os"

// 	// "github.com/charmbracelet/bubbles/key"
// 	tea "github.com/charmbracelet/bubbletea"
// 	// "github.com/charmbracelet/lipgloss"
// 	"github.com/favetelinguis/bfg-go/betfair"
// )

// // TODO
// // Tabs and composable views split editors, how to compose models?
// // bubble tea components, multiple component in a view?
// // widget is the name of a component and it has its own model internally?

// func quitSession(session *betfair.Session) tea.Cmd {
// 	return func() tea.Msg {
// 		val, err := session.Logout()
// 		if err != nil {
// 			return errMsg{err}
// 		}
// 		return statusMsg(val.Status)
// 	}
// }

// type statusMsg string
// type errMsg struct{ err error }

// func (e errMsg) Error() string { return e.err.Error() }

// type model struct {
// 	bfSession *betfair.Session
// 	status    string
// 	err       error
// }

// func initalModel() model {
// 	return model{
// 		bfSession: betfair.NewSession(),
// 	}
// }

// func (m model) Init() tea.Cmd {
// 	// TODO here i should do call to BF to setup all base UI elementes.
// 	return quitSession(m.bfSession)
// }

// func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	switch msg := msg.(type) {
// 	case statusMsg:
// 		m.status = string(msg)
// 		return m, nil
// 	case errMsg:
// 		m.err = msg
// 		return m, nil

// 	case tea.KeyMsg:
// 		switch msg.String() {
// 		case "ctrl+c", "q":
// 			return m, tea.Quit
// 		}
// 	}
// 	return m, nil
// }

// func (m model) View() string {
// 	if m.err != nil {
// 		return fmt.Sprintf("\nWe had some trouble: %v\n\n", m.err)
// 	}

// 	s := fmt.Sprintf("Checking %s ...", "logout")

// 	if m.status != "" {
// 		s += fmt.Sprintf("%s!", m.status)
// 	}

// 	return "\n" + s + "\n\n"
// }

// func main() {
// 	// Login and get a token
// 	session := betfair.NewSession()
// 	defer session.Logout()
// 	// Setup stream and auth to stream
// 	stream := betfair.NewStream(session)
// 	defer stream.StopListen()

// 	// Start tui
// 	p := tea.NewProgram(initalModel(), tea.WithAltScreen())

// 	// Connect the tui with the betfair stream so we can send messages to the update fn
// 	stream.StartListen(p.Send)

// 	if _, err := p.Run(); err != nil {
// 		fmt.Printf("Error starting %v", err)
// 		os.Exit(1)
// 	}
// }
