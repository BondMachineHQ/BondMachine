package bondmachine

const (
	basicAXIStreamBram = `
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
 
    /*
        NOW START THE AXIS SLAVE SECTION
    */

 
    localparam samples = {{ .Samples }}; // number of samples that I expect from the client
    localparam bminputs = {{ .InputNum }};  // number of bminputs for each sample (or bminputs)
    localparam bmoutputs = {{ .OutputNum }}; // number of output for the classification
	localparam NUMBER_OF_INPUTS  = samples*bminputs;                                     
    localparam NUMBER_OF_OUTPUTS = samples*bmoutputs;
	{{- if le $.Rsize 8 }}
	localparam precision = 8;
	{{- else if le $.Rsize 16 }}
	localparam precision = 16;
	{{- else }}
	localparam precision = 32;
	{{- end }}
	localparam maxfifoloop = (C_S00_AXIS_TDATA_WIDTH / precision) - 1;
 
	// Machine state for the slave stream part
	parameter [1:0] IDLE = 1'b0,
	                WRITE_FIFO  = 1'b1; 
 
	wire  	   axis_tready;
	reg        mst_exec_state;     
	wire       fifo_wren;
	reg        fifo_full_flag;
	reg        writes_done;
    wire       test;
	reg [31:0] read_pointer;
	reg [31:0] read_pointer_output;
 
	reg [31:0] read_state;
	reg [31:0] out_read_state;
	reg  [31:0] outputs_counter = 0;
    reg  [31:0] outputs_counter_incr = 0;
	reg  [31:0] outputs_counter_pointer = 0;
    reg  [31:0] stream_output_counter = 0;
	reg  [31:0] maxfifoloopcounter = 0;
	reg [31:0] write_pointer;
 
	assign s00_axis_tready	= axis_tready;
 
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
	  end
	  else begin
	  if(!s00_axis_aresetn)
	    begin
	      writes_done <= 1'b0;
	    end  
	  else
	    if (write_pointer <= (NUMBER_OF_INPUTS/(maxfifoloop+1))-1)
	      begin
	        if (fifo_wren)
	          begin
	            write_pointer <= write_pointer + 1;
	            writes_done <= 1'b0;
	          end
	          if ((write_pointer == (NUMBER_OF_INPUTS/(maxfifoloop+1))-1)|| s00_axis_tlast)
	            begin
	              writes_done <= 1'b1;
	            end
	      end
	      end
	end 
 

	localparam NUM_WORDS = NUMBER_OF_INPUTS / (C_S00_AXIS_TDATA_WIDTH / precision);
	assign fifo_wren = s00_axis_tvalid && axis_tready;
 
	(* ram_style = "block" *)
    reg [(C_S00_AXIS_TDATA_WIDTH)-1:0] stream_data_fifo [0 : NUM_WORDS - 1];
 
	(* ram_style = "block" *)
    reg [(C_S00_AXIS_TDATA_WIDTH)-1:0] stream_data_fifo_backup [0 : NUM_WORDS - 1];

	(* ram_style = "block" *)
    reg [(C_S00_AXIS_TDATA_WIDTH)-1:0] stream_data_fifo_backup_2 [0 : NUM_WORDS - 1];

	(* ram_style = "block" *)
    reg [(C_S00_AXIS_TDATA_WIDTH)-1:0] stream_data_fifo_backup_3 [0 : NUM_WORDS - 1];
 
 
	always @( posedge s00_axis_aclk )
    begin
		if (tx_done) begin
	       maxfifoloopcounter <= 0;
	  end
      if (fifo_wren)
        begin
          if (precision == 32) begin
				stream_data_fifo[write_pointer] <= s00_axis_tdata;
          end
          else if (precision == 16) begin
		  	stream_data_fifo[write_pointer][15:0]   <= s00_axis_tdata[15:0];
			stream_data_fifo_backup[write_pointer][15:0]  <= s00_axis_tdata[31:16];
          end else if (precision == 8) begin
		 	 stream_data_fifo[write_pointer][7:0]   <= s00_axis_tdata[7:0];
			 stream_data_fifo_backup[write_pointer][7:0]  <= s00_axis_tdata[15:8];
			 stream_data_fifo_backup_2[write_pointer][7:0]  <= s00_axis_tdata[23:16];
			 stream_data_fifo_backup_3[write_pointer][7:0]  <= s00_axis_tdata[31:24];
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
    wire  	tx_en;
	reg [C_M00_AXIS_TDATA_WIDTH-1 : 0] 	stream_data_out;
	reg  	tx_done;
    wire     bm_done;
 
	(* ram_style = "block" *)
    reg  [(C_S00_AXIS_TDATA_WIDTH)-1:0] output_stream_data_fifo [0 : NUMBER_OF_OUTPUTS-1];
 
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
 
	assign axis_tvalid = ((mst_exec_state_M == SEND_STREAM_M) && (writes_done) && (bm_done));
    assign axis_tlast = (stream_output_counter == (NUMBER_OF_OUTPUTS/(maxfifoloop+1)) - 1);
 
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
 
	reg [2:0] counter;
	reg [15:0] output_mutex = {{ if eq $.Rsize 8 }}0{{ else }}1{{ end }};
	reg [15:0] input_reader = 1;
	reg [15:0] input_reader_index = 0;
 
	always @( posedge m00_axis_aclk )                  
	begin        
        if (tx_done) begin
            outputs_counter <= 0;
            outputs_counter_incr <= 0;
			outputs_counter_pointer <= 0;
			{{- if .Outputs }}
            {{- range $i, $output := .Outputs }}           
			{{ $output }}_received_r <= 32'b0;
            {{ $output }}_valid_r <= 32'b0;            
			{{- end }}
            {{- end }}
			read_pointer <= 1'b0;
			read_pointer_output  <= 1'b0;
			output_mutex <= {{ if eq $.Rsize 8 }}0{{ else }}1{{ end }};
			input_reader_index <= 0;
        end
        else begin
            if (writes_done && !bm_done) begin
				if (send == 2'b00) begin
					send <= 2'b01;
				end
				else if (send == 2'b01) begin
 
					{{ if eq $.Rsize 32 }}
						case (read_state)
 
							{{- if .Inputs }}
							{{- $inputsLen := len .Inputs }}
							{{- range $i, $input := .Inputs }}
							{{- if ne (inc $i) $inputsLen }}
							32'd{{ $i }}: begin
								{{ $input }}_r <= stream_data_fifo[outputs_counter+{{ $i }}];
								read_pointer <= read_pointer + 1;	
								//{{ $input }}_valid_r <= 1'b1;
								read_state <= 32'd{{ (inc $i) }};
							end
							{{- else }}
							32'd{{ $i }}: begin
								{{ $input }}_r <= stream_data_fifo[outputs_counter+{{ $i }}];
								read_pointer <= read_pointer + 1;	
								//{{ $input }}_valid_r <= 1'b1;
								read_state <= 32'd0;
							end
							{{- end }}
							{{- end }}
							{{- end }}
						endcase
 
						if (read_pointer >= (bminputs)) begin
							read_state <= 32'd0;
							read_pointer <= 0;
							send <= 2'b10;
							{{- if .Inputs }}
							{{- $inputsLen := len .Inputs }}
							{{- range $i, $input := .Inputs }}
							{{ $input }}_valid_r <= 1'b1;
							{{- end }}
							{{- end }}
						end
 
					{{- end }}
					{{ if eq $.Rsize 16 }}
						{{- if eq (len .Inputs) 1 }}
						case (read_state)
							32'd0: begin
								if (outputs_counter_incr[0] == 1'b0) begin
									{{ index .Inputs 0 }}_r <= stream_data_fifo[outputs_counter];
								end else begin
									{{ index .Inputs 0 }}_r <= stream_data_fifo_backup[outputs_counter];
								end
								read_pointer <= read_pointer + 1;	
								read_state <= 32'd1;
							end
						endcase
						{{- else }}
						case (read_state)

							{{- if .Inputs }}
							{{- $inputsLen := len .Inputs }}
							{{- range $i, $input := .Inputs }}

							32'd{{ $i }}: begin
								{{- if eq (mod $i 2) 0 }}
									{{ $input }}_r <= stream_data_fifo[outputs_counter+{{ div $i 2 }}];
								{{- else }}
									{{ $input }}_r <= stream_data_fifo_backup[outputs_counter+{{ div $i 2 }}];
								{{- end }}

								read_pointer <= read_pointer + 1;	
								read_state <= 32'd{{ (inc $i) }};
							end

							{{- end }}
							{{- end }}
						endcase
						{{- end }}

						if (read_pointer >= (bminputs)) begin
							read_state <= 32'd0;
							read_pointer <= 0;
							send <= 2'b10;
							{{- if .Inputs }}
							{{- $inputsLen := len .Inputs }}
							{{- range $i, $input := .Inputs }}
							{{ $input }}_valid_r <= 1'b1;
							{{- end }}
							{{- end }}
						end
					{{- end }}

					{{ if eq $.Rsize 8 }}
					{{- if eq (len .Inputs) 1 }}
					// Special case for single input: cycle through 4 FIFOs based on sample index
					case (read_state)
						32'd0: begin
							case (outputs_counter_incr[1:0])
								2'b00: {{ index .Inputs 0 }}_r <= stream_data_fifo[outputs_counter][7:0];
								2'b01: {{ index .Inputs 0 }}_r <= stream_data_fifo_backup[outputs_counter][7:0];
								2'b10: {{ index .Inputs 0 }}_r <= stream_data_fifo_backup_2[outputs_counter][7:0];
								2'b11: {{ index .Inputs 0 }}_r <= stream_data_fifo_backup_3[outputs_counter][7:0];
							endcase
							read_pointer <= read_pointer + 1;
							read_state <= 32'd1;
						end
					endcase
					{{- else }}
					// Multiple inputs: use global input index
					case (read_state)
						{{- if .Inputs }}
						{{- $inputsLen := len .Inputs }}
						{{- range $i, $input := .Inputs }}
						32'd{{ $i }}: begin
							// Global input index = outputs_counter_incr * bminputs + {{ $i }}
							// For 8-bit: address = global_index / 4, selector = global_index % 4
							case ((outputs_counter_incr * bminputs + {{ $i }}) % 4)
								2'd0: {{ $input }}_r <= stream_data_fifo[(outputs_counter_incr * bminputs + {{ $i }}) / 4][7:0];
								2'd1: {{ $input }}_r <= stream_data_fifo_backup[(outputs_counter_incr * bminputs + {{ $i }}) / 4][7:0];
								2'd2: {{ $input }}_r <= stream_data_fifo_backup_2[(outputs_counter_incr * bminputs + {{ $i }}) / 4][7:0];
								2'd3: {{ $input }}_r <= stream_data_fifo_backup_3[(outputs_counter_incr * bminputs + {{ $i }}) / 4][7:0];
							endcase
							read_pointer <= read_pointer + 1;
							read_state <= 32'd{{ inc $i }};
						end
						{{- end }}
						{{- end }}
					endcase
					{{- end }}

					if (read_pointer >= bminputs) begin
						read_state <= 32'd0;
						read_pointer <= 0;
						send <= 2'b10;
						{{- if .Inputs }}
						{{- range .Inputs }}
						{{ . }}_valid_r <= 1'b1;
						{{- end }}
						{{- end }}
					end
				{{- end }}
				end
				else if (send == 2'b10) begin
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
 
						send <= 2'b11;
					end
				end
				else if (send == 2'b11) begin

					{{ if eq $.Rsize 32 }}
					if (precision == 32) begin
						case (out_read_state)
 
							{{- if .Outputs }}
							{{- $outputsLen := len .Outputs }}
							{{- range $i, $output := .Outputs }}
							{{- if ne (inc $i) $outputsLen }}
							32'd{{ $i }}: begin
								if ( {{ $output }}_valid && !{{ $output }}_received_r) begin
									{{ $output }}_valid_r <= 1'b1;
									{{ $output }}_received_r <= 1'b1;
									output_stream_data_fifo[read_pointer_output] <= {{ $output }};
									read_pointer_output <= read_pointer_output + 1;
									out_read_state <= 32'd{{ (inc $i) }};
								end
							end
							{{- else }}
							32'd{{ $i }}: begin
								if ( {{ $output }}_valid && !{{ $output }}_received_r) begin
									{{ $output }}_valid_r <= 1'b1;
									{{ $output }}_received_r <= 1'b1;
									output_stream_data_fifo[read_pointer_output] <= {{ $output }};
									read_pointer_output <= read_pointer_output + 1;
									out_read_state <= 32'd0;
								end
							end
							{{- end }}
							{{- end }}
							{{- end }}
						endcase
					end
					{{ end }}
					{{ if eq $.Rsize 16 }}
					if (precision == 16) begin
						case (out_read_state)
 
							{{- if .Outputs }}
							{{- $outputsLen := len .Outputs }}
							{{- range $i, $output := .Outputs }}
							{{- if ne (inc $i) $outputsLen }}
							32'd{{ $i }}: begin
								if ( {{ $output }}_valid && !{{ $output }}_received_r) begin
									{{ $output }}_valid_r <= 1'b1;
									{{ $output }}_received_r <= 1'b1;
									if (output_mutex == 1) begin
										output_stream_data_fifo[read_pointer_output][15:0] <= {{ $output }};
									end
									else begin 
										output_stream_data_fifo[read_pointer_output][31:16] <= {{ $output }};
									end
 
									out_read_state <= 32'd{{ (inc $i) }};
 
									if (output_mutex == 2) begin
										read_pointer_output <= read_pointer_output + 1;
										output_mutex <= 1;
									end 
									else begin 
										output_mutex <= output_mutex + 1;
									end
								end
							end
							{{- else }}
							32'd{{ $i }}: begin
								if ( {{ $output }}_valid && !{{ $output }}_received_r) begin
									{{ $output }}_valid_r <= 1'b1;
									{{ $output }}_received_r <= 1'b1;
									if (output_mutex == 1) begin
										output_stream_data_fifo[read_pointer_output][15:0] <= {{ $output }};
									end
									else begin 
										output_stream_data_fifo[read_pointer_output][31:16] <= {{ $output }};
									end
 
									out_read_state <= 32'd0;
 
									if (output_mutex == 2) begin
										read_pointer_output <= read_pointer_output + 1;
										output_mutex <= 1;
									end 
									else begin 
										output_mutex <= output_mutex + 1;
									end
								end
							end
							{{- end }}
							{{- end }}
							{{- end }}
						endcase
					end
					{{ end }}
					{{ if eq $.Rsize 8 }}
					if (precision == 8) begin
						case (out_read_state)

							{{- if .Outputs }}
							{{- $outputsLen := len .Outputs }}
							{{- range $i, $output := .Outputs }}
							{{- if ne (inc $i) $outputsLen }}
							32'd{{ $i }}: begin
								if ( {{ $output }}_valid && !{{ $output }}_received_r) begin
									{{ $output }}_valid_r <= 1'b1;
									{{ $output }}_received_r <= 1'b1;

									// Store the output value into an appropriate 8-bit slot
									// if (output_mutex == 0) begin
									// 	output_stream_data_fifo[read_pointer_output][7:0] <= {{ $output }};
									// end else if (output_mutex == 1) begin
									// 	output_stream_data_fifo[read_pointer_output][15:8] <= {{ $output }};
									// end else if (output_mutex == 2) begin
									// 	output_stream_data_fifo[read_pointer_output][23:16] <= {{ $output }};
									// end else if (output_mutex == 3) begin
									// 	output_stream_data_fifo[read_pointer_output][31:24] <= {{ $output }};
									// end

									output_stream_data_fifo[read_pointer_output] <= { 
										(output_mutex == 3 ? {{ $output }} : output_stream_data_fifo[read_pointer_output][31:24]),
										(output_mutex == 2 ? {{ $output }} : output_stream_data_fifo[read_pointer_output][23:16]),
										(output_mutex == 1 ? {{ $output }} : output_stream_data_fifo[read_pointer_output][15:8]),
										(output_mutex == 0 ? {{ $output }} : output_stream_data_fifo[read_pointer_output][7:0])
									};
									// else begin
									// 	output_stream_data_fifo[read_pointer_output][23:16] <= {{ $output }};
									// end

									out_read_state <= 32'd{{ (inc $i) }};

									if (output_mutex == 3) begin
										read_pointer_output <= read_pointer_output + 1;
										output_mutex <= 0;
									end 
									else begin 
										output_mutex <= output_mutex + 1;
									end
								end
							end
							{{- else }}
							32'd{{ $i }}: begin
								if ( {{ $output }}_valid && !{{ $output }}_received_r) begin
									{{ $output }}_valid_r <= 1'b1;
									{{ $output }}_received_r <= 1'b1;

									// Store the output value into an appropriate 8-bit slot
									// if (output_mutex == 0) begin
									// 	output_stream_data_fifo[read_pointer_output][7:0] <= {{ $output }};
									// end else if (output_mutex == 1) begin
									// 	output_stream_data_fifo[read_pointer_output][15:8] <= {{ $output }};
									// end else if (output_mutex == 2) begin
									// 	output_stream_data_fifo[read_pointer_output][23:16] <= {{ $output }};
									// end else if (output_mutex == 3) begin
									// 	output_stream_data_fifo[read_pointer_output][31:24] <= {{ $output }};
									// end
									output_stream_data_fifo[read_pointer_output] <= { 
										(output_mutex == 3 ? {{ $output }} : output_stream_data_fifo[read_pointer_output][31:24]),
										(output_mutex == 2 ? {{ $output }} : output_stream_data_fifo[read_pointer_output][23:16]),
										(output_mutex == 1 ? {{ $output }} : output_stream_data_fifo[read_pointer_output][15:8]),
										(output_mutex == 0 ? {{ $output }} : output_stream_data_fifo[read_pointer_output][7:0])
									};
									// else begin
									// 	output_stream_data_fifo[read_pointer_output][23:16] <= {{ $output }};
									// end

									//output_stream_data_fifo[read_pointer_output][7:0] <= {{ $output }};

									out_read_state <= 32'd0;

									if (output_mutex == 3) begin
										read_pointer_output <= read_pointer_output + 1;
										output_mutex <= 0;
									end 
									else begin 
										output_mutex <= output_mutex + 1;
									end
								end
							end
							{{- end }}
							{{- end }}
							{{- end }}
						endcase
					end
					{{ end }}


					if ( 
 
				{{- if .Outputs }}
				{{- $outputsLen := len .Outputs }}
				{{- range $i, $output := .Outputs }}
				{{- if eq (inc $i) $outputsLen }}
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
					outputs_counter_pointer <= outputs_counter_pointer + bmoutputs;
					outputs_counter_incr <= outputs_counter_incr + 1;
 
					if (precision == 32) begin
						outputs_counter <= bminputs*(outputs_counter_incr+1);
					end else if (precision == 16) begin
						outputs_counter <= (bminputs*(outputs_counter_incr+1)) / 2;
					end else if (precision == 8) begin
						outputs_counter <= (bminputs*(outputs_counter_incr+1)) / 4;
					end
					send <= 2'b01;
			end
			end
			end
	    end
	end
 
	assign tx_en = m00_axis_tready && axis_tvalid;  
 
	reg [10:0] maxfifoloopcounteroutput = 0;
 
	always @( posedge m00_axis_aclk )                  
    begin        
       if (tx_done) begin
		maxfifoloopcounteroutput <= 0;
        stream_output_counter <= 0;
        tx_done <= 1'b0;
       end
 
      if(!m00_axis_aresetn)                            
        begin   
          stream_data_out <= 1;                      
        end                                          
      else if (tx_en)
        begin
           if (stream_output_counter <= (NUMBER_OF_OUTPUTS/(maxfifoloop+1)) - 1) begin
              if (precision == 32) begin
                    stream_data_out <= output_stream_data_fifo[stream_output_counter];
              end
              else if (precision == 16) begin
					stream_data_out[15:0] <= output_stream_data_fifo[stream_output_counter][15:0];
					stream_data_out[31:16] <= output_stream_data_fifo[stream_output_counter][31:16];
              end
			  else if (precision == 8) begin
					stream_data_out[7:0] <= output_stream_data_fifo[stream_output_counter][7:0];
					stream_data_out[15:8] <= output_stream_data_fifo[stream_output_counter][15:8];
					stream_data_out[23:16] <= output_stream_data_fifo[stream_output_counter][23:16];
					stream_data_out[31:24] <= output_stream_data_fifo[stream_output_counter][31:24];
              end
              stream_output_counter <= stream_output_counter + 1;
              if (stream_output_counter == (NUMBER_OF_OUTPUTS/(maxfifoloop+1)) - 1) begin
                    tx_done <= 1'b1;
              end
          end
        end                                          
    end
 
endmodule
 
`
)
