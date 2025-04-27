package config

import (
	"context"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/andrewbytecoder/k9fyne/utils"
	"go.uber.org/zap"
)

type Ctx struct {
	Ctx    context.Context
	App    fyne.App
	Window fyne.Window
	Config *Cfg
	Log    *zap.Logger
}

func NewCtx() *Ctx {
	log, err := utils.GetZapLog(utils.NewLogConfig(utils.LogLevel("debug"), utils.FileName("k9fyne.log")))
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	return &Ctx{
		Ctx:    ctx,
		App:    app.NewWithID("k9fyne"),
		Config: NewConfig(log, ctx),
		Log:    log,
	}
}

func (ctx *Ctx) GetLogger() *zap.Logger {
	return ctx.Log
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
