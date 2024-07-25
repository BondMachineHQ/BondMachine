
The main purpose of the `procbuilder` command is to provide a command-line interface for interacting with the `procbuilder` package. The `procbuilder` package is the library that contains the core functionality for building and interacting with a CP processor architecture. The command accepts various flags that control its behavior and perform different operations based on those flags. This document provides an overview of the command's functionality, the flags it accepts, and the operations it performs based on those flags.

This command performs various operations based on the provided flags. It can create a new CP, add instructions to the CP, assemble a program, disassemble a program and produce the HDL code for the CP. 

### Usage

The two main workflows for the `procbuilder` command are:

- **Loading an existing CP** and program from JSON files, inspecting the CP and program

- **Creating a new CP**, adding instructions to the CP, assembling a program, disassembling the program and saving the machine state to a JSON file

### Loading an Existing CP and Program

The flag to load an existing CP and program is `load_machine`. If this flag is set to a non-empty string, the command checks if a file with the specified name exists. If the file exists, the command reads the machine state from the file and initializes the CP and program with the read data. If any errors occur during file operations or JSON unmarshaling, the command panics.

### TODO - Finish this section

### Creating a New CP and Program

The flags to create a new CP and program are `save_machine`. It is used to save the machine state to a JSON file. If this flag is set to a non-empty string, the command saves the machine state to a JSON file with the specified name.

### TODO - Finish this section