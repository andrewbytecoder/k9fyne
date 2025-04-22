package config

import (
	"context"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

type Ctx struct {
	Ctx    context.Context
	App    fyne.App
	Window fyne.Window
	Config *Cfg
}

func NewCtx() *Ctx {
	return &Ctx{
		Ctx:    context.Background(),
		App:    app.NewWithID("k9fyne"),
		Config: NewConfig(),
	}
}

func (ctx *Ctx) SetWindow(window fyne.Window) *Ctx {
	ctx.Window = window
	return ctx
}
func (ctx *Ctx) GetWindow() fyne.Window {
	return ctx.Window
}
func (ctx *Ctx) SetApp(app fyne.App) *Ctx {
	ctx.App = app
	return ctx
}
func (ctx *Ctx) SetConfig(config *Cfg) *Ctx {
	ctx.Config = config
	return ctx
}
func (ctx *Ctx) GetApp() fyne.App {
	return ctx.App
}
func (ctx *Ctx) GetConfig() *Cfg {
	return ctx.Config
}
