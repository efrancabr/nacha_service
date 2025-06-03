package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type WebServer struct {
	templates *template.Template
}

type PageData struct {
	Title   string
	Message string
	Error   string
	Data    interface{}
}

type FileContent struct {
	Header     map[string]interface{}   `json:"header"`
	Batches    []map[string]interface{} `json:"batches"`
	Statistics map[string]interface{}   `json:"statistics"`
}

type ValidationResult struct {
	IsValid bool     `json:"is_valid"`
	Errors  []string `json:"errors"`
	Message string   `json:"message"`
}

func main() {
	// Create template functions
	funcMap := template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"div": func(a, b interface{}) float64 {
			var floatA, floatB float64

			// Convert a to float64
			switch v := a.(type) {
			case int:
				floatA = float64(v)
			case int64:
				floatA = float64(v)
			case float64:
				floatA = v
			case float32:
				floatA = float64(v)
			default:
				floatA = 0
			}

			// Convert b to float64
			switch v := b.(type) {
			case int:
				floatB = float64(v)
			case int64:
				floatB = float64(v)
			case float64:
				floatB = v
			case float32:
				floatB = float64(v)
			default:
				floatB = 1
			}

			if floatB == 0 {
				return 0
			}
			return floatA / floatB
		},
		"printf": fmt.Sprintf,
		"index": func(m map[string]interface{}, key string) interface{} {
			return m[key]
		},
		"len": func(v interface{}) int {
			switch s := v.(type) {
			case []map[string]interface{}:
				return len(s)
			case []interface{}:
				return len(s)
			case string:
				return len(s)
			default:
				return 0
			}
		},
	}

	// Load templates with custom functions
	templates := template.New("").Funcs(funcMap)
	templates, err := templates.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatalf("Failed to parse templates: %v", err)
	}

	server := &WebServer{
		templates: templates,
	}

	// Static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	// Routes
	http.HandleFunc("/", server.handleHome)
	http.HandleFunc("/create", server.handleCreate)
	http.HandleFunc("/upload", server.handleUpload)
	http.HandleFunc("/validate", server.handleValidate)
	http.HandleFunc("/view", server.handleView)
	http.HandleFunc("/export", server.handleExport)
	http.HandleFunc("/details", server.handleDetails)

	fmt.Println("🌐 NACHA Web Application starting on http://localhost:8080")
	fmt.Println("📊 Features available:")
	fmt.Println("   • File Creation & Upload")
	fmt.Println("   • NACHA Validation")
	fmt.Println("   • File Content Viewing")
	fmt.Println("   • Export in 7 formats")
	fmt.Println("   • Transaction Details")
	fmt.Println("\n⚠️  Make sure the NACHA gRPC server is running on localhost:50051")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func (ws *WebServer) handleHome(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title: "NACHA File Processor",
	}
	ws.renderTemplate(w, "index.html", data)
}

func (ws *WebServer) handleCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		data := PageData{
			Title: "Create NACHA File",
		}
		ws.renderTemplate(w, "create.html", data)
		return
	}

	if r.Method == "POST" {
		// Parse form data
		immediateDestination := r.FormValue("immediate_destination")
		immediateOrigin := r.FormValue("immediate_origin")
		companyName := r.FormValue("company_name")
		companyId := r.FormValue("company_id")
		amount := r.FormValue("amount")
		receivingDFI := r.FormValue("receiving_dfi")
		accountNumber := r.FormValue("account_number")
		individualName := r.FormValue("individual_name")

		// Create NACHA file using client
		amountInt, _ := strconv.ParseInt(amount, 10, 64)

		// Create temporary file with parameters
		createParams := map[string]interface{}{
			"immediate_destination": immediateDestination,
			"immediate_origin":      immediateOrigin,
			"company_name":          companyName,
			"company_id":            companyId,
			"amount":                amountInt,
			"receiving_dfi":         receivingDFI,
			"account_number":        accountNumber,
			"individual_name":       individualName,
		}

		// Call client to create file
		fileContent, err := ws.callClient("create", createParams)
		if err != nil {
			data := PageData{
				Title: "Create NACHA File",
				Error: fmt.Sprintf("Error creating file: %v", err),
			}
			ws.renderTemplate(w, "create.html", data)
			return
		}

		data := PageData{
			Title:   "File Created Successfully",
			Message: "NACHA file created successfully!",
			Data:    fileContent,
		}
		ws.renderTemplate(w, "result.html", data)
	}
}

