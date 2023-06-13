package bmanalysis

const (
	notebookMLSim = `
	{
		"cells": [
		 {
		  "cell_type": "code",
		  "execution_count": null,
		  "metadata": {},
		  "outputs": [],
		  "source": [
		   "import pandas as pd\n",
		   "import json\n",
		   "import numpy as np\n",
		   "import matplotlib.pyplot as plt\n",
		   "import re"
		  ]
		 },
		 {
			"cell_type": "code",
			"execution_count": null,
			"metadata": {},
			"outputs": [],
			"source": [
					"def fix_values(value):\n",
					"    try:\n",
					"        if isinstance(value, float) or isinstance(value, int):\n",
					"            return value\n",
					"        toclean_value = value.replace(\"0lq<16.1>\", \"\") #0lq<16.1>12.542724609375\n",
					"        toclean_value = toclean_value.replace(\"0lq<8.1>\", \"\") #0lq<16.1>12.542724609375\n",
					"        toclean_value = toclean_value.replace(\"0f<16>\", \"\")\n",
					"        toclean_value = toclean_value.replace(\"0f<32>\", \"\")\n",
					"        cleaned_value = re.sub(r'[^0-9.]', '', str(toclean_value))\n",
					"        return float(cleaned_value)\n",
					"    except Exception as e:\n",
					"        error_message = str(e)\n",
					"        if \"+Inf\" in value:\n",
					"            return 9999999999999999\n",
					"        else:\n",
					"            print(\"Unable to convert this value: \", toclean_value)\n",
					"            return 0"
			]
		 },
		 {
		  "cell_type": "code",
		  "execution_count": null,
		  "metadata": {},
		  "outputs": [],
		  "source": [
		   "def load_data(filename, toPlot):\n",
		   "    \n",
		   "    runinfo = {}\n",
		   "    \n",
		   "    rundata = pd.read_csv(filename+'.csv')\n",
		   "    columns = rundata.columns.tolist()\n",
		   "    columns = columns[1:len(columns)]\n",
		   "    \n",
		   "    rundata.set_index('tick', inplace=True)\n",
		   "    \n",
		   "    for c in columns:\n",
		   "        rundata[c] = rundata[c].apply(fix_values)\n",
		   "    \n",
		   "    if not toPlot:\n",
		   "        runinfo[\"name\"] = filename\n",
		   "        runinfo[\"rundata\"] = rundata\n",
		   "        runinfo[\"columns\"] = columns\n",
		   "        \n",
		   "        return runinfo\n",
		   "    \n",
		   "    for c in columns:\n",
		   "        \n",
		   "        plt.figure(figsize=(10, 6))\n",
		   "        rundata[c].plot()\n",
		   "        plt.xlabel('Tick', fontsize=16)\n",
		   "        plt.ylabel(c, fontsize=16)\n",
		   "        plt.xticks(fontsize=14)\n",
		   "        plt.yticks(fontsize=14)\n",
		   "        plt.title('Plot of '+c+' against Tick', fontsize=14)\n",
		   "        #plt.show()\n",
		   "        \n",
		   "    runinfo[\"name\"] = filename\n",
		   "    runinfo[\"rundata\"] = rundata \n",
		   "    runinfo[\"columns\"] = columns\n",
		   "    \n",
		   "    return runinfo\n",
		   "    "
		  ]
		 },
		 {
		  "cell_type": "code",
		  "execution_count": null,
		  "metadata": {},
		  "outputs": [],
		  "source": [
		   "def analyze(runs, pivot, pl1, pl2, specific_resources):\n",
		   "    \n",
		   "    resources = runs[0][\"columns\"]\n",
		   "    if specific_resources != None:\n",
		   "        if len(specific_resources) > 0:\n",
		   "            resources = specific_resources\n",
		   "            \n",
		   "    if pl1:\n",
		   "        for res in resources:\n",
		   "            to_plot = []\n",
		   "            ticks = []\n",
		   "            \n",
		   "            for r in runs:\n",
		   "                ticks = r[\"rundata\"].index\n",
		   "                values_in_scope = r[\"rundata\"][res]\n",
		   "                \n",
		   "                to_plot.append({\n",
		   "                    \"name\": r[\"name\"],\n",
		   "                    \"resource\": res,\n",
		   "                    \"values\": values_in_scope.tolist()\n",
		   "                })\n",
		   "            \n",
		   "            plt.figure(figsize=(10, 6))\n",
		   "                \n",
		   "            for t in to_plot:\n",
		   "                plt.plot(ticks, t[\"values\"], label=t[\"name\"])\n",
		   "                last_value = t[\"values\"][-1]\n",
		   "                plt.text(ticks[-1], last_value, t[\"name\"], ha='center', va='bottom', fontsize=14)\n",
		   "                \n",
		   "            plt.ylabel(t[\"resource\"],fontsize=14)    \n",
		   "            plt.xlabel('Ticks',fontsize=14)\n",
		   "    \n",
		   "    if pl2:\n",
		   "        for res in resources:        \n",
		   "            for i in range(0, len(runs)):\n",
		   "                ticks = runs[i][\"rundata\"].index\n",
		   "                if i != pivot:\n",
		   "                    diff = abs((runs[i][\"rundata\"][res] - runs[pivot][\"rundata\"][res]).mean())\n",
		   "                    \n",
		   "                    if diff > 0:\n",
		   "                        print(\" For resource: \", res, \" mean difference of value between \", runs[i][\"name\"], \"and \", runs[pivot][\"name\"] , \" is \", diff)\n",
		   "        \n",
		   "    "
		  ]
		 },
		{{- if .ProjectsList }}
		{{- range .ProjectsList }}
		{
			"cell_type": "code",
			"execution_count": null,
			"id": "976bd531-2a31-4a80-b1d9-0cc7f4cc9ce8",
			"metadata": {},
			"outputs": [],
			"source": [
			"{{ . }}=load_data(\"{{ . }}_simreport\", False)"
			]
		},
		{{- end }}
		{{- end }}
		 {
		  "cell_type": "code",
		  "execution_count": null,
		  "metadata": {},
		  "outputs": [],
		  "source": [
		   "analyze([{{- if .ProjectsList }} {{- $projectLen := len .ProjectsList }} {{- range $i, $project := .ProjectsList }} {{- if eq (inc $i) $projectLen }} {{ $project }} {{- else }} {{ $project }}, {{- end }} {{- end }} {{- end }}], 0, True, True, [])"
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
		  "version": "3.8.10"
		 },
		 "orig_nbformat": 4
		},
		"nbformat": 4,
		"nbformat_minor": 2
	   }
`
)
