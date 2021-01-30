type Context2D = CanvasRenderingContext2D;


class Canvas {
    canvas: HTMLCanvasElement;
    context: Context2D
    width: number;
    height: number;
    centerX: number;
    centerY: number;

    constructor(id: string) {
        this.canvas = <HTMLCanvasElement>document.getElementById(id);
        this.context = <Context2D>this.canvas.getContext('2d');
        this.width = this.canvas.getBoundingClientRect().width;
        this.height = this.canvas.getBoundingClientRect().height;
        this.centerX = this.width / 2;
        this.centerY = this.height / 2;

        window.onresize = this.update;
    }

    update() {
        this.width = this.canvas.getBoundingClientRect().width;
        this.height = this.canvas.getBoundingClientRect().height;
    }
}

export default Canvas