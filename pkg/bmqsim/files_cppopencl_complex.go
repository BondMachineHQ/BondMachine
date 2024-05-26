package bmqsim

const (
	CppOpenCLComplex = `#include "xcl2.hpp"
#include <vector>

#define BATCH_SIZE 16
#define BM_INPUT {{ mult 2 .MatrixRows }}
#define BM_OUTPUT {{ mult 2 .MatrixRows }}
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

        std::string bmFirmwareFile = "firmware.xclbin";

        // Allocate Memory in Host Memory
        auto vector_size_bytes_in = sizeof(float) * TOTIN;
        auto vector_size_bytes_out = sizeof(float) * TOTOUT;

        std::vector<float, aligned_allocator<float>> source_input(TOTIN);
        std::vector<float, aligned_allocator<float>> source_output(TOTOUT);

        // OpenCL context, device, program, kernel
        cl_int err;
        cl::CommandQueue q;
        cl::Context context;
        cl::Program program;

        auto devices = xcl::get_xil_devices();

        // load the BM firmware and return the pointer to file buffer.
        auto fileBuf = xcl::read_binary_file(bmFirmwareFile);
        cl::Program::Binaries bins{{"{{"}}fileBuf.data(), fileBuf.size(){{"}}"}};
        bool valid_device = false;
        for (unsigned int i = 0; i < devices.size(); i++)
        {
                auto device = devices[i];

                // Creating Context and Command Queue for selected Device
                OCL_CHECK(err, context = cl::Context(device, nullptr, nullptr, nullptr, &err));
                OCL_CHECK(err, q = cl::CommandQueue(context, device, CL_QUEUE_OUT_OF_ORDER_EXEC_MODE_ENABLE | CL_QUEUE_PROFILING_ENABLE, &err));

                // std::cout << "Trying to program device[" << i << "]: " << device.getInfo<CL_DEVICE_NAME>() << std::endl;
                program = cl::Program(context, {device}, bins, nullptr, &err);
                if (err != CL_SUCCESS)
                {
                        std::cout << "Failed to program device[" << i << "] with xclbin file!\n";
                }
                else
                {
                        std::cout << "Device[" << i << "]: program successful!\n";
                        valid_device = true;
                        break; // we break because we found a valid device
                }
        }
        if (!valid_device)
        {
                std::cout << "Failed to program any device found, exit!\n";
                exit(EXIT_FAILURE);
        }

        OCL_CHECK(err, cl::Kernel krnl_bondmachine(program, "krnl_bondmachine_rtl", &err));

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

                float *data = (float *)malloc(sizeof(float) * i);
                i = 0;

                while ((read = getline(&line, &len, fp)) != -1)
                {
                        data[i] = atof(line);
                        // printf("%f\n", data[i]);
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
                                        source_input[i] = data[i + offset];
                                else
                                {
                                        if (threshold == BATCH_SIZE)
                                                threshold = i / BM_INPUT;
                                        source_input[i] = 0;
                                }

                        // Allocate Buffer in Global Memory
                        OCL_CHECK(err, cl::Buffer buffer_input(context, CL_MEM_USE_HOST_PTR | CL_MEM_READ_ONLY, vector_size_bytes_in, source_input.data(), &err));
                        OCL_CHECK(err, cl::Buffer buffer_output(context, CL_MEM_USE_HOST_PTR | CL_MEM_WRITE_ONLY, vector_size_bytes_out, source_output.data(), &err));

                        // Set the Kernel Arguments
                        OCL_CHECK(err, err = krnl_bondmachine.setArg(0, buffer_input));
                        OCL_CHECK(err, err = krnl_bondmachine.setArg(1, buffer_output));

                        // Copy input data to device global memory
                        cl::Event write_event;
                        OCL_CHECK(err, err = q.enqueueMigrateMemObjects({buffer_input}, 0 /* 0 means from host*/, nullptr, &write_event));

                        // Launch the Kernel
                        std::vector<cl::Event> eventVec;
                        eventVec.push_back(write_event);
                        OCL_CHECK(err, err = q.enqueueTask(krnl_bondmachine, &eventVec));

                        // wait for all kernels to finish their operations
                        OCL_CHECK(err, err = q.finish());

                        // Copy Result from Device Global Memory to Host Local Memory
                        OCL_CHECK(err, err = q.enqueueMigrateMemObjects({buffer_output}, CL_MIGRATE_MEM_OBJECT_HOST));
                        OCL_CHECK(err, err = q.finish());
                        for (int i = 0; i < TOTOUT; i++)
                                if (i < threshold * BM_OUTPUT)
                                        fprintf(fp, "%.12f\n", source_output[i]);
                }

                fclose(fp);
        }

        else // If a file is not provided, compute the zero state
        {
                for (int i = 0; i < TOTIN; i++)
                        source_input[i] = 0;
                source_input[0] = 1.0;

                printf("%f %f %f %f %f %f %f %f \n", source_input[0], source_input[1], source_input[2], source_input[3], source_input[4], source_input[5], source_input[6], source_input[7]);

                // Allocate Buffer in Global Memory
                OCL_CHECK(err, cl::Buffer buffer_input(context, CL_MEM_USE_HOST_PTR | CL_MEM_READ_ONLY, vector_size_bytes_in, source_input.data(), &err));
                OCL_CHECK(err, cl::Buffer buffer_output(context, CL_MEM_USE_HOST_PTR | CL_MEM_WRITE_ONLY, vector_size_bytes_out, source_output.data(), &err));

                // Set the Kernel Arguments
                OCL_CHECK(err, err = krnl_bondmachine.setArg(0, buffer_input));
                OCL_CHECK(err, err = krnl_bondmachine.setArg(1, buffer_output));

                // Copy input data to device global memory
                cl::Event write_event;
                OCL_CHECK(err, err = q.enqueueMigrateMemObjects({buffer_input}, 0 /* 0 means from host*/, nullptr, &write_event));

                // Launch the Kernel
                std::vector<cl::Event> eventVec;
                eventVec.push_back(write_event);
                OCL_CHECK(err, err = q.enqueueTask(krnl_bondmachine, &eventVec));

                // wait for all kernels to finish their operations
                OCL_CHECK(err, err = q.finish());

                // Copy Result from Device Global Memory to Host Local Memory
                OCL_CHECK(err, err = q.enqueueMigrateMemObjects({buffer_output}, CL_MIGRATE_MEM_OBJECT_HOST));
                OCL_CHECK(err, err = q.finish());

	        for (int i = 0; i < BM_OUTPUT; i++) {
	                std::cout << source_output[i] << std::endl;
                }
        }
        return 0;
}
`
)
