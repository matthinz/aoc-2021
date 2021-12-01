let prevLine = "";
let prevValue;
let increases = 0;

process.stdin.on("data", (chunk) => {
  const lines = chunk.toString().split("\n");
  lines[0] = prevLine + lines[0];
  prevLine = lines.pop();
  lines.forEach(processLine);
});

process.stdin.on("close", () => {
  processLine(prevLine);
  console.log("\n\n %d increases", increases);
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
