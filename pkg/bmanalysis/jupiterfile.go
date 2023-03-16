package bmanalysis

const (
	notebook = `
	{
		"cells": [
		 {
		  "cell_type": "code",
		  "execution_count": 9,
		  "id": "6293e6e8-1433-4161-9e04-6c3d670b0e8c",
		  "metadata": {},
		  "outputs": [],
		  "source": [
		   "!pip install pandas\n",
		   "!pip install matplotlib\n",
		   "!pip install scipy"
		  ]
		 },
		{
			"cell_type": "code",
			"execution_count": 17,
			"id": "ff7dc30b-f076-446f-977d-5534a4120f20",
			"metadata": {},
			"outputs": [],
			"source": [
			 "import pandas as pd\n",
			 "import json\n",
			 "import numpy as np\n",
			 "import matplotlib.pyplot as plt\n",
			 "import matplotlib.image as mpimg\n",
			 "from scipy.stats import norm"
			]
		},
		{
			"cell_type": "code",
			"execution_count": 18,
			"id": "1cf6e3c0-5f9f-48ac-99c5-a074ffc955a0",
			"metadata": {},
			"outputs": [],
			"source": [
			 "class colors:\n",
			 "    WHITE = '\\033[97m'\n",
			 "    BLACK = '\\033[90m'\n",
			 "    HEADER = '\\033[95m'\n",
			 "    BLUE = '\\033[94m'\n",
			 "    CYAN = '\\033[96m'\n",
			 "    GREEN = '\\033[92m'\n",
			 "    YELLOW = '\\033[93m'\n",
			 "    RED = '\\033[91m'\n",
			 "    ENDC = '\\033[0m'\n",
			 "    BOLD = '\\033[1m'\n",
			 "    UNDERLINE = '\\033[4m'\n",
			 "    \n",
			 "def blue(text):\n",
			 "    return colors.BLUE+text+colors.BLACK\n",
			 "\n",
			 "def green(text):\n",
			 "    return colors.GREEN+text+colors.BLACK\n",
			 "\n",
			 "def red(text):\n",
			 "    return colors.RED+text+colors.BLACK\n",
			 "\n",
			 "def yellow(text):\n",
			 "    return colors.YELLOW+text+colors.BLACK"
			]
		},
		{
			"cell_type": "code",
			"execution_count": 122,
			"id": "c9269fa0-2479-47ee-a261-1981ab82ddf2",
			"metadata": {},
			"outputs": [],
			"source": [
			 "def load_run(runname):\n",
			 "    rundata=pd.read_csv(runname+'.csv')\n",
			 "    \n",
			 "    with open(runname+'.json') as f:\n",
			 "        runinfo=json.load(f)\n",
			 "\n",
			 "    if 'clock_freq' in runinfo:\n",
			 "        clock_freq=runinfo['clock_freq']\n",
			 "    else:\n",
			 "        clock_freq=None\n",
			 "\n",
			 "    if 'clock_cycles' in rundata:\n",
			 "        clock_mean=rundata['clock_cycles'].mean()\n",
			 "        clock_std=rundata['clock_cycles'].std()\n",
			 "        inference_time_mean=clock_mean/(clock_freq)\n",
			 "        \n",
			 "        fig, ax = plt.subplots(1, 3, gridspec_kw={'width_ratios': [3, 1, 2]}, figsize=(15, 4))\n",
			 "            \n",
			 "        ax[0].set_title(\"Clock cicles distribution\")\n",
			 "        ax[0].set_facecolor('#aaaaaa')\n",
			 "        rundata.plot.hist(ax=ax[0],column=[\"clock_cycles\"], bins=50,density=True)\n",
			 "        x = np.linspace(0, 20000, 100)\n",
			 "        p = norm.pdf(x, clock_mean, clock_std)\n",
			 "        ax[0].plot(x, p, 'k', linewidth=2)\n",
			 "        \n",
			 "        ax[1].set_title(\"Run report\")\n",
			 "        ax[1].set_axis_off()\n",
			 "        pos=0.9\n",
			 "        for key,value in sorted(runinfo.items()):\n",
			 "            ax[1].text(0.0, pos, key+\": \"+str(value))\n",
			 "            pos-=0.1\n",
			 "        ax[1].text(0.0, pos, \"clock_mean: \"+str(clock_mean)[:10])\n",
			 "        pos-=0.1\n",
			 "        ax[1].text(0.0, pos, \"clock_std: \"+str(clock_std)[:10])\n",
			 "        pos-=0.1\n",
			 "        ax[1].text(0.0, pos, \"inference_time_mean: \"+str(inference_time_mean)[:10]+\" us\")\n",
			 "\n",
			 "        ax[2].set_title(\"BondMachine diagram\")\n",
			 "        ax[2].set_axis_off()\n",
			 "        img = mpimg.imread(runname+'.png')\n",
			 "        ax[2].imshow(img)\n",
			 "        \n",
			 "    else:\n",
			 "        clock_mean=None\n",
			 "        clock_std=None\n",
			 "        inference_time_mean=None\n",
			 "        \n",
			 "    runinfo[\"runname\"]=runname\n",
			 "    runinfo[\"rundata\"]=rundata\n",
			 "    runinfo[\"clock_freq\"]=clock_freq\n",
			 "    runinfo[\"clock_mean\"]=clock_mean\n",
			 "    runinfo[\"clock_std\"]=clock_std\n",
			 "    runinfo[\"inference_time_mean\"]=inference_time_mean\n",
			 "    \n",
			 "    return runinfo"
			]
		},
		{
			"cell_type": "code",
			"execution_count": 216,
			"id": "5d8e2b9f-f5d0-41b7-9a44-e46f5e1462b8",
			"metadata": {},
			"outputs": [],
			"source": [
			 "def runplot(runs, pivot_run, x_param, y_param, title):\n",
			 "    x=[]\n",
			 "    y=[]\n",
			 "    names=[]\n",
			 "    for i in range(len(runs)):\n",
			 "        if x_param in runs[i] and y_param in runs[i]:\n",
			 "            x.append(runs[i][x_param])\n",
			 "            y.append(runs[i][y_param])\n",
			 "            names.append(runs[i][\"runname\"])\n",
			 "\n",
			 "    if len(names) < 2:\n",
			 "        print (\"Not enought data to plot\")\n",
			 "        return\n",
			 "            \n",
			 "    fig, ax = plt.subplots(1, 2, gridspec_kw={'width_ratios': [3, 1]}, figsize=(15, 4))\n",
			 "    \n",
			 "    style = dict(size=10, color='black')\n",
			 "    \n",
			 "    if title != None:\n",
			 "        ax[0].set_title(title)\n",
			 "    else:\n",
			 "        ax[0].set_title(x_param + \" vs \"+ y_param)\n",
			 "    ax[0].set_facecolor('#aaaaaa')\n",
			 "    for i in range(len(runs)):\n",
			 "        if x_param in runs[i] and y_param in runs[i]:            \n",
			 "            ax[0].text(runs[i][x_param], runs[i][y_param], runs[i][\"runname\"], **style)\n",
			 "    ax[0].plot(x, y ,linewidth=2)\n",
			 "    \n",
			 "    ax[1].set_title(\"Data\")\n",
			 "    ax[1].set_axis_off()\n",
			 "    pos=0.9\n",
			 "    ax[1].text(0.0, pos, \"run name - \"+x_param+\" - \"+y_param)\n",
			 "    pos-=0.1\n",
			 "    for i in range(len(names)):\n",
			 "        ax[1].text(0.0, pos, names[i]+\" - \"+str(x[i])+\" - \"+str(y[i]))\n",
			 "        pos-=0.1\n",
			 "        \n",
			 "    display(fig)\n",
			 "    plt.close()"
			]
		},
		   {
			"cell_type": "code",
			"execution_count": 217,
			"id": "ddcd0fd8-3352-437b-8560-aa4865c2aeeb",
			"metadata": {},
			"outputs": [],
			"source": [
			 "def consistency(runs, pivot_run):\n",
			 "    # Result consistency, compare to pivot run every other run by mean and std\n",
			 "    for i in range(len(runs)):\n",
			 "        if i!= pivot_run:\n",
			 "            runname=runs[i][\"runname\"]\n",
			 "            rundata0=abs(runs[i][\"rundata\"][\"probability_0\"]-runs[pivot_run][\"rundata\"][\"probability_0\"])\n",
			 "            rundata0_mean=rundata0.mean()\n",
			 "            rundata0_std=rundata0.std()\n",
			 "            rundata1=abs(runs[i][\"rundata\"][\"probability_1\"]-runs[pivot_run][\"rundata\"][\"probability_1\"])\n",
			 "            rundata1_mean=rundata1.mean()\n",
			 "            rundata1_std=rundata1.std()\n",
			 "            prediction=runs[i][\"rundata\"][\"classification\"]==runs[pivot_run][\"rundata\"][\"classification\"]\n",
			 "           \n",
			 "            print (blue(runname)+\" - probability_0 - mean: \"+str(rundata0_mean)+ \" std: \"+str(rundata0_std))\n",
			 "            print (blue(runname)+\" - probability_1 - mean: \"+str(rundata1_mean)+ \" std: \"+str(rundata1_std))\n",
			 "            print (blue(runname)+\" - prediction    - \"+ green(str(prediction.value_counts(normalize=True)[True]*100)+ \"%\")+\" equal\")\n",
			 "            print ()"
			]
		},
		   {
			"cell_type": "code",
			"execution_count": 234,
			"id": "3fdf043c-fba9-41bf-ad47-11b2cbcb9454",
			"metadata": {},
			"outputs": [],
			"source": [
			 "def analyze(runs, pivot_run):\n",
			 "    pivot_name=runs[pivot_run][\"runname\"]\n",
			 "    print (\"Pivot run \"+ blue(pivot_name))\n",
			 "    print ()\n",
			 "    print (red(\"Consistency, compare to pivot run every other run\"))\n",
			 "    consistency(runs, pivot_run)\n",
			 "    \n",
			 "    print (red(\"Inference time mean\"))\n",
			 "    runplot(runs, pivot_run, \"expprec\", \"inference_time_mean\", \"\")\n",
			 "    \n",
			 "#     print (red(\"Occupancy\"))\n",
			 "#     runplot(runs, pivot_run, \"cps\", \"luts\", \"Occupancy (cps vs luts)\")\n",
			 "\n",
			 "#     print (red(\"Power consumption\"))\n",
			 "#     runplot(runs, pivot_run, \"cps\", \"power\", \"Power consumption (cps vs power)\")"
			]
		   },
		   {
			"cell_type": "code",
			"execution_count": 226,
			"id": "a2f43b12-5cf2-4e63-b646-c525a573adc9",
			"metadata": {},
			"outputs": [],
			"source": [
			 "soft=load_run(\"sw\")\n",
			 "sim=load_run(\"sim\")"
			]
		},
		{{- if .ProjectLists }}
    	{{- range .ProjectLists }}
		{
			"cell_type": "code",
			"execution_count": 232,
			"id": "976bd531-2a31-4a80-b1d9-0cc7f4cc9ce8",
			"metadata": {},
			"outputs": [],
			"source": [
			 "expanded_{{ . }}=load_run(\"expanded_{{ . }}\")"
			]
		},
		{{- end }}
    	{{- end }}
		{
			"cell_type": "code",
			"execution_count": 250,
			"id": "42a2b462-4534-4475-83d4-cb691c758894",
			"metadata": {},
			"outputs": [],
			"source": [
			 "analyze([soft, sim {{- if .ProjectLists }} {{- range .ProjectLists }},expanded_{{ . }} {{- end }} {{- end }}],0)"
			]
		},
		{
			"cell_type": "code",
			"execution_count": null,
			"id": "059800fc-01dc-49cd-a64e-1aa506d08db0",
			"metadata": {},
			"outputs": [],
			"source": []
		   }
		  ],
		  "metadata": {
		   "kernelspec": {
			"display_name": "Python 3 (ipykernel)",
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
			"version": "3.10.6"
		   },
		   "vscode": {
			"interpreter": {
			 "hash": "97b4c6de4695f1fe6bfabfed1e13c8c763fcd1e15536b82de4b874cda1fc74b7"
			}
		   }
		  },
		  "nbformat": 4,
		  "nbformat_minor": 5
	}
`
)
