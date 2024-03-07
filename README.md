# K8 Cluster service

## Description

 This is a Cluster service written with gin framework and gorm ORM.
 It is a simple service that can be used to create, read, update and delete 
        1. Deployments
        2. Services
        3. Pods
        4. Namespaces
        5. Nodes
        6. Ingresses
        7. ConfigMaps
        8. Secrets
# DistSysProj_ClusterService
# DistSysProj_ClusterService



<!-- go installation in master commmands -->
curl -OL https://go.dev/dl/go1.20.2.linux-amd64.tar.gz
sha256sum go1.20.3.linux-386.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.20.2.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
source ~/.profile


kubeadm join 172.31.24.23:6443 --token cm17g7.rkdge59x78pr4qbz \
        --discovery-token-ca-cert-hash sha256:b77d9eb39c7e02b4d68dbd3313606d9230066bc046fce2343a45720a16341826





steps to cover now:

1. make the static consensus when a call happens to worker 
from the api the worker 
initiates the consensus process and 
waits for the result 
and if the number of results are more than or equal to 3 we stop

2. deploy a add set with 1 worker and 4 replicas
3. how to access the pod env and details from the go code
4. create a service which acts as the entry point to the worker pod

the parser will be calling the load balancer
the load balancer passes the request to worker and the worker triggers 
the primary replica and the primary replica starts the consensus 
and the primary replica gets the result

if the primary gets 3 result from the replicas serves 
the request otherwise 
waits 
till the response is coming
then worker sends back the response to scheduler







