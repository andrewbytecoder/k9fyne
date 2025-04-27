package widgets

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/andrewbytecoder/k9fyne/kube/pod"
	corev1 "k8s.io/api/core/v1"
	"strconv"
)

type PodWidgetsInfo struct {
	currentNameSpace int
	namespace        []string
	namespaceSelect  *widget.Select // 命名空间名字
	currentPod       int
}

var podInfoCols = []string{
	"Name",
	"Status",
	"PodIp",
	"NodeIp",
	"NodeName",
}

func (p *PodWidgetsInfo) SetCurrentNameSpace(idx int) {
	// 超过正常的范围
	if idx >= len(p.namespace) || idx < 0 {
		return
	}
	// 保证这里只设置合法的 namespace 索引
	p.currentNameSpace = idx
	// 更新list链表里面的数据
	p.namespaceSelect.SetSelected(p.namespace[idx])
}

func makePodList(win fyne.Window, d interface{}) fyne.CanvasObject {

	podInterface, ok := d.(pod.KubePodInfoInterface)
	if !ok {
		return container.NewHSplit(widget.NewButton("podInterface is null", func() {}),
			widget.NewButton("Set ssh config first", func() {}))
	}

	b := &PodWidgetsInfo{}
	b.namespace = podInterface.GetAllNamespace()

	prev := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		b.SetCurrentNameSpace(b.currentNameSpace - 1)
	})
	next := widget.NewButtonWithIcon("", theme.NavigateNextIcon(), func() {
		b.SetCurrentNameSpace(b.currentNameSpace + 1)
	})
	// 选择并更新 namespace
	b.namespaceSelect = widget.NewSelect(b.namespace, func(name string) {
		for i, namespace := range b.namespace {
			if namespace == name {
				if b.currentNameSpace != i {
					b.SetCurrentNameSpace(i)
				}
				break
			}
		}
	})
	b.namespaceSelect.SetSelected(b.namespace[b.currentNameSpace])
	buttons := container.NewHBox(prev, next)
	bar := container.NewBorder(nil, nil, buttons, nil, b.namespaceSelect)

	podList, err := podInterface.GetPodInfoByNamespace(b.namespace[b.currentNameSpace])
	if err != nil {
		return container.NewHSplit(widget.NewButton("podInterface is null", func() {}),
			widget.NewButton("Set ssh config first", func() {}))
	}

	podTable := makePodInfoTable(win, podList)

	return container.NewBorder(bar, nil, nil, nil, podTable)
}

func makePodInfoTable(_ fyne.Window, podList *corev1.PodList) fyne.CanvasObject {

	rows := len(podList.Items)
	cols := len(podInfoCols)
	t := widget.NewTableWithHeaders(
		func() (int, int) { return rows, cols },
		func() fyne.CanvasObject {
			return widget.NewLabel("Cell")
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)
			switch id.Col {
			case 0:
				label.SetText(podList.Items[id.Row].Name)
			case 1:
				label.SetText(string(podList.Items[id.Row].Status.Phase))
			case 2:
				label.SetText(podList.Items[id.Row].Status.PodIP)
			case 3:
				label.SetText(podList.Items[id.Row].Status.HostIP)
			case 4:
				label.SetText(podList.Items[id.Row].Spec.NodeName)
			default:
				label.SetText(fmt.Sprintf("Cell %d, %d", id.Row+1, id.Col+1))
			}
		})
	t.CreateHeader = func() fyne.CanvasObject {
		return widget.NewLabel("Header")
	}

	t.UpdateHeader = func(id widget.TableCellID, cell fyne.CanvasObject) {
		l := cell.(*widget.Label)
		if id.Row < 0 {
			// Col 这里从0开始
			l.SetText(podInfoCols[id.Col])
		} else if id.Col < 0 {
			l.SetText(strconv.Itoa(id.Row + 1))
		} else {
			l.SetText("")
		}

		//label.SetText(podInfoCols[id.Col])
	}

	t.StickyRowCount = 0

	t.SetColumnWidth(0, 420)
	t.SetColumnWidth(1, 132)
	t.SetColumnWidth(2, 320)
	t.SetColumnWidth(3, 320)
	t.SetColumnWidth(4, 260)
	t.SetRowHeight(2, 50)

	return t
}
