import { events, ServerEvent, TickCounter } from './events';
import { Game } from './game';

let socket = new WebSocket("ws://localhost:8080/ws")

socket.onopen = (): any => {
    console.log("Websocket opened!");
}

socket.onclose = (): any => {
    console.log("Websocket closed :(");
    socket.close();
}

socket.onerror = (ev: Event): any => {
    console.log(ev);
    socket.close();
}
// 
let game: Game;
let canvas = <HTMLCanvasElement>document.getElementById("canvas");
let ticks = new TickCounter("tps");

socket.onmessage = (message: MessageEvent<string>) => {
    let event: ServerEvent<any> = JSON.parse(message.data);
    switch (event.name) {
        case "setup":
            let setup: events.Setup = event.data;
            game = new Game(setup, canvas, socket);
            break;

        case "tick":
            let tick: events.Tick = event.data;
            game.tick(tick);
            ticks.update(event);
            break;

        case "update-teams":
            let update: events.TeamUpdate = event.data;
            game.updateTeams(update);
            break;

        default:
            console.log("EVENT NOT PROCESSED", event);
            break;
    }
}

