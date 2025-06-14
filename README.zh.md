# k9fyne

#### k9fyne - 使用 Fyne 框架实现的 Kubernetes 管理工具

![img.png](resources/k9fyne.png)

欢迎使用 **k9fyne**，这是一个强大且用户友好的 Kubernetes 集群管理工具。通过 Fyne 框架构建，k9fyne 提供了一个直观的图形界面，简化了 Kubernetes 资源的管理。

---

#### 功能

- **Pod 信息**: 查看集群中所有 Pod 的详细信息。
- **服务信息**: 管理和检查在 Kubernetes 环境中运行的服务。
- **部署信息**: 监控和控制部署，确保应用程序顺利发布。
- **DaemonSet 信息**: 获取 DaemonSet 的洞察，确保关键系统组件始终运行。
- **StatefulSet 信息**: 管理具有持久存储需求的状态化应用的 StatefulSet。
- **文档**: 直接从应用程序访问全面的文档以快速参考。
- **深色/浅色模式**: 切换深色和浅色主题以获得最佳观看舒适度。

---

#### 开始使用

1. **安装**:
    - 克隆仓库: `git clone https://github.com/yourusername/k9fyne.git`
    - 构建应用程序: `go build -o k9fyne main.go`

2. **处理依赖**
   - 如果你在国内使用设置了代理，可以关闭GOSUM验证
   ```bash
   go env -w GOPROXY=https://goproxy.cn,direct
   go env -w GOSUMDB=off
   ```
   这样会禁用模块哈希数据库验证，适用于可信网络环境或开发调试。
3. **使用**:
    - 运行应用程序: `./k9fyne`
    - 通过配置必要的凭据连接到您的 Kubernetes 集群。
    - 通过菜单选项导航来管理不同的 Kubernetes 资源。

---

#### 作者

- Andrew Wang

非常感谢我们众多慷慨的赞助者。

---

#### 文档与赞助

- [文档](#)
- [赞助](#)

