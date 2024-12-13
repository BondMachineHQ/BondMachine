package bmmatrix

const (
	templateMult = `
{{- range $i := n 0 (rows .Mtx1) }}
{{- range $j := n 0 (cols $.Mtx2) }}
{{- range $k := n 0 (cols $.Mtx1) }}
%section mult_{{ $i }}_{{ $j }}_{{ $k }} .romtext iomode:{{ $.Iomode }}
	entry _start    ; Entry point
_start:
	mov	r1, {{"{{"}} .Params.typeprefix {{"}}"}}{{ index $.Mtx1 $i $k }}
loop:
	mov	r0, i0
	{{"{{"}} .Params.multop {{"}}"}}	r0, r1
	mov	o0, r0
	j loop
%endsection

%meta cpdef	mult{{ $i }}_{{ $j }}_{{ $k }}	romcode: mult_{{ $i }}_{{ $j }}_{{ $k }}, execmode: ha, typeprefix:0f, multop:multf

%meta iodef     toadd{{ $i }}_{{ $j }}_{{ $k }} type:io
%meta ioatt	toadd{{ $i }}_{{ $j }}_{{ $k }} cp:mult{{ $i }}_{{ $j }}_{{ $k }}, type:output, index:0
%meta ioatt	toadd{{ $i }}_{{ $j }}_{{ $k }} cp:add{{ $i }}_{{ $j }}, type:input, index:  {{ $k }}


%meta iodef     fromin{{ $i }}_{{ $j }}_{{ $k }} type:io
%meta ioatt	fromin{{ $i }}_{{ $j }}_{{ $k }} cp:mult{{ $i }}_{{ $j }}_{{ $k }}, type:input, index:0
%meta ioatt	fromin{{ $i }}_{{ $j }}_{{ $k }} cp:bm, type:input, index: {{ sum (mult $k (cols $.Mtx2) ) $j }}

{{ end }}

%section add_{{ $i }}_{{ $j }} .romtext iomode:{{ $.Iomode }}
	entry _start    ; Entry point
_start:
	clr	r1
{{- range $k := n 0 (cols $.Mtx1) }}
	mov	r0, i{{ $k }}
	{{"{{"}} .Params.addop {{"}}"}}	r1, r0
{{- end }}
	mov	o0, r1
	j _start
%endsection

%meta cpdef	add{{ $i }}_{{ $j }}	romcode: add_{{ $i }}_{{ $j }}, execmode: ha, addop:addf, typeprefix:0f

%meta iodef     toout{{ $i }}_{{ $j }} type:io
%meta ioatt	toout{{ $i }}_{{ $j }} cp:add{{ $i }}_{{ $j }}, type:output, index:0
%meta ioatt	toout{{ $i }}_{{ $j }} cp:bm, type:output, index: {{ sum (mult $i (cols $.Mtx2) ) $j }}
{{ end }}
{{- end }}

%meta bmdef	global registersize:32
`
)
