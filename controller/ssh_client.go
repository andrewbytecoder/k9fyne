package controller

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/andrewbytecoder/k9fyne/utils"
	"github.com/melbahja/goph"
	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"net"
	"net/url"
	"os"
	"strings"
	"time"
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
func VerifyHost(host string, remote net.Addr, key ssh.PublicKey) error {

	hostFound, err := goph.CheckKnownHost(host, remote, key, "")
	if hostFound && err != nil {
		return err
	}
	if hostFound && err == nil {
		return nil
	}

	// Add the new host to known hosts file.
	return goph.AddKnownHost(host, remote, key, "")
}

// CreateSSHClient create ssh config panel
func (c *SSHClient) CreateSSHClient(win fyne.Window) {
	username := widget.NewEntry()
	username.Validator = validation.NewRegexp(`^[A-Za-z0-9_-]+$`, "username can only contain letters, numbers, '_', and '-'")
	username.SetText(c.UserName)
	address := widget.NewEntry()
	address.SetText(c.Address)
	password := widget.NewPasswordEntry()
	//password.Validator = validation.NewRegexp(`^[A-Za-z0-9_-@]+$`, "password can only contain letters, numbers, '_', and '-'")
	password.SetText(c.Password)

	remember := false
	items := []*widget.FormItem{
		widget.NewFormItem("Address", address),
		widget.NewFormItem("Username", username),
		widget.NewFormItem("Password", password),
		widget.NewFormItem("Remember me", widget.NewCheck("", func(checked bool) {
			remember = checked
		})),
	}

	dialog.ShowForm("Login...", "Log In", "Cancel", items, func(b bool) {
		if !b {
			return
		}
		var rememberText string
		if remember {
			rememberText = "and remember this login"
		}

		log.Println("Please Authenticate", username.Text, password.Text, rememberText)
		// save the ssh config info
		c.UserName = username.Text
		c.Password = password.Text
		c.Address = address.Text
		if remember {
			fyne.CurrentApp().Preferences().SetString("username", c.UserName)
			fyne.CurrentApp().Preferences().SetString("password", c.Password)
			fyne.CurrentApp().Preferences().SetString("address", c.Address)
		} else {
			fyne.CurrentApp().Preferences().RemoveValue("username")
			fyne.CurrentApp().Preferences().RemoveValue("password")
			fyne.CurrentApp().Preferences().RemoveValue("address")
		}

		if c.client != nil {
			c.client.Close()
			c.client = nil
		}

		parseAddress, port, err := utils.ParseAddress(c.Address)
		if err != nil {
			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Title:   "SSH Connect Failed",
				Content: err.Error(),
			})
			return
		}

		config := &goph.Config{
			User:     username.Text,
			Addr:     parseAddress,
			Auth:     goph.Password(c.Password),
			Port:     uint(port),
			Timeout:  20 * time.Second,
			Callback: VerifyHost,
		}

		client, err := goph.NewConn(config)
		if err != nil {
			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Title:   "SSH Connect Failed",
				Content: err.Error(),
			})
		} else {
			c.client = client

		}
		err = c.GetKubeConfig()
		if err != nil {
			fyne.CurrentApp().SendNotification(&fyne.Notification{
				Title:   "Get kube config failed",
				Content: err.Error(),
			})
		}

	}, win)
}

func (c *SSHClient) GetClient() *goph.Client {
	return c.client
}

// KubeConfig 定义结构体以解析 YAML 数据
type KubeConfig struct {
	APIVersion string `yaml:"apiVersion"`
	Clusters   []struct {
		Name    string `yaml:"name"`
		Cluster struct {
			Server                   string `yaml:"server"`
			CertificateAuthorityData string `yaml:"certificate-authority-data"`
		} `yaml:"cluster"`
	} `yaml:"clusters"`
	Contexts []struct {
		Name    string `yaml:"name"`
		Context struct {
			Cluster string `yaml:"cluster"`
			User    string `yaml:"user"`
		} `yaml:"context"`
	} `yaml:"contexts"`
	CurrentContext string   `yaml:"current-context"`
	Kind           string   `yaml:"kind"`
	Preferences    struct{} `yaml:"preferences"`
	Users          []struct {
		Name string `yaml:"name"`
		User struct {
			ClientCertificateData string `yaml:"client-certificate-data"`
			ClientKeyData         string `yaml:"client-key-data"`
		} `yaml:"user"`
	} `yaml:"users"`
}

func replaceHostInURL(ipWithPort, originalURL string) (string, error) {
	// 解析原始 URL
	parsedURL, err := url.Parse(originalURL)
	if err != nil {
		return "", fmt.Errorf("无法解析原始 URL: %v", err)
	}

	// 提取 IP 和端口
	ipPortParts := strings.Split(ipWithPort, ":")
	if len(ipPortParts) != 2 {
		return "", fmt.Errorf("无效的 IP 和端口格式，应为 'IP:Port'")
	}
	ip := ipPortParts[0]

	// 替换 Host（IP 和端口）
	parsedURL.Host = fmt.Sprintf("%s:%s", ip, parsedURL.Port())

	// 返回修改后的 URL 字符串
	return parsedURL.String(), nil
}
func (c *SSHClient) GetKubeConfig() error {
	if c.client == nil {
		return fmt.Errorf("ssh client is nil")
	}

	sftp, err := c.client.NewSftp()
	if err != nil {
		return err
	}
	defer sftp.Close()

	kubeConfigFile, err := sftp.OpenFile("/root/.kube/config", os.O_RDONLY)
	if err != nil {
		return err
	}
	defer kubeConfigFile.Close()

	fileIfo, err := kubeConfigFile.Stat()
	if err != nil {
		return err
	}

	var buffer = make([]byte, fileIfo.Size())
	if _, err := kubeConfigFile.Read(buffer); err != nil && err != io.EOF {
		return err
	}
	// 解析 YAML 内容到结构体
	var kubeConfig KubeConfig
	err = yaml.Unmarshal([]byte(buffer), &kubeConfig)
	if err != nil {
		log.Fatalf("无法解析 YAML: %v", err)
		return err
	}

	// 修改 server 地址
	for i := range kubeConfig.Clusters {
		if kubeConfig.Clusters[i].Cluster.Server == "https://apiserver.cluster.local:6443" {
			kubeConfig.Clusters[i].Cluster.Server, err = replaceHostInURL(c.Address, kubeConfig.Clusters[i].Cluster.Server)
			if err != nil {
				return err
			}
		}
	}

	// 将修改后的结构体序列化回 YAML
	modifiedYAML, err := yaml.Marshal(&kubeConfig)
	if err != nil {
		log.Fatalf("无法序列化 YAML: %v", err)
		return err
	}

	// 将修改后的内容保存到文件
	err = os.WriteFile("config.yaml", modifiedYAML, 0644)
	if err != nil {
		log.Fatalf("无法写入文件: %v", err)
		return err
	}

	return nil
}
