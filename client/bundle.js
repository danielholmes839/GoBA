(function(){function r(e,n,t){function o(i,f){if(!n[i]){if(!e[i]){var c="function"==typeof require&&require;if(!f&&c)return c(i,!0);if(u)return u(i,!0);var a=new Error("Cannot find module '"+i+"'");throw a.code="MODULE_NOT_FOUND",a}var p=n[i]={exports:{}};e[i][0].call(p.exports,function(r){var n=e[i][1][r];return o(n||r)},p,p.exports,r,e,n,t)}return n[i].exports}for(var u="function"==typeof require&&require,i=0;i<t.length;i++)o(t[i]);return o}return r})()({1:[function(require,module,exports){
let socket = new WebSocket("ws://localhost:8080/ws");
console.log("Attempting Connection...");

class HealthBar {
    constructor(health, maxHealth) {
        this.health = health
        this.maxHealth = maxHealth
    }

    draw(ctx, x, y, width, height) {
        ctx.fillStyle = "#000000";
        ctx.fillRect(x, y, width, height)
        ctx.fillStyle = "#32CD32";
        ctx.fillRect(x, y, width * (this.health / this.maxHealth), height)
    }
}

class Game {
    constructor() {
        this.teams = {};
        this.walls = [];
        this.champions = [];
    }

    draw(ctx) {
        ctx.fillStyle = "#fafafa";
        ctx.fillRect(0, 0, canvas.width, canvas.height);

        this.walls.forEach(wall => wall.draw(ctx))
        this.champions.forEach(champion => champion.draw(ctx))
        this.champions.filter(champion => champion.visible).forEach(champion => {
            this.champions.filter(champion => champion.visible).forEach(other => {
                if (other.visible) {
                    ctx.strokeStyle = '#ff0000';
                    ctx.beginPath();
                    ctx.moveTo(champion.x, champion.y);
                    ctx.lineTo(other.x, other.y);
                    ctx.stroke();
                }
            })
        })
    }
}

class Champion {
    constructor({ id, x, y, health, visible }) {
        this.id = id;
        this.x = x;
        this.y = y;
        this.visible = visible;
        this.health = new HealthBar(health, 100);
        this.radius = 30;
    }

    draw(ctx) {
        if (this.visible) {
            ctx.fillStyle = "#000000";
            ctx.fillText(this.id, this.x - this.radius, this.y - (this.radius + 20));
            ctx.strokeStyle = '#000000';
            ctx.beginPath();
            ctx.arc(this.x, this.y, this.radius, 0, 2 * Math.PI);
            ctx.stroke();

            this.health.draw(ctx, this.x - this.radius, this.y - (this.radius + 15), this.radius * 2, 10)
        }
    }
}
class Wall {
    constructor({ x, y, w, h }) {
        this.x = x
        this.y = y
        this.w = w
        this.h = h
    }

    draw(ctx) {
        ctx.fillStyle = "#000000";
        ctx.fillRect(this.x, this.y, this.w, this.h)
    }
}

class TPS {
    constructor() {
        this.ticks = 0;
        this.timestamp = 0;
        this.display = document.getElementById("tps");
    }

    tick(event) {
        if (event.timestamp != this.timestamp) {
            this.timestamp = event.timestamp;
            this.display.innerHTML = `TPS: ${this.ticks}`;
            this.ticks = 0;
        }
        this.ticks++;
    }
}


socket.onopen = () => {
    console.log("Successfully Connected");
};

socket.onclose = event => {
    console.log("Socket Closed Connection: ", event);
};

socket.onerror = error => {
    console.log("Socket Error: ", error);
};

let canvas = document.getElementById("game")
let ctx = canvas.getContext("2d")
ctx.font = "12px Arial";

let tps = new TPS();
let game = new Game();

socket.onmessage = message => {
    message = JSON.parse(message.data)
    console.log(message);

    switch (message.event) {
        case "walls":
            game.walls = message.data.map(wall => new Wall(wall));
        case "update":
            tps.tick(message)
            game.champions = message.data.champions.map(champion => new Champion(champion))
    }

    game.draw(ctx);
};

canvas.addEventListener('contextmenu', e => {
    e.preventDefault()
}, false)

canvas.addEventListener('mousedown', e => {
    e.preventDefault()
    e.stopPropagation()
    const rect = canvas.getBoundingClientRect()
    const x = Math.round(e.clientX - rect.left)
    const y = Math.round(e.clientY - rect.top)
    console.log(x, y)

    let update = { category: "game", event: "move-event", timstamp: Date.now(), data: { x, y } }
    console.log(update)
    socket.send(JSON.stringify(update))

    e.preventDefault()
    e.stopPropagation()

})
},{}]},{},[1]);
