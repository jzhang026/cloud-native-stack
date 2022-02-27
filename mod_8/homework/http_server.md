### pod yaml

```yaml
apiVersion: v1
kind: Pod
metadata:
        name: http-server
spec:
        containers:
                - name: http-server
                  image: localhost:5000/httpserver
                  imagePullPolicy: IfNotPresent
                  resources:
                          limits:
                                  cpu: 700m
                                  memory: "500Mi"
                          requests:
                                  cpu: 700m
                                  memory: "500Mi"
.
                  readinessProbe:
                          httpGet:        
                                  path: /healthz
                                  port: 80
                          initialDelaySeconds: 15
                          periodSeconds: 5
                  lifecycle:
                          preStop:
                                  exec:
                                          command: ["PID=`ps -ef | grep httpser | grep -v sh | grep -v grep | awk '{print $2}'` && kill -9 $PID"]
```

### 1 YAML参数备注

#### 1.1  配置本地镜像

##### imagePullPolicy: IfNotPresent

```
优先读取该节点本地image镜像
```

##### 上传本地http image至docker仓库

##### 配置k8s本地image方法:

```
1. 搭建Docker Registry
This image contains an implementation of the Docker Registry HTTP API V2 for use with Docker 1.6+. See github.com/docker/distribution for more details about what it is.

Run a local registry: Quick Version

$ docker run -d -p 5000:5000 --restart=always --name registry registry

这将使用官方的 registry 镜像来启动私有仓库。默认情况下，仓库会被创建在容器的 /var/lib/registry 目录下。你可以通过 -v 参数来将镜像文件存放在本地的指定路径。例如下面的例子将上传的镜像放到本地的 /opt/data/registry 目录。
$ docker run -d \
    -p 5000:5000 \
    -v /opt/data/registry:/var/lib/registry \
    registry


$ docker tag httpserver localhost:5000/httpserver
$ docker push localhost:5000/httpserver:latest
查看镜像:

curl http://localhost:5000/v2/_catalog
{"repositories":["httpserver"]}
确认有httpserver镜像
```

#### 设置qos

##### 当limits=requests时，pods的qosClass为Guaranteed

```
resources:
    limits:
        cpu: 700m
        memory: "500Mi"
        
    requests:
    
        cpu: 700m
        memory: "500Mi"
```

##### readinessProbe 优雅启动pod+探活，通过httpGet实现

```
                  readinessProbe:
                          httpGet:        
                                  path: /healthz
                                  port: 80
                          initialDelaySeconds: 15
                          periodSeconds: 5
```

##### preStop 优化停止pod，通过执行kill pid实现

```
                  lifecycle:
                          preStop:
                                  exec:
                                          command: ["PID=`ps -ef | grep httpser | grep -v sh | grep -v grep | awk '{print $2}'` && kill -9 $PID"]
```



### 运行验证

#### 启动pod

kubectl create -f http.yaml

```
$ kubectl create -f http.yaml 
pod/http-server created

```

watch查看container在1秒内完成创建，处于Ready状态为0/1，pod起了但是由于优雅启动机制在，因此没有立刻对外提供服务。延迟启动15秒后，探活机制生效，过了5秒httpGet获取到了200请求，Ready状态变为1/1， 对外提供服务

```
$ kubectl get pod -w
NAME                               READY   STATUS    RESTARTS      AGE
envoy-6958c489d9-7jf79             1/1     Running   1 (26h ago)   3d4h
ng-deployment                      1/1     Running   2 (26h ago)   25d
nginx-deployment-8d545c96d-ls7sm   1/1     Running   1 (26h ago)   3d4h
http-server                         0/1     Pending   0             0s
http-server                         0/1     Pending   0             0s
http-server                         0/1     ContainerCreating   0             0s
http-server                         0/1     ContainerCreating   0             1s
http-server                         0/1     Running             0             1s
http-server                         1/1     Running             0             20s

```



> 参考其他同学的优秀作业: https://github.com/ducknightii/geek-cloud/tree/main/module8

