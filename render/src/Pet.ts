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
    BODYPARTS = (await Promise.all(promises)).map((a: PartFrame[]) => a.sort((a, b) => a.frame - b.frame))
}

function getBody() {
    return BODIES[Math.floor(Math.random() * BODIES.length)]
}


function getPart(type: string, dont: string[] = []): PartFrame[] {
    const parts = BODYPARTS.filter((o) => o[0].type === type && dont.indexOf(o[0].name) === -1)
    return parts[Math.floor(Math.random() * parts.length)]
}

function getOtherPart(part: PartFrame[]): PartFrame[] {
    const parts = BODYPARTS.filter((o) => o[0].type !== part[0].type)
    return parts.find((p) => p[0].name === part[0].name)!
}


enum PetAnimationState {
    IDLE,
    EATING,
    WALKING
}

type BodyPart = "mouth" | "leg1" | "leg2" | "arm1" | "arm2" | "eye1" | "eye2"

type AnimationConfig = {
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
     * Does animation allow animated body
     */
    body: boolean
};

const ANIMATIONS: Record<PetAnimationState, AnimationConfig> = {
    [PetAnimationState.IDLE]: {
        parts: [],
        speed: 1,
        loop: true,
        body: true
    },
    [PetAnimationState.WALKING]: {
        parts: ['leg1', 'leg2', 'arm1', 'arm2'],
        speed: 20,
        loop: true,
        body: true
    },
    [PetAnimationState.EATING]: {
        parts: ['mouth'],
        speed: 20,
        loop: true,
        body: false
    }
};


/**
 * Number of tick blinking
 */
const BLINK_DURATION = 10

/**
 * Number of ticks before another blinking can occure 
 */
const BLINK_COOLDOWN = 300

/**
 * 1/BLINK_PROBABILITY chance to blink per tick 
 */
const BLINK_PROBABILITY = 100


export class Pet {

    private bodys: BodyFrame[]
    private body: BodyFrame;
    private mouths: PartFrame[]
    public mouth: PartFrame;
    private leg1s: PartFrame[]
    public leg1: PartFrame;
    private leg2s: PartFrame[]
    public leg2: PartFrame;
    private arm1s: PartFrame[]
    public arm1: PartFrame;
    private arm2s: PartFrame[]
    public arm2: PartFrame;
    private eye1s: PartFrame[]
    private eye2s: PartFrame[]
    public eye1: PartFrame;
    public eye2: PartFrame;

    public colors = ["#004990", "#f7c8dd", "#a1d6dd"]

    /**
     * Speed of body frame change
     */
    public bodySpeed = 100
    /**
     * negative = blinking cooldown / positive = blinking 
     */
    public blink = 0
    /** 
     * Level of zoom
     */
    public zoom = 2;
    public x = 10;
    public y = 10;

    /**
     * Currently playing animation
     */
    public animationState: PetAnimationState = PetAnimationState.WALKING

    private bodyPath!: Path2D
    private outPath!: Path2D
    private inPath!: Path2D

    private static CLOSED_EYE_FRAME: PartFrame;

    constructor() {
        this.bodys = getBody()
        this.body = this.bodys[0]

        this.mouths = getPart("mouth")
        this.mouth = this.mouths[0]
        this.leg1s = getPart("leg1")
        this.leg1 = this.leg1s[0]
        this.leg2s = getOtherPart(this.leg1s)
        this.leg2 = this.leg2s[0]
        this.arm1s = getPart("arm1")
        this.arm1 = this.arm1s[0]
        this.arm2s = getOtherPart(this.arm1s)
        this.arm2 = this.arm2s[0]
        this.eye1s = getPart("eye", ["CLOSED"])
        this.eye2s = this.eye1s
        this.eye1 = this.eye1s[0]
        this.eye2 = this.eye2s[0]

        console.log(this)

        this.buildPaths()

        Pet.CLOSED_EYE_FRAME = BODYPARTS.find((d) => d[0].name === "CLOSED")![0]
    }

    private pinPart(path: Path2D, points: Point[], part: PartFrame) {
        const index = points.findIndex((po) => po.type === part.type)
        if (index < 0) return
        const anchor = points[index]
        points.splice(index, 1)
        const translate = new DOMMatrix()
        translate.translateSelf(anchor?.x, anchor?.y)
        translate.rotateSelf(anchor?.t)
        const subpath = new Path2D(part.path)
        path.addPath(subpath, translate)
    }

    private buildPaths() {
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

        this.bodyPath = body
        this.inPath = inParts
        this.outPath = outParts
    }

    private frameCounter = 0;

    private animate() {
        this.frameCounter = (this.frameCounter + 1);

        const anim = ANIMATIONS[this.animationState];
        if (!anim) return

        let needRebuild = false;

        // animate body
        if (anim.body && this.frameCounter % this.bodySpeed === 0) {
            this.body = this.bodys[(this.body.frame + 1) % this.bodys.length]
            needRebuild = true
        }

        if (anim && anim.parts.length > 0 && this.frameCounter % anim.speed === 0) {
            for (let partName of anim.parts) {
                const frames = this[`${partName}s`] as PartFrame[];
                const current = this[partName as keyof Pet] as PartFrame;

                let nextIdx = current.frame + 1;
                if (nextIdx >= frames.length) {
                    nextIdx = anim.loop ? 0 : frames.length - 1;
                    if (!anim.loop) {
                        this.animationState = PetAnimationState.IDLE;
                    }
                }

                let nextFrame = frames[nextIdx]

                if (this[partName as keyof Pet] !== nextFrame) {
                    (this[partName as keyof Pet] as any) = nextFrame;
                    needRebuild = true;
                }
            }
        }

        // blinking management 
        if(this.blink > 0) {
            this.blink--
            if(this.blink === 1) {
                needRebuild = true
                this.blink = BLINK_COOLDOWN * -1
            }
        }
        if(this.blink < 0) this.blink++
        if (this.blink == 0 && Math.floor(Math.random() * BLINK_PROBABILITY) === 0) {
            this.blink = BLINK_DURATION
            needRebuild = true
        } 
        if(needRebuild) {
            if(this.blink > 0) {
                this.eye1 = Pet.CLOSED_EYE_FRAME
                this.eye2 = Pet.CLOSED_EYE_FRAME
            } else if(this.blink < 0) {
                this.eye1 = this.eye1s[0]
                this.eye2 = this.eye2s[0]
            }
        }


        if (needRebuild) {
            this.buildPaths();
        }
    }

    render() {
        if (!Context || !this.bodyPath || !this.outPath || !this.inPath) return
        const transform = new DOMMatrix()
        transform.scaleSelf(this.zoom, this.zoom)
        transform.translateSelf(this.x, this.y)

        this.animate()



        Context.save()
        Context.save()
        Context.setTransform(transform)
        Context.fillStyle = "red"
        Context.fillRect(0, 0, this.body.size.x, this.body.size.y);
        Context.lineWidth = 1
        Context.fillStyle = this.colors[2]
        Context.fill(this.outPath)
        Context.strokeStyle = this.colors[0]
        Context.stroke(this.outPath)
        Context.globalCompositeOperation = "destination-out"
        Context.fill(this.bodyPath)
        Context.restore()
        Context.setTransform(transform)
        Context.fillStyle = this.colors[1]
        Context.fill(this.bodyPath)
        Context.lineWidth = 1
        Context.strokeStyle = this.colors[0]
        Context.stroke(this.bodyPath)
        Context.stroke(this.inPath)
        Context.restore()
    }
}