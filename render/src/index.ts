import { Canvas, Context } from "./Canvas";
import { PetEntity } from "./Pet/Entity";
import { Entity } from "./types";
import { LoadBodies, LoadBodyparts } from "./utils";

const entities: Entity[] = []

function render() {
    if(!Canvas || !Context) return
    Context.fillStyle = "white"
    Context.fillRect(0, 0, Canvas.clientWidth, Canvas.clientHeight)

    entities.forEach((e) => e.Render())

    requestAnimationFrame(render)
}

window.onload = async () => {
    await LoadBodies()
    await LoadBodyparts()
    if(!Canvas || !Context) return
    const pet = new PetEntity()
    pet.Move(0, 0, 2)
    pet.Move(Canvas.clientWidth / 2 - pet.W() / 2, Canvas.clientHeight / 2 - pet.H() / 2)
    entities.push(
        pet
    ) 
    requestAnimationFrame(render)
}
