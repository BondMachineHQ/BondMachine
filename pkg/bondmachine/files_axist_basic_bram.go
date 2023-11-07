package bondmachine

const (
	basicAXIStreamBram = `

	module bram_oneportsync
	#(
			parameter RAM_WIDTH 		= 32,
			parameter RAM_DEPTH         = 256
		)(
		input clk, 
		input [RAM_WIDTH-1:0] addr,
		input [RAM_WIDTH-1:0] waddr,
		input cs_n,
		input wr_n, 
		input rd_n,
		input [RAM_WIDTH-1:0] bram_data_in,
		output reg [RAM_WIDTH-1:0] bram_data_out
	);
    
    (* RAM_STYLE="BLOCK" *)
    reg [RAM_WIDTH-1:0] mem [RAM_DEPTH-1:0];

    always @(posedge clk)
        if (cs_n == 1'b0) begin
            begin
                if (wr_n == 1'b1) mem[(addr)] <= bram_data_in;
                if (rd_n == 1'b1) bram_data_out <= mem[waddr];
            end
        end
    endmodule
	
	module bram
	#(
		parameter RAM_WIDTH 		= 32,
		parameter RAM_ADDR_BITS 	= 9
	)
	(
	input							clock,
	input							ram_enable,
	input							write_enable,
    input 		[31:0]	waddress,
    input 		[31:0]	raddress,
    input 		[31:0] 	input_data,
	output reg 	[31:0] 	output_data
	);
	
   
   (* RAM_STYLE="BLOCK" *)
   reg [RAM_WIDTH-1:0] ram_name [RAM_ADDR_BITS:0];
    
   always @(posedge clock) begin
      if (ram_enable) begin
         if (write_enable) begin
            ram_name[waddress] <= input_data;
         end
            
         output_data <= ram_name[raddress];
      end
    end
    
    
	endmodule

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
		input wire  s00_axis_tlast,
		input wire  s00_axis_tvalid,

		input wire  m00_axis_aclk,
		input wire  m00_axis_aresetn,
		output wire  m00_axis_tvalid,
		output wire [C_M00_AXIS_TDATA_WIDTH-1 : 0] m00_axis_tdata,
		output wire  m00_axis_tlast,
		input wire  m00_axis_tready
	);

	reg tx_done;

    /*
        NOW START THE AXIS SLAVE SECTION
    */

    localparam samples = {{ .Samples }}; // number of samples that I expect from the client
    localparam bminputs = {{ .InputNum }};  // number of bminputs for each sample (or bminputs)
    localparam bmoutputs = {{ .OutputNum }}; // number of output for the classification
	localparam NUMBER_OF_INPUTS  = samples*bminputs;                                     
    localparam NUMBER_OF_OUTPUTS = samples*bmoutputs;
	localparam precision = {{ $.Rsize }}; // precision bit
	localparam maxfifoloop = (C_S00_AXIS_TDATA_WIDTH / precision) - 1;
	localparam waitCycle = 2;
    localparam raddress_size = (precision == 32) ? NUMBER_OF_INPUTS : ((NUMBER_OF_INPUTS % 2 == 0) ? (NUMBER_OF_INPUTS / 2) : (NUMBER_OF_INPUTS / 2 + 1));
    localparam waddressn_size = (precision == 32) ? NUMBER_OF_OUTPUTS : ((NUMBER_OF_OUTPUTS % 2 == 0) ? (NUMBER_OF_OUTPUTS / 2) : (NUMBER_OF_OUTPUTS / 2 + 1));
    localparam waitStream = 10;

	// Machine state for the slave stream part
	parameter [1:0] IDLE = 1'b0,
	                WRITE_FIFO  = 1'b1; 

	wire  	   axis_tready;
	reg        mst_exec_state;     
	wire       fifo_wren;
	reg        fifo_full_flag;
	reg [{{ .CountersBits }}:0] write_pointer;
	reg        writes_done;
    wire       test;

	assign s00_axis_tready	= axis_tready;

	/*
		ADD BLOCK RAM MODULE HERE
	*/
	
    reg							ram_enable = 1;
    reg							input_stream_write_enable = 1;
    reg							input_stream_ram_enable = 1;
    reg 	[C_S00_AXIS_TDATA_WIDTH-1:0]	input_stream_waddress = 0;
    reg 	[C_S00_AXIS_TDATA_WIDTH-1:0]	input_stream_raddress = raddress_size; 
    reg 	[C_S00_AXIS_TDATA_WIDTH-1:0] 	input_stream_data = 0;
    wire	[C_S00_AXIS_TDATA_WIDTH-1:0] 	input_output_stream_data;
    
    reg [31:0] blockram_state_machine = 0;
    reg [31:0] half_precision_input_index = 0; 
    
    
	bram
    #(
        .RAM_WIDTH 		(C_S00_AXIS_TDATA_WIDTH 		),
        .RAM_ADDR_BITS 	(NUMBER_OF_INPUTS 	)
    )
    bram_inst
    (
        .clock			(s00_axis_aclk	),
        .ram_enable		(input_stream_ram_enable		),
        .write_enable	(input_stream_write_enable	),
        .waddress		(input_stream_waddress ),
        .raddress		(input_stream_raddress		),
        .input_data		(input_stream_data		),
        .output_data    (input_output_stream_data	)
    );
    
    reg [31:0] data_in = 32'd0;
    reg [31:0] rwaddress = 0;
    reg [31:0] waddressn = waddressn_size;
    reg  wr_n = 1'b0;
    reg  rd_n = 1'b0;
    wire [31:0] data_out;
    
    
    bram_oneportsync
    #(
        .RAM_WIDTH 		(C_S00_AXIS_TDATA_WIDTH ),
        .RAM_DEPTH 	    (NUMBER_OF_OUTPUTS 	)
    )
    bram_oneportsync_inst
    (
        .clk(s00_axis_aclk), 
        .addr(rwaddress), 
        .waddr(waddressn),
        .cs_n(0),
        .wr_n(wr_n), 
        .rd_n(rd_n),
        .bram_data_in(data_in),
        .bram_data_out(data_out)
    );

	always @(posedge s00_axis_aclk) 
	begin  
	  if (!s00_axis_aresetn) 
	    begin
	      mst_exec_state <= IDLE;
	    end  
	  else
	    case (mst_exec_state)
	      IDLE:
	          if (s00_axis_tvalid)
	            begin
	              mst_exec_state <= WRITE_FIFO;
	            end
	          else
	            begin
	              mst_exec_state <= IDLE;
	            end
	      WRITE_FIFO:
	        if (writes_done)
	          begin
	            mst_exec_state <= IDLE;
	          end
	        else
	          begin
	            mst_exec_state <= WRITE_FIFO;
	          end

	    endcase
	end

	assign axis_tready = ((mst_exec_state == WRITE_FIFO) && (write_pointer <= (NUMBER_OF_INPUTS/(maxfifoloop+1))-1));

	always@(posedge s00_axis_aclk)
	begin

	  if (tx_done) begin
	       write_pointer <= 0;
	       writes_done <= 1'b0;
		   input_stream_waddress <= 0;
	  end
	  else begin
	  if(!s00_axis_aresetn)
	    begin
	      write_pointer <= 0;
	      writes_done <= 1'b0;
	    end  
	  else
	    if (write_pointer <= (NUMBER_OF_INPUTS/(maxfifoloop+1))-1)
	      begin
	        if (fifo_wren)
	          begin
	            write_pointer <= write_pointer + 1;
				input_stream_waddress <= input_stream_waddress + 1;
	            writes_done <= 1'b0;
	          end
	          if ((write_pointer == (NUMBER_OF_INPUTS/(maxfifoloop+1))-1)|| s00_axis_tlast)
	            begin
	              writes_done <= 1'b1;
	            end
	      end
	      end
	end 


	assign fifo_wren = s00_axis_tvalid && axis_tready;

    always @( posedge s00_axis_aclk )
    begin
      if (fifo_wren)
        begin
          if (precision == 32) begin
			input_stream_data <= s00_axis_tdata;
          end
          else if (precision == 16) begin
			input_stream_data[31:16]  <= s00_axis_tdata[15:0];
			input_stream_data[15:0] <= s00_axis_tdata[31:16];
          end
        end   
     end      

    /*
        NOW START THE MASTER AXIS SECTION
    */
                             
	parameter [1:0] IDLE_M = 2'b00,                                             
	                INIT_COUNTER_M  = 2'b01, 
	                PROCESS_BM = 2'B10,   
	                SEND_STREAM_M   = 2'b11; 

	reg [1:0]   mst_exec_state_M;
    reg [{{ .CountersBits }}:0] 	count;

    wire  	axis_tvalid;
    reg  	axis_tvalid_delay;
    wire  	axis_tlast;
    reg  	axis_tlast_delay;
    reg [C_M00_AXIS_TDATA_WIDTH-1 : 0] 	stream_data_out;
    wire  	tx_en;
    wire     bm_done;

    reg  [{{ .CountersBits }}:0] outputs_counter = 0;
    reg  [{{ .CountersBits }}:0] outputs_counter_incr = 0;
	reg  [{{ .CountersBits }}:0] outputs_counter_pointer = 0;
    reg  [{{ .CountersBits }}:0] stream_output_counter = 0;

	{{- if .Outputs }}
	{{- range .Outputs }}
	reg [1:0] write_{{ . }};
	{{- end }}
	{{- end }}

    assign m00_axis_tvalid	= axis_tvalid_delay;
	assign m00_axis_tdata	= stream_data_out;
	assign m00_axis_tlast	= axis_tlast_delay;

	always @(posedge m00_axis_aclk)                                             
	begin                                                                     
	  if (!m00_axis_aresetn)                                                  
	    begin                                                                 
	      mst_exec_state_M <= IDLE_M;                                             
	      count    <= 0;                                                      
	    end                                                                   
	  else                                                                    
	    case (mst_exec_state_M)                                                 
	      IDLE_M:                                                         
	            mst_exec_state_M  <= INIT_COUNTER_M; 

	      INIT_COUNTER_M:                              
	        if ( count == 32 - 1 )                               
	          begin                                                           
	            mst_exec_state_M  <= PROCESS_BM;                               
	          end                                                             
	        else                                                              
	          begin                                                           
	            count <= count + 1;                                           
	            mst_exec_state_M  <= INIT_COUNTER_M;                              
	          end                                                             

	      PROCESS_BM:
	           if (!bm_done) 
	           begin
	               mst_exec_state_M <= PROCESS_BM;
	           end
	           else
	           begin
	               mst_exec_state_M <= SEND_STREAM_M;   
	           end

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

	reg [31:0] endStream = NUMBER_OF_OUTPUTS/(maxfifoloop+1) - 1;
    reg [31:0] wait_stream = 0;

	assign axis_tvalid = ((mst_exec_state_M == SEND_STREAM_M) && (writes_done) && (bm_done) && (sendData == waitCycle) && (rd_n) && (wait_stream >= waitStream));
    assign axis_tlast = (waddressn == 1);

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

	{{- if .Inputs }}
	{{- range .Inputs }}
	reg [{{ dec $.Rsize }}:0] {{ . }}_r = {{ $.Rsize }}'b0;
	{{- end }}
	{{- end }}

	{{- if .Outputs }}
	{{- range .Outputs }}
	reg [{{ dec $.Rsize }}:0] {{ . }}_received_r;
	{{- end }}
	{{- end }}

	{{- if .Inputs }}
	{{- range .Inputs }}
	wire [{{ dec $.Rsize }}:0] {{ . }};
	wire {{ . }}_valid;
	wire {{ . }}_received;
	{{- end }}
	{{- end }}

	{{- if .Inputs }}
	{{- range .Inputs }}
	reg {{ . }}_valid_r = 1'b0;
	{{- end }}
	{{- end }}

	{{- if .Outputs }}
	{{- range .Outputs }}
	wire [{{ dec $.Rsize }}:0] {{ . }};
	wire {{ . }}_valid;
	wire {{ . }}_received;
	reg  {{ . }}_valid_r = 1'b0;
	{{- end }}
	{{- end }}

	{{- if .Inputs }}
	{{- range .Inputs }}
	assign {{ . }} = {{ . }}_r;
	{{- end }}
	{{- end }}

	{{- if .Inputs }}
	{{- range .Inputs }}
	assign {{ . }}_valid = {{ . }}_valid_r;
	{{- end }}
	{{- end }}

	{{- if .Outputs }}
	{{- range .Outputs }}
	assign {{ . }}_received = {{ . }}_received_r;
	{{- end }}
	{{- end }}

    assign bm_done = (outputs_counter_incr == samples);
	reg[1:0] send = 2'b00;

	bondmachine bm(.clk(m00_axis_aclk),
	.reset(!m00_axis_aresetn),
	{{- if .Inputs }}
	{{- range .Inputs }}
	.{{ . }}({{ . }}),
	.{{ . }}_valid({{ . }}_valid),
	.{{ . }}_received({{ . }}_received),
	{{- end }}
	{{- end }}
	{{- if .Outputs }}
	{{- $outputsLen := len .Outputs }}
	{{- range $i, $output := .Outputs }}
	.{{ $output }}({{ $output }}),
	.{{ $output }}_valid({{ $output }}_valid),
	{{- if eq (inc $i) $outputsLen }}
	.{{ $output }}_received({{ $output }}_received)
	{{- else }}
	.{{ $output }}_received({{ $output }}_received),
	{{- end }}
	{{- end }}
	{{- end }}
	);

	assign tx_en = m00_axis_tready && ((mst_exec_state_M == SEND_STREAM_M) && (writes_done) && (bm_done) && (rd_n) && (wait_stream >= waitStream));  

	reg [31:0] outputIndex = 0;
	always @( posedge m00_axis_aclk )                  
	begin
	   if (tx_done) begin
	       rwaddress <= 0;
	       outputIndex <= 0;
	   end
	   
	   if (mst_exec_state_M == PROCESS_BM) begin
	       if (precision == 32) begin
				{{- if .Outputs }}
				{{- $OutputsLen := len .Outputs }}
				{{- range $i, $output := .Outputs }}
				{{- if eq ($i) 0 }}
				if ({{ . }}_valid && !{{ . }}_received_r) begin
					rwaddress <= rwaddress +1;
				end
				{{- else }}
				else if ({{ . }}_valid && !{{ . }}_received_r) begin
					rwaddress <= rwaddress +1;
				end
				{{- end }}
				{{- end }}
				{{- end }}

           end 
           else if (precision == 16) begin
		   
		   {{- $OutputsLen := len .Outputs }}
		   {{- range $i, $output := .Outputs }}
		   {{- if eq ($OutputsLen) 1 }}
		   if ({{ . }}_valid && !{{ . }}_received_r) begin
		   {{- else }}
		   {{- if eq ($i) 0 }}
		   if (({{ . }}_valid && !{{ . }}_received_r) {{- else }} {{- if eq $i (sub $OutputsLen 1) }} || ({{ . }}_valid && !{{ . }}_received_r)) begin {{- else }} || ({{ . }}_valid && !{{ . }}_received_r)
		   {{- end }}  
		   {{- end }}
		   {{- end }}
		   {{- end }}

		   			if (outputIndex == 1) begin
                        outputIndex <= 0;
                    end 
                    else begin
                        outputIndex <= outputIndex + 1;
                        rwaddress <= rwaddress + 1;
                    end
			end
           end
	   end
	end

	reg [1:0] f = 2'd1;
    reg [31:0] precision_index = 0;
    reg [31:0] prev_value = 0;



	always @( posedge m00_axis_aclk )                  
	begin        
    
        if (tx_done) begin
            outputs_counter <= 0;
            outputs_counter_incr <= 0;
			outputs_counter_pointer <= 0;
			{{- if .Outputs }}
			{{- range .Outputs }}
			{{ . }}_received_r <= 1'b0;
			{{ . }}_valid_r <= 1'b0;
			{{- end }}
			{{- end }}
			blockram_state_machine <= 0;
			rd_n <= 1'b0;
			{{- if .Outputs }}
			{{- range .Outputs }}
			write_{{ . }} <= 0;
			{{- end }}
			input_stream_raddress <= raddress_size;
			send <= 2'b00;
			input_stream_write_enable <= 1;
        end
        else begin
		if (writes_done && !bm_done) begin     
		if (send == 2'b00) begin
			if (precision == 32) begin
					if (blockram_state_machine == 32'd0) begin
						 input_stream_write_enable <= 0;
						 blockram_state_machine <= 32'd1;
					end else if (blockram_state_machine == 32'd1) begin
						 blockram_state_machine <= 32'd2;
						 input_stream_raddress <= input_stream_raddress - 1;
					end {{- if .Inputs }} {{- $inputsLen := len .Inputs }} {{- range $i, $input := .Inputs }} {{- if ne (inc $i) $inputsLen }} else if (blockram_state_machine == 32'd{{ add $i 2 }}) begin
						blockram_state_machine <= 32'd{{ add $i 3 }};
						{{ $input }}_r <= input_output_stream_data;
						input_stream_raddress <= input_stream_raddress - 1;
					end {{- else }} else if (blockram_state_machine == 32'd{{ add $i 2 }}) begin
						blockram_state_machine <= 32'd{{ add $i 3 }};
						{{ $input }}_r <= input_output_stream_data;
					{{- end }}
					{{- end }}
					{{- end }}
					{{- end }}
					end{{ $inputsLength := len .Inputs }} else if (blockram_state_machine == 32'd{{ add $inputsLength 2 }}) begin
						blockram_state_machine <= 32'd{{ add $inputsLength 3 }};
						{{- if .Inputs }}
						{{- range .Inputs }}
						{{ . }}_valid_r <= 1'b1;
						{{- end }}
						{{- end }}
					end
					else begin
						 blockram_state_machine <= 32'd0;
						 send <= 2'b01;
					 end
				end
			else if (precision == 16) begin
					if (blockram_state_machine == 32'd0) begin
						input_stream_write_enable <= 0;
						blockram_state_machine <= 32'd1;
					end else if (blockram_state_machine == 32'd1) begin
						blockram_state_machine <= 32'd2;
						if (f == 1) begin
							f <= 0;
							input_stream_raddress <= input_stream_raddress - 1;
						end
					end {{- if .Inputs }} {{- $inputsLen := len .Inputs }} {{- range $i, $input := .Inputs }} {{- if ne (inc $i) $inputsLen }} else if (blockram_state_machine == 32'd{{ add $i 2 }}) begin
						blockram_state_machine <= 32'd{{ add $i 3 }};
						if (half_precision_input_index == 0) begin
							{{ $input }}_r <= input_output_stream_data[15:0];
							half_precision_input_index <= half_precision_input_index + 1;
                            prev_value <= input_output_stream_data;
						end else begin
							{{ $input }}_r <= prev_value[31:16];
							half_precision_input_index <= 0;
							input_stream_raddress <= input_stream_raddress - 1;
						end
					end {{- else }} else if (blockram_state_machine == 32'd{{ add $i 2 }}) begin
						blockram_state_machine <= 32'd{{ add $i 3 }};
						if (half_precision_input_index == 0) begin
							{{ $input }}_r <= input_output_stream_data[15:0];
							half_precision_input_index <= half_precision_input_index + 1;
							prev_value <= input_output_stream_data;
						end else begin
							{{ $input }}_r <= prev_value[31:16];
							half_precision_input_index <= 0;
							f <= 1;
						end
					end{{- end }} {{- end}} {{- end}} {{ $inputsLength := len .Inputs }} else if (blockram_state_machine == 32'd{{ add $inputsLength 2 }}) begin
						blockram_state_machine <= 32'd{{ add $inputsLength 3 }};
						{{- if .Inputs }}
						{{- range .Inputs }}
						{{ . }}_valid_r <= 1'b1;
						{{- end }}
						{{- end }}
					end
					else begin
						 blockram_state_machine <= 32'd0;
						 send <= 2'b01;
					end
				end
			end
		else if (send == 2'b01) begin
			if (
				{{- if .Inputs }}
				{{- $InputsLen := len .Inputs }}
				{{- range $i, $input := .Inputs }}
				{{- if eq (inc $i) $InputsLen }}
				{{ $input }}_received
				{{- else }}
				{{ $input }}_received &&
				{{- end }}
				{{- end }}
				{{- end }}
			) begin
				{{- if .Inputs }}
				{{- range .Inputs }}
				{{ . }}_valid_r <= 1'b0;
				{{- end }}
				{{- end }}
				send <= 2'b10;
				wr_n <= 1;
			end
		end
		else if (send == 2'b10) begin
				{{- if .Outputs }}
				{{- $OutputsLen := len .Outputs }}
				{{- range $i, $output := .Outputs }}
				{{- if eq ($i) 0 }}
				if ({{ . }}_valid && !{{ . }}_received_r) begin
					{{ . }}_received_r <= 1'b1;
					write_{{ . }} <= 1;
				end
				{{- else }}
				else if ({{ . }}_valid && !{{ . }}_received_r) begin
					{{ . }}_received_r <= 1'b1;
					write_{{ . }} <= 1;
				end
				{{- end }}
				{{- end }}
				{{- end }}
		
				if (precision == 32) begin
				{{- if .Outputs }}
				{{- $OutputsLen := len .Outputs }}
				{{- range $i, $output := .Outputs }}
				{{- if eq ($i) 0 }}
				if (write_{{ . }}) begin
						write_{{ . }} <= 0;
						{{ . }}_valid_r <= 1'b1;
						data_in <= {{ . }};
				end
				{{- else }}
				else if (write_{{ . }}) begin
						write_{{ . }} <= 0;
						{{ . }}_valid_r <= 1'b1;
						data_in <= {{ . }};
				end
				{{- end }}
				{{- end }}
				{{- end }}
				end 
				else if (precision == 16) begin

				{{- if .Outputs }}
				{{- $OutputsLen := len .Outputs }}
				{{- range $i, $output := .Outputs }}
				{{- if eq ($i) 0 }}
				if (write_{{ . }}) begin
						write_{{ . }} <= 0;
						{{ . }}_valid_r <= 1'b1;
						if (outputIndex == 0) begin
							data_in[15:0] <= {{ . }};
						end
						else if(outputIndex == 1) begin
							data_in[31:16] <= {{ . }};
						end
				end
				{{- else }}
				else if (write_{{ . }}) begin
						write_{{ . }} <= 0;
						{{ . }}_valid_r <= 1'b1;
						if (outputIndex == 0) begin
							data_in[15:0] <= {{ . }};
						end
						else if(outputIndex == 1) begin
							data_in[31:16] <= {{ . }};
						end
				end
				{{- end }}
				{{- end}}
				{{- end}}
				end 
				if (
					{{- if .Outputs }}
					{{- $OutputsLen := len .Outputs }}
					{{- range $i, $output := .Outputs }}
					{{- if eq (inc $i) $OutputsLen }}
					{{ $output }}_valid_r
					{{- else }}
					{{ $output }}_valid_r &&
					{{- end }}
					{{- end }}
					{{- end }}
					 ) begin
				{{- if .Outputs }}
						{{- range .Outputs }}
							{{ . }}_valid_r <= 1'b0;
						{{- end }}
						{{- end }}

				outputs_counter_pointer <= outputs_counter_pointer + bmoutputs;
				outputs_counter_incr <= outputs_counter_incr + 1;
				outputs_counter <= bminputs*(outputs_counter_incr+1);
			
				 end
				else if(
					{{- if .Outputs }}
					{{- $outputsLen := len .Outputs }}
					{{- range $i, $output := .Outputs }}
					{{- if eq (inc $i) $outputsLen }}
					!{{ $output }}_valid && {{ $output }}_received_r
					{{- else }}
					!{{ $output }}_valid && {{ $output }}_received_r &&
					{{- end }}
					{{- end }}
					{{- end }}
				) 
					begin
						{{- if .Outputs }}
						{{- range .Outputs }}
							{{ . }}_received_r <= 1'b0;
						{{- end }}
						{{- end }}
							
							send <= 2'b00;
					end
			end
		end
		//end 
		else
		begin
		 if (bm_done) begin
			 wr_n <= 1'b0;
			 rd_n <= 1'b1;
		 end
		end
	    end
	end
	
	reg [31:0] sendData = waitCycle;
    
	always @( posedge m00_axis_aclk )                  
    begin        
       if (tx_done) begin
        stream_output_counter <= 0;
        tx_done <= 1'b0;
        waddressn <= waddressn_size;
        sendData <= waitCycle;
        wait_stream <= 0;
       end
      
      if (rd_n == 1) begin
            if (wait_stream < waitStream) begin
                wait_stream <= wait_stream + 1;
            end
      end
                                     
      if(!m00_axis_aresetn)                            
        begin   
          stream_data_out <= 1;                      
        end                                          
      else if (tx_en)
        begin
              if (waddressn >= 0) begin
                  if (precision == 32) begin
                        if (sendData < waitCycle) begin
                               if(sendData == 0) begin
                                    waddressn <= waddressn - 1;
                                end
                                sendData <= sendData + 1;
                        end
                       else if (sendData == waitCycle) begin
                            sendData <= 0;
                            stream_data_out <= data_out;
                        end
                  end
                  else if (precision == 16) begin
                        if (sendData < waitCycle) begin
                               if(sendData == 0) begin
                                    waddressn <= waddressn - 1;
                                end
                                sendData <= sendData + 1;
                         end
                       else if (sendData == waitCycle) begin
                            sendData <= 0;
                            stream_data_out[15:0] <= data_out[15:0];
                            stream_data_out[31:16] <= data_out[31:16];
                       end
                  end
                  if (waddressn == 0) begin
                        tx_done <= 1'b1;
                  end
            end
        end                                          
    end

endmodule

`
)