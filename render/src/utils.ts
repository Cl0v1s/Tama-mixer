import { BodyFrame, PartFrame } from "./types"

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