const WINDOW_SIZE = 3;

function readInput(stream, handler) {
    let currentWindow = [];

    return new Promise(resolve => {
        let prevLine = "";

        stream.on("data", (chunk) => {
            const lines = chunk.toString().split('\n');
            lines[0] = prevLine + lines[0];
            prevLine = lines.pop();
            lines.forEach(processLine);
        }); 

        stream.on("close", () => {
            processLine(prevLine);
            resolve();
        });
    });

    function processLine(line) {
        const value = parseInt(line, 10);
        if (isNaN(value)) {
            return;
        }
        currentWindow.push(value);
        if (currentWindow.length === WINDOW_SIZE) {
            handler(currentWindow);
            currentWindow.shift();
        }
    }
}

let increases = 0;
let prevSum;

readInput(process.stdin, (win) => {
    const sum = win.reduce((prev, val) => prev+val, 0);
    const didIncrease = prevSum != null && sum > prevSum;

    increases += didIncrease ? 1 : 0;

    console.log("%s = %d (%s)", win.join(' + '), sum, didIncrease ? "increase" : "decrease");

    prevSum = sum;
}).then(() => {
    console.log("\n\n\n%d increases", increases);
});


function processLine(line) {
    const value = parseInt(line, 10);
    if (isNaN(value)) {
        return;
    }

    if (prevValue == null) {
        console.log("%d (N/A - no previous measurement)", value);
    } else {
        const didIncrease = value > prevValue;
        if (didIncrease) {
            increases++;
            console.log("%d (increased)", value);
        } else {
            console.log("%d (decreased)", value);
        }
    }

    prevValue = value;
}
