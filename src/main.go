package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type LineId int

const (
	Empty = iota
	Comment
	Class
	Struct
	Enum
	Function
	Property
	AccessModifier
	OpenIgnore
	CloseIgnore
	OpenBracket
	CloseBracket
)

type AccessType int

const (
	Public = iota
	Protected
	Private
)

// type DataType int
// const (
// 	Class = iota
// 	Struct
// 	Enum
// )

type DataInfo struct {
	Name       string
	Parents    []string
	Comments   []string
	Properties []PropertyInfo
	Functions  []FunctionInfo
	IsStruct   bool
	IsEnum     bool
}

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

type FileInfo struct {
	Path string
	Name string
	Data []DataInfo
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: program <source_folder> <destination_folder>")
		os.Exit(1)
	}

	sourceFolder := os.Args[1]
	destFolder := os.Args[2]

	ignoreFiles := []string{
		"FlowPilotModule.h",
		"FlowPilotCustomVersion.h",
	}

	var fileInfoList []FileInfo

	err := filepath.Walk(sourceFolder, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if strings.HasSuffix(info.Name(), ".h") || strings.HasSuffix(info.Name(), ".hpp") {

			skipFile := false
			for _, ignore := range ignoreFiles {
				if ignore == info.Name() {
					skipFile = true
					break
				}
			}

			if !skipFile {
				fOutput := FileInfo{
					path,
					info.Name(),
					[]DataInfo{},
				}

				fileInfoList = append(fileInfoList, fOutput)
			}
		}

		return nil
	})

	for i := 0; i < len(fileInfoList); i++ {
		processFile(&fileInfoList[i], destFolder, i)
	}

	if err != nil {
		fmt.Printf("Error walking through directory: %v\n", err)
	}
}

func processFile(fileInfo *FileInfo, destFolder string, fileIndex int) {
	file, err := os.Open(fileInfo.Path)
	if err != nil {
		fmt.Printf("Error opening file %s: %v\n", fileInfo.Path, err)
		return
	}
	defer file.Close()

	extractInfo(file, fileInfo)

	if fileInfo.Name == "" {
		fmt.Printf("No class found in file %s\n", fileInfo.Path)
		return
	}

	outputMarkdown(fileInfo, destFolder)
}

