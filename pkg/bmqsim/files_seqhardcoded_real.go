package bmqsim

const (
	SeqHardcodedReal = `%section matrixmulel .romtext iomode:sync
        entry _start    ; Entry point
_start:
	rset	r2, {{ .NumGates }} ; the number of matrices
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

{{- range $i := n 0 .MatrixRows }}
{{- range $j := n 0 $.MatrixRows }}
%section matrixelement{{ $i }}_{{ $j }} .ramdata
{{- range $k := n 0 $.NumGates }}
	step{{ $k }} dd 0f{{ index (index (index $.MtxReal $k) $i) $j }}
{{- end }}
%endsection
{{ end }}
{{- end }}


%section matrixaddel .romtext iomode:sync
        entry _start    ; Entry point
_start:
	clr	r1
{{- range $i := n 0 .MatrixRows }}
	mov	r0, i{{ $i }}
	{{"{{"}} .Params.addop {{"}}"}}	r1, r0
{{- end }}
	mov	o0, r1
	j _start
%endsection

%section main .romtext iomode:sync
	entry _start    ; Entry point
_start:

	clr	r1
{{- range $i := n 0 .MatrixRows }}
	mov r0, i{{ $i }}
	mov ram:[r1], r0
	inc r1
{{- end }}
	rset	r2, {{ .NumGates }} ; the number of matrices
mainloop:
	jz	r2, endloop

	clr r1
{{- range $i := n 0 .MatrixRows }}
	mov r0, ram:[r1]
	mov o{{ sum $.MatrixRows $i }}, r0
	inc r1
{{- end }}

	clr	r1
{{- range $i := n 0 .MatrixRows }}
	mov r0, i{{ sum $.MatrixRows $i }}
	mov ram:[r1], r0
	inc r1
{{- end }}
	dec r2
	j mainloop

endloop:
	clr r1
{{- range $i := n 0 .MatrixRows }}	
	mov r0, ram:[r1]
	mov o{{ $i }}, r0
	inc r1
{{- end }}
	j _start
%endsection

%section mainram .ramdata
{{- range $i := n 0 .NumGates }}
        res{{ $i }} dd 0f0
{{- end }}
%endsection

{{- range $i := n 0 .MatrixRows }}
{{- range $j := n 0 $.MatrixRows }}
%meta cpdef	mult{{ $i }}_{{ $j }}	romcode: matrixmulel, ramdata: matrixelement{{ $i }}_{{ $j }}, execmode: ha, multop: multf
{{- end }}
{{- end }}

{{- range $i := n 0 .MatrixRows }}
%meta cpdef	add{{ $i }}	romcode: matrixaddel, execmode: ha, addop: addf
{{- end }}

%meta cpdef	main	romcode: main, ramdata:mainram, execmode: ha

{{- range $i := n 0 .MatrixRows }}
{{- range $j := n 0 $.MatrixRows }}
%meta ioatt	toadd{{ $i }}_{{ $j }}	cp: mult{{ $i }}_{{ $j }}, index:0, type:output
%meta ioatt     toadd{{ $i }}_{{ $j }}	cp: add{{ $i }}, index:{{ $j }}, type:input
{{- end }}
{{- end }}

{{- range $i := n 0 .MatrixRows }}
{{- range $j := n 0 $.MatrixRows }}
%meta ioatt     tomult{{ $i }}_{{ $j }}	cp: main, index:{{ sum $.MatrixRows $j }}, type:output
%meta ioatt     tomult{{ $i }}_{{ $j }}	cp: mult{{ $i }}_{{ $j }}, index:0, type:input
{{- end }}
{{- end }}

{{- range $i := n 0 .MatrixRows }}
%meta ioatt     tomain{{ $i }}	cp: add{{ $i }}, index:0, type:output
%meta ioatt     tomain{{ $i }}	cp: main, index:{{ sum $.MatrixRows $i }}, type:input
{{- end }}

{{- range $i := n 0 .MatrixRows }}
%meta ioatt     bmi{{ $i }}	cp: main, index:{{ $i }}, type:input
%meta ioatt     bmi{{ $i }}	cp: bm, index:{{ $i }}, type:input
{{- end }}

{{- range $i := n 0 .MatrixRows }}
%meta ioatt     bmo{{ $i }}	cp: main, index:{{ $i }}, type:output
%meta ioatt     bmo{{ $i }}	cp: bm, index:{{ $i }}, type:output
{{- end }}

%meta bmdef	global registersize:32
`
)