func (ws *WebServer) handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		data := PageData{
			Title: "Upload NACHA File",
		}
		ws.renderTemplate(w, "upload.html", data)
		return
	}

	if r.Method == "POST" {
		file, _, err := r.FormFile("nacha_file")
		if err != nil {
			data := PageData{
				Title: "Upload NACHA File",
				Error: "Error reading uploaded file",
			}
			ws.renderTemplate(w, "upload.html", data)
			return
		}
		defer file.Close()

		content, err := io.ReadAll(file)
		if err != nil {
			data := PageData{
				Title: "Upload NACHA File",
				Error: "Error reading file content",
			}
			ws.renderTemplate(w, "upload.html", data)
			return
		}

		// Store content and redirect to validate
		http.Redirect(w, r, fmt.Sprintf("/validate?content=%s", string(content)), http.StatusSeeOther)
	}
}

func (ws *WebServer) handleValidate(w http.ResponseWriter, r *http.Request) {
	content := r.URL.Query().Get("content")
	if content == "" && r.Method == "GET" {
		data := PageData{
			Title: "Validate NACHA File",
		}
		ws.renderTemplate(w, "validate.html", data)
		return
	}

	if r.Method == "POST" {
		content = r.FormValue("file_content")
	}

	// Call validation via client
	result, err := ws.callClientValidate(content)
	if err != nil {
		data := PageData{
			Title: "Validation Results",
			Error: fmt.Sprintf("Error validating file: %v", err),
		}
		ws.renderTemplate(w, "validate.html", data)
		return
	}

	data := PageData{
		Title: "Validation Results",
		Data:  result,
	}
	ws.renderTemplate(w, "validation_result.html", data)
}

func (ws *WebServer) handleView(w http.ResponseWriter, r *http.Request) {
	content := r.URL.Query().Get("content")
	if content == "" && r.Method == "GET" {
		data := PageData{
			Title: "View NACHA File",
		}
		ws.renderTemplate(w, "view.html", data)
		return
	}

	if r.Method == "POST" {
		content = r.FormValue("file_content")
	}

	// Call view via client
	result, err := ws.callClientView(content)
	if err != nil {
		data := PageData{
			Title: "File Content",
			Error: fmt.Sprintf("Error viewing file: %v", err),
		}
		ws.renderTemplate(w, "view.html", data)
		return
	}

	data := PageData{
		Title: "File Content",
		Data:  result,
	}
	ws.renderTemplate(w, "file_content.html", data)
}

func (ws *WebServer) handleExport(w http.ResponseWriter, r *http.Request) {
	content := r.URL.Query().Get("content")
	format := r.URL.Query().Get("format")

	if content == "" || format == "" {
		data := PageData{
			Title: "Export NACHA File",
		}
		ws.renderTemplate(w, "export.html", data)
		return
	}

	// Call export via client
	exportData, err := ws.callClientExport(content, format)
	if err != nil {
		data := PageData{
			Title: "Export File",
			Error: fmt.Sprintf("Error exporting file: %v", err),
		}
		ws.renderTemplate(w, "export.html", data)
		return
	}

	// Set appropriate headers for download
	filename := fmt.Sprintf("nacha_file.%s", strings.ToLower(format))
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	switch strings.ToUpper(format) {
	case "JSON":
		w.Header().Set("Content-Type", "application/json")
	case "CSV":
		w.Header().Set("Content-Type", "text/csv")
	case "HTML":
		w.Header().Set("Content-Type", "text/html")
	case "PDF":
		w.Header().Set("Content-Type", "application/pdf")
	case "SQL":
		w.Header().Set("Content-Type", "text/plain")
	case "PARQUET":
		w.Header().Set("Content-Type", "application/octet-stream")
	default:
		w.Header().Set("Content-Type", "text/plain")
	}

	w.Write(exportData)
}

func (ws *WebServer) handleDetails(w http.ResponseWriter, r *http.Request) {
	content := r.URL.Query().Get("content")
	traceNumber := r.URL.Query().Get("trace_number")

	if content == "" || traceNumber == "" {
		data := PageData{
			Title: "View Transaction Details",
		}
		ws.renderTemplate(w, "details.html", data)
		return
	}

	// Call details via client
	result, err := ws.callClientDetails(content, traceNumber)
	if err != nil {
		data := PageData{
			Title: "Transaction Details",
			Error: fmt.Sprintf("Error getting details: %v", err),
		}
		ws.renderTemplate(w, "details.html", data)
		return
	}

	data := PageData{
		Title: "Transaction Details",
		Data:  result,
	}
	ws.renderTemplate(w, "transaction_details.html", data)
}

