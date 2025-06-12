package widgets

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/andrewbytecoder/k9fyne/kube/pod"
	"github.com/andrewbytecoder/k9fyne/utils"
	corev1 "k8s.io/api/core/v1"
	"strconv"
	"strings"
)

type PodWidgetsInfo struct {
	currentNameSpace int
	namespace        []string
	namespaceSelect  *widget.Select // 命名空间名字
	currentPod       int
	table            *widget.Table // pod info table
	infoInterface    pod.KubePodInfoInterface
	container        *fyne.Container
}

var podInfoCols = []string{
	"Name",
	"Status",
	"PodIp",
	"NodeIp",
	"NodeName",
	"Age",
	"Containers",
}

func (b *PodWidgetsInfo) SetCurrentNameSpace(idx int) {
	// 超过正常的范围
	if idx >= len(b.namespace) || idx < 0 {
		return
	}

	bFlush := false
	if b.currentNameSpace != idx {
		bFlush = true
	}

	// 保证这里只设置合法的 namespace 索引
	b.currentNameSpace = idx
	// 更新list链表里面的数据
	b.namespaceSelect.SetSelected(b.namespace[idx])

	if bFlush {
		//	 表数据
		podList, err := b.infoInterface.GetPodInfoByNamespace(b.namespace[b.currentNameSpace])
		if err != nil {
			return
		}
		table := makePodInfoTable(nil, podList)
		b.container.Remove(b.table)
		b.container.Add(table)
		b.table = table

		b.container.Refresh()
	}
}

func makePodList(win fyne.Window, d interface{}) fyne.CanvasObject {

	podInterface, ok := d.(pod.KubePodInfoInterface)
	if !ok {
		return container.NewHSplit(widget.NewButton("podInterface is null", func() {}),
			widget.NewButton("Set ssh config first", func() {}))
	}

	b := &PodWidgetsInfo{}
	b.namespace = podInterface.GetAllNamespace()
	b.infoInterface = podInterface
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

	b.table = makePodInfoTable(nil, podList)
	b.container = container.NewBorder(bar, nil, nil, nil, b.table)
	return b.container
}

func makePodInfoTable(_ fyne.Window, podList *corev1.PodList) *widget.Table {
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
			case 5:
				label.SetText(utils.TimeFormat(podList.Items[id.Row].ObjectMeta.CreationTimestamp.Time))
			case 6:
				label.SetText(GetContainerInfo(&podList.Items[id.Row]))
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

	t.SetColumnWidth(0, 380)
	t.SetColumnWidth(1, 100)
	t.SetColumnWidth(2, 130)
	t.SetColumnWidth(3, 130)
	t.SetColumnWidth(4, 130)
	t.SetColumnWidth(5, 200)
	t.SetColumnWidth(6, 500)
	t.SetRowHeight(2, 50)

	return t
}

func GetContainerInfo(pod *corev1.Pod) string {
	var result strings.Builder

	// 遍历所有容器状态
	for _, containerStatus := range pod.Status.ContainerStatuses {
		result.WriteString(fmt.Sprintf("容器名称: %s", containerStatus.Name))
		result.WriteString(fmt.Sprintf(":状态: %v", GetContainerState(&containerStatus.State)))
		result.WriteString(fmt.Sprintf(":重启次数: %d", containerStatus.RestartCount))
	}

	return result.String()
}

func GetContainerState(state *corev1.ContainerState) string {

	if state.Running != nil {
		return "Running"
	}
	if state.Terminated != nil {
		return "Terminated"
	}
	if state.Waiting != nil {
		return "Waiting"
	}
	return ""
}
