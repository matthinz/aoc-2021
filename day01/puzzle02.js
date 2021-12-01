const WINDOW_SIZE = 3;

function readInputInWindows(stream, handler, callback) {
  let currentWindow = [];
  let prevLine = "";

  stream.on("data", (chunk) => {
    const lines = chunk.toString().split("\n");
    lines[0] = prevLine + lines[0];
    prevLine = lines.pop();
    lines.forEach(processLine);
  });

  stream.on("end", () => {
    processLine(prevLine);
    callback();
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

process.stdin.resume();
process.stdin.setEncoding("utf8");

readInputInWindows(
  process.stdin,
  (win) => {
    const sum = win.reduce((prev, val) => prev + val, 0);
    const didIncrease = prevSum != null && sum > prevSum;
    const didDecrease = prevSum != null && sum < prevSum;

    if (didIncrease) {
      increases++;
      console.log("%s = %d (%s)", win.join(" + "), sum, "increase");
    } else if (didDecrease) {
      console.log("%s = %d (%s)", win.join(" + "), sum, "decrease");
    }

    prevSum = sum;
  },
  () => {
    console.log("\n\n\n%d increases", increases);
  }
);

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
