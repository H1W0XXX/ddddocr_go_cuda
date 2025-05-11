package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"

	"ddddocr_go/charmap"
	"github.com/nfnt/resize"
	ort "github.com/yalue/onnxruntime_go"
)

func main() {
	// 1. 初始化 ONNX Runtime
	os.Setenv("log_severity_level", "0")
	os.Setenv("ORT_LOG_VERBOSITY_LEVEL", "4")
	ort.SetSharedLibraryPath("onnxruntime.dll")
	if err := ort.InitializeEnvironment(); err != nil {
		log.Fatal("InitializeEnvironment:", err)
	}
	defer ort.DestroyEnvironment()

	// 2. SessionOptions + CUDA
	so, err := ort.NewSessionOptions()
	if err != nil {
		log.Fatal("NewSessionOptions:", err)
	}
	cudaOpts, _ := ort.NewCUDAProviderOptions()
	so.AppendExecutionProviderCUDA(cudaOpts)

	// 3. 创建 DynamicAdvancedSession
	dynSess, err := ort.NewDynamicAdvancedSession(
		"common.onnx",
		[]string{"input1"},
		[]string{"output"},
		so,
	)
	if err != nil {
		log.Fatal("NewDynamicAdvancedSession:", err)
	}
	defer dynSess.Destroy()

	// 4. 遍历当前目录下所有 .png 文件
	dir := "." // 或者指定其他目录
	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal("ReadDir:", err)
	}

	const H = 64
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(e.Name()))
		if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
			continue
		}
		imgPath := filepath.Join(dir, e.Name())

		// —— 4.1 打开并 Decode
		f, err := os.Open(imgPath)
		if err != nil {
			log.Printf("[%s] Open error: %v\n", e.Name(), err)
			continue
		}
		img, _, err := image.Decode(f)
		f.Close()
		if err != nil {
			log.Printf("[%s] Decode error: %v\n", e.Name(), err)
			continue
		}

		// —— 4.2 Resize & 灰度归一化
		W := int(float64(img.Bounds().Dx()) * float64(H) / float64(img.Bounds().Dy()))
		resized := resize.Resize(uint(W), uint(H), img, resize.Lanczos2)

		data := make([]float32, H*W)
		for y := 0; y < H; y++ {
			for x := 0; x < W; x++ {
				r, g, b, _ := resized.At(x, y).RGBA()
				data[y*W+x] = float32((r>>8+g>>8+b>>8)/3) / 255.0
			}
		}

		// —— 4.3 构造 Tensor
		inShape := ort.NewShape(1, 1, int64(H), int64(W))
		inTensor, err := ort.NewTensor(inShape, data)
		if err != nil {
			log.Printf("[%s] NewTensor error: %v\n", e.Name(), err)
			continue
		}
		// 每张图用完后销毁
		defer inTensor.Destroy()

		seqLen := int(math.Ceil(float64(W) / 4.0))
		outShape := ort.NewShape(1, int64(seqLen))
		outTensor, err := ort.NewEmptyTensor[int64](outShape)
		if err != nil {
			inTensor.Destroy()
			log.Printf("[%s] NewEmptyTensor error: %v\n", e.Name(), err)
			continue
		}
		defer outTensor.Destroy()

		// —— 4.4 推理
		if err := dynSess.Run(
			[]ort.Value{inTensor},
			[]ort.Value{outTensor},
		); err != nil {
			log.Printf("[%s] Run error: %v\n", e.Name(), err)
			continue
		}

		// —— 4.5 CTC 解码
		outs := outTensor.GetData()
		var sb strings.Builder
		var last int64
		for _, v := range outs {
			if v != 0 && v != last && int(v) < len(charmap.BetaCharset) {
				sb.WriteString(charmap.BetaCharset[int(v)])
			}
			last = v
		}

		fmt.Printf("[%s] => %s\n", e.Name(), sb.String())
	}
}

//export PKG_CONFIG_PATH=/mingw64/lib/pkgconfig
//export CGO_CPPFLAGS="$(pkg-config --cflags opencv4)"
//export CGO_LDFLAGS="$(pkg-config --libs opencv4)"
//cd /d/go/ddddocr_go/
//go build -tags=customenv -o /z/onnxruntime-win-x64-gpu-1.22.0/onnxruntime-win-x64-gpu-1.22.0/lib/ddddocr_go.exe
