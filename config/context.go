package config

import (
	"context"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/andrewbytecoder/k9fyne/kube"
	"github.com/andrewbytecoder/k9fyne/utils"
	"go.uber.org/zap"
)

type Ctx struct {
	Ctx    context.Context
	App    fyne.App
	Window fyne.Window
	Config *Cfg
	Log    *zap.Logger
	k9Info *kube.K9Info
}

func NewCtx() *Ctx {
	log, err := utils.GetZapLog(utils.NewLogConfig(utils.LogLevel("debug"), utils.FileName("k9fyne.log")))
	if err != nil {
		panic(err)
	}
	return &Ctx{
		Ctx:    context.Background(),
		App:    app.NewWithID("k9fyne"),
		Config: NewConfig(log),
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

func (ctx *Ctx) GetK9Info() *kube.K9Info {
	return ctx.k9Info
}

func (ctx *Ctx) SetK9Info(k9Info *kube.K9Info) *Ctx {
	ctx.k9Info = k9Info
	return ctx
}
