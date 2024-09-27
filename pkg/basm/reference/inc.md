# inc

**Instruction**: inc

**Length**: opBits + regBits

**Description**:

The INC instruction increments a register by 1. It assumes that the register is a positive integer and it will overflow to 0 if the maximum value is reached. The length of the register is defined by the architecture.        

**Snippet**:


```asm
%section code .romtext  ; The section name is code and it is a ROM program
        entry _start    ; Entry point
_start:

        clr     r0      ; Clear register r0
loop:
        inc	r0      ; Increment register r0
	r2o	r0,o0   ; Output register r0 to the special output register o0
        j	loop    ; Jump to loop

%endsection

%meta cpdef	cpu	romcode: code, ramsize:8        ; Define a CPU with the code section as program 256 bytes of RAM
%meta ioatt     testio cp: cpu, index:0, type:output    ; Define an output bond to the CP (outgoing, CP endpoint)
%meta ioatt     testio cp: bm, index:0, type:output     ; Define an output bond to the BM (outgoing, BM endpoint)
%meta bmdef	global registersize:8                   ; Define the register size of the BM to 8 bits
```

INC example
