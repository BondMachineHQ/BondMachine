#!/bin/bash

WEBDIR=$1
if [[ "$WEBDIR" == "" ]]
then
	echo "Usage: $0 <webdir>"
	exit 1
fi

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

{{#length}}
**Length**: {{ length }}
{{/length}}

{{#desc}}
**Description**:

{{ desc }} {{ desc1 }} {{ desc2 }} {{ desc3 }} {{ desc4 }} {{ desc5 }} {{ desc6 }} {{ desc7 }} {{ desc8 }} {{ desc9 }}
{{/desc}}

{{#snippet}}
**Snippet**:

{{/snippet}}

EOF

cat > /tmp/dynoptemplate << EOF
# {{ name }}

**Instruction Group**: {{ name }}

{{#length}}
**Length**: {{ length }}
{{/length}}

{{#desc}}
**Description**:

{{ desc }} {{ desc1 }} {{ desc2 }} {{ desc3 }} {{ desc4 }} {{ desc5 }} {{ desc6 }} {{ desc7 }} {{ desc8 }} {{ desc9 }}
{{/desc}}

{{#snippet}}
**Snippet**:

{{/snippet}}
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
	if [[ "`jq .snippet <<<\"$YAMLDATA\"`" != "null" ]]
	then
		SNAME=`jq -r .snippet <<<"$YAMLDATA"`
		TITLE=`cat $WEBDIR/snippets/$SNAME/title`
		CODE=`cat $WEBDIR/snippets/$SNAME/code`
		DESC=`cat $WEBDIR/snippets/$SNAME/desc`
		echo "\`\`\`asm" >> ./reference/$opname.md
		echo "$CODE" >> ./reference/$opname.md
		echo "\`\`\`" >> ./reference/$opname.md
		echo "" >> ./reference/$opname.md
		echo "$DESC" >> ./reference/$opname.md
	fi
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
	if [[ "`jq .snippet <<<\"$YAMLDATA\"`" != "null" ]]
	then
		SNAME=`jq -r .snippet <<<"$YAMLDATA"`
		TITLE=`cat $WEBDIR/snippets/$SNAME/title`
		CODE=`cat $WEBDIR/snippets/$SNAME/code`
		DESC=`cat $WEBDIR/snippets/$SNAME/desc`
		echo "\`\`\`asm" >> ./reference/$opname.md
		echo "$CODE" >> ./reference/$opname.md
		echo "\`\`\`" >> ./reference/$opname.md
		echo "" >> ./reference/$opname.md
		echo "$DESC" >> ./reference/$opname.md
	fi
done

cat > ./reference/matrix.md << EOF
# Support Matrix

The following tables show the feature state in term of level of development for each instruction and instruction group.
For each of them the support of the features is shown.

The features are the following:
| Feature | Description |
| --- | --- |
| hdl | The instruction can be translated to hardware description language |
| hwopt | The instruction has hardware optimizations available |
| asm | The instruction can be assembled by the assembler |
| disasm | The instruction can be disassembled by the disassembler |
| hlasm | The instruction can be assembled by the high-level assembler (Basm) |
| asmeta | The instruction has metadata for the assembler |
| gosim | The instruction can be simulated in the Go-based simulator |
| gosimlat | The instruction simulation has latency setup aligned to the hardware |
| hdlsim | The instruction can be simulated in the hardware description language simulator |
| mt | Instruction multi-thread support |

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
for feature in $(echo "${!SupportArray[@]}" | tr ' ' '\n' | sort)
do
	echo -n " $feature |" >> ./reference/matrix.md
done
echo "" >> ./reference/matrix.md

echo -n "| --- |" >> ./reference/matrix.md
for feature in $(echo "${!SupportArray[@]}" | tr ' ' '\n' | sort)
do
	echo -n " --- |" >> ./reference/matrix.md
done
echo "" >> ./reference/matrix.md

for i in `ls ../procbuilder/op_* | sort`
do
	opname=`basename $i | cut -d_ -f2- | cut -d. -f1`
	echo -n "| [$opname]($opname.md) |" >> ./reference/matrix.md
	for feature in $(echo "${!SupportArray[@]}" | tr ' ' '\n' | sort)
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
for feature in $(echo "${!SupportArrayDyn[@]}" | tr ' ' '\n' | sort)
do
	echo -n " $feature |" >> ./reference/matrix.md
done

echo "" >> ./reference/matrix.md

echo -n "| --- |" >> ./reference/matrix.md
for feature in $(echo "${!SupportArrayDyn[@]}" | tr ' ' '\n' | sort)
do
	echo -n " --- |" >> ./reference/matrix.md
done
echo "" >> ./reference/matrix.md

for i in `ls ../procbuilder/dynop_* | sort`
do
	opname=`basename $i | cut -d_ -f2- | cut -d. -f1`
	echo -n "| [$opname]($opname.md) |" >> ./reference/matrix.md
	for feature in $(echo "${!SupportArrayDyn[@]}" | tr ' ' '\n' | sort)
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