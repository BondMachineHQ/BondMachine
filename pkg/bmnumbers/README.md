# bmnumbers

bmstack is part of BondMachine project. bmnumbers is both a command line tool to convert or cast numbers to and from different formats and a library to do the same. It is used within the BondMachine every time numbers are handled.

## Supported number types

The supported number types are listed in the following table.

| Type Name | Prefix | Description | Static |
| ---- | ------- | ----------- | ------ |
| unsigned | | Unsigned integer | yes |
| signed | | Signed integer | yes |
| bin | | Binary number |yes |
| hex | 0x | Hexadecimal number | yes |
| float16 | 0f<16> | IEEE 754 half precision floating point number | yes |
| float32 | 0f | IEEE 754 single precision floating point number | yes |
| lqs[s]t[t] | 0lq<s.t> | Linear quantized number with size s and type t | no |
| fps[s]f[f] | 0fp<s.f> | Fixed point number with size s and fraction f | no |
