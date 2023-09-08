package bondmachine

const (
	aximmPynqExample = `
	{
		"cells": [
		 {
		  "cell_type": "code",
		  "execution_count": 107,
		  "metadata": {},
		  "outputs": [],
		  "source": [
		   "from pynq import Overlay\n",
		   "from pynq import MMIO\n",
		   "import os\n",
		   "import numpy as np\n",
		   "import struct\n",
		   "import time\n",
		   "import requests\n",
		   "import json"
		  ]
		 },
		 {
		  "cell_type": "code",
		  "execution_count": 108,
		  "metadata": {},
		  "outputs": [],
		  "source": [
		   "# PERSONALIZE HERE\n",
		   "use_preload_dataset = True\n",
		   "n_input       = 4 # number of input (features) to send for prediction\n",
		   "n_output      = 2 # number of outputs (without benchcore)\n",
		   "benchcore     = True # set to True if the BondMachine has the benchcore enabled\n",
		   "data_type     = \"{{ .DataType }}\"\n",
		   "reg_size      = {{ $.Rsize }}"
		  ]
		 },
		 {
			"cell_type": "code",
			"execution_count": null,
			"metadata": {},
			"outputs": [],
			"source": [
			 "# SETTINGS\n",
			 "firmware_name = \"firmware.bit\"\n",
			 "\n",
			 "dataset_info = {\n",
			 "    \"preloaded\": use_preload_dataset, # set to True if you want to use a dataset, if set to False the dataset is generated randomly\n",
			 "    \"x_test\": \"banknote-authentication_X_test.npy\",\n",
			 "    \"y_test\": \"banknote-authentication_y_test.npy\",\n",
			 "    \"features\": n_input\n",
			 "}\n",
			 "\n",
			 "precision_info = {\n",
			 "    \"type\": data_type,\n",
			 "    \"e\": 0,\n",
			 "    \"f\": 0,\n",
			 "    \"s\": 0,\n",
			 "    \"t\": 0,\n",
			 "    \"host\": \"10.2.129.49\",\n",
			 "    \"port\": \"80\"\n",
			 "}\n",
			 "padding_info = {\n",
			 "    \"enabled\": False,\n",
			 "    \"size\": reg_size\n",
			 "}\n",
			 "\n",
			 "if (precision_info[\"type\"][:3] == \"flp\"):\n",
			 "    exp = precision_info[\"type\"][4:precision_info[\"type\"].rindex(\"f\")]\n",
			 "    mant = precision_info[\"type\"][precision_info[\"type\"].rindex(\"f\")+1:len(precision_info[\"type\"])]\n",
			 "    precision_info[\"e\"] = int(exp)\n",
			 "    precision_info[\"f\"] = int(mant)\n",
			 "elif (precision_info[\"type\"][:3] == \"fps\"):\n",
			 "    s = precision_info[\"type\"][precision_info[\"type\"].rindex(\"s\")+1:precision_info[\"type\"].rindex(\"f\")]\n",
			 "    f = precision_info[\"type\"][precision_info[\"type\"].rindex(\"f\")+1:len(precision_info[\"type\"])]\n",
			 "    precision_info[\"s\"] = s\n",
			 "    precision_info[\"f\"] = f\n",
			 "elif (precision_info[\"type\"][:3] == \"lqs\"):\n",
			 "    s = precision_info[\"type\"][precision_info[\"type\"].rindex(\"s\")+1:precision_info[\"type\"].rindex(\"t\")]\n",
			 "    t = precision_info[\"type\"][precision_info[\"type\"].rindex(\"t\")+1:len(precision_info[\"type\"])]\n",
			 "    precision_info[\"s\"] = s\n",
			 "    precision_info[\"t\"] = t\n",
			 "    \n",
			 "conversion_url ='http://'+precision_info[\"host\"]+':'+str(precision_info[\"port\"])+'/bmnumbers'"
			]
		   },
		   {
			"cell_type": "code",
			"execution_count": null,
			"metadata": {},
			"outputs": [],
			"source": [
			 "def cast_float_to_binary(num):\n",
			 "    \n",
			 "    if (precision_info[\"type\"][:3] == \"flp\"):\n",
			 "        strNum = \"0flp<\"+str(precision_info[\"e\"])+\".\"+str(precision_info[\"f\"])+\">\"+str(num)\n",
			 "    elif (precision_info[\"type\"][:3] == \"fps\"):\n",
    		 "        strNum = \"0fp<\"+str(precision_info[\"s\"])+\".\"+str(precision_info[\"f\"])+\">\"+str(num)\n",
			 "    elif (precision_info[\"type\"][:3] == \"lqs\"):\n",
			 "        strNum = \"0lq<\"+str(precision_info[\"s\"])+\".\"+str(precision_info[\"t\"])+\">\"+str(num)\n",
			 "        \n",
			 "    reqBody = {'action': 'cast', 'numbers': [strNum], 'reqType': 'bin', 'viewMode': 'native'}\n",
			 "    xReq = requests.post(conversion_url, json = reqBody)\n",
			 "    convertedNumber = json.loads(xReq.text)[\"numbers\"][0]\n",
			 "    \n",
			 "    if (precision_info[\"type\"][:3] == \"flp\"):\n",
			 "        strNumber = convertedNumber[0]+convertedNumber[convertedNumber.rindex(\">\")+1:len(convertedNumber)]\n",
			 "    elif (precision_info[\"type\"][:3] == \"lqs\" or precision_info[\"type\"][:3] == \"fps\"):\n",
			 "        strNumber = convertedNumber[convertedNumber.rindex(\">\")+1:len(convertedNumber)]\n",
			 "        \n",
			 "        if len(strNumber) < int(precision_info[\"s\"]):\n",
			 "            diff = int(precision_info[\"s\"]) - len(strNumber)\n",
			 "            \n",
			 "            for i in range(0, diff):\n",
			 "                strNumber = \"0\"+strNumber\n",
			 "             \n",
			 "        if padding_info[\"enabled\"] == True:\n",
			 "            n_zeros = padding_info[\"size\"] - len(strNumber)\n",
			 "            \n",
			 "            for i in range(0, n_zeros):\n",
			 "                strNumber = \"0\"+strNumber\n",
			 "                \n",
			 "            padding_info[\"n_zeros\"] = n_zeros\n",
			 "                \n",
			 "    return strNumber"
			]
		   },
		   {
			"cell_type": "code",
			"execution_count": null,
			"metadata": {},
			"outputs": [],
			"source": [
			 "def cast_binary_to_float(num):\n",
			 "    \n",
			 "    if padding_info[\"enabled\"] == True:\n",
			 "        num = num[-padding_info[\"n_zeros\"]:]\n",
			 "    \n",
			 "    if (precision_info[\"type\"][:3] == \"flp\"):\n",
			 "        totLen = precision_info[\"e\"] + precision_info[\"f\"] + 3\n",
			 "    elif (precision_info[\"type\"][:3] == \"fps\"):\n",
    		 "        totLen = int(precision_info[\"s\"])\n",
			 "    elif (precision_info[\"type\"][:3] == \"lqs\"):\n",
			 "        totLen = int(precision_info[\"s\"])\n",
			 "        \n",
			 "    if len(num) > totLen:\n",
			 "        num = num[1:]\n",
			 "    elif len(num) < totLen:\n",
			 "        diff_len = totLen - len(num)\n",
			 "    \n",
			 "    strNum = \"0b<\"+str(totLen)+\">\"+str(num)\n",
			 "    \n",
			 "    if (precision_info[\"type\"][:3] == \"flp\"):\n",
			 "        newType = \"flpe\"+str(precision_info[\"e\"])+\"f\"+str(precision_info[\"f\"])\n",
			 "    elif (precision_info[\"type\"][:3] == \"fps\"):\n",
    		 "        newType = \"fps\"+str(precision_info[\"s\"])+\"f\"+str(precision_info[\"f\"])\n",
			 "    elif (precision_info[\"type\"][:3] == \"lqs\"):\n",
			 "        newType = \"lqs\"+str(precision_info[\"s\"])+\"t\"+str(precision_info[\"t\"])\n",
			 "    \n",
			 "    reqBody = {'action': 'cast', 'numbers': [strNum], 'reqType': newType, 'viewMode': 'native'}\n",
			 "    xReq = requests.post(conversion_url, json = reqBody)\n",
			 "    try:\n",
			 "        convertedNumber = json.loads(xReq.text)[\"numbers\"][0]\n",
			 "        strNumber = convertedNumber[convertedNumber.rindex(\">\")+1:len(convertedNumber)]\n",
			 "        return float(strNumber)\n",
			 "    except Exception as e:\n",
			 "        print(e)"
			]
		   },
		   {
			"cell_type": "code",
			"execution_count": null,
			"metadata": {},
			"outputs": [],
			"source": [
			 "def get_binary_from_float(num):\n",
			 "    return bin(struct.unpack('!i',struct.pack('!f',num))[0])"
			]
		   },
		   {
			"cell_type": "code",
			"execution_count": null,
			"metadata": {},
			"outputs": [],
			"source": [
			 "def bin_to_float(binary):\n",
			 "    return struct.unpack('!f',struct.pack('!I', int(binary, 2)))[0]"
			]
		   },
		   {
			"cell_type": "code",
			"execution_count": null,
			"metadata": {},
			"outputs": [],
			"source": [
			 "def read_output():\n",
			 "    starting_offset = (n_input*4)+(2*4)\n",
			 "    result_from_bm_ml = []\n",
			 "    offset = starting_offset\n",
			 "    borderRange = n_output + 1 if benchcore == True else n_output\n",
			 "    for i in range(0, borderRange):\n",
			 "        raw_res = spi0.read(offset, 4)\n",
			 "        #print(\" * RAW    OUT:  \", raw_res)\n",
			 "        bin_res = bin(raw_res)\n",
			 "        bin_res_str = str(bin_res).replace(\"b\", \"\")\n",
			 "        #print(\" * BINARY  OUT: \", binToSend)\n",
			 "        if benchcore == True:\n",
			 "            if i != borderRange-1:\n",
			 "                if precision_info[\"type\"][:3] == \"flp\" or precision_info[\"type\"][:3] == \"lqs\" or precision_info[\"type\"][:3] == \"fps\":\n",
			 "                    output = cast_binary_to_float(bin_res_str)\n",
			 "                else:\n",
			 "                    output = bin_to_float(bin_res_str)\n",
			 "            else:\n",
			 "                output = int(str(bin_res), 2)\n",
			 "        else:\n",
			 "            if precision_info[\"type\"][:3] == \"flp\" or precision_info[\"type\"][:3] == \"lqs\" or precision_info[\"type\"][:3] == \"fps\":\n",
			 "                output = cast_binary_to_float(bin_res_str)\n",
			 "            else:\n",
			 "                output = bin_to_float(bin_res_str)\n",
			 "            \n",
			 "        result_from_bm_ml.append(output) # APPEND THE OUTPUT\n",
			 "        offset = offset + 4\n",
			 "    \n",
			 "    return result_from_bm_ml"
			]
		   },
		   {
			"cell_type": "code",
			"execution_count": null,
			"metadata": {},
			"outputs": [],
			"source": [
			 "# LOAD OVERLAY\n",
			 "overlay = Overlay(os.getcwd()+\"/\"+firmware_name)"
			]
		   },
		   {
			"cell_type": "code",
			"execution_count": null,
			"metadata": {},
			"outputs": [],
			"source": [
			 "# GET MEMORY ADDRESS OF IP\n",
			 "bm_starting_address = (overlay.ip_dict[\"bondmachineip_0\"][\"phys_addr\"])\n",
			 "print(\" Starting memory address of Bondmachine IP is (in dec): \", bm_starting_address)\n",
			 "print(\" Starting memory address of Bondmachine IP is (in hex): \", hex(bm_starting_address))"
			]
		   },
		   {
			"cell_type": "code",
			"execution_count": null,
			"metadata": {},
			"outputs": [],
			"source": [
			 "# GET THE OBJECT NECESSARY TO INTERACT WITH ML IP\n",
			 "spi0 = MMIO(bm_starting_address, 128)"
			]
		   },
		   {
			"cell_type": "code",
			"execution_count": null,
			"metadata": {},
			"outputs": [],
			"source": [
			 "# LOAD BANKNOTE TESTSET\n",
			 "if dataset_info[\"preloaded\"] == True:\n",
			 "    X_test = np.load(dataset_info[\"x_test\"])\n",
			 "    Y_test = np.load(dataset_info[\"y_test\"])\n",
			 "    X_test = X_test.reshape(-1)\n",
			 "    Y_test = X_test.reshape(-1)\n",
			 "else:\n",
			 "    X_test = np.random.uniform(-20, 20, size=10).reshape(-1)\n",
			 "    Y_test = np.random.rand(2).reshape(-1)\n",
			 "#X_test"
			]
		   },
		   {
			"cell_type": "code",
			"execution_count": null,
			"metadata": {},
			"outputs": [],
			"source": [
			 "# IN THIS CASE I WANT TO SEND ONLY THE FIRST X INPUT SAMPLE\n",
			 "idx = 0\n",
			 "results_to_dump = []\n",
			 "cnt = 0\n",
			 "offset = 0\n",
			 "\n",
			 "for f in range(0,len(X_test)):\n",
			 "    if precision_info[\"type\"][:3] == \"flp\" or precision_info[\"type\"][:3] == \"lqs\" or precision_info[\"type\"][:3] == \"fps\":\n",
			 "        binToSend = cast_float_to_binary(X_test[f])\n",
			 "    else:\n",
			 "        binToSend = get_binary_from_float(X_test[f])\n",
			 "        \n",
			 "    #print(\" * FLOAT   IN: \", X_test[f])\n",
			 "    #print(\" * BINARY  IN: \", binToSend)\n",
			 "    \n",
			 "    decToSend = int(binToSend, 2)\n",
			 "    spi0.write_mm(offset, decToSend) # WRITE THE FEATURE TO THE CORRESPONDING INPUT\n",
			 "    offset = offset + 4 # 4 BYTE = 32 BIT\n",
			 "    cnt = cnt + 1\n",
			 "    \n",
			 "    if dataset_info[\"preloaded\"] == True:\n",
			 "        if cnt == dataset_info[\"features\"]:\n",
			 "            time.sleep(1)\n",
			 "            out = np.asarray(read_output())\n",
			 "            \n",
			 "            if (benchcore == True):\n",
			 "                classification = np.argmax(out[0:2])\n",
			 "                results_to_dump.append([out[0], out[1], classification, out[2]])\n",
			 "            else:\n",
			 "                for i in range(0, n_output):\n",
			 "                    results_to_dump.append(out[i])\n",
			 "                    \n",
			 "                classification = np.argmax(out[0:n_output])\n",
			 "                results_to_dump.append(classification)\n",
			 "            idx = idx + 1\n",
			 "            cnt = 0\n",
			 "            offset = 0\n",
			 "            print(results_to_dump)\n",
			 "    else:\n",
			 "        if cnt == n_input:\n",
			 "            offset = 0\n",
			 "            cnt = 0\n",
			 "            time.sleep(1)\n",
			 "            out = np.asarray(read_output())\n",
			 "            print(\" * FLOAT OUT: \", out)\n",
			 "            print(\" *************************************************** \")"
			]
		   },
		   {
			"cell_type": "code",
			"execution_count": null,
			"metadata": {},
			"outputs": [],
			"source": [
			 "if use_preload_dataset == True:\n",
			 "    import csv\n",
			 "    \n",
			 "    fields = []\n",
			 "    for i in range(0, n_output):\n",
			 "        fields.append('probability_'+str(i))\n",
			 "        \n",
			 "    fields.append('classification')\n",
			 "    if benchcore == True:\n",
			 "        fields.append('clock_cycles')\n",
			 "    \n",
			 "    with open(\"firmware.csv\", 'w') as f:\n",
			 "        write = csv.writer(f)\n",
			 "        write.writerow(fields)\n",
			 "        write.writerows(results_to_dump)"
			]
		   },
		   {
			"cell_type": "code",
			"execution_count": null,
			"metadata": {},
			"outputs": [],
			"source": []
		   }
		  ],
		  "metadata": {
		   "kernelspec": {
			"display_name": "Python 3",
			"language": "python",
			"name": "python3"
		   },
		   "language_info": {
			"codemirror_mode": {
			 "name": "ipython",
			 "version": 3
			},
			"file_extension": ".py",
			"mimetype": "text/x-python",
			"name": "python",
			"nbconvert_exporter": "python",
			"pygments_lexer": "ipython3",
			"version": "3.6.5"
		   }
		  },
		  "nbformat": 4,
		  "nbformat_minor": 2
		 }
			  
	`
)
