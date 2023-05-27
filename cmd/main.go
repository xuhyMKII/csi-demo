package main

import (
	"flag"
	"log"
	"net"
	"os"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"google.golang.org/grpc"

	"csi-demo/pkg"
)

// 定义了一些全局的命令行参数。包括 CSI 的 unix socket 路径，驱动名，节点 ID，版本号，以及 NFS 服务器的相关信息。
var (
	endpoint   = flag.String("endpoint", "/csi/csi.sock", "CSI unix socket")
	driverName = flag.String("driver-name", "csi-demo", "CSI driver name")
	nodeId     = flag.String("nodeid", "nfs-csi-node", "node id of CSI")
	version    = flag.String("version", "N/A", "version of CSI")
	server     = flag.String("server", "", "nfs server ip")
	serverPath = flag.String("serverPath", "", "nfs server root mount path")
	mountPath  = flag.String("mountPath", "/mount", "local mount path")
)

// 程序的主入口点
func main() {
	// 解析命令行参数
	flag.Parse()

	// 打印一些初始化的信息
	log.Printf("driverName: %s, version: %s, nodeID: %s", *driverName, *version, *nodeId)

	// 如果存在，移除 unix socket 文件，防止启动 gRPC server 时因为文件已存在而失败
	if err := os.Remove(*endpoint); err != nil && !os.IsNotExist(err) {
		log.Fatalf("remove endpoint %s error", *endpoint)
	}

	// 创建一个 unix socket 服务端
	ln, err := net.Listen("unix", *endpoint)
	if err != nil {
		log.Fatalf("listen unix endpoint %s error: %s", *endpoint, err.Error())
	}

	// 当程序退出时关闭 unix socket 服务端，并且移除 unix socket 文件
	defer func() {
		_ = ln.Close()
		_ = os.Remove(*endpoint)
	}()

	// 创建一个 gRPC server
	grpcServer := grpc.NewServer()

	// 创建一个 NFS 驱动
	nfsDriver := pkg.NewNFSDriver(&pkg.Options{
		Name:         *driverName,
		Version:      *version,
		NodeID:       *nodeId,
		NFSServer:    *server,
		NFSRootPath:  *serverPath,
		NFSMountPath: *mountPath,
	})

	// 在 gRPC server 上注册 CSI 相关的服务
	csi.RegisterIdentityServer(grpcServer, nfsDriver)
	csi.RegisterControllerServer(grpcServer, nfsDriver)
	csi.RegisterNodeServer(grpcServer, nfsDriver)

	// 启动 gRPC server
	log.Println("grpc server start")
	defer log.Println("grpc server exit")

	// 开始接收来自客户端的请求
	if err = grpcServer.Serve(ln); err != nil {
		log.Fatalf("grpc serve error: %s", err.Error())
	}
}
