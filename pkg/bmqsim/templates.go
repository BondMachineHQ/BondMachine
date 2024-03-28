package bmqsim

import (
	"bytes"
	"fmt"
	"math"
	"text/template"
)

// Modes:
// seq_hardcoded
// seq
// full_hw_hardcoded

var HardwareFlavors = map[string]string{
	"seq_hardcoded_real":            SeqHardcodedReal,
	"seq_hardcoded_complex":         SeqHardcodedComplex,
	"seq_hardcoded_addtree_complex": SeqHardcodedAddTreeComplex,
}

var HardwareFlavorsTags = map[string][]string{
	"seq_hardcoded_real":            {"real"},
	"seq_hardcoded_complex":         {"complex"},
	"seq_hardcoded_addtree_complex": {"complex"},
}

var AppFlavors = map[string]string{
	"python_pynq_real":    PythonPynqReal,
	"python_pynq_complex": PythonPynqComplex,
	"c_pynqapi_real":      CPynqApiReal,
	"c_pynqapi_complex":   CPynqApiComplex,
	"cpp_opencl_real":     CppOpenCLReal,
	"cpp_opencl_complex":  CppOpenCLComplex,
}

var AppFlavorsTags = map[string][]string{
	"python_pynq_real":    {"real"},
	"python_pynq_complex": {"complex"},
	"c_pynqapi_real":      {"real"},
	"c_pynqapi_complex":   {"complex"},
	"cpp_opencl_real":     {"real"},
	"cpp_opencl_complex":  {"complex"},
}

type templateData struct {
	Qbits      int
	NumGates   int
	MatrixRows int
	MtxReal    [][][]float32
	MtxImag    [][][]float32
	funcMap    template.FuncMap
}

func (sim *BmQSimulator) createBasicTemplateData() *templateData {
	result := new(templateData)
	result.Qbits = len(sim.qbits)
	result.MatrixRows = int(math.Pow(float64(2), float64(len(sim.qbits))))
	result.NumGates = len(sim.Mtx)
	result.MtxReal = make([][][]float32, result.NumGates)
	result.MtxImag = make([][][]float32, result.NumGates)
	for g, m := range sim.Mtx {
		result.MtxReal[g] = make([][]float32, m.N)
		result.MtxImag[g] = make([][]float32, m.N)
		for i := 0; i < m.N; i++ {
			result.MtxReal[g][i] = make([]float32, m.N)
			result.MtxImag[g][i] = make([]float32, m.N)
			for j := 0; j < m.N; j++ {
				result.MtxReal[g][i][j] = m.Data[i][j].Real
				result.MtxImag[g][i][j] = m.Data[i][j].Imag
			}
		}
	}

	funcMap := template.FuncMap{
		"inc": func(i int) int {
			return i + 1
		},
		"dec": func(i int) int {
			return i - 1
		},
		"n": func(start, end int) []int {
			var result []int
			for i := start; i < end; i++ {
				result = append(result, i)
			}
			return result
		},
		"ns": func(start, end, step int) []int {
			var result []int
			for i := start; i < end; i += step {
				result = append(result, i)
			}
			return result
		},
		"sum": func(a, b int) int {
			return a + b
		},
		"div": func(a, b int) int {
			return a / b
		},
		"mult": func(a, b int) int {
			return a * b
		},
	}
	result.funcMap = funcMap
	return result
}

func (sim *BmQSimulator) ApplyTemplate(mode string) (string, error) {
	var data string
	if d, ok := HardwareFlavors[mode]; ok {
		data = d
	}
	if d, ok := AppFlavors[mode]; ok {
		data = d
	}

	templateData := sim.createBasicTemplateData()
	t, err := template.New(mode).Funcs(templateData.funcMap).Parse(data)
	if err != nil {
		return "", err
	}
	var f bytes.Buffer
	err = t.Execute(&f, *templateData)
	if err != nil {
		return "", err
	}
	return f.String(), nil
}

func (sim *BmQSimulator) VerifyConditions(mode string) error {
	switch mode {
	case "seq_hardcoded_real":
		// Check if there are complex numbers within the matrices
		for _, m := range sim.Mtx {
			for i := 0; i < m.N; i++ {
				for j := 0; j < m.N; j++ {
					if m.Data[i][j].Imag != 0 {
						return fmt.Errorf("complex numbers are not supported in this mode")
					}
				}
			}
		}
	case "seq_hardcoded_complex":
	}
	return nil
}
