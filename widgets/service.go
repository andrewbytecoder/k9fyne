package widgets

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/andrewbytecoder/k9fyne/kube/service"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strconv"
	"strings"
	"time"
)

type ServiceWidgetsInfo struct {
	currentNameSpace     int
	namespace            []string
	namespaceSelect      *widget.Select // 命名空间名字
	currentPod           int
	table                *widget.Table // pod info table
	serviceInfoInterface service.KubeServiceInfoInterface
	container            *fyne.Container
}

var serviceInfoCols = []string{
	"Name",
	"Type",
	"ClusterIp",
	"ExternalIp",
	"Port(s)",
	"Age",
}

func (b *ServiceWidgetsInfo) SetServiceCurrentNameSpace(idx int) {
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
		list, err := b.serviceInfoInterface.GetServiceInfoByNamespace(b.namespace[b.currentNameSpace])
		if err != nil {
			return
		}
		table := makeServiceInfoTable(nil, list)
		b.container.Remove(b.table)
		b.container.Add(table)
		b.table = table
		b.container.Refresh()
	}
}

func makeServiceList(win fyne.Window, d interface{}) fyne.CanvasObject {
	serviceInterface, ok := d.(service.KubeServiceInfoInterface)
	if !ok {
		return container.NewHSplit(widget.NewButton("podInterface is null", func() {}),
			widget.NewButton("Set ssh config first", func() {}))
	}

	b := &ServiceWidgetsInfo{}
	b.namespace = serviceInterface.GetAllNamespace()
	b.serviceInfoInterface = serviceInterface
	prev := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		b.SetServiceCurrentNameSpace(b.currentNameSpace - 1)
	})
	next := widget.NewButtonWithIcon("", theme.NavigateNextIcon(), func() {
		b.SetServiceCurrentNameSpace(b.currentNameSpace + 1)
	})
	// 选择并更新 namespace
	b.namespaceSelect = widget.NewSelect(b.namespace, func(name string) {
		for i, namespace := range b.namespace {
			if namespace == name {
				if b.currentNameSpace != i {
					b.SetServiceCurrentNameSpace(i)
				}
				break
			}
		}
	})
	b.namespaceSelect.SetSelected(b.namespace[b.currentNameSpace])
	buttons := container.NewHBox(prev, next)
	bar := container.NewBorder(nil, nil, buttons, nil, b.namespaceSelect)

	podList, err := serviceInterface.GetServiceInfoByNamespace(b.namespace[b.currentNameSpace])
	if err != nil {
		return container.NewHSplit(widget.NewButton("serviceInterface is null", func() {}),
			widget.NewButton("Set ssh config first", func() {}))
	}

	b.table = makeServiceInfoTable(nil, podList)
	b.container = container.NewBorder(bar, nil, nil, nil, b.table)
	return b.container
}

func makeServiceInfoTable(_ fyne.Window, list *corev1.ServiceList) *widget.Table {

	rows := len(list.Items)
	cols := len(serviceInfoCols)
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
				label.SetText(string(list.Items[id.Row].Spec.Type))
			case 2:
				label.SetText(list.Items[id.Row].Spec.ClusterIP)
			case 3:
				label.SetText(strings.Join(list.Items[id.Row].Spec.ExternalIPs, ","))
			case 4:
				label.SetText(generateServicePortsInfo(&list.Items[id.Row].Spec))
			case 5:
				label.SetText(calculateAge(list.Items[id.Row].CreationTimestamp))
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
			l.SetText(serviceInfoCols[id.Col])
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
func generateServicePortsInfo(spec *corev1.ServiceSpec) string {
	var portStrings []string
	for _, port := range spec.Ports {
		var portStr string
		if port.NodePort == 0 {
			portStr = fmt.Sprintf("%d/%s", port.Port, port.Protocol)
		} else {
			portStr = fmt.Sprintf("%d:%d/%s", port.Port, port.NodePort, port.Protocol)
		}
		portStrings = append(portStrings, portStr)
	}
	return strings.Join(portStrings, ",")
}
func calculateAge(creationTime metav1.Time) string {
	d := time.Since(creationTime.Time).Round(time.Second)
	switch {
	case d < time.Minute:
		return fmt.Sprintf("%ds", int(d.Seconds()))
	case d < time.Hour:
		return fmt.Sprintf("%dm", int(d.Minutes()))
	default:
		hours := int(d.Hours())
		minutes := int(d.Minutes()) % 60
		return fmt.Sprintf("%dh%dm", hours, minutes)
	}
}
