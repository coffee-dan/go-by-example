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

type Model struct {
	server bool
	piece  Piece
	curX   int
	curY   int
	turn   Piece
	board  [3][3]Piece
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

// func handleConnection(conn net.Conn) {
// 	buf := make([]byte, 13)
// 	conn.Read(buf)
// 	fmt.Printf("%s", buf)
// 	conn.Close()
// }

// func startServer() {
// 	lnr, err := net.Listen("tcp", ":8080")
// 	if err != nil {
// 		panic(err)
// 	}

// 	defer lnr.Close()

// 	for {
// 		conn, err := lnr.Accept()
// 		if err != nil {
// 			panic(err)
// 		}
// 		go handleConnection(conn)
// 	}
// }

// func startClient() {
// 	var err error

// 	rdr := bufio.NewReader(os.Stdin)
// 	var msg string

// 	for {
// 		fmt.Printf("Huh?: ")
// 		msg, err = rdr.ReadString('\n')
// 		if err != nil {
// 			panic(err)
// 		}
// 		msg = strings.TrimSpace(msg)

// 		if len(msg) > 0 {
// 			conn, err := net.Dial("tcp", ":8080")
// 			if err != nil {
// 				panic(err)
// 			}

// 			fmt.Fprint(conn, msg)
// 		}
// 	}
// }

func New(isServer bool, pi Piece) *Model {
	var tr Piece
	if isServer {
		tr = pi
	} else {
		tr = otherPiece(pi)
	}

	return &Model{
		server: isServer,
		piece:  pi,
		turn:   tr,
		curX:   0,
		curY:   0,
		board: [3][3]Piece{
			{NoPieces, NoPieces, NoPieces},
			{NoPieces, NoPieces, NoPieces},
			{NoPieces, NoPieces, NoPieces},
		},
	}
}

func (m *Model) Init() tea.Cmd {
	m.startServer()

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
				m.board[m.curX][m.curY] = m.piece
				m.turn = otherPiece(m.piece)
				return m, m.gameTurnComplete
			}
		}
		return m, nil
	}

	return m, nil
}

func (m *Model) gameTurnComplete() tea.Msg {
	return GameMsg(GameTurnComplete)
}

// rdr := bufio.NewReader(os.Stdin)
// 	fmt.Print("What?: ")
// 	cmd, err := rdr.ReadString('\n')
// 	if err != nil {
// 		panic(err)
// 	}

// 	cmd = strings.TrimSpace(cmd)

// 	switch cmd {
// 	case "server", "s":
// 		startServer()
// 	case "client", "c":
// 		startClient()
// 	}

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

func (m *Model) handleConnection(conn net.Conn) {
	var buf []byte
	conn.Read(buf)

	var mv Move

	json.Unmarshal(buf, &mv)

	conn.Close()
}

func (m *Model) startServer() {
	lnr, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	defer lnr.Close()

	for {
		conn, err := lnr.Accept()
		if err != nil {
			panic(err)
		}
		go m.handleConnection(conn)
	}
}

func main() {
	args := os.Args[1:]

	var isServer bool
	switch args[0] {
	case "s", "server":
		isServer = true
	case "c", "client":
		isServer = false
	default:
		panic("unknown role")
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

	p := tea.NewProgram(New(isServer, piece))

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
