package asr

type ASRPool struct {
	pool []*ASR
}

func NewASRPool(poolSize int) *ASRPool {
	asrPool := &ASRPool{
		pool: make([]*ASR, poolSize),
	}

	for i := 0; i < poolSize; i++ {
		asr := NewASR()
		asr.Start()
		asrPool.pool[i] = asr
	}

	return asrPool
}

func (p *ASRPool) ProcessAudioChunk(chunk []byte) (string, error) {
	asr := p.getNextAvailableASR()

	return asr.ProcessAudioChunk(chunk)
}

func (p *ASRPool) getNextAvailableASR() *ASR {
	for _, asr := range p.pool {
		if asr.IsAvailable() {
			return asr
		}
	}

	asr := NewASR()
	asr.Start()
	p.pool = append(p.pool, asr)
	return asr
}

func (p *ASRPool) Stop() {
	for _, asr := range p.pool {
		asr.Stop()
	}
}
