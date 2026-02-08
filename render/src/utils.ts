import { BodyFrame, PartFrame } from "./types"

export const COLOR_PALETTE = {
    pet: {
        stroke: "#004c84",
        fill1: "#fff79c",
        fill2: "#4192cd",
    },
    sky: [
        { pos: 0.00, color: "#c2f5ff" },
        { pos: 0.49, color: "#c2f5ff" },
        { pos: 0.50, color: "#80ECFF" },
        { pos: 0.69, color: "#80ECFF" },
        { pos: 0.70, color: "#00b2d6" },
        { pos: 0.84, color: "#00b2d6" },
        { pos: 0.85, color: "#0077b8" },
        { pos: 0.94, color: "#0077b8" },
        { pos: 0.95, color: "#000261" },
        { pos: 1.00, color: "#000261" }
    ],
    ground: {
        stroke: "#0A2A57",
        fill: "#C8F9ED"
    }
}

let BODIES: BodyFrame[][] = []
export async function LoadBodies() {
    // @ts-expect-error no typedef
    const promises = Object.values(import.meta.glob('./../../mixer/out/bodies/*.json')).map((m: any) => m().then((m: any) => m.default))
    BODIES = await Promise.all(promises)
}

let BODYPARTS: PartFrame[][] = []
export async function LoadBodyparts() {
    // @ts-expect-error no typedef
    const promises = Object.values(import.meta.glob('./../../mixer/out/bodyparts/*.json')).map((m: any) => m().then((m: any) => m.default))
    BODYPARTS = (await Promise.all(promises)).map((a: PartFrame[]) => a.sort((a, b) => a.frame - b.frame))
}

export function getBody() {
    return BODIES[Math.floor(Math.random() * BODIES.length)]
}


export function getPart(type: string, dont: string[] = []): PartFrame[] {
    const parts = BODYPARTS.filter((o) => o[0].type === type && dont.indexOf(o[0].name) === -1)
    return parts[Math.floor(Math.random() * parts.length)]
}

export function getOtherPart(part: PartFrame[]): PartFrame[] {
    const parts = BODYPARTS.filter((o) => o[0].type !== part[0].type)
    return parts.find((p) => p[0].name === part[0].name)!
}

export function getClosedEyeFrame() {
    return BODYPARTS.find((d) => d[0].name === "CLOSED")![0]
}

export function getYOffset(x: number) {
    return -0.4 * x + 0.002 * x * x;
}