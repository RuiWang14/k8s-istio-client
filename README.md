# k8s-istio-client

a go client for operating istio in k8s cluster

## 原理

![image](docs/imgs/customer controller.png)

Informer 有两个作用：
1. 通过 listWatch 机制与 kube api server 同步，更新维护本地缓存；
2. 根据增量事件类型触发 event handler；

每经过 resyncPeriod 事件，Informer 维护的本地缓存都会使用最近一次 LIST 返回的结果强制更新一次，从而保证缓存的有效性。

WorkQueue 作用:
1. 同步 Informer 和 控制循环之间的数据；

control loop 控制循环:
1. 从 WorkQueue 出队一个成员；
2. 从本地缓存拿到该成员对应的对象；
3. 对比“期望状态”和“实际状态”，做相关操作；

- “期望状态”：APIServer 里保存的，用户提交到 APIServer 里的信息（已经缓存在了本地）；
- “实际状态”：集群/业务实际的状态；
