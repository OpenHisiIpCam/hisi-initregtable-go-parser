#!/usr/bin/env python3

import sys

from systemrdl import RDLCompiler, RDLListener, RDLWalker, RDLCompileError
from systemrdl.node import FieldNode

import argparse

consts = []
r32 = ""

class MyGoListener(RDLListener):
    #def __init__(self):
        #global regs
        #global r32

    def enter_Addrmap(self, node):
        global r32
        r32 = r32 + str("var registers = [...]register32 {\n") #print("var registers = [...]register32 {")

    def exit_Addrmap(self, node):
        global r32
        r32 = r32 + str("}\n") #print("}")

    def enter_Reg(self, node):
        global r32
        global consts
        consts.append([node.type_name, node.absolute_address])
        #regs.append(node.type_name)
        #regs.append(123)

        r32 = r32 + str("register32 {\n") #print("\tregister32 {")
        r32 = r32 + str("addr: " + hex(node.absolute_address) + ",\n") #print("\t\taddr: " + hex(node.absolute_address) + ",")
        #r32 = r32 + str("reset: 0,\n") #print("\t\treset: 0,")
        r32 = r32 + str('name: "' + node.type_name + '",\n') #print('\t\tname: "' + node.type_name + '",')
        r32 = r32 + str('desc: "' + node.get_property("name") + '",\n') #print('\t\tdesc: "' + node.get_property("name") + '",')
        r32 = r32 + str("fields: []field {\n") #print("\t\tfields: []field {")

    def exit_Reg(self, node):
        global r32
        r32 = r32 + str("},\n") #print("\t\t},")
        r32 = r32 + str("},\n") #print("\t},")

    def enter_Field(self, node):
        global r32
        r32 = r32 + str("field {\n") #print("\t\t\tfield {")
        r32 = r32 + str("bitStart: " + str(node.low) + ",\n") #print("\t\t\t\tbitStart: " + str(node.low) + ","),
        r32 = r32 + str("bitEnd: " + str(node.high) + ",\n") #print("\t\t\t\tbitEnd: " + str(node.high) + ","),
        #r32 = r32 + str("size: 2,\n") #print("\t\t\t\tsize: 2,"),
        r32 = r32 + str('name: "' + node.type_name  + '",\n') #print('\t\t\t\tname: "' + node.type_name  + '",')
        r32 = r32 + str('desc: "' + node.get_property("name") + '",\n') #print('\t\t\t\tdesc: "' + node.get_property("desc") + '",')
        if (node.get_property("encode")):
            r32 = r32 + str("values: []value {\n") #print("\t\t\t\tvalues: []value {")
            for item in node.get_property("encode"):
                r32 = r32 + str("value {\n") #print("\t\t\t\t\tvalue {")
                name_tmp = str(item).split(".")
                name = name_tmp[1]
                r32 = r32 + str('name: "' + name + '",\n') #print('\t\t\t\t\t\tname: "' + name + '",')
                r32 = r32 + str("value: " + str(item.value) + ",\n") #print("\t\t\t\t\t\tvalue: " + str(item.value) + ",")
                #r32 = r32 + str("desc: " + )
                r32 = r32 + str("},\n") #print("\t\t\t\t\t},")
            r32 = r32 + str("},\n") #print("\t\t\t\t},")

    def exit_Field(self, node):
        global r32
        r32 = r32 + str("},\n") #print("\t\t\t},")

# Collect input files from the command line arguments
input_files = sys.argv[1:]

# Create an instance of the compiler
rdlc = RDLCompiler()

try:
    # Compile all the files provided
    for input_file in input_files:
        rdlc.compile_file(input_file)

    # Elaborate the design
    root = rdlc.elaborate()
except RDLCompileError:
    # A compilation error occurred. Exit with error code
    sys.exit(1)


# Traverse the register model!
walker = RDLWalker(unroll=True)
#listener = MyModelPrintingListener()
listener = MyGoListener()
walker.walk(root, listener)


#print("//+build mytag")
print("package main")
print("")

#print(consts)
if len(consts)>0:
    print("const (")
    for rname, raddr in consts:
        print("\t" + rname + "\t=\t" + hex(raddr))
    print(")")

print(r32)