&{&Pod{ObjectMeta:{
    add-worker  default  17cf7858-0a43-4218-9470-12c00355b280 1251904 0 2023-05-13 04:24:20 -0400 
    EDT <nil> <nil> map[] map[] [] [] 
    [{main Update v1 2023-05-13 04:24:20 -0400 EDT FieldsV1 {"f:spec":{"f:containers":{"k:{\"name\":\"add-worker\"}":{".":{},"f:image":{},"f:imagePullPolicy":{},"f:name":{},"f:ports":{".":{},"k:{\"containerPort\":50050,\"protocol\":\"TCP\"}":{".":{},"f:containerPort":{},"f:protocol":{}}},"f:resources":{},"f:terminationMessagePath":{},"f:terminationMessagePolicy":{}}},"f:dnsPolicy":{},"f:enableServiceLinks":{},"f:restartPolicy":{},"f:schedulerName":{},"f:securityContext":{},"f:terminationGracePeriodSeconds":{}}} }]},Spec:PodSpec{Volumes:[]Volume{Volume{Name:kube-api-access-xhsc9,VolumeSource:VolumeSource{HostPath:nil,EmptyDir:nil,GCEPersistentDisk:nil,AWSElasticBlockStore:nil,GitRepo:nil,Secret:nil,NFS:nil,ISCSI:nil,Glusterfs:nil,PersistentVolumeClaim:nil,RBD:nil,FlexVolume:nil,Cinder:nil,CephFS:nil,Flocker:nil,DownwardAPI:nil,FC:nil,AzureFile:nil,ConfigMap:nil,VsphereVolume:nil,Quobyte:nil,AzureDisk:nil,PhotonPersistentDisk:nil,PortworxVolume:nil,ScaleIO:nil,Projected:&ProjectedVolumeSource{Sources:[]VolumeProjection{VolumeProjection{Secret:nil,DownwardAPI:nil,ConfigMap:nil,ServiceAccountToken:&ServiceAccountTokenProjection{Audience:,ExpirationSeconds:*3607,Path:token,},},VolumeProjection{Secret:nil,DownwardAPI:nil,ConfigMap:&ConfigMapProjection{LocalObjectReference:LocalObjectReference{Name:kube-root-ca.crt,},Items:[]KeyToPath{KeyToPath{Key:ca.crt,Path:ca.crt,Mode:nil,},},Optional:nil,},ServiceAccountToken:nil,},VolumeProjection{Secret:nil,DownwardAPI:&DownwardAPIProjection{Items:[]DownwardAPIVolumeFile{DownwardAPIVolumeFile{Path:namespace,FieldRef:&ObjectFieldSelector{APIVersion:v1,FieldPath:metadata.namespace,},ResourceFieldRef:nil,Mode:nil,},},},ConfigMap:nil,ServiceAccountToken:nil,},},DefaultMode:*420,},StorageOS:nil,CSI:nil,Ephemeral:nil,},},},Containers:[]Container{Container{Name:add-worker,Image:rajeshreddyt/grpcserveradd-worker:latest,Command:[],Args:[],WorkingDir:,Ports:[]ContainerPort{ContainerPort{Name:,HostPort:0,ContainerPort:50050,Protocol:TCP,HostIP:,},},Env:[]EnvVar{},Resources:ResourceRequirements{Limits:ResourceList{},Requests:ResourceList{},Claims:[]ResourceClaim{},},VolumeMounts:[]VolumeMount{VolumeMount{Name:kube-api-access-xhsc9,ReadOnly:true,MountPath:/var/run/secrets/kubernetes.io/serviceaccount,SubPath:,MountPropagation:nil,SubPathExpr:,},},LivenessProbe:nil,ReadinessProbe:nil,Lifecycle:nil,TerminationMessagePath:/dev/termination-log,ImagePullPolicy:Always,SecurityContext:nil,Stdin:false,StdinOnce:false,TTY:false,EnvFrom:[]EnvFromSource{},TerminationMessagePolicy:File,VolumeDevices:[]VolumeDevice{},StartupProbe:nil,},},RestartPolicy:Always,TerminationGracePeriodSeconds:*30,ActiveDeadlineSeconds:nil,DNSPolicy:ClusterFirst,NodeSelector:map[string]string{},ServiceAccountName:default,DeprecatedServiceAccount:default,NodeName:,HostNetwork:false,HostPID:false,HostIPC:false,SecurityContext:&PodSecurityContext{SELinuxOptions:nil,RunAsUser:nil,RunAsNonRoot:nil,SupplementalGroups:[],FSGroup:nil,RunAsGroup:nil,Sysctls:[]Sysctl{},WindowsOptions:nil,FSGroupChangePolicy:nil,SeccompProfile:nil,},ImagePullSecrets:[]LocalObjectReference{},Hostname:,Subdomain:,Affinity:nil,SchedulerName:default-scheduler,InitContainers:[]Container{},AutomountServiceAccountToken:nil,Tolerations:[]Toleration{Toleration{Key:node.kubernetes.io/not-ready,Operator:Exists,Value:,Effect:NoExecute,TolerationSeconds:*300,},Toleration{Key:node.kubernetes.io/unreachable,Operator:Exists,Value:,Effect:NoExecute,TolerationSeconds:*300,},},HostAliases:[]HostAlias{},PriorityClassName:,Priority:*0,DNSConfig:nil,ShareProcessNamespace:nil,ReadinessGates:[]PodReadinessGate{},RuntimeClassName:nil,EnableServiceLinks:*true,PreemptionPolicy:*PreemptLowerPriority,Overhead:ResourceList{},TopologySpreadConstraints:[]TopologySpreadConstraint{},EphemeralContainers:[]EphemeralContainer{},SetHostnameAsFQDN:nil,OS:nil,HostUsers:nil,SchedulingGates:[]PodSchedulingGate{},ResourceClaims:[]PodResourceClaim{},},Status:PodStatus{Phase:Pending,Conditions:[]PodCondition{},Message:,Reason:,HostIP:,PodIP:,StartTime:<nil>,ContainerStatuses:[]ContainerStatus{},QOSClass:BestEffort,InitContainerStatuses:[]ContainerStatus{},NominatedNodeName:,PodIPs:[]PodIP{},EphemeralContainerStatuses:[]ContainerStatus{},},}}