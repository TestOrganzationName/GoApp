# GoApp

A terminal-based implementation of the classic board game Go, written in Go.

## Features

- 9x9 board size
- Full Go rules implementation including:
  - Stone capture mechanics
  - Liberty detection
  - Suicide rule prevention
  - Pass system
- Interactive command-line interface
- Clear visual board display with Unicode stones (● for Black, ○ for White)

## Installation

```bash
go mod download
```

## Usage

Run the game:

```bash
go run go.go
```

## How to Play

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

To build an executable:

```bash
go build -o goapp go.go
```

Then run:

```bash
./goapp
```