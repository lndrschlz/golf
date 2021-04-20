package main

import (
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
	"reflect"
)

type golf struct {
	interpreter *interp.Interpreter
}

func NewGolf(gopath string) *golf {
	interpreter := interp.New(interp.Options{
		GoPath: gopath,
	})
	interpreter.Use(stdlib.Symbols)
	// default imports
	doDotImport(interpreter, "fmt")
	doDotImport(interpreter, "strings")
	doImport(interpreter, "os")
	builtIns(interpreter)
	return &golf{interpreter: interpreter}
}

func builtIns(interpreter *interp.Interpreter) {
	statements := []string{
		"var arg = make(map[string]interface{})",
		"func isSet(key string) bool { _, ok := arg[key]; return ok }",
	}
	for _, statement := range statements {
		_, err := interpreter.Eval(statement)
		checkFail(err, "Failed to initialize the built-in variables and function. Detail: " + statement)
	}
}

func doImport(interpreter *interp.Interpreter, pkg string) {
	_, err := interpreter.Eval(`import "` + pkg + `"`)
	checkFail(err, "Package not found: "+pkg)
}

func doDotImport(interpreter *interp.Interpreter, pkg string) {
	_, err := interpreter.Eval(`import . "` + pkg + `"`)
	checkFail(err, "Package not found: "+pkg)
}

func (g *golf) eval(call string, line string) (result string, ok bool, repeat bool, err error) {
	// disable the repeat instruction
	_, err = g.interpreter.Eval("repeat := false")
	if err != nil {
		return "", false, false, err
	}
	_, err = g.interpreter.Eval("line := `" + line + "`")
	if err != nil {
		return "", false, false, err
	}
	_, err = g.interpreter.Eval("token := Split(line, ` `)")
	if err != nil {
		return "", false, false, err
	}
	resultValue, err := g.interpreter.Eval(call)
	result = resultValue.String()
	ok = resultValue.Kind() == reflect.String
	if err != nil {
		return "", false, false, err
	}
	repeatValue, err := g.interpreter.Eval("repeat")
	repeat = repeatValue.Bool()
	return
}
