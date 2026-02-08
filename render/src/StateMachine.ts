import { IStateMachine, State } from "./types";

export class StateMachine implements IStateMachine {
    currentState: State
    timer?: number
    canChange = true;

    constructor(initState: State) {
        this.currentState = initState
        if(this.currentState.OnEnter) this.currentState.OnEnter(this)
    }

    State() {
        return this.currentState
    }

    private OnEnterManual(nextState: State): boolean {
        this.currentState = nextState
        if(this.currentState.OnEnter) this.currentState.OnEnter(this)
        return true;
    }

    private onEnterDelay(nextState: State): boolean {
        if(nextState.time === undefined || nextState.time === null) {
            console.warn("Delay state must have a time")
            return false
        }
        this.canChange = false;
        clearTimeout(this.timer)
        this.timer = setTimeout(() => { this.canChange = true }, nextState.time)
        this.currentState = nextState
        if(this.currentState.OnEnter) this.currentState.OnEnter(this)
        return true;
    }

    private onEnterTimeout(nextState: State): boolean {
        if(!nextState.OnExit) {
            console.warn("Timeout state must have a callback")
            return false
        }
        if(nextState.time === undefined || nextState.time === null) {
            console.warn("Timeout state must have a time")
            return false
        }
        this.currentState = nextState
        if(this.currentState.OnEnter) this.currentState.OnEnter(this)
        clearTimeout(this.timer)
        this.timer = setTimeout(() => {
            if(this.currentState.OnExit) this.currentState.OnExit(this)
        }, nextState.time)
        return true;
    }


    Enter(nextState: State): boolean {
        if (this.currentState.next.indexOf(nextState.key) === -1 || !this.canChange) {
            return false;
        }
        switch (nextState.condition) {
            case "timeout": return this.onEnterTimeout(nextState)
            case "delay": return this.onEnterDelay(nextState)
            case "manual":
            default:
                return this.OnEnterManual(nextState)
        }
    }
}