let socket = new WebSocket("ws://localhost:8080/ws");
console.log("Attempting Connection...");

let canvas = document.getElementById("game")
let ctx = canvas.getContext("2d")

socket.onopen = () => {
    console.log("Successfully Connected");
    socket.send("Hi From the Client!")
    socket.send(JSON.stringify({ "event": "move", data: { x: 10, y: 10 } }))
};

socket.onclose = event => {
    console.log("Socket Closed Connection: ", event);
};

socket.onerror = error => {
    console.log("Socket Error: ", error);
};

socket.onmessage = message => {
    message.data.text().then((text) => {
        data = JSON.parse(text)
        ctx.fillStyle = "#fafafa";
        ctx.fillRect(0, 0, canvas.width, canvas.height);

        for (let point of data.points) {
            ctx.beginPath();
            ctx.arc(point.x, point.y, 10, 0, 2 * Math.PI);
            ctx.stroke();
        }
    });
};
