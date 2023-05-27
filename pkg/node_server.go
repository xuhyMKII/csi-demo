package pkg

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	mount "k8s.io/mount-utils"
)

//这段代码主要用于实现 CSI 的节点服务接口，包括卷的挂载、卸载、获取卷统计信息、
//扩展卷的大小以及获取节点的信息和能力等操作。其中，挂载和卸载操作已经实现，
//但获取卷统计信息和扩展卷的大小等操作目前还没有实现。

// NodeStageVolume 是 CSI 规范中定义的一个 RPC，用于在节点上准备卷以供容器使用。
// 这个方法目前并没有实现，所以返回了 Unimplemented 错误。
func (nfs *NFSDriver) NodeStageVolume(context.Context, *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse, error) {
	log.Println("NodeStageVolume request")

	return nil, status.Error(codes.Unimplemented, "")
}

// NodeUnstageVolume 是 CSI 规范中定义的一个 RPC，用于在节点上清理卷，以便它可以被卸载。
// 这个方法目前并没有实现，所以返回了 Unimplemented 错误。
func (nfs *NFSDriver) NodeUnstageVolume(context.Context, *csi.NodeUnstageVolumeRequest) (*csi.NodeUnstageVolumeResponse, error) {
	log.Println("NodeUnstageVolume request")

	return nil, status.Error(codes.Unimplemented, "")
}

// NodePublishVolume 是 CSI 规范中定义的一个 RPC，用于在节点上挂载卷。
// 这个方法会检查请求中的参数，然后挂载 NFS 卷到指定的路径。
func (nfs *NFSDriver) NodePublishVolume(_ context.Context, req *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
	log.Println("NodePublishVolume request")

	// 检查卷的容量是否存在
	capacity := req.GetVolumeCapability()
	if capacity == nil {
		return nil, errors.Errorf("capacity is nill")
	}

	// 获取挂载选项
	options := capacity.GetMount().GetMountFlags()
	if req.Readonly {
		options = append(options, "ro")
	}

	// 检查目标路径是否存在
	targetPath := req.GetTargetPath()
	if targetPath == "" {
		return nil, errors.Errorf("target path is nill")
	}

	// 构造 NFS 卷的源路径
	source := fmt.Sprintf("%s:%s", nfs.nfsServer, filepath.Join(nfs.nfsRootPath, req.GetVolumeId()))

	// 检查目标路径是否已经被挂载
	notMnt, err := nfs.mounter.IsLikelyNotMountPoint(targetPath)
	if err != nil {
		if os.IsNotExist(err) {
			// 如果目标路径不存在，就创建它
			if err := os.MkdirAll(targetPath, os.FileMode(0755)); err != nil {
				return nil, status.Error(codes.Internal, err.Error())
			}
			notMnt = true
		} else {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	if !notMnt {
		// 如果目标路径已经被挂载，那么直接返回成功
		return &csi.NodePublishVolumeResponse{}, nil
	}

	// 挂载 NFS 卷到目标路径
	log.Printf("source: %s, targetPath: %s, options: %v", source, targetPath, options)
	if err := nfs.mounter.Mount(source, targetPath, "nfs", options); err != nil {
		return nil, errors.Wrap(err, "mount nfs path error")
	}

	return &csi.NodePublishVolumeResponse{}, nil
}

// NodeUnpublishVolume 是 CSI 规范中定义的一个 RPC，用于在节点上卸载卷。
// 这个方法会卸载请求中指定的卷。
func (nfs *NFSDriver) NodeUnpublishVolume(_ context.Context, req *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	log.Println("NodeUnpublishVolume request")

	// 获取目标路径
	targetPath := req.GetTargetPath()

	// 卸载目标路径上的卷
	if err := mount.CleanupMountPoint(targetPath, nfs.mounter, true); err != nil {
		return nil, errors.Wrap(err, "clean mount point error")
	}

	return &csi.NodeUnpublishVolumeResponse{}, nil
}

// NodeGetVolumeStats 是 CSI 规范中定义的一个 RPC，用于获取卷的使用统计信息。
// 这个方法目前并没有实现，所以返回了 Unimplemented 错误。
func (nfs *NFSDriver) NodeGetVolumeStats(context.Context, *csi.NodeGetVolumeStatsRequest) (*csi.NodeGetVolumeStatsResponse, error) {
	log.Println("NodeGetVolumeStats request")

	return nil, status.Error(codes.Unimplemented, "")
}

// NodeExpandVolume 是 CSI 规范中定义的一个 RPC，用于在节点上扩展卷的大小。
// 这个方法目前并没有实现，所以返回了 Unimplemented 错误。
func (nfs *NFSDriver) NodeExpandVolume(context.Context, *csi.NodeExpandVolumeRequest) (*csi.NodeExpandVolumeResponse, error) {
	log.Println("NodeExpandVolume request")

	return nil, status.Error(codes.Unimplemented, "")
}

// NodeGetCapabilities 是 CSI 规范中定义的一个 RPC，用于获取节点服务的能力。
// 这个方法会返回 NFSDriver 的节点服务能力。
func (nfs *NFSDriver) NodeGetCapabilities(context.Context, *csi.NodeGetCapabilitiesRequest) (*csi.NodeGetCapabilitiesResponse, error) {
	log.Println("NodeGetCapabilities request")

	return &csi.NodeGetCapabilitiesResponse{
		Capabilities: nfs.nodeServiceCapabilities,
	}, nil
}

// NodeGetInfo 是 CSI 规范中定义的一个 RPC，用于获取节点的信息。
// 这个方法会返回 NFSDriver 的节点 ID。
func (nfs *NFSDriver) NodeGetInfo(context.Context, *csi.NodeGetInfoRequest) (*csi.NodeGetInfoResponse, error) {
	log.Println("NodeGetInfo request")

	return &csi.NodeGetInfoResponse{
		NodeId: nfs.nodeID,
	}, nil
}
