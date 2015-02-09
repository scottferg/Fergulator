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

func (handler *EventHandler) ReloadFile(filename string) {
	fmt.Println("Reloading", filename)

	js, err := ioutil.ReadFile(filename)

	if err != nil {
		fmt.Fprintln(os.Stderr, filename, "not readable, not loaded")
		return
	}

	// Clear out handlers so we don't end up with double callbacks.
	// Keep variables though, they are useful.
	handler.handlers = map[string][]otto.Value{}

	_, err = handler.vm.Run(js)

	if err != nil {
		fmt.Println(err)
	}
	handler.Handle("reload")
}

func NewEventHandler(filename string) *EventHandler {
	handler := EventHandler{
		handlers: map[string][]otto.Value{},
		vm:       otto.New(),
	}

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
	handler.vm = vm
	handler.ReloadFile(filename)
	handler.Handle("init")

	return &handler
}

func (handler *EventHandler) Handle(event string) {
	state := map[string]interface{}{
		"ram": func(call otto.FunctionCall) otto.Value {
			ram, _ := handler.vm.ToValue(Ram[0:0x800])
			return ram
		},
		"writeRam": func(call otto.FunctionCall) otto.Value {
			idx, _ := call.Argument(0).ToInteger()
			val, _ := call.Argument(1).ToInteger()

			err := Ram.Write(Word(idx), Word(val))
			if err != nil {
				fmt.Println(err)
			}

			return otto.Value{}
		},
	}

	ottoState, _ := handler.vm.ToValue(state)
	for _, x := range handler.handlers[event] {
		_, err := x.Call(otto.Value{}, ottoState)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}
