package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	args := os.Args
	if len(args) < 3 {
		fmt.Println("Please provide the component name and directory")
		os.Exit(1)
	}

	name := args[1]
	currentDir := args[2]

	err := createComponent(name, currentDir)
	if err != nil {
		fmt.Println("Error creating component:", err)
		os.Exit(1)
	}

	err = updateMainIndex(name, currentDir)
	if err != nil {
		fmt.Println("Error updating main index:", err)
		os.Exit(1)
	}
}

func capitalizeFirstLetter(s string) string {
	return strings.ToUpper(s[:1]) + s[1:]
}

func createComponent(name, currentDir string) error {
	capitalizedName := capitalizeFirstLetter(name)
	dirPath := fmt.Sprintf("%s/%s", currentDir, name)
	componentPath := fmt.Sprintf("%s/%s.tsx", dirPath, name)
	sassPath := fmt.Sprintf("%s/%s.module.scss", dirPath, name)
	indexPath := fmt.Sprintf("%s/index.ts", dirPath)
	storyPath := fmt.Sprintf("%s/%s.stories.tsx", dirPath, name)

	// Create directory if not exists
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		err := os.Mkdir(dirPath, 0755)
		if err != nil {
			return err
		}
	}

	// Component content
	componentContent := fmt.Sprintf(`import React from 'react'
import s from './%s.module.scss'
export type %sProps = {}

export const %s: React.FC<%sProps> = ({}) => {
  return <div className={s.container}>%s</div>
}
`, name, capitalizedName, capitalizedName, capitalizedName, capitalizedName)

	// SASS content
	sassContent := `.container {
  // styles go here
}`

	// Index content
	indexContent := fmt.Sprintf(`export * from './%s'`, name)

	// Story content
	storyContent := fmt.Sprintf(`import type { Meta, StoryObj } from '@storybook/react'
import { %s } from './'

const meta = {
  component: %s,
  tags: ['autodocs'],
  title: 'Components/%s',
} satisfies Meta<typeof %s>

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {
  args: {},
}
`, capitalizedName, capitalizedName, capitalizedName, capitalizedName)

	// Write files
	err := writeToFile(componentPath, componentContent)
	if err != nil {
		return err
	}
	err = writeToFile(sassPath, sassContent)
	if err != nil {
		return err
	}
	err = writeToFile(indexPath, indexContent)
	if err != nil {
		return err
	}
	err = writeToFile(storyPath, storyContent)
	if err != nil {
		return err
	}

	// Execute formatting and linting commands
	runCommand(currentDir, "pnpm", "run", "format:file", dirPath)
	runCommand(currentDir, "pnpm", "run", "lint:file", dirPath+"/**")

	return nil
}

func writeToFile(path, content string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(content)
	if err != nil {
		return err
	}

	return writer.Flush()
}

func updateMainIndex(name, currentDir string) error {
	mainIndexPath := fmt.Sprintf("%s/index.ts", currentDir)
	lineToAdd := fmt.Sprintf(`export * from './%s'`, name)

	// Read and update file content
	content, err := os.ReadFile(mainIndexPath)
	if err != nil {
		return err
	}

	if !strings.Contains(string(content), lineToAdd) {
		newContent := lineToAdd + "\n" + string(content)
		err = os.WriteFile(mainIndexPath, []byte(newContent), 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

func runCommand(currentDir, command string, args ...string) {
	cmd := exec.Command(command, args...)
	cmd.Dir = currentDir // Set the working directory
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error executing command:", err)
	}
	fmt.Println(string(output))
}
