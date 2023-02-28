package bondmachine

const (
	basicAXIStream = `

	module test01_v1_0 #
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

	assign axis_tready = ((mst_exec_state == WRITE_FIFO) && (write_pointer <= NUMBER_OF_INPUTS-1));

	always@(posedge s00_axis_aclk)
	begin

	  if (tx_done) begin
	       write_pointer <= 0;
	       writes_done <= 1'b0;
	  end
	  else begin
	  if(!s00_axis_aresetn)
	    begin
	      write_pointer <= 0;
	      writes_done <= 1'b0;
	    end  
	  else
	    if (write_pointer <= NUMBER_OF_INPUTS-1)
	      begin
	        if (fifo_wren)
	          begin
	            write_pointer <= write_pointer + 1;
	            writes_done <= 1'b0;
	          end
	          if ((write_pointer == NUMBER_OF_INPUTS-1)|| s00_axis_tlast)
	            begin
	              writes_done <= 1'b1;
	            end
	      end
	      end
	end 


	assign fifo_wren = s00_axis_tvalid && axis_tready;

	reg  [(C_S00_AXIS_TDATA_WIDTH)-1:0] stream_data_fifo [0 : NUMBER_OF_INPUTS-1];
    always @( posedge s00_axis_aclk )
    begin
      if (fifo_wren)
        begin
          stream_data_fifo[write_pointer] <= s00_axis_tdata;
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
    reg  	tx_done;
    wire     bm_done;

    reg  [{{ .CountersBits }}:0] outputs_counter = 0;
    reg  [{{ .CountersBits }}:0] outputs_counter_incr = 0;
	reg  [{{ .CountersBits }}:0] outputs_counter_pointer = 0;
    reg  [{{ .CountersBits }}:0] stream_output_counter = 0;
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
    assign axis_tlast = (stream_output_counter == NUMBER_OF_OUTPUTS - 1);

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
	reg [{{ dec $.Rsize }}:0] {{ . }}_r = 32'b0;
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

    always @( posedge m00_axis_aclk )                  
	begin        
    
        if (tx_done) begin
            outputs_counter <= 0;
            outputs_counter_incr <= 0;
			outputs_counter_pointer <= 0;
			{{- if .Outputs }}
			{{- range .Outputs }}
			{{ . }}_received_r = 1'b0;
			{{ . }}_valid_r = 1'b0;
			{{- end }}
			{{- end }}
        end
        else begin
            if (writes_done && !bm_done) begin
                
			{{- if .Inputs }}
			{{- range $index, $element := .Inputs }}
			{{ $element }}_r <= stream_data_fifo[outputs_counter+{{ $index }}];
			{{- end }}
			{{- end }}

            {{- if .Inputs }}
			{{- range .Inputs }}
			{{ . }}_valid_r = 1'b1;
			{{- end }}
			{{- end }}


			{{- if .Outputs }}
			{{- $outputsLen := len .Outputs }}
			{{- range $i, $output := .Outputs }}
			if ( {{ $output }}_valid && !{{ $output }}_received_r) begin
				{{ $output }}_valid_r <= 1'b1;
				{{ $output }}_received_r <= 1'b1;
				output_stream_data_fifo[outputs_counter_pointer+{{ $i }}] <= {{ $output }};
			end
			{{- end }}
			{{- end }}

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
				{{ . }}_valid_r = 1'b0;
				{{- end }}
				{{- end }}

				outputs_counter_pointer <= outputs_counter_pointer + bmoutputs;
                outputs_counter_incr <= outputs_counter_incr + 1;
                outputs_counter <= bminputs*(outputs_counter_incr+1);

				{{- if .Inputs }}
				{{- range $index, $element := .Inputs }}
				{{ $element }}_valid_r <= 1'b0;
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
			end
			end
	    end
	end

    assign tx_en = m00_axis_tready && axis_tvalid;   

    always @( posedge m00_axis_aclk )                  
    begin        
       if (tx_done) begin
        stream_output_counter <= 0;
        tx_done <= 1'b0;
       end
                                     
      if(!m00_axis_aresetn)                            
        begin   
          stream_data_out <= 1;                      
        end                                          
      else if (tx_en)
        begin
          if (stream_output_counter <= NUMBER_OF_OUTPUTS - 1) begin
              stream_data_out <= output_stream_data_fifo[stream_output_counter];
              stream_output_counter <= stream_output_counter + 1;
              if (stream_output_counter == NUMBER_OF_OUTPUTS - 1) begin
                    tx_done <= 1'b1;
              end
          end
        end                                          
    end

endmodule

`
)
