import { Canvas, Context } from "../Canvas";
import { Rect, Renderer, RendererListener } from "../types";
import { COLOR_PALETTE } from "../utils";

export class WorldRenderer implements Renderer {
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
        if(!Context || !Canvas) return;
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
        Context.moveTo(x, y);
        Context.quadraticCurveTo(x+100, y - 40, x+200, y);
        Context.fill();
        Context.fillRect(0, y, Canvas.clientWidth, Canvas.clientHeight - y)
        Context.restore()
    }
}