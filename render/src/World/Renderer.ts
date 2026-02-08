import { Canvas, Context } from "../Canvas";
import { Point, Rect, Renderer, RendererListener } from "../types";
import { COLOR_PALETTE, getYOffset } from "../utils";


const STARS_NUMBER = 10
const GRASS_NUMBER = 3


export class WorldRenderer implements Renderer {

    private stars?: Point[]
    private grass?: Point[]

    Subscribe(listener: RendererListener): void {
        throw new Error("Method not implemented.");
    }
    UnSubscribe(listener: RendererListener): void {
        throw new Error("Method not implemented.");
    }
    BoundingBox(): Rect {
        throw new Error("Method not implemented.");
    }
    Render(x: number, y: number, z: number): void {
        if (!Context || !Canvas) return;
        Context.save()
        var gradient = Context.createLinearGradient(0, 200, 0, 0);
        COLOR_PALETTE.sky.forEach((p) => {
            gradient.addColorStop(p.pos, p.color);
        })
        Context.fillStyle = gradient;
        Context.fillRect(0, 0, 200, 200);
        Context.lineWidth = 2
        Context.strokeStyle = COLOR_PALETTE.ground.stroke
        Context.fillStyle = COLOR_PALETTE.ground.fill
        Context.beginPath();
        Context.moveTo(0, y);
        Context.quadraticCurveTo(0 + 100, y - 40, x + 200, y);
        Context.fill();
        Context.fillRect(0, y, Canvas.clientWidth, Canvas.clientHeight - y)

        // stars management 
        Context.fillStyle = "white"
        if (!this.stars) {
            this.stars = []
            for (let i = 0; i < STARS_NUMBER; i++) {
                this.stars.push(
                    {
                        x: Math.floor(Math.random() * Canvas.clientWidth),
                        y: Math.floor(Math.random() * Math.min(y, Canvas.clientHeight / 2)),
                        t: Math.random() * 3 + 1,
                    }
                )
            }
        }
        this.stars.forEach((star) => {
            Context?.fillRect(star.x - 0.03 * star.t * x, star.y, star.t, star.t)
        })
        // grass management 
        Context.strokeStyle = "rgba(0, 150, 0, 0.3)"
        if (!this.grass) {
            this.grass = []
            for (let i = 0; i < GRASS_NUMBER; i++) {
                this.grass.push(
                    {
                        x: Math.floor(Math.random() * Canvas.clientWidth),
                        y: Math.floor(Math.random() * (Canvas.clientHeight - y)) + y,
                        t: 0,
                    }
                )
            }
        }
        this.grass.forEach((grass) => {
            Context?.beginPath()
            Context?.moveTo(grass.x - 3, grass.y -3)
            Context?.bezierCurveTo(grass.x - 3, grass.y - 3, grass.x, grass.y - 1.5, grass.x, grass.y)
            Context?.moveTo(grass.x + 3, grass.y - 3)
            Context?.bezierCurveTo(grass.x + 3, grass.y - 3, grass.x, grass.y + 1.5, grass.x, grass.y)
            Context?.stroke()
        })

        Context.restore()


    }
}