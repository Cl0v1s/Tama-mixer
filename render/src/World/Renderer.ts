import { Context } from "../Canvas";
import { Rect, Renderer, RendererListener } from "../types";

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
        if(!Context) return;
        Context.save()
        // var gradient = Context.createLinearGradient(0, 200, 0, 0);
        // gradient.addColorStop(0, "#7b2ca98e");
        // gradient.addColorStop(1, "#532d9f90");
        // Context.fillStyle = gradient;
        // Context.fillRect(0, 0, 200, 200);

        Context.lineWidth = 1
        Context.strokeStyle = "#0c3761"
        Context.beginPath();
        Context.moveTo(x, y);
        Context.quadraticCurveTo(x+100, y - 40, x+200, y);
        Context.stroke();
        Context.restore()
    }
}