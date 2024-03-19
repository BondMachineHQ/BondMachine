package bmqsim

const (
	CPynqApiReal = `#include <stdio.h>
#include <pynq_api.h>

#define BATCH_SIZE 16
#define BM_INPUT {{ .MatrixRows }}
#define BM_OUTPUT {{ .MatrixRows }}
#define PRECISION 32

#define TOTIN BATCH_SIZE*BM_INPUT
#define TOTOUT BATCH_SIZE*BM_OUTPUT

int main() {
        PYNQ_loadBitstream("firmware.bit");

        PYNQ_SHARED_MEMORY shared_memory_1, shared_memory_2;
        PYNQ_allocatedSharedMemory(&shared_memory_1, sizeof(float)*TOTIN, 1);
        PYNQ_allocatedSharedMemory(&shared_memory_2, sizeof(float)*TOTOUT, 1);
  
        float * d1=(float*)shared_memory_1.pointer;
        float * d2=(float*)shared_memory_2.pointer;

        for (int i=0;i<TOTIN;i++) {
                d2[i]=0;
        }
  
        d1[0]=1.0;

        printf("{{- range $i := n 0 .MatrixRows }}%f {{ end }}\n"{{- range $i := n 0 .MatrixRows }},d1[{{ $i }}]{{ end }});

        PYNQ_AXI_DMA dma;
        PYNQ_openDMA(&dma, 0x40400000);
  
        PYNQ_writeDMA(&dma, &shared_memory_1, 0, sizeof(float)*TOTIN);
        PYNQ_readDMA(&dma, &shared_memory_2, 0, sizeof(float)*TOTOUT);

        PYNQ_waitForDMAComplete(&dma, AXI_DMA_WRITE);
        PYNQ_waitForDMAComplete(&dma, AXI_DMA_READ);
        
        printf("{{- range $i := n 0 .MatrixRows }}%f {{ end }}\n"{{- range $i := n 0 .MatrixRows }},d2[{{ $i }}]{{ end }});
        
        PYNQ_closeDMA(&dma);
        PYNQ_freeSharedMemory(&shared_memory_1);
        PYNQ_freeSharedMemory(&shared_memory_2);
        
        return 0;
}
`
)
