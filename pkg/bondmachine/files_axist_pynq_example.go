package bondmachine

const (
	axistPynqExample = `
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
       "import time\n",
       "import sys\n",
       "import csv\n",
       "from datetime import datetime\n",
       "import statistics\n",
       "\n",
       "np.set_printoptions(threshold=sys.maxsize)"
      ]
     },
     {
      "cell_type": "code",
      "execution_count": null,
      "metadata": {},
      "outputs": [],
      "source": [
       "overlay = Overlay(os.getcwd()+\"/firmware.bit\")"
      ]
     },
     {
      "cell_type": "code",
      "execution_count": null,
      "metadata": {},
      "outputs": [],
      "source": [
       "from pynq import DefaultHierarchy, DefaultIP, allocate\n",
       "sendchannel = overlay.axi_dma_0.sendchannel\n",
       "recvchannel = overlay.axi_dma_0.recvchannel"
      ]
     },
     {
      "cell_type": "code",
      "execution_count": null,
      "metadata": {},
      "outputs": [],
      "source": [
       "X_test = np.load('banknote-authentication_X_test.npy')\n",
       "Y_test = np.load('banknote-authentication_y_test.npy')"
      ]
     },
     {
      "cell_type": "code",
      "execution_count": null,
      "metadata": {},
      "outputs": [],
      "source": [
       "# COSTANTS\n",
       "SAMPLES = len(X_test)\n",
       " \n",
       "BATCH_SIZE = 16  # SIZE OF THE BATCH TO SEND\n",
       "BM_INPUTS  = 4   # N. OF INPUTS OF THE BONDMACHINE MODULE \n",
       "BM_OUTPUTS = 3   # N. OF OUTPUTS OF THE BONDMACHINE MODULE\n",
       "PRECISION  = 16\n",
       "\n",
       "INPUT_SHAPE  = (BATCH_SIZE, BM_INPUTS)\n",
       "OUTPUT_SHAPE = (BATCH_SIZE, BM_OUTPUTS)"
      ]
     },
     {
      "cell_type": "code",
      "execution_count": null,
      "metadata": {},
      "outputs": [],
      "source": [
       "n_batches = 0\n",
       "fill = False\n",
       "if (SAMPLES/BATCH_SIZE % 2 != 0):\n",
       "    n_batches = int(SAMPLES/BATCH_SIZE) + 1\n",
       "    fill = True\n",
       "else:\n",
       "    n_batches = int(SAMPLES/BATCH_SIZE)"
      ]
     },
     {
      "cell_type": "code",
      "execution_count": null,
      "metadata": {},
      "outputs": [],
      "source": [
       "def random_pad(vec, pad_width, *_, **__):\n",
       "    vec[:pad_width[0]] = np.random.uniform(0, 1, size=pad_width[0])\n",
       "    vec[vec.size-pad_width[1]:] = np.random.uniform(0,1, size=pad_width[1])"
      ]
     },
     {
      "cell_type": "code",
      "execution_count": null,
      "metadata": {},
      "outputs": [],
      "source": [
       "batches = []\n",
       "last_batch_size = 0\n",
       "for i in range(0, n_batches):\n",
       "    new_batch = X_test[i*BATCH_SIZE:(i+1)*BATCH_SIZE]\n",
       "    if (len(new_batch) < BATCH_SIZE):\n",
       "        last_batch_size = len(new_batch)\n",
       "        new_batch = np.pad(new_batch,  [(0, BATCH_SIZE-len(new_batch)), (0,0)], mode=random_pad)\n",
       "    batches.append(new_batch)"
      ]
     },
     {
      "cell_type": "code",
      "execution_count": null,
      "metadata": {},
      "outputs": [],
      "source": [
       "outputs = []\n",
       "inference_times = []\n",
       "for i in range(0, len(batches)):\n",
       "    input_buffer = allocate(shape=INPUT_SHAPE, dtype=np.float16)\n",
       "    output_buffer = allocate(shape=OUTPUT_SHAPE, dtype=np.uint16)\n",
       "    input_buffer[:]=batches[i]\n",
       "    start_time = datetime.now()\n",
       "    sendchannel.transfer(input_buffer)\n",
       "    recvchannel.transfer(output_buffer)\n",
       "    sendchannel.wait()\n",
       "    recvchannel.wait()\n",
       "    end_time = datetime.now()\n",
       "    diff_in_ms = (end_time - start_time).total_seconds() * 1000\n",
       "    inference_times.append(diff_in_ms)\n",
       "    if fill == True and i == len(batches) - 1:\n",
       "        outputs.append(output_buffer[0:last_batch_size])\n",
       "    else:\n",
       "        outputs.append(output_buffer)"
      ]
     },
     {
      "cell_type": "code",
      "execution_count": null,
      "metadata": {},
      "outputs": [],
      "source": [
       "ms_median = statistics.median(inference_times)\n",
       "print(\"median:            \", ms_median, \" ms\")\n",
       "print(\"median time for a sample: \", statistics.median(inference_times)/BATCH_SIZE, \" ms\")"
      ]
     },
     {
      "cell_type": "code",
      "execution_count": null,
      "metadata": {},
      "outputs": [],
      "source": [
       "def bin_to_float16(binary_str):\n",
       "    binary_bytes = int(binary_str, 2).to_bytes(2, byteorder='big')\n",
       "    float_val = struct.unpack('>e', binary_bytes)[0]\n",
       "    return float_val"
      ]
     },
     {
      "cell_type": "code",
      "execution_count": null,
      "metadata": {},
      "outputs": [],
      "source": [
       "def bin_to_float32(binary):\n",
       "    return struct.unpack('!f',struct.pack('!I', int(binary, 2)))[0]"
      ]
     },
     {
      "cell_type": "code",
      "execution_count": null,
      "metadata": {},
      "outputs": [],
      "source": [
       "import struct\n",
       "results_to_dump = []\n",
       "\n",
       "for outcome in outputs:\n",
       "    for out in outcome:\n",
       "        \n",
       "        if PRECISION == 16:\n",
       "            binary_str_prob0 = bin(out[0])[2:]\n",
       "            prob_0_float_value = bin_to_float16(binary_str_prob0)\n",
       "\n",
       "            binary_str_prob1 = bin(out[1])[2:]\n",
       "            prob_1_float_value = bin_to_float16(binary_str_prob1)\n",
       "            \n",
       "        elif PRECISION == 32:\n",
       "            bin_str = bin(out[0])[2:].zfill(32)\n",
       "            byte_str = int(bin_str, 2).to_bytes(4, byteorder='big')\n",
       "            prob_0_float_value = struct.unpack('>f', byte_str)[0]\n",
       "\n",
       "            bin_str = bin(out[1])[2:].zfill(32)\n",
       "            byte_str = int(bin_str, 2).to_bytes(4, byteorder='big')\n",
       "            prob_1_float_value = struct.unpack('>f', byte_str)[0]\n",
       "        \n",
       "        classification = np.argmax([prob_0_float_value, prob_1_float_value])\n",
       "        results_to_dump.append([prob_0_float_value, prob_1_float_value, classification, out[2]])"
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
       "with open(\"predictions.csv\", 'w') as f:\n",
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
