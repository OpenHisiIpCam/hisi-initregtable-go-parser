package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fatih/color"
)

func bytesToUInt32(ba []byte) uint32 {
	var value uint32
	value |= uint32(ba[0])
	value |= uint32(ba[1]) << 8
	value |= uint32(ba[2]) << 16
	value |= uint32(ba[3]) << 24
	return value
}

func parseRegFields(r *register32, value uint32, bitStart, bitNum uint8) {
	regName := r.getName()

	if regName != "" {
		fmt.Printf("\t%s[%d:%d] (%s)", regName, bitStart, bitStart+bitNum, r.getDesc())
	} else {
		color.Set(color.FgRed)
		fmt.Printf("\tunknown[%d:%d]", bitStart, bitStart+bitNum)
	}

	fmt.Println()
	color.Set(color.FgCyan)

	fields := r.findFields(bitStart, bitStart+bitNum)
	if len(fields) > 0 {
		for _, field := range fields {
			fmt.Printf("\t\t\t\t\t\t\tfield %s[%d:%d] (%s) ", field.getName(), field.bitStart, field.bitEnd, field.getDesc())

			fieldValue := ((value << bitStart) << (31 - field.bitEnd) >> (31 - field.bitEnd)) >> field.bitStart
			recognizedValue := field.getValueName(fieldValue)
			if recognizedValue != "" {
				fmt.Printf("val = 0x%X (%s)\n", fieldValue, recognizedValue)
			} else {
				fmt.Printf("val = 0x%X\n", fieldValue)
			}
		}
	} else {
		fmt.Printf("\t\t\t\t\t\t\tfields not found\n")
	}

}

func main() {
	var rbFile = flag.String("file", "u-boot.bin", "file with reg data")
	var rbOffset = flag.Int("offset", 0, "offset of reg data in file")
	var rbSize = flag.Int("size", 4016, "size of reg data in file")
	var rbChip = flag.String("chip", "hiXXX", "HiSilicon chip family name")
	flag.Parse()

	color.Set(color.FgRed)

	if (*rbSize % 16) != 0 {
		fmt.Println("Size shoud be muptiple of 16 bytes.")
		os.Exit(1)
	}

	if addrMaps[*rbChip] == nil {
		fmt.Print("WARNING! No such chip ", *rbChip, " in the internal database (avalible: ")
		for key, _ := range addrMaps {
			fmt.Print(key, " ")
		}
		//fmt.Println(")")
		fmt.Println("). All registers will be unknown!")
	}

	file, err := os.Open(*rbFile)
	if err != nil {
		fmt.Println("File open error. ", err)
		os.Exit(1)
	}

	_, err = file.Seek(int64(*rbOffset), 0)

	regBin := make([]byte, int64(*rbSize))
	read, err := file.Read(regBin)
	if err != nil {
		fmt.Println("File read error. ", err)
		os.Exit(1)

	}

	if read < *rbSize {
		fmt.Printf("WARNING! Read %d bytes from expected %d\n", read, *rbSize)
	}
	color.Unset()

	fmt.Printf("reg        value      delay      attrs\n")

	var i int
	for i = 0; i < read; i = i + 16 {

		reg := bytesToUInt32(regBin[i : i+4])
		value := bytesToUInt32(regBin[i+4 : i+8])
		delay := bytesToUInt32(regBin[i+8 : i+12])
		attrs := bytesToUInt32(regBin[i+12 : i+16])

		if reg == 0x0 {
			//finish
			fmt.Println("NULL reg found, aborting.")
			break
		}

		fmt.Printf("0x%08X", reg)
		fmt.Printf(" 0x%08X", value)
		fmt.Printf(" 0x%08X", delay)
		fmt.Printf(" 0x%08X", attrs)

		regi := findByAddr(addrMaps[*rbChip], reg)

		write := attrs & 0b111
		read := (attrs & 0b1110000000000000000) >> 16

		if write&0b100 == 0b100 {
			color.Set(color.FgBlue)
			fmt.Print("\tWRITE")

			writeNum := uint8((attrs & 0b11111000) >> 3)
			writeStart := uint8(attrs&0b1111100000000000) >> 11

			parseRegFields(regi, value, writeStart, writeNum)
			color.Unset()
		} else if read&0b100 == 0b100 {
			color.Set(color.FgGreen)
			fmt.Print("\tREAD")

			readNum := uint8((attrs & 0b111110000000000000000000) >> 19)
			readStart := uint8(attrs&0b11111000000000000000000000000000) >> 27

			parseRegFields(regi, value, readStart, readNum)
			color.Unset()
		} else {
			if delay > 0 {
				color.Set(color.FgYellow)
				fmt.Println("\tDELAY  ", delay)
				color.Unset()
			} else {
				color.Set(color.FgRed)
				fmt.Println("\tUNKNOWN")
				color.Unset()
			}
		}
	}
	fmt.Println("Total ", (i/16)-1, "lines processed.")

}
