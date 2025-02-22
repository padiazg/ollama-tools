package ollama

type MemoryEstimation struct {
	BaseModelSize float64
	KVCacheSize   float64
	GPURAM        float64
	SystemRAM     float64
}
