package bmqsim

const (
	SeqHardcodedAddTreeComplex = `%section matrixmulel .romtext iomode:sync
        entry _start    ; Entry point
_start:
	rset	r2, {{ .NumGates }} ; the number of matrices
	rset	r1, 0 ; counter
mainloop:
	jz	r2, _start
	mov	r3, ram:[r1]
	inc	r1
	mov	r4, ram:[r1]
	mov	r0, i0
	mov	r5, i1

	{{"{{"}} .Params.multop {{"}}"}}	r0, r3
	mov	r6, r0
	mov	r7, 0f-1
	{{"{{"}} .Params.multop {{"}}"}}	r7, r5
	{{"{{"}} .Params.multop {{"}}"}}	r7, r4
	{{"{{"}} .Params.addop {{"}}"}}	r6, r7 ; real part now in r6

	{{"{{"}} .Params.multop {{"}}"}}	r0, r4 ; r0 now is a1*b2
	mov	r7, r3
	{{"{{"}} .Params.multop {{"}}"}}	r7, r5 ; r7 now is a2*b1
	{{"{{"}} .Params.addop {{"}}"}}	r0, r7 ; r0 now is a1*b2 + a2*b1 (imaginary part)

	mov	o0, r6
	mov	o1, r0
	dec	r2
	inc	r1
	j mainloop
	mov	ram:[r1], r0
%endsection

{{- range $i := n 0 .MatrixRows }}
{{- range $j := n 0 $.MatrixRows }}
%section matrixelement{{ $i }}_{{ $j }} .ramdata
{{- range $k := n 0 $.NumGates }}
	step{{ $k }}r dd 0f{{ index (index (index $.MtxReal $k) $i) $j }}
	step{{ $k }}i dd 0f{{ index (index (index $.MtxImag $k) $i) $j }}
{{- end }}
%endsection
{{ end }}
{{- end }}


%section matrixaddel .romtext iomode:sync
        entry _start    ; Entry point
_start:
	clr	r0 ; real part zero
	clr	r1 ; imag part zero
	mov	r2, i0
	mov	r3, i1
	{{"{{"}} .Params.addop {{"}}"}}	r0, r2
	{{"{{"}} .Params.addop {{"}}"}}	r1, r3
	mov	r2, i2
	mov	r3, i3
	{{"{{"}} .Params.addop {{"}}"}}	r0, r2
	{{"{{"}} .Params.addop {{"}}"}}	r1, r3
	mov	o0, r0 ; real part
	mov	o1, r1 ; imag part
	j _start
%endsection

%section main .romtext iomode:sync
	entry _start    ; Entry point
_start:

	clr	r1
{{- range $i := n 0 (mult 2 .MatrixRows) }}
	mov r0, i{{ $i }}
	mov ram:[r1], r0
	inc r1
{{- end }}
	rset	r2, {{ .NumGates }} ; the number of matrices
mainloop:
	jz	r2, endloop

	clr r1
{{- range $i := n 0 (mult 2 .MatrixRows) }}
	mov r0, ram:[r1]
	mov o{{ sum (mult 2 $.MatrixRows) $i }}, r0
	inc r1
{{- end }}

	clr	r1
{{- range $i := n 0 (mult 2 .MatrixRows) }}
	mov r0, i{{ sum (mult 2 $.MatrixRows) $i }}
	mov ram:[r1], r0
	inc r1
{{- end }}
	dec r2
	j mainloop

endloop:
	clr r1
{{- range $i := n 0 (mult 2 .MatrixRows) }}	
	mov r0, ram:[r1]
	mov o{{ $i }}, r0
	inc r1
{{- end }}
	j _start
%endsection

%section mainram .ramdata
{{- range $i := n 0 .MatrixRows }}
        res{{ $i }}r dd 0f0
	res{{ $i }}i dd 0f0
{{- end }}
%endsection

{{- range $i := n 0 .MatrixRows }}
{{- range $j := n 0 $.MatrixRows }}
%meta cpdef	mult{{ $i }}_{{ $j }}	romcode: matrixmulel, ramdata: matrixelement{{ $i }}_{{ $j }}, execmode: ha, multop: multf, addop: addf
{{- end }}
{{- end }}

{{/*
{{- define "addline" -}}
{{- $level := .level -}}
{{- $prefix := .prefix -}}
 	{{- if eq $level 0 }}
 %meta cpdef	add {{ $prefix }} romcode: matrixaddel, execmode: ha, addop: addf
 	{{- else -}}
 		{{- range $i := n 0 (pow $level 2) -}}
 			{{- template "addline" dict "level" (dec $level) "prefix" (printf "%s_%d" $prefix $i) -}}
 		{{- end -}}
 	{{- end -}}
 {{- end -}}

%meta cpdef	add{{- template "addline" dict "level" 3 "prefix" "g" -}}	romcode: matrixaddel, execmode: ha, addop: addf
*/}}

{{- range $r := n 0 .MatrixRows }}
{{- range $i := n 0 $.Qbits }}
{{- range $j := n 0 (pow 2 $i) }}
%meta cpdef	add{{ $r }}_{{ $i }}_{{ $j }} romcode: matrixaddel, execmode: ha, addop: addf
{{- end -}}
{{- end -}}
{{ end }}

%meta cpdef	main	romcode: main, ramdata:mainram, execmode: ha

{{- range $r := n 0 .MatrixRows }}
{{- range $i := n 0 $.Qbits}}
{{- range $j := ns 0 (pow 2 $i) 2 }}
{{- if ne $i 0 }}
%meta ioatt	addtree{{ $r }}_{{ $i }}_{{ $j }}_r	cp: add{{ $r }}_{{ $i }}_{{ $j }}, index:0, type:output
%meta ioatt	addtree{{ $r }}_{{ $i }}_{{ $j }}_r	cp: add{{ $r }}_{{ dec $i }}_{{ div $j 2 }}, index:0, type:input
%meta ioatt	addtree{{ $r }}_{{ $i }}_{{ $j }}_i	cp: add{{ $r }}_{{ $i }}_{{ $j }}, index:1, type:output
%meta ioatt	addtree{{ $r }}_{{ $i }}_{{ $j }}_i	cp: add{{ $r }}_{{ dec $i }}_{{ div $j 2 }}, index:1, type:input
%meta ioatt	addtree{{ $r }}_{{ $i }}_{{ inc $j }}_r	cp: add{{ $r }}_{{ $i }}_{{ inc $j }}, index:0, type:output
%meta ioatt	addtree{{ $r }}_{{ $i }}_{{ inc $j }}_r	cp: add{{ $r }}_{{ dec $i }}_{{ div $j 2 }}, index:2, type:input
%meta ioatt	addtree{{ $r }}_{{ $i }}_{{ inc $j }}_i	cp: add{{ $r }}_{{ $i }}_{{ inc $j }}, index:1, type:output
%meta ioatt	addtree{{ $r }}_{{ $i }}_{{ inc $j }}_i	cp: add{{ $r }}_{{ dec $i }}_{{ div $j 2 }}, index:3, type:input
{{- end -}}
{{- end -}}
{{- end -}}
{{ end }}

{{- range $r := n 0 .MatrixRows }}
{{- range $j := ns 0 $.MatrixRows 2 }}
%meta ioatt	toadd{{ $r }}_{{ $j }}_r	cp: mult{{ $r }}_{{ $j }}, index:0, type:output
%meta ioatt     toadd{{ $r }}_{{ $j }}_r	cp: add{{ $r }}_{{ dec $.Qbits }}_{{ div $j 2 }}, index:0, type:input
%meta ioatt	toadd{{ $r }}_{{ $j }}_i	cp: mult{{ $r }}_{{ $j }}, index:1, type:output
%meta ioatt     toadd{{ $r }}_{{ $j }}_i	cp: add{{ $r }}_{{ dec $.Qbits }}_{{ div $j 2 }}, index:1, type:input
%meta ioatt	toadd{{ $r }}_{{ inc $j }}_r	cp: mult{{ $r }}_{{ inc $j }}, index:0, type:output
%meta ioatt     toadd{{ $r }}_{{ inc $j }}_r	cp: add{{ $r }}_{{ dec $.Qbits }}_{{ div $j 2 }}, index:2, type:input
%meta ioatt	toadd{{ $r }}_{{ inc $j }}_i	cp: mult{{ $r }}_{{ inc $j }}, index:1, type:output
%meta ioatt     toadd{{ $r }}_{{ inc $j }}_i	cp: add{{ $r }}_{{ dec $.Qbits }}_{{ div $j 2 }}, index:3, type:input
{{- end }}
{{- end }}

{{- range $i := n 0 .MatrixRows }}
{{- range $j := n 0 $.MatrixRows }}
%meta ioatt     tomult{{ $i }}_{{ $j }}_r	cp: main, index:{{ sum (mult 2 $.MatrixRows) (mult 2 $j) }}, type:output
%meta ioatt     tomult{{ $i }}_{{ $j }}_i	cp: main, index:{{ sum (mult 2 $.MatrixRows) (sum 1 (mult 2 $j)) }}, type:output
%meta ioatt     tomult{{ $i }}_{{ $j }}_r	cp: mult{{ $i }}_{{ $j }}, index:0, type:input
%meta ioatt     tomult{{ $i }}_{{ $j }}_i	cp: mult{{ $i }}_{{ $j }}, index:1, type:input
{{- end }}
{{- end }}

{{- range $i := n 0 .MatrixRows }}
%meta ioatt     tomain{{ $i }}_r	cp: add{{ $i }}_0_0, index:0, type:output
%meta ioatt     tomain{{ $i }}_i	cp: add{{ $i }}_0_0, index:1, type:output
%meta ioatt     tomain{{ $i }}_r	cp: main, index:{{ sum (mult 2 $.MatrixRows) (mult 2 $i) }}, type:input
%meta ioatt     tomain{{ $i }}_i	cp: main, index:{{ sum (mult 2 $.MatrixRows) (sum 1 (mult 2 $i)) }}, type:input
{{- end }}

{{- range $i := n 0 (mult 2 .MatrixRows ) }}
%meta ioatt     bmi{{ $i }}	cp: main, index:{{ $i }}, type:input
%meta ioatt     bmi{{ $i }}	cp: bm, index:{{ $i }}, type:input
{{- end }}

{{- range $i := n 0 (mult 2 .MatrixRows ) }}
%meta ioatt     bmo{{ $i }}	cp: main, index:{{ $i }}, type:output
%meta ioatt     bmo{{ $i }}	cp: bm, index:{{ $i }}, type:output
{{- end }}

%meta bmdef	global registersize:32
`
)
