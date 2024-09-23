![bmqsim](bmqsim.png)

bmqsim is part of BondMachine project. It is a library and a CLI tool to produce HDL code and test benches to simulate a quantum computer on FPGA using the BondMachine architecture framework. It also provides a sotfware simulator to be used in a CPU.
The installation of bmqsim is part of the standard installation of BondMachine framework. Please refer to the [BondMachine quick start](http://bondmachine.fisica.unipg.it/docs/#quickstart) for the instructions.

## Usage

To use the bmqsim CLI tool, you have to choose two things: the backend and the quantum algorithm to simulate.
The backend creates target hardware architecture of a specific type and its flavor (the subtype as for example the `only real numbers` or the `complex numbers` flavor) if applicable.
The quantum algorithm is a file containing the quantum circuit to simulate. Let's begin the latter.

## Quantum algorithm

The quantum circuit is a file containing the quantum algorithm to simulate. The file is written in the the same format as the one used by the Basm assembler and has the bmq extension. The file is a sequence of instructions, each one in a line inserted in a **block** section. To instruct the simulator, outside the block section, you can use the `%meta` directive to specify the simulator metadata.

The following is an example (*program.bmq*) of a quantum algorithm file implementing the Bell state:
```bmq
%block code1 .sequential
	qbits	q0,q1
	zero	q0,q1
	h	q0
	cx	q0,q1
%endblock

%meta bmdef global main:code1
```
## Backend

Once you have the quantum algorithm, you can choose the backend to use. The backend creates the target hardware architecture of a specific type and flavor. 
The subsequent tables show the available backend and architectures types, their flavors and the relevant command line options:

| Architecture type (backend) | Command line option | Description | Flavors | 
| --- | --- | --- | --- |
| `Software` | `--software-simulation` | Software simulator | `None` |
| `MatrixSeqHardcoded` | `--build-matrix-seq-hardcoded` | Sequence of matrices hardcoded in the FPGA | `seq_hardcoded_real`, `seq_hardcoded_complex` |
| `MatrixSeq` | `--build-matrix-seq` | Sequence of matrices in the FPGA | *Not yet implemented* |
| `FullHardwareHardcoded` | `--build-full-hardware-hardcoded` | Full hardware with hardcoded matrices | *Not yet implemented* |

Each architecture type is associated to a different set of command line options and activates a different simulator backend. Some of the options are common to all the backends, while others are specific to a single backend. The following sections describe the possible backends and the options.

### Common options

| Option | Description |
| --- | --- |
| `-hw-flavor-list` | List the available flavors for the chosen hardware architecture. |
| `-hw-flavor` | Choose the flavor of the hardware architecture. |
| `-save-basm` | Save the BondMachine Assembly code for the Basm assembler. |
| `-build-app` | Build the application template to use the selected hardware architecture. |
| `-app-flavor-list` | List the available flavors for the chosen application template. |
| `-app-flavor` | Choose the flavor of the application template. |
| `-app-file` | Choose the file name of the application template. |

### Software simulator

The software simulator backend is used to simulate the quantum algorithm on a CPU. The program accept a json input file and produces a json output file.
The following table shows the command line options specific to the software simulator:

| Option | Description |
| --- | --- |
| `-software-simulation-input` | Choose the input file for the software simulator. |
| `-software-simulation-output` | Choose the output file for the software simulator. |

An example of usage of the software simulator with the Bell state algorithm previously described (*program.bmq*) is the following:
```bash
bmqsim --software-simulation --software-simulation-input input.json --software-simulation-output output.json program.bmq
```
where `input.json`, the input file for the software simulator, is:
```json
{
	{"Vector":[{"Real":1,"Imag":0},{"Real":0,"Imag":0},{"Real":0,"Imag":0},{"Real":0,"Imag":0}]}
}
```
and `output.json`, the output file produced by the simulator, is:
```json
{
	{"Vector":[{"Real":0.7071067,"Imag":0},{"Real":0,"Imag":0},{"Real":0,"Imag":0},{"Real":0.7071067,"Imag":0}]}
}
```

### Hardcoded sequence of matrices

### Loadable sequence of matrices

### Full hardware with hardcoded matrices

## Other options

Alongside these main options, there are other options that does not depend on the backend chosen. These options are listed in the following table:

| Option | Description |
| --- | --- |


## Examples
