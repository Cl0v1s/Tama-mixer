import { Canvas, Context } from "./Canvas";
import { PetEntity } from "./Pet/Entity";
import { Entity } from "./types";
import { LoadBodies, LoadBodyparts } from "./utils";


const entities: Entity[] = []


function render() {
    if(!Canvas || !Context) return
    Context.fillStyle = "white"
    Context.fillRect(0, 0, Canvas.clientWidth, Canvas.clientHeight)

    entities.forEach((e) => e.render())

    requestAnimationFrame(render)
}

window.onload = async () => {
    await LoadBodies()
    await LoadBodyparts()
    if(!Canvas || !Context) return
    entities.push(
        new PetEntity(Canvas.clientWidth/2, Canvas.clientHeight/2)
    ) 
    requestAnimationFrame(render)
}
