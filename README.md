# consol-1d

consol-1d 是饱和土一维固结核算（太沙基，∂u/∂t = cv·∂²u/∂z²）的命令行工具：给定固结系数 cv、层厚 H、排水条件（单面/双面）与超静孔压初值 u0（均匀或线性），沿时间和深度输出孔压消散剖面、平均固结度 U 与沉降比 s/s_ult。输入来自一个 JSON 场景文件，`profile` 子命令沿深度打印剖面表并报告层平均固结度、中点孔压与沉降。

- 输入：JSON 场景（`cv`、`thickness`、`drainage`、`initial_pressure`、`time`，可选 `mv`/`delta_sigma` 算 s_ult），示例见 `example/double-drain.json`
- 输出：深度剖面（u、u/u0、消散度 1−u/u0）、平均固结度 U、中点孔压、沉降 s 与 s_ult（可算时）、沉降比 s/s_ult
- 边界：cv>0、H>0、t≥0（t=0 合法：U=0 且孔压=初值）、初值非负且均值>0；NaN/Inf 一律拒绝；非法参数以 error 返回，CLI 打印 stderr 并非零退出
- 不做：不建施工监测看板、不做摩尔圆破坏包线、不做地质编录台账，只做固结内核

## 钉死的约定

- **排水路径**：单面排水排水距离 Hdr = H（z=0 排水、z=H 不排水）；双面排水 Hdr = H/2（z=0 与 z=H 均排水）。
- **时间因子**：Tv = cv·t/Hdr²。全仓只有一个 Tv 来源：Hdr、Tv 与级数指数 exp(−(2n+1)²π²Tv/4) 同源，深度剖面与平均固结度不许各用一套 H。
- **平均固结度（均匀初值）**：双面与单面用同一个经典级数，Tv 按各自 Hdr 定义代入：

  U = 1 − (8/π²)·Σ_{n≥0} exp(−(2n+1)²·π²·Tv/4) / (2n+1)²

  求解器并不把上式当黑盒：U 由层平均孔压定义 U = 1 − ū(t)/ū(0) 得出，均匀初值下与上式恒等（测试交叉验证）。
- **孔压级数（均匀初值）**：双面排水 u(z,t) = (4u0/π)·Σ_{n≥0} sin((2n+1)πz/H)·exp(−(2n+1)²π²Tv/4)/(2n+1)；单面排水特征函数为 sin((2n+1)πz/(2H))，Tv 用 Hdr=H。
- **线性初值**：u(z,0) = uA + (uB−uA)·z/H。模态系数为解析闭式（双面 A_n = 2(uA+uB)/((2n+1)π)；单面 A_n = 4uA/((2n+1)π) + (−1)ⁿ·8(uB−uA)/((2n+1)²π²)），均匀初值是其 uA=uB 特例。U 仍由层平均孔压积分定义，与孔压剖面共用同一套系数与 Tv。
- **沉降**：s_ult = mv·Δσ·H，s = U·s_ult，沉降比 s/s_ult = U。`mv` 与 `delta_sigma` 必须同时给出或同时缺省；缺省时绝对沉降不输出（n/a），沉降比仍按定义等于 U。
- **级数截断**：项数自适应增长，直到绝对余项上界 ≤ 1e-9（用互补误差函数 erfc 做尾和上界）；报告打印实际项数与余项上界，不声称超过该界之外的精度。项数上限 2²⁵，超出报 error。
- **交叉规则**：Tv=0 时 U=0 且孔压=初值；Tv 很大时 U→1 且孔压→0；cv 加倍等效 Tv 加倍（同 H、同 t、同排水）；双面排水比同厚单面快；层中面孔压消散慢于排水面（排水面孔压恒 0）。

## 可测契约

`example/double-drain.json`（H=10 m、cv=1e-7 m²/s、双面排水、均匀 u0=100 kPa、t=5e7 s ⇒ Tv=0.2）的平均固结度 U 落在教科书 0.5 附近（0.5±0.05），测试断言 [0.45, 0.55]。

## 构建 / 运行 / 测试

```text
go build ./...                     # 编译（纯标准库）
go test ./...                      # 全部测试（series / model / report / cli）
go run . profile example/double-drain.json
go run . profile example/double-drain.json --t 1e6      # 覆盖时间
go run . profile example/double-drain.json --nodes 21   # 剖面点数
go run . profile example/double-drain.json -out profile.csv   # 导出剖面 CSV
go run . profile example/double-drain.json -json        # 结构化输出
go run . curve example/double-drain.json                # U-t 固结曲线
go run . curve example/double-drain.json -times 1e5,1e6,1e7
go run . settle example/double-drain.json -target 0.9   # 达到 90% 固结的时间
```
