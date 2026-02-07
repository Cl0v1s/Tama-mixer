import { AVAILABLE_PET_STATES } from "./States";

export type BodyPart = "mouth" | "leg1" | "leg2" | "arm1" | "arm2" | "eye1" | "eye2"


export type AnimationConfig = {
    /**
     * Moving parts in this animations
     */
    parts: Array<BodyPart>;
    /**
     * 1 = every tick, 4 = every 4 tick etc...
     */
    speed: number;
    /**
     * Does animation loops
     */
    loop: boolean;
    /**
     * Does animation allow animated body 0 -> no
     */
    body: number
};



export const PET_ANIMATIONS: Record<typeof AVAILABLE_PET_STATES[number], AnimationConfig> = {
    Idle: {
        parts: [],
        speed: 1,
        loop: true,
        body: 100,
    },
    Walking: {
        parts: ['leg1', 'leg2', 'arm1', 'arm2'],
        speed: 10,
        loop: true,
        body: 20
    },
    Eating: {
        parts: ['mouth'],
        speed: 20,
        loop: true,
        body: 0
    }
};
