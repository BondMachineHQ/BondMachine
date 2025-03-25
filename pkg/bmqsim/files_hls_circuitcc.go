package bmqsim

const (
	HLSCircuitCC = `#include "circuit.h"

#include "ap_float.h"  // include the header for ap_float if necessary

void gate0(unsigned int N, float *vectorState) {
    #pragma HLS PIPELINE
    // Ensure that we have at least 4 elements in vectorState.
    if (N < 4) return;

    // Define a 4x4 matrix with example values.
    // You can replace these with any values you need.
    float matrix[4][4] = {
        {(float)1.0f, (float)2.0f, (float)3.0f, (float)4.0f},
        {(float)5.0f, (float)6.0f, (float)7.0f, (float)8.0f},
        {(float)9.0f, (float)10.0f, (float)11.0f, (float)12.0f},
        {(float)13.0f, (float)14.0f, (float)15.0f, (float)16.0f}
    };

    // Temporary array to store the result of the multiplication.
    float result[4] = {(float)0.0f, (float)0.0f, (float)0.0f, (float)0.0f};

    // Perform the matrix-vector multiplication.
    // For each row of the matrix, compute the dot product with the vector.
    for (int i = 0; i < 4; i++) {
        for (int j = 0; j < 4; j++) {
            result[i]+= matrix[i][j] * vectorState[j];
        }
    }

    // Copy the result back into the original vectorState array.
    for (int i = 0; i < 4; i++) {
        vectorState[i] = result[i];
    }
}

void circuit(unsigned int N, float *vectorState){
    #pragma hls interface mode=m_axi port=vectorState offset=slave bundle=gmem

    gate0(N, vectorState);
}

`
)
