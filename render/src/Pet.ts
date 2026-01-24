import Eggplant from './../../mixer/out/bodies/eggplant.json'
import PeanutMouth from './../../mixer/out/bodyparts/mouth-peanutmouth.json'
import { Context } from './Canvas'
import { BodyFrame, PartFrame, Point } from './types'


let BODIES: BodyFrame[][] = []
export async function LoadBodies() {
    const promises = Object.values(import.meta.glob('./../../mixer/out/bodies/*.json')).map((m: any) => m().then((m: any) => m.default))
    BODIES = await Promise.all(promises)
}

let BODYPARTS: PartFrame[][] = []
export async function LoadBodyparts() {
    const promises = Object.values(import.meta.glob('./../../mixer/out/bodyparts/*.json')).map((m: any) => m().then((m: any) => m.default))
    BODYPARTS = await Promise.all(promises)
}


function getPart(type: string): PartFrame {
    const parts = BODYPARTS.map((o) => o[0]).filter((o) => o.type === type)
    return parts[Math.floor(Math.random() * parts.length)]
}

function getOtherPart(part: PartFrame): PartFrame {
    const parts = BODYPARTS.map((o) => o[0]).filter((o) => o.type !== part.type)
    return parts.find((p) => p.name === part.name)!
}

export class Pet {

    private body: BodyFrame;
    public mouth: PartFrame;
    public leg1: PartFrame;
    public leg2: PartFrame;
    public arm1: PartFrame;
    public arm2: PartFrame;
    public eye1: PartFrame;
    public eye2: PartFrame;

    public zoom = 3;
    public x = 10;
    public y = 10;

    private bodyPath: Path2D
    private outPath: Path2D
    private inPath: Path2D

    constructor() {
        this.body = BODIES[1][0]

        this.mouth = getPart("mouth")
        this.leg1 = getPart("leg1")
        this.leg2 = getOtherPart(this.leg1)
        this.arm1 = getPart("arm1")
        this.arm2 = getOtherPart(this.arm1)
        this.eye1 = getPart("eye")
        this.eye2 = this.eye1

        const paths = this.buildPaths()
        this.bodyPath = paths[0]
        this.inPath = paths[1]
        this.outPath = paths[2]
    }

    pinPart(path: Path2D, points: Point[], part: PartFrame) {
            const index = points.findIndex((po) => po.type === part.type)
            if(index < 0) return
            const anchor = points[index]
            points.splice(index, 1)
            const translate = new DOMMatrix()
            translate.translateSelf(anchor?.x, anchor?.y)
            translate.rotateSelf(anchor?.t)
            const subpath = new Path2D(part.path)
            path.addPath(subpath, translate)
    }

    buildPaths(): Path2D[] {
        const points: BodyFrame["points"] = JSON.parse(JSON.stringify(this.body.points))
        const ins = [
            this.mouth,
            this.eye1,
            this.eye2
        ].filter((n) => !!n)

        const outs = [
            this.leg1,
            this.leg2,
            this.arm1,
            this.arm2,
        ].filter((n) => !!n)

        const inParts = new Path2D()
        ins.forEach(this.pinPart.bind(this, inParts, points))
        const outParts = new Path2D()
        outs.forEach(this.pinPart.bind(this, outParts, points))
        const body = new Path2D(this.body.path);
        return [body, inParts, outParts]
    }



    render() {
        if (!Context || !this.bodyPath || !this.outPath || !this.inPath) return
        Context.save()
        Context.stroke(this.outPath)
        Context.globalCompositeOperation = "destination-out"
        Context.fill(this.bodyPath)
        Context.restore()
        Context.stroke(this.bodyPath)
        Context.stroke(this.inPath)
    }
}