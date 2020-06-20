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
        switch (result.command) {
            case 'display':
                displayMessage(result.message);
                break;
        }
        setTimeout(executeStep, 3000)
    } catch {
        $('#message').text('Unknown error occurred');
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
