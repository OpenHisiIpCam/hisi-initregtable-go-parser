package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fatih/color"
	_ "math"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

/*
func reverse(numbers []byte) {
	for i, j := 0, len(numbers)-1; i < j; i, j = i+1, j-1 {
		numbers[i], numbers[j] = numbers[j], numbers[i]
	}
}
*/
func bytesToUInt32(ba []byte) uint32 {
	var value uint32
	value |= uint32(ba[0])
	value |= uint32(ba[1]) << 8
	value |= uint32(ba[2]) << 16
	value |= uint32(ba[3]) << 24
	return value
}

func main() {
	var rbFile = flag.String("file", "u-boot.bin", "file with reg data")
	var rbOffset = flag.Int("offset", 0, "offset of reg data in file")
	var rbSize = flag.Int("size", 4016, "size of reg data in file")
	var rbChip = flag.String("chip", "hi3516av200", "HiSilicon chip family name")
	flag.Parse()

	if (*rbSize % 16) != 0 {
		panic("size shoud be muptiple of 16 bytes")
	}

	//fmt.Println("test")
	//color.Cyan("Prints text in cyan.")

	//blue := color.New(color.FgCyan, color.Bold)

	///////////////////////////////////////////////////////////////////
	//indexByAddr()

	///////////////////////////////////////////////////////////////////

	file, err := os.Open(*rbFile)
	check(err)

	_, err = file.Seek(int64(*rbOffset), 0)
	check(err)

	regBin := make([]byte, int64(*rbSize))
	read, err := file.Read(regBin)
	check(err)

	fmt.Printf("Read %d bytes from expected %d\n", read, *rbSize)
	//fmt.Printf("%v\n", string(b2[:n2]))

	fmt.Printf("reg        value      delay      attrs\n")
	//0x12010034 0x12010034 0x12010034 0x12010034
	//var i int
	for i := 0; i < *rbSize; i = i + 16 {
		//fmt.Print("offset: ", i, " ")
		reg := bytesToUInt32(regBin[i : i+4])
		value := bytesToUInt32(regBin[i+4 : i+8])
		delay := bytesToUInt32(regBin[i+8 : i+12])
		attrs := bytesToUInt32(regBin[i+12 : i+16])

		if reg == 0x0 {
			//finish
			fmt.Println("NULL reg found, aborting.")
			//panic("aborting")
			break
		}

		//fmt.Printf("%d:", i/16)
		fmt.Printf("0x%08X", reg)
		fmt.Printf(" 0x%08X", value)
		fmt.Printf(" 0x%08X", delay)
		fmt.Printf(" 0x%08X", attrs)

		regi := findByAddr(reg)
		//if regi != nil {
		//	regi.print()
		//}

		write := attrs & 0b111
		writeNum := ((attrs & 0b11111000) >> 3) //+ 1
		writeStart := (attrs & 0b1111100000000000) >> 11
		read := (attrs & 0b1110000000000000000) >> 16
		readNum := ((attrs & 0b111110000000000000000000) >> 19) //+ 1
		readStart := (attrs & 0b11111000000000000000000000000000) >> 27
		//fmt.Println("write ", write, " writeNum ", writeNum, " writeStart ", writeStart)
		//fmt.Println("read ", read, " readNum ", readNum, " readStart ", readStart)

		if write&0b100 == 0b100 {
			//this is write
			color.Set(color.FgBlue)
			fmt.Print("\tWRITE")
			regName := regi.getName()
			if regName != "" {
				fmt.Printf("\t%s[%d:%d] (%s)", regName, writeStart, writeStart+writeNum, regi.getDesc())
			} else {
				color.Set(color.FgRed)
				fmt.Printf("\tunknown[%d:%d]", writeStart, writeStart+writeNum)
				color.Set(color.FgBlue)
			}
			//fmt.Printf("[%d:%d]", writeStart, writeStart+writeNum)

			if delay > 0 {
				fmt.Print(" AND DELAY", delay)
			}

			fmt.Println()
			color.Set(color.FgCyan)
			fields := regi.findFields(writeStart, writeStart+writeNum)
			if len(fields) > 0 {
				tvalue := value << writeStart
				for _, field := range fields {
					fmt.Printf("\t\t\t\t\t\t\tfield %s[%d:%d] (%s) ", field.getName(), field.bitStart, field.bitEnd, field.getDesc())
					fvalue := (tvalue << (31 - field.bitEnd) >> (31 - field.bitEnd)) >> field.bitStart
					recognized_value, err := field.getValueName(fvalue)
					if err == 0 {
						fmt.Printf("val = 0x%X (%s)\n", fvalue, recognized_value)
					} else {
						fmt.Printf("val = 0x%X\n", fvalue)
					}
					//fmt.Printf("debug: (value << (31 - %d) >> (31 - %d)) >> %d\n", field.bitEnd, field.bitEnd, field.bitStart)
				}
			} else {
				color.Set(color.FgRed)
				fmt.Printf("\t\t\t\t\t\t\tfields not found\n")
			}

			//fmt.Println()
			color.Unset()
		} else if read&0b100 == 0b100 {
			//this is read
			color.Set(color.FgGreen)
			fmt.Print("\tREAD")
			regName := regi.getName()
			if regName != "" {
				fmt.Printf("\t%s[%d:%d] (%s)", regName, readStart, readStart+readNum, regi.getDesc())
			} else {
				color.Set(color.FgRed)
				fmt.Printf("\tunknown[%d:%d]", readStart, readStart+readNum)
				color.Set(color.FgGreen)
			}
			//fmt.Printf("[%d:%d]", readStart, readStart+readNum)
			if delay > 0 {
				fmt.Print("\tAND DELAY", delay)
			}

			fmt.Println()
			//fmt.Printf("\t\t\t\t\t\t\t")
			color.Set(color.FgCyan)
			fields := regi.findFields(readStart, readStart+readNum)
			if len(fields) > 0 {
				tvalue := value << writeStart
				for _, field := range fields {
					fmt.Printf("\t\t\t\t\t\t\tfield %s[%d:%d] (%s) ", field.getName(), field.bitStart, field.bitEnd, field.getDesc())
					fvalue := (tvalue << (31 - field.bitEnd) >> (31 - field.bitEnd)) >> field.bitStart
					recognized_value, err := field.getValueName(fvalue)
					if err == 0 {
						fmt.Printf("val = 0x%X (%s)\n", fvalue, recognized_value)
					} else {
						fmt.Printf("val = 0x%X\n", fvalue)
					}

				}
			} else {
				color.Set(color.FgRed)
				fmt.Printf("\t\t\t\t\t\t\tfields not found\n")
			}

			//fmt.Println()
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
