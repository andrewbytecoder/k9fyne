// Package main provides various examples of Fyne API capabilities.
package main

import (
	"fmt"
	"github.com/andrewbytecoder/k9fyne/config"
	"github.com/andrewbytecoder/k9fyne/kube/topo"
	"github.com/andrewbytecoder/k9fyne/resources"
	"github.com/andrewbytecoder/k9fyne/widgets"
	"go.uber.org/zap"
	"log"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/cmd/fyne_settings/settings"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const preferenceCurrentTutorial = "currentTutorial"

var topWindow fyne.Window

func main() {
	ctx := config.NewCtx()
	ctx.App.SetIcon(resources.K9FyneFireLogo)
	makeTray(ctx.App)
	logLifecycle(ctx.App, ctx.GetLogger())
	topWindow = ctx.App.NewWindow("k9fyne")
	ctx.SetWindow(topWindow)

	ctx.GetWindow().SetMainMenu(makeMenu(ctx.App, topWindow, ctx.GetConfig()))
	ctx.GetWindow().SetMaster()

	content := container.NewStack()
	title := widget.NewLabel("Component name")
	intro := widget.NewLabel("An introduction would probably go\nhere, as well as a")
	intro.Wrapping = fyne.TextWrapWord

	top := container.NewVBox(title, widget.NewSeparator(), intro)
	setContent := func(t widgets.Tutorial) {
		title.SetText(t.Title)
		isMarkdown := len(t.Intro) == 0
		if !isMarkdown {
			intro.SetText(t.Intro)
		}

		if topo.K9InfoHandler != nil {
			err := topo.K9InfoHandler.FetchData(&t)
			if err != nil {
				ctx.Log.Error("Failed to fetch data", zap.Error(err))
				return
			}
		}

		if t.Title == "Kubernetes" || isMarkdown {
			top.Hide()
		} else {
			top.Show()
		}

		content.Objects = []fyne.CanvasObject{t.View(ctx.GetWindow(), t.Data)}
		content.Refresh()
	}

	tutorial := container.NewBorder(
		top, nil, nil, nil, content)
	split := container.NewHSplit(makeNav(setContent, true), tutorial)
	split.Offset = 0.2
	ctx.GetWindow().SetContent(split)
	ctx.GetWindow().Resize(fyne.NewSize(640, 460))
	// connect window
	ctx.GetConfig().GetSSH().CreateSSHClient(topWindow)
	ctx.GetWindow().ShowAndRun()
}

func logLifecycle(a fyne.App, l *zap.Logger) {
	a.Lifecycle().SetOnStarted(func() {
		log.Println("Lifecycle: Started")
	})
	a.Lifecycle().SetOnStopped(func() {
		log.Println("Lifecycle: Stopped")
	})
	a.Lifecycle().SetOnEnteredForeground(func() {
		log.Println("Lifecycle: Entered Foreground")
	})
	a.Lifecycle().SetOnExitedForeground(func() {
		log.Println("Lifecycle: Exited Foreground")
	})
}

func makeMenu(a fyne.App, w fyne.Window, c *config.Cfg) *fyne.MainMenu {
	newItem := fyne.NewMenuItem("New", nil)
	sshSettings := func() {
		c.GetSSH().CreateSSHClient(w)
	}

	sshItem := fyne.NewMenuItem("SSH Config", sshSettings)
	sshItem.Icon = theme.ComputerIcon()
	newItem.ChildMenu = fyne.NewMenu("",
		sshItem,
	)

	openSettings := func() {
		w := a.NewWindow("K9Fyne Settings")
		w.SetContent(settings.NewSettings().LoadAppearanceScreen(w))
		w.Resize(fyne.NewSize(440, 520))
		w.Show()
	}
	settingsItem := fyne.NewMenuItem("Settings", openSettings)
	settingsShortcut := &desktop.CustomShortcut{KeyName: fyne.KeyComma, Modifier: fyne.KeyModifierShortcutDefault}
	settingsItem.Shortcut = settingsShortcut
	w.Canvas().AddShortcut(settingsShortcut, func(shortcut fyne.Shortcut) {
		openSettings()
	})

	cutShortcut := &fyne.ShortcutCut{Clipboard: a.Clipboard()}
	cutItem := fyne.NewMenuItem("Cut", func() {
		shortcutFocused(cutShortcut, a.Clipboard(), w.Canvas().Focused())
	})
	cutItem.Shortcut = cutShortcut
	copyShortcut := &fyne.ShortcutCopy{Clipboard: a.Clipboard()}
	copyItem := fyne.NewMenuItem("Copy", func() {
		shortcutFocused(copyShortcut, a.Clipboard(), w.Canvas().Focused())
	})
	copyItem.Shortcut = copyShortcut
	pasteShortcut := &fyne.ShortcutPaste{Clipboard: a.Clipboard()}
	pasteItem := fyne.NewMenuItem("Paste", func() {
		shortcutFocused(pasteShortcut, a.Clipboard(), w.Canvas().Focused())
	})
	pasteItem.Shortcut = pasteShortcut
	performFind := func() { fmt.Println("Menu Find") }
	findItem := fyne.NewMenuItem("Find", performFind)
	findItem.Shortcut = &desktop.CustomShortcut{KeyName: fyne.KeyF, Modifier: fyne.KeyModifierShortcutDefault | fyne.KeyModifierAlt | fyne.KeyModifierShift | fyne.KeyModifierControl | fyne.KeyModifierSuper}
	w.Canvas().AddShortcut(findItem.Shortcut, func(shortcut fyne.Shortcut) {
		performFind()
	})

	helpMenu := fyne.NewMenu("Help",
		fyne.NewMenuItem("Documentation", func() {
			u, _ := url.Parse("https://developer.fyne.io")
			_ = a.OpenURL(u)
		}),
		fyne.NewMenuItem("Support", func() {
			u, _ := url.Parse("https://fyne.io/support/")
			_ = a.OpenURL(u)
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Sponsor", func() {
			u, _ := url.Parse("https://fyne.io/sponsor/")
			_ = a.OpenURL(u)
		}))

	// a quit item will be appended to our first (File) menu
	file := fyne.NewMenu("File", newItem)
	file.Items = append(file.Items, fyne.NewMenuItemSeparator(), settingsItem)
	mainMenu := fyne.NewMainMenu(
		file,
		fyne.NewMenu("Edit", cutItem, copyItem, pasteItem, fyne.NewMenuItemSeparator(), findItem),
		helpMenu,
	)

	return mainMenu
}

func makeTray(a fyne.App) {
	if desk, ok := a.(desktop.App); ok {
		h := fyne.NewMenuItem("Hello", func() {})
		h.Icon = theme.HomeIcon()
		menu := fyne.NewMenu("Hello World", h)
		h.Action = func() {
			log.Println("System tray menu tapped")
			h.Label = "Welcome"
			menu.Refresh()
		}
		desk.SetSystemTrayMenu(menu)
	}
}

func makeNav(setTutorial func(tutorial widgets.Tutorial), loadPrevious bool) fyne.CanvasObject {
	a := fyne.CurrentApp()

	tree := &widget.Tree{
		ChildUIDs: func(uid string) []string {
			return widgets.TutorialIndex[uid]
		},
		IsBranch: func(uid string) bool {
			children, ok := widgets.TutorialIndex[uid]

			return ok && len(children) > 0
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("Collection Widgets")
		},
		UpdateNode: func(uid string, branch bool, obj fyne.CanvasObject) {
			t, ok := widgets.Tutorials[uid]
			if !ok {
				fyne.LogError("Missing tutorial panel: "+uid, nil)
				return
			}
			obj.(*widget.Label).SetText(t.Title)
		},
		OnSelected: func(uid string) {
			if t, ok := widgets.Tutorials[uid]; ok {
				for _, f := range widgets.OnChangeFuncs {
					f()
				}
				widgets.OnChangeFuncs = nil // Loading a page registers a new cleanup.

				a.Preferences().SetString(preferenceCurrentTutorial, uid)
				setTutorial(t)
			}
		},
	}

	if loadPrevious {
		currentPref := a.Preferences().StringWithFallback(preferenceCurrentTutorial, "welcome")
		tree.Select(currentPref)
	}

	themes := container.NewGridWithColumns(2,
		widget.NewButton("Dark", func() {
			a.Settings().SetTheme(&forcedVariant{Theme: theme.DefaultTheme(), variant: theme.VariantDark})
		}),
		widget.NewButton("Light", func() {
			a.Settings().SetTheme(&forcedVariant{Theme: theme.DefaultTheme(), variant: theme.VariantLight})
		}),
	)

	return container.NewBorder(nil, themes, nil, nil, tree)
}

func shortcutFocused(s fyne.Shortcut, cb fyne.Clipboard, f fyne.Focusable) {
	switch sh := s.(type) {
	case *fyne.ShortcutCopy:
		sh.Clipboard = cb
	case *fyne.ShortcutCut:
		sh.Clipboard = cb
	case *fyne.ShortcutPaste:
		sh.Clipboard = cb
	}
	if focused, ok := f.(fyne.Shortcutable); ok {
		focused.TypedShortcut(s)
	}
}
