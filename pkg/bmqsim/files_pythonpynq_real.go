package bmqsim

const (
	PythonPynqReal = `
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
    "# COSTANTS\n",
    "#SAMPLES = len(X_test)\n",
    " \n",
    "BATCH_SIZE = 16  # SIZE OF THE BATCH TO SEND\n",
    "BM_INPUTS  = 2   # N. OF INPUTS OF THE BONDMACHINE MODULE \n",
    "BM_OUTPUTS = 2   # N. OF OUTPUTS OF THE BONDMACHINE MODULE\n",
    "PRECISION  = 32\n",
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
    "input_buffer = allocate(shape=INPUT_SHAPE, dtype=np.float32)\n",
    "output_buffer = allocate(shape=OUTPUT_SHAPE, dtype=np.float32)"
   ]
  },
  {
   "cell_type": "code",
   "execution_count": null,
   "metadata": {},
   "outputs": [],
   "source": [
    "random_array = np.random.randint(low=0, high=100, size=(16, 2))"
   ]
  },
  {
   "cell_type": "code",
   "execution_count": null,
   "metadata": {},
   "outputs": [],
   "source": [
    "input_buffer[:] = random_array"
   ]
  },
  {
   "cell_type": "code",
   "execution_count": null,
   "metadata": {},
   "outputs": [],
   "source": [
    "input_buffer"
   ]
  },
  {
   "cell_type": "code",
   "execution_count": null,
   "metadata": {},
   "outputs": [],
   "source": [
    "sendchannel.transfer(input_buffer)\n",
    "recvchannel.transfer(output_buffer)"
   ]
  },
  {
   "cell_type": "code",
   "execution_count": null,
   "metadata": {},
   "outputs": [],
   "source": [
    "recvchannel.idle"
   ]
  },
  {
   "cell_type": "code",
   "execution_count": null,
   "metadata": {},
   "outputs": [],
   "source": [
    "output_buffer"
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
