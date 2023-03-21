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
		   "import time"
		  ]
		 },
		 {
		  "cell_type": "code",
		  "execution_count": null,
		  "metadata": {},
		  "outputs": [],
		  "source": [
		   "# SETTINGS\n",
		   "project_name  = os.getcwd()+\"firmware\"\n",
		   "firmware_name = project_name+\".bit\"\n",
		   "n_input       = 4\n",
		   "n_output      = 2\n",
		   "benchcore     = True"
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
		   "        bin_res = bin(spi0.read(offset, 4))\n",
		   "        \n",
		   "        if benchcore == True:\n",
		   "            if i != borderRange-1:\n",
		   "                output = bin_to_float(str(bin_res).replace(\"b\", \"\"))\n",
		   "            else:\n",
		   "                output = int(str(bin_res), 2)\n",
		   "        else:\n",
		   "            output = bin_to_float(str(bin_res).replace(\"b\", \"\"))\n",
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
		   "X_test = np.load('banknote-authentication_X_test.npy')\n",
		   "Y_test = np.load('banknote-authentication_y_test.npy')\n",
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
		   "        binToSend = get_binary_from_float(feature)\n",
		   "        decToSend = int(binToSend, 2)\n",
		   "        spi0.write_mm(offset, decToSend) # WRITE THE FEATURE TO THE CORRESPONDING INPUT\n",
		   "        offset = offset + 4 # 4 BYTE = 32 BIT\n",
		   "    time.sleep(0.5)\n",
		   "    out = np.asarray(read_output())\n",
		   "    classification = np.argmax(out[0:2])\n",
		   "    if (benchcore == True):\n",
		   "        results_to_dump.append([out[0], out[1], classification, out[2]])\n",
		   "    else: \n",
		   "        results_to_dump.append([out[0], out[1], classification])\n",
		   "    idx = idx + 1"
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
