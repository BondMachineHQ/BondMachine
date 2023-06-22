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
					"def read_mk_file(file_path):\n",
					"    dictionary = {}\n",
					"    with open(file_path, 'r') as file:\n",
					"        for line in file:\n",
					"            line = line.strip()\n",
					"            if line and not line.startswith('#') and '=' in line:\n",
					"                key, value = line.split('=', 1)\n",
					"                dictionary[key.strip()] = value.strip()\n",
					"    return dictionary"
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
					"    directory_name = filename.replace(\"_simreport\", \"\")\n",
					"    runinfo[\"templateinfo\"] = read_mk_file(directory_name+\"/generated.mk\")\n",
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
					"                r[\"templateinfo\"].update({\"template\":r[\"name\"]})\n",
					"                to_plot.append({\n",
					"                    \"name\": r[\"name\"],\n",
					"                    \"resource\": res,\n",
					"                    \"values\": values_in_scope.tolist(),\n",
					"                    \"info\": r[\"templateinfo\"]\n",
					"                })\n",
					"            \n",
					"            \n",
					"            fig_height = len(to_plot)*2\n",
					"            fig_width = len(to_plot)*4\n",
					"            \n",
					"            #plt.figure(figsize=(20, 6))\n",
					"            fig, (ax1, ax2) = plt.subplots(1, 2, gridspec_kw={'width_ratios': [4, 2]}, figsize=(fig_width, fig_height))\n",
					"            ax2.set_axis_off()\n",
					"            info_to_plot = r[\"templateinfo\"]\n",
					"            \n",
					"            pos=1\n",
					"            for t in to_plot:\n",
					"                ax1.plot(ticks, t[\"values\"], label=t[\"name\"])\n",
					"                last_value = t[\"values\"][-1]\n",
					"                #ax1.text(ticks[-1], last_value, t[\"name\"], ha='center', va='bottom', fontsize=14)\n",
					"                \n",
					"                for key,value in t[\"info\"].items():\n",
					"                    ax2.text(0.0, pos, key+\": \"+str(value))\n",
					"                    pos-=0.05\n",
					"                \n",
					"                pos -= 0.05\n",
					"            ax1.legend()\n",
					"            \n",
					"    \n",
					"    if pl2:\n",
					"        for res in resources:        \n",
					"            for i in range(0, len(runs)):\n",
					"                ticks = runs[i][\"rundata\"].index\n",
					"                if i != pivot:\n",
					"                    diff = abs((runs[i][\"rundata\"][res] - runs[pivot][\"rundata\"][res]).mean())\n",
					"                    \n",
					"                    if diff > 10:\n",
					"                        print(\" For resource: \", res, \" mean difference of value between \", runs[i][\"name\"], \"and \", runs[pivot][\"name\"] , \" is \", diff)\n",
					"        \n",
					"    "
				]
			},
			{
				"cell_type": "code",
				"execution_count": null,
				"id": "976bd531-2a31-4a80-b1d9-0cc7f4cc9ce8",
				"metadata": {},
				"outputs": [],
				"source": [
					"template_1d6c78b2_d2ea_403e_b1c2_d6226b866705=load_data(\"template_1d6c78b2_d2ea_403e_b1c2_d6226b866705_simreport\", False)"
				]
			},
			{
				"cell_type": "code",
				"execution_count": null,
				"id": "976bd531-2a31-4a80-b1d9-0cc7f4cc9ce8",
				"metadata": {},
				"outputs": [],
				"source": [
					"template_5b0b5dd1_8bbb_40cf_b5eb_447232f1eb88=load_data(\"template_5b0b5dd1_8bbb_40cf_b5eb_447232f1eb88_simreport\", False)"
				]
			},
			{
				"cell_type": "code",
				"execution_count": null,
				"id": "976bd531-2a31-4a80-b1d9-0cc7f4cc9ce8",
				"metadata": {},
				"outputs": [],
				"source": [
					"template_d4b87f3a_7764_4aa7_af73_66730e79f5b6=load_data(\"template_d4b87f3a_7764_4aa7_af73_66730e79f5b6_simreport\", False)"
				]
			},
			{
				"cell_type": "code",
				"execution_count": null,
				"id": "976bd531-2a31-4a80-b1d9-0cc7f4cc9ce8",
				"metadata": {},
				"outputs": [],
				"source": [
					"template_032a97ef_10f7_48d6_9c5a_8f5b56b32995=load_data(\"template_032a97ef_10f7_48d6_9c5a_8f5b56b32995_simreport\", False)"
				]
			},
			{
				"cell_type": "code",
				"execution_count": null,
				"metadata": {},
				"outputs": [],
				"source": [
					"analyze([ template_1d6c78b2_d2ea_403e_b1c2_d6226b866705, template_5b0b5dd1_8bbb_40cf_b5eb_447232f1eb88, template_d4b87f3a_7764_4aa7_af73_66730e79f5b6, template_032a97ef_10f7_48d6_9c5a_8f5b56b32995], 0, True, True, [])"
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
