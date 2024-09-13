#!/bin/bash

cat > ./reference/README.md << EOF
# Assembly Reference

This directory contains reference documentation for the BASM assembly language.
The instructions listed here are only the 1 on 1 mapping of the machine code instructions
that are available in the BondMachine project. BASM supports these instructions and
pseudo-instructions as well. The pseudo-instructions are not listed here.

## Instructions

EOF

cat > /tmp/opcodetemplate << EOF
# {{ name }}

**Instruction**: {{ name }}

**Length**: {{ length }}

**Description**:

{{ desc }} {{ desc1 }} {{ desc2 }} {{ desc3 }} {{ desc4 }} {{ desc5 }} {{ desc6 }} {{ desc7 }} {{ desc8 }} {{ desc9 }}

EOF

cat > /tmp/dynoptemplate << EOF
Name {{ name }}
EOF

declare -A SupportArray

for i in `ls ../procbuilder/op_* | sort`
do
	opname=`basename $i | cut -d_ -f2- | cut -d. -f1`
	echo "Building instruction $opname reference"
	echo "[$opname]($opname.md)" >> ./reference/README.md
	YAMLDATA=`echo "{ \"name\": \"$opname\" }"`
	export IFS=$'\n'
	for ijson in `cat ../procbuilder/op_$opname.go | grep -Eo '"reference":.*{.*}' | cut -d: -f2-`
	do
		key=`jq -r 'keys[]' <<<"$ijson"`
		if [[ "$key" =~ "support_" ]]
		then
			feature=`echo $key | cut -d_ -f2-`
			SupportArray[$feature]="true"
		fi
		YAMLDATA=`jq '. + '$ijson <<<"$YAMLDATA"`
	done
	unset IFS
	# echo $YAMLDATA
	echo $YAMLDATA | mustache /tmp/opcodetemplate > ./reference/$opname.md 
done

cat >> ./reference/README.md << EOF

## Dynamical Instructions

EOF

declare -A SupportArrayDyn

for i in `ls ../procbuilder/dynop_* | sort`
do
	opname=`basename $i | cut -d_ -f2- | cut -d. -f1`
	echo "Building dynamical instruction group $opname reference"
	echo "[$opname]($opname.md)" >> ./reference/README.md
	YAMLDATA=`echo "{ \"name\": \"$opname\" }"`
	export IFS=$'\n'
	for ijson in `cat ../procbuilder/dynop_$opname.go | grep -Eo '"reference":.*{.*}' | cut -d: -f2-`
	do
		key=`jq -r 'keys[]' <<<"$ijson"`
		if [[ "$key" =~ "support_" ]]
		then
			feature=`echo $key | cut -d_ -f2-`
			SupportArrayDyn[$feature]="true"
		fi
		YAMLDATA=`jq '. + '$ijson <<<"$YAMLDATA"`
	done
	unset IFS
	# echo $YAMLDATA
	echo $YAMLDATA | mustache /tmp/dynoptemplate > ./reference/$opname.md
done

cat > ./reference/matrix.md << EOF
# Support Matrix

The following tables show the feature support of the instructions and the dynamical instruction in the BondMachine project.
For each of them the support of the features is shown.

The features are the following:
| Feature | Description |
| --- | --- |
| hdl | The instruction can be translated to hardware description language |
| gosim | The instruction can be simulated in the Go-based simulator |

The possible support values are shown below:

| Value | Meaning |
| --- | --- |
| ![ok](iconok.png) | The feature is fully implemented |
| ![no](iconno.png) | The feature is not yet implemented |
| ![testing](icontesting.png) | The feature is being tested |
| ![partial](iconpartial.png) | The feature is partially implemented |
| ![notapplicable](iconnotapplicable.png) | The feature is not applicable to the instruction |

## Support Matrix for Static Instructions

EOF

echo -n "| Instruction |" >> ./reference/matrix.md
for feature in "${!SupportArray[@]}"
do
	echo -n " $feature |" >> ./reference/matrix.md
done
echo "" >> ./reference/matrix.md

echo -n "| --- |" >> ./reference/matrix.md
for feature in "${!SupportArray[@]}"
do
	echo -n " --- |" >> ./reference/matrix.md
done
echo "" >> ./reference/matrix.md

for i in `ls ../procbuilder/op_* | sort`
do
	opname=`basename $i | cut -d_ -f2- | cut -d. -f1`
	echo -n "| [$opname]($opname.md) |" >> ./reference/matrix.md
	for feature in "${!SupportArray[@]}"
	do
		valueok="false"
		export IFS=$'\n'
		for ijson in `cat ../procbuilder/op_$opname.go | grep -Eo '"reference":.*{.*}' | cut -d: -f2-`
		do
			key=`jq -r 'keys[]' <<<"$ijson"`
			if [[ "$key" == "support_$feature" ]]
			then
				value=`jq -r ".$key" <<<"$ijson"`
				case $value in
					"ok") value="![ok](iconok.png)" ;;
					"no") value="![no](iconno.png)" ;;
					"testing") value="![testing](icontesting.png)" ;;
					"partial") value="![partial](iconpartial.png)" ;;
					"notapplicable") value="![notapplicable](iconnotapplicable.png)" ;;
				esac
				echo -n " $value |" >> ./reference/matrix.md
				valueok="true"
				break
			fi

		done
		unset IFS
		if [[ "$valueok" == "false" ]]
		then
			echo -n " - |" >> ./reference/matrix.md
		fi
	done
	echo "" >> ./reference/matrix.md
done

cat >> ./reference/matrix.md << EOF

## Support Matrix for Dynamical Instructions

EOF

echo -n "| Instruction |" >> ./reference/matrix.md
for feature in "${!SupportArrayDyn[@]}"
do
	echo -n " $feature |" >> ./reference/matrix.md
done

echo "" >> ./reference/matrix.md

echo -n "| --- |" >> ./reference/matrix.md
for feature in "${!SupportArrayDyn[@]}"
do
	echo -n " --- |" >> ./reference/matrix.md
done
echo "" >> ./reference/matrix.md

for i in `ls ../procbuilder/dynop_* | sort`
do
	opname=`basename $i | cut -d_ -f2- | cut -d. -f1`
	echo -n "| [$opname]($opname.md) |" >> ./reference/matrix.md
	for feature in "${!SupportArrayDyn[@]}"
	do
		valueok="false"
		export IFS=$'\n'
		for ijson in `cat ../procbuilder/dynop_$opname.go | grep -Eo '"reference":.*{.*}' | cut -d: -f2-`
		do
			key=`jq -r 'keys[]' <<<"$ijson"`
			if [[ "$key" == "support_$feature" ]]
			then
				value=`jq -r ".$key" <<<"$ijson"`
				case $value in
					"ok") value="![ok](iconok.png)" ;;
					"no") value="![no](iconno.png)" ;;
					"testing") value="![testing](icontesting.png)" ;;
					"partial") value="![partial](iconpartial.png)" ;;
					"notapplicable") value="![notapplicable](iconnotapplicable.png)" ;;
				esac
				echo -n " $value |" >> ./reference/matrix.md
				valueok="true"
				break
			fi

		done
		unset IFS
		if [[ "$valueok" == "false" ]]
		then
			echo -n " - |" >> ./reference/matrix.md
		fi
	done
	echo "" >> ./reference/matrix.md
done

cat >> ./reference/README.md << EOF

Not all instructions are fully supported by the BondMachine project. Some instructions are in the process of being implemented.
The support matrix for the instructions is available [here](matrix.md).

EOF

rm /tmp/opcodetemplate
rm /tmp/dynoptemplate