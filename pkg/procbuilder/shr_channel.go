package procbuilder

import (
	"strconv"
	"strings"
)

type Channel struct{}

func (op Channel) Shr_get_name() string {
	return "channel"
}

func (op Channel) Shortname() string {
	return "ch"
}

func (op Channel) Get_header(arch *Arch, shared_constraint string, seq int) string {
	chname := "ch" + strconv.Itoa(seq)
	return ", " + chname + "in, " + chname + "wwr, " + chname + "wrd, " + chname + "ack_ch_ready, " + chname + "op_check_ready, " + chname + "finish_channel, " + chname + "out, " + chname + "ack_wwr, " + chname + "ack_wrd, " + chname + "ch_ready, " + chname + "ch_w_r_ready"
}

func (op Channel) Get_params(arch *Arch, shared_constraint string, seq int) string {

	chname := "ch" + strconv.Itoa(seq)

	result := ""
	result += "	output [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] " + chname + "in;\n"
	result += "	output " + chname + "wwr;\n"
	result += "	output " + chname + "wrd;\n"
	result += "	output " + chname + "ack_ch_ready;\n"
	result += "	output " + chname + "op_check_ready;\n"
	result += "	input " + chname + "finish_channel;\n"
	result += "	input [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] " + chname + "out;\n"
	result += "	input " + chname + "ack_wwr;\n"
	result += "	input " + chname + "ack_wrd;\n"
	result += "	input " + chname + "ch_ready;\n"
	result += "	input [1:0] " + chname + "ch_w_r_ready;\n"

	return result
}

func (op Channel) Get_internal_params(arch *Arch, shared_constraint string, seq int) string {
	channel_num := 0
	if arch.Shared_constraints != "" {
		constraints := strings.Split(arch.Shared_constraints, ",")
		for _, constraint := range constraints {
			values := strings.Split(constraint, ":")
			soname := values[0]
			if soname == "channel" {
				channel_num++
			}
		}
	}

	chname := "ch" + strconv.Itoa(seq)

	result := ""
	result += "	output [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] " + chname + "in;\n"
	result += "	output " + chname + "wwr;\n"
	result += "	output " + chname + "wrd;\n"
	result += "	output " + chname + "ack_ch_ready;\n"
	result += "	output " + chname + "op_check_ready;\n"
	result += "	input " + chname + "finish_channel;\n"
	result += "	input [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] " + chname + "out;\n"
	result += "	input " + chname + "ack_wwr;\n"
	result += "	input " + chname + "ack_wrd;\n"
	result += "	input " + chname + "ch_ready;\n"
	result += "	input [1:0] " + chname + "ch_w_r_ready;\n"

	if seq == 0 {
		result += "\twire [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] ch2proc_i[" + strconv.Itoa(channel_num-1) + ":0];\n"
		result += "\treg [" + strconv.Itoa(int(arch.Rsize)-1) + ":0] proc2ch_i [" + strconv.Itoa(channel_num-1) + ":0];\n"
		result += "\twire [" + strconv.Itoa(channel_num-1) + ":0] ch_wwr_i;\n"
		result += "\twire [" + strconv.Itoa(channel_num-1) + ":0] ch_wrd_i;\n"
		result += "\twire [" + strconv.Itoa(channel_num-1) + ":0] ack_wwr_i;\n"
		result += "\twire [" + strconv.Itoa(channel_num-1) + ":0] ack_wrd_i;\n"
		result += "\twire [" + strconv.Itoa(channel_num-1) + ":0] ch_ready_i;\n"
		result += "\twire [1:0] ch_w_r_ready_i [" + strconv.Itoa(channel_num-1) + ":0];\n"
		result += "\treg [" + strconv.Itoa(channel_num-1) + ":0] ack_ch_ready_i;\n"
		result += "\treg [" + strconv.Itoa(channel_num-1) + ":0] ch_op_ready_i;\n"
		result += "\twire [" + strconv.Itoa(channel_num-1) + ":0] finish_channel_i;\n"
	}

	result += "\tassign ch2proc_i[" + strconv.Itoa(seq) + "] = " + chname + "out;\n"
	result += "\tassign " + chname + "in = proc2ch_i[" + strconv.Itoa(seq) + "];\n"
	result += "\tassign " + chname + "wwr = ch_wwr_i[" + strconv.Itoa(seq) + "];\n"
	result += "\tassign " + chname + "wrd = ch_wrd_i[" + strconv.Itoa(seq) + "];\n"
	result += "\tassign " + chname + "ack_ch_ready = ack_ch_ready_i[" + strconv.Itoa(seq) + "];\n"
	result += "\tassign " + chname + "op_check_ready = ch_op_ready_i[" + strconv.Itoa(seq) + "];\n"
	result += "\tassign ack_wwr_i[" + strconv.Itoa(seq) + "] = " + chname + "ack_wwr;\n"
	result += "\tassign ack_wrd_i[" + strconv.Itoa(seq) + "] = " + chname + "ack_wrd;\n"
	result += "\tassign ch_ready_i[" + strconv.Itoa(seq) + "] = " + chname + "ch_ready;\n"
	result += "\tassign ch_w_r_ready_i[" + strconv.Itoa(seq) + "] = " + chname + "ch_w_r_ready;\n"
	result += "\tassign finish_channel_i[" + strconv.Itoa(seq) + "] = " + chname + "finish_channel;\n"
	return result
}
