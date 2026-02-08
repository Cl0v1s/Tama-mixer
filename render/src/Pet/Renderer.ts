import { COLOR_PALETTE, getBody, getClosedEyeFrame, getOtherPart, getPart, getYOffset } from '../utils';
import { Context } from './../Canvas'
import { BodyFrame, PartFrame, Point, Rect, Renderer, RendererListener } from './../types'
import { AnimationConfig } from './Animations';

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


export class PetRenderer implements Renderer {

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

    public colors = [COLOR_PALETTE.pet.stroke, COLOR_PALETTE.pet.fill1, COLOR_PALETTE.pet.fill2]

    private boundingBox!: Rect;
    /**
     * negative = blinking cooldown / positive = blinking 
     */
    private blink = 0

    /**
     * Currently playing animation
     */
    private animation?: AnimationConfig;

    public revert: boolean = false;

    private listeners: RendererListener[] = [];


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

        PetRenderer.CLOSED_EYE_FRAME = getClosedEyeFrame()
    }

    private pinPart(path: Path2D, points: Point[], part: PartFrame): Rect {
        try {
            const index = points.findIndex((po) => po.type === part.type);
            if (index < 0) {
                return { x: 0, y: 0, width: 0, height: 0 };
            }

            const anchor = points[index];
            points.splice(index, 1);

            // Matrice de transformation : translation + rotation autour de l'ancre
            const matrix = new DOMMatrix();
            matrix.translateSelf(anchor.x, anchor.y);
            matrix.rotateSelf(anchor.t);           // rotation en degrés

            // Ajout du sous-chemin avec la transformation
            const subpath = new Path2D(part.path);
            path.addPath(subpath, matrix);

            const radians = anchor.t * (Math.PI / 180)
            const bounds = [part.boundingBox.topLeft, { x: part.boundingBox.bottomRight.x, y: part.boundingBox.topLeft.y }, part.boundingBox.bottomRight, { x: part.boundingBox.topLeft.x, y: part.boundingBox.bottomRight.y }]
                .map((point) => ({ ...point, x: point.x * Math.cos(radians) - point.y * Math.sin(radians), y: point.x * Math.sin(radians) + point.y * Math.cos(radians) } as Point))
                .map((point) => ({ ...point, x: point.x + anchor.x, y: point.y + anchor.y } as Point)) // translate 
                .reduce(
                    (acc, p) => ({
                        minX: Math.min(acc.minX, p.x),
                        minY: Math.min(acc.minY, p.y),
                        maxX: Math.max(acc.maxX, p.x),
                        maxY: Math.max(acc.maxY, p.y),
                    }),
                    { minX: Infinity, minY: Infinity, maxX: -Infinity, maxY: -Infinity }
                );

            return {
                x: bounds.minX,
                y: bounds.minY,
                width: bounds.maxX - bounds.minX,
                height: bounds.maxY - bounds.minY,
            };
        } catch (e) {
            console.error('Error building body with ' + part.name + part.frame)
            throw e
        }
    }

    private buildPaths() {
        let boxes: Rect[] = []
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
        boxes = [...boxes, ...ins.map(this.pinPart.bind(this, inParts, points))]
        const outParts = new Path2D()
        boxes = [...boxes, ...outs.map(this.pinPart.bind(this, outParts, points))]
        const body = new Path2D(this.body.path);
        this.bodyPath = body
        this.inPath = inParts
        this.outPath = outParts


        let minX = 0;
        let minY = 0;
        let maxX = this.body.size.x;
        let maxY = this.body.size.y;

        boxes.forEach(b => {
            minX = Math.min(minX, b.x);
            minY = Math.min(minY, b.y);
            maxX = Math.max(maxX, b.x + b.width);
            maxY = Math.max(maxY, b.y + b.height);
        });

        this.boundingBox = {
            x: minX,
            y: minY,
            width: (maxX - minX),
            height: (maxY - minY)
        };
    }

    public UnSubscribe(listener: RendererListener): void {
        this.listeners = this.listeners.filter((b) => b !== listener)
    }

    public Subscribe(listener: RendererListener): void {
        this.listeners.push(listener);
    }

    private onAnimationEnd() {
        this.listeners.forEach((l) => l.OnRenderer(this))
        this.animation = undefined;
    }


    private frameCounter = 0;
    private animate() {
        this.frameCounter = (this.frameCounter + 1);

        const anim = this.animation
        if (!anim) return

        let needRebuild = false;

        // animate body
        if (anim.body && this.frameCounter % anim.body === 0) {
            this.body = this.bodys[(this.body.frame + 1) % this.bodys.length]
            needRebuild = true
        }

        if (anim && anim.parts.length > 0 && this.frameCounter % anim.speed === 0) {
            for (let partName of anim.parts) {
                const frames = this[`${partName}s`] as PartFrame[];
                const current = this[partName as keyof PetRenderer] as PartFrame;

                let nextIdx = current.frame + 1;
                if (nextIdx >= frames.length) {
                    nextIdx = anim.loop ? 0 : frames.length - 1;
                    if (!anim.loop) {
                        this.onAnimationEnd()
                    }
                }

                let nextFrame = frames[nextIdx]

                if (this[partName as keyof PetRenderer] !== nextFrame) {
                    (this[partName as keyof PetRenderer] as any) = nextFrame;
                    needRebuild = true;
                }
            }
        }

        // blinking management 
        if (this.blink > 0) {
            this.blink--
            if (this.blink === 1) {
                needRebuild = true
                this.blink = BLINK_COOLDOWN * -1
            }
        }
        if (this.blink < 0) this.blink++
        if (this.blink == 0 && Math.floor(Math.random() * BLINK_PROBABILITY) === 0) {
            this.blink = BLINK_DURATION
            needRebuild = true
        }
        if (needRebuild) {
            if (this.blink > 0) {
                this.eye1 = PetRenderer.CLOSED_EYE_FRAME
                this.eye2 = PetRenderer.CLOSED_EYE_FRAME
            } else if (this.blink < 0) {
                this.eye1 = this.eye1s[0]
                this.eye2 = this.eye2s[0]
            }
        }


        if (needRebuild) {
            this.buildPaths();
        }
    }

    Play(animation: AnimationConfig) {
        this.animation = animation
    }

    Render(x: number, y: number, z: number) {
        if (!Context || !this.bodyPath || !this.outPath || !this.inPath) return

        y += getYOffset(x + (this.boundingBox.width * z) / 2 )

        let transform = new DOMMatrix()
        if (this.revert) {
            transform.translateSelf(x + this.boundingBox.width * z, y)
            transform = transform.flipX()
        } else {
            transform.translateSelf(x, y)
        }
        transform.scaleSelf(z, z)


        this.animate()

        const ctx = Context;
        const t = transform;        
        const [border, bodyFill, partsFill] = this.colors;  // ex: border=#000, body=#fff, parts=#f00

        ctx.save();
        ctx.setTransform(t);          

        ctx.lineWidth = 1;
        ctx.strokeStyle = border;

        ctx.fillStyle = partsFill;
        ctx.fill(this.outPath);
        ctx.stroke(this.outPath);

        ctx.globalCompositeOperation = "destination-out";
        ctx.fill(this.bodyPath);     
        ctx.globalCompositeOperation = "source-over"; 

        ctx.fillStyle = bodyFill;
        ctx.fill(this.bodyPath);
        ctx.stroke(this.bodyPath);

        ctx.fillStyle = partsFill;
        ctx.fill(this.inPath);
        ctx.stroke(this.inPath);

        ctx.restore();
    }

    BoundingBox(): Rect {
        return this.boundingBox
    }
}