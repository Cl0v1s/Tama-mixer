import { LoadBodies, LoadBodyparts, Pet } from "./Pet";


window.onload = async () => {
    await LoadBodies()
    await LoadBodyparts()
    
    const pet = new Pet()

    pet.render()
}