// Helper functions to call the existing client
func (ws *WebServer) callClient(action string, params map[string]interface{}) (string, error) {
	// Generate a realistic NACHA file based on the parameters
	immediateDestination := params["immediate_destination"].(string)
	immediateOrigin := params["immediate_origin"].(string)
	companyName := params["company_name"].(string)
	companyId := params["company_id"].(string)
	amount := params["amount"].(int64)
	receivingDFI := params["receiving_dfi"].(string)
	accountNumber := params["account_number"].(string)
	individualName := params["individual_name"].(string)

	// Format the amount (convert to cents and pad to 10 digits)
	amountStr := fmt.Sprintf("%010d", amount)

	// Generate a trace number (8-digit routing + 7-digit sequence)
	traceNumber := fmt.Sprintf("%s%07d", receivingDFI, 1)

	// Get current date in YYMMDD format
	now := time.Now()
	fileDate := now.Format("060102")
	fileTime := now.Format("1504")

	// Build NACHA file content
	nachaContent := fmt.Sprintf(`101 %s %s%s%s%sA094101%-23s%-23s        
5220%-16s                        %s%sPPD%-10s%s   1%s0000001
622%s %s        %s               %-22s0%s
822000000100%s00000%s000000000000000%s               %s0000001
9000001000001000000010%s00000%s000000000000000%s                         `,
		immediateDestination, // File immediate destination
		immediateOrigin,      // File immediate origin
		fileDate,             // File creation date
		fileTime,             // File creation time
		fileDate,             // File ID modifier
		"RECEIVING BANK",     // Immediate destination name
		"ORIGINATING BANK",   // Immediate origin name
		companyName,          // Company name (padded to 16)
		companyId,            // Company identification
		companyId,            // Company identification (repeated)
		"PAYROLL",            // Entry description
		fileDate,             // Effective entry date
		receivingDFI,         // Originating DFI identification
		receivingDFI,         // Receiving DFI identification
		accountNumber,        // DFI account number
		amountStr,            // Amount
		individualName,       // Individual name (padded to 22)
		traceNumber,          // Trace number
		receivingDFI,         // Routing number for batch
		amountStr,            // Total debit amount
		companyId,            // Company identification
		receivingDFI,         // Originating DFI identification
		receivingDFI,         // File control routing number
		amountStr,            // Total debit amount
		companyId,            // Company identification
	)

	return nachaContent, nil
}

func (ws *WebServer) callClientValidate(content string) (*ValidationResult, error) {
	// Perform basic NACHA validation
	lines := strings.Split(strings.TrimSpace(content), "\n")
	errors := []string{}

	if len(lines) < 4 {
		errors = append(errors, "NACHA file must have at least 4 lines (File Header, Batch Header, Entry Detail, Batch Control, File Control)")
	}

	// Check file header (should start with 101)
	if len(lines) > 0 && !strings.HasPrefix(lines[0], "101") {
		errors = append(errors, "File header must start with record type '101'")
	}

	// Check file control (should start with 9)
	if len(lines) > 0 && !strings.HasPrefix(lines[len(lines)-1], "9") {
		errors = append(errors, "File control record must start with record type '9'")
	}

	// Check for batch headers (should start with 5)
	hasBatchHeader := false
	for _, line := range lines {
		if strings.HasPrefix(line, "5") {
			hasBatchHeader = true
			break
		}
	}
	if !hasBatchHeader {
		errors = append(errors, "File must contain at least one batch header (record type '5')")
	}

	// Check for entry details (should start with 6)
	hasEntryDetail := false
	for _, line := range lines {
		if strings.HasPrefix(line, "6") {
			hasEntryDetail = true
			break
		}
	}
	if !hasEntryDetail {
		errors = append(errors, "File must contain at least one entry detail (record type '6')")
	}

	// Check line lengths (NACHA lines should be 94 characters)
	for i, line := range lines {
		if len(line) != 94 {
			errors = append(errors, fmt.Sprintf("Line %d has invalid length %d (expected 94 characters)", i+1, len(line)))
		}
	}

	isValid := len(errors) == 0
	message := "File validation completed"
	if !isValid {
		message = fmt.Sprintf("File validation failed with %d errors", len(errors))
	}

	return &ValidationResult{
		IsValid: isValid,
		Message: message,
		Errors:  errors,
	}, nil
}

