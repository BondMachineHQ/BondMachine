package bondmachine

const (
	aximmPynqExample = `
	{
		"cells": [
		 {
		  "cell_type": "code",
		  "execution_count": null,
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
		  "execution_count": null,
		  "metadata": {},
		  "outputs": [],
		  "source": [
		   "# SETTINGS\n",
		   "project_name  = \"firmware\"\n",
		   "firmware_name = project_name+\".bit\"\n",
		   "n_input       = 4\n",
		   "n_output      = 2\n",
		   "benchcore     = True\n",
		   "precision_info = {\n",
		   "    \"type\": \"{{ .DataType }}\",\n",
		   "    \"e\": 0,\n",
		   "    \"f\": 0,\n",
		   "    \"s\": 0,\n",
		   "    \"t\": 0,\n",
		   "    \"host\": \"10.2.129.49\",\n",
		   "    \"port\": \"8080\"\n",
		   "}\n",
		   "padding = False\n",
		   "\n",
		   "if (precision_info[\"type\"][:3] == \"flp\"):\n",
		   "    exp = precision_info[\"type\"][4:precision_info[\"type\"].rindex(\"f\")]\n",
		   "    mant = precision_info[\"type\"][precision_info[\"type\"].rindex(\"f\")+1:len(precision_info[\"type\"])]\n",
		   "    precision_info[\"e\"] = int(exp)\n",
		   "    precision_info[\"f\"] = int(mant)\n",
		   "elif (precision_info[\"type\"][:3] == \"lqs\"):\n",
		   "    s = precision_info[\"type\"][precision_info[\"type\"].rindex(\"s\")+1:precision_info[\"type\"].rindex(\"t\")]\n",
		   "    t = precision_info[\"type\"][precision_info[\"type\"].rindex(\"t\")+1:len(precision_info[\"type\"])]\n",
		   "    precision_info[\"s\"] = s\n",
		   "    precision_info[\"t\"] = t\n",
		   "    \n",
		   "print(precision_info)\n",
		   "conversion_url ='http://'+precision_info[\"host\"]+':'+str(precision_info[\"port\"])+'/bmnumbers'"
		  ]
		 },
		 {
		  "cell_type": "code",
		  "execution_count": null,
		  "metadata": {},
		  "outputs": [],
		  "source": [
		   "def convert_flopoco_float_to_binary(num):\n",
		   "    \n",
		   "    print(\"float number to convert is: \", num)\n",
		   "    \n",
		   "    if (precision_info[\"type\"][:3] == \"flp\"):\n",
		   "        strNum = \"0flp<\"+str(precision_info[\"e\"])+\".\"+str(precision_info[\"f\"])+\">\"+str(num)\n",
		   "    elif (precision_info[\"type\"][:3] == \"lqs\"):\n",
		   "        strNum = \"0lq<\"+str(precision_info[\"s\"])+\".\"+str(precision_info[\"t\"])+\">\"+str(num)\n",
		   "        \n",
		   "    reqBody = {'action': 'cast', 'numbers': [strNum], 'reqType': 'bin', 'viewMode': 'native'}\n",
		   "    xReq = requests.post(conversion_url, json = reqBody)\n",
		   "    convertedNumber = json.loads(xReq.text)[\"numbers\"][0]\n",
		   "    print(convertedNumber)\n",
		   "    if (precision_info[\"type\"][:3] == \"flp\"):\n",
		   "        strNumber = convertedNumber[0]+convertedNumber[6:len(convertedNumber)]\n",
		   "    elif (precision_info[\"type\"][:3] == \"lqs\"):\n",
		   "        strNumber = convertedNumber[5:len(convertedNumber)]\n",
		   "    \n",
		   "        if len(strNumber) < int(precision_info[\"s\"]):\n",
		   "            diff = int(precision_info[\"s\"]) - len(strNumber)\n",
		   "            \n",
		   "            for i in range(0, diff):\n",
		   "                strNumber = \"0\"+strNumber\n",
		   "             \n",
		   "        if padding == True:\n",
		   "            n_zeros = 16 - len(strNumber)\n",
		   "            \n",
		   "            for i in range(0, n_zeros):\n",
		   "                strNumber = \"0\"+strNumber\n",
		   "    # print(\"binary number returned from  bmnumbers is: \", strNumber)\n",
		   "    # convert again from binary to float\n",
		   "    \n",
		   "    # converted_again = convert_flopoco_binary_to_float(strNumber)\n",
		   "    \n",
		   "    # print(\"float number from bmnumbers is: \", converted_again)\n",
		   "    # print(\"\\n\")\n",
		   "    \n",
		   "    return strNumber"
		  ]
		 },
		 {
		  "cell_type": "code",
		  "execution_count": null,
		  "metadata": {},
		  "outputs": [],
		  "source": [
		   "def convert_flopoco_binary_to_float(num):\n",
		   "    \n",
		   "    if padding == True:\n",
		   "        num = num[-8:]\n",
		   "    \n",
		   "    if (precision_info[\"type\"][:3] == \"flp\"):\n",
		   "        totLen = precision_info[\"e\"] + precision_info[\"f\"] + 3\n",
		   "    elif (precision_info[\"type\"][:3] == \"lqs\"):\n",
		   "        totLen = int(precision_info[\"s\"])\n",
		   "        \n",
		   "    if len(num) > totLen:\n",
		   "        num = num[1:]\n",
		   "    elif len(num) < totLen:\n",
		   "        diff_len = totLen - len(num)\n",
		   "        print(\"diff_len \", diff_len)\n",
		   "    \n",
		   "    strNum = \"0b<\"+str(totLen)+\">\"+str(num)\n",
		   "    \n",
		   "    print(\"string number\", strNum)\n",
		   "    if (precision_info[\"type\"][:3] == \"flp\"):\n",
		   "        newType = \"flpe\"+str(precision_info[\"e\"])+\"f\"+str(precision_info[\"f\"])\n",
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
		   "        print(spi0.read(offset, 4))\n",
		   "        bin_res = bin(spi0.read(offset, 4))\n",
		   "        print(\"binary received crudo from bm\", bin_res)\n",
		   "        if benchcore == True:\n",
		   "            if i != borderRange-1:\n",
		   "                if precision_info[\"type\"][:3] == \"flp\" or precision_info[\"type\"][:3] == \"lqs\":\n",
		   "                    output = convert_flopoco_binary_to_float(str(bin_res).replace(\"b\", \"\"))\n",
		   "                else:\n",
		   "                    output = bin_to_float(str(bin_res).replace(\"b\", \"\"))\n",
		   "            else:\n",
		   "                output = int(str(bin_res), 2)\n",
		   "        else:\n",
		   "            output = bin_to_float(str(bin_res).replace(\"b\", \"\"))\n",
		   "            \n",
		   "        result_from_bm_ml.append(output) # APPEND THE OUTPUT\n",
		   "        offset = offset + 4\n",
		   "    \n",
		   "    # print(result_from_bm_ml) # VECTOR OF PROBABIBILITES\n",
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
		   "X_test = np.load('banknote-authentication_X_test.npy')[:1]\n",
		   "Y_test = np.load('banknote-authentication_y_test.npy')[:1]\n",
		   "print(\" Example of first two input:  \", X_test[:1])\n",
		   "print(\" Example of first two output: \", Y_test[:1])"
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
		   "for xSample in X_test:\n",
		   "    offset = 0\n",
		   "    for feature in list(xSample): #\n",
		   "        if precision_info[\"type\"][:3] == \"flp\" or precision_info[\"type\"][:3] == \"lqs\":\n",
		   "            binToSend = convert_flopoco_float_to_binary(feature)\n",
		   "            print(\"binary to send for prediction: \", binToSend)\n",
		   "            print(\"binary to send for prediction: \", len(binToSend))\n",
		   "        else:\n",
		   "            binToSend = get_binary_from_float(feature)\n",
		   "        decToSend = int(binToSend, 2)\n",
		   "        spi0.write_mm(offset, decToSend) # WRITE THE FEATURE TO THE CORRESPONDING INPUT\n",
		   "        offset = offset + 4 # 4 BYTE = 32 BIT\n",
		   "    time.sleep(0.5)\n",
		   "    out = np.asarray(read_output())\n",
		   "    #print(\" #\",idx,\" -> classification: \", np.argmax(out[0:2]))\n",
		   "    classification = np.argmax(out[0:2])\n",
		   "    if (benchcore == True):\n",
		   "        results_to_dump.append([out[0], out[1], classification, out[2]])\n",
		   "    else: \n",
		   "        results_to_dump.append([out[0], out[1], classification])\n",
		   "    idx = idx + 1\n",
		   "    #break\n",
		   "\n",
		   "# 01001010100\n",
		   "print(results_to_dump)"
		  ]
		 },
		 {
		  "cell_type": "code",
		  "execution_count": null,
		  "metadata": {},
		  "outputs": [],
		  "source": [
		   "import csv\n",
		   "fields = ['probability_0', 'probability_1', 'classification', 'clock_cycles'] \n",
		   "\n",
		   "with open(project_name+\".csv\", 'w') as f:\n",
		   "    write = csv.writer(f)\n",
		   "    write.writerow(fields)\n",
		   "    write.writerows(results_to_dump)"
		  ]
		 },
		 {
		  "cell_type": "code",
		  "execution_count": null,
		  "metadata": {},
		  "outputs": [],
		  "source": []
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
