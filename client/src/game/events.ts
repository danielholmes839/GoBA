import { Wall, Bush, Team, Champion, Projectile, Score } from './gameplay';

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
        scores: { [key: string]: Score };
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



