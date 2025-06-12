package resources

// This file embeds all the resources used by the program.

import (
	_ "embed"
	"fyne.io/fyne/v2"
)

//go:embed fire.png
var embedFirePng []byte
var K9FyneFireLogo = fyne.NewStaticResource("k9fyneFireLogo", embedFirePng)

//go:embed weixin.png
var weChatPic []byte
var WeChat = fyne.NewStaticResource("wechat", weChatPic)

//go:embed k8s.png
var embedLogoIconPng []byte
var K9FyneLogo = fyne.NewStaticResource("k9fynelogo", embedLogoIconPng)
