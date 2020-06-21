async function main() {
    await executeStep();
}

async function executeStep() {
    try {
        const result = await $.ajax({
            url: '/get-step',
            data: {},
            // success: success,
            dataType: 'json'
        });
        switch (result.Command) {
            case 'display':
                displayMessage(result.Message);
                break;
            case 'fill-boards':
                fillBoards(result.Boards);
                break;
            case 'the-end':
                displayMessage(result.Message);
                return;
        }
        setTimeout(executeStep, 3000)
    } catch (e) {
        $('#message').text('Unknown error occurred: ' + e.message);
    }
}

function fillBoards(boards) {
    fillBoard(boards.White.Board, '#whiteBoard', '#343a40');
    fillBoard(boards.Black.Board, '#blackBoard', 'white');
}

function fillBoard(board, selector, color) {
    const boardObject = $(selector).children('tbody');
    for (let x in board) {
        for (let y in board[x]) {
            if (+board[x][y] !== 0 && +board[x][y] !== 99) {
                const cellPosition = `tr:nth-child(${+y+2}) > td:nth-child(${+x+2})`;
                const cell = boardObject.find(cellPosition);
                cell.css({ backgroundColor: color });
            }
        }
    }
}

function displayMessage(message) {
    $('#message').text(message);
    clearMessage();
}

function clearMessage() {
    setTimeout(() => { $('#message').text(''); }, 2000);
}

$(document).ready(async () => {
    await main();
});
