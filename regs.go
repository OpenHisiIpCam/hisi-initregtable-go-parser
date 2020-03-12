package main

import (
	_ "fmt"
)

type value struct {
	value uint32
	name  string
	desc  string
}

type field struct {
	bitStart uint32
	bitEnd   uint32
	size     uint32
	name     string
	desc     string
	values   []value
}

func (f *field) getName() string {
	if f == nil {
		return ""
	}
	return f.name
}

func (f *field) getDesc() string {
	if f == nil {
		return ""
	}
	return f.desc
}

func (f *field) getValueName(value uint32) (string, int) {

	if f == nil {
		return "", -1
	}
	if len(f.values) == 0 {
		return "", -2
	}
	for i := 0; i < len(f.values); i++ {
		if f.values[i].value == value {
			return f.values[i].name, 0
		}
	}
	return "unknown", -3
}

/*
func (f *field) getBits() (bitStart, bitEnd uint32) {
	if f == nil {
		return 0, 0
	}
	return f.bitStart, f.bitEnd
}
*/

type register32 struct {
	//base   uint32 //base address
	//offset uint32 //offset address
	addr   uint32 //full address
	reset  uint32 //reset value
	name   string
	desc   string
	fields []field
}

func (r *register32) getName() string {
	if r == nil {
		return ""
	}
	return r.name
}

func (r *register32) getDesc() string {
	if r == nil {
		return ""
	}
	return r.desc
}

/*
func (r *register32) print() {
	if r != nil {
		fmt.Printf("addr: 0x%08X ", r.addr)
		fmt.Println(*r)
	} else {
		panic("reg pointer is nil!")
	}
}
*/

func (r *register32) findField(bitStart, size uint32) *field {
	var f *field

	if r == nil {
		return nil
	}

	for i := 0; i < len(r.fields); i++ {
		if r.fields[i].bitStart == bitStart {
			f = &r.fields[i]
		}
	}
	if f != nil {
		if f.size != size {
			f = nil
		}
	}
	return f
}

func (r *register32) findFields(bitStart, bitEnd uint32) []*field {
	fields := make([]*field, 0)

	if r == nil {
		return fields
	}

	for i := 0; i < len(r.fields); i++ {
		field := &r.fields[i]
		if field.bitStart >= bitStart {
			if field.bitEnd <= bitEnd {
				fields = append(fields, field)
			}
		}
	}
	return fields
}

/*
var test = [...]register32{
	register32{
		//base:   0,
		//offset: 0,
		addr:  0,
		reset: 0,
		name:  "test",
		desc:  "test",
		fields: []field{
			field{
				bitStart: 0,
				bitEnd:   1,
				size:     2,
				name:     "ftest",
				desc:     "ftest",
				values: []value{
					value{
						value:  0,
						name: "vtest",
					},
				},
			},
		},
	},
}
*/
var regsByAddr map[uint32]*register32

//var regsByName map[string]*register32

func indexByAddr() {
	regsByAddr = make(map[uint32]*register32)
	/*some bug with &r
	for i, r := range registers {
		//fmt.Println("Adding ", reg.addr)
		fmt.Printf("adding 0x%08X %p %p\n", r.addr, &r, &registers[i])
		regsByAddr[r.addr] = &r
	}
	*/
	for i := 0; i < len(registers); i++ {
		//fmt.Printf("adding 0x%08X %p\n", registers[i].addr, &registers[i])
		regsByAddr[registers[i].addr] = &registers[i]
	}
	/*
		fmt.Println(regsByAddr)
		fmt.Println(regsByAddr[0x12010034])
		fmt.Println(regsByAddr[302055476])
	*/
}

func findByAddr(r uint32) *register32 {
	var found *register32

	if len(regsByAddr) > 0 {
		//fmt.Println("looking for ", r)
		found = regsByAddr[r]
	} else {
		//TODO interation search
		for _, reg := range registers {
			if reg.addr == r {
				found = &reg
				break
				//fmt.Println(i, reg)
			}
		}
		//if found != nil {
		//	fmt.Println(*found)
		//}
	}
	//fmt.Println("FOUND: ", found)
	return found
}
