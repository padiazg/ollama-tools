# The theory
In fact it's more my notes about how to come to a number for an estimation on RAM requirements for running a model. All of them were taken from other projects and looking around on several pages.

## Tables
Here are the values used in the calculations based on the quantization level.
$$
\begin{aligned}
quantizationBits &=  \begin{cases} 
  4 & \text{if } quantizationLevel = \text{Q4} \\
  5 & \text{if } quantizationLevel = \text{Q5} \\
  8 & \text{if } quantizationLevel = \text{Q8} \\
  16 & \text{if } quantizationLevel = \text{F16} \\
  32 & \text{if } quantizationLevel = \text{F32} \\
  12 & otherwise
\end{cases} \\
\\
bytesPerParameter &= \begin{cases} 
4.0/8.0 & \text{if } quantizationBits = 4 & \text{4-bit quantization ≈ 0.5 bytes per parameter}\\
0.625 & \text{if } quantizationBits = 5 & \text{5-bit quantization ≈ 0.625 bytes per parameter} \\
1.0 & \text{if } quantizationBits = 8 & \text{8-bit quantization = 1 byte per parameter} \\
2.0 & \text{if } quantizationBits = 16 & \text{16-bit floating point = 2 bytes per parameter} \\
4.0 & \text{if } quantizationBits = 32 & \text{32-bit floating point = 4 bytes per parameter}  \\
1.5 & otherwise & \text{For GGUF/GGML models with unspecified quantization, default to 1.5 bytes average} 
\end{cases} \\
\\
systemRAMMultiplier &= \begin{cases}
1.1 & \text{if } quantizationBits &= 4 & \text{INT4 most efficient} \\
1.15 & \text{if } quantizationBits &= 5 & \\
1.0 & \text{if } quantizationBits &= 8  & \text{INT8 more efficient} \\
2.0 & \text{if } quantizationBits &= 16 & \text{FP16 baseline} \\
4.0 & \text{if } quantizationBits &= 32 & \text{FP32 needs more headroom} \\
1.5 & otherwise
\end{cases} \\
\\
1Gb &= 1073741824
\end{aligned}
$$

## Formulas:
For the GPU VRAM estimation lets sum the base model size, the key-value cache size and the gpu overhead.  
For the total system ram we need to apply a multiplier based on the quantization level.
$$
\begin{aligned}
BaseModelSize &= (ParametersCount * BytesPerParameter) / 1Gb \\ \\
HiddenSize &= \sqrt{ParametersCount/6} \\ \\
KVCacheSize &= (4 * HiddenSize * ContextLength * BytesPerParameter) / 1Gb \\ \\
GPUOverhead &= BaseModelSize * .1 \\ \\
TotalGPURAM &= BaseModelSizeGB + KVCacheSize + GPUOverhead \\ \\
TotalSystemRAM &= TotalGPURAM * SystemRAMMultiplier \\ 
\end{aligned}
$$
