export enum PetState {
    IDLE,
    EATING,
    WALKING
}

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



export const ANIMATIONS: Record<PetState, AnimationConfig> = {
    [PetState.IDLE]: {
        parts: [],
        speed: 1,
        loop: true,
        body: 100,
    },
    [PetState.WALKING]: {
        parts: ['leg1', 'leg2', 'arm1', 'arm2'],
        speed: 10,
        loop: true,
        body: 20
    },
    [PetState.EATING]: {
        parts: ['mouth'],
        speed: 20,
        loop: true,
        body: 0
    }
};
