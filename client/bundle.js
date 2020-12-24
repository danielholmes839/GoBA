(function(){function r(e,n,t){function o(i,f){if(!n[i]){if(!e[i]){var c="function"==typeof require&&require;if(!f&&c)return c(i,!0);if(u)return u(i,!0);var a=new Error("Cannot find module '"+i+"'");throw a.code="MODULE_NOT_FOUND",a}var p=n[i]={exports:{}};e[i][0].call(p.exports,function(r){var n=e[i][1][r];return o(n||r)},p,p.exports,r,e,n,t)}return n[i].exports}for(var u="function"==typeof require&&require,i=0;i<t.length;i++)o(t[i]);return o}return r})()({1:[function(require,module,exports){
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

},{}]},{},[1]);
