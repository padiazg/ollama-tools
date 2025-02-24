package tools

import (
	"fmt"
	"math"

	"github.com/padiazg/ollama-tools/models/ollama"
)

const (
	ONE_GB = 1_073_741_824 // 1024 * 1024 * 1024
)

func EstimateMemory(parameter_count int64, context_length int, quantization_level string) *ollama.MemoryEstimation {
	var (
		mem                   = &ollama.MemoryEstimation{}
		hiddenSize            = math.Sqrt(float64(parameter_count) / 6)
		quantization_bits     = QuantizationBits(NormalizeQuantizationLevel(quantization_level))
		bytes_per_parameter   = BytesPerParameter(quantization_bits)
		system_ram_multiplier = SystemRAMMultiplier(quantization_bits)
	)

	mem.BaseModelSize = (float64(parameter_count) * bytes_per_parameter) / ONE_GB
	mem.KVCacheSize = (4 * hiddenSize * float64(context_length) * bytes_per_parameter) / ONE_GB
	gpuOverhead := mem.BaseModelSize * .1
	mem.GPURAM = mem.BaseModelSize + mem.KVCacheSize + gpuOverhead
	mem.SystemRAM = mem.GPURAM * system_ram_multiplier

	return mem
}

func PrintEstimatedMemoryPlain(mem *ollama.MemoryEstimation) {
	fmt.Printf("\n  Memory Breakdown:\n")
	fmt.Printf("    Model Weights Memory: %s\n", FormatMemorySize(mem.BaseModelSize))
	fmt.Printf("    KV Cache (for context): %s\n", FormatMemorySize(mem.KVCacheSize))
	fmt.Printf("    GPU VRAM: %s\n", FormatMemorySize(mem.GPURAM))
	fmt.Printf("    System RAM: %s\n", FormatMemorySize(mem.SystemRAM))
}
