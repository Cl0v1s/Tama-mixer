import { State } from "../types";
import { PET_ANIMATIONS } from "./Animations";
import { PetEntity } from "./Entity";

export const AVAILABLE_PET_STATES = ["Idle", "Walking", "Eating"] as const

type IPetState = State & {
    currentPet?: PetEntity,
}

export const PET_STATES: Record<typeof AVAILABLE_PET_STATES[number], IPetState> = {
    Idle: {
        key: "Idle",
        condition: "manual",
        next: [
            "Walking",
            "Eating"
        ],
        OnEnter: function () {
            this.currentPet?.renderer.Play(PET_ANIMATIONS.Idle)
        },
        Update: function () {
            if (Math.floor(Math.random() * 100) === 0) {
                this.currentPet?.SetState(PET_STATES.Walking)
            }
        }
    },
    Walking: {
        key: "Walking",
        currentPet: undefined,
        condition: "timeout",
        time: 2000,
        next: [
            "Idle"
        ],
        OnEnter: function () {
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

