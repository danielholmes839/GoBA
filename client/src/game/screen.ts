type Context2D = CanvasRenderingContext2D;

interface CanvasResize {
    onCanvasResize: (width: number, height: number) => void
}

export const setupCanvasResize = (canvas: HTMLCanvasElement, observers: CanvasResize[]) => {
    const notify = () => {
        for (let observer of observers) {
            canvas.width = canvas.getBoundingClientRect().width;
            canvas.height = canvas.getBoundingClientRect().height;
            observer.onCanvasResize(canvas.width, canvas.height);
        }
    }

    window.onresize = notify
    notify()
}

export class Screen {
    public canvas: HTMLCanvasElement;
    public ctx: Context2D;
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

        this.zoom = 0.5;
        this.dx = 100;
        this.dy = 100;

    }

    onCanvasResize(width: number, height: number) {
        this.centerX = width / 2;
        this.centerY = height / 2;
    }

    clear() {
        this.ctx.fillStyle = "#f5f5dc";
        this.ctx.fillRect(0, 0, this.canvas.width, this.canvas.height);
    }

    drawStealth() {
        this.drawRect(0, 0, this.canvas.width, this.canvas.height, "#000000", 0.1, false)
    }

    center(x: number, y: number) {
        this.dx = x;
        this.dy = y;
    }

    // Convert the position in the game to the position on the canvas
    transformGamePosition(x: number, y: number): [number, number] {
        x = this.zoom * (x - this.dx) + this.centerX;
        y = this.zoom * (y - this.dy) + this.centerY;
        return [x, y]
    }

    // Convert the position on the canvas to the position in the game
    transformCanvasPosition(x: number, y: number): [number, number] {
        x = Math.round((x - this.centerX) / this.zoom) + this.dx;
        y = Math.round((y - this.centerY) / this.zoom) + this.dy;
        return [x, y]
    }

    // Draw a rectangle on the canvas
    drawRect(x: number, y: number, w: number, h: number, color: string, alpha: number = 1, transform: boolean = true) {
        if (transform) {
            [x, y] = this.transformGamePosition(x, y);
            w *= this.zoom;
            h *= this.zoom;
        }

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
        let size = Math.round(30 * this.zoom)
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

export class Map {
    private ctx: Context2D
    private zoom: number;
    private static levelSize = 3000;        // game units. distance from 0
    public static size: number = 200;                   // pixels
    public static x: number = 0;
    public static y: number = 0;

    constructor(canvas: HTMLCanvasElement) {
        this.zoom = Map.size / Map.levelSize;
        this.ctx = <Context2D>canvas.getContext("2d");
    }

    drawBackground() {
        this.ctx.fillStyle = "#f5f5dc";
        this.ctx.fillRect(Map.x, Map.y, Map.size, Map.size);
    }

    drawBorder() {
        this.ctx.strokeStyle = "#000000";
        this.ctx.strokeRect(Map.x, Map.y, Map.size, Map.size);
    }

    /* Draw rect with game position and scale */
    drawRect(x: number, y: number, w: number, h: number, color: string) {
        x = Map.x + (x * this.zoom);
        y = Map.y + (y * this.zoom);

        if (x < Map.x || y < Map.y || x > Map.x + Map.size || y > Map.y + Map.size) {
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
        x = Map.x + (x * this.zoom);
        y = Map.y + (y * this.zoom);

        if (x < Map.x || y < Map.y || x > Map.x + Map.size || y > Map.y + Map.size) {
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

    static onCanvasResize(width: number, height: number) {
        Map.x = width - Map.size;
        Map.y = height - Map.size;
    }
}

export class Ability {
    private static x: number = 20;
    private static y: number = 0;
    private static size: number = 50;
    private static separation: number = 20;

    public name: string;
    public maxTicks: number;
    public currentTicks: number;

    constructor(name: string, cooldownSeconds: number) {
        let tps = 64
        this.name = name;
        this.maxTicks = Math.round(tps * cooldownSeconds);
        this.currentTicks = this.maxTicks;
    }

    tick() {
        this.currentTicks = Math.min(this.maxTicks, this.currentTicks + 1);
    }

    start() {
        if (this.currentTicks === this.maxTicks) {
            this.currentTicks = 0;
        }
    }

    draw(screen: Screen, i: number) {
        let x = (i * (Ability.size + Ability.separation)) + Ability.x;
        let percent = this.currentTicks / this.maxTicks;
        let offset = Ability.size * (1 - percent);

        screen.drawRect(x, Ability.y, Ability.size, Ability.size, "#000000", 1, false);
        screen.drawRect(x, Ability.y + offset, Ability.size, Ability.size * percent, "#FFFF00", 1, false);
        screen.ctx.font = `40px 'Courier New', monospace`;
        screen.ctx.strokeStyle = "#333333";
        screen.ctx.fillStyle = "#333333";
        screen.ctx.fillText(this.name, x + 10, Ability.y + 35)
    }

    static onCanvasResize(width: number, height: number) {
        Ability.y = height - (Ability.separation + Ability.size);
    }
}