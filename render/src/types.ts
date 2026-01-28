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