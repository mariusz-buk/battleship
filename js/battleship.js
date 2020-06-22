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

function zeroFill( number, width )
{
    width -= number.toString().length;
    if ( width > 0 )
    {
        return new Array( width + (/\./.test( number ) ? 2 : 1) ).join( '0' ) + number;
    }
    return number + ""; // always return a string
}

function displayMessage(message) {
    const msgField = $('#message');
    const paragraph = document.createElement('p');
    const date = new Date();
    const dateString = '[' + zeroFill(date.getHours(), 2) + ':'
        + zeroFill(date.getMinutes(), 2) + ':'
        + zeroFill(date.getSeconds(), 2) + '] ';
    const text = document.createTextNode(dateString + message);
    paragraph.append(text);
    msgField.prepend(paragraph);
}

$(document).ready(async () => {
    await main();
});
