# consol-1d — Go 太沙基一维固结核算命令行与 HTTP API 服务

consol-1d 是饱和土一维固结核算（太沙基，∂u/∂t = cv·∂²u/∂z²）的 Go 服务：给定固结系数 cv、层厚 H、排水条件（单面/双面）与超静孔压初值 u0（均匀或线性），沿时间和深度输出孔压消散剖面、平均固结度 U 与沉降比 s/s_ult。

## 构建 / 运行 / 测试

```text
go build ./...     # 编译
go run . profile example/double-drain.json
go test ./...      # 测试
```

## 评测镜像

本目录评测专用文件（勿覆盖项目自带 Dockerfile/README）：

- `benzhi.Dockerfile`
- `build_benzhi_docker.sh`
- `BENZHI_README.md`（本文件）

两种架构都要构建并进容器验证：

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh <image-name> linux/arm64
./build_benzhi_docker.sh <image-name> linux/amd64
docker run -it <image-name>:latest
```
