

https://github.com/ollama/ollama/issues/1388
```
New quant types (recommended):
- Q2_K: smallest, extreme quality loss - not recommended
- Q3_K: alias for Q3_K_M
- Q3_K_S: very small, very high quality loss
- Q3_K_M: very small, very high quality loss
- Q3_K_L: small, substantial quality loss
- Q4_K: alias for Q4_K_M
- Q4_K_S: small, significant quality loss
- Q4_K_M: medium, balanced quality - recommended
- Q5_K: alias for Q5_K_M
- Q5_K_S: large, low quality loss - recommended
- Q5_K_M: large, very low quality loss - recommended
- Q6_K: very large, extremely low quality loss
- Q8_0: very large, extremely low quality loss - not recommended
- F16: extremely large, virtually no quality loss - not recommended
- F32: absolutely huge, lossless - not recommended
```

https://thoughtbot.com/blog/understanding-open-source-llms


https://unfoldai.com/gpu-memory-requirements-for-llms/
Common quantization formats include:
- FP32 (32-bit floating point): 4 bytes per parameter
- FP16 (16-bit floating point): 2 bytes per parameter
- INT8 (8-bit integer): 1 byte per parameter
- INT4 (4-bit integer): 0.5 bytes per parameter


https://github.com/aleibovici/ollama-gpu-calculator/blob/main/src/OllamaGPUCalculator.js

