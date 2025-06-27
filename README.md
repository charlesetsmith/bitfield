Go module to set bits as described in the JSON file. You can describe 8, 16, 32 and 64 bit fields and any subfields within those. 
Still work to do getting this all into a module, but the basics are there.

For each type (uint8, uint16, uint32, uint64) a subfield is defined with a name, Length (number of bits), Msb (Most significant bit within the 
uint and the options that field can have.

You use numb := set8(0, "fieldname", "option") to set the option within numb.

