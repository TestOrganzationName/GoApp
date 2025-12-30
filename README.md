# GoApp

A complete implementation of the classic board game Go, written in Go, available in both web and CLI versions.

## Features

- 9x9 board size
- Full Go rules implementation including:
  - Stone capture mechanics
  - Liberty detection
  - Suicide rule prevention
  - Pass system
- Two ways to play:
  - **Web Version**: Beautiful browser-based UI with click-to-play interface
  - **CLI Version**: Terminal-based with Unicode stones (● for Black, ○ for White)

## Installation

```bash
go mod download
```

## Usage

### Web Version (Recommended)

Run the web server:

```bash
go run main.go
```

Then open your browser and navigate to:
```
http://localhost:8080
```

The web version features:
- Visual board with click-to-play interface
- Beautiful gradient UI design
- Real-time game state updates
- Pass and reset buttons
- Responsive design for mobile devices

### CLI Version

Run the terminal game:

```bash
go run go.go
```

## How to Play

### Web Version
1. The game starts with an empty 9x9 board
2. Black plays first
3. Click on any intersection to place a stone
4. Click "Pass Turn" to skip your turn
5. Click "New Game" to reset the board
6. The game ends when both players pass consecutively

### CLI Version
1. The game starts with an empty 9x9 board
2. Black plays first
3. Enter moves in the format: `row col` (e.g., `3 4`)
4. Enter `pass` to pass your turn
5. Enter `quit` to exit the game
6. The game ends when both players pass consecutively

## Rules

- Stones are captured when they have no liberties (empty adjacent spaces)
- You cannot place a stone that would have no liberties (suicide rule)
- Players alternate turns between Black (●) and White (○)
- The game ends when both players pass

## Build

### Build Web Version

```bash
go build -o goapp-web main.go
```

Then run:

```bash
./goapp-web
```

Open your browser to `http://localhost:8080`

### Build CLI Version

```bash
go build -o goapp go.go
```

Then run:

```bash
./goapp
```

## Project Structure

```
GoApp/
├── main.go           # Web server and API endpoints
├── go.go             # CLI version
├── go.mod            # Go module configuration
├── README.md         # This file
├── .gitignore        # Git ignore rules
└── static/           # Web assets
    ├── index.html    # Web UI
    ├── styles.css    # Styling
    └── app.js        # Client-side JavaScript
```