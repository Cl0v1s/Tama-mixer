import { Entity, Physics, Point } from "../types";
import { PetEntity } from "./Entity";


const PET_GROUND = 190
const G = 2

export class PetPhysics implements Physics {

    private entity: Entity
    private vector: Point

    constructor(pet: PetEntity) {
        this.entity = pet
        this.vector = { x: 0, y: 0, t: 0}
    }

    Vector(): Point {
        return this.vector;
    }

    Entity(): Entity {
        return this.entity
    }

    public applyForce(v: Point) {
        this.vector.x += v.x;
        this.vector.y += v.y
    }

    private applyLimits() {
        if(this.entity.X() < 0) { 
            this.vector.x *= -1
            this.entity.Move(0, this.entity.Y())
        }
        if(this.entity.X() + this.entity.W() > 200) {
            this.vector.x *= -1
            this.entity.Move(200 - this.entity.W(), this.entity.Y())
        }
        if(this.entity.Y() + this.entity.H() > 200) {
            this.vector.y *= -1
            this.entity.Move(this.entity.X(), 200 - this.entity.W())
        }
        if(this.entity.Y() < 0) {
            this.vector.y *= -1
            this.entity.Move(this.entity.X(), 0)
        }
    }

    private applyGravity() {
        if(this.entity.Y() + this.entity.H() < PET_GROUND) this.vector.y += G
        else if(this.vector.y >= 0) {
            this.vector.y *= -0.5
        }
    }

    tick(alpha: number): void {
        this.applyLimits()
        this.applyGravity()
        if(Math.abs(this.vector.x) <= 1) this.vector.x = 0
        if(Math.abs(this.vector.y) <= 1) this.vector.y = 0
        this.entity.Move(this.entity.X() + this.vector.x * alpha, this.entity.Y() + this.vector.y * alpha)
    }
}