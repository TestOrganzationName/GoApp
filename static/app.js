let boardState = null;
let boardSize = 9;

// Initialize the game
async function init() {
    await fetchBoard();
    renderBoard();
    setupEventListeners();
}

// Fetch the current board state from the server
async function fetchBoard() {
    try {
        const response = await fetch('/api/board');
        boardState = await response.json();
        boardSize = boardState.size;
        updateUI();
    } catch (error) {
        console.error('Error fetching board:', error);
    }
}

// Render the board
function renderBoard() {
    const boardElement = document.getElementById('board');
    boardElement.innerHTML = '';
    boardElement.style.gridTemplateColumns = `repeat(${boardSize}, 40px)`;
    boardElement.style.gridTemplateRows = `repeat(${boardSize}, 40px)`;

    for (let row = 0; row < boardSize; row++) {
        for (let col = 0; col < boardSize; col++) {
            const intersection = createIntersection(row, col);
            boardElement.appendChild(intersection);
        }
    }
}

// Create an intersection element
function createIntersection(row, col) {
    const intersection = document.createElement('button');
    intersection.className = 'intersection';
    intersection.dataset.row = row;
    intersection.dataset.col = col;

    // Add edge classes for border intersections
    if (row === 0) intersection.classList.add('edge-top');
    if (row === boardSize - 1) intersection.classList.add('edge-bottom');
    if (col === 0) intersection.classList.add('edge-left');
    if (col === boardSize - 1) intersection.classList.add('edge-right');

    // Add stone if present
    if (boardState && boardState.grid[row][col] !== 0) {
        const stone = document.createElement('div');
        stone.className = 'stone';
        stone.classList.add(boardState.grid[row][col] === 1 ? 'black' : 'white');
        intersection.appendChild(stone);
    }

    // Add click handler
    intersection.addEventListener('click', () => handleMove(row, col));

    return intersection;
}

// Handle a move
async function handleMove(row, col) {
    if (boardState.gameOver) {
        alert('Game is over! Start a new game.');
        return;
    }

    try {
        const response = await fetch('/api/move', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ row, col }),
        });

        const result = await response.json();

        if (result.success) {
            boardState = result.board;
            renderBoard();
            updateUI();

            if (result.gameOver) {
                setTimeout(() => {
                    alert('Game Over! Both players passed. Click "New Game" to play again.');
                }, 100);
            }
        } else {
            alert(result.message);
        }
    } catch (error) {
        console.error('Error making move:', error);
        alert('Error making move. Please try again.');
    }
}

// Handle pass
async function handlePass() {
    if (boardState.gameOver) {
        alert('Game is over! Start a new game.');
        return;
    }

    try {
        const response = await fetch('/api/pass', {
            method: 'POST',
        });

        const result = await response.json();
        boardState = result.board;
        renderBoard();
        updateUI();

        if (result.gameOver) {
            setTimeout(() => {
                alert('Game Over! Both players passed. Click "New Game" to play again.');
            }, 100);
        }
    } catch (error) {
        console.error('Error passing:', error);
    }
}

// Handle reset
async function handleReset() {
    if (confirm('Are you sure you want to start a new game?')) {
        try {
            const response = await fetch('/api/reset', {
                method: 'POST',
            });

            const result = await response.json();
            boardState = result.board;
            renderBoard();
            updateUI();
        } catch (error) {
            console.error('Error resetting game:', error);
        }
    }
}

// Update UI elements
function updateUI() {
    const playerIndicator = document.getElementById('currentPlayer');
    const gameStatus = document.getElementById('gameStatus');

    if (boardState) {
        const currentPlayer = boardState.turn === 1 ? 'Black' : 'White';
        playerIndicator.textContent = currentPlayer;
        playerIndicator.className = 'player-indicator ' + currentPlayer.toLowerCase();

        if (boardState.passes >= 2) {
            gameStatus.textContent = 'Game Over';
            gameStatus.style.color = '#dc3545';
        } else if (boardState.passes === 1) {
            gameStatus.textContent = 'Last player passed';
            gameStatus.style.color = '#ffc107';
        } else {
            gameStatus.textContent = 'Game in progress';
            gameStatus.style.color = '#667eea';
        }
    }
}

// Setup event listeners
function setupEventListeners() {
    document.getElementById('passBtn').addEventListener('click', handlePass);
    document.getElementById('resetBtn').addEventListener('click', handleReset);
}

// Start the game when the page loads
document.addEventListener('DOMContentLoaded', init);
