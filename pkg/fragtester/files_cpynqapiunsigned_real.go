package fragtester

const (
	CPynqApiUnsigned = `#include <stdio.h>
#include <pynq_api.h>

#define BATCH_SIZE 16
#define BM_INPUT {{ len .Inputs }}
#define BM_OUTPUT {{ len .Outputs }}
#define PRECISION 32

#define TOTIN BATCH_SIZE*BM_INPUT
#define TOTOUT BATCH_SIZE*BM_OUTPUT

int main(int argc, char *argv[])
{

        FILE *fp;
        char *line = NULL;
        size_t len = 0;
        ssize_t read;
        int i = 0;
        int total = 0;
        int offset = 0;
        int threshold = 0;

        PYNQ_loadBitstream("firmware.bit");

        PYNQ_SHARED_MEMORY shared_memory_1, shared_memory_2;
        PYNQ_allocatedSharedMemory(&shared_memory_1, 4 * TOTIN, 1);
        PYNQ_allocatedSharedMemory(&shared_memory_2, 4 * TOTOUT, 1);

        u_int32_t *d1 = (u_int32_t *)shared_memory_1.pointer;
        u_int32_t *d2 = (u_int32_t *)shared_memory_2.pointer;

        PYNQ_AXI_DMA dma;
        PYNQ_openDMA(&dma, 0x40400000);

        if (argc == 3) // argv[1] is the input file, argv[2] is the output file
        {
                fp = fopen(argv[1], "r");
                if (fp == NULL)
                        printf("Input file %s not found\n", argv[1]), exit(EXIT_FAILURE);
                while ((read = getline(&line, &len, fp)) != -1)
                        i++;

                if (i % BM_INPUT != 0)
                        printf("Input file %s has wrong number of elements\n", argv[1]), exit(EXIT_FAILURE);

                rewind(fp);

                u_int32_t *data = (u_int32_t *)malloc(4 * i);
                i = 0;

                while ((read = getline(&line, &len, fp)) != -1)
                {
                        data[i] = atoi(line);
                        printf("%d\n", data[i]);
                        i++;
                }
                fclose(fp);
                total = i;

                fp = fopen(argv[2], "w");
                if (fp == NULL)
                        printf("Output file %s cannot be opened\n", argv[2]), exit(EXIT_FAILURE);

                threshold = BATCH_SIZE;
                for (offset = 0; offset < total; offset += TOTIN)
                {
                        for (int i = 0; i < TOTIN; i++)
                                if (i + offset < total)
                                        d1[i] = data[i + offset];
                                else
                                {
                                        if (threshold == BATCH_SIZE)
                                                threshold = i/BM_INPUT;
                                        d1[i] = 0;
                                }

                        PYNQ_writeDMA(&dma, &shared_memory_1, 0, 4 * TOTIN);
                        PYNQ_readDMA(&dma, &shared_memory_2, 0, 4 * TOTOUT);

                        PYNQ_waitForDMAComplete(&dma, AXI_DMA_WRITE);
                        PYNQ_waitForDMAComplete(&dma, AXI_DMA_READ);

                        for (int i = 0; i < TOTOUT; i++)
                                if (i < threshold * BM_OUTPUT) {
						unsigned int val = *(unsigned int *)&d2[i];
						fprintf(fp,"%u\n", val);
		                                //for (int bit = 31; bit >= 0; bit--)
                                                //       fprintf(fp, "%d", (val >> bit) & 1);
                                                //     fprintf(fp, "\n");
				}
                }

                fclose(fp);
        }
       
	PYNQ_closeDMA(&dma);
        PYNQ_freeSharedMemory(&shared_memory_1);
        PYNQ_freeSharedMemory(&shared_memory_2);

        return 0;
}
`
)