```js
const OllamaGPUCalculator = () => {
    const [parameters, setParameters] = useState('');
    const [quantization, setQuantization] = useState('16');
    const [contextLength, setContextLength] = useState(4096);
    const [gpuConfigs, setGpuConfigs] = useState([{ gpuModel: '', count: '1' }]);
    const [results, setResults] = useState(null);

    useEffect(() => {
        if (parameters && gpuConfigs.some(config => config.gpuModel)) {
            calculateOllamaRAM();
        }
    }, [
        parameters,
        quantization,
        contextLength,
        gpuConfigs,
        // Add any other state variables that should trigger recalculation
    ]);

    const unsortedGpuSpecs = {
        // GPU specifications with TFLOPS values in FP16/mixed precision and TDP in watts
        'h100': { name: 'H100', vram: 80, generation: 'Hopper', tflops: 1979, tdp: 700 },  // Correct: 700W SXM
        'a100-80gb': { name: 'A100 80GB', vram: 80, generation: 'Ampere', tflops: 312, tdp: 400 },  // Correct: 400W SXM
        'a100-40gb': { name: 'A100 40GB', vram: 40, generation: 'Ampere', tflops: 312, tdp: 400 },  // Correct: 400W SXM
        'a40': { name: 'A40', vram: 48, generation: 'Ampere', tflops: 149.8, tdp: 300 },  // Correct: 300W
        'v100-32gb': { name: 'V100 32GB', vram: 32, generation: 'Volta', tflops: 125, tdp: 300 },  // Correct: 300W SXM2
        'v100-16gb': { name: 'V100 16GB', vram: 16, generation: 'Volta', tflops: 125, tdp: 300 },  // Correct: 300W SXM2
        'rtx4090': { name: 'RTX 4090', vram: 24, generation: 'Ada Lovelace', tflops: 82.6, tdp: 450 },  // Correct: 450W
        'rtx4080': { name: 'RTX 4080', vram: 16, generation: 'Ada Lovelace', tflops: 65, tdp: 320 },  // Correct: 320W
        'rtx3090ti': { name: 'RTX 3090 Ti', vram: 24, generation: 'Ampere', tflops: 40, tdp: 450 },  // Correct: 450W
        'rtx3090': { name: 'RTX 3090', vram: 24, generation: 'Ampere', tflops: 35.6, tdp: 350 },  // Correct: 350W
        'rtx3080ti': { name: 'RTX 3080 Ti', vram: 12, generation: 'Ampere', tflops: 34.1, tdp: 350 },  // Correct: 350W
        'rtx3080': { name: 'RTX 3080', vram: 10, generation: 'Ampere', tflops: 29.8, tdp: 320 },  // Correct: 320W
        'a6000': { name: 'A6000', vram: 48, generation: 'Ampere', tflops: 38.7, tdp: 300 },  // Correct: 300W
        'a5000': { name: 'A5000', vram: 24, generation: 'Ampere', tflops: 27.8, tdp: 230 },  // Correct: 230W
        'a4000': { name: 'A4000', vram: 16, generation: 'Ampere', tflops: 19.2, tdp: 140 },  // Correct: 140W
        'rtx4060ti': { name: 'RTX 4060 Ti', vram: 8, generation: 'Ada Lovelace', tflops: 22.1, tdp: 160 },  // Correct: 160W
        'gtx1080ti': { name: 'GTX 1080 Ti', vram: 11, generation: 'Pascal', tflops: 11.3, tdp: 250 },  // Correct: 250W
        'gtx1070ti': { name: 'GTX 1070 Ti', vram: 8, generation: 'Pascal', tflops: 8.1, tdp: 180 },  // Correct: 180W
        'teslap40': { name: 'Tesla P40', vram: 24, generation: 'Pascal', tflops: 12, tdp: 250 },  // Correct: 250W
        'teslap100': { name: 'Tesla P100', vram: 16, generation: 'Pascal', tflops: 9.3, tdp: 250 },  // Correct: 250W PCIe
        'gtx1070': { name: 'GTX 1070', vram: 8, generation: 'Pascal', tflops: 6.5, tdp: 150 },  // Correct: 150W
        'gtx1060': { name: 'GTX 1060', vram: 6, generation: 'Pascal', tflops: 4.4, tdp: 120 },  // Correct: 120W
        'm4': { name: 'Apple M4', vram: 16, generation: 'Apple Silicon', tflops: 4.6, tdp: 30 },  // Estimated: Not released yet
        'm3-max': { name: 'Apple M3 Max', vram: 40, generation: 'Apple Silicon', tflops: 4.5, tdp: 92 },  // Updated: ~92W max package power
        'm3-pro': { name: 'Apple M3 Pro', vram: 18, generation: 'Apple Silicon', tflops: 4.3, tdp: 67 },  // Updated: ~67W max package power
        'm3': { name: 'Apple M3', vram: 8, generation: 'Apple Silicon', tflops: 4.1, tdp: 45 },  // Updated: ~45W max package power
        'rx7900xtx': { name: 'Radeon RX 7900 XTX', vram: 24, generation: 'RDNA3', tflops: 61, tdp: 355 },  // Correct: 355W
        'rx7900xt': { name: 'Radeon RX 7900 XT', vram: 20, generation: 'RDNA3', tflops: 52, tdp: 315 },  // Correct: 315W
        'rx7900gre': { name: 'Radeon RX 7900 GRE', vram: 16, generation: 'RDNA3', tflops: 46, tdp: 260 },  // Correct: 260W
        'rx7800xt': { name: 'Radeon RX 7800 XT', vram: 16, generation: 'RDNA3', tflops: 37, tdp: 263 },  // Correct: 263W
        'rx7700xt': { name: 'Radeon RX 7700 XT', vram: 12, generation: 'RDNA3', tflops: 35, tdp: 245 },  // Correct: 245W
    };

    const gpuSpecs = Object.fromEntries(
        Object.entries(unsortedGpuSpecs)
            .sort(([, a], [, b]) => {
                // First sort by name prefix (A, GTX, RTX, etc.)
                const nameA = a.name.split(' ')[0];
                const nameB = b.name.split(' ')[0];
                if (nameA !== nameB) return nameA.localeCompare(nameB);
                // Then sort by VRAM if names are the same
                return a.vram - b.vram;
            })
    );

    const calculateRAMRequirements = (paramCount, quantBits, contextLength, gpuConfigs) => {
        // Add model size-based RAM requirements per Ollama docs
        const getMinimumRAM = (paramCount) => {
            if (paramCount <= 3) return 8;  // 3B models need 8GB
            if (paramCount <= 7) return 16; // 7B models need 16GB
            if (paramCount <= 13) return 32; // 13B models need 32GB
            return 64; // 70B models need 64GB
        };

        const minimumSystemRAM = getMinimumRAM(paramCount);
        
        // Calculate base model size in GB
        const baseModelSizeGB = (paramCount * quantBits * 1000000000) / (8 * 1024 * 1024 * 1024);

        // Calculate hidden size (d_model)
        const hiddenSize = Math.sqrt(paramCount * 1000000000 / 6);

        // Calculate KV cache size in GB
        const kvCacheSize = (2 * hiddenSize * contextLength * 2 * quantBits / 8) / (1024 * 1024 * 1024);

        // Add GPU overhead
        const gpuOverhead = baseModelSizeGB * 0.1;
        const totalGPURAM = baseModelSizeGB + kvCacheSize + gpuOverhead;

        // Calculate system RAM requirements
        const systemRAMMultiplier = getSystemRAMMultiplier(quantBits);
        const totalSystemRAM = totalGPURAM * systemRAMMultiplier;

        // Calculate total available VRAM across all GPU configs
        let totalAvailableVRAM = 0;
        gpuConfigs.forEach(config => {
            if (config.gpuModel) {
                const numGPUs = parseInt(config.count);
                const gpuVRAM = gpuSpecs[config.gpuModel].vram * numGPUs;
                totalAvailableVRAM += gpuVRAM;
            }
        });

        // Fix: Check if using multiple GPUs by comparing against first GPU's VRAM
        const firstGpuVRAM = gpuConfigs[0].gpuModel ? gpuSpecs[gpuConfigs[0].gpuModel].vram : 0;
        const multiGpuEfficiency = totalAvailableVRAM > firstGpuVRAM ? 0.9 : 1;
        const effectiveVRAM = totalAvailableVRAM * multiGpuEfficiency;

        // Add storage requirement calculation (approximately 10GB base + model size)
        const storageRequired = 10 + baseModelSizeGB;
        
        // Add CPU core requirements
        const recommendedCores = paramCount > 13 ? 8 : 4;
        
        return {
            baseModelSizeGB,
            kvCacheSize,
            totalGPURAM,
            totalSystemRAM,
            totalAvailableVRAM,
            effectiveVRAM,
            vramMargin: totalAvailableVRAM - totalGPURAM,
            minimumSystemRAM,
            storageRequired,
            recommendedCores,
            // Add warning if system requirements not met
            systemRequirementsMet: totalSystemRAM >= minimumSystemRAM
        };
    };

    const calculateTokensPerSecond = (paramCount, numGPUs, gpuModel, quantization) => {
        if (!gpuModel) return null;

        const selectedGPU = gpuSpecs[gpuModel];
        const baseTPS = (selectedGPU.tflops * 1e12) / (6 * paramCount * 1e9) * 0.05;
        
        // More accurate quantization factors based on research
        let quantizationFactor = 1;  // FP16 baseline
        switch(quantization) {
            case '32':
                quantizationFactor = 0.5;  // FP32 is slower
                break;
            case '8':
                quantizationFactor = 1.8;  // INT8 is significantly faster
                break;
            case '4':
                quantizationFactor = 2.2;  // INT4 provides highest throughput
                break;
        }

        let totalTPS = baseTPS * quantizationFactor;
        for(let i = 1; i < numGPUs; i++) {
            totalTPS += baseTPS * 0.9 * quantizationFactor;
        }
        
        return Math.round(Math.min(totalTPS, 200));
    };
```