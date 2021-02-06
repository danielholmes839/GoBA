import { ServerEvent} from './events';

export class TickCounter {
    ticks: number;
    timestamp: number;
    app: any;

    constructor(app: any) {
        this.ticks = 0;
        this.timestamp = 0;
        this.app = app;
    }

    update(event: ServerEvent<void>) {
        if (event.timestamp != this.timestamp) {
            this.app.updateTPS(this.ticks);
            this.timestamp = event.timestamp;
            this.ticks = 0;
        }
        this.ticks++;
    }
}