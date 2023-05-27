# NFS CSI 插件部署和测试指南

这个项目提供了一个基于NFS的CSI插件，用于在Kubernetes集群中动态供应NFS卷并进行挂载。以下是部署和测试的步骤：

## 前置条件

- 一个运行中的Kubernetes集群（版本1.13以上）
- 配置好的NFS服务器

## 部署步骤

1. 首先，将NFS服务器的IP和根路径填写到`Node.yaml`和`provisioner.yaml`中对应的位置。

2. 部署`Node.yaml`：
   ```
   kubectl apply -f Node.yaml
   ```
   这将在每个节点上部署一个DaemonSet，负责处理卷的挂载和卸载。

3. 部署`provisioner.yaml`：
   ```
   kubectl apply -f provisioner.yaml
   ```
   这将部署一个Deployment，负责处理卷的动态供应。

## 测试步骤

1. 创建一个PVC来测试动态供应：
   ```
   kubectl apply -f - <<EOF
   apiVersion: v1
   kind: PersistentVolumeClaim
   metadata:
     name: test-pvc
   spec:
     accessModes:
       - ReadWriteOnce
     resources:
       requests:
         storage: 1Gi
     storageClassName: csi-demo
   EOF
   ```
   检查PVC的状态，直到它变为Bound：
   ```
   kubectl get pvc test-pvc
   ```

2. 创建一个Pod来测试卷的挂载：
   ```
   kubectl apply -f - <<EOF
   apiVersion: v1
   kind: Pod
   metadata:
     name: test-pod
   spec:
     containers:
     - name: test-container
       image: busybox
       command: ["sh", "-c", "echo 'Hello, World!' > /data/hello; sleep 3600"]
       volumeMounts:
       - name: test-volume
         mountPath: /data
     volumes:
     - name: test-volume
       persistentVolumeClaim:
         claimName: test-pvc
   EOF
   ```
   检查Pod的状态，直到它变为Running：
   ```
   kubectl get pod test-pod
   ```

3. 进入Pod，检查文件是否写入到NFS卷：
   ```
   kubectl exec -it test-pod -- cat /data/hello
   ```
   如果一切正常，你应该能看到输出 "Hello, World!"。

