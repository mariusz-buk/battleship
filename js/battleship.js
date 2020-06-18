async function main() {
    await createConnectionToken();
}

async function createConnectionToken() {
    const result = await $.get('http://battleship.local:8080/getConntectionToken');
    console.log(result);
}

$(document).ready(async () => {
    await main();
});
