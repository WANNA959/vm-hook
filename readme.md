# kube-vm label webhook
## 流程

```shell
# 1. 创建csr.conf，创建key+cert
# 2. 创建CertificateSigningRequest
# 3. apiserver approve csr
# 4. apiserver创建secret(webhook-demo-secret)
# 5. 存储key-cert到本地
sh scripts/webhook_create_apiserver_certs.sh 
```

- 部署webhook

```shell
# delete & rebuild & re-deploy
sh scripts/rebuild_webhook.sh 
```

- 测试

```shell
# 没有host label，有其他label
kubectl apply -f test/resources/create-vm.json

# 有host label，value!="vm"
kubectl apply -f test/resources/create-vm-2.json
```

## mutating逻辑

#### mutating webhook

- 作用域范围
  - apiVersion: doslab.io/v1
  - Kind: VirtualMachine
  
- /mutate/add-label：采用patch的方式
  - 检查是否存在key=host的label
    - 无则add
    - 有且value!="vm"，则update
