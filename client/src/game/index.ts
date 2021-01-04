import { ServerEvent, events, createEvent } from '../events';

type Context2D = CanvasRenderingContext2D;

type Rectangle = {
    x: number;
    y: number;
    w: number;
    h: number;
}
export type Team = {
    color: string;
    size: number;
}

export type Wall = Rectangle;

export type Bush = Rectangle;

class Screen {
    private canvas: HTMLCanvasElement;
    private ctx: Context2D;
    public centerX: number;
    public centerY: number;
    public dx: number;
    public dy: number;
    public zoom: number;

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
        this.ctx.fillStyle = "#f5f5dc";
        this.ctx.fillRect(0, 0, this.canvas.width, this.canvas.height);
    }

    drawStealth() {
        this.ctx.fillStyle = "#000000";
        this.ctx.globalAlpha = 0.1;
        this.ctx.fillRect(0, 0, this.canvas.width, this.canvas.height);
        this.ctx.globalAlpha = 1;
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
    drawRect(x: number, y: number, w: number, h: number, color: string, alpha: number = 1) {
        [x, y] = this.transformGamePosition(x, y);
        w *= this.zoom;
        h *= this.zoom;

        this.ctx.strokeStyle = color;
        this.ctx.globalAlpha = alpha;
        this.ctx.fillStyle = color;
        this.ctx.fillRect(x, y, w, h);
        this.ctx.globalAlpha = 1;
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
        this.ctx.strokeStyle = "#000000";
        this.ctx.fillStyle = "#000000";
        this.ctx.fillText(text, x, y)
    }

    drawLine(x1: number, y1: number, x2: number, y2: number, color: string) {
        [x1, y1] = this.transformGamePosition(x1, y1);
        [x2, y2] = this.transformGamePosition(x2, y2);

        this.ctx.strokeStyle = color;
        this.ctx.beginPath();
        this.ctx.moveTo(x1, y1);
        this.ctx.lineTo(x2, y2);
        this.ctx.stroke();
    }
}

class Map {
    private static levelSize = 4000;        // game units. distance from 0

    private ctx: Context2D
    private x: number;
    private y: number;
    private mapSize: number;                // pixels
    private zoom: number;

    constructor(x: number, y: number, mapSize: number, canvas: HTMLCanvasElement) {
        this.x = x;
        this.y = y;
        this.mapSize = mapSize;
        this.zoom = this.mapSize / Map.levelSize;
        this.ctx = <Context2D> canvas.getContext("2d");
    }

    draw(game: Game) {
        // MINI MAP
        this.ctx.fillStyle = "#f5f5dc";
        this.ctx.fillRect(this.x, this.y, this.mapSize, this.mapSize);

        game.bushes.forEach(({ x, y, w, h }) => {
            this.drawRect(x, y, w, h, "#32CD32");
        })

        game.walls.forEach(({ x, y, w, h }) => {
            this.drawRect(x, y, w, h, "#8C8F8F");
        })

        game.champions.forEach(champion => {
            let teamColor = game.championGetTeamColor(champion);
            this.drawCircle(champion.x, champion.y, champion.r, teamColor);
        });

        this.ctx.strokeStyle = "#000000";
        this.ctx.strokeRect(this.x, this.y, this.mapSize, this.mapSize);
        
    }
    /* Draw rect with game position and scale */
    drawRect(x: number, y: number, w: number, h: number, color: string) {
        x = this.x + (x * this.zoom);
        y = this.y + (y * this.zoom);

        if (x < this.x || y < this.y || x > this.x + this.mapSize || y > this.y + this.mapSize) {
            return;
        }

        w *= this.zoom;
        h *= this.zoom;

        this.ctx.strokeStyle = color;
        this.ctx.fillStyle = color;
        this.ctx.fillRect(x, y, w, h);
    }

