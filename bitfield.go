// Bitfield handlers

package main

import (
	"encoding/json"
//	"errors"
	"log"
	"fmt"
	"os"
//	"reflect"
//	"strings"
)

type Flagtype16 struct {
	Length	uint16			// Bit Length of the flag within the 16 bits
	Msb	uint16			// Most significant bit if the flag within the 16 bits
	Option	map[string]uint16	// What are the Options for the Flag
}
var Flag16 map[string]*Flagtype16

func copy16(flag string, x *Flagtype16) {
	// fmt.Printf("%s Length=%d Msb=%d Options=%s\n", flag, x.Length, x.Msb, x.Option)
	Flag16[flag] = x
	Flag16[flag].Length = x.Length
	Flag16[flag].Msb = x.Msb

	// Flag16[flag].Option = x.Option
	for k, _ := range x.Option {
 		Flag16[flag].Option[k] = x.Option[k]
 	}
}

func parse16(key string, val interface{}) (error) {
	flags := val.(map[string]interface{})
	// fmt.Println("We are parsing ", key)
	for okey, oval := range flags {
		// fmt.Println("We now have field ", okey)
		var tmp *Flagtype16
		tmp = new(Flagtype16)
		info := oval.(map[string]interface{})
		for ikey, ival := range info {
			// fmt.Println("We now have ikey ", ikey)
			switch ikey {
			case "length":
				tmp.Length = uint16(ival.(float64))
				// fmt.Println("Length:", tmp.Length)
			case "msb":
				tmp.Msb = uint16(ival.(float64))
				// fmt.Println("Msb:", tmp.Msb)
			case "options":
				opinfo := ival.(map[string]interface{})
				tmp.Option = make(map[string]uint16)
				for opkey, opval := range opinfo {
					tmp.Option[opkey] = uint16(opval.(float64))
					// fmt.Println("Option:",opkey, "=", tmp.Option[opkey])
				}
			default:
				log.Fatal("Invalid Field:", okey)
			}
		}
		// Copy to the global Flag16 map
		s := key + ":" + okey
		copy16(s, tmp)
	}
	return nil
}

// Given a Current curval bitfield set new bits according to
// the field and option for that bitfield you want set
// return the new bitfield
func Set16(curval uint16, field string, option string) (uint16) {
	shiftbits := uint16(16 - Flag16[field].Length - Flag16[field].Msb)
	maskbits := uint16((1 << Flag16[field].Length) - 1)
	setbits := uint16(maskbits << shiftbits)
	result := uint16(((curval) & (^setbits)))
	result |= (Flag16[field].Option[option] << shiftbits)
	return result
}

// ============================================================================================

type Flagtype8 struct {
	Length	uint8			// Bit Length of the flag within the 8 bits
	Msb	uint8			// Most significant bit if the flag within the 8 bits
	Option	map[string]uint8	// What are the Options for the Flag
}
var Flag8 map[string]*Flagtype8

func copy8(flag string, x *Flagtype8) {
	// fmt.Printf("%s Length=%d Msb=%d Options=%s\n", flag, x.Length, x.Msb, x.Option)
	Flag8[flag] = x
	Flag8[flag].Length = x.Length
	Flag8[flag].Msb = x.Msb

	// Flag8[flag].Option = x.Option
	for k, _ := range x.Option {
 		Flag8[flag].Option[k] = x.Option[k]
 	}
}

func parse8(key string, val interface{}) (error) {
	flags := val.(map[string]interface{})
	// fmt.Println("We are parsing ", key)
	for okey, oval := range flags {
		// fmt.Println("We now have field ", okey)
		var tmp *Flagtype8
		tmp = new(Flagtype8)
		info := oval.(map[string]interface{})
		for ikey, ival := range info {
			// fmt.Println("We now have ikey ", ikey)
			switch ikey {
			case "length":
				tmp.Length = uint8(ival.(float64))
				// fmt.Println("Length:", tmp.Length)
			case "msb":
				tmp.Msb = uint8(ival.(float64))
				// fmt.Println("Msb:", tmp.Msb)
			case "options":
				opinfo := ival.(map[string]interface{})
				tmp.Option = make(map[string]uint8)
				for opkey, opval := range opinfo {
					tmp.Option[opkey] = uint8(opval.(float64))
					// fmt.Println("Option:",opkey, "=", tmp.Option[opkey])
				}
			default:
				log.Fatal("Invalid Field:", okey)
			}
		}
		// Copy to the global Flag8 map
		s := key + ":" + okey
		copy8(s, tmp)
	}
	return nil
}

// Given a Current curval bitfield set new bits according to
// the field and option for that bitfield you want set
// return the new bitfield
func Set8(curval uint8, field string, option string) (uint8) {
	shiftbits := uint8(8 - Flag8[field].Length - Flag8[field].Msb)
	maskbits := uint8((1 << Flag8[field].Length) - 1)
	setbits := uint8(maskbits << shiftbits)
	result := uint8(((curval) & (^setbits)))
	result |= (Flag8[field].Option[option] << shiftbits)
	return result
}

//
// msb = most significant bit in the bitfield where
// a msb of 0 is the left most bit
// length = the number of bits in the field
// options = the name used to describe that bitfield and its value
// Reads in json describing all the Bitfield Flags
func ReadBitField(fname string) error {
	var err error
	var confdata []byte
	
	// BitField := make(map[string][]string)
	if confdata, err = os.ReadFile(fname); err != nil {
		fmt.Println("Cant open json bitfiled config file ", fname, ":", err)
		return err
	}

	var bitconfdata map[string]interface{}
	if err = json.Unmarshal([]byte(confdata), &bitconfdata); err != nil {
		fmt.Println("Cant unmarshal json from config file ", fname, ":", err)
		return err
	}

	for key, val := range bitconfdata {
		switch key {
		case "bitfield8":
			// fmt.Println("We have json for ", key)
			err := parse8(key, val)
			if (err != nil) {
				log.Fatal(err)
			}
		case "bitfield16":
			err := parse16(key, val)
			if (err != nil) {
				log.Fatal(err)
			}
		}
	}
	return nil
}

func main() {

	// Create a Global 8 Bit Flag
	Flag8 = make(map[string]*Flagtype8)
	// Create a Glogal 16 Bit Flag
	Flag16 = make(map[string]*Flagtype16)

	// Read in the json describing all the flags
	err := ReadBitField("bitfield.json")
	if (err != nil) {
		log.Fatal(err)
	}
	
	newval8 := Set8(0, "bitfield8:bits876", "v8")
	fmt.Printf("bitfield8:bits876:v8=%08b\n", newval8);
	newval8 = Set8(newval8, "bitfield8:bit3", "a2")
	fmt.Printf("bitfield8:bit3:a2=%08b\n", newval8);
	
	newval16 := Set16(0, "bitfield16:bit16", "yes")
	fmt.Printf("bitfield16:bit16:yes=%016b\n", newval16);
	
	newval16 = Set16(newval16, "bitfield16:bit16", "no")
	fmt.Printf("bitfield16:bit16:no=%016b\n", newval16);

	newval16 = Set16(newval16, "bitfield16:bit1", "true")
	fmt.Printf("bitfield16:bit1:true=%016b\n", newval16);

	newval16 = Set16(newval16, "bitfield16:bit16", "yes")
	fmt.Printf("bitfield16:bit16:yes=%016b\n", newval16);
}

