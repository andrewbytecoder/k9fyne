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
		"service": {"Service info",
			"Show service info.",
			makeServiceList, nil,
		},
		"deployment": {"Deployment info",
			"Show deployment info.",
			makeDeploymentList, nil,
		},
		"daemonSets": {"DaemonSets info",
			"Show daemonSets info.",
			makeDaemonSetsList, nil,
		},
		"statefulSets": {"StatefulSets info",
			"Show statefulSets info.",
			makeStatefulSetsList, nil,
		},
		"image": {"Image",
			"镜像管理.",
			makeToolbarTab,
			nil,
		},
		"pullImage": {"PullImage info",
			"Show pullImage info.",
			makePullImageList, nil,
		},
		"doc": {"Doc",
			"使用指导文档和常见命令.",
			k9Doc,
			nil,
		},
		"readme":     loadDefinition("Readme", "docs/readme.md"),
		"canvastext": loadDefinition("Text", "docs/text.md"),
		"line":       loadDefinition("Line", "docs/line.md"),
		"rectangle":  loadDefinition("Rectangle", "docs/rectangle.md"),
		"circle":     loadDefinition("Circle", "docs/circle.md"),
		"raster":     loadDefinition("Raster", "docs/raster.md"),
		"gradient":   loadDefinition("Gradient", "docs/gradient.md"),
	}

	// TutorialIndex  defines how our tutorials should be laid out in the index tree
	TutorialIndex = map[string][]string{
		"":      {"k8s", "image", "doc"},
		"k8s":   {"pods", "service", "deployment", "daemonSets", "statefulSets"},
		"image": {"pullImage"},
		"doc":   {"canvastext", "line", "rectangle", "circle", "raster", "gradient"},
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

		style := styles.Get("solarized-dark")
		bg := styleBackgroundColor(chroma.Background, style)
		prop := canvas.NewRectangle(color.Transparent)
		prop.SetMinSize(fyne.NewSize(18, 18)) // big enough for canvas elements to show, but smaller than most widgets

		details.Add(container.NewPadded(container.NewPadded(
			container.NewStack(
				canvas.NewRectangle(bg),
				container.NewPadded(usage,
					container.NewHBox(layout.NewSpacer()))))))

		details.Add(container.NewPadded(container.NewStack(prop)))
	}

	return container.NewBorder(top, nil, nil, nil, container.NewScroll(details))
}
