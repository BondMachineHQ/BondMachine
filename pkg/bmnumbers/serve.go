// TRY WITH CURL
// curl -X POST -d '{"action": "cast", "numbers": ["0f0.45", "0f0.30"], "reqType": "bin", "dumpMode": "native"}' http://127.0.0.1:8080/bmnumbers

package bmnumbers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type BMNumbersRequest struct {
	Action   string   `json:"action"`
	Numbers  []string `json:"numbers"`
	ReqType  string   `json:"reqType"`
	DumpMode string   `json:"dumpMode"`
}

type BMNumbersResponse struct {
	Numbers []string `json:"numbers"`
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the BM Numbers server!\n")
	fmt.Fprintf(w, "There is one endpoint available: \n")
	fmt.Fprintf(w, "[POST] bmnumbers \n")
	fmt.Fprintf(w, "[POST] body example: {'action': 'cast', 'numbers': ['0f0.45', '0f0.30'], 'reqType': 'bin', 'dumpMode': 'native'} \n")
}

func ExecRequest(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var bmNumbersRequest BMNumbersRequest
	err = json.Unmarshal(reqBody, &bmNumbersRequest)
	if err != nil {
		http.Error(w, "Failed to unmarshal request body", http.StatusBadRequest)
		return
	}

	var newType BMNumberType

	if v := GetType(bmNumbersRequest.ReqType); v == nil {
		http.Error(w, "Unknown type", http.StatusInternalServerError)
		return
	} else {
		newType = v
	}

	var results []string

	for i := 0; i < len(bmNumbersRequest.Numbers); i++ {
		if output, err := ImportString(bmNumbersRequest.Numbers[i]); err != nil {
			http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
			return
		} else {

			if bmNumbersRequest.Action == "convert" {
				if err := newType.Convert(output); err != nil {
					http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
					return
				}
			} else if bmNumbersRequest.Action == "cast" {
				if err := CastType(output, newType); err != nil {
					http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
					return
				}
			}

			switch bmNumbersRequest.DumpMode {
			case "native":
				if value, err := output.ExportString(); err != nil {
					http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
					return
				} else {
					results = append(results, value) // add a string to the slice
				}
			case "bin":
				if value, err := output.ExportBinary(false); err != nil {
					http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
					return
				} else {
					results = append(results, value) // add a string to the slice
				}
			case "unsigned":
				if value, err := output.ExportUint64(); err != nil {
					http.Error(w, fmt.Sprintf("Error: %v", err), http.StatusInternalServerError)
					return
				} else {
					results = append(results, strconv.FormatUint(value, 10)) // add a string to the slice
				}
			default:
				http.Error(w, "Unknown dump format", http.StatusInternalServerError)
				return
			}
		}
	}

	w.WriteHeader(http.StatusOK) // set HTTP status code to 200
	json.NewEncoder(w).Encode(BMNumbersResponse{Numbers: results})
}

// REST API to convert numbers
func Serve() {

	http.HandleFunc("/", HomePage)
	http.HandleFunc("/bmnumbers", ExecRequest)

	fmt.Println("Server started on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
