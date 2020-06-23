let soundFire;
let soundHit;
let soundMissed;
let isDemoRunning = false;

async function main() {
    displayMessage('Welcome to Battleship demo.');
    await loadSoundEffects();
    displayMessage('Armies do not see each other of course. They will shoot blindly.');
    $('#startDemo').on('click', startOrPauseDemo);
}

async function startOrPauseDemo() {
    if (isDemoRunning) {
        isDemoRunning = false;
        $('#startDemo').text('Continue');
    } else {
        isDemoRunning = true;
        $('#startDemo').text('Pause');
        await executeStep();
    }
}

async function executeStep() {
    if (!isDemoRunning) {
        return;
    }

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
            case 'new-target':
                newTarget(result);
                break;
            case 'hit-report':
                hitReport(result);
                break;
            case 'the-end':
                displayMessage(result.Message);
                return;
        }
        setTimeout(executeStep, 3000)
    } catch (e) {
        displayMessage('Unknown error occurred: ' + e.message);
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

function getCell(board, x, y) {
    let selector;
    if (board === 'black') {
        selector = '#blackBoard';
    } else if (board === 'white') {
        selector = '#whiteBoard';
    }
    const boardObject = $(selector).children('tbody');
    const cellPosition = `tr:nth-child(${+y+2}) > td:nth-child(${+x+2})`;
    return boardObject.find(cellPosition);
}

function hitReport(report) {
    const x = report.X;
    const y = report.Y;
    const translatedX = report.TranslatedX;
    const translatedY = report.TranslatedY;
    const sank = report.Sank;
    const hit = report.Hit;
    const board = report.Board;

    let message = '';
    if (hit) {
        if (board === 'black') {
            message = 'Black ';
        } else if (board === 'white') {
            message = 'White ';
        }
        message += `ship at position ${translatedY}:${translatedX} has been hit. `;
        if (sank) {
            message += 'It sank.';
        }

        const cell = getCell(board, x, y);
        cell.css({
            backgroundImage: 'url(/images/explosion.gif)',
            backgroundPosition: 'center',
            backgroundSize: 'contain',
        });

        soundHit.play();
    } else {
        const cell = getCell(board, x, y);
        cell.html('&#x1F4A3;');

        message = 'Misses';

        soundMissed.play();
    }

    displayMessage(message);
}

function newTarget(report) {
    displayMessage(report.Message);

    soundFire.play();
}

async function loadSoundEffects() {
    soundFire = new Audio('/sounds/fire.ogg');
    await soundFire.load();
    soundMissed = new Audio('/sounds/missed.ogg');
    await soundMissed.load();
    soundHit = new Audio('/sounds/hit.ogg');
    await soundHit.load();
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
