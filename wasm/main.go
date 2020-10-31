package main

import (
	"fmt"
	"syscall/js"

	"github.com/domino14/macondo/analyzer"
	"github.com/domino14/macondo/cache"
)

func precache(this js.Value, args []js.Value) interface{} {
	cache.Precache(args[0].String(), readJSBytes(args[1]))
	return nil
}

func analyze(this js.Value, args []js.Value) interface{} {
	// JS doesn't use utf8, but it converts automatically if we take/return strings.
	jsonBoard := []byte(args[0].String())

	an := analyzer.NewDefaultAnalyzer()
	jsonMoves, err := an.Analyze(jsonBoard)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	jsonMovesStr := string(jsonMoves)
	return jsonMovesStr
}

func registerCallbacks() {
	js.Global().Get("resMacondo").Invoke(map[string]interface{}{
		"precache": js.FuncOf(precache),
		"analyze":  js.FuncOf(analyze),
	})
}

func main() {
	registerCallbacks()
	// Keep Go "program" running.
	select {}
}
