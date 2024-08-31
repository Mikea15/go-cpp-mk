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
	Function
	Property
	AccessModifier
)

type FileOutput struct {
	Path string
	Name string
}

type ClassInfo struct {
	Name       string
	ParentName string
	Comments   string
	Properties []PropertyInfo
	Functions  []FunctionInfo
}

type PropertyInfo struct {
	Declaration string
	Comments    string
}

type FunctionInfo struct {
	Name        string
	Declaration string
	Comments    string
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: program <source_folder> <destination_folder>")
		os.Exit(1)
	}

	sourceFolder := os.Args[1]
	destFolder := os.Args[2]

	var fileList []FileOutput
	// var classList []ClassInfo

	err := filepath.Walk(sourceFolder, func(path string, info os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if strings.HasSuffix(info.Name(), ".h") || strings.HasSuffix(info.Name(), ".hpp") {
			fOutput := FileOutput{path, info.Name()}
			fileList = append(fileList, fOutput)
		}

		return nil
	})

	for i := 0; i < len(fileList); i++ {
		processFile(fileList[i].Path, destFolder, i)
	}

	if err != nil {
		fmt.Printf("Error walking through directory: %v\n", err)
	}
}

func processFile(filePath, destFolder string, fileIndex int) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Error opening file %s: %v\n", filePath, err)
		return
	}
	defer file.Close()

	classInfo := extractInfo(file)
	if classInfo.Name == "" {
		fmt.Printf("No class found in file %s\n", filePath)
		return
	}

	// classInfo.FileIndex = fileIndex

	outputMarkdown(classInfo, filePath, destFolder)
}

func extractInfo(file *os.File) ClassInfo {
	var info ClassInfo
	scanner := bufio.NewScanner(file)

	var commentStack string
	var functionStack string
	var propertyStack string

	ignoreList := []string{
		"#",
		"// UFlowPilotTask",
		"//~UFlowPilotTask",
	}

	var prevId LineId = Empty
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		var skip = false
		for i := 0; i < len(ignoreList); i++ {
			if strings.Contains(line, ignoreList[i]) {
				skip = true
				break
			}
		}

		if skip {
			continue
		}

		id := idLine(line, prevId)

		// fmt.Printf("[%d] %s\n", id, line)

		switch id {
		case Empty:
			continue
		case Comment:
			// stack comments
			commentStack += cleanComment(line) + "\n"
		case Class:
			if !isClassMacro(line) {
				info.Comments = commentStack
				commentStack = ""

				var name, parent = extractClassInfo(line)
				info.Name = name
				info.ParentName = parent
			}
		case Struct:

		case Function:
			if !isFunctionMacro(line) {
				var data = FunctionInfo{
					Name:        extractFunctionName(line),
					Declaration: functionStack + line,
					Comments:    commentStack,
				}
				commentStack = ""
				functionStack = ""

				info.Functions = append(info.Functions, data)
			} else {
				functionStack = line + "\n"
			}
		case Property:
			propertyStack += line + "\n"

			if !isPropertyMacro(line) {
				var data = PropertyInfo{
					Declaration: propertyStack,
					Comments:    commentStack,
				}
				commentStack = ""
				propertyStack = ""

				info.Properties = append(info.Properties, data)
			}
		case AccessModifier:

		default:
		}

		prevId = id
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
	}

	return info
}

func idLine(line string, prevId LineId) LineId {

	if len(line) == 0 || isCopy(line) {
		return Empty
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

	if isFunctionMacro(line) || isFunction(line) {
		return Function
	}

	if isPropertyMacro(line) || isProperty(line) {
		return Property
	}

	if isAccessModifier(line) {
		return AccessModifier
	}

	return Empty
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

func extractClassInfo(line string) (class, parent string) {
	parts := strings.Fields(line)
	if len(parts) > 1 {
		if strings.Contains(parts[1], "_API") {
			class = strings.TrimSpace(parts[2])
		}

		if len(parts) > 3 && parts[3] == ":" {
			parent = strings.TrimSpace(parts[5])
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

func isStructMacro(line string) bool {
	return strings.HasPrefix(line, "USTRUCT")
}

func isStruct(line string) bool {
	return strings.HasPrefix(line, "struct")
}

func isPropertyMacro(line string) bool {
	return strings.HasPrefix(line, "UPROPERTY")
}

func isProperty(line string) bool {
	return !strings.Contains(line, "(") && strings.HasSuffix(line, ";")
}

func isFunctionMacro(line string) bool {
	return strings.HasPrefix(line, "UFUNCTION")
}

func isFunction(line string) bool {
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
		fmt.Printf("Error opening output file %s: %v\n", outputPath, err)
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

		if strings.Contains(line, "## Class Info") {
			hasDefinitionHeader = true
			i++
			break
		}

		i++
	}

	return
}

func outputMarkdown(info ClassInfo, sourceFile, destFolder string) {
	fileName := filepath.Base(sourceFile)
	outputPath := filepath.Join(destFolder, strings.TrimSuffix(fileName, filepath.Ext(fileName))+".mdx")

	var keepContent, hasDefinitionHeader = keepExistingMarkdown(sourceFile, destFolder)

	file, err := os.Create(outputPath)
	if err != nil {
		fmt.Printf("Error opening output file %s: %v\n", outputPath, err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	// Header Page
	writer.WriteString("---\n")
	writer.WriteString("title: " + info.Name + "\n")
	writer.WriteString("description: Reference page for " + info.Name + " class\n")
	writer.WriteString("---\n")

	// Keep Existing Content
	writer.WriteString(keepContent)

	// Go on to Autogenerated Content
	if !hasDefinitionHeader {
		writer.WriteString("\n## Class Info\n\n")
	}

	if info.ParentName != "" {
		writer.WriteString(fmt.Sprintf("__Parent Class:__ `%s`\n", info.ParentName))
	}

	// writer.WriteString("## Properties\n\n")
	// writer.WriteString("| Property | Description |\n")
	// writer.WriteString("|----------|-------------|\n")
	// for _, prop := range info.Properties {
	// 	writer.WriteString(fmt.Sprintf("| `%s` | %s |\n", prop.Declaration, prop.Comments))
	// }
	// writer.WriteString("\n")

	writer.WriteString("### Functions\n\n")
	for _, function := range info.Functions {
		writer.WriteString("#### `" + function.Name + "` \n")
		writer.WriteString(function.Comments + "\n")
		writer.WriteString("```cpp\n")
		writer.WriteString(function.Declaration + "\n")
		writer.WriteString("```\n\n")
	}

	writer.Flush()
	fmt.Printf("Generated markdown file: %s\n", outputPath)
}
