package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/shenma/code-completions-go/pkg/tokenizers"
)

// StatsResult represents the statistics result for a file
type StatsResult struct {
	FilePath   string `json:"file_path"`
	TokenCount int    `json:"token_count"`
	LineCount  int    `json:"line_count"`
	Error      string `json:"error,omitempty"`
}

// Config represents the application configuration
type Config struct {
	DataDir       string `json:"data_dir"`
	Output        string `json:"output"`
	TokenizerPath string `json:"tokenizer_path"`
}

/**
 *	获取costrict目录结构设定
 */
func getCostrictDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}
	return filepath.Join(homeDir, ".costrict")
}

var rootCmd = &cobra.Command{
	Use:   "token-stats",
	Short: "Token statistics tool for files",
	Long:  `Token-stats is a tool that counts tokens in files using a specified tokenizer.`,
	Run: func(cmd *cobra.Command, args []string) {
		dataDir, _ := cmd.Flags().GetString("data-dir")
		output, _ := cmd.Flags().GetString("output")

		config := Config{
			DataDir:       dataDir,
			Output:        output,
			TokenizerPath: filepath.Join(getCostrictDir(), "config", "tokenizer.json"),
		}

		if err := runStats(config); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.Flags().StringP("data-dir", "d", "", "Directory containing files to analyze")
	rootCmd.Flags().StringP("output", "o", "", "Output file for results (default: stdout)")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

/**
 * Run token statistics on files in the specified directory
 * @param {Config} config - Application configuration including data directory and output file
 * @returns {error} Returns error if any step fails, nil on success
 * @description
 * - Creates a tokenizer instance using the specified tokenizer path
 * - Walks through the data directory to find all files
 * - Counts tokens for each file using the tokenizer
 * - Outputs results to specified file or stdout
 * @throws
 * - Tokenizer creation failure (NewTokenizer)
 * - Directory access errors (filepath.Walk)
 * - File reading errors (ioutil.ReadFile)
 * - JSON marshaling errors (json.Marshal)
 * - File writing errors (ioutil.WriteFile)
 */
func runStats(config Config) error {
	// Create tokenizer
	tokenizer, err := tokenizers.NewTokenizer(config.TokenizerPath)
	if err != nil {
		return fmt.Errorf("failed to create tokenizer: %v", err)
	}
	defer tokenizer.Close()

	// Collect results
	var results []StatsResult

	// Walk through the directory
	err = filepath.Walk(config.DataDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Process file
		result := processFile(path, tokenizer)
		results = append(results, result)

		return nil
	})

	if err != nil {
		return fmt.Errorf("error walking through directory: %v", err)
	}

	// Output results
	return outputResults(results, config.Output)
}

/**
 * Process a single file to count tokens and lines
 * @param {string} filePath - Path to the file to process
 * @param {*tokenizers.Tokenizer} tokenizer - Tokenizer instance to use for counting
 * @returns {StatsResult} Returns statistics result containing file path, token count, line count, and any error
 * @description
 * - Reads the file content
 * - Counts tokens using the provided tokenizer
 * - Counts lines in the file
 * - Returns result with error if file reading or tokenization fails
 * @throws
 * - File reading errors (ioutil.ReadFile)
 */
func processFile(filePath string, tokenizer *tokenizers.Tokenizer) StatsResult {
	result := StatsResult{
		FilePath: filePath,
	}

	// Read file
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		result.Error = fmt.Sprintf("failed to read file: %v", err)
		return result
	}

	// Convert content to string and count tokens
	text := string(content)
	tokenCount := tokenizer.GetTokenCount(text)
	result.TokenCount = tokenCount

	// Count lines
	lines := countLines(text)
	result.LineCount = lines

	return result
}

/**
 * Count lines in a text string
 * @param {string} text - Text content to count lines for
 * @returns {int} Returns the number of lines in the text
 * @description
 * - Splits text by newline characters
 * - Counts resulting lines
 * - Handles both Unix (LF) and Windows (CRLF) line endings
 */
func countLines(text string) int {
	if text == "" {
		return 0
	}

	// Count lines by splitting on newline characters
	// This handles both Unix (LF) and Windows (CRLF) line endings
	lines := strings.Split(text, "\n")

	// If the last line is empty, don't count it
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		return len(lines) - 1
	}

	return len(lines)
}

/**
 * Output results to specified file or stdout
 * @param {[]StatsResult} results - Array of statistics results to output
 * @param {string} outputPath - Path to output file, empty for stdout
 * @returns {error} Returns error if JSON marshaling or file writing fails, nil on success
 * @description
 * - Marshals results to JSON format
 * - Writes to specified file if outputPath is provided
 * - Prints to stdout if outputPath is empty
 * @throws
 * - JSON marshaling errors (json.MarshalIndent)
 * - File writing errors (ioutil.WriteFile)
 */
func outputResults(results []StatsResult, outputPath string) error {
	// Marshal results to JSON
	jsonData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal results: %v", err)
	}

	// Add newline at the end
	jsonData = append(jsonData, '\n')

	// Output to file or stdout
	if outputPath != "" {
		err = ioutil.WriteFile(outputPath, jsonData, 0644)
		if err != nil {
			return fmt.Errorf("failed to write output file: %v", err)
		}
		fmt.Printf("Results written to %s\n", outputPath)
	} else {
		fmt.Println(string(jsonData))
	}

	return nil
}

// PreloadedTokenizerJSON contains the preloaded tokenizer.json content
// This constant will be replaced with the actual tokenizer.json content
