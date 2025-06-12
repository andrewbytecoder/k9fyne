package pullimage

import (
	"context"
	"fmt"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/google/go-containerregistry/pkg/v1/layout"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"log"
	"os"
	"strings"
)

type K9PullImageInfo struct {
	ctx     context.Context
	log     *zap.Logger
	kc      *kubernetes.Clientset
	podList *corev1.PodList
	podsMap map[string]*KubePullImageInfo // namespace to pods
}

func NewK9PullImageInfo(ctx context.Context, k *kubernetes.Clientset, log *zap.Logger) *K9PullImageInfo {
	return &K9PullImageInfo{
		ctx:     ctx,
		log:     log,
		kc:      k,
		podList: nil,
		podsMap: make(map[string]*KubePullImageInfo),
	}
}

type KubePullImageInfo struct {
	podMap map[string]corev1.Pod // Pod name 2 Pod
}

type KubePullImageInfoInterface interface {
	GetAllTags(ref string) []string
	GetPullImageInfoByNamespace(ref string, tag string) error
}

func (p *K9PullImageInfo) GetAllTags(ref string) []string {
	repo, err := name.NewRepository(ref)
	if err != nil {
		log.Fatalf("Error constructing repo name %q: %v", ref, err)
		return []string{}
	}

	tags, err := remote.List(repo)
	if err != nil {
		log.Fatalf("Error listing tags for %q: %v", ref, err)
		return []string{}
	}

	return tags
}

func (p *K9PullImageInfo) GetPullImageInfoByNamespace(imageName string, tag string) error {
	refStr := strings.Join([]string{imageName, tag}, ":")
	// 2. 解析镜像引用（包含 registry、repository 和 tag）
	ref, err := name.ParseReference(refStr)
	if err != nil {
		return err
	}

	// 3. 使用 remote 包从远程获取镜像
	img, err := remote.Image(ref, remote.WithAuth(authn.Anonymous))
	if err != nil {
		return err
	}

	// 4. 创建一个临时布局目录（用于中间存储 OCI layout 格式）
	layoutPath := strings.Join([]string{imageName, "nginx-layout", tag}, "-")
	err = os.RemoveAll(layoutPath) // 清除旧数据（如果存在）
	if err != nil {
		return err
	}

	l, err := layout.Write(layoutPath, empty.Index)
	if err != nil {
		return err
	}

	// 5. 将镜像写入 layout 目录
	err = l.AppendImage(img)
	if err != nil {
		return err
	}

	// 6. 将 layout 转换为 tar 文件（Docker archive 格式）
	tarFileName := strings.Join([]string{imageName, tag}, "-")
	tarFile := tarFileName + ".tar"
	fmt.Println("Saving image as:", tarFile)

	f, err := os.Create(tarFile)
	if err != nil {
		return err
	}
	defer f.Close()

	err = tarball.Write(ref.Context().Tag("latest"), img, f)
	if err != nil {
		return err
	}

	fmt.Println("✅ Image saved successfully.")
	return nil
}
