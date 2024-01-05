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

type StatusMessage string

func (sm StatusMessage) Msg() {}

func (sm StatusMessage) String() string {
	return string(sm)
}

func (m *Model) sendMessage(x int, y int) tea.Msg {
	var conn net.Conn
	var err error
	for {
		conn, err = net.Dial("tcp", ":3030")
		if err == nil {
			break
		}
	}

	requestData, err := json.Marshal(GameMessage{Name: m.name, MoveX: x, MoveY: y})
	if err != nil {
		panic(err)
	}

	data, err := json.Marshal(Request{
		GameID:        m.gameID,
		RecipientName: otherName(m.name),
		RequestType:   "GameMessage",
		RequestData:   requestData,
	})
	if err != nil {
		panic(err)
	}

	conn.Write(data)

	m.board[x][y] = m.piece
	m.turn = otherPiece(m.piece)
	return CustomMsg(MessageReceived)
}

type Model struct {
	gameID  string
	address string
	name    string
	piece   Piece
	curX    int
	curY    int
	turn    Piece
	board   [3][3]Piece
	status  string
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
		gameID:  "0",
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
		status: "starting...",
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
	case StatusMessage:
		fmt.Print(msg.String())
		m.status = msg.String()
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
	return fmt.Sprintf(
		"%s\n%s\n\n%s", m.name, m.renderBoard(), m.status,
	)
}

type Move struct {
	X     int
	Y     int
	Piece string
}

func receiveMessage(program *tea.Program, conn net.Conn) {
	defer conn.Close()

	// fmt.Print("receiving message")
	program.Send(StatusMessage("receiving message"))

	buf := make([]byte, 256)
	bytesRead, _ := conn.Read(buf)

	var gm GameMessage
	err := json.Unmarshal(buf[0:bytesRead], &gm)
	if err != nil {
		panic(
			fmt.Sprintf("could not parse received move\n%v\n%e", gm, err),
		)
	}

	program.Send(gm)
}

func startServer(program *tea.Program, address string) {
	lnr, err := net.Listen("tcp", address)
	if err != nil {
		panic(err)
	}
	defer lnr.Close()

	program.Send(StatusMessage("sanity check..."))

	for {
		conn, err := lnr.Accept()
		if err != nil {
			panic(err)
		}
		program.Send(StatusMessage("got a connection..."))
		go receiveMessage(program, conn)
	}
}

type Game struct {
	ID             string
	Player1Name    string
	Player1Piece   Piece
	Player1Address string
	Player2Name    string
	Player2Piece   Piece
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

func NewRelay() *Server {
	return &Server{
		GameList: []Game{
			{
				ID:             "0",
				Player1Name:    "alice",
				Player1Piece:   XPiece,
				Player1Address: ":8080",
				Player2Name:    "bob",
				Player2Piece:   OPiece,
				Player2Address: ":4040",
			},
		},
	}
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

	buf := make([]byte, 512)
	bytesRead, _ := conn.Read(buf)

	// fmt.Printf("bytes: %d\n%s\n\n", bytesRead, string(buf))

	var req Request
	err := json.Unmarshal(buf[0:bytesRead], &req)
	if err != nil {
		// should tell sender request failed (new connection? just a response?)
		panic(
			fmt.Sprintf("send request failed, bad json?\n%s\n%e", string(buf), err),
		)
	}

	var game Game
	game, err = s.findGameById(req.GameID)
	if err != nil {
		// tell sender could not find game
		panic(
			fmt.Sprintf("sender could not find game\n%v\n%v\n%e", req, req.GameID, err),
		)
	}

	switch req.RequestType {
	case "GameMessage":
		var gm GameMessage
		err = json.Unmarshal(req.RequestData, &gm)
		if err != nil {
			// should tell sender request failed (new connection? just a response?)
			panic(
				fmt.Sprintf("request data bad json?\n%s\n%e", string(req.RequestData), err),
			)
		}

		s.sendGameMessage(game, req.RecipientName, gm)
	}
}

func otherName(name string) string {
	switch name {
	case "alice":
		return "bob"
	case "bob":
		return "alice"
	default:
		panic("who in the world?")
	}
}

func (s *Server) sendGameMessage(game Game, recipientName string, gm GameMessage) {
	fmt.Printf("sending to %s", recipientName)
	var address string
	switch recipientName {
	case game.Player1Address:
		address = game.Player1Address
	case game.Player2Name:
		address = game.Player2Address
	default:
		// tell sender player is not part of this game
		panic("lole")
	}

	var conn net.Conn
	var err error
	for {
		conn, err = net.Dial("tcp", address)
		if err == nil {
			break
		}
	}

	fmt.Printf("%v", gm)

	data, err := json.Marshal(gm)
	if err != nil {
		panic(err)
	}

	_, err = conn.Write(data)
	if err != nil {
		panic(err)
	}

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

// func otherAddress(name string) string {
// 	var other string
// 	switch name {
// 	case "alice":
// 		other = ":4040"
// 	case "bob":
// 		other = ":8080"
// 	}
// 	return other
// }

func main() {
	args := os.Args[1:]

	var address string
	var name string
	var piece Piece
	switch args[0] {
	case "x", "X":
		address = ":8080"
		name = "alice"
		piece = XPiece
	case "o", "O":
		address = ":4040"
		name = "bob"
		piece = OPiece
	default:
		NewRelay().startRelay()
	}

	if piece != NoPiece {
		p := tea.NewProgram(New(address, name, piece))

		go startServer(p, ":3030")

		if _, err := p.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}

	}

}
