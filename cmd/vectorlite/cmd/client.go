package cmd

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var serverURL string
var selectedDatabase string

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Interactive client for VectorLite server",
	Long:  `Start an interactive shell to connect and interact with a VectorLite server.`,
	Run: func(cmd *cobra.Command, args []string) {
		runInteractiveClient()
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)
	clientCmd.Flags().StringVar(&serverURL, "server", "http://localhost:9123", "VectorLite server URL")
}

func runInteractiveClient() {
	fmt.Printf("VectorLite Interactive Client\n")
	fmt.Printf("Connected to: %s\n", serverURL)
	fmt.Printf("Type 'help' for commands or 'quit' to exit\n\n")

	scanner := bufio.NewScanner(os.Stdin)
	
	for {
		prompt := "vectorlite"
		if selectedDatabase != "" {
			prompt += "[" + selectedDatabase + "]"
		}
		fmt.Print(prompt + "> ")
		if !scanner.Scan() {
			break
		}
		
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		
		if line == "quit" || line == "exit" {
			fmt.Println("Goodbye!")
			break
		}
		
		handleCommand(line)
	}
}

func handleCommand(input string) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return
	}
	
	command := parts[0]
	
	switch command {
	case "help":
		showHelp()
	case "status":
		checkServerStatus()
	case "create-db":
		handleCreateDatabase(parts[1:])
	case "use-db":
		handleUseDatabase(parts[1:])
	case "list-dbs":
		handleListDatabases()
	case "add":
		handleAddEntry(parts[1:])
	case "query":
		handleQuery(parts[1:])
	case "list":
		handleListEntries()
	case "import":
		handleImport(parts[1:])
	default:
		fmt.Printf("Unknown command: %s. Type 'help' for available commands.\n", command)
	}
}

func showHelp() {
	fmt.Println("Available commands:")
	fmt.Println("  help                          - Show this help message")
	fmt.Println("  status                        - Check server connection")
	fmt.Println("  create-db <name> <algorithm>  - Create a new database")
	fmt.Println("    Example: create-db mydb bruteforce")
	fmt.Println("    Algorithms: bruteforce, hnsw")
	fmt.Println("  use-db <name>                 - Select database to use")
	fmt.Println("    Example: use-db mydb")
	fmt.Println("  list-dbs                      - List all databases")
	fmt.Println("  add <vector> <metadata>       - Add vector entry")
	fmt.Println("    Example: add [1.0,2.0,3.0] name=test,type=example")
	fmt.Println("  query <vector> <k> <metric>   - Query similar vectors")
	fmt.Println("    Example: query [1.0,2.0,3.0] 5 cosine")
	fmt.Println("    Metrics: cosine, dot_product, euclidean")
	fmt.Println("  import <file>                 - Import vectors from file")
	fmt.Println("    Example: import vectors.csv")
	fmt.Println("    Supported formats: CSV")
	fmt.Println("  list                          - List all entries")
	fmt.Println("  quit/exit                     - Exit the client")
}

func checkServerStatus() {
	resp, err := http.Get(serverURL + "/entries")
	if err != nil {
		fmt.Printf("Error connecting to server: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	if resp.StatusCode == http.StatusOK {
		fmt.Println("Server is running and accessible")
	} else {
		fmt.Printf("Server responded with status: %d\n", resp.StatusCode)
	}
}

func handleAddEntry(args []string) {
	if selectedDatabase == "" {
		fmt.Println("Error: No database selected. Use 'use-db <name>' to select a database first.")
		return
	}
	
	if len(args) < 2 {
		fmt.Println("Usage: add <vector> <metadata>")
		fmt.Println("Example: add [1.0,2.0,3.0] name=test,type=example")
		return
	}
	
	// Parse vector
	vectorStr := args[0]
	vector, err := parseVector(vectorStr)
	if err != nil {
		fmt.Printf("Error parsing vector: %v\n", err)
		return
	}
	
	// Parse metadata
	metadataStr := args[1]
	metadata, err := parseMetadata(metadataStr)
	if err != nil {
		fmt.Printf("Error parsing metadata: %v\n", err)
		return
	}
	
	// Create request
	reqData := map[string]interface{}{
		"database":  selectedDatabase,
		"vectors":   [][]float64{vector},
		"metadatas": []map[string]string{metadata},
	}
	
	err = makePostRequest("/entries", reqData)
	if err != nil {
		fmt.Printf("Error adding entry: %v\n", err)
		return
	}
	
	fmt.Println("Entry added successfully")
}

func handleQuery(args []string) {
	if selectedDatabase == "" {
		fmt.Println("Error: No database selected. Use 'use-db <name>' to select a database first.")
		return
	}
	
	if len(args) < 3 {
		fmt.Println("Usage: query <vector> <k> <metric>")
		fmt.Println("Example: query [1.0,2.0,3.0] 5 cosine")
		fmt.Println("Metrics: cosine, dot_product, euclidean")
		return
	}
	
	// Parse vector
	vector, err := parseVector(args[0])
	if err != nil {
		fmt.Printf("Error parsing vector: %v\n", err)
		return
	}
	
	// Parse k
	k, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Printf("Error parsing k: %v\n", err)
		return
	}
	
	// Get metric
	metric := args[2]
	if metric != "cosine" && metric != "dot_product" && metric != "euclidean" {
		fmt.Printf("Invalid metric: %s. Use: cosine, dot_product, or euclidean\n", metric)
		return
	}
	
	// Create request
	reqData := map[string]interface{}{
		"database": selectedDatabase,
		"vector":   vector,
		"k":        k,
		"metric":   metric,
	}
	
	resp, err := makePostRequestWithResponse("/query", reqData)
	if err != nil {
		fmt.Printf("Error querying: %v\n", err)
		return
	}
	
	// Parse and display results
	var result map[string]interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		fmt.Printf("Error parsing response: %v\n", err)
		return
	}
	
	entries, ok := result["entries"].([]interface{})
	if !ok {
		fmt.Println("No entries found in response")
		return
	}
	
	fmt.Printf("Found %d similar entries:\n", len(entries))
	for i, entry := range entries {
		entryMap := entry.(map[string]interface{})
		fmt.Printf("%d. ID: %.0f\n", i+1, entryMap["id"])
		fmt.Printf("   Vector: %v\n", entryMap["vector"])
		fmt.Printf("   Metadata: %v\n", entryMap["metadata"])
		fmt.Println()
	}
}

