package bmqsim

const (
	CppOpenCLReal = `#include "xcl2.hpp"
#include <vector>

#define BATCH_SIZE 16
#define BM_INPUT {{ .MatrixRows }}
#define BM_OUTPUT {{ .MatrixRows }}
#define PRECISION 32

#define TOTIN BATCH_SIZE*BM_INPUT
#define TOTOUT BATCH_SIZE*BM_OUTPUT

int main() {
        std::string bmFirmwareFile = "firmware.xclbin";

        // Allocate Memory in Host Memory
        auto vector_size_bytes_in = sizeof(float) * TOTIN;
        auto vector_size_bytes_out = sizeof(float) * TOTOUT;

        std::vector<float, aligned_allocator<float> > source_input(TOTIN);
        std::vector<float, aligned_allocator<float> > source_output(TOTOUT);

        // Create zero state
        for (int i = 0; i < TOTIN ; i++) {
                source_input[i] = 0.0;
        }
        source_input[0] = 1.0;

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
        for (unsigned float i = 0; i < devices.size(); i++) {
                auto device = devices[i];

                // Creating Context and Command Queue for selected Device
                OCL_CHECK(err, context = cl::Context(device, nullptr, nullptr, nullptr, &err));
                OCL_CHECK(err, q = cl::CommandQueue(context, device, CL_QUEUE_OUT_OF_ORDER_EXEC_MODE_ENABLE | CL_QUEUE_PROFILING_ENABLE, &err));

                // std::cout << "Trying to program device[" << i << "]: " << device.getInfo<CL_DEVICE_NAME>() << std::endl;
                program = cl::Program(context, {device}, bins, nullptr, &err);
                if (err != CL_SUCCESS) {
                        std::cout << "Failed to program device[" << i << "] with xclbin file!\n";
                } else {
                        std::cout << "Device[" << i << "]: program successful!\n";
                        valid_device = true;
                        break; // we break because we found a valid device
                }
        }
        if (!valid_device) {
                std::cout << "Failed to program any device found, exit!\n";
                exit(EXIT_FAILURE);
        }

        OCL_CHECK(err, cl::Kernel krnl_bondmachine(program, "krnl_bondmachine_rtl", &err));

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

        for (int i = 0; i < TOTOUT; i++) {
                std::cout << " Device result = " << source_output[i] << std::endl;
        }

        return 0; 
}
`
)
