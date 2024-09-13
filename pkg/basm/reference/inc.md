# inc

**Instruction**: inc

**Length**: opBits + regBits

**Description**:

The INC instruction increments a register by 1. It assumes that the register is a positive integer and it will overflow to 0 if the maximum value is reached. The length of the register is defined by the architecture.        

**Snippet**:


```asm
%section code .romtext
        entry _start    ; Entry point
_start:

        clr     r0
loop:
        inc	r0
	r2o	r0,o0
        j	loop

%endsection

%meta cpdef	cpu	romcode: code, ramsize:8
%meta ioatt     testio cp: cpu, index:0, type:output
%meta ioatt     testio cp: bm, index:0, type:output
%meta bmdef	global registersize:8
```

INC example
