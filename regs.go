package main

import (
	_ "fmt"
)

var addrMaps map[string][]register32

func addAddrMap(name string, regs []register32) {
	if len(addrMaps) == 0 {
		addrMaps = make(map[string][]register32)
	}
	addrMaps[name] = regs
}

type value struct {
	value uint32
	name  string
	desc  string
}

type field struct {
	bitStart uint8
	bitEnd   uint8
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

func (f *field) getValueName(value uint32) string {

	if f == nil {
		return ""
	}
	if len(f.values) == 0 {
		return ""
	}
	for i := 0; i < len(f.values); i++ {
		if f.values[i].value == value {
			return f.values[i].name
		}
	}
	return "unknown"
}

type register32 struct {
	addr   uint32 //full address
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

func (r *register32) findFields(bitStart, bitEnd uint8) []*field {
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

func findByAddr(registers []register32, r uint32) *register32 {
	var found *register32

	for _, reg := range registers {
		if reg.addr == r {
			found = &reg
			break
		}
	}
	return found
}
