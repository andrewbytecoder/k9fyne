package widgets

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	pullimage "github.com/andrewbytecoder/k9fyne/kube/pull_image"
	"github.com/andrewbytecoder/k9fyne/kube/statefulsets"
	"github.com/andrewbytecoder/k9fyne/utils"
	appsv1 "k8s.io/api/apps/v1"
	"strconv"
)

type PullImageWidgetsInfo struct {
	currentNameSpace int
	tags             []string
	currentImageName string
	namespaceSelect  *widget.Select // 命名空间名字
	currentPod       int
	table            *widget.Table // pod info table
	infoInterface    pullimage.KubePullImageInfoInterface
	container        *fyne.Container
}

var kubePullImageInfoCols = []string{
	"Name",
	"Ready",
	"Age",
}

func (b *PullImageWidgetsInfo) SetPullImageCurrentNameSpace(idx int) {
	// 超过正常的范围
	if idx >= len(b.tags) || idx < 0 {
		return
	}

	bFlush := false
	if b.currentNameSpace != idx {
		bFlush = true
	}

	// 保证这里只设置合法的 namespace 索引
	b.currentNameSpace = idx
	// 更新list链表里面的数据
	b.namespaceSelect.SetSelected(b.tags[idx])

	if bFlush {
		////	 表数据
		//list, err := b.infoInterface.GetPullImageInfoByNamespace(b.currentImageName, b.tags[b.currentNameSpace])
		//if err != nil {
		//	return
		//}
		//table := makePullImageInfoTable(nil, list)
		//b.container.Remove(b.table)
		//b.container.Add(table)
		//b.table = table
		//b.container.Refresh()
	}
}

func makePullImageList(win fyne.Window, d interface{}) fyne.CanvasObject {
	serviceInterface, ok := d.(pullimage.KubePullImageInfoInterface)
	if !ok {
		return container.NewHSplit(widget.NewButton("podInterface is null", func() {}),
			widget.NewButton("Set ssh config first", func() {}))
	}

	b := &PullImageWidgetsInfo{}
	b.tags = serviceInterface.GetAllTags(b.currentImageName)
	b.infoInterface = serviceInterface
	prev := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		b.SetPullImageCurrentNameSpace(b.currentNameSpace - 1)
	})
	next := widget.NewButtonWithIcon("", theme.NavigateNextIcon(), func() {
		b.SetPullImageCurrentNameSpace(b.currentNameSpace + 1)
	})
	// 选择并更新 namespace
	b.namespaceSelect = widget.NewSelect(b.namespace, func(name string) {
		for i, namespace := range b.namespace {
			if namespace == name {
				if b.currentNameSpace != i {
					b.SetPullImageCurrentNameSpace(i)
				}
				break
			}
		}
	})
	b.namespaceSelect.SetSelected(b.namespace[b.currentNameSpace])
	buttons := container.NewHBox(prev, next)
	bar := container.NewBorder(nil, nil, buttons, nil, b.namespaceSelect)

	podList, err := serviceInterface.GetPullImageInfoByNamespace(b.namespace[b.currentNameSpace])
	if err != nil {
		return container.NewHSplit(widget.NewButton("serviceInterface is null", func() {}),
			widget.NewButton("Set ssh config first", func() {}))
	}

	b.table = makePullImageInfoTable(nil, podList)
	b.container = container.NewBorder(bar, nil, nil, nil, b.table)
	return b.container
}

func makePullImageInfoTable(_ fyne.Window, list *appsv1.StatefulSetList) *widget.Table {
	rows := len(list.Items)
	cols := len(statefulSetsInfoCols)
	t := widget.NewTableWithHeaders(
		func() (int, int) { return rows, cols },
		func() fyne.CanvasObject {
			return widget.NewLabel("Cell")
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)
			item := list.Items[id.Row]
			switch id.Col {
			case 0:
				label.SetText(item.Name)
			case 1:
				label.SetText(strconv.Itoa(int(item.Status.ReadyReplicas)))
			case 2:
				label.SetText(utils.TimeFormat(list.Items[id.Row].CreationTimestamp.Time))
			default:
				label.SetText(fmt.Sprintf("Cell %d, %d", id.Row+1, id.Col+1))
			}
		})
	t.CreateHeader = func() fyne.CanvasObject {
		return widget.NewLabel("Header")
	}
	// 更新table的标题头
	t.UpdateHeader = func(id widget.TableCellID, cell fyne.CanvasObject) {
		l := cell.(*widget.Label)
		if id.Row < 0 {
			// Col 这里从0开始
			l.SetText(statefulSetsInfoCols[id.Col])
		} else if id.Col < 0 {
			l.SetText(strconv.Itoa(id.Row + 1))
		} else {
			l.SetText("")
		}

		//label.SetText(podInfoCols[id.Col])
	}

	t.StickyRowCount = 0

	t.SetColumnWidth(0, 320)
	t.SetColumnWidth(1, 80)
	t.SetColumnWidth(2, 230)

	t.SetRowHeight(2, 50)

	return t
}
