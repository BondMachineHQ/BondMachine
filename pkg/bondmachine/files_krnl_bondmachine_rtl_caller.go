package bondmachine

const (
	krnlBondmachineRTLCaller = `
/**
* Copyright (C) 2019-2021 Xilinx, Inc
*
* Licensed under the Apache License, Version 2.0 (the "License"). You may
* not use this file except in compliance with the License. A copy of the
* License is located at
*
*     http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
* WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
* License for the specific language governing permissions and limitations
* under the License.
*/

////////////////////////////////////////////////////////////////////////////////
// Description: Basic Adder, no overflow. Unsigned. Combinatorial.
////////////////////////////////////////////////////////////////////////////////

` + "`default_nettype none" + `

module krnl_bondmachine_rtl_caller #(
  parameter integer C_DATA_WIDTH   = 32, // Data width of both input and output data
  parameter integer C_NUM_CHANNELS = 1   // Number of input channels.  Only a value of 2 implemented.
)
(
  input wire                                         aclk,
  input wire                                         areset,

  input wire                     s_tvalid,
  input wire  [C_NUM_CHANNELS-1:0][C_DATA_WIDTH-1:0] s_tdata,
  output wire                    s_tready,

  output wire                                        m_tvalid,
  output wire [C_DATA_WIDTH-1:0]                     m_tdata,
  input  wire                                        m_tready

);

timeunit 1ps; 
timeprecision 1ps; 


wire fready;
wire fmvalid;
wire f_saxistready;
wire f_maxistready;

logic [C_DATA_WIDTH-1:0] acc;
logic [C_DATA_WIDTH-1:0] bmacc;
wire last;
wire slast;

wire reset;
// reg  [31:0] counter;

// always @ (posedge clk) begin
//   counter <= counter + 1;
// end

// always_comb begin 
//   acc = s_tdata[0];
// end

bmaccelerator_v1_0 #( 
  .C_S00_AXIS_TDATA_WIDTH   ( C_DATA_WIDTH ) ,
  .C_M00_AXIS_TDATA_WIDTH   ( C_DATA_WIDTH ),
  .C_M00_AXIS_START_COUNT   ( 32)
)
inst_bmaccelerator_v1_0 ( 
  .s00_axis_aclk      ( aclk            ),
  .s00_axis_aresetn   ( !reset            ),
  .s00_axis_tready    ( s_tready    ),
  .s00_axis_tdata     ( s_tdata[0]     ),
  .s00_axis_tlast     ( slast),
  .s00_axis_tvalid    ( s_tvalid ), // s_tvalid

  .m00_axis_aclk    ( aclk),
  .m00_axis_aresetn ( !reset),
  .m00_axis_tvalid  ( m_tvalid      ),
  .m00_axis_tdata   ( m_tdata       ), // output
  .m00_axis_tlast   ( last ), // output
  .m00_axis_tready  ( m_tready   )
);

/*
TO THERE
*/

//assign s_tready = 1; // cosi quasi funziona


endmodule : krnl_bondmachine_rtl_caller

` + "`default_nettype wire" + `

`
)
