import { events, ServerEvent, TickCounter } from './events';
import { Game } from './gameplay';

export const setup = (socket: WebSocket, app: any) => {
    let game: Game;
    let ticks = new TickCounter(app);

    // Setup Events
    socket.onclose = () => {
        location.reload();
    }

    socket.onmessage = (message: MessageEvent<string>) => {
        let event: ServerEvent<any> = JSON.parse(message.data);
        let canvas: HTMLCanvasElement = <HTMLCanvasElement>document.getElementById("canvas");

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
                break;
        }
    }
}
