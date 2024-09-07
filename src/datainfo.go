package main

import (
	"bufio"
	"fmt"
)

type AccessType int

const (
	Public = iota
	Protected
	Private
)

type PropertyInfo struct {
	Macro       string
	Declaration string
	Comments    []string
	Access      AccessType
}

type FunctionInfo struct {
	Name        string
	Macro       string
	Declaration string
	Comments    []string
	Access      AccessType
}

type DataInfo struct {
	Name       string
	Parents    []string
	Comments   []string
	Properties []PropertyInfo
	Functions  []FunctionInfo
	IsStruct   bool
	IsEnum     bool
}

func (d *DataInfo) OutputHeader(writer *bufio.Writer) {
	writer.WriteString("\n")
	writer.WriteString(fmt.Sprintf("\n## `" + d.Name + "` \n\n"))
}

func (d *DataInfo) OutputParents(writer *bufio.Writer) {
	if len(d.Parents) > 0 {
		writer.WriteString("\n")
		writer.WriteString("__Parent Classes:__\n")
		writer.WriteString("[ ")
		for i, parent := range d.Parents {
			isLast := i == len(d.Parents)-1
			if !isLast {
				writer.WriteString(fmt.Sprintf("`%s`, ", parent))
			} else {
				writer.WriteString(fmt.Sprintf("`%s`", parent))
			}
		}
		writer.WriteString(" ]\n")
	}
}

func (d *DataInfo) OutputDescription(writer *bufio.Writer) {
	if len(d.Comments) > 0 {
		writer.WriteString("\n")
		for i, com := range d.Comments {
			isLast := i == max(len(d.Comments)-1, 0)
			if isLast {
				writer.WriteString("" + cleanComment(com) + " \n")
			} else {
				writer.WriteString("" + cleanComment(com) + " \\\n")
			}
		}
	}
}

func (d *DataInfo) OutputProperties(writer *bufio.Writer) {
	if d.HasDocumentedProperties() {
		writer.WriteString("\n")
		writer.WriteString("### Properties\n\n")

		writer.WriteString("```cpp\n")
		for _, prop := range d.Properties {
			if len(prop.Comments) == 0 {
				continue
			}
			for _, comm := range prop.Comments {
				writer.WriteString("// " + cleanComment(comm) + " \n")
			}
			if prop.Macro != "" {
				writer.WriteString(prop.Macro + "\n")
			}
			writer.WriteString(prop.Declaration + "\n\n")
		}
		writer.WriteString("```\n")
	}
}

func (d *DataInfo) OutputFunctions(writer *bufio.Writer) {
	if d.HasDocumentedFunctions() {
		writer.WriteString("\n")
		writer.WriteString("### Functions\n\n")

		for _, function := range d.Functions {
			if len(function.Comments) == 0 {
				continue
			}
			writer.WriteString("#### `" + function.Name + "`\n")
			for i, comm := range function.Comments {
				isLast := i == max(len(function.Comments)-1, 0)
				if isLast {
					writer.WriteString("> " + cleanComment(comm) + " \n")
				} else {
					writer.WriteString("> " + cleanComment(comm) + " \\\n")
				}
			}
			writer.WriteString("```cpp\n")
			writer.WriteString(function.Declaration + "\n")
			writer.WriteString("```\n")
		}
	}
}

func (d *DataInfo) HasDocumentation() bool {
	return len(d.Comments) > 0 && d.HasDocumentedProperties() || d.HasDocumentedFunctions()
}

func (d *DataInfo) HasDocumentedProperties() bool {
	for _, prop := range d.Properties {
		if len(prop.Comments) > 0 {
			return true
		}
	}
	return false
}

func (d *DataInfo) HasDocumentedFunctions() bool {
	for _, function := range d.Functions {
		if len(function.Comments) > 0 {
			return true
		}
	}
	return false
}
