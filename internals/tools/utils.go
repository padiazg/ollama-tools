package tools

import "fmt"

func QuantizationBits(quantization_level string) int {
	switch quantization_level {
	case "Q4": // 4-bit quantization ≈ 0.5 bytes per parameter
		return 4
	case "Q5": // 5-bit quantization ≈ 0.625 bytes per parameter
		return 5
	case "Q8": // 8-bit quantization = 1 byte per parameter
		return 8
	case "F16", "BF16": // 16-bit floating point = 2 bytes per parameter
		return 16
	case "F32": // 32-bit floating point = 4 bytes per parameter
		return 32
	default:
		// For GGUF/GGML models with unspecified quantization, default to 1.5 bytes average
		if quantization_level != "" {
			fmt.Printf("QuantizationBits %s not recognized, using defaults", quantization_level)
		}
		return 12
	}
}

func BytesPerParameter(quantization_bits int) float64 {
	switch quantization_bits {
	case 4: // 4-bit quantization ≈ 0.5 bytes per parameter
		return 4.0 / 8.0
	case 5: // 5-bit quantization ≈ 0.625 bytes per parameter
		return 5.0 / 8.0
	case 8: // 8-bit quantization = 1 byte per parameter
		return 1.0
	case 16: // 16-bit floating point = 2 bytes per parameter
		return 2.0
	case 32: // 32-bit floating point = 4 bytes per parameter
		return 4.0
	default:
		// For GGUF/GGML models with unspecified quantization, default to 1.5 bytes average
		fmt.Printf("BytesPerParameter %d not recognized, using defaults", quantization_bits)
		return 1.5
	}
}

func SystemRAMMultiplier(quantization_bits int) float64 {
	switch quantization_bits {
	case 4:
		return 1.1 // INT4 most efficient
	case 5:
		return 1.15
	case 8:
		return 1.0 // INT8 more efficient
	case 16:
		return 2.0 // FP16 baseline
	case 32:
		return 4.0 // FP32 needs more headroom
	default:
		fmt.Printf("SystemRAMMultiplier %d not recognized, using defaults", quantization_bits)
		return 1.5 // For GGUF/GGML models with unspecified quantization, default to 1.5 bytes average
	}
}

// NormalizeQuantizationLevel returns a normalized quantization level
func NormalizeQuantizationLevel(quantization_level string) string {
	switch quantization_level[0] {
	case 'Q': // ex: Q4_K_M
		return quantization_level[0:2]
	case 'F': // ex: F16, F32
		return quantization_level[0:3]
	case 'B': // ex: BF16
		return quantization_level[1:4]
	default:
		fmt.Printf("NormalizeQuantizationLevel prefix not recognized: %s", quantization_level)
		return ""
	}
}

func FormatParamCount(parameter_count int64) string {
	switch {
	case parameter_count >= 1_000_000_000:
		return fmt.Sprintf("%.2fB", float64(parameter_count)/1_000_000_000)
	case parameter_count >= 1_000_000:
		return fmt.Sprintf("%.2fM", float64(parameter_count)/1_000_000)
	case parameter_count >= 1_000:
		return fmt.Sprintf("%.2fK", float64(parameter_count)/1_000)
	default:
		return fmt.Sprintf("%d", parameter_count)
	}
}

// FormatMemorySize format memory size to a human-readable string
func FormatMemorySize(memoryMB float64) string {
	if memoryMB >= 1024 {
		return fmt.Sprintf("%.2f GB", memoryMB/1024)
	}
	return fmt.Sprintf("%.2f MB", memoryMB)
}
