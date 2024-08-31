package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ClassInfo struct {
	Name       string
	ParentName string
	Properties []PropertyInfo
	Methods    []MethodInfo
}

type PropertyInfo struct {
	Declaration string
	Macro       string
	Comments    []string
}

type MethodInfo struct {
	Declaration string
	Comments    []string
}

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: program <source_folder> <destination_folder>")
		os.Exit(1)
	}

	sourceFolder := os.Args[1]
	destFolder := os.Args[2]

	err := filepath.Walk(sourceFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && (strings.HasSuffix(info.Name(), ".h") || strings.HasSuffix(info.Name(), ".hpp")) {
			processFile(path, destFolder)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking through directory: %v\n", err)
	}
}

func processFile(filePath, destFolder string) {
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

	outputMarkdown(classInfo, filePath, destFolder)
}

func extractInfo(file *os.File) ClassInfo {
	var info ClassInfo
	scanner := bufio.NewScanner(file)

	var currentComments []string
	var inMultiLineComment bool
	var multiLineComment string
	var inClassDeclaration bool
	var lastUPROPERTY string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if inMultiLineComment {
			if strings.Contains(line, "*/") {
				multiLineComment += line[:strings.Index(line, "*/")+2]
				currentComments = append(currentComments, multiLineComment)
				inMultiLineComment = false
				multiLineComment = ""
			} else {
				multiLineComment += line + "\n"
			}
			continue
		}

		if line == "" {
			continue
		}

		if strings.Contains(line, "Copyright") {
			continue
		}

		if strings.Contains(line, "// UFlowPilotTask") {
			continue
		}

		if strings.Contains(line, "//~UFlowPilotTask") {
			continue
		}

		if strings.HasPrefix(line, "//") {
			currentComments = append(currentComments, line)
		} else if strings.HasPrefix(line, "/*") {
			if strings.Contains(line, "*/") {
				currentComments = append(currentComments, line)
			} else {
				inMultiLineComment = true
				multiLineComment = line + "\n"
			}
		} else if strings.HasPrefix(line, "UPROPERTY(") {
			lastUPROPERTY = line
		} else if strings.HasPrefix(line, "class") {
			parts := strings.Fields(line)
			if len(parts) > 1 {
				if strings.Contains(parts[1], "_API") {
					info.Name = parts[2]
				}

				if len(parts) > 3 && parts[3] == ":" {
					info.ParentName = parts[5]
				}
				inClassDeclaration = true
			}
		} else if inClassDeclaration {
			if isProperty(line) && lastUPROPERTY != "" {
				info.Properties = append(info.Properties, PropertyInfo{
					Declaration: line,
					Macro:       lastUPROPERTY,
					Comments:    currentComments,
				})
				currentComments = nil
				lastUPROPERTY = ""
			} else if isMethod(line) {
				info.Methods = append(info.Methods, MethodInfo{
					Declaration: line,
					Comments:    currentComments,
				})
				currentComments = nil
			}
			if !isProperty(line) && !isMethod(line) {
				lastUPROPERTY = ""
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
	}

	return info
}

func isProperty(line string) bool {
	return !strings.Contains(line, "(") && strings.HasSuffix(line, ";")
}

func isMethod(line string) bool {
	return strings.Contains(line, "(") && strings.HasSuffix(line, ";")
}

func outputMarkdown(info ClassInfo, sourceFile, destFolder string) {
	fileName := filepath.Base(sourceFile)
	outputPath := filepath.Join(destFolder, strings.TrimSuffix(fileName, filepath.Ext(fileName))+".md")

	file, err := os.Create(outputPath)
	if err != nil {
		fmt.Printf("Error creating output file %s: %v\n", outputPath, err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	if info.ParentName != "" {
		writer.WriteString(fmt.Sprintf("# `%s` : `%s`\n\n", info.Name, info.ParentName))
	} else {
		writer.WriteString(fmt.Sprintf("# `%s`\n\n", info.Name))
	}

	writer.WriteString("## Properties\n\n")
	writer.WriteString("| Property | Description |\n")
	writer.WriteString("|----------|-------------|\n")
	for _, prop := range info.Properties {
		description := strings.Join(prop.Comments, " ")
		description = strings.TrimPrefix(description, "//")
		description = strings.TrimPrefix(description, "/*")
		description = strings.TrimSuffix(description, "*/")
		description = strings.TrimSpace(description)

		writer.WriteString(fmt.Sprintf("| `%s`<br>`%s` | %s |\n", prop.Macro, prop.Declaration, description))
	}
	writer.WriteString("\n")

	writer.WriteString("## Methods\n\n")
	for _, method := range info.Methods {
		writer.WriteString("```cpp\n")
		writer.WriteString(method.Declaration + "\n")
		writer.WriteString("```\n\n")
		for _, comment := range method.Comments {
			writer.WriteString(strings.TrimPrefix(comment, "//") + "\n")
		}
		writer.WriteString("\n")
	}

	writer.Flush()
	fmt.Printf("Generated markdown file: %s\n", outputPath)
}