func extractInfo(file *os.File, fileInfo *FileInfo) {
	scanner := bufio.NewScanner(file)

	var currentClassIndex IntStack = IntStack{}
	var commentStack []string = []string{}
	var fnMacro string = ""
	var propMacro string = ""

	ignorePrefixList := []string{
		"// UFlowPilotTask",
		"//~UFlowPilotTask",
	}

	var ignoreBlock = false
	var currentAccessType AccessType = Private
	var prevId LineId = Empty
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		var skip = false
		for i := 0; i < len(ignorePrefixList); i++ {
			if strings.HasPrefix(line, ignorePrefixList[i]) {
				skip = true
				break
			}
		}

		if skip {
			continue
		}

		id := idLine(line, prevId)

		// fmt.Printf("[%d][%d][%d] %s\n", id, currentAccessType, ignoreBlock, line)

		if ignoreBlock && id == CloseIgnore {
			ignoreBlock = false
			continue
		}

		if !ignoreBlock && id == OpenIgnore {
			ignoreBlock = true
			continue
		}

		if ignoreBlock {
			continue
		}

		switch id {
		case Empty:
			continue
		case Comment:
			// stack comments
			commentStack = append(commentStack, line)
		case CloseBracket:
			currentClassIndex.Pop()
		case OpenBracket:

		case Enum:
			if !isEnumMacro(line) {
				var name = extractEnumInfo(line)
				var info = DataInfo{
					Name:     name,
					IsStruct: false,
					IsEnum:   true,
				}
				fileInfo.Data = append(fileInfo.Data, info)
				currentClassIndex.Push(len(fileInfo.Data) - 1)

				commentStack = []string{}
			}
		case Class:
			if !isClassMacro(line) {
				var name, parents, _ = extractClassInfo(line)

				var info = DataInfo{
					Name:     name,
					Parents:  parents,
					Comments: commentStack,
					IsStruct: false,
					IsEnum:   false,
				}
				fileInfo.Data = append(fileInfo.Data, info)
				currentClassIndex.Push(len(fileInfo.Data) - 1)

				commentStack = []string{}
			}
		case Struct:
			if !isStructMacro(line) {
				var name, parents, _ = extractStructInfo(line)

				var info = DataInfo{
					Name:     name,
					Parents:  parents,
					Comments: commentStack,
					IsStruct: true,
					IsEnum:   false,
				}
				fileInfo.Data = append(fileInfo.Data, info)
				currentClassIndex.Push(len(fileInfo.Data) - 1)

				commentStack = []string{}
			}
		case Function:
			if isFunctionMacro(line) {
				fnMacro = line
			} else {
				var data = FunctionInfo{
					Name:        extractFunctionName(line),
					Macro:       fnMacro,
					Declaration: line,
					Comments:    commentStack,
					Access:      currentAccessType,
				}

				fnMacro = ""
				commentStack = []string{}

				fileInfo.Data[currentClassIndex.Top()].Functions = append(fileInfo.Data[currentClassIndex.Top()].Functions, data)
			}
		case Property:
			if isPropertyMacro(line) {
				propMacro = line
			} else {
				var data = PropertyInfo{
					Macro:       propMacro,
					Declaration: line,
					Comments:    commentStack,
					Access:      currentAccessType,
				}

				propMacro = ""
				commentStack = []string{}

				fileInfo.Data[currentClassIndex.Top()].Properties = append(fileInfo.Data[currentClassIndex.Top()].Properties, data)
			}
		case AccessModifier:
			accessStr := strings.TrimRight(line, ":")
			switch accessStr {
			case "public":
				currentAccessType = Public
			case "protected":
				currentAccessType = Protected
			case "private":
				currentAccessType = Private
			}
		default:
		}

		prevId = id
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
	}
}

func idLine(line string, prevId LineId) LineId {

	if len(line) == 0 || isCopy(line) {
		return Empty
	}

	if isOpenIgnore(line) {
		return OpenIgnore
	}

	if isCloseIgnore(line) {
		return CloseIgnore
	}

	if isCloseBracket(line) {
		return CloseBracket
	}

	if isOpenBracket(line) {
		return OpenBracket
	}

	if isComment(line, prevId == Comment) {
		return Comment
	}

	if isClassMacro(line) || isClass(line) {
		return Class
	}

	if isStructMacro(line) || isStruct(line) {
		return Struct
	}

	if isEnumMacro(line) || isEnum(line) {
		return Enum
	}

	if isPropertyMacro(line) || isProperty(line, prevId == Property) {
		return Property
	}

	if isFunctionMacro(line) || isFunction(line, prevId == Function) {
		return Function
	}

	if isAccessModifier(line) {
		return AccessModifier
	}

	return Empty
}

func isOpenIgnore(line string) bool {
	return strings.Contains(line, "#if")
}

func isCloseIgnore(line string) bool {
	return strings.Contains(line, "#endif")
}

func isComment(line string, prevIsComment bool) bool {
	if prevIsComment {
		if strings.HasPrefix(line, "*/") {
			return false
		}

		if strings.HasPrefix(line, "*") {
			return true
		}
	}

	return strings.HasPrefix(line, "/*") || strings.HasPrefix(line, "//")
}

func extractClassInfo(line string) (class string, parent []string, foundClassOpenBracket bool) {
	parts := strings.Fields(line)
	foundParentDelimiter := false
	foundClassOpenBracket = false
	for _, part := range parts {
		if strings.HasPrefix(part, "class") || strings.Contains(part, "_API") {
			continue
		}
		if strings.HasPrefix(part, "public") || strings.HasPrefix(part, "private") || strings.HasPrefix(part, "protected") {
			continue
		}
		if strings.HasPrefix(part, ":") {
			foundParentDelimiter = true
			continue
		}
		if strings.HasPrefix(part, "U") || strings.HasPrefix(part, "I") {
			if foundParentDelimiter {
				parent = append(parent, part)
			} else {
				class = part
			}
		}
		if isOpenBracket(part) {
			foundClassOpenBracket = true
		}
	}
	return
}

