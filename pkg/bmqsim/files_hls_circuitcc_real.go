package bmqsim

const (
	HLSCircuitCCReal = `#include "circuit.h"

#include "ap_float.h"  // include the header for ap_float if necessary

{{- range $k := n 0 $.NumGates }}
void gate{{ $k }}(unsigned int N, float *vectorState) {
    #pragma HLS PIPELINE
    // Ensure that we have at least 4 elements in vectorState.
    if (N < {{ $.MatrixRows }}) return;

    // Define a matrix with gate values.
    float matrix[{{ $.MatrixRows }}][{{ $.MatrixRows }}] = {
{{- range $i := n 0 $.MatrixRows }}
        { {{- range $j := n 0 $.MatrixRows }}{{- if ne $j 0 }}, {{- end }}(float){{ index (index $.MtxReal $k) $i $j }}{{- end }} }{{- if ne $i (dec $.MatrixRows) }},{{- end }}
{{- end }}
    };

    // Temporary array to store the result of the multiplication.
    float result[{{ $.MatrixRows }}] = { {{- range $i := n 0 $.MatrixRows }}{{- if ne $i 0 }},{{- end }} (float)0.0{{- end }} };

    // Perform the matrix-vector multiplication.
    // For each row of the matrix, compute the dot product with the vector.
    for (int i = 0; i < {{ $.MatrixRows }}; i++) {
        for (int j = 0; j < {{ $.MatrixRows }}; j++) {
            result[i]+= matrix[i][j] * vectorState[j];
        }
    }

    // Copy the result back into the original vectorState array.
    for (int i = 0; i < {{ $.MatrixRows }}; i++) {
        vectorState[i] = result[i];
    }
}
{{- end }}

void circuit(unsigned int N, float *vectorState) {
    #pragma hls interface mode=m_axi port=vectorState offset=slave bundle=gmem

{{- range $k := n 0 $.NumGates }}
    gate{{ $k }}(N, vectorState);
{{- end }}
}

`
)
