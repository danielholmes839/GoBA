import Vue from 'vue';
import { events, ServerEvent, TickCounter } from './events';
import { Game } from './gameplay';

export const setup = (socketUrl: string, canvasId: string, app: any) => {
    let game: Game;
    let canvas = <HTMLCanvasElement>document.getElementById(canvasId);
    let ticks = new TickCounter(app);
    let socket = new WebSocket(socketUrl)

    // Setup Events
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
                app.updateTeams(update.teams);
                break;
    
            default:
                console.log("EVENT NOT PROCESSED", event);
                break;
        }
    }
}
