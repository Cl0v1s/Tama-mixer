import { Canvas, Context } from "./Canvas";
import { LoadBodies, LoadBodyparts, Pet } from "./Pet";


let pet: Pet;


function render() {
    if(!Canvas || !Context) return
    Context.fillStyle = "white"
    Context.fillRect(0, 0, Canvas.clientWidth, Canvas.clientHeight)
    pet.render()
    requestAnimationFrame(render)
}

window.onload = async () => {
    await LoadBodies()
    await LoadBodyparts()
    pet = new Pet()
    requestAnimationFrame(render)
}