func handleListEntries() {
	if selectedDatabase == "" {
		fmt.Println("Error: No database selected. Use 'use-db <name>' to select a database first.")
		return
	}
	
	url := fmt.Sprintf("%s/entries?database=%s", serverURL, selectedDatabase)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error listing entries: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}
	
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Printf("Error parsing response: %v\n", err)
		return
	}
	
	entries, ok := result["entries"].([]interface{})
	if !ok {
		fmt.Println("No entries found")
		return
	}
	
	fmt.Printf("Total entries: %d\n", len(entries))
	for i, entry := range entries {
		entryMap := entry.(map[string]interface{})
		fmt.Printf("%d. ID: %.0f\n", i+1, entryMap["id"])
		fmt.Printf("   Vector: %v\n", entryMap["vector"])
		fmt.Printf("   Metadata: %v\n", entryMap["metadata"])
		fmt.Println()
	}
}

func parseVector(vectorStr string) ([]float64, error) {
	// Remove brackets and split by comma
	vectorStr = strings.Trim(vectorStr, "[]")
	parts := strings.Split(vectorStr, ",")
	
	vector := make([]float64, len(parts))
	for i, part := range parts {
		val, err := strconv.ParseFloat(strings.TrimSpace(part), 64)
		if err != nil {
			return nil, fmt.Errorf("invalid float value: %s", part)
		}
		vector[i] = val
	}
	
	return vector, nil
}

func parseMetadata(metadataStr string) (map[string]string, error) {
	metadata := make(map[string]string)
	
	pairs := strings.Split(metadataStr, ",")
	for _, pair := range pairs {
		kv := strings.Split(pair, "=")
		if len(kv) != 2 {
			return nil, fmt.Errorf("invalid metadata format: %s", pair)
		}
		metadata[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
	}
	
	return metadata, nil
}

func makePostRequest(endpoint string, data interface{}) error {
	_, err := makePostRequestWithResponse(endpoint, data)
	return err
}

func makePostRequestWithResponse(endpoint string, data interface{}) ([]byte, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	
	resp, err := http.Post(serverURL+endpoint, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("server error: %d - %s", resp.StatusCode, string(body))
	}
	
	return body, nil
}

func handleImport(args []string) {
	if selectedDatabase == "" {
		fmt.Println("Error: No database selected. Use 'use-db <name>' to select a database first.")
		return
	}
	
	if len(args) < 1 {
		fmt.Println("Usage: import <file>")
		fmt.Println("Example: import vectors.csv")
		fmt.Println("Supported formats: CSV")
		return
	}
	
	filename := args[0]
	
	// Determine file format (extensible for future formats)
	format := detectFileFormat(filename)
	if format == "" {
		fmt.Printf("Unsupported file format: %s\n", filepath.Ext(filename))
		fmt.Println("Supported formats: .csv")
		return
	}
	
	// Import based on format
	switch format {
	case "csv":
		err := importCSV(filename)
		if err != nil {
			fmt.Printf("Error importing CSV: %v\n", err)
			return
		}
	default:
		fmt.Printf("Import for format '%s' not implemented yet\n", format)
		return
	}
	
	fmt.Printf("Successfully imported vectors from %s\n", filename)
}

func detectFileFormat(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	
	switch ext {
	case ".csv":
		return "csv"
	// Future formats can be added here:
	// case ".json":
	//     return "json"
	// case ".jsonl":
	//     return "jsonl"
	default:
		return ""
	}
}

func importCSV(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()
	
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV: %v", err)
	}
	
	if len(records) == 0 {
		return fmt.Errorf("CSV file is empty")
	}
	
	// Parse records into vectors and metadata
	var vectors [][]float64
	var metadatas []map[string]string
	
	// Check if first row is header by trying to parse as numbers
	startRow := 0
	if len(records) > 0 {
		isHeader := false
		for _, cell := range records[0] {
			if _, err := strconv.ParseFloat(cell, 64); err != nil {
				isHeader = true
				break
			}
		}
		if isHeader {
			startRow = 1
			fmt.Println("Detected header row, skipping...")
		}
	}
	
	fmt.Printf("Processing %d rows...\n", len(records)-startRow)
	
	for i := startRow; i < len(records); i++ {
		record := records[i]
		if len(record) == 0 {
			continue
		}
		
		vector, metadata, err := parseCSVRecord(record)
		if err != nil {
			fmt.Printf("Warning: skipping row %d: %v\n", i+1, err)
			continue
		}
		
		vectors = append(vectors, vector)
		metadatas = append(metadatas, metadata)
		
		// Progress reporting for large files
		if (i-startRow+1)%100 == 0 {
			fmt.Printf("Processed %d rows...\n", i-startRow+1)
		}
	}
	
	if len(vectors) == 0 {
		return fmt.Errorf("no valid vectors found in CSV")
	}
	
	// Send in batches for better performance
	batchSize := 100
	totalBatches := (len(vectors) + batchSize - 1) / batchSize
	
	fmt.Printf("Importing %d vectors in %d batches...\n", len(vectors), totalBatches)
	
	for i := 0; i < len(vectors); i += batchSize {
		end := i + batchSize
		if end > len(vectors) {
			end = len(vectors)
		}
		
		batchVectors := vectors[i:end]
		batchMetadatas := metadatas[i:end]
		
		reqData := map[string]interface{}{
			"database":  selectedDatabase,
			"vectors":   batchVectors,
			"metadatas": batchMetadatas,
		}
		
		err := makePostRequest("/entries", reqData)
		if err != nil {
			return fmt.Errorf("failed to import batch %d-%d: %v", i+1, end, err)
		}
		
		fmt.Printf("Imported batch %d/%d (%d vectors)\n", (i/batchSize)+1, totalBatches, end-i)
	}
	
	return nil
}

