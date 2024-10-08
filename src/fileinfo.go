package main

import (
	"bufio"
	"fmt"
	"strings"
)

type LineId int

const (
	Empty = iota
	Comment
	Class
	Struct
	Enum
	EnumProp
	Function
	Property
	AccessModifier
	OpenIgnore
	CloseIgnore
	OpenBracket
	CloseBracket
)

type FileInfo struct {
	Path string
	Name string
	Data []DataInfo
}

func (f *FileInfo) OutputInfo(writer *bufio.Writer) (enums, structs, classes []DataInfo) {

	writer.WriteString("\n__FileName:__ `" + f.Name + "`\n")

	for _, data := range f.Data {
		if data.IsEnum {
			enums = append(enums, data)
		} else if data.IsStruct {
			structs = append(structs, data)
		} else {
			classes = append(classes, data)
		}
	}

	if len(enums) > 0 {
		writer.WriteString("- __Enum List:__ \n")

		writer.WriteString("[ ")
		for i, e := range enums {
			isLast := i == len(enums)-1
			if !isLast {
				writer.WriteString(fmt.Sprintf("[`" + e.Name + "`](#" + strings.ToLower(e.Name) + ") | "))
			} else {
				writer.WriteString(fmt.Sprintf("[`" + e.Name + "`](#" + strings.ToLower(e.Name) + ")"))
			}
		}
		writer.WriteString(" ]\n")
	}

	if len(structs) > 0 {
		writer.WriteString("- __Struct List:__ \n")

		writer.WriteString("[ ")
		for i, s := range structs {
			if s.HasDocumentation() {
				isLast := i == len(structs)-1
				if !isLast {
					writer.WriteString(fmt.Sprintf("[`" + s.Name + "`](#" + strings.ToLower(s.Name) + ") | "))
				} else {
					writer.WriteString(fmt.Sprintf("[`" + s.Name + "`](#" + strings.ToLower(s.Name) + ")"))
				}
			}
		}
		writer.WriteString(" ]\n")
	}

	if len(classes) > 0 {
		writer.WriteString("- __Class List:__ \n")

		writer.WriteString("[ ")
		for i, c := range classes {
			if c.HasDocumentation() {
				isLast := i == len(classes)-1
				if !isLast {
					writer.WriteString(fmt.Sprintf("[`" + c.Name + "`](#" + strings.ToLower(c.Name) + ") | "))
				} else {
					writer.WriteString(fmt.Sprintf("[`" + c.Name + "`](#" + strings.ToLower(c.Name) + ")"))
				}
			}
		}
		writer.WriteString(" ]\n")
	}

	return
}
