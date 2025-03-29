package bmqsim

const (
	HLSPythonPynqComplex = `import pynq
import numpy as np 
from pynq import allocate
import time

ol = pynq.Overlay("./firmware.xclbin")
keys = ol._keys()
print(keys)


circuit = ol.circuit_1
dimension = {{ mult 2 .MatrixRows }}

z_buf = allocate(shape=(dimension) , dtype=np.float32)
z_buf[0] = 1.0
for i in range(1, dimension):
    z_buf[i] = 0.0

start_fpga = time.time()

z_buf.sync_to_device()

circuit.call(dimension, z_buf)
z_buf.sync_from_device()
print("\nRisultato FPGA:  " , z_buf, "\n")

end_fpga = time.time()
fpga_time = end_fpga - start_fpga
print(f"FPGA completed in {fpga_time:.6f} seconds.\n")

try:
    ol.free()
    z_buf.freebuffer()
except Exception as e:
    print("error deleting buffers", e)

`
)
