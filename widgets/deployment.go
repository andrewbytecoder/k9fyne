package widgets

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/andrewbytecoder/k9fyne/kube/deployment"
	"github.com/andrewbytecoder/k9fyne/utils"
	appsv1 "k8s.io/api/apps/v1"
	"strconv"
)

type DeploymentWidgetsInfo struct {
	currentNameSpace        int
	namespace               []string
	namespaceSelect         *widget.Select // 命名空间名字
	currentPod              int
	table                   *widget.Table // pod info table
	deploymentInfoInterface deployment.KubeDeploymentInfoInterface
	container               *fyne.Container
}

var deploymentInfoCols = []string{
	"Name",
	"Ready",
	"UP-TO-DATE",
	"Available",
	"Age",
}

func (b *DeploymentWidgetsInfo) SetDeploymentCurrentNameSpace(idx int) {
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
		list, err := b.deploymentInfoInterface.GetDeploymentInfoByNamespace(b.namespace[b.currentNameSpace])
		if err != nil {
			return
		}
		table := makeDeploymentInfoTable(nil, list)
		b.container.Remove(b.table)
		b.container.Add(table)
		b.table = table
		b.container.Refresh()
	}
}

func makeDeploymentList(win fyne.Window, d interface{}) fyne.CanvasObject {
	serviceInterface, ok := d.(deployment.KubeDeploymentInfoInterface)
	if !ok {
		return container.NewHSplit(widget.NewButton("podInterface is null", func() {}),
			widget.NewButton("Set ssh config first", func() {}))
	}

	b := &DeploymentWidgetsInfo{}
	b.namespace = serviceInterface.GetAllNamespace()
	b.deploymentInfoInterface = serviceInterface
	prev := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		b.SetDeploymentCurrentNameSpace(b.currentNameSpace - 1)
	})
	next := widget.NewButtonWithIcon("", theme.NavigateNextIcon(), func() {
		b.SetDeploymentCurrentNameSpace(b.currentNameSpace + 1)
	})
	// 选择并更新 namespace
	b.namespaceSelect = widget.NewSelect(b.namespace, func(name string) {
		for i, namespace := range b.namespace {
			if namespace == name {
				if b.currentNameSpace != i {
					b.SetDeploymentCurrentNameSpace(i)
				}
				break
			}
		}
	})
	b.namespaceSelect.SetSelected(b.namespace[b.currentNameSpace])
	buttons := container.NewHBox(prev, next)
	bar := container.NewBorder(nil, nil, buttons, nil, b.namespaceSelect)

	podList, err := serviceInterface.GetDeploymentInfoByNamespace(b.namespace[b.currentNameSpace])
	if err != nil {
		return container.NewHSplit(widget.NewButton("serviceInterface is null", func() {}),
			widget.NewButton("Set ssh config first", func() {}))
	}

	b.table = makeDeploymentInfoTable(nil, podList)
	b.container = container.NewBorder(bar, nil, nil, nil, b.table)
	return b.container
}

func makeDeploymentInfoTable(_ fyne.Window, list *appsv1.DeploymentList) *widget.Table {

	rows := len(list.Items)
	cols := len(deploymentInfoCols)
	t := widget.NewTableWithHeaders(
		func() (int, int) { return rows, cols },
		func() fyne.CanvasObject {
			return widget.NewLabel("Cell")
		},
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			label := cell.(*widget.Label)
			switch id.Col {
			case 0:
				label.SetText(list.Items[id.Row].Name)
			case 1:
				ready := fmt.Sprintf("%d/%d", list.Items[id.Row].Status.ReadyReplicas, *list.Items[id.Row].Spec.Replicas)
				label.SetText(ready)
			case 2:
				label.SetText(strconv.Itoa(int(list.Items[id.Row].Status.UpdatedReplicas)))
			case 3:
				label.SetText(strconv.Itoa(int(list.Items[id.Row].Status.AvailableReplicas)))
			case 4:
				label.SetText(utils.TimeFormat(list.Items[id.Row].CreationTimestamp.Time))
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
			l.SetText(deploymentInfoCols[id.Col])
		} else if id.Col < 0 {
			l.SetText(strconv.Itoa(id.Row + 1))
		} else {
			l.SetText("")
		}

		//label.SetText(podInfoCols[id.Col])
	}

	t.StickyRowCount = 0

	t.SetColumnWidth(0, 320)
	t.SetColumnWidth(1, 100)
	t.SetColumnWidth(2, 130)
	t.SetColumnWidth(3, 80)
	t.SetColumnWidth(4, 230)
	t.SetColumnWidth(5, 230)
	t.SetRowHeight(2, 50)

	return t
}