func extractStructInfo(line string) (class string, parent []string, foundClassOpenBracket bool) {
	parts := strings.Fields(line)
	foundParentDelimiter := false
	foundClassOpenBracket = false
	for _, part := range parts {
		if strings.HasPrefix(part, "struct") || strings.Contains(part, "_API") {
			continue
		}
		if strings.HasPrefix(part, "public") || strings.HasPrefix(part, "private") || strings.HasPrefix(part, "protected") {
			continue
		}
		if strings.HasPrefix(part, ":") {
			foundParentDelimiter = true
			continue
		}
		if strings.HasPrefix(part, "F") {
			if foundParentDelimiter {
				parent = append(parent, part)
			} else {
				class = part
			}
		}
		if isOpenBracket(part) {
			foundClassOpenBracket = true
		}
	}
	return
}

func extractEnumInfo(line string) (name string) {
	parts := strings.Fields(line)
	for _, part := range parts {
		if strings.HasPrefix(part, "enum") || strings.HasPrefix(part, "class") {
			continue
		}

		if strings.HasPrefix(part, "E") {
			name = part
		}
	}
	return
}

func extractFunctionName(line string) (fnName string) {
	openBracketIndex := strings.Index(line, "(")
	if openBracketIndex == -1 {
		fnName = "INVALID METHOD NAME"
		return
	}

	lastIndex := strings.LastIndex(line[0:openBracketIndex], " ")
	if lastIndex >= 0 {
		fnName = strings.TrimSpace(line[lastIndex:openBracketIndex])
		return
	}

	fnName = strings.TrimSpace(line[0:openBracketIndex])
	return
}

func cleanComment(line string) (comment string) {
	comment = strings.TrimLeft(line, "//")
	comment = strings.TrimLeft(comment, "/*")
	comment = strings.TrimLeft(comment, "/**")
	comment = strings.TrimLeft(comment, "*")
	comment = strings.TrimRight(comment, "*/")
	comment = strings.TrimSpace(comment)
	return
}

func isCopy(line string) bool {
	return strings.HasPrefix(line, "// Copy")
}

func isClassMacro(line string) bool {
	return strings.HasPrefix(line, "UCLASS")
}

func isClass(line string) bool {
	return strings.HasPrefix(line, "class")
}

func isEnumMacro(line string) bool {
	return strings.HasPrefix(line, "UENUM")
}

func isEnum(line string) bool {
	return strings.HasPrefix(line, "enum")
}

func isStructMacro(line string) bool {
	return strings.HasPrefix(line, "USTRUCT")
}

func isStruct(line string) bool {
	return strings.HasPrefix(line, "struct")
}

func isPropertyMacro(line string) bool {
	return strings.HasPrefix(line, "UPROPERTY")
}

func isProperty(line string, prevIsProp bool) bool {
	if prevIsProp {
		return true
	}
	return !strings.Contains(line, "(") && strings.HasSuffix(line, ";")
}

func isFunctionMacro(line string) bool {
	return strings.HasPrefix(line, "UFUNCTION")
}

func isOpenBracket(line string) bool {
	return strings.EqualFold(line, "{")
}

func isCloseBracket(line string) bool {
	return strings.EqualFold(line, "};")
}

func isFunction(line string, prevIsFn bool) bool {
	if prevIsFn {
		return true
	}
	return strings.Contains(line, "(") && strings.Contains(line, ")") && strings.HasSuffix(line, ";")
}

func isAccessModifier(line string) bool {
	return strings.HasPrefix(line, "public:") || strings.HasPrefix(line, "protected:") || strings.HasPrefix(line, "private:")
}

