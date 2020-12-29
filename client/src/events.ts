import { Wall, Team, Champion } from './game';

export type ServerEvent<T> = {
    subscription: string;
    name: string;
    timestamp: number;
    data: T;
}

export class TickCounter {
    ticks: number;
    timestamp: number;
    paragraph: HTMLParagraphElement;

    constructor(id: string) {
        this.ticks = 0;
        this.timestamp = 0;
        this.paragraph = <HTMLParagraphElement> document.getElementById(id);
    }

    update(event: ServerEvent<void>) {
        if (event.timestamp != this.timestamp) {
            this.paragraph.innerHTML = `TPS: ${this.ticks}`;
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
    }

    // Update of teams
    export type TeamUpdate = {
        clients: { [key: string]: string };   // client-id: team-id
        teams: { [key: string]: Team };     // team-id: team
    }

    export type Tick = {
        champions: Champion[];
    }
}



