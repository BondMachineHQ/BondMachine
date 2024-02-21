package bmqsim

const (
	SeqHardcodedReal = `%section matrixmulel .romtext iomode:sync
        entry _start    ; Entry point
_start:
	rset	r2, 2 ; the number of matrices
	rset	r1, 0 ; counter
mainloop:
	jz	r2, _start
	mov	r3, ram:[r1]
	mov	r0, i0
	{{"{{"}} .Params.multop {{"}}"}}	r0, r3
	mov	o0, r0
	dec	r2
	inc	r1
	j mainloop
	mov	ram:[r1], r0
%endsection

%section matrixelement00 .ramdata
        step0 dd 0f2
	step1 dd 0f-1
%endsection

%section matrixelement01 .ramdata
        step0 dd 0f3
	step1 dd 0f1.1
%endsection

%section matrixelement10 .ramdata
        step0 dd 0f3
	step1 dd 0f1.1
%endsection

%section matrixelement11 .ramdata
        step0 dd 0f3
	step1 dd 0f1.1
%endsection


%section matrixaddel .romtext iomode:sync
        entry _start    ; Entry point
_start:
	clr	r1
	mov	r0, i0
	{{"{{"}} .Params.addop {{"}}"}}	r1, r0
	mov	r0, i1
	{{"{{"}} .Params.addop {{"}}"}}	r1, r0
	mov	o0, r1
	j _start
%endsection

%section main .romtext iomode:sync
	entry _start    ; Entry point
_start:
	;clr	r1
	mov r0, i0
	;mov ram:[r1], r0
	;inc r1
	mov r1, i1
	;mov ram:[r1], r0
;	j endloop
	
	rset	r2, 2 ; the number of matrices
mainloop:
	jz	r2, endloop
	;clr r1
	;mov r0, ram:[r1]
	mov o2, r0
	;inc r1
	;mov r0, ram:[r1]
	mov o3, r1

	;clr r1
	mov r0, i2
	;mov ram:[r1], r0
	;inc r1
	mov r1, i3
	;mov ram:[r1], r0

	dec r2
	j mainloop

endloop:
	;clr r1
	;mov r0, ram:[r1]
	mov o0, r0
	;inc r1
	;mov r0, ram:[r1]
	mov o1, r1

	j _start
%endsection

%section mainram .ramdata
        res0 dd 0f0
        res1 dd 0f0
%endsection

%meta cpdef	mult00	romcode: matrixmulel, ramdata: matrixelement00, execmode: ha, multop: multf
%meta cpdef	mult01	romcode: matrixmulel, ramdata: matrixelement01, execmode: ha, multop: multf
%meta cpdef	mult10	romcode: matrixmulel, ramdata: matrixelement10, execmode: ha, multop: multf
%meta cpdef	mult11	romcode: matrixmulel, ramdata: matrixelement11, execmode: ha, multop: multf

%meta cpdef	add0	romcode: matrixaddel, execmode: ha, addop: addf
%meta cpdef	add1	romcode: matrixaddel, execmode: ha, addop: addf

%meta cpdef	main	romcode: main, ramdata:mainram, execmode: ha

%meta ioatt     toadd00	cp: mult00, index:0, type:output
%meta ioatt     toadd01	cp: mult01, index:0, type:output
%meta ioatt     toadd10	cp: mult10, index:0, type:output
%meta ioatt     toadd11	cp: mult11, index:0, type:output

%meta ioatt     toadd00	cp: add0, index:0, type:input
%meta ioatt     toadd01	cp: add0, index:1, type:input
%meta ioatt     toadd10	cp: add1, index:0, type:input
%meta ioatt     toadd11	cp: add1, index:1, type:input

%meta ioatt     tomult00	cp: main, index:2, type:output
%meta ioatt     tomult01	cp: main, index:3, type:output
%meta ioatt     tomult00	cp: mult00, index:0, type:input
%meta ioatt     tomult01	cp: mult01, index:0, type:input
%meta ioatt     tomult10	cp: main, index:2, type:output
%meta ioatt     tomult11	cp: main, index:3, type:output
%meta ioatt     tomult10	cp: mult10, index:0, type:input
%meta ioatt     tomult11	cp: mult11, index:0, type:input

%meta ioatt     tomain0	cp: add0, index:0, type:output
%meta ioatt     tomain1	cp: add1, index:0, type:output
%meta ioatt     tomain0	cp: main, index:2, type:input
%meta ioatt     tomain1	cp: main, index:3, type:input

%meta ioatt     bmi0	cp: main, index:0, type:input
%meta ioatt     bmi1	cp: main, index:1, type:input
%meta ioatt     bmi0	cp: bm, index:0, type:input
%meta ioatt     bmi1	cp: bm, index:1, type:input

%meta ioatt     bmo0	cp: main, index:0, type:output
%meta ioatt     bmo1	cp: main, index:1, type:output
%meta ioatt     bmo0	cp: bm, index:0, type:output
%meta ioatt     bmo1	cp: bm, index:1, type:output

%meta bmdef	global registersize:32
`
)
