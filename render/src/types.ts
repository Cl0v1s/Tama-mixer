export type Point = {
    x: number, 
    y: number,
    t: number,
    type?: string
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
    name: string
}