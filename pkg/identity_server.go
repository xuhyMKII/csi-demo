package pkg

import (
	"context"
	"log"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/golang/protobuf/ptypes/wrappers"
)

//这些方法是实现 CSI IdentityServer 接口的部分。IdentityServer 接口提供了获取插件信息、
//获取插件能力和探测插件健康状态等方法，这些方法主要供 Kubernetes 或其他 CSI 插件调用，以了解插件的基本信息和状态。

// GetPluginInfo 方法返回此 CSI 插件的名称和版本信息
func (nfs *NFSDriver) GetPluginInfo(context.Context, *csi.GetPluginInfoRequest) (*csi.GetPluginInfoResponse, error) {
	log.Println("GetPluginInfo request")

	// 返回插件的名称和版本信息
	return &csi.GetPluginInfoResponse{
		Name:          nfs.name,
		VendorVersion: nfs.version,
	}, nil
}

// GetPluginCapabilities 方法返回此插件所支持的能力
func (nfs *NFSDriver) GetPluginCapabilities(context.Context, *csi.GetPluginCapabilitiesRequest) (*csi.GetPluginCapabilitiesResponse, error) {
	log.Println("GetPluginCapabilities request")

	//CSI，即 Container Storage Interface，它是一个标准的接口定义，目标是为了提供一种通用的方式以使存储系统能够对接到各种容器编排系统中。
	//
	//在 CSI 规范中，定义了三类服务，分别是：
	//
	//Identity Service: 提供插件的一些基本信息，例如插件的名称、版本以及所支持的服务和特性等。
	//Controller Service: 包含一些控制器相关的操作，例如创建和删除卷、创建和删除快照等。
	//Node Service: 包含一些节点相关的操作，例如在一个特定的节点上挂载或者卸载卷等。
	//
	//csi.GetPluginCapabilitiesResponse 是用来响应 GetPluginCapabilities 请求的。这个响应包含一个 Capabilities 列表，每一项都是一个 PluginCapability 结构体，描述插件的一种能力。
	//
	//每个 PluginCapability 结构体都有一个 Type 字段，表示这个能力的类型。在这里，Type 字段的值是 PluginCapability_Service_，表示这是一个服务类型的能力。
	//
	//PluginCapability_Service_ 结构体里面有一个 Service 字段，这是一个 PluginCapability_Service 结构体，用来描述服务类型的能力的具体信息。
	//
	//在 PluginCapability_Service 结构体里，有一个 Type 字段，用来表示服务的类型。在这里，Type 字段的值是 PluginCapability_Service_CONTROLLER_SERVICE，表示这个插件支持控制器服务（Controller Service）。
	//
	//所以，这段代码的意思是：这个 NFS CSI 插件支持控制器服务。
	// 此插件支持控制器服务，对应的是 CreateVolume 和 DeleteVolume 方法
	return &csi.GetPluginCapabilitiesResponse{
		Capabilities: []*csi.PluginCapability{
			{
				Type: &csi.PluginCapability_Service_{
					Service: &csi.PluginCapability_Service{
						Type: csi.PluginCapability_Service_CONTROLLER_SERVICE,
					},
				},
			},
		},
	}, nil
}

// Probe 方法用于检查插件是否准备就绪
func (nfs *NFSDriver) Probe(context.Context, *csi.ProbeRequest) (*csi.ProbeResponse, error) {
	log.Println("Probe request")

	// 此插件总是准备就绪
	return &csi.ProbeResponse{
		Ready: &wrappers.BoolValue{
			Value: true,
		},
	}, nil
}
