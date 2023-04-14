package bmnumbers

func testDataSet(dataSet string) interface{} {

	dataSetUint := make(map[string]uint64)
	dataSetBin := make(map[string]string)
	dataSetBinNoSize := make(map[string]string)

	dataSetBin["0flp<4.4>1.45"] = "0b<11>1001110111"
	dataSetBin["0flp<5.7>-4.3"] = "0b<15>11100010001010"

	dataSetUint["56"] = 56
	dataSetBin["56"] = "0b<64>111000"
	dataSetBinNoSize["56"] = "111000"

	dataSetUint["0x901"] = 2305
	dataSetBin["0x901"] = "0b<16>100100000001"
	dataSetBinNoSize["0x901"] = "100100000001"

	dataSetUint["0b10101"] = 21
	dataSetBin["0b10101"] = "0b<5>10101"
	dataSetBinNoSize["0b10101"] = "10101"

	dataSetUint["0d56"] = 56
	dataSetBin["0d56"] = "0b<64>111000"
	dataSetBinNoSize["0d56"] = "111000"

	dataSetUint["0f56"] = 1113587712
	dataSetBin["0f56"] = "0b<32>1000010011000000000000000000000"
	dataSetBinNoSize["0f56"] = "1000010011000000000000000000000"

	dataSetUint["0finfinity"] = 2139095040
	dataSetBin["0finfinity"] = "0b<32>1111111100000000000000000000000"
	dataSetBinNoSize["0finfinity"] = "1111111100000000000000000000000"

	dataSetUint["0fNaN"] = 2143289344
	dataSetBin["0fNaN"] = "0b<32>1111111110000000000000000000000"
	dataSetBinNoSize["0fNaN"] = "1111111110000000000000000000000"

	dataSetUint["0f4e-4"] = 970045207
	dataSetBin["0f4e-4"] = "0b<32>111001110100011011011100010111"
	dataSetBinNoSize["0f4e-4"] = "111001110100011011011100010111"

	switch dataSet {
	case "uintData":
		return dataSetUint
	case "binData":
		return dataSetBin
	case "binDataNoSize":
		return dataSetBinNoSize
	}

	return nil
}
