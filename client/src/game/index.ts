import { ServerEvent, events } from '../events';

type Context2D = CanvasRenderingContext2D;

class Screen {
    private canvas: HTMLCanvasElement;
    private ctx: Context2D;
    centerX: number;
    centerY: number;
    dx: number;
    dy: number;
    zoom: number;

    constructor(canvas: HTMLCanvasElement) {
        this.canvas = canvas;
        this.ctx = <Context2D>this.canvas.getContext("2d", { alpha: false });
        this.centerX = this.canvas.width / 2;
        this.centerY = this.canvas.height / 2;

        this.zoom = 2;
        this.dx = 100;
        this.dy = 100;

    }

    clear() {
        this.ctx.fillStyle = "#ffffff";
        this.ctx.fillRect(0, 0, this.canvas.width, this.canvas.height);
    }

    center(x: number, y: number) {
        this.dx = x;
        this.dy = y;
    }

    /* Game position to canvas position */
    transformGamePosition(x: number, y: number): [number, number] {
        x = this.zoom * (x - this.dx) + this.centerX;
        y = this.zoom * (y - this.dy) + this.centerY;
        return [x, y]
    }

    /* Canvas position to game position */
    transformCanvasPosition(x: number, y: number): [number, number] {
        x = Math.round((x - this.centerX) / this.zoom) + this.dx;  //+ this.centerX;
        y = Math.round((y - this.centerY) / this.zoom) + this.dy; //+ this.centerY;
        return [x, y]
    }

    /* Draw rect with game position and scale */
    drawRect(x: number, y: number, w: number, h: number, color: string) {
        [x, y] = this.transformGamePosition(x, y);
        w *= this.zoom;
        h *= this.zoom;

        this.ctx.strokeStyle = color;
        this.ctx.fillStyle = color;
        this.ctx.fillRect(x, y, w, h);
    }

    /* Draw circle with game position and scale */
    drawCircle(x: number, y: number, r: number, color: string) {
        [x, y] = this.transformGamePosition(x, y);
        r *= this.zoom;

        this.ctx.strokeStyle = color;
        this.ctx.fillStyle = color;
        this.ctx.beginPath();
        this.ctx.arc(x, y, r, 0, 2 * Math.PI);
        this.ctx.stroke();
        this.ctx.fill();
    }

    /* Draw text with game position */
    drawText(x: number, y: number, text: string) {
        [x, y] = this.transformGamePosition(x, y);
        let size = Math.round(18 * this.zoom)
        this.ctx.font = `${size}px 'Courier New', monospace`;
        this.ctx.fillStyle = "#000000";
        this.ctx.fillText(text, x, y)
    }
}

export type Team = {
    color: string;
    size: number;
}

export type Wall = {
    x: number;
    y: number;
    w: number;
    h: number;
}

export class Champion {
    private static allyHealth: string = "#00ff00";
    private static enemyHealth: string = "#ff0000";
    private static clientHealth: string = "#ffff00";

    public id: string;
    public x: number;
    public y: number;
    public r: number = 50;
    public health: number;
    public maxHealth: number;


    constructor(id: string, health: number, x: number, y: number) {
        this.id = id;
        this.x = x;
        this.y = y;
        this.health = health;
        this.maxHealth = 100;
    }

    draw(screen: Screen, isAlly: boolean, isClient: boolean, teamColor: string) {
        screen.drawCircle(this.x, this.y, this.r, teamColor);

        let color: string;
        if (isClient) {
            color = Champion.clientHealth;
        } else if (isAlly) {
            color = Champion.allyHealth;
        } else {
            color = Champion.enemyHealth;
        }

        screen.drawRect(this.x - this.r, this.y - (this.r + 30), this.r * 2.5, 20, color)
        screen.drawText(this.x - this.r, this.y - (this.r + 40), this.id);
    }
}

export class Game {
    private canvas: HTMLCanvasElement;
    private socket: WebSocket;
    private screen: Screen;
    private client: string;                     // client id
    private clients: { [key: string]: string }  // client id: team id
    private teams: { [key: string]: Team }      // team id: team info

    private grid: any;
    private walls: Wall[];
    private champions: Champion[];

    constructor(setup: events.Setup, canvas: HTMLCanvasElement, socket: WebSocket) {
        this.canvas = canvas;
        this.socket = socket;
        this.screen = new Screen(canvas)
        this.client = setup.id;
        this.walls = setup.walls;
        this.champions = [];
        this.teams = {};
        this.clients = {};
        this.setOnClicks();
    }

    tick(t: events.Tick) {
        this.champions = t.champions.map(c => {
            if (this.championIsClient(c)) {
                this.screen.center(c.x, c.y);
            }
            return Object.assign(new Champion("", 0, 0, 0), c)
        });
        this.draw();
    }

    updateTeams({ teams, clients }: events.TeamUpdate) {
        this.teams = teams;
        this.clients = clients;
    }

    draw() {
        this.screen.clear();

        this.walls.forEach(({ x, y, w, h }) => {
            this.screen.drawRect(x, y, w, h, "#dddddd");
        })
        this.champions.forEach(champion => {
            let isAlly = this.championIsAlly(champion);
            let isClient = this.championIsClient(champion);
            let teamColor = this.championGetTeamColor(champion);

            champion.draw(this.screen, isAlly, isClient, teamColor);
        });
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
        this.canvas.addEventListener('contextmenu', e => {
            e.preventDefault()
        }, false)

        this.canvas.addEventListener('mousedown', e => {
            // e.preventDefault()
            // e.stopPropagation()
            const rect = this.canvas.getBoundingClientRect();
            let x = Math.round(e.clientX - rect.left);
            let y = Math.round(e.clientY - rect.top);
            [x, y] = this.screen.transformCanvasPosition(x, y);
            let update = { category: "game", event: "move-event", timstamp: Date.now(), data: { x, y } }
            this.socket.send(JSON.stringify(update))
        });

        this.canvas.addEventListener('wheel', e => {
            this.screen.zoom = Math.min(2, Math.max(0.5, this.screen.zoom + (e.deltaY / -1250)))
            console.log(e);
        });
    }
}