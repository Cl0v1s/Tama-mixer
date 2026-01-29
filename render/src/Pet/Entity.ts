import { Entity } from "../types";
import { PetRenderer } from "./Renderer";


export class PetEntity implements Entity {
    private x: number;
    private y: number;
    private z: number;

    private renderer = new PetRenderer();

    constructor(x: number, y: number) {
        this.x = x;
        this.y = y;
        this.z = 1; 
    }

    X(): number {
        return this.x
    }

    Y(): number {
        return this.y
    }

    Z(): number {
        return this.z
    }

    render(): void {
        this.renderer.render(this.x ,this.y,this.z)
    }
}