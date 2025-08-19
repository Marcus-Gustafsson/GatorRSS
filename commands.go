package main

import "fmt"

// run executes the given command using the associated handler, if it exists.
func (cmdsPtr *cmds) run(stPtr *state, cmd command) error{

    function, ok := cmdsPtr.FunctionMap[cmd.Name]
    // check if a handler exists before calling it, prevents a panic from a nil map entry.
    if !ok {
        return fmt.Errorf("unknown command: %s", cmd.Name)
    }
    return function(stPtr, cmd)
}

// register adds a command handler function for a given command name to used cmds structs map.
func (cmdsPtr *cmds) register(name string, function func(*state, command) error){
    cmdsPtr.FunctionMap[name] = function
}