export interface RendererListener {
    OnRenderer(renderer: Renderer): void;
}

export interface Renderer { 
    Subscribe(listener: RendererListener): void
    UnSubscribe(listener: RendererListener): void
    BoundingBox(): Rect
    Render(x: number, y: number, z: number): void
}

export interface Physics {
    Vector(): Point
    Entity(): Entity
    Tick(alpha: number): void
}

export interface Entity extends RendererListener { 
    X(): number 
    Y(): number 
    W(): number
    H(): number;
    Z(): number
    Render(): void
    Move(x: number, y: number, z?: number): Entity 
    Destroy(): void 
}

export type Point = {
    x: number, 
    y: number,
    t: number,
    type?: string
}

export type Rect = {
    x: number,
    y: number,
    width: number,
    height: number,
}

export type BodyFrame = {
    path: string,
    points: Point[],
    frame: number,
    name: string,
    size: Point
}

export type PartFrame = {
    path: string,
    type: string,
    frame: number,
    name: string,
    boundingBox: {
        topLeft: Point,
        bottomRight: Point
    }
}

export interface IStateMachine {
    Enter(nextState: State): boolean
}

export type StateLeaveCondition = "manual" | "timeout" | "delay"

export type State = {
    condition: StateLeaveCondition,
    time?: number,
    callback?: (caller: IStateMachine) => void,
    next: State[]
}