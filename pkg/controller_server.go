package pkg

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateVolume 方法是用来创建一个新的卷的。这里的实现是在 NFS 服务器上创建一个新的目录。
// 这个目录的名字由 CSI 发送的请求参数决定。
func (nfs *NFSDriver) CreateVolume(_ context.Context, req *csi.CreateVolumeRequest) (*csi.CreateVolumeResponse, error) {
	// 打印一条日志，表示收到了 CreateVolume 请求
	log.Println("CreateVolume request")

	// 打印要创建的卷的名字
	log.Println("req name: ", req.GetName())
	// 构造 NFS 上的目录路径
	mountPath := filepath.Join(nfs.nfsMountPath, req.GetName())
	// 创建目录
	if err := os.Mkdir(mountPath, 0755); err != nil {
		log.Printf("mkdir %s error: %s", mountPath, err.Error())
		// 如果创建目录失败，返回错误
		return nil, errors.Wrap(err, "mkdir error")
	}

	// 返回 CreateVolumeResponse。VolumeId 是卷的名字，CapacityBytes 是卷的容量。
	// 这里 CapacityBytes 是 0，表示没有限制。
	return &csi.CreateVolumeResponse{
		Volume: &csi.Volume{
			VolumeId:      req.Name,
			CapacityBytes: 0,
		},
	}, nil
}

// DeleteVolume 方法是用来删除一个卷的。这里的实现是删除 NFS 服务器上的一个目录。
func (nfs *NFSDriver) DeleteVolume(_ context.Context, req *csi.DeleteVolumeRequest) (*csi.DeleteVolumeResponse, error) {
	// 打印一条日志，表示收到了 DeleteVolume 请求
	log.Println("DeleteVolume request")

	// 打印要删除的卷的 ID
	log.Println("volumeID: ", req.GetVolumeId())

	// 删除 NFS 上的目录
	return nil, os.Remove(filepath.Join(nfs.nfsMountPath, req.GetVolumeId()))
}

// 下面的方法都还没有实现，所以都返回了一个 Unimplemented 的错误。
// ControllerPublishVolume 方法是用来将一个卷附加到一个节点上的。
// ControllerUnpublishVolume 方法是用来将一个卷从一个节点上卸载的。
// ValidateVolumeCapabilities 方法是用来验证一个卷是否支持指定的访问模式，访问类型，和挂载选项的。
// ListVolumes 方法是用来列出所有的卷的。
// GetCapacity 方法是用来获取存储池的容量的。
// ControllerGetCapabilities 方法是用来获取控制器的功能的。
// CreateSnapshot 方法是用来创建一个卷的快照的。
// DeleteSnapshot 方法是用来删除一个卷的快照的。
// ListSnapshots 方法是用来列出所有的卷快照的。
// ControllerExpandVolume 方法是用来扩容一个卷的。
// ControllerGetVolume 方法是用来获取一个卷的详细信息的。
func (nfs *NFSDriver) ControllerPublishVolume(context.Context, *csi.ControllerPublishVolumeRequest) (*csi.ControllerPublishVolumeResponse, error) {
	log.Println("ControllerPublishVolume request")

	return nil, status.Error(codes.Unimplemented, "")
}

func (nfs *NFSDriver) ControllerUnpublishVolume(context.Context, *csi.ControllerUnpublishVolumeRequest) (*csi.ControllerUnpublishVolumeResponse, error) {
	log.Println("ControllerUnpublishVolume request")

	return nil, status.Error(codes.Unimplemented, "")
}

func (nfs *NFSDriver) ValidateVolumeCapabilities(context.Context, *csi.ValidateVolumeCapabilitiesRequest) (*csi.ValidateVolumeCapabilitiesResponse, error) {
	log.Println("ValidateVolumeCapabilities request")

	return nil, status.Error(codes.Unimplemented, "")
}

func (nfs *NFSDriver) ListVolumes(context.Context, *csi.ListVolumesRequest) (*csi.ListVolumesResponse, error) {
	log.Println("ListVolumes request")

	return nil, status.Error(codes.Unimplemented, "")
}

func (nfs *NFSDriver) GetCapacity(context.Context, *csi.GetCapacityRequest) (*csi.GetCapacityResponse, error) {
	log.Println("GetCapacity request")

	return nil, status.Error(codes.Unimplemented, "")
}

func (nfs *NFSDriver) ControllerGetCapabilities(context.Context, *csi.ControllerGetCapabilitiesRequest) (*csi.ControllerGetCapabilitiesResponse, error) {
	log.Println("ControllerGetCapabilities request")

	return &csi.ControllerGetCapabilitiesResponse{
		Capabilities: nfs.controllerServiceCapabilities,
	}, nil
}

func (nfs *NFSDriver) CreateSnapshot(context.Context, *csi.CreateSnapshotRequest) (*csi.CreateSnapshotResponse, error) {
	log.Println("CreateSnapshot request")

	return nil, status.Error(codes.Unimplemented, "")
}

func (nfs *NFSDriver) DeleteSnapshot(context.Context, *csi.DeleteSnapshotRequest) (*csi.DeleteSnapshotResponse, error) {
	log.Println("DeleteSnapshot request")

	return nil, status.Error(codes.Unimplemented, "")
}

func (nfs *NFSDriver) ListSnapshots(context.Context, *csi.ListSnapshotsRequest) (*csi.ListSnapshotsResponse, error) {
	log.Println("ListSnapshots request")

	return nil, status.Error(codes.Unimplemented, "")
}

func (nfs *NFSDriver) ControllerExpandVolume(context.Context, *csi.ControllerExpandVolumeRequest) (*csi.ControllerExpandVolumeResponse, error) {
	log.Println("ControllerExpandVolume request")

	return nil, status.Error(codes.Unimplemented, "")
}

func (nfs *NFSDriver) ControllerGetVolume(context.Context, *csi.ControllerGetVolumeRequest) (*csi.ControllerGetVolumeResponse, error) {
	log.Println("ControllerGetVolume request")

	return nil, status.Error(codes.Unimplemented, "")
}
