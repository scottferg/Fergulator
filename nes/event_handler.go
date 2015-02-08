package nes

import (
	"fmt"
	"github.com/robertkrimen/otto"
	"io/ioutil"
	"os"
)

type EventHandler struct {
	handlers map[string][]otto.Value
	vm       *otto.Otto
}

func NewEventHandler(filename string) *EventHandler {
	js, err := ioutil.ReadFile(filename)

	handler := EventHandler{
		handlers: map[string][]otto.Value{},
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, filename, "not readable, not loaded")
		return &handler
	}
	fmt.Println(filename, "loaded")

	vm := otto.New()
	vm.Set("handle", func(call otto.FunctionCall) otto.Value {
		event, err := call.Argument(0).ToString()
		if err != nil {
			// TODO: Handle this
			panic(err)
		}
		handler.handlers[event] = append(handler.handlers[event], call.Argument(1))
		return otto.Value{}
	})
	_, err = vm.Run(js)

	if err != nil {
		fmt.Println(err)
	}
	handler.vm = vm

	return &handler
}

func (handler *EventHandler) HandlePause() {
	state := map[string]interface{}{
		"ram": func(call otto.FunctionCall) otto.Value {
			ram, _ := handler.vm.ToValue(Ram)
			return ram
		},
	}

	ottoState, _ := handler.vm.ToValue(state)
	for _, x := range handler.handlers["pause"] {
		_, err := x.Call(otto.Value{}, ottoState)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}

func (handler *EventHandler) HandleUnpause() {
	state := map[string]interface{}{
		"ram": func(call otto.FunctionCall) otto.Value {
			ram, _ := handler.vm.ToValue(Ram)
			return ram
		},
	}

	ottoState, _ := handler.vm.ToValue(state)
	for _, x := range handler.handlers["unpause"] {
		_, err := x.Call(otto.Value{}, ottoState)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}
