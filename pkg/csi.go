package pkg

import (
	"github.com/container-storage-interface/spec/lib/go/csi"
	mount "k8s.io/mount-utils"
)

//这部分代码的主要功能是创建和初始化一个 NFSDriver 实例，
//该实例实现了 CSI 中的三个接口，并添加了一些基本的服务能力。
//这使得 NFSDriver 实例可以作为一个 CSI 插件提供服务。

// Options 是传递给 NewNFSDriver 函数的参数类型
type Options struct {
	Name    string
	Version string
	NodeID  string

	NFSServer    string
	NFSRootPath  string
	NFSMountPath string
}

// NFSDriver 是实现 CSI 插件的主要结构体
type NFSDriver struct {
	name    string
	version string
	nodeID  string

	nfsServer    string
	nfsRootPath  string
	nfsMountPath string

	mounter mount.Interface // 用于执行挂载操作的接口

	// 这两个字段用于存储此 CSI 插件支持的控制器和节点服务能力
	controllerServiceCapabilities []*csi.ControllerServiceCapability
	nodeServiceCapabilities       []*csi.NodeServiceCapability
}

// 确保 NFSDriver 实现了 CSI 中定义的三个接口
var _ csi.IdentityServer = &NFSDriver{}
var _ csi.ControllerServer = &NFSDriver{}
var _ csi.NodeServer = &NFSDriver{}

// NewNFSDriver 函数根据传入的 Options 创建一个 NFSDriver 实例
func NewNFSDriver(opt *Options) *NFSDriver {
	nfs := &NFSDriver{
		name:         opt.Name,
		version:      opt.Version,
		nodeID:       opt.NodeID,
		nfsServer:    opt.NFSServer,
		nfsRootPath:  opt.NFSRootPath,
		nfsMountPath: opt.NFSMountPath,
		mounter:      mount.New(""), // 使用 mount 包创建一个新的挂载器
	}

	// 为 NFSDriver 实例添加支持的控制器服务能力，例如创建和删除卷
	nfs.addControllerServiceCapabilities([]csi.ControllerServiceCapability_RPC_Type{
		csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
	})

	// 此处添加支持的节点服务能力，当前没有添加任何能力
	nfs.addNodeServiceCapabilities([]csi.NodeServiceCapability_RPC_Type{})

	return nfs
}

// newControllerServiceCapability 函数创建一个新的 ControllerServiceCapability
func newControllerServiceCapability(cap csi.ControllerServiceCapability_RPC_Type) *csi.ControllerServiceCapability {
	return &csi.ControllerServiceCapability{
		Type: &csi.ControllerServiceCapability_Rpc{
			Rpc: &csi.ControllerServiceCapability_RPC{
				Type: cap,
			},
		},
	}
}

// addControllerServiceCapabilities 方法为 NFSDriver 实例添加一组控制器服务能力
func (nfs *NFSDriver) addControllerServiceCapabilities(capabilities []csi.ControllerServiceCapability_RPC_Type) {
	var csc = make([]*csi.ControllerServiceCapability, 0, len(capabilities))
	for _, c := range capabilities {
		csc = append(csc, newControllerServiceCapability(c))
	}
	nfs.controllerServiceCapabilities = csc
}

// newNodeServiceCapability 函数创建一个新的 NodeServiceCapability
func newNodeServiceCapability(cap csi.NodeServiceCapability_RPC_Type) *csi.NodeServiceCapability {
	return &csi.NodeServiceCapability{
		Type: &csi.NodeServiceCapability_Rpc{
			Rpc: &csi.NodeServiceCapability_RPC{
				Type: cap,
			},
		},
	}
}

// addNodeServiceCapabilities 函数：这个方法是为 NFSDriver 实例添加一组节点服务能力。
// 它首先创建一个空的 NodeServiceCapability 切片，然后遍历传入的能力类型列表，对于每个能力类型，
// 都调用 newNodeServiceCapability 函数创建一个新的 NodeServiceCapability 并添加到切片中。
// 最后，将这个切片赋值给 NFSDriver 实例的 nodeServiceCapabilities 字段。
func (nfs *NFSDriver) addNodeServiceCapabilities(capabilities []csi.NodeServiceCapability_RPC_Type) {
	var nsc = make([]*csi.NodeServiceCapability, 0, len(capabilities))
	for _, n := range capabilities {
		nsc = append(nsc, newNodeServiceCapability(n))
	}
	nfs.nodeServiceCapabilities = nsc
}