func keepExistingMarkdown(sourceFile, destFolder string) (existingMarkdown string, hasDefinitionHeader bool) {
	fileName := filepath.Base(sourceFile)
	outputPath := filepath.Join(destFolder, strings.TrimSuffix(fileName, filepath.Ext(fileName))+".mdx")

	file, err := os.Open(outputPath)
	if err != nil {
		fmt.Printf("Could not open file %s: %v. Will be created as new\n", outputPath, err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	i := 0
	for scanner.Scan() {

		line := strings.TrimSpace(scanner.Text())
		if i < 4 {
			i++
			continue
		}

		existingMarkdown += line + "\n"

		if strings.Contains(line, "## File Info") {
			hasDefinitionHeader = true
			i++
			break
		}

		i++
	}

	return
}

func accessModifierString(accessType AccessType) string {
	if accessType == Public {
		return "Public"
	} else if accessType == Protected {
		return "Protected"
	}
	return "Private"
}

func outputMarkdown(fileInfo *FileInfo, destFolder string) {
	fileName := filepath.Base(fileInfo.Path)
	outputPath := filepath.Join(destFolder, strings.TrimSuffix(fileName, filepath.Ext(fileName))+".mdx")

	var keepContent, hasDefinitionHeader = keepExistingMarkdown(fileInfo.Path, destFolder)

	file, err := os.Create(outputPath)
	if err != nil {
		fmt.Printf("Error opening output file %s: %v\n", outputPath, err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// Header Page
	writer.WriteString("---\n")
	writer.WriteString("title: " + fileInfo.Name + "\n")
	writer.WriteString("description: Reference page for " + fileInfo.Name + "\n")
	writer.WriteString("---\n")

	// Keep Existing Content
	writer.WriteString(keepContent)

	// Go on to Autogenerated Content
	if !hasDefinitionHeader {
		writer.WriteString("\n## File Info\n\n")
	}

	writer.WriteString("\n __FileName:__ `" + fileName + "`\n\n")

	for _, class := range fileInfo.Data {
		writer.WriteString(fmt.Sprintf("- [`" + class.Name + "`](#" + class.Name + ") \n"))
	}

	for _, class := range fileInfo.Data {
		writer.WriteString(fmt.Sprintf("\n### `" + class.Name + "` \n\n"))

		for _, parent := range class.Parents {
			writer.WriteString(fmt.Sprintf("- __Parent Class:__ `%s`\n", parent))
		}

		// for i, comm := range info.Comments {
		// 	isLast := i == max(len(info.Comments)-1, 0)
		// 	if isLast {
		// 		writer.WriteString("> " + cleanComment(comm) + " \n")
		// 	} else {
		// 		writer.WriteString("> " + cleanComment(comm) + " \\\n")
		// 	}
		// }

		writer.WriteString("\n## Properties\n\n")
		if len(class.Properties) == 0 {
			writer.WriteString("No documented properties available\n")
		} else {
			writer.WriteString("```cpp\n")
			for _, prop := range class.Properties {
				for _, comm := range prop.Comments {
					writer.WriteString("// " + cleanComment(comm) + " \n")
				}
				writer.WriteString(prop.Macro + "\n")
				writer.WriteString(prop.Declaration + "\n\n")
			}
			writer.WriteString("```\n")
		}

		writer.WriteString("\n## Functions\n\n")
		if len(class.Functions) == 0 {
			writer.WriteString("No documented functions available\n")
		} else {
			for _, function := range class.Functions {
				writer.WriteString("### `" + function.Name + "`\n")
				for i, comm := range function.Comments {
					isLast := i == max(len(function.Comments)-1, 0)
					if isLast {
						writer.WriteString("> " + cleanComment(comm) + " \n")
					} else {
						writer.WriteString("> " + cleanComment(comm) + " \\\n")
					}
				}
				writer.WriteString("```cpp\n")
				writer.WriteString(function.Declaration + "\n\n")
				writer.WriteString("```\n\n")
			}
		}
	}

	writer.Flush()
	fmt.Printf("Generated markdown file: %s\n", outputPath)
}
