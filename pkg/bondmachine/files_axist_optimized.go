package bondmachine

const (
	optimizedAXIStream = `

module {{ .ModuleName }}_v1_0 #
(
	parameter integer C_S00_AXIS_TDATA_WIDTH	= 32,
	parameter integer C_M00_AXIS_TDATA_WIDTH	= 32,
	parameter integer C_M00_AXIS_START_COUNT	= 32
)
(

	input wire  s00_axis_aclk,
	input wire  s00_axis_aresetn,
	output wire  s00_axis_tready,
	input wire [C_S00_AXIS_TDATA_WIDTH-1 : 0] s00_axis_tdata,
	input wire [(C_S00_AXIS_TDATA_WIDTH/8)-1 : 0] s00_axis_tstrb,
	input wire  s00_axis_tlast,
	input wire  s00_axis_tvalid,

	input wire  m00_axis_aclk,
	input wire  m00_axis_aresetn,
	output wire  m00_axis_tvalid,
	output wire [C_M00_AXIS_TDATA_WIDTH-1 : 0] m00_axis_tdata,
	output wire [(C_M00_AXIS_TDATA_WIDTH/8)-1 : 0] m00_axis_tstrb,
	output wire  m00_axis_tlast,
	input wire  m00_axis_tready
);

    
localparam [{{ .CountersBits }}:0] samples = {{ .Samples }}; // number of samples that I expect from the client
localparam [{{ .CountersBits }}:0] fifodepth = {{ .FifoDepth }}; // Should be multiple of samples
    
localparam [{{ .CountersBits }}:0] bminputs = {{ .InputNum }};  // number of bminputs for each sample
localparam [{{ .CountersBits }}:0] BATCH_IN_ELEMENTS = samples*bminputs;
localparam [{{ .CountersBits }}:0] TOT_IN_ELEMENTS = fifodepth*bminputs;                                     
	
localparam [{{ .CountersBits }}:0] bmoutputs = {{ .OutputNum }};
localparam [{{ .CountersBits }}:0] BATCH_OUT_ELEMENTS = samples*bmoutputs;
localparam [{{ .CountersBits }}:0] TOT_OUT_ELEMENTS = fifodepth*bmoutputs;

// Receiver object
localparam [1:0] IDLE = 1'b0,
                 WRITE_FIFO  = 1'b1; 
	                                    
wire  	   axis_tready;
reg        mst_exec_state;     
wire       fifo_wren;
reg        fifo_full_flag;
reg [{{ .CountersBits }}:0] fifo2bm_write_pointer;
reg [{{ .CountersBits }}:0] fifo2bm_read_pointer;
reg [{{ .CountersBits }}:0] fifo2bm_batch_pointer;
reg [{{ .CountersBits }}:0] fifo2bm_count;
reg        writes_done;

reg fifo_wren_old;
reg  [(C_S00_AXIS_TDATA_WIDTH)-1:0] stream_data_fifo [0 : TOT_IN_ELEMENTS-1];


{{- if .Inputs }}
{{- range .Inputs }}
wire [{{ dec $.Rsize }}:0] {{ . }};
wire {{ . }}_valid;
wire {{ . }}_received;
{{- end }}
{{- end }}

reg fifo2bm_impulse;
wire fifo2bm_ready;
reg [(C_S00_AXIS_TDATA_WIDTH)-1:0] fifo2bm_data; 
	
// Sender objects
localparam [1:0] IDLE_M = 2'b00,                                              
                 SEND_STREAM_M = 2'b01,
                 WAIT_M = 2'b10; 
                                                          
reg [1:0]   mst_exec_state_M;

wire  	axis_tvalid;
reg  	axis_tvalid_delay;
wire  	axis_tlast;
reg  	axis_tlast_delay;
reg [C_M00_AXIS_TDATA_WIDTH-1 : 0] 	stream_data_out;
wire  	tx_en;
reg  	tx_done;

{{- if .Outputs }}
{{- range .Outputs }}
wire [{{ dec $.Rsize }}:0] {{ . }};
wire {{ . }}_valid;
wire {{ . }}_received;
{{- end }}
{{- end }}

reg bm2fifo_ack;
wire bm2fifo_ready;
wire [C_M00_AXIS_TDATA_WIDTH-1 : 0] bm2fifo_data;

reg [{ .CountersBits }}:0] bm2fifo_write_pointer;
reg [{ .CountersBits }}:0] bm2fifo_read_pointer;
reg [{ .CountersBits }}:0] bm2fifo_batch_pointer;
reg [{ .CountersBits }}:0] bm2fifo_count;

reg  [(C_S00_AXIS_TDATA_WIDTH)-1:0] outstream_data_fifo [0 : TOT_OUT_ELEMENTS-1];

integer i,j;

initial begin
	bm2fifo_write_pointer=0;
	bm2fifo_read_pointer=0;
	bm2fifo_batch_pointer=0;
	bm2fifo_count=0;
	bm2fifo_ack=0;
	fifo2bm_data=0;
	fifo2bm_write_pointer=0;
	fifo2bm_read_pointer=0;
	fifo2bm_batch_pointer=0;
	fifo2bm_count=0;
	writes_done=0;
	mst_exec_state=0;
	fifo_wren_old=0;
	fifo2bm_impulse=0;
	fifo2bm_data=0;
	tx_done=0;
	bm2fifo_ack=0;
	for (i=0;i<TOT_OUT_ELEMENTS;i=i+1) begin
		outstream_data_fifo[i]<=0;
	end
	for (j=0;j<TOT_IN_ELEMENTS;j=j+1) begin
		stream_data_fifo[j]<=0;
	end
end

/*
    NOW START THE AXIS SLAVE SECTION, corresponding to BM Inputs
*/

assign m00_axis_tvalid	= axis_tvalid_delay;
assign m00_axis_tdata[(C_S00_AXIS_TDATA_WIDTH)-1:0] = stream_data_out[(C_S00_AXIS_TDATA_WIDTH)-1:0];
assign m00_axis_tlast	= axis_tlast_delay;


assign s00_axis_tready	= axis_tready;

// FSM for the AXI stream receiver
always @(posedge s00_axis_aclk) begin  
	if (!s00_axis_aresetn) begin
		mst_exec_state <= IDLE;
	end  
	else begin
		case (mst_exec_state)
		IDLE: begin
	        	if (s00_axis_tvalid && (fifo_wren == fifo_wren_old) && (fifo2bm_count + BATCH_IN_ELEMENTS <= TOT_IN_ELEMENTS)) begin
	        		mst_exec_state <= WRITE_FIFO;
	            	end
	          	else begin
	              		mst_exec_state <= IDLE;
	            	end
		end
		WRITE_FIFO: begin
			if (writes_done || !s00_axis_tvalid) begin
				mst_exec_state <= IDLE;
			end
			else begin
				mst_exec_state <= WRITE_FIFO;
			end
		end
		endcase
	end
end

assign axis_tready = ((mst_exec_state == WRITE_FIFO) && (fifo2bm_batch_pointer <= BATCH_IN_ELEMENTS-1));

always@(posedge s00_axis_aclk) begin
	if(!s00_axis_aresetn) begin
		fifo2bm_batch_pointer <= 0;
		writes_done <= 1'b0;
	end  
	else begin
		if (!writes_done && fifo2bm_batch_pointer <= BATCH_IN_ELEMENTS-1) begin
			if (fifo_wren) begin
				fifo2bm_batch_pointer <= fifo2bm_batch_pointer + 1;
				writes_done <= 1'b0;
			end

			if (fifo2bm_batch_pointer == BATCH_IN_ELEMENTS-1) begin
				writes_done <= 1'b1;
				fifo2bm_batch_pointer <= 1'b0;
			end
		end
		else begin
			writes_done <= 1'b0;
		end
	end
end 
	       
	       
assign fifo_wren = s00_axis_tvalid && axis_tready;

// old value of fifo_wren	
always @( posedge s00_axis_aclk ) begin
	fifo_wren_old <= fifo_wren;
end
	

    always @( posedge s00_axis_aclk )
    begin
      if (fifo_wren)
        begin
          stream_data_fifo[fifo2bm_write_pointer+fifo2bm_batch_pointer] <= s00_axis_tdata[31:0];
        end  
    end      
	               
    always @( posedge s00_axis_aclk )
    begin
        if (!fifo_wren && fifo_wren_old) begin
            fifo2bm_impulse <= 1'b0;
            if (fifo2bm_write_pointer+BATCH_IN_ELEMENTS == TOT_IN_ELEMENTS) begin
                fifo2bm_write_pointer <= 0;
                fifo2bm_count <= TOT_IN_ELEMENTS - fifo2bm_read_pointer;
            end
            else begin
                fifo2bm_write_pointer <= fifo2bm_write_pointer + BATCH_IN_ELEMENTS;
                if (fifo2bm_write_pointer+BATCH_IN_ELEMENTS > fifo2bm_read_pointer) begin
                    fifo2bm_count <= fifo2bm_write_pointer - fifo2bm_read_pointer + BATCH_IN_ELEMENTS;
                end
                else begin
                    fifo2bm_count <= TOT_IN_ELEMENTS - fifo2bm_read_pointer + fifo2bm_write_pointer + BATCH_IN_ELEMENTS;
                end
            end
        end
        else begin
            if (fifo2bm_count > 0) begin
                if (fifo2bm_ready && !fifo2bm_impulse) begin
                    fifo2bm_impulse <= 1'b1;
                    fifo2bm_data <= stream_data_fifo[fifo2bm_read_pointer];
                    
                    if (fifo2bm_read_pointer+1 == TOT_IN_ELEMENTS) begin
                        fifo2bm_read_pointer <= 0;
                        fifo2bm_count <= fifo2bm_write_pointer;
                    end
                    else begin
                        fifo2bm_read_pointer <= fifo2bm_read_pointer +1;
                        if (fifo2bm_write_pointer < fifo2bm_read_pointer + 1) begin
                            fifo2bm_count <= TOT_IN_ELEMENTS - fifo2bm_read_pointer + fifo2bm_write_pointer - 1;
                        end
                        else begin
                            fifo2bm_count <= fifo2bm_write_pointer - fifo2bm_read_pointer - 1;
                        end
                    end
                end
                else begin
                    fifo2bm_impulse <= 1'b0;
                end
            end
            else begin
                fifo2bm_impulse <= 1'b0;
            end
        end
    end
   
    
bmdeserialize fifo2bm(.clk(m00_axis_aclk),
	.impulse(fifo2bm_impulse),
	.data(fifo2bm_data[31:0]),
	.ready(fifo2bm_ready),
    .reset(!m00_axis_aresetn),
    .i0(i0),
    .i0_valid(i0_valid),
    .i0_recv(i0_received),
    .i1(i1[31:0]),
    .i1_valid(i1_valid),
    .i1_recv(i1_received),
    .i2(i2[31:0]),
    .i2_valid(i2_valid),
    .i2_recv(i2_received),
    .i3(i3[31:0]),
    .i3_valid(i3_valid),
    .i3_recv(i3_received)
);

bondmachine bm(.clk(m00_axis_aclk),
    .reset(!m00_axis_aresetn),
    .i0(i0[31:0]),
    .i0_valid(i0_valid),
    .i0_received(i0_received),
    .i1(i1[31:0]),
    .i1_valid(i1_valid),
    .i1_received(i1_received),
    .i2(i2[31:0]),
    .i2_valid(i2_valid),
    .i2_received(i2_received),
    .i3(i3[31:0]),
    .i3_valid(i3_valid),
    .i3_received(i3_received),
    .o0(o0[31:0]),
    .o0_valid(o0_valid),
    .o0_received(o0_received)
    );

    
bmserialize bm2fifo(.clk(m00_axis_aclk),
	.ack(bm2fifo_ack),
	.data(bm2fifo_data[31:0]),
	.ready(bm2fifo_ready),
    .reset(!m00_axis_aresetn),
    .o0(o0[31:0]),
    .o0_valid(o0_valid),
    .o0_recv(o0_received)
);    
      	               
    /*
        NOW START THE MASTER AXIS SECTION corrponfing to BM outputs
    */


	always @(posedge m00_axis_aclk)                                             
	begin                                                                     
	  if (!m00_axis_aresetn)                                                  
	    begin
	           bm2fifo_count <= 0;                                                            
               bm2fifo_ack <= 1'b0;
               bm2fifo_write_pointer <= 0;
               bm2fifo_read_pointer <= 0;
               // TODO reset fifo                              
	    end                                                                   
	  else
	   begin
	       if (tx_done) begin
	            bm2fifo_ack <= 1'b0;
//	            bm2fifo_count <= bm2fifo_count - BATCH_OUT_ELEMENTS;
                if (bm2fifo_read_pointer+BATCH_OUT_ELEMENTS == TOT_OUT_ELEMENTS) begin
                    bm2fifo_read_pointer <= 0;
                    bm2fifo_count <= bm2fifo_write_pointer;
                end
                else begin
                    bm2fifo_read_pointer <= bm2fifo_read_pointer + BATCH_OUT_ELEMENTS;
                    if (bm2fifo_write_pointer < bm2fifo_read_pointer + BATCH_OUT_ELEMENTS) begin
                        bm2fifo_count <= TOT_OUT_ELEMENTS - bm2fifo_read_pointer - BATCH_OUT_ELEMENTS + bm2fifo_write_pointer;
                    end
                    else
                    begin
                        bm2fifo_count <= bm2fifo_write_pointer - bm2fifo_read_pointer - BATCH_OUT_ELEMENTS;
                    end
                end
	       end
	       else
	       begin
	           if (bm2fifo_ready && !bm2fifo_ack) begin
    	           outstream_data_fifo[bm2fifo_write_pointer] <= bm2fifo_data[31:0];
    	           bm2fifo_ack <= 1'b1;
//	               bm2fifo_count <= bm2fifo_count + 1;
    	           if (bm2fifo_write_pointer + 1 == TOT_OUT_ELEMENTS) begin
    	               bm2fifo_write_pointer <= 0;
    	               bm2fifo_count <= TOT_OUT_ELEMENTS - bm2fifo_read_pointer;
    	           end
    	           else
    	           begin
    	               bm2fifo_write_pointer <= bm2fifo_write_pointer + 1;
    	               if (bm2fifo_write_pointer + 1 > bm2fifo_read_pointer) begin
    	                   bm2fifo_count <= bm2fifo_write_pointer + 1 - bm2fifo_read_pointer;
    	               end
    	               else
    	               begin
    	                   bm2fifo_count <= TOT_OUT_ELEMENTS - bm2fifo_read_pointer + bm2fifo_write_pointer + 1;
    	               end
    	           end
    	       end
    	       else
    	       begin
    	           bm2fifo_ack <= 1'b0;
    	       end
           end
	   end
	end  
                              

	always @(posedge m00_axis_aclk)                                             
	begin                                                                     
	  if (!m00_axis_aresetn)                                                  
	    begin                                                                 
	      mst_exec_state_M <= IDLE_M;                                                                                                   
	    end                                                                   
	  else                                                                    
	    case (mst_exec_state_M)                                                 
	      IDLE_M: 
	            if (bm2fifo_count >= BATCH_OUT_ELEMENTS)                                                        
	               mst_exec_state_M  <= SEND_STREAM_M;                                                        
	      SEND_STREAM_M:                          
	        if (tx_done)                                                      
	          begin                                                           
	            mst_exec_state_M <= IDLE_M;                                       
	          end                                                             
	        else                                                              
	          begin                                                           
	            mst_exec_state_M <= SEND_STREAM_M;                                
	          end                                                             
	    endcase                                                               
	end

	assign axis_tvalid = (mst_exec_state_M == SEND_STREAM_M);
    assign axis_tlast = (bm2fifo_batch_pointer == BATCH_OUT_ELEMENTS-1);

    always @(posedge m00_axis_aclk)                                                                  
	begin        
	if (tx_done) begin
	       axis_tvalid_delay <= 1'b0;                                                               
	      axis_tlast_delay <= 1'b0;         
	end     
	else begin                                                                             
	  if (!m00_axis_aresetn)                                                                         
	    begin                                                                                      
	      axis_tvalid_delay <= 1'b0;                                                               
	      axis_tlast_delay <= 1'b0;                                                                
	    end                                                                                        
	  else                                                                                         
	    begin                                                                                      
	      axis_tvalid_delay <= axis_tvalid;                                                        
	      axis_tlast_delay <= axis_tlast;                                                          
	    end                                                                                        
	end 
	end

    always@(posedge m00_axis_aclk)                                               
	begin     
	  
	  if (tx_done) begin
	       bm2fifo_batch_pointer <= 0;
	       tx_done <= 1'b0;
	  end
	  else begin                                                      
          if(!m00_axis_aresetn)                                                            
          begin                                                                        
              bm2fifo_batch_pointer <= 0;                                                         
              tx_done <= 1'b0;                                                           
          end                                                                          
          else
          begin                                                                           
//            if (bm2fifo_batch_pointer <= BATCH_OUT_ELEMENTS-1)                                
//            begin                                                                      
            if (tx_en)                                                                
              begin                                                                  
                bm2fifo_batch_pointer <= bm2fifo_batch_pointer + 1;                                    
                if (bm2fifo_batch_pointer + 1 == BATCH_OUT_ELEMENTS) 
                begin
                   tx_done <= 1'b1;
                end
//                else
//                begin
//                   tx_done <= 1'b0;
//                end                             
              end                                                                       
//            end 
         end
       end
    end
       
    assign tx_en = m00_axis_tready && axis_tvalid;   
	                                                     
	    always @( posedge m00_axis_aclk )                  
	    begin                                         
	      if(!m00_axis_aresetn)                            
	        begin                                        
	          stream_data_out <= 1;                      
	        end                                          
	      else if (tx_en)
	        begin                           
                stream_data_out[31:0] <= outstream_data_fifo[bm2fifo_read_pointer+bm2fifo_batch_pointer];
	        end                                          
	    end

	endmodule
	
module bmdeserialize(
	   input clk,
	   input reset,
	   input impulse,
	   input [31:0] data,
	   output reg ready,
	   
	   output wire [31:0] i0,
	   output i0_valid,
	   input i0_recv,
	   
	   output wire [31:0] i1,
	   output i1_valid,
	   input i1_recv,
	   
	   output wire [31:0] i2,
	   output i2_valid,
	   input i2_recv,
	   
	   output wire [31:0] i3,
	   output i3_valid,
	   input i3_recv	   	   
	);
	
	reg [2:0] input_index;
    reg [0:0] SM;
    
    reg [31:0] localdata;

    reg [3:0] valids;

    wire [3:0] recvs;
    
    localparam SMIDLE=1'b0,
                SMBM=1'b1;
	
	always @( posedge clk) begin
	   if (reset) begin
	       ready <= 1'b0;
	       input_index <= 3'b000;
	       SM<=SMIDLE;
	       localdata[31:0] <= 32'b0;
	       valids[3:0] <= 4'b0000;
	   end 
	   else begin
	       case (SM)
	       SMIDLE: begin
	               if (impulse) begin
	                   ready <= 1'b0;
	                   localdata[31:0] <= data[31:0];
	                   SM<=SMBM;
	               end
	               else begin
	                   ready <= 1'b1;
	               end
	           end
	       SMBM: begin
	           if (recvs[input_index] == 1'b0) begin
	               valids[input_index] <= 1'b1;
	               ready <= 1'b0;
	           end
	           else begin
	               valids[input_index] <= 1'b0;
	               if (input_index + 1 == 4'd4) begin
	                   input_index <= 0;
	               end
	               else begin
	                   input_index <= input_index + 1;
	               end
	               SM<=SMIDLE;
	               ready <= 1'b1;
	           end
	       end
	       endcase
	   end
	end
	
	assign i0[31:0] = localdata[31:0];
	assign i1[31:0] = localdata[31:0];
	assign i2[31:0] = localdata[31:0];
	assign i3[31:0] = localdata[31:0];
	
	assign i0_valid=valids[0];
	assign i1_valid=valids[1];
	assign i2_valid=valids[2];
	assign i3_valid=valids[3];
	
	assign recvs[3:0] = {i3_recv, i2_recv, i1_recv, i0_recv};
	
	endmodule
	
module bmserialize(
	   input clk,
	   input reset,
	   input ack,
	   output [31:0] data,
	   output reg ready,
	   
	   input [31:0] o0,
	   input o0_valid,
	   output o0_recv   	   
	);
	
	
	reg [1:0] output_index;
    reg [1:0] SM;
    
    reg [31:0] localdata;

    wire [1:0] valids;

    reg [1:0] recvs;
    
    wire [31:0]outputs;
    
    localparam SMIDLE=2'b00,
                SMRES=2'b01,
                SMBM=2'b10;
	
	always @( posedge clk) begin
	   if (reset) begin
	       ready <= 1'b0;
	       output_index <= 2'b00;
	       SM<=SMIDLE;
	       recvs[1:0] <= 2'b00;
	   end 
	   else begin
	       case (SM)
	       SMIDLE: begin
	               if (valids[output_index]) begin
	                   ready <= 1'b1;
	                   localdata[31:0] <= outputs[31:0];
	                   SM<=SMRES;
	               end
	               else begin
	                   ready <= 1'b0;
	               end
	           end
	       SMRES: begin
	           if (ack) begin
	               ready <= 1'b0;
	               SM<=SMBM;
	           end   
	       end
	       SMBM: begin
	           if (!valids[output_index]) begin
	               if (output_index + 1 == 2'd1) begin
	                   output_index <= 0;
	               end
	               else begin
	                   output_index <= output_index + 1;
	               end
	               recvs[output_index] <= 1'b0;
	               SM<=SMIDLE;           
	           end
	           else
	           begin
	               recvs[output_index] <= 1'b1;
	           end
	       end
	       endcase
	   end
	end
	
	assign data[31:0] = localdata[31:0];
	assign outputs={o0[31:0]};
	assign o0_recv=recvs[0];

	assign valids[0] = o0_valid;	
	
endmodule 
`
)
