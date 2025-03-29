package bmqsim

const (
	HLSCircuitCCComplex = `#include "circuit.h"

#include "ap_float.h"  // include the header for ap_float if necessary

#define vectorStateReal(i) vectorState[i*2] 
#define vectorStateImag(i) vectorState[i*2 + 1] 

{{- range $k := n 0 $.NumGates }}
void gate{{ $k }}(unsigned int N, float *vectorState) {
    #pragma HLS PIPELINE
    // Ensure that we have at least 4 elements in vectorState.
    if (N < {{ mult 2 $.MatrixRows }}) return;

    // Define a matrix with gate values.
    float matrixReal[{{ $.MatrixRows }}][{{ $.MatrixRows }}] = {
{{- range $i := n 0 $.MatrixRows }}
        { {{- range $j := n 0 $.MatrixRows }}{{- if ne $j 0 }}, {{- end }}(float){{ index (index $.MtxReal $k) $i $j }}{{- end }} }{{- if ne $i (dec $.MatrixRows) }},{{- end }}
{{- end }}
    };
    float matrixImag[{{ $.MatrixRows }}][{{ $.MatrixRows }}] = {
{{- range $i := n 0 $.MatrixRows }}
        { {{- range $j := n 0 $.MatrixRows }}{{- if ne $j 0 }}, {{- end }}(float){{ index (index $.MtxImag $k) $i $j }}{{- end }} }{{- if ne $i (dec $.MatrixRows) }},{{- end }}
{{- end }}
    };

    // Temporary array to store the result of the multiplication.
    float resultReal[{{ $.MatrixRows }}] = { {{- range $i := n 0 $.MatrixRows }}{{- if ne $i 0 }},{{- end }} (float)0.0{{- end }} };
    float resultImag[{{ $.MatrixRows }}] = { {{- range $i := n 0 $.MatrixRows }}{{- if ne $i 0 }},{{- end }} (float)0.0{{- end }} };

    // Perform the matrix-vector multiplication.
    // For each row of the matrix, compute the dot product with the vector.
    for (int i = 0; i < {{ $.MatrixRows }}; i++) {
        for (int j = 0; j < {{ $.MatrixRows }}; j++) {
            resultReal[i] += matrixReal[i][j] * vectorStateReal(j) - matrixImag[i][j] * vectorStateImag(j);
            resultImag[i] += matrixReal[i][j] * vectorStateImag(j) + matrixImag[i][j] * vectorStateReal(j);
        }
    }

    // Copy the result back into the original vectorState array.
    for (int i = 0; i < {{ mult 2 $.MatrixRows }}; i++) {
        if (i%2 == 0) {
            vectorState[i] = resultReal[i/2];
        } else {
            vectorState[i] = resultImag[i/2];
        }
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