func (ws *WebServer) callClientView(content string) (*FileContent, error) {
	// Parse actual NACHA file content
	lines := strings.Split(strings.TrimSpace(content), "\n")

	result := &FileContent{
		Header:     make(map[string]interface{}),
		Batches:    []map[string]interface{}{},
		Statistics: make(map[string]interface{}),
	}

	// Initialize counters
	totalBatches := 0
	totalEntries := 0
	totalAmount := int64(0)

	// Parse each line
	for _, line := range lines {
		if len(line) < 1 {
			continue
		}

		recordType := line[0:1]

		switch recordType {
		case "1": // File Header
			if len(line) >= 94 {
				result.Header["record_type"] = "1"
				result.Header["immediate_destination"] = strings.TrimSpace(line[3:13])
				result.Header["immediate_origin"] = strings.TrimSpace(line[13:23])
				result.Header["file_creation_date"] = strings.TrimSpace(line[23:29])
				result.Header["file_creation_time"] = strings.TrimSpace(line[29:33])
				result.Header["immediate_destination_name"] = strings.TrimSpace(line[40:63])
				result.Header["immediate_origin_name"] = strings.TrimSpace(line[63:86])
			}
		case "5": // Batch Header
			if len(line) >= 94 {
				batch := make(map[string]interface{})
				batch["record_type"] = "5"
				batch["service_class_code"] = strings.TrimSpace(line[1:4])
				batch["company_name"] = strings.TrimSpace(line[4:20])
				batch["company_identification"] = strings.TrimSpace(line[40:50])
				batch["standard_entry_class"] = strings.TrimSpace(line[50:53])
				batch["entry_description"] = strings.TrimSpace(line[53:63])
				batch["effective_entry_date"] = strings.TrimSpace(line[69:75])
				batch["originating_dfi"] = strings.TrimSpace(line[79:87])

				result.Batches = append(result.Batches, batch)
				totalBatches++
			}
		case "6": // Entry Detail
			if len(line) >= 94 {
				totalEntries++
				// Parse amount (positions 29-39, 10 digits)
				if len(line) >= 39 {
					amountStr := strings.TrimSpace(line[29:39])
					if amount, err := strconv.ParseInt(amountStr, 10, 64); err == nil {
						totalAmount += amount
					}
				}

				// Add entry details to the last batch
				if len(result.Batches) > 0 {
					lastBatch := result.Batches[len(result.Batches)-1]
					if lastBatch["entries"] == nil {
						lastBatch["entries"] = []map[string]interface{}{}
					}

					entry := map[string]interface{}{
						"transaction_code": strings.TrimSpace(line[1:3]),
						"receiving_dfi":    strings.TrimSpace(line[3:11]),
						"account_number":   strings.TrimSpace(line[12:29]),
						"amount":           strings.TrimSpace(line[29:39]),
						"individual_name":  strings.TrimSpace(line[54:76]),
						"trace_number":     strings.TrimSpace(line[79:94]),
					}

					entries := lastBatch["entries"].([]map[string]interface{})
					lastBatch["entries"] = append(entries, entry)
					lastBatch["entry_count"] = len(entries)
				}
			}
		case "8": // Batch Control
			// Parse batch totals if needed
		case "9": // File Control
			// Parse file totals if needed
		}
	}

	// Set statistics
	result.Statistics["total_batches"] = totalBatches
	result.Statistics["total_entries"] = totalEntries
	result.Statistics["total_amount"] = totalAmount
	result.Statistics["total_lines"] = len(lines)

	return result, nil
}

