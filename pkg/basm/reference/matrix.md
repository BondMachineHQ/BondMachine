# Support Matrix

The following tables show the feature state in term of level of development for each instruction and instruction group.
For each of them the support of the features is shown.

The features are the following:
| Feature | Description |
| --- | --- |
| hdl | The instruction can be translated to hardware description language |
| asm | The instruction can be assembled by the assembler |
| disasm | The instruction can be disassembled by the disassembler |
| hlasm | The instruction can be assembled by the high-level assembler (Basm) |
| asmeta | The instruction has metadata for the assembler |
| gosim | The instruction can be simulated in the Go-based simulator |
| hdlsim | The instruction can be simulated in the hardware description language simulator |
| mt | Instruction thread support |

The possible support values are shown below:

| Value | Meaning |
| --- | --- |
| ![ok](iconok.png) | The feature is fully implemented |
| ![no](iconno.png) | The feature is not yet implemented |
| ![testing](icontesting.png) | The feature is being tested |
| ![partial](iconpartial.png) | The feature is partially implemented |
| ![notapplicable](iconnotapplicable.png) | The feature is not applicable to the instruction |

## Support Matrix for Static Instructions

| Instruction | disasm | gosim | mt | hlasm | sim | hdl | asmeta | asm |
| --- | --- | --- | --- | --- | --- | --- | --- | --- |
| [adc](adc.md) | - | - | - | - | - | - | - | - |
| [addf16](addf16.md) | - | - | - | - | - | - | - | - |
| [addf](addf.md) | - | - | - | - | - | - | - | - |
| [add](add.md) | - | - | - | - | - | - | - | - |
| [addi](addi.md) | - | - | - | - | - | - | - | - |
| [addp](addp.md) | - | - | - | - | - | - | - | - |
| [and](and.md) | - | - | - | - | - | - | - | - |
| [chc](chc.md) | - | - | - | - | - | - | - | - |
| [chw](chw.md) | - | - | - | - | - | - | - | - |
| [cilc](cilc.md) | - | - | - | - | - | - | - | - |
| [cil](cil.md) | - | - | - | - | - | - | - | - |
| [cir](cir.md) | - | - | - | - | - | - | - | - |
| [cirn](cirn.md) | - | - | - | - | - | - | - | - |
| [clc](clc.md) | - | - | - | - | - | - | - | - |
| [clr](clr.md) | - | - | - | - | - | - | - | - |
| [cmpr](cmpr.md) | - | - | - | - | - | - | - | - |
| [cmprlt](cmprlt.md) | - | - | - | - | - | - | - | - |
| [cmpv](cmpv.md) | - | - | - | - | - | - | - | - |
| [cpy](cpy.md) | - | - | - | - | - | - | - | ![ok](iconok.png) |
| [cset](cset.md) | - | - | - | - | - | - | - | - |
| [dec](dec.md) | - | - | - | - | - | - | - | - |
| [divf16](divf16.md) | - | - | - | - | - | - | - | - |
| [divf](divf.md) | - | - | - | - | - | - | - | - |
| [div](div.md) | - | - | - | - | - | - | - | - |
| [divp](divp.md) | - | - | - | - | - | - | - | - |
| [dpc](dpc.md) | - | - | - | - | - | - | - | - |
| [expf](expf.md) | - | - | - | - | - | - | - | - |
| [hit](hit.md) | - | - | - | - | - | - | - | - |
| [hlt](hlt.md) | - | - | - | - | - | - | - | - |
| [i2r](i2r.md) | - | - | - | - | - | - | - | - |
| [i2rw](i2rw.md) | - | ![ok](iconok.png) | - | - | - | - | - | ![ok](iconok.png) |
| [incc](incc.md) | - | - | - | - | - | - | - | - |
| [inc](inc.md) | ![ok](iconok.png) | ![ok](iconok.png) | ![ok](iconok.png) | ![ok](iconok.png) | - | ![ok](iconok.png) | ![ok](iconok.png) | ![ok](iconok.png) |
| [ja](ja.md) | - | - | - | - | - | - | - | - |
| [jc](jc.md) | - | - | - | - | - | - | - | - |
| [jcmpa](jcmpa.md) | - | - | - | - | - | - | - | - |
| [jcmpl](jcmpl.md) | - | - | - | - | - | - | - | - |
| [jcmpo](jcmpo.md) | - | - | - | - | - | - | - | - |
| [jcmpria](jcmpria.md) | - | - | - | - | - | - | - | - |
| [jcmprio](jcmprio.md) | - | - | - | - | - | - | - | - |
| [je](je.md) | - | - | - | - | - | - | - | - |
| [j](j.md) | - | - | - | - | - | - | - | - |
| [jgt0f](jgt0f.md) | - | - | - | - | - | - | - | - |
| [jo](jo.md) | - | - | - | - | - | - | - | - |
| [jria](jria.md) | - | - | - | - | - | - | - | - |
| [jri](jri.md) | - | - | - | - | - | - | - | - |
| [jrio](jrio.md) | - | - | - | - | - | - | - | - |
| [jz](jz.md) | - | - | - | - | - | - | - | - |
| [k2r](k2r.md) | - | - | - | - | - | - | - | - |
| [lfsr82r](lfsr82r.md) | - | - | - | - | - | - | - | - |
| [m2r](m2r.md) | - | - | - | - | - | - | - | - |
| [m2rri](m2rri.md) | - | - | - | - | - | - | - | - |
| [mod](mod.md) | - | - | - | - | - | - | - | - |
| [mulc](mulc.md) | - | - | - | - | - | - | - | - |
| [multf16](multf16.md) | - | - | - | - | - | - | - | - |
| [multf](multf.md) | - | - | - | - | - | - | - | - |
| [mult](mult.md) | - | - | - | - | - | - | - | - |
| [multp](multp.md) | - | - | - | - | - | - | - | - |
| [nand](nand.md) | - | - | - | - | - | - | - | - |
| [nop](nop.md) | ![ok](iconok.png) | ![ok](iconok.png) | ![testing](icontesting.png) | ![ok](iconok.png) | - | ![ok](iconok.png) | ![ok](iconok.png) | ![ok](iconok.png) |
| [nor](nor.md) | - | - | - | - | - | - | - | - |
| [not](not.md) | - | - | - | - | - | - | - | - |
| [or](or.md) | - | - | - | - | - | - | - | - |
| [q2r](q2r.md) | - | - | - | - | - | - | - | - |
| [r2m](r2m.md) | - | - | - | - | - | - | - | - |
| [r2mri](r2mri.md) | - | - | - | - | - | - | - | - |
| [r2o](r2o.md) | - | - | - | - | - | - | - | - |
| [r2owaa](r2owaa.md) | - | - | - | - | - | - | - | - |
| [r2owa](r2owa.md) | - | ![ok](iconok.png) | - | - | - | - | - | - |
| [r2q](r2q.md) | - | - | - | - | - | - | - | - |
| [r2s](r2s.md) | - | - | - | - | - | - | - | - |
| [r2t](r2t.md) | - | - | - | - | - | - | - | - |
| [r2u](r2u.md) | - | - | - | - | - | - | - | - |
| [r2v](r2v.md) | - | - | - | - | - | - | - | - |
| [r2vri](r2vri.md) | - | - | - | - | - | - | - | - |
| [ro2r](ro2r.md) | - | - | - | - | - | - | - | - |
| [ro2rri](ro2rri.md) | - | - | - | - | - | - | - | - |
| [rsc](rsc.md) | - | - | - | - | - | - | - | - |
| [rset](rset.md) | - | - | - | - | - | - | - | - |
| [s2r](s2r.md) | - | - | - | - | - | - | - | - |
| [saj](saj.md) | - | - | - | - | - | - | - | - |
| [sbc](sbc.md) | - | - | - | - | - | - | - | - |
| [sic](sic.md) | - | - | - | - | - | - | - | - |
| [sicv2](sicv2.md) | - | - | - | ![notapplicable](iconnotapplicable.png) | ![testing](icontesting.png) | - | ![notapplicable](iconnotapplicable.png) | ![ok](iconok.png) |
| [sicv3](sicv3.md) | ![ok](iconok.png) | ![ok](iconok.png) | - | yes | - | - | - | ![ok](iconok.png) |
| [sub](sub.md) | - | - | - | - | - | - | - | - |
| [t2r](t2r.md) | - | - | - | - | - | - | - | - |
| [tsp](tsp.md) | - | - | - | - | - | - | - | - |
| [u2r](u2r.md) | - | - | - | - | - | - | - | - |
| [wrd](wrd.md) | - | - | - | - | - | - | - | - |
| [wwr](wwr.md) | - | - | - | - | - | - | - | - |
| [xnor](xnor.md) | - | - | - | - | - | - | - | - |
| [xor](xor.md) | - | - | - | - | - | - | - | - |

## Support Matrix for Dynamical Instructions

| Instruction | disasm | gosim | asm |
| --- | --- | --- | --- |
| [call](call.md) | - | - | - |
| [fixed_point](fixed_point.md) | - | - | - |
| [flopoco](flopoco.md) | ![ok](iconok.png) | ![notapplicable](iconnotapplicable.png) | ![ok](iconok.png) |
| [linear_quantizer](linear_quantizer.md) | - | - | - |
| [rsets](rsets.md) | - | - | - |
| [stack](stack.md) | - | - | - |
