//go:generate go run gen.go

package widgets

import (
	"image/color"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/styles"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// OnChangeFuncs is a slice of functions that can be registered
// to run when the user switches tutorial.
var OnChangeFuncs []func()

// Tutorial defines the data structure for a tutorial
type Tutorial struct {
	Title, Intro string
	View         func(w fyne.Window, data interface{}) fyne.CanvasObject
	Data         interface{}
}

type GeneratedTutorial struct {
	title string

	content []string
	code    []func() fyne.CanvasObject
}

var (
	// Tutorials defines the metadata for each tutorial
	Tutorials = map[string]Tutorial{
		"k8s": {"Kubernetes", "", welcomeScreen, nil},
		"pods": {"Pod info",
			"Show pods info.",
			makePodList, nil,
		},
		"Doc": {"Doc",
			"使用指导文档和常见命令.",
			k9Doc,
			nil,
		},
		"canvastext": loadDefinition("Text", "canvas/text.md"),
		"line":       loadDefinition("Line", "canvas/line.md"),
		"rectangle":  loadDefinition("Rectangle", "canvas/rectangle.md"),
		"circle":     loadDefinition("Circle", "canvas/circle.md"),
		"image":      loadDefinition("Image", "canvas/image.md"),
		"raster":     loadDefinition("Raster", "canvas/raster.md"),
		"gradient":   loadDefinition("Gradient", "canvas/gradient.md"),
		"animations": {"Animations",
			"See how to animate components.",
			makeAnimationScreen,
			nil,
		},
		"icons": {"Theme Icons",
			"Browse the embedded icons.",
			iconScreen,
			nil,
		},
		"containers": {"Containers",
			"Containers group other widgets and canvas objects, organising according to their layout.\n" +
				"Standard containers are illustrated in this section, but developers can also provide custom " +
				"layouts using the fyne.NewContainerWithLayout() constructor.",
			containerScreen,
			nil,
		},
		"apptabs": {"AppTabs",
			"A container to help divide up an application into functional areas.",
			makeAppTabsTab,
			nil,
		},
		"border": {"Border",
			"A container that positions items around a central content.",
			makeBorderLayout,
			nil,
		},
		"box": {"Box",
			"A container arranges items in horizontal or vertical list.",
			makeBoxLayout,
			nil,
		},
		"center": {"Center",
			"A container to that centers child elements.",
			makeCenterLayout,
			nil,
		},
		"doctabs": {"DocTabs",
			"A container to display a single document from a set of many.",
			makeDocTabsTab,
			nil,
		},
		"grid": {"Grid",
			"A container that arranges all items in a grid.",
			makeGridLayout,
			nil,
		},
		"split": {"Split",
			"A split container divides the container in two pieces that the user can resize.",
			makeSplitTab,
			nil,
		},
		"scroll": {"Scroll",
			"A container that provides scrolling for its content.",
			makeScrollTab,
			nil,
		},
		"innerwindow": {"InnerWindow",
			"A window that can be used inside a traditional window to contain a document or content.",
			makeInnerWindowTab,
			nil,
		},
		"widgets": {"Widgets",
			"In this section you can see the features available in the toolkit widget set.\n" +
				"Expand the tree on the left to browse the individual tutorial elements.",
			widgetScreen,
			nil,
		},
		"accordion": loadDefinition("Accordion", "widgets/accordion.md"),
		"activity":  loadDefinition("Activity", "widgets/activity.md"),
		"button":    loadDefinition("Button", "widgets/button.md"),
		"card": {"Card",
			"Group content and widgets.",
			makeCardTab,
			nil,
		},
		"entry": {"Entry",
			"Different ways to use the entry widget.",
			makeEntryTab,
			nil,
		},
		"form": {"Form",
			"Gathering input widgets for data submission.",
			makeFormTab,
			nil,
		},
		"input": {"Input",
			"A collection of widgets for user input.",
			makeInputTab,
			nil,
		},
		"text": {"Text",
			"Text handling widgets.",
			makeTextTab,
			nil,
		},
		"toolbar": {"Toolbar",
			"A row of shortcut icons for common tasks.",
			makeToolbarTab,
			nil,
		},
		"progress": {"Progress",
			"Show duration or the need to wait for a task.",
			makeProgressTab,
			nil,
		},
		"collections": {"Collections",
			"Collection widgets provide an efficient way to present lots of content.\n" +
				"The List, Table, and Tree provide a cache and re-use mechanism that make it possible to scroll through thousands of elements.\n" +
				"Use this for large data sets or for collections that can expand as users scroll.",
			collectionScreen,
			nil,
		},
		"list": {"List",
			"A vertical arrangement of cached elements with the same styling.",
			makeListTab,
			nil,
		},

		"table": {"Table",
			"A two dimensional cached collection of cells.",
			makeTableTab,
			nil,
		},
		"tree": {"Tree",
			"A tree based arrangement of cached elements with the same styling.",
			makeTreeTab,
			nil,
		},
		"gridwrap": {"GridWrap",
			"A grid based arrangement of cached elements that wraps rows to fit.",
			makeGridWrapTab,
			nil,
		},
		"dialogs": {"Dialogs",
			"Work with dialogs.",
			dialogScreen,
			nil,
		},
		"windows": {"Windows",
			"Window function demo.",
			windowScreen,
			nil,
		},
		"binding": {"Data Binding",
			"Connecting widgets to a data source.",
			bindingScreen,
			nil,
		},
		"advanced": {"Advanced",
			"Debug and advanced information.",
			advancedScreen,
			nil,
		},
	}

	// TutorialIndex  defines how our tutorials should be laid out in the index tree
	TutorialIndex = map[string][]string{
		"":            {"k8s", "canvas", "widgets", "collections", "containers", "animations", "dialogs", "windows", "icons", "binding", "advanced"},
		"k8s":         {"pods"},
		"canvas":      {"canvastext", "line", "rectangle", "circle", "image", "raster", "gradient"},
		"collections": {"list", "table", "tree", "gridwrap"},
		"containers":  {"apptabs", "border", "box", "center", "doctabs", "grid", "scroll", "split", "innerwindow"},
		"widgets":     {"accordion", "activity", "button", "card", "entry", "form", "input", "progress", "text", "toolbar"},
	}
)

func loadDefinition(title, file string) Tutorial {
	return Tutorial{title,
		"",
		func(_ fyne.Window, Data interface{}) fyne.CanvasObject { return makeNewTutorial(file) },
		nil,
	}
}

func makeNewTutorial(file string) fyne.CanvasObject {
	tutorial := tutorials[file]
	top := container.NewVBox(
		widget.NewLabel(tutorial.title), widget.NewSeparator())
	details := container.NewVBox()

	for i, p := range tutorial.content {
		if p == "" {
			continue
		}
		if i%2 == 0 {
			text := widget.NewRichTextFromMarkdown(p)
			text.Wrapping = fyne.TextWrapWord

			details.Add(text)
			continue
		}

		usage := widget.NewTextGridFromString(p + "\n")
		usage.ShowLineNumbers = true
		usage.Scroll = fyne.ScrollHorizontalOnly
		highlightTextGrid(usage)

		codeID := (i - 1) / 2
		preview := tutorial.code[codeID]()

		tools := widget.NewToolbar(
			widget.NewToolbarAction(theme.ContentCopyIcon(), func() {
				fyne.CurrentApp().Clipboard().SetContent(usage.Text())
			}),
			widget.NewToolbarAction(theme.MediaPlayIcon(), func() {
				w := fyne.CurrentApp().NewWindow(tutorial.title + " preview")
				w.SetContent(tutorial.code[codeID]())
				w.Show()
			}),
		)

		style := styles.Get("solarized-dark")
		bg := styleBackgroundColor(chroma.Background, style)
		prop := canvas.NewRectangle(color.Transparent)
		prop.SetMinSize(fyne.NewSize(18, 18)) // big enough for canvas elements to show, but smaller than most widgets

		details.Add(container.NewPadded(container.NewPadded(
			container.NewStack(
				canvas.NewRectangle(bg),
				container.NewPadded(usage,
					container.NewHBox(layout.NewSpacer(), container.NewVBox(
						tools)))))))
		details.Add(widget.NewRichTextFromMarkdown("*Preview:*"))
		details.Add(container.NewPadded(container.NewStack(prop, preview)))
	}

	return container.NewBorder(top, nil, nil, nil, container.NewScroll(details))
}