    /* Draw circle with game position and scale */
    drawCircle(x: number, y: number, r: number, color: string) {
        x = this.x + (x * this.zoom);
        y = this.y + (y * this.zoom);
        
        if (x < this.x || y < this.y || x > this.x + this.mapSize || y > this.y + this.mapSize) {
            return;
        }

        r *= this.zoom;

        this.ctx.strokeStyle = color;
        this.ctx.fillStyle = color;
        this.ctx.beginPath();
        this.ctx.arc(x, y, r, 0, 2 * Math.PI);
        this.ctx.stroke();
        this.ctx.fill();
    }
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
        // screen.drawRect(this.x, this.y, 50, 50, teamColor);

        let color: string;
        if (isClient) {
            color = Champion.clientHealth;
        } else if (isAlly) {
            color = Champion.allyHealth;
        } else {
            color = Champion.enemyHealth;
        }

        screen.drawRect(this.x - this.r, this.y - (this.r + 30), this.r * 2.5, 20, color);
        screen.drawText(this.x - this.r, this.y - (this.r + 40), this.id);
    }
}

export class Game {
    private canvas: HTMLCanvasElement;
    private socket: WebSocket;
    private screen: Screen;
    private map: Map;
    private client: string;                     // client id
    private clients: { [key: string]: string }  // client id: team id
    private teams: { [key: string]: Team }      // team id: team info
    private cameraLockOn: boolean;

    public walls: Wall[];
    public bushes: Bush[];
    public champions: Champion[];

    constructor(setup: events.Setup, canvas: HTMLCanvasElement, socket: WebSocket) {
        this.canvas = canvas;
        this.socket = socket;
        this.screen = new Screen(canvas);
        this.map = new Map(canvas.width-200, canvas.height-200, 200, canvas);
        this.client = setup.id;
        this.walls = setup.walls;
        this.bushes = setup.bushes;
        this.cameraLockOn = true;
        this.champions = [];
        this.teams = {};
        this.clients = {};
        this.setOnClicks();
    }

    tick(t: events.Tick) {
        let client: Champion = <Champion>t.champions.find(champion => this.championIsClient(champion));
        let bush: Bush | undefined = this.bushes.find(bush => {
            return (bush.x < client.x) && (client.x < bush.x + bush.w) && (bush.y < client.y) && (client.y < bush.y + bush.h);
        })
        let stealthed = bush !== undefined;

        this.champions = t.champions.map(c => Object.assign(new Champion("", 0, 0, 0), c));

        if (this.cameraLockOn) {
            this.screen.center(client.x, client.y);
        }

        this.draw(stealthed);
        this.map.draw(this);
        // 
        this.champions.forEach(({ x, y }) => {
            this.screen.drawLine(x, y, client.x, client.y, "#000000")
        })


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

        this.champions.forEach(champion => {
            let isAlly = this.championIsAlly(champion);
            let isClient = this.championIsClient(champion);
            let teamColor = this.championGetTeamColor(champion);

            champion.draw(this.screen, isAlly, isClient, teamColor);
        });

        if (stealthed) {
            this.screen.drawStealth();
        }



        // this.fog.draw();
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
        window.onmousemove = (e: MouseEvent ) => {
            x = Math.round(e.clientX - rect.left);
            y = Math.round(e.clientY - rect.top);
        }

        window.onkeydown = (e: KeyboardEvent) => {
            console.log(x, y);
            switch (e.keyCode) {
                case 89:
                    this.cameraLockOn = !this.cameraLockOn
                case 81:
                    this.socket.send(createEvent("game", "q", { x, y }));
                case 87:
                    this.socket.send(createEvent("game", "w", { x, y }));
                case 69:
                    this.socket.send(createEvent("game", "e", { x, y }));
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
            let zoom = this.screen.zoom + (-.25 * (e.deltaY / Math.abs(e.deltaY)));
            this.screen.zoom = Math.min(2, Math.max(0.50, + (zoom)));
        });


    }
}