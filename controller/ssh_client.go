package controller

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/melbahja/goph"
)

// SSHClient ssh client info
type SSHClient struct {
	UserName string
	Password string
	Address  string
	client   *goph.Client
}

func NewSSHClient() *SSHClient {
	return &SSHClient{
		UserName: fyne.CurrentApp().Preferences().String("username"),
		Password: fyne.CurrentApp().Preferences().String("password"),
		Address:  fyne.CurrentApp().Preferences().String("address"),
	}
}

// CreateSSHClient create ssh config panel
func (c *SSHClient) CreateSSHClient() fyne.CanvasObject {
	username := widget.NewEntry()
	username.SetText(c.UserName)
	username.Show()
	address := widget.NewEntry()
	address.SetText(c.Address)
	address.Show()
	password := widget.NewPasswordEntry()
	password.SetText(c.Password)
	password.Show()

	form := widget.NewForm(widget.NewFormItem("Address", address),
		widget.NewFormItem("Username", username),
		widget.NewFormItem("Password", password))

	form.OnCancel = func() {

	}
	form.OnSubmit = func() {
		// save the ssh config info
		c.UserName = username.Text
		c.Password = password.Text
		c.Address = address.Text
		fyne.CurrentApp().Preferences().SetString("username", c.UserName)
		fyne.CurrentApp().Preferences().SetString("password", c.Password)
		fyne.CurrentApp().Preferences().SetString("address", c.Address)

		if c.client != nil {
			c.client.Close()
			c.client = nil
		}

		auth := goph.Password(c.Password)
		client, err := goph.New(c.UserName, c.Address, auth)
		if err != nil {
			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Title:   "SSH Connect Failed",
				Content: err.Error(),
			})
		} else {
			c.client = client
			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Title:   "SSH Connect Success",
				Content: "SSH Connect Success",
			})
		}
	}
	content := container.NewVBox(widget.NewLabelWithStyle("SSH Panel", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		form)

	return content
}
