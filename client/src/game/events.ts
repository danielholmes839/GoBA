import { Wall, Bush, Team, Champion, Projectile } from './gameplay';

export type ClientEvent<T> = {
    category: string;
    event: string;
    timestamp: number;
    data: T;
}

export type ServerEvent<T> = {
    subscription: string;
    name: string;
    timestamp: number;
    data: T;
}

export class TickCounter {
    ticks: number;
    timestamp: number;
    app: any;

    constructor(app: any) {
        this.ticks = 0;
        this.timestamp = 0;
        this.app = app;
        console.log(app);
    }

    update(event: ServerEvent<void>) {
        if (event.timestamp != this.timestamp) {
            this.app.updateTPS(`TPS: ${this.ticks}`);
            this.timestamp = event.timestamp;
            this.ticks = 0;
        }
        this.ticks++;
    }
}

export namespace events {
    export type Setup = {
        id: string;         // the clients id
        walls: Wall[];      // the walls
        bushes: Bush[];      // the walls
    }

    // Update of teams
    export type TeamUpdate = {
        clients: { [key: string]: string };   // client-id: team-id
        teams: { [key: string]: Team };     // team-id: team
    }

    export type Tick = {
        champions: Champion[];
        projectiles: Projectile[];
    }
}

export const createEvent = (category: string, event: string, data: any): string => {
    return JSON.stringify({
        category: category,
        event: event,
        timestamp: Date.now(),
        data: data
    });
}



