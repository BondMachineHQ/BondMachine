# bmnumbers

bmstack is part of BondMachine project. bmnumbers is both a command line tool to convert or cast numbers to and from different formats and a library to do the same. It is used within the BondMachine every time numbers are handled.

## Supported number types

The supported number types are listed in the following table.

| Type Name | Prefixes | Description | Static | Lenght |
| ---- | ------- | ----------- | ------ | ------ |
| unsigned | none <br> 0u <br> 0d | Unsigned integer | yes | any |
| signed | 0s <br> 0sd | Signed integer | yes | any |
| bin | 0b <br> 0b\<s\> | Binary number | yes | any <br> s bits|
| hex | 0x | Hexadecimal number | yes | any |
| float16 | 0f<16> | IEEE 754 half precision floating point number | yes | 16 bits |
| float32 | 0f <br> 0f<32> | IEEE 754 single precision floating point number | yes | 32 bits |
| lqs[s]t[t] | 0lq\<s.t\> | Linear quantized number with size s and type t | no | s bits |
| fps[s]f[f] | 0fp\<s.f\> | Fixed point number with size s and fraction f | no | s bits |
| flp[e]f[f] | 0flp\<e.f\> | FloPoCo floating point number with exponent e and mantissa f | no | e+f+3 bits |