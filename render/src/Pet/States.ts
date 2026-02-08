import { State } from "../types";
import { PET_ANIMATIONS } from "./Animations";
import { PetEntity } from "./Entity";

export const AVAILABLE_PET_STATES = ["Idle", "Walking", "Eating", "Jumping"] as const

type IPetState = State & {
    currentPet?: PetEntity,
    [x: string]: unknown
}

export const PET_STATES: Record<typeof AVAILABLE_PET_STATES[number], IPetState> = {
    Idle: {
        key: "Idle",
        condition: "manual",
        next: [
            "Walking",
            "Eating",
            "Jumping"
        ],
        OnEnter: function () {
            this.currentPet?.renderer.Play(PET_ANIMATIONS.Idle)
        },
        Update: function () {
            if (Math.floor(Math.random() * 500) === 0) {
                this.currentPet?.SetState(PET_STATES.Jumping)
            } else if (Math.floor(Math.random() * 500) === 0) {
                this.currentPet?.SetState(PET_STATES.Walking)
            }
        }
    },
    Jumping: {
        key: "Jumping",
        currentPet: undefined,
        condition: "manual",
        timeout: undefined,
        jumped: false,
        next: [
            "Idle"
        ],
        OnEnter: function () {
            this.currentPet?.renderer.ApplyPose({ body: 1, leg1: 1, leg2: 1 })
            clearTimeout(this.timeout as number)
            this.timeout = setTimeout(() => {
                this.currentPet?.renderer.ApplyPose({ body: 0, leg1: 0, leg2: 0 })
                this.currentPet?.physics.ApplyForce({ x: 0, y: -10, t: 0 })
                this.jumped = true
            }, 1000)
        },
        Update: function() {
            if(!this.currentPet) return
            if(this.jumped && Math.abs(this.currentPet.physics.Vector().y) < 1) {
                this.currentPet?.SetState(PET_STATES.Idle)
            }
        },
        OnExit: function() {
            this.currentPet?.physics.Stop();
            clearTimeout(this.timeout as number)
            this.timeout = undefined
        }
    },
    Walking: {
        key: "Walking",
        currentPet: undefined,
        condition: "timeout",
        time: 2000, // this will be randomized on enter
        next: [
            "Idle"
        ],
        OnEnter: function () {
            this.time = Math.round(Math.random() * 3000) + 1000
            const dir = Math.round(Math.random()) === 1 ? -1 : 1
            this.currentPet?.renderer.Play(PET_ANIMATIONS.Walking)
            this.currentPet?.physics.ApplyForce({ x: 2 * dir, y: 0, t: 0 })
        },
        OnExit: function () {
            this.currentPet?.physics.Stop();
            this.currentPet?.SetState(PET_STATES.Idle)
        }
    },
    Eating: {
        key: "Eating",
        currentPet: undefined,
        condition: "timeout",
        time: 2000,
        next: [
            "Idle"
        ],
        OnEnter: function () {
            this.currentPet?.renderer.Play(PET_ANIMATIONS.Eating)
        },
        OnExit: function () {
            this.currentPet?.SetState(PET_STATES.Idle)
        }
    }
}

