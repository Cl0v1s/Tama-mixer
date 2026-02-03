import { Entity, Renderer } from "../types"
import { PetRenderer } from "./Renderer"

export class PetEntity implements Entity {
    private x: number
    private y: number
    private z: number

    private renderer = new PetRenderer()

    constructor(x: number, y: number, z: number) {
        this.x = x;
        this.y = y;
        this.z = z;
    }

    Renderer(): Renderer {
        return this.renderer
    }

    public X() {
        return this.x
    }

    public Y() {
        return this.y
    }

    public Z() {
        return this.z
    }
}