import { Entity } from "../types";
import { PetState } from "./Animations";
import { PetPhysics } from "./Physics";
import { PetRenderer } from "./Renderer";


export class PetEntity implements Entity {
    private x: number;
    private y: number;
    private z: number;

    private state: PetState = PetState.IDLE
    private renderer = new PetRenderer();
    private physics;

    constructor() {
        this.x = 0;
        this.y = 0;
        this.z = 1; 
        this.physics = new PetPhysics(this)
    }

    Move(x: number, y: number, z?: number): Entity {
        this.z = z === undefined ? this.z : z;
        this.x = x;
        this.y = y;
        return this;
    }

    W(): number {
        return this.renderer.BoundingBox().width * this.z
    }
    H(): number {
        return this.renderer.BoundingBox().height * this.z
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

    tick(): void {
        if(this.state === PetState.IDLE) {
            if(Math.floor(Math.random() * 10) === 0) {
                // this.physics.applyForce({ x: 2, y: 0, t: 0})
            }
        }
        if(Math.abs(this.physics.Vector().x) + Math.abs(this.physics.Vector().y) > 0) {
            this.state = PetState.WALKING
        } else {
            this.state = PetState.IDLE
        }
        if(this.physics.Vector().x > 0) {
            this.renderer.revert = true
        } else {
            this.renderer.revert = false
        }

        this.renderer.animationState = this.state
    }

    render(): void {
        this.physics.tick(1)
        this.tick()
        this.renderer.render(this.x ,this.y, this.z)
    }
}