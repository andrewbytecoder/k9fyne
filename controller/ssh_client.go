package controller

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

// SSHClient ssh client info
type SSHClient struct {
	UserName string
	Password string
	Address  string
}

func NewSSHClient() *SSHClient {
	return &SSHClient{
		UserName: fyne.CurrentApp().Preferences().String("username"),
		Password: fyne.CurrentApp().Preferences().String("password"),
		Address:  fyne.CurrentApp().Preferences().String("address"),
	}
}

// CreateSSHClient create ssh config panel
func CreateSSHClient(c *SSHClient) fyne.CanvasObject {
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

	}
	content := container.NewVBox(widget.NewLabelWithStyle("SSH Panel", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		form)

	return content
}
