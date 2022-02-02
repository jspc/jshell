package main

import (
	"github.com/manifoldco/promptui"

	"github.com/jspc/jshell/apps/hello-world"
	"github.com/jspc/jshell/apps/hex"
	"github.com/jspc/jshell/apps/quit"
	"github.com/jspc/jshell/apps/wordle"
)

var (
	Apps = []App{
		new(helloworld.HelloWorld),
		new(wordle.Wordle),
		new(hex.Hex),

		// Place at end, seems less jumbled
		new(quit.Quit),
	}

	appMenuTemplate = &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "\U0001F336  {{ .Name | cyan }}",
		Inactive: "   {{ .Name | cyan }}",
		Selected: "\U0001F336  {{ .Name | red | cyan }}",
		Details: `

{{ "Name:" | faint }}	{{ .Name }}
{{ "Description:" | faint }}	{{ .Desc }}`,
	}
)

// App is an interface all jshell apps must implement
// in order to be run.
type App interface {
	// Name should return the application name in a way
	// that our menu system can read
	Name() string

	// Description should return a string describing what
	// this application does
	Description() string

	// Run will *run* the application, in the sense that
	// this function should be the entrypoint to whatever the
	// app does
	Run() error

	// Cleanup performs some kind of action once an application
	// finishes.
	//
	// It's up to an application what that means (if anything).
	//
	// jshell will *always* clear screen after an application finishes
	Cleanup() error
}

// AppMenu turns our list of Apps into a menu for PromptUI to
// display, returning the relevant App to run
func AppMenu() (a App, err error) {
	type menuItem struct{ Name, Desc string }

	menuData := make([]menuItem, len(Apps))

	for i, app := range Apps {
		menuData[i] = menuItem{app.Name(), app.Description()}
	}

	prompt := promptui.Select{
		Label:     "Application Menu",
		Items:     menuData,
		Templates: appMenuTemplate,
		Size:      4,
	}

	idx, _, err := prompt.Run()

	a = Apps[idx]

	return
}
