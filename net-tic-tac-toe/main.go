package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Piece uint8
type GameMsg int

const (
	GameTurnComplete = iota
)

type CustomMsg int

func (cm CustomMsg) Msg() {}

const (
	MessageReceived = iota
)

type GameMessage struct {
	Name  string `json:"name"`
	MoveX int    `json:"moveX"`
	MoveY int    `json:"moveY"`
}

func (gm GameMessage) Msg() {}

func (m *Model) sendMessage(x int, y int) tea.Msg {
	var conn net.Conn
	var err error
	for {
		conn, err = net.Dial("tcp", otherAddress(m.name))
		if err == nil {
			break
		}
	}

	gm := GameMessage{Name: m.name, MoveX: x, MoveY: y}
	data, err := json.Marshal(gm)
	if err != nil {
		panic(err)
	}
	conn.Write(data)

	m.board[x][y] = m.piece
	m.turn = otherPiece(m.piece)
	return CustomMsg(MessageReceived)
}

type Model struct {
	address string
	name    string
	piece   Piece
	curX    int
	curY    int
	turn    Piece
	board   [3][3]Piece
}

const (
	NoPiece = iota
	XPiece
	OPiece
)

const (
	NoPieces = iota
	Xs
	Os
)

func (gm GameMsg) Msg() {}

func New(address string, name string, pi Piece) *Model {
	var tr Piece
	if name == "alice" {
		tr = pi
	} else {
		tr = otherPiece(pi)
	}

	return &Model{
		address: address,
		name:    name,
		piece:   pi,
		turn:    tr,
		curX:    0,
		curY:    0,
		board: [3][3]Piece{
			{NoPieces, NoPieces, NoPieces},
			{NoPieces, NoPieces, NoPieces},
			{NoPieces, NoPieces, NoPieces},
		},
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func otherPiece(p Piece) Piece {
	switch p {
	case XPiece:
		return OPiece
	case OPiece:
		return XPiece
	default:
		return NoPiece
	}
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyUp:
			if m.curX <= 0 {
				m.curX = 2
			} else {
				m.curX -= 1
			}
		case tea.KeyDown:
			if m.curX >= 2 {
				m.curX = 0
			} else {
				m.curX += 1
			}
		case tea.KeyRight:
			if m.curY >= 2 {
				m.curY = 0
			} else {
				m.curY += 1
			}
		case tea.KeyLeft:
			if m.curY <= 0 {
				m.curY = 2
			} else {
				m.curY -= 1
			}
		case tea.KeyEnter, tea.KeySpace:
			if m.turn != m.piece {
				return m, nil
			}

			if m.board[m.curX][m.curY] == NoPiece {
				// m.board[m.curX][m.curY] = m.piece
				m.turn = NoPiece
				return m, func() tea.Msg { return m.sendMessage(m.curX, m.curY) }
			}
		}
		return m, nil
	case GameMessage:
		m.board[msg.MoveX][msg.MoveY] = otherPiece(m.piece)
		m.turn = m.piece
	}

	return m, nil
}

// func (m *Model) gameTurnComplete() tea.Msg {
// 	return GameMsg(GameTurnComplete)
// }

func (p Piece) String() string {
	switch p {
	case NoPiece:
		return " "
	case XPiece:
		return "x"
	case OPiece:
		return "o"
	}
	return " "
}

var (
	// white       = lipgloss.CompleteColor{TrueColor: "#FFFFFF", ANSI256: "15", ANSI: "15"}
	black = lipgloss.CompleteColor{TrueColor: "#000000", ANSI256: "0", ANSI: "0"}
	// magenta     = lipgloss.CompleteColor{TrueColor: "#AF48B6", ANSI256: "13", ANSI: "5"}
	// cyan        = lipgloss.CompleteColor{TrueColor: "#4DA5C9", ANSI256: "14", ANSI: "6"}
	// green       = lipgloss.CompleteColor{TrueColor: "#0dbc79", ANSI256: "2", ANSI: "2"}
	brightgreen = lipgloss.CompleteColor{TrueColor: "#23d18b", ANSI256: "10", ANSI: "10"}
)

var highlightStyle = lipgloss.NewStyle().
	Background(brightgreen).
	Foreground(black)

func (m *Model) renderBoard() string {
	var str string
	for i, row := range m.board {
		for j, box := range row {

			if m.curX == i && m.curY == j {
				str += highlightStyle.Render(box.String())
			} else {
				str += box.String()
			}

			if j < len(row)-1 {
				str += "|"
			}
		}
		if i < len(m.board)-1 {
			str += "\n-----\n"
		} else {
			str += "\n"
		}
	}
	return str
}

func (m *Model) View() string {
	return m.renderBoard()
}

type Move struct {
	X     int
	Y     int
	Piece string
}

func receiveMessage(program *tea.Program, conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 256)
	bytesRead, _ := conn.Read(buf)

	var gm GameMessage
	json.Unmarshal(buf[0:bytesRead], &gm)

	program.Send(gm)
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

type Game struct {
	ID             string
	Player1Name    string
	Player1Address string
	Player2Name    string
	Player2Address string
}

type Server struct {
	GameList []Game
}

type Request struct {
	GameID        string `json:"gameID"`
	RecipientName string `json:"recipientName"`
	RequestType   string `json:"requestType"`
	RequestData   []byte `json:"requestData"`
}

func (s *Server) startRelay() {
	lnr, err := net.Listen("tcp", ":3030")
	if err != nil {
		panic(err)
	}
	defer lnr.Close()

	for {
		conn, err := lnr.Accept()
		if err != nil {
			panic(err)
		}
		go s.routeIncomingRequest(conn)
	}
}

func (s *Server) routeIncomingRequest(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 256)
	bytesRead, _ := conn.Read(buf)

	var req Request
	err := json.Unmarshal(buf[0:bytesRead], &req)
	if err != nil {
		// should tell sender request failed (new connection? just a response?)
		panic(err)
	}

	var game Game
	game, err = s.findGameById(req.GameID)
	if err != nil {
		// tell sender could not find game
		panic(err)
	}

	switch req.RequestType {
	case "Move":
		var mov Move
		err = json.Unmarshal(req.RequestData, &mov)
		if err != nil {
			// should tell sender request failed (new connection? just a response?)
			panic(err)
		}

		s.sendMove(game, req.RecipientName, mov)
	}

	// program.Send(gm)
}

func (s *Server) sendMove(game Game, recipientName string, move Move) {
	var address string
	switch recipientName {
	case game.Player1Address:
		address = game.Player1Address
	case game.Player2Name:
		address = game.Player2Address
	default:
		// tell sender player is not part of this game
		panic(err)
	}

	var conn net.Conn
	var err error
	for {
		conn, err = net.Dial("tcp", address)
		if err == nil {
			break
		}
	}

	data, err := json.Marshal(move)
	if err != nil {
		panic(err)
	}
	conn.Write(data)
}

func (s *Server) findGameById(id string) (Game, error) {
	for _, gm := range s.GameList {
		if gm.ID == id {
			return gm, nil
		}
	}
	return Game{}, CustomError{desc: "Not found"}
}

type CustomError struct {
	desc string
}

func (ce CustomError) Error() string {
	return ce.desc
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

	var piece Piece
	switch args[1] {
	case "x", "X":
		piece = XPiece
	case "o", "O":
		piece = OPiece
	default:
		panic("unknown piece")
	}

	p := tea.NewProgram(New(address, name, piece))

	go startServer(p, address)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
