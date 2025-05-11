# ddddocr\_go

一个基于 Go、ONNX Runtime 和 OpenCV 的 OCR 工具，用于快速批量识别图片中的文字。

## 项目引用

* [https://github.com/FeilongTest/go-ddddocr/tree/master](https://github.com/FeilongTest/go-ddddocr/tree/master)
* [https://github.com/sml2h3/ddddocr/tree/master](https://github.com/sml2h3/ddddocr/tree/master)

## 环境依赖

1. **CUDA 12.9**

   * 从 NVIDIA 官网下载并安装：`https://developer.nvidia.com/cuda-12-9-0-downloads`
   * 按照安装程序提示完成安装。

2. **cuDNN v9**

   * 从 NVIDIA 官网下载：`https://developer.nvidia.com/rdp/cudnn-download`
   * 解压后将所有 `cudnn*.dll` 文件复制到 CUDA 安装目录下的 `bin` 文件夹，例如：

     ```text
     C:\Program Files\NVIDIA GPU Computing Toolkit\CUDA\v12.9\bin
     ```

3. **Go (>= 1.20)**

   * 安装 Go 开发环境：`https://golang.org/dl/`

4. **OpenCV (4.x) + Go 绑定**

   * 在 Mingw64 环境中安装 OpenCV，并确保头文件和库文件位于 `/mingw64` 路径。
   * 安装 Go 绑定：

     ```bash
     go get -u gocv.io/x/gocv
     ```

5. **Mingw64**

   * Windows 下使用 Mingw64 终端进行编译。

## 编译步骤

在 Mingw64 终端中执行以下命令：

```bash
# 1. 设置 pkg-config 路径
export PKG_CONFIG_PATH=/mingw64/lib/pkgconfig

# 2. 设置 CGO 编译标志，链接 OpenCV
export CGO_CPPFLAGS="$(pkg-config --cflags opencv4)"
export CGO_LDFLAGS="$(pkg-config --libs opencv4)"

# 3. 进入项目目录
cd /d/go/ddddocr_go/

# 4. 执行编译，生成可执行文件
go build -tags=customenv -o /z/onnxruntime-win-x64-gpu-1.22.0/onnxruntime-win-x64-gpu-1.22.0/lib/ddddocr_go.exe
```

编译完成后，请将以下动态库与可执行文件放在同一目录下：

* `onnxruntime.dll`
* `onnxruntime_providers_cuda.dll`
* `onnxruntime_providers_shared.dll`
* 各种 `cudnn*.dll`

## 使用示例

```bash
# 将需识别的图片（支持 .png、.jpg、.jpeg）放在可执行文件同目录
# 运行程序
ddddocr_go.exe
```

程序会遍历同目录下所有图片并输出识别结果。

---

