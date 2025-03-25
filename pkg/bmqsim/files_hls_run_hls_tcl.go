package bmqsim

const (
	HLSRunHlsTcl = `open_project -reset proj
set_top circuit
add_files src/circuit.cc
# add_files -tb testbench.cc

# reset the solution
open_solution -reset "solution" -flow_target vitis
set_part {xcu55c-fsvh2892-2L-e}
create_clock -period 3.5

# just check that the C++ compiles
# csim_design -argv "2 6"

# synthethize the algorithm
csynth_design

export_design -format xo

exit
`
)
