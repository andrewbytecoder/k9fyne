package widgets

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	pullimage "github.com/andrewbytecoder/k9fyne/kube/pull_image"
	"time"
)

type PullImageWidgetsInfo struct {
	currentNameSpace int
	tags             []string
	searchImageName  string // 查找的镜像名称，可能是非法的，真正下载需要重新构造
	tagIndex         int
	namespaceSelect  *widget.Select // 命名空间名字
	currentPod       int
	split            *container.Split
	infoInterface    pullimage.KubePullImageInfoInterface
	container        *fyne.Container
	progressBar      *widget.ProgressBar
}

func makePullImageList(win fyne.Window, d interface{}) fyne.CanvasObject {
	serviceInterface, ok := d.(pullimage.KubePullImageInfoInterface)
	if !ok {
		return container.NewHSplit(widget.NewButton("podInterface is null", func() {}),
			widget.NewButton("Set ssh config first", func() {}))
	}

	b := &PullImageWidgetsInfo{}
	b.tags = serviceInterface.GetAllTags(b.searchImageName)
	b.infoInterface = serviceInterface
	imageInput := widget.NewEntry()
	imageInput.SetPlaceHolder("Image Name")
	imageInput.OnChanged = func(s string) {
		b.searchImageName = s
	}
	b.progressBar = widget.NewProgressBar()
	searchButton := widget.NewButtonWithIcon("Search", theme.SearchIcon(), func() {
		b.tags = serviceInterface.GetAllTags(b.searchImageName)
		b.container.Remove(b.split)
		b.split = b.makeList(b.tags)
		b.container.Add(b.split)
	})
	bar := container.NewBorder(nil, nil, nil, searchButton, imageInput)

	b.split = b.makeList(b.tags)
	b.container = container.NewBorder(bar, nil, nil, nil, b.split)
	return b.container
}

func (puller *PullImageWidgetsInfo) makeList(data []string) *container.Split {
	if data == nil || len(data) == 0 {
		return container.NewHSplit(widget.NewButton("put image name first", func() {}),
			widget.NewButtonWithIcon("pull images", theme.DownloadIcon(), func() {
				fmt.Println("put image name first")
			}))
	}

	icon := widget.NewIcon(nil)
	label := widget.NewLabel("Select An Item From The List")
	downloadButton := widget.NewButtonWithIcon("pull images", theme.DownloadIcon(), func() {
		fmt.Println("put image name first")
		progressAnimation := fyne.Animation{
			Curve:    fyne.AnimationLinear,
			Duration: 10 * time.Second,
			Tick: func(percentage float32) {
				value := float64(percentage)
				if value >= 99 {
					puller.progressBar.SetValue(99)
				} else {
					puller.progressBar.SetValue(value)
				}
			},
		}
		progressAnimation.Start()
		go func() {
			err := puller.infoInterface.PullImage(puller.searchImageName, data[puller.tagIndex], puller.progressBar)
			if err != nil {
				fmt.Println(err)
				return
			}
		}()

	})

	hbox := container.NewVBox(container.NewHBox(label, downloadButton), puller.progressBar)

	list := widget.NewList(
		func() int {
			return len(data)
		},
		func() fyne.CanvasObject {
			return container.NewHBox(widget.NewIcon(theme.DocumentIcon()), widget.NewLabel("Template Object"))
		},
		func(id widget.ListItemID, item fyne.CanvasObject) {
			item.(*fyne.Container).Objects[1].(*widget.Label).SetText(data[id])
		},
	)
	list.OnSelected = func(id widget.ListItemID) {
		label.SetText(data[id])
		icon.SetResource(theme.DocumentIcon())
		// 下载程序
		puller.tagIndex = id
		puller.progressBar.SetValue(0)
	}
	list.OnUnselected = func(id widget.ListItemID) {
		label.SetText("Select An Item From The List")
		icon.SetResource(nil)
	}
	list.Select(125)
	list.SetItemHeight(5, 50)
	list.SetItemHeight(6, 50)

	return container.NewHSplit(list, container.NewCenter(hbox))
}

func (puller *PullImageWidgetsInfo) makeProgressTab() fyne.CanvasObject {
	progress := widget.NewProgressBar()

	fprogress := widget.NewProgressBar()
	fprogress.TextFormatter = func() string {
		return fmt.Sprintf("%.2f out of %.2f", fprogress.Value, fprogress.Max)
	}

	infProgress := widget.NewProgressBarInfinite()

	progressAnimation := fyne.Animation{
		Curve:    fyne.AnimationLinear,
		Duration: 10 * time.Second,
		Tick: func(percentage float32) {
			value := float64(percentage)
			progress.SetValue(value)
			fprogress.SetValue(value)
		},
	}

	progressAnimation.Start()

	OnChangeFuncs = append(OnChangeFuncs, progressAnimation.Stop, infProgress.Stop)

	return container.NewVBox(
		widget.NewLabel("Percent"), progress,
		widget.NewLabel("Formatted"), fprogress,
		widget.NewLabel("Infinite"), infProgress)
}

func makeToolbarTab(_ fyne.Window, d interface{}) fyne.CanvasObject {
	t := widget.NewToolbar(widget.NewToolbarAction(theme.FileImageIcon(), func() { fmt.Println("New") }),
		widget.NewToolbarSeparator(),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.SearchIcon(), func() { fmt.Println("find") }),
	)

	return container.NewBorder(t, nil, nil, nil)
}