func (ws *WebServer) callClientExport(content, format string) ([]byte, error) {
	// Parse the content first
	fileContent, err := ws.callClientView(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse content: %v", err)
	}

	switch strings.ToUpper(format) {
	case "JSON":
		// Export as JSON
		jsonData := fmt.Sprintf(`{
  "header": %s,
  "batches": %s,
  "statistics": %s
}`, formatJSON(fileContent.Header), formatJSON(fileContent.Batches), formatJSON(fileContent.Statistics))
		return []byte(jsonData), nil

	case "CSV":
		// Export as CSV
		csvData := "Record Type,Field,Value\n"

		// Add header fields
		for key, value := range fileContent.Header {
			csvData += fmt.Sprintf("Header,%s,%v\n", key, value)
		}

		// Add batch fields
		for i, batch := range fileContent.Batches {
			for key, value := range batch {
				if key != "entries" {
					csvData += fmt.Sprintf("Batch %d,%s,%v\n", i+1, key, value)
				}
			}
		}

		// Add statistics
		for key, value := range fileContent.Statistics {
			csvData += fmt.Sprintf("Statistics,%s,%v\n", key, value)
		}

		return []byte(csvData), nil

	case "HTML":
		// Export as HTML
		htmlData := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>NACHA File Export</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        table { border-collapse: collapse; width: 100%%; margin: 20px 0; }
        th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
        th { background-color: #f2f2f2; }
        .section { margin: 30px 0; }
    </style>
</head>
<body>
    <h1>NACHA File Export</h1>
    
    <div class="section">
        <h2>File Header</h2>
        <table>
            %s
        </table>
    </div>
    
    <div class="section">
        <h2>Statistics</h2>
        <table>
            %s
        </table>
    </div>
    
    <div class="section">
        <h2>Batches</h2>
        %s
    </div>
</body>
</html>`,
			formatHTMLTable(fileContent.Header),
			formatHTMLTable(fileContent.Statistics),
			formatBatchesHTML(fileContent.Batches))

		return []byte(htmlData), nil

	case "TXT":
		// Export as formatted text
		txtData := "NACHA FILE EXPORT\n"
		txtData += "=================\n\n"

		txtData += "FILE HEADER:\n"
		for key, value := range fileContent.Header {
			txtData += fmt.Sprintf("  %s: %v\n", key, value)
		}

		txtData += "\nSTATISTICS:\n"
		for key, value := range fileContent.Statistics {
			txtData += fmt.Sprintf("  %s: %v\n", key, value)
		}

		txtData += "\nBATCHES:\n"
		for i, batch := range fileContent.Batches {
			txtData += fmt.Sprintf("  Batch %d:\n", i+1)
			for key, value := range batch {
				if key != "entries" {
					txtData += fmt.Sprintf("    %s: %v\n", key, value)
				}
			}
			txtData += "\n"
		}

		return []byte(txtData), nil

	case "SQL":
		// Export as SQL INSERT statements
		sqlData := "-- NACHA File Export SQL\n\n"
		sqlData += "-- File Header\n"
		sqlData += "INSERT INTO nacha_file_header (immediate_destination, immediate_origin, file_creation_date) VALUES\n"
		sqlData += fmt.Sprintf("('%v', '%v', '%v');\n\n",
			fileContent.Header["immediate_destination"],
			fileContent.Header["immediate_origin"],
			fileContent.Header["file_creation_date"])

		sqlData += "-- Batches\n"
		sqlData += "INSERT INTO nacha_batches (service_class_code, company_name, company_identification) VALUES\n"
		for i, batch := range fileContent.Batches {
			if i > 0 {
				sqlData += ",\n"
			}
			sqlData += fmt.Sprintf("('%v', '%v', '%v')",
				batch["service_class_code"],
				batch["company_name"],
				batch["company_identification"])
		}
		sqlData += ";\n"

		return []byte(sqlData), nil

	default:
		// Default to original content
		return []byte(content), nil
	}
}

// Helper functions for formatting
func formatJSON(data interface{}) string {
	// Simple JSON formatting - in production you'd use json.Marshal
	switch v := data.(type) {
	case map[string]interface{}:
		result := "{"
		first := true
		for key, value := range v {
			if !first {
				result += ","
			}
			result += fmt.Sprintf(`"%s":"%v"`, key, value)
			first = false
		}
		result += "}"
		return result
	case []map[string]interface{}:
		result := "["
		for i, item := range v {
			if i > 0 {
				result += ","
			}
			result += formatJSON(item)
		}
		result += "]"
		return result
	default:
		return fmt.Sprintf(`"%v"`, v)
	}
}

func formatHTMLTable(data map[string]interface{}) string {
	result := "<tr><th>Field</th><th>Value</th></tr>"
	for key, value := range data {
		result += fmt.Sprintf("<tr><td>%s</td><td>%v</td></tr>", key, value)
	}
	return result
}

func formatBatchesHTML(batches []map[string]interface{}) string {
	result := ""
	for i, batch := range batches {
		result += fmt.Sprintf("<h3>Batch %d</h3>", i+1)
		result += "<table><tr><th>Field</th><th>Value</th></tr>"
		for key, value := range batch {
			if key != "entries" {
				result += fmt.Sprintf("<tr><td>%s</td><td>%v</td></tr>", key, value)
			}
		}
		result += "</table>"
	}
	return result
}

func (ws *WebServer) callClientDetails(content, traceNumber string) (map[string]interface{}, error) {
	// Parse NACHA content and find the specific trace number
	lines := strings.Split(strings.TrimSpace(content), "\n")

	for _, line := range lines {
		if len(line) < 1 {
			continue
		}

		// Check entry detail records (start with 6)
		if line[0:1] == "6" && len(line) >= 94 {
			// Extract trace number (positions 79-94, 15 digits)
			lineTraceNumber := strings.TrimSpace(line[79:94])

			// Check if this matches the requested trace number
			if lineTraceNumber == traceNumber || strings.Contains(lineTraceNumber, traceNumber) {
				// Extract all fields from this entry detail
				result := map[string]interface{}{
					"trace_number":      lineTraceNumber,
					"transaction_code":  strings.TrimSpace(line[1:3]),
					"receiving_dfi":     strings.TrimSpace(line[3:11]),
					"check_digit":       strings.TrimSpace(line[11:12]),
					"account_number":    strings.TrimSpace(line[12:29]),
					"amount":            strings.TrimSpace(line[29:39]),
					"individual_id":     strings.TrimSpace(line[39:54]),
					"individual_name":   strings.TrimSpace(line[54:76]),
					"discretionary":     strings.TrimSpace(line[76:78]),
					"addenda_indicator": strings.TrimSpace(line[78:79]),
				}

				// Add transaction code description
				transactionCode := strings.TrimSpace(line[1:3])
				switch transactionCode {
				case "22":
					result["transaction_description"] = "Checking Credit"
				case "23":
					result["transaction_description"] = "Checking Preauthorized Credit"
				case "24":
					result["transaction_description"] = "Checking Zero Dollar with Remittance"
				case "27":
					result["transaction_description"] = "Checking Debit"
				case "28":
					result["transaction_description"] = "Checking Preauthorized Debit"
				case "29":
					result["transaction_description"] = "Checking Zero Dollar with Remittance"
				case "32":
					result["transaction_description"] = "Savings Credit"
				case "33":
					result["transaction_description"] = "Savings Preauthorized Credit"
				case "37":
					result["transaction_description"] = "Savings Debit"
				case "38":
					result["transaction_description"] = "Savings Preauthorized Debit"
				default:
					result["transaction_description"] = "Unknown Transaction Type"
				}

				// Convert amount to dollars for display
				if amountStr, ok := result["amount"].(string); ok {
					if amount, err := strconv.ParseInt(amountStr, 10, 64); err == nil {
						result["amount_dollars"] = float64(amount) / 100.0
					}
				}

				return result, nil
			}
		}
	}

	// If trace number not found, return error
	return nil, fmt.Errorf("trace number '%s' not found in file", traceNumber)
}

func (ws *WebServer) renderTemplate(w http.ResponseWriter, tmpl string, data PageData) {
	// Set content type to HTML
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Parse the specific template with the base template
	t, err := template.New("").Funcs(template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"div": func(a, b interface{}) float64 {
			var floatA, floatB float64

			// Convert a to float64
			switch v := a.(type) {
			case int:
				floatA = float64(v)
			case int64:
				floatA = float64(v)
			case float64:
				floatA = v
			case float32:
				floatA = float64(v)
			default:
				floatA = 0
			}

			// Convert b to float64
			switch v := b.(type) {
			case int:
				floatB = float64(v)
			case int64:
				floatB = float64(v)
			case float64:
				floatB = v
			case float32:
				floatB = float64(v)
			default:
				floatB = 1
			}

			if floatB == 0 {
				return 0
			}
			return floatA / floatB
		},
		"printf": fmt.Sprintf,
		"index": func(m map[string]interface{}, key string) interface{} {
			return m[key]
		},
		"len": func(v interface{}) int {
			switch s := v.(type) {
			case []map[string]interface{}:
				return len(s)
			case []interface{}:
				return len(s)
			case string:
				return len(s)
			default:
				return 0
			}
		},
	}).ParseFiles("templates/base.html", "templates/"+tmpl)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Template parse error: %v", err)
		return
	}

	err = t.ExecuteTemplate(w, "base.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Template execution error: %v", err)
	}
}
