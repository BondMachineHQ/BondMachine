package bondmachine

const (
	auxfilesAXIOutRegs = `
    {{- $smindex:= 0 }}
    {{- if .Inputs }}
    {{- range .Inputs }}
    {{- $smindex = inc $smindex }}
    {{- end }}
    {{- end }}
    {{- $smindex = inc $smindex }}
    {{- $smindex = inc $smindex }}
    {{- if .Outputs }}
    {{- range .Outputs }}
    slv_reg{{ $smindex }}
    {{- $smindex = inc $smindex }}
    {{- end }}
    {{- end }}   
    slv_reg{{ $smindex }}
    {{- $smindex = inc $smindex }}
    slv_reg{{ $smindex }}
`
)
