#!/usr/bin/env python3

import sys

from systemrdl import RDLCompiler, RDLListener, RDLWalker, RDLCompileError
from systemrdl.node import FieldNode

import argparse

################################################################################
class RdlToGoCodeListener(RDLListener):
    r32 = ""
    name = ""

    def enter_Addrmap(self, node):
        self.name = node.type_name
        self.r32 += "var " + self.name + "Registers = [...]register32 {\n"

    def exit_Addrmap(self, node):
        self.r32 += "}\n"

    def enter_Reg(self, node):
        self.r32 += "register32 {\n"
        self.r32 += "addr: " + hex(node.absolute_address) + ",\n"
        self.r32 += 'name: "' + node.type_name + '",\n'
        self.r32 += 'desc: "' + node.get_property("name") + '",\n'
        self.r32 += "fields: []field {\n"

    def exit_Reg(self, node):
        self.r32 += "},\n"
        self.r32 += "},\n"

    def enter_Field(self, node):
        self.r32 += "field {\n"
        self.r32 += "bitStart: " + str(node.low) + ",\n"
        self.r32 += "bitEnd: " + str(node.high) + ",\n" 
        self.r32 += 'name: "' + node.type_name  + '",\n'
        self.r32 += 'desc: "' + node.get_property("name") + '",\n'
        if (node.get_property("encode")):
            self.r32 += "values: []value {\n"
            for item in node.get_property("encode"):
                self.r32 += "value {\n"
                name = str(item).split(".")
                self.r32 += 'name: "' + name[1] + '",\n'
                self.r32 += "value: " + str(item.value) + ",\n"
                self.r32 += "},\n"
            self.r32 += "},\n"

    def exit_Field(self, node):
        self.r32 += "},\n"

################################################################################
parser = argparse.ArgumentParser(description='dfdfdfd.')

#parser.add_argument('infile', nargs='?', type=argparse.FileType('r') )
parser.add_argument('rdl',     type=str,       help='SystemRDL source file')
parser.add_argument('chip',    type=str,       help='Chip model name')

#parser.add_argument('--', type=int,   help='Full (default) or Half duplex mode',               required=False)
#parser.add_argument('--mc21', type=int,   help='To remove',               required=False)
#parser.add_argument('--servercamera', type=int, help='Camera number',               required=False)


args = parser.parse_args()

#print(args)

rdlc = RDLCompiler()

try:
    rdlc.compile_file(args.rdl)

    root = rdlc.elaborate(args.chip)
except RDLCompileError:
    sys.exit(1)

#sys.exit(10)

walker = RDLWalker(unroll=True)
golistener = RdlToGoCodeListener()
walker.walk(root, golistener)

print("//THIS FILE WAS AUTO GENERATED. PLEASE DO NOT EDIT")
print("package main")

print('func init() { addAddrMap("%s", %sRegisters[:]) }' % (golistener.name, golistener.name))

print(golistener.r32)
