import { StateMachine } from "../StateMachine";
import { Entity, Renderer, State } from "../types";
import { PetPhysics } from "./Physics";
import { PetRenderer } from "./Renderer";
import { AVAILABLE_PET_STATES, PET_STATES } from "./States";


export class PetEntity implements Entity {
    private x: number;
    private y: number;
    private z: number;

    private stateMachine: StateMachine

    public renderer: PetRenderer;
    public physics: PetPhysics;

    constructor() {
        this.x = 0;
        this.y = 0;
        this.z = 1; 
        this.physics = new PetPhysics(this)
        this.renderer = new PetRenderer();
        this.renderer.Subscribe(this);

        const s = Object.assign({}, PET_STATES.Idle)
        s.currentPet = this
        this.stateMachine = new StateMachine(s)
    }

    Destroy(): void {
        this.renderer.UnSubscribe(this)
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
        if(this.stateMachine.currentState.Update) this.stateMachine.currentState.Update();
        if(this.physics.Vector().x > 0) {
            this.renderer.revert = true
        } else {
            this.renderer.revert = false
        }
    }

    OnRenderer(renderer: Renderer): void {
        throw new Error("Method not implemented.");
    }

    Render(): void {
        this.physics.Tick(1)
        this.tick()
        this.renderer.Render(this.x ,this.y, this.z)
    }

    SetState(state: typeof PET_STATES[typeof AVAILABLE_PET_STATES[number]]): void {
        const s = Object.assign({}, state)
        s.currentPet = this
        this.stateMachine.Enter(s);
    }
}