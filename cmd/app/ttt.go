package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
	_ "github.com/mattn/go-sqlite3"
	"github.com/muesli/termenv"
)

const (
	host = "localhost"
	port = "2222"
)

func main() {
	store, err := NewSQLiteStore("tictactoe.db")
	if err != nil {
		log.Fatal("database error", "error", err)
	}
	defer store.Close()

	shared := NewSharedState(store)

	go shared.StartCleanupLoop(30 * time.Second)

	handler := func(sess ssh.Session) *tea.Program {
		userID := sess.User()
		sessID := sess.Context().Value(ssh.ContextKeySessionID).(string)
		pty, _, _ := sess.Pty()

		model := NewRootModel(shared, userID, sessID)

		opts := bubbletea.MakeOptions(sess)
		opts = append(opts, tea.WithAltScreen())
		p := tea.NewProgram(model, opts...)

		model.program = p
		shared.AddToLobby(sessID, userID, p)

		go func() {
			<-sess.Context().Done()
			shared.HandleDisconnect(sessID)
		}()

		go func() {
			time.Sleep(100 * time.Millisecond)
			p.Send(tea.WindowSizeMsg{Width: pty.Window.Width, Height: pty.Window.Height})
			p.Send(RoomListUpdateMsg{Rooms: shared.Rooms.List()})
		}()

		return p
	}

	s, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(".ssh/id_ed25519"),
		wish.WithMiddleware(
			bubbletea.MiddlewareWithProgramHandler(handler, termenv.ANSI256),
			activeterm.Middleware(),
			logging.Middleware(),
		),
	)
	if err != nil {
		log.Fatal("server creation failed", "error", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Info("Starting SSH server", "host", host, "port", port)
	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Error("server error", "error", err)
			done <- nil
		}
	}()

	<-done
	log.Info("Stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Error("shutdown error", "error", err)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// SQLite Store
// ─────────────────────────────────────────────────────────────────────────────

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(path string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS players (
			name TEXT PRIMARY KEY,
			wins INTEGER DEFAULT 0
		)
	`)
	if err != nil {
		return nil, err
	}

	return &SQLiteStore{db: db}, nil
}

func (s *SQLiteStore) RecordWin(name string) {
	_, _ = s.db.Exec(`
		INSERT INTO players (name, wins) VALUES (?, 1)
		ON CONFLICT(name) DO UPDATE SET wins = wins + 1
	`, name)
}

func (s *SQLiteStore) GetPlayerScore(name string) int {
	var wins int
	_ = s.db.QueryRow("SELECT wins FROM players WHERE name = ?", name).Scan(&wins)
	return wins
}

func (s *SQLiteStore) Close() error { return s.db.Close() }

// ─────────────────────────────────────────────────────────────────────────────
// Game Logic
// ─────────────────────────────────────────────────────────────────────────────

type GameState struct {
	Cells       [9]rune
	CurrentTurn rune
	MoveCount   int
}

func NewGameState() *GameState {
	var cells [9]rune
	for i := range cells {
		cells[i] = ' '
	}
	return &GameState{Cells: cells, CurrentTurn: 'X'}
}

func (g *GameState) MakeMove(position int) error {
	if position < 0 || position > 8 {
		return errors.New("position out of range")
	}
	if g.Cells[position] != ' ' {
		return errors.New("square already taken")
	}
	if g.IsOver() {
		return errors.New("game is over")
	}

	g.Cells[position] = g.CurrentTurn
	g.MoveCount++

	if g.CurrentTurn == 'X' {
		g.CurrentTurn = 'O'
	} else {
		g.CurrentTurn = 'X'
	}

	return nil
}

func (g *GameState) Winner() rune {
	lines := [8][3]int{
		{0, 1, 2},
		{3, 4, 5},
		{6, 7, 8},
		{0, 3, 6},
		{1, 4, 7},
		{2, 5, 8},
		{0, 4, 8},
		{2, 4, 6},
	}
	for _, line := range lines {
		a, b, c := g.Cells[line[0]], g.Cells[line[1]], g.Cells[line[2]]
		if a != ' ' && a == b && b == c {
			return a
		}
	}
	return ' '
}

func (g *GameState) IsDraw() bool {
	return g.MoveCount == 9 && g.Winner() == ' '
}

func (g *GameState) IsOver() bool {
	return g.Winner() != ' ' || g.MoveCount == 9
}

func (g *GameState) CurrentPlayerString() string {
	return string(g.CurrentTurn)
}

// ─────────────────────────────────────────────────────────────────────────────
// Player Roles
// ─────────────────────────────────────────────────────────────────────────────

type PlayerRole int

const (
	RoleSpectator PlayerRole = iota
	RolePlayerX
	RolePlayerO
)

func (r PlayerRole) String() string {
	switch r {
	case RolePlayerX:
		return "X"
	case RolePlayerO:
		return "O"
	default:
		return "Spectator"
	}
}

func (r PlayerRole) Mark() rune {
	switch r {
	case RolePlayerX:
		return 'X'
	case RolePlayerO:
		return 'O'
	default:
		return ' '
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Messages (cross-session communication)
// ─────────────────────────────────────────────────────────────────────────────

type RoomInfo struct {
	ID      string
	Players int
	Status  string
}

type (
	RoomListUpdateMsg struct{ Rooms []RoomInfo }
	JoinRoomMsg       struct{ RoomID string }
	LeaveRoomMsg      struct{}
	RoomChatMsg       struct{ Sender, Text string }
	PlayerJoinedMsg   struct {
		Name string
		Role PlayerRole
	}
	PlayerLeftMsg   struct{ Name string }
	RoleAssignedMsg struct{ Role PlayerRole }
	GameUpdateMsg   struct {
		Cells       [9]rune
		CurrentTurn string
		IsOver      bool
		Winner      string
	}
)

// ─────────────────────────────────────────────────────────────────────────────
// Client
// ─────────────────────────────────────────────────────────────────────────────

type Client struct {
	Program *tea.Program
	Role    PlayerRole
	UserID  string
}

// ─────────────────────────────────────────────────────────────────────────────
// Room
// ─────────────────────────────────────────────────────────────────────────────

type Room struct {
	mu      sync.RWMutex
	ID      string
	clients map[string]*Client
	game    *GameState
	started bool
	store   *SQLiteStore
}

func NewRoom(id string, store *SQLiteStore) *Room {
	return &Room{
		ID:      id,
		clients: make(map[string]*Client),
		game:    NewGameState(),
		store:   store,
	}
}

func (r *Room) Join(sessID, userID string, p *tea.Program) PlayerRole {
	r.mu.Lock()
	defer r.mu.Unlock()

	role := r.assignRole()

	r.clients[sessID] = &Client{Program: p, Role: role, UserID: userID}

	if role == RolePlayerO && !r.started {
		r.started = true
	}

	r.broadcastLocked(PlayerJoinedMsg{Name: userID, Role: role})
	go p.Send(RoleAssignedMsg{Role: role})

	if r.started {
		go p.Send(r.gameSnapshot())
	}

	return role
}

func (r *Room) assignRole() PlayerRole {
	hasX, hasO := false, false
	for _, c := range r.clients {
		if c.Role == RolePlayerX {
			hasX = true
		}
		if c.Role == RolePlayerO {
			hasO = true
		}
	}

	switch {
	case !hasX:
		return RolePlayerX
	case !hasO:
		return RolePlayerO
	default:
		return RoleSpectator
	}
}

func (r *Room) Leave(sessID string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	client, ok := r.clients[sessID]
	if !ok {
		return
	}

	delete(r.clients, sessID)
	r.broadcastLocked(PlayerLeftMsg{Name: client.UserID})
}

func (r *Room) HandleMove(sessID string, position int) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	client, ok := r.clients[sessID]
	if !ok || client.Role == RoleSpectator || !r.started || r.game.IsOver() {
		return false
	}

	if string(r.game.CurrentTurn) != client.Role.String() {
		return false
	}

	if err := r.game.MakeMove(position); err != nil {
		return false
	}

	if r.game.IsOver() {
		r.recordResult()
	}

	r.broadcastLocked(r.gameSnapshot())
	return true
}

func (r *Room) recordResult() {
	winner := r.game.Winner()
	if winner == ' ' {
		return
	}
	for _, c := range r.clients {
		if c.Role.Mark() == winner {
			r.store.RecordWin(c.UserID)
			break
		}
	}
}

func (r *Room) gameSnapshot() GameUpdateMsg {
	winnerStr := ""
	if w := r.game.Winner(); w != ' ' {
		winnerStr = string(w)
	}

	return GameUpdateMsg{
		Cells:       r.game.Cells,
		CurrentTurn: r.game.CurrentPlayerString(),
		IsOver:      r.game.IsOver(),
		Winner:      winnerStr,
	}
}

func (r *Room) BroadcastChat(sender, text string) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	r.broadcastLocked(RoomChatMsg{Sender: sender, Text: text})
}

func (r *Room) broadcastLocked(msg tea.Msg) {
	for _, c := range r.clients {
		p := c.Program
		go p.Send(msg)
	}
}

func (r *Room) PlayerCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.clients)
}

func (r *Room) Status() string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	switch {
	case r.game.IsOver():
		return "finished"
	case r.started:
		return "playing"
	default:
		return "waiting"
	}
}

func (r *Room) Info() RoomInfo {
	return RoomInfo{ID: r.ID, Players: r.PlayerCount(), Status: r.Status()}
}

// ─────────────────────────────────────────────────────────────────────────────
// Room Manager
// ─────────────────────────────────────────────────────────────────────────────

type RoomManager struct {
	mu    sync.RWMutex
	rooms map[string]*Room
	store *SQLiteStore
}

func NewRoomManager(store *SQLiteStore) *RoomManager {
	return &RoomManager{rooms: make(map[string]*Room), store: store}
}

func (rm *RoomManager) Create(id string) *Room {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	room := NewRoom(id, rm.store)
	rm.rooms[id] = room
	return room
}

func (rm *RoomManager) GetOrCreate(id string) *Room {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	if r, ok := rm.rooms[id]; ok {
		return r
	}
	room := NewRoom(id, rm.store)
	rm.rooms[id] = room
	return room
}

func (rm *RoomManager) List() []RoomInfo {
	rm.mu.RLock()
	defer rm.mu.RUnlock()
	list := make([]RoomInfo, 0, len(rm.rooms))
	for _, r := range rm.rooms {
		list = append(list, r.Info())
	}
	return list
}

func (rm *RoomManager) CleanupEmpty() {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	for id, r := range rm.rooms {
		if r.PlayerCount() == 0 {
			delete(rm.rooms, id)
		}
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Shared State (lobby + rooms)
// ─────────────────────────────────────────────────────────────────────────────

type LobbyPlayer struct {
	Program *tea.Program
	UserID  string
}

type SharedState struct {
	Rooms   *RoomManager
	lobbyMu sync.RWMutex
	lobby   map[string]*LobbyPlayer
}

func NewSharedState(store *SQLiteStore) *SharedState {
	return &SharedState{
		Rooms: NewRoomManager(store),
		lobby: make(map[string]*LobbyPlayer),
	}
}

func (s *SharedState) AddToLobby(sessID, userID string, p *tea.Program) {
	s.lobbyMu.Lock()
	defer s.lobbyMu.Unlock()
	s.lobby[sessID] = &LobbyPlayer{Program: p, UserID: userID}
}

func (s *SharedState) RemoveFromLobby(sessID string) {
	s.lobbyMu.Lock()
	defer s.lobbyMu.Unlock()
	delete(s.lobby, sessID)
}

func (s *SharedState) BroadcastLobby() {
	msg := RoomListUpdateMsg{Rooms: s.Rooms.List()}
	s.lobbyMu.RLock()
	defer s.lobbyMu.RUnlock()
	for _, lp := range s.lobby {
		p := lp.Program
		go p.Send(msg)
	}
}

func (s *SharedState) HandleDisconnect(sessID string) {
	s.RemoveFromLobby(sessID)

	s.Rooms.mu.RLock()
	rooms := make([]*Room, 0, len(s.Rooms.rooms))
	for _, r := range s.Rooms.rooms {
		rooms = append(rooms, r)
	}
	s.Rooms.mu.RUnlock()

	for _, room := range rooms {
		room.Leave(sessID)
	}

	s.BroadcastLobby()
}

func (s *SharedState) StartCleanupLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for range ticker.C {
		s.Rooms.CleanupEmpty()
		s.BroadcastLobby()
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Root Model (lobby ↔ room switching)
// ─────────────────────────────────────────────────────────────────────────────

type viewState int

const (
	viewLobby viewState = iota
	viewRoom
)

type rootModel struct {
	state   viewState
	lobby   lobbyModel
	room    *roomModel
	shared  *SharedState
	program *tea.Program
	userID  string
	sessID  string
	width   int
	height  int
}

func NewRootModel(shared *SharedState, userID, sessID string) *rootModel {
	return &rootModel{
		state:  viewLobby,
		lobby:  newLobbyModel(shared, userID),
		shared: shared,
		userID: userID,
		sessID: sessID,
	}
}

func (m rootModel) Init() tea.Cmd { return nil }

func (m rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.lobby.width, m.lobby.height = msg.Width, msg.Height
		if m.room != nil {
			m.room.width, m.room.height = msg.Width, msg.Height
		}
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	case JoinRoomMsg:
		return m.joinRoom(msg.RoomID)

	case LeaveRoomMsg:
		return m.leaveRoom()
	}

	switch m.state {
	case viewLobby:
		var cmd tea.Cmd
		m.lobby, cmd = m.lobby.Update(msg)
		return m, cmd

	case viewRoom:
		if m.room != nil {
			var cmd tea.Cmd
			*m.room, cmd = m.room.Update(msg)
			return m, cmd
		}
	}

	return m, nil
}

func (m rootModel) joinRoom(roomID string) (tea.Model, tea.Cmd) {
	room := m.shared.Rooms.GetOrCreate(roomID)
	m.shared.RemoveFromLobby(m.sessID)

	role := room.Join(m.sessID, m.userID, m.program)
	rm := newRoomModel(room, m.sessID, m.userID, role, m.shared, m.width, m.height)

	m.room = &rm
	m.state = viewRoom
	m.shared.BroadcastLobby()

	return m, m.room.Init()
}

func (m rootModel) leaveRoom() (tea.Model, tea.Cmd) {
	if m.room != nil {
		m.room.room.Leave(m.sessID)
		m.room = nil
	}

	m.state = viewLobby
	m.shared.AddToLobby(m.sessID, m.userID, m.program)
	m.lobby.rooms = m.shared.Rooms.List()
	m.shared.BroadcastLobby()

	return m, nil
}

func (m rootModel) View() string {
	if m.state == viewRoom && m.room != nil {
		return m.room.View()
	}
	return m.lobby.View()
}

// ─────────────────────────────────────────────────────────────────────────────
// Lobby Model
// ─────────────────────────────────────────────────────────────────────────────

type lobbyMode int

const (
	lobbyBrowse lobbyMode = iota
	lobbyCreate
)

type lobbyModel struct {
	rooms  []RoomInfo
	cursor int
	mode   lobbyMode
	input  textinput.Model
	shared *SharedState
	userID string
	width  int
	height int
}

func newLobbyModel(shared *SharedState, userID string) lobbyModel {
	ti := textinput.New()
	ti.Placeholder = "Room name..."
	ti.CharLimit = 20
	ti.Width = 20

	return lobbyModel{
		rooms:  shared.Rooms.List(),
		shared: shared,
		userID: userID,
		input:  ti,
	}
}

func (m lobbyModel) Init() tea.Cmd { return nil }

func (m lobbyModel) Update(msg tea.Msg) (lobbyModel, tea.Cmd) {
	switch msg := msg.(type) {
	case RoomListUpdateMsg:
		m.rooms = msg.Rooms
		if m.cursor >= len(m.rooms) && len(m.rooms) > 0 {
			m.cursor = len(m.rooms) - 1
		}
		return m, nil

	case tea.KeyMsg:
		if m.mode == lobbyCreate {
			return m.handleCreateInput(msg)
		}
		return m.handleBrowseInput(msg)
	}

	return m, nil
}

func (m lobbyModel) handleBrowseInput(msg tea.KeyMsg) (lobbyModel, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}

	case "down", "j":
		if m.cursor < len(m.rooms)-1 {
			m.cursor++
		}

	case "enter":
		if len(m.rooms) > 0 {
			id := m.rooms[m.cursor].ID
			return m, func() tea.Msg { return JoinRoomMsg{RoomID: id} }
		}

	case "c":
		m.mode = lobbyCreate
		m.input.Reset()
		m.input.Focus()
		return m, textinput.Blink
	}

	return m, nil
}

func (m lobbyModel) handleCreateInput(msg tea.KeyMsg) (lobbyModel, tea.Cmd) {
	switch msg.String() {
	case "enter":
		name := strings.TrimSpace(m.input.Value())
		if name != "" {
			m.shared.Rooms.Create(name)
			m.shared.BroadcastLobby()
			m.mode = lobbyBrowse
			return m, func() tea.Msg { return JoinRoomMsg{RoomID: name} }
		}

	case "esc":
		m.mode = lobbyBrowse
		return m, nil
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

// ── Lobby Styles ─────────────────────────────────────────────────────────────

var (
	lobbyTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("170")).
			MarginBottom(1)

	lobbyItemStyle    = lipgloss.NewStyle().PaddingLeft(2)
	lobbySelectedItem = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170")).Bold(true)

	lobbyStatusWaiting  = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	lobbyStatusPlaying  = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	lobbyStatusFinished = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

	lobbyHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).MarginTop(1)
)

func (m lobbyModel) View() string {
	var b strings.Builder

	b.WriteString(lobbyTitleStyle.Render("♟  TicTacToe Lobby"))
	b.WriteString("\n\n")

	if m.mode == lobbyCreate {
		b.WriteString("  Enter room name:\n\n")
		b.WriteString("  " + m.input.View() + "\n\n")
		b.WriteString(lobbyHelpStyle.Render("  enter: create  esc: cancel"))
		return b.String()
	}

	if len(m.rooms) == 0 {
		b.WriteString("  No rooms yet. Press 'c' to create one.\n")
	} else {
		for i, room := range m.rooms {
			cursor := "  "
			style := lobbyItemStyle
			if i == m.cursor {
				cursor = "▸ "
				style = lobbySelectedItem
			}

			status := renderStatus(room.Status)
			line := fmt.Sprintf("%s%s  [%d/2]  %s", cursor, room.ID, room.Players, status)
			b.WriteString(style.Render(line) + "\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(lobbyHelpStyle.Render("  ↑/↓: navigate  enter: join  c: create  ctrl+c: quit"))

	return b.String()
}

func renderStatus(status string) string {
	switch status {
	case "waiting":
		return lobbyStatusWaiting.Render(status)
	case "playing":
		return lobbyStatusPlaying.Render(status)
	default:
		return lobbyStatusFinished.Render(status)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Room Model (game board + chat, split pane)
// ─────────────────────────────────────────────────────────────────────────────

type focusPane int

const (
	paneGame focusPane = iota
	paneChat
)

type roomModel struct {
	room   *Room
	shared *SharedState
	sessID string
	userID string
	role   PlayerRole

	focus     focusPane
	cursorRow int
	cursorCol int

	cells       [9]rune
	currentTurn string
	gameOver    bool
	winner      string
	gameStarted bool

	chatViewport viewport.Model
	chatInput    textinput.Model
	chatLog      []string

	width  int
	height int
}

func newRoomModel(room *Room, sessID, userID string, role PlayerRole, shared *SharedState, w, h int) roomModel {
	vp := viewport.New(30, 10)
	vp.SetContent("Waiting for players...")

	ti := textinput.New()
	ti.Placeholder = "Chat..."
	ti.CharLimit = 200
	ti.Width = 28

	var cells [9]rune
	for i := range cells {
		cells[i] = ' '
	}

	focus := paneGame
	if role == RoleSpectator {
		focus = paneChat
		ti.Focus()
	}

	return roomModel{
		room: room, shared: shared,
		sessID: sessID, userID: userID, role: role,
		focus: focus, cells: cells,
		chatViewport: vp, chatInput: ti, chatLog: []string{},
		width: w, height: h,
	}
}

func (m roomModel) Init() tea.Cmd { return nil }

func (m roomModel) Update(msg tea.Msg) (roomModel, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.chatViewport.Width = 30
		m.chatViewport.Height = max(m.height-8, 5)

	case RoleAssignedMsg:
		m.role = msg.Role

	case GameUpdateMsg:
		m.cells = msg.Cells
		m.currentTurn = msg.CurrentTurn
		m.gameOver = msg.IsOver
		m.winner = msg.Winner
		m.gameStarted = true

	case RoomChatMsg:
		m.appendChat(fmt.Sprintf("%s: %s", msg.Sender, msg.Text))

	case PlayerJoinedMsg:
		m.appendChat(fmt.Sprintf("* %s joined as %s", msg.Name, msg.Role))

	case PlayerLeftMsg:
		m.appendChat(fmt.Sprintf("* %s left", msg.Name))

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			return m, func() tea.Msg { return LeaveRoomMsg{} }
		case "tab":
			m.toggleFocus()
			return m, nil
		}

		if m.focus == paneGame {
			return m.handleGameInput(msg)
		}
		return m.handleChatInput(msg)
	}

	return m, nil
}

func (m *roomModel) appendChat(line string) {
	m.chatLog = append(m.chatLog, line)
	m.chatViewport.SetContent(strings.Join(m.chatLog, "\n"))
	m.chatViewport.GotoBottom()
}

func (m *roomModel) toggleFocus() {
	if m.focus == paneGame {
		m.focus = paneChat
		m.chatInput.Focus()
	} else {
		m.focus = paneGame
		m.chatInput.Blur()
	}
}

func (m roomModel) handleGameInput(msg tea.KeyMsg) (roomModel, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.cursorRow = max(0, m.cursorRow-1)
	case "down", "j":
		m.cursorRow = min(2, m.cursorRow+1)
	case "left", "h":
		m.cursorCol = max(0, m.cursorCol-1)
	case "right", "l":
		m.cursorCol = min(2, m.cursorCol+1)
	case "enter", " ":
		if m.role != RoleSpectator && m.gameStarted && !m.gameOver {
			pos := m.cursorRow*3 + m.cursorCol
			m.room.HandleMove(m.sessID, pos)
		}
	}
	return m, nil
}

func (m roomModel) handleChatInput(msg tea.KeyMsg) (roomModel, tea.Cmd) {
	if msg.String() == "enter" {
		text := strings.TrimSpace(m.chatInput.Value())
		if text != "" {
			m.room.BroadcastChat(m.userID, text)
			m.chatInput.Reset()
		}
		return m, nil
	}

	var cmd tea.Cmd
	m.chatInput, cmd = m.chatInput.Update(msg)
	return m, cmd
}

// ── Room Styles ──────────────────────────────────────────────────────────────

var (
	boardBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(1, 2).
			MarginRight(2)

	chatBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(0, 1)

	cellDefault   = lipgloss.NewStyle().Width(3).Height(1).Align(lipgloss.Center)
	cellHighlight = lipgloss.NewStyle().Width(3).Height(1).Align(lipgloss.Center).
			Background(lipgloss.Color("62")).Foreground(lipgloss.Color("0"))

	markX = lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Bold(true)
	markO = lipgloss.NewStyle().Foreground(lipgloss.Color("4")).Bold(true)

	focusLabel   = lipgloss.NewStyle().Foreground(lipgloss.Color("170")).Bold(true)
	roomStatus   = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).MarginTop(1)
	roomHelpText = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

func (m roomModel) View() string {
	game := m.viewGamePanel()
	chat := m.viewChatPanel()

	panels := lipgloss.JoinHorizontal(lipgloss.Top, game, chat)
	status := m.viewStatusBar()
	help := m.viewHelp()

	return lipgloss.JoinVertical(lipgloss.Left, panels, status, help)
}

func (m roomModel) viewGamePanel() string {
	var b strings.Builder

	if m.focus == paneGame {
		b.WriteString(focusLabel.Render("▸ GAME") + "\n\n")
	} else {
		b.WriteString("  GAME\n\n")
	}

	for row := 0; row < 3; row++ {
		var rendered []string
		for col := 0; col < 3; col++ {
			idx := row*3 + col
			ch := m.cells[idx]

			display := " "
			if ch == 'X' {
				display = markX.Render("X")
			} else if ch == 'O' {
				display = markO.Render("O")
			}

			isCursor := row == m.cursorRow && col == m.cursorCol && m.focus == paneGame
			if isCursor {
				rendered = append(rendered, cellHighlight.Render(display))
			} else {
				rendered = append(rendered, cellDefault.Render(display))
			}
		}

		b.WriteString("  " + strings.Join(rendered, "│") + "\n")
		if row < 2 {
			b.WriteString("  ───┼───┼───\n")
		}
	}

	return boardBorder.Render(b.String())
}

func (m roomModel) viewChatPanel() string {
	var b strings.Builder

	if m.focus == paneChat {
		b.WriteString(focusLabel.Render("▸ CHAT") + "\n")
	} else {
		b.WriteString("  CHAT\n")
	}

	b.WriteString(m.chatViewport.View() + "\n")
	b.WriteString(m.chatInput.View())

	return chatBorder.Render(b.String())
}

func (m roomModel) viewStatusBar() string {
	parts := []string{fmt.Sprintf("You: %s", m.role)}

	switch {
	case !m.gameStarted:
		parts = append(parts, "Waiting for opponent...")
	case m.gameOver && m.winner != "":
		parts = append(parts, fmt.Sprintf("Winner: %s!", m.winner))
	case m.gameOver:
		parts = append(parts, "Draw!")
	default:
		parts = append(parts, fmt.Sprintf("Turn: %s", m.currentTurn))
	}

	return roomStatus.Render(strings.Join(parts, "  "))
}

func (m roomModel) viewHelp() string {
	if m.focus == paneGame {
		return roomHelpText.Render("↑/↓/←/→: move  enter: place  tab: chat  esc: leave")
	}
	return roomHelpText.Render("type to chat  enter: send  tab: game  esc: leave")
}
