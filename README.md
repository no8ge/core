# core

核心控制器

## 要求

- Kubernetes
- Helm

## 快速开始

```shell
# 查看应用
helm show all oci://registry-1.docker.io/no8ge/core --version 1.0.0

# 下载应用
helm pull oci://registry-1.docker.io/no8ge/core --version 1.0.0

# 安装应用
helm install core oci://registry-1.docker.io/no8ge/core --version 1.0.0

# 升级应用
helm upgrade core oci://registry-1.docker.io/no8ge/core --version 1.0.0
```
