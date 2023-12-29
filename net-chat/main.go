package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	address     string
	name        string
	chatHistory []string
	chatView    viewport.Model
	chatInput   textinput.Model
}

func New(address string, name string) *Model {
	view := viewport.New(20, 10)
	input := textinput.New()
	input.Focus()
	input.CharLimit = 20 - len(input.Prompt)
	input.Width = 20

	return &Model{
		chatView:  view,
		chatInput: input,
		address:   address,
		name:      name,
	}
}

func (m *Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var viewCmd tea.Cmd
	var inputCmd tea.Cmd

	m.chatView, viewCmd = m.chatView.Update(msg)
	m.chatInput, inputCmd = m.chatInput.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			message := m.chatInput.Value()

			return m, func() tea.Msg { return m.sendMessage(message) }
		}
		return m, nil
	case CustomMsg:
		switch msg {
		case MessageReceived:
			m.chatInput.Reset()
		}
	case ChatMessage:
		m.appendChatHistory(msg)
		return m, tea.Batch(viewCmd, inputCmd)
	}

	return m, tea.Batch(viewCmd, inputCmd)
}

func (m *Model) appendChatHistory(cm ChatMessage) {
	chatLine := fmt.Sprintf("%s: %s", cm.Name, cm.Text)
	m.chatHistory = append(m.chatHistory, chatLine)
	m.chatView.SetContent(m.renderChatHistory())
}

type CustomMsg int

func (cm CustomMsg) Msg() {}

const (
	MessageReceived = iota
)

type ChatMessage struct {
	Name string `json:"name"`
	Text string `json:"text"`
}

func (cm ChatMessage) Msg() {}

func (m *Model) sendMessage(msg string) tea.Msg {
	var conn net.Conn
	var err error
	for {
		conn, err = net.Dial("tcp", otherAddress(m.name))
		if err == nil {
			break
		}
	}

	cm := ChatMessage{Name: m.name, Text: msg}
	data, err := json.Marshal(cm)
	if err != nil {
		panic(err)
	}
	conn.Write(data)

	m.appendChatHistory(cm)
	return CustomMsg(MessageReceived)
}

func (m *Model) renderChatHistory() string {
	var str string
	for idx, msg := range m.chatHistory {
		str += msg
		if idx < len(m.chatHistory) {
			str += "\n"
		}
	}
	return str
}

func receiveMessage(program *tea.Program, conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 256)
	bytesRead, _ := conn.Read(buf)
	// text := string(buf[0:bytesRead])

	var cm ChatMessage
	json.Unmarshal(buf[0:bytesRead], &cm)

	program.Send(cm)
}

func startServer(program *tea.Program, address string) {
	lnr, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}
	defer lnr.Close()

	for {
		conn, err := lnr.Accept()
		if err != nil {
			panic(err)
		}
		go receiveMessage(program, conn)
	}
}

func (m *Model) View() string {
	return lipgloss.JoinVertical(
		lipgloss.Left,
		m.chatView.View(),
		m.chatInput.View(),
	)
}

func otherAddress(name string) string {
	var other string
	switch name {
	case "alice":
		other = ":4040"
	case "bob":
		other = ":8080"
	}
	return other
}

func main() {
	args := os.Args[1:]

	var name string
	var address string
	switch args[0] {
	case "a", "A", "alice":
		name = "alice"
		address = ":8080"
	case "b", "B", "bob":
		name = "bob"
		address = ":4040"
	}

	p := tea.NewProgram(New(address, name))

	go startServer(p, address)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