func parseCSVRecord(record []string) ([]float64, map[string]string, error) {
	if len(record) == 0 {
		return nil, nil, fmt.Errorf("empty record")
	}
	
	var vector []float64
	metadata := make(map[string]string)
	
	// Parse each field - numeric values go to vector, others to metadata
	for i, cell := range record {
		cell = strings.TrimSpace(cell)
		
		if val, err := strconv.ParseFloat(cell, 64); err == nil {
			// It's a number, add to vector
			vector = append(vector, val)
		} else {
			// It's not a number, treat as metadata
			if strings.Contains(cell, "=") {
				// Parse key=value format
				kv := strings.SplitN(cell, "=", 2)
				if len(kv) == 2 {
					metadata[strings.TrimSpace(kv[0])] = strings.TrimSpace(kv[1])
				}
			} else {
				// Use column index as key for non key=value format
				metadata[fmt.Sprintf("col_%d", i)] = cell
			}
		}
	}
	
	// If no vector values found, return error
	if len(vector) == 0 {
		return nil, nil, fmt.Errorf("no numeric values found for vector")
	}
	
	// Set default metadata if none provided
	if len(metadata) == 0 {
		metadata["imported"] = "true"
		metadata["source"] = "csv"
	}
	
	return vector, metadata, nil
}

func handleCreateDatabase(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: create-db <name> <algorithm>")
		fmt.Println("Example: create-db mydb bruteforce")
		fmt.Println("Algorithms: bruteforce, hnsw")
		return
	}
	
	name := args[0]
	algorithm := args[1]
	
	if algorithm != "bruteforce" && algorithm != "hnsw" {
		fmt.Printf("Invalid algorithm: %s. Use: bruteforce or hnsw\n", algorithm)
		return
	}
	
	reqData := map[string]interface{}{
		"name":      name,
		"algorithm": algorithm,
	}
	
	err := makePostRequest("/databases", reqData)
	if err != nil {
		fmt.Printf("Error creating database: %v\n", err)
		return
	}
	
	fmt.Printf("Database '%s' created successfully with algorithm '%s'\n", name, algorithm)
}

func handleUseDatabase(args []string) {
	if len(args) < 1 {
		fmt.Println("Usage: use-db <name>")
		fmt.Println("Example: use-db mydb")
		return
	}
	
	name := args[0]
	selectedDatabase = name
	fmt.Printf("Now using database: %s\n", name)
}

func handleListDatabases() {
	resp, err := http.Get(serverURL + "/databases")
	if err != nil {
		fmt.Printf("Error listing databases: %v\n", err)
		return
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}
	
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Printf("Error parsing response: %v\n", err)
		return
	}
	
	databases, ok := result["databases"].([]interface{})
	if !ok {
		fmt.Println("No databases found")
		return
	}
	
	fmt.Printf("Available databases (%d):\n", len(databases))
	for i, db := range databases {
		dbName := db.(string)
		fmt.Printf("%d. %s", i+1, dbName)
		if dbName == selectedDatabase {
			fmt.Print(" (selected)")
		}
		fmt.Println()
	}
}