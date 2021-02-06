import { ServerEvent, events, createEvent } from './events';
import { Screen, Map, Ability, setupCanvasResize } from './screen'


type Point = {
    x: number;
    y: number;
}

type Rectangle = {
    x: number;
    y: number;
    w: number;
    h: number;
}

export type Score = {
    kills: number;
    deaths: number;
    assists: number;
}

export type Team = {
    color: string;
    size: number;
}

export type Projectile = {
    team: string;
    r: number;
    x: number;
    y: number;
}

export type Wall = Rectangle;

export type Bush = Rectangle;

export class Champion {
    private static allyHealth: string = "#00ff00";
    private static enemyHealth: string = "#ff0000";
    private static clientHealth: string = "#ffff00";

    public id: string;
    public name: string;
    public x: number;
    public y: number;
    public r: number;
    public health: number;
    public maxHealth: number;


    constructor(id: string, name: string, health: number, x: number, y: number, r: number) {
        this.id = id;
        this.name = name;
        this.x = x;
        this.y = y;
        this.r = r;
        this.health = health;
        this.maxHealth = 100;
    }

    draw(screen: Screen, isAlly: boolean, isClient: boolean, teamColor: string) {
        screen.drawCircle(this.x, this.y, this.r, teamColor);
        // screen.drawRect(this.x, this.y, 50, 50, teamColor);

        let color: string;
        if (isClient) {
            color = Champion.clientHealth;
        } else if (isAlly) {
            color = Champion.allyHealth;
        } else {
            color = Champion.enemyHealth;
        }

        screen.drawRect(this.x - this.r, this.y - (this.r + 30), this.r * 2, 20, "#000000");
        screen.drawRect(this.x - this.r, this.y - (this.r + 30), this.r * 2 * (this.health / this.maxHealth), 20, color);
        screen.drawText(this.x - this.r, this.y - (this.r + 40), this.name);
    }
}

export class Game {
    private socket: WebSocket;
    private canvas: HTMLCanvasElement;
    private screen: Screen;
    private map: Map;
    private client: string;                     // client id
    private clients: { [key: string]: string }  // client id: team id
    private teams: { [key: string]: Team }      // team id: team info
    private cameraLockOn: boolean;

    public walls: Wall[];
    public bushes: Bush[];
    public champions: Champion[];
    public projectiles: Projectile[];
    public abilities: Ability[];

    constructor(setup: events.Setup, canvas: HTMLCanvasElement, socket: WebSocket) {
        this.socket = socket;
        this.canvas = canvas;
        this.screen = new Screen(canvas);
        this.map = new Map(canvas);
        this.client = setup.id;
        this.walls = setup.walls;
        this.bushes = setup.bushes;
        this.cameraLockOn = true;
        this.champions = [];
        this.projectiles = [];
        this.teams = {};
        this.clients = {};
        this.setOnClicks();
        this.abilities = [new Ability("Q", 0.1), new Ability("W", 2)];

        setupCanvasResize(canvas, [this.screen, Map, Ability])
    }
    /**
     * Processes ticks
     * @param t Incoming tick from the server
     */
    tick(t: events.Tick) {
        // Find the clients champion
        let client: Champion = <Champion>t.champions.find(champion => this.championIsClient(champion));

        // Check 
        let bush: Bush | undefined = this.bushes.find(bush => {
            return (bush.x < client.x) && (client.x < bush.x + bush.w) && (bush.y < client.y) && (client.y < bush.y + bush.h);
        })
        let stealthed = bush !== undefined;

        // Update information from the tick
        this.projectiles = t.projectiles;
        this.champions = t.champions.map(c => {
            return Object.assign(new Champion("", "", 0, 0, 0, 0), c)
        });

        if (this.cameraLockOn) {
            this.screen.center(client.x, client.y);
        }

        this.draw(stealthed);
    }

    updateTeams({ teams, clients }: events.TeamUpdate) {
        this.teams = teams;
        this.clients = clients;
    }

    draw(stealthed: boolean) {
        this.screen.clear();

        this.bushes.forEach(({ x, y, w, h }) => {
            this.screen.drawRect(x, y, w, h, "#32CD32");
        })

        this.walls.forEach(({ x, y, w, h }) => {
            this.screen.drawRect(x, y, w, h, "#8C8F8F");
        })

        this.projectiles.forEach(({ x, y, r, team }) => {
            this.screen.drawCircle(x, y, r, this.teams[team].color);
        })

        this.champions.forEach(champion => {
            let isAlly = this.championIsAlly(champion);
            let isClient = this.championIsClient(champion);
            let teamColor = this.championGetTeamColor(champion);

            champion.draw(this.screen, isAlly, isClient, teamColor);
        });

        if (stealthed) {
            this.screen.drawStealth();
        }

        this.abilities.forEach((ability, index) => {
            ability.tick();
            ability.draw(this.screen, index);
        })

        this.drawMap()
    }

    drawMap() {
        this.map.drawBackground();

        this.bushes.forEach(({ x, y, w, h }) => {
            this.map.drawRect(x, y, w, h, "#32CD32");
        })

        this.walls.forEach(({ x, y, w, h }) => {
            this.map.drawRect(x, y, w, h, "#8C8F8F");
        })

        this.champions.forEach(champion => {
            let teamColor = this.championGetTeamColor(champion);
            this.map.drawCircle(champion.x, champion.y, champion.r, teamColor);
        });

        this.map.drawBorder();
    }

    championIsClient(champion: Champion): boolean {
        return champion.id === this.client;
    }

    championIsAlly(champion: Champion): boolean {
        // Matching teams
        return this.teams[this.clients[champion.id]].color === this.teams[this.clients[this.client]].color;
    }

    championGetTeamColor(champion: Champion): string {
        // Matching teams
        return this.teams[this.clients[champion.id]].color;
    }

    private setOnClicks() {
        let x: number, y: number;
        let rect = this.canvas.getBoundingClientRect();

        // update position of the mouse
        window.onmousemove = (e: MouseEvent) => {
            x = Math.round(e.clientX - rect.left);
            y = Math.round(e.clientY - rect.top);
        }

        window.onkeydown = (e: KeyboardEvent) => {
            switch (e.keyCode) {
                case 89:
                    this.cameraLockOn = !this.cameraLockOn;
                    break;
                case 81:
                    this.abilities[0].start()
                    let [x_new, y_new] = this.screen.transformCanvasPosition(x, y);
                    this.socket.send(createEvent("game", "shoot", { x: x_new, y: y_new }));
                    break;
                case 87:
                    this.abilities[1].start()
                    this.socket.send(createEvent("game", "dash", {}))
                    break;
            }
        };

        // stop the context menu from showing up on right click
        this.canvas.addEventListener('contextmenu', e => {
            e.preventDefault()
        }, false)

        // clicks
        this.canvas.addEventListener('mousedown', e => {
            [x, y] = this.screen.transformCanvasPosition(x, y);
            this.socket.send(createEvent("game", "move", { x, y }));
        });

        // zoom in / out
        this.canvas.addEventListener('wheel', e => {
            e.preventDefault();
            e.stopPropagation();
            let zoom = this.screen.zoom + (-.15 * (e.deltaY / Math.abs(e.deltaY)));
            this.screen.zoom = Math.min(2, Math.max(0.35, + (zoom)));
        });
    }
}