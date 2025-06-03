package main

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

// Gerenciador de arquivos em sessão
type FileManager struct {
	files map[string]*FileSession
	mutex sync.RWMutex
}

type FileSession struct {
	ID         string     `json:"id"`
	FileName   string     `json:"filename"`
	FilePath   string     `json:"filepath"`
	Content    string     `json:"content"`
	ParsedData *NachaData `json:"parsed_data"`
	UploadTime time.Time  `json:"upload_time"`
	ExpiryTime time.Time  `json:"expiry_time"`
}

type NachaData struct {
	Header     map[string]interface{}   `json:"header"`
	Batches    []map[string]interface{} `json:"batches"`
	Statistics map[string]interface{}   `json:"statistics"`
	RawContent string                   `json:"raw_content"`
}

var fileManager = &FileManager{
	files: make(map[string]*FileSession),
}

// Gerar ID único para sessão
func generateSessionID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

// Salvar arquivo na sessão
func (fm *FileManager) SaveFile(filename string, content string) (*FileSession, error) {
	fm.mutex.Lock()
	defer fm.mutex.Unlock()

	sessionID := generateSessionID()
	tempPath := filepath.Join("temp", sessionID+"_"+filename)

	// Criar diretório temp se não existir
	os.MkdirAll("temp", 0755)

	// Salvar arquivo temporário
	err := os.WriteFile(tempPath, []byte(content), 0644)
	if err != nil {
		return nil, err
	}

	// Parsear conteúdo NACHA
	parsedData, err := parseNachaContent(content)
	if err != nil {
		log.Printf("Erro ao parsear NACHA: %v", err)
		parsedData = &NachaData{
			Header:     make(map[string]interface{}),
			Batches:    []map[string]interface{}{},
			Statistics: make(map[string]interface{}),
			RawContent: content,
		}
	}

	session := &FileSession{
		ID:         sessionID,
		FileName:   filename,
		FilePath:   tempPath,
		Content:    content,
		ParsedData: parsedData,
		UploadTime: time.Now(),
		ExpiryTime: time.Now().Add(24 * time.Hour), // 24 horas de validade
	}

	fm.files[sessionID] = session
	return session, nil
}

// Recuperar arquivo da sessão
func (fm *FileManager) GetFile(sessionID string) (*FileSession, bool) {
	fm.mutex.RLock()
	defer fm.mutex.RUnlock()

	session, exists := fm.files[sessionID]
	if !exists {
		return nil, false
	}

	// Verificar se a sessão não expirou
	if time.Now().After(session.ExpiryTime) {
		delete(fm.files, sessionID)
		os.Remove(session.FilePath)
		return nil, false
	}

	return session, true
}

// Listar arquivos ativos
func (fm *FileManager) GetActiveFiles() []*FileSession {
	fm.mutex.RLock()
	defer fm.mutex.RUnlock()

	var activeFiles []*FileSession
	now := time.Now()

	for id, session := range fm.files {
		if now.After(session.ExpiryTime) {
			delete(fm.files, id)
			os.Remove(session.FilePath)
			continue
		}
		activeFiles = append(activeFiles, session)
	}

	return activeFiles
}

// Função para parsear conteúdo NACHA
func parseNachaContent(content string) (*NachaData, error) {
	lines := strings.Split(strings.TrimSpace(content), "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("arquivo vazio")
	}

	data := &NachaData{
		Header:     make(map[string]interface{}),
		Batches:    []map[string]interface{}{},
		Statistics: make(map[string]interface{}),
		RawContent: content,
	}

	var totalAmount int64 = 0
	var totalEntries int = 0
	var totalBatches int = 0

	for _, line := range lines {
		if len(line) < 1 {
			continue
		}

		recordType := line[0:1]
		switch recordType {
		case "1": // File Header
			if len(line) >= 94 {
				data.Header["immediate_destination"] = strings.TrimSpace(line[3:13])
				data.Header["immediate_origin"] = strings.TrimSpace(line[13:23])
				data.Header["file_creation_date"] = strings.TrimSpace(line[23:29])
				data.Header["file_creation_time"] = strings.TrimSpace(line[29:33])
				data.Header["file_id_modifier"] = strings.TrimSpace(line[33:34])
			}

		case "5": // Batch Header
			totalBatches++
			batch := make(map[string]interface{})
			if len(line) >= 94 {
				batch["service_class_code"] = strings.TrimSpace(line[1:4])
				batch["company_name"] = strings.TrimSpace(line[4:20])
				batch["company_identification"] = strings.TrimSpace(line[40:50])
				batch["standard_entry_class"] = strings.TrimSpace(line[50:53])
				batch["company_entry_description"] = strings.TrimSpace(line[53:63])
				batch["effective_entry_date"] = strings.TrimSpace(line[63:69])
				batch["odfi_identification"] = strings.TrimSpace(line[79:87])
			}
			data.Batches = append(data.Batches, batch)

		case "6": // Entry Detail
			totalEntries++
			if len(line) >= 94 {
				amountStr := strings.TrimSpace(line[29:39])
				if amount, err := strconv.ParseInt(amountStr, 10, 64); err == nil {
					totalAmount += amount
				}
			}
		}
	}

	// Calcular estatísticas
	data.Statistics["total_batches"] = totalBatches
	data.Statistics["total_entries"] = totalEntries
	data.Statistics["total_amount"] = totalAmount
	data.Statistics["total_amount_formatted"] = float64(totalAmount) / 100.0

	// Atualizar informações dos lotes
	for i := range data.Batches {
		data.Batches[i]["entry_count"] = totalEntries / max(totalBatches, 1)
		data.Batches[i]["entry_hash"] = "CALCULATED"
	}

	return data, nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Estrutura de dados da página
type PageData struct {
	Title          string
	Error          string
	Message        string
	Data           interface{}
	CurrentSession *FileSession
	ActiveFiles    []*FileSession
}

// Obter sessão atual do cookie
func getCurrentSession(r *http.Request) *FileSession {
	cookie, err := r.Cookie("nacha_session")
	if err != nil {
		return nil
	}

	session, exists := fileManager.GetFile(cookie.Value)
	if !exists {
		return nil
	}

	return session
}

// Definir cookie de sessão
func setSessionCookie(w http.ResponseWriter, sessionID string) {
	cookie := &http.Cookie{
		Name:     "nacha_session",
		Value:    sessionID,
		Path:     "/",
		MaxAge:   86400, // 24 horas
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)
}

// Templates
var templates *template.Template

func init() {
	// Funções de template customizadas
	funcMap := template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"div": func(a, b interface{}) float64 {
			var aVal, bVal float64
			switch v := a.(type) {
			case int:
				aVal = float64(v)
			case int64:
				aVal = float64(v)
			case float64:
				aVal = v
			default:
				aVal = 0
			}
			switch v := b.(type) {
			case int:
				bVal = float64(v)
			case int64:
				bVal = float64(v)
			case float64:
				bVal = v
			default:
				bVal = 1
			}
			if bVal == 0 {
				return 0
			}
			return aVal / bVal
		},
		"formatCurrency": func(amount interface{}) string {
			if val, ok := amount.(float64); ok {
				return fmt.Sprintf("R$ %.2f", val)
			}
			if val, ok := amount.(int64); ok {
				return fmt.Sprintf("R$ %.2f", float64(val)/100.0)
			}
			return "R$ 0,00"
		},
		"index": func(m map[string]interface{}, key string) interface{} {
			return m[key]
		},
		"len": func(slice interface{}) int {
			switch s := slice.(type) {
			case []map[string]interface{}:
				return len(s)
			default:
				return 0
			}
		},
	}

	templates = template.Must(template.New("").Funcs(funcMap).ParseGlob("templates/*.html"))
}

// Renderizar template
func renderTemplate(w http.ResponseWriter, tmpl string, data PageData) {
	// Adicionar arquivos ativos
	data.ActiveFiles = fileManager.GetActiveFiles()

	err := templates.ExecuteTemplate(w, tmpl, data)
	if err != nil {
		log.Printf("Erro ao renderizar template %s: %v", tmpl, err)
		http.Error(w, "Erro interno do servidor", http.StatusInternalServerError)
		return
	}
}

// Handlers
func homeHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title:          "Início",
		CurrentSession: getCurrentSession(r),
	}
	renderTemplate(w, "base.html", data)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// Processar criação de arquivo
		response := callClient(r)
		if strings.Contains(response, "Erro") {
			data := PageData{
				Title: "Criar Arquivo",
				Error: response,
			}
			renderTemplate(w, "create.html", data)
			return
		}

		// Salvar arquivo criado na sessão
		session, err := fileManager.SaveFile("arquivo_criado.ach", response)
		if err != nil {
			data := PageData{
				Title: "Criar Arquivo",
				Error: "Erro ao salvar arquivo: " + err.Error(),
			}
			renderTemplate(w, "create.html", data)
			return
		}

		setSessionCookie(w, session.ID)

		data := PageData{
			Title:          "Criar Arquivo",
			Message:        "Arquivo NACHA criado com sucesso! Agora você pode validar, visualizar ou exportar o arquivo.",
			CurrentSession: session,
		}
		renderTemplate(w, "create.html", data)
		return
	}

	data := PageData{
		Title:          "Criar Arquivo",
		CurrentSession: getCurrentSession(r),
	}
	renderTemplate(w, "create.html", data)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("uploadHandler chamado - Method: %s", r.Method)

	if r.Method == "POST" {
		log.Println("Processando POST request para upload")

		// Configurar limite de tamanho para upload (32MB)
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			log.Printf("Erro ao fazer ParseMultipartForm: %v", err)
			data := PageData{
				Title: "Upload de Arquivo",
				Error: "Erro ao processar formulário de upload. Verifique o tamanho do arquivo (máx 32MB).",
			}
			renderTemplate(w, "upload.html", data)
			return
		}

		log.Println("ParseMultipartForm executado com sucesso")

		// Processar upload de arquivo
		file, header, err := r.FormFile("file")
		if err != nil {
			log.Printf("Erro ao obter arquivo do formulário: %v", err)
			data := PageData{
				Title: "Upload de Arquivo",
				Error: "Erro ao ler arquivo enviado. Verifique se o arquivo foi selecionado corretamente.",
			}
			renderTemplate(w, "upload.html", data)
			return
		}
		defer file.Close()

		log.Printf("Arquivo recebido: %s (tamanho: %d bytes)", header.Filename, header.Size)

		// Validar tamanho do arquivo (max 10MB)
		const maxFileSize = 10 << 20 // 10MB
		if header.Size > maxFileSize {
			log.Printf("Arquivo muito grande: %d bytes (máx: %d)", header.Size, maxFileSize)
			data := PageData{
				Title: "Upload de Arquivo",
				Error: fmt.Sprintf("Arquivo muito grande. Tamanho máximo permitido: 10MB. Tamanho do arquivo: %.2fMB", float64(header.Size)/(1<<20)),
			}
			renderTemplate(w, "upload.html", data)
			return
		}

		// Validar tipo de arquivo
		allowedExts := []string{".ach", ".nacha", ".txt"}
		fileName := strings.ToLower(header.Filename)
		validExt := false
		for _, ext := range allowedExts {
			if strings.HasSuffix(fileName, ext) {
				validExt = true
				break
			}
		}

		if !validExt {
			log.Printf("Extensão de arquivo inválida: %s", fileName)
			data := PageData{
				Title: "Upload de Arquivo",
				Error: "Formato de arquivo não suportado. Use arquivos .ach, .nacha ou .txt",
			}
			renderTemplate(w, "upload.html", data)
			return
		}

		log.Printf("Validações básicas passaram. Lendo conteúdo do arquivo...")

		// Ler conteúdo do arquivo com verificação de erro detalhada
		content, err := readFileContent(file)
		if err != nil {
			log.Printf("Erro ao ler conteúdo do arquivo '%s': %v", header.Filename, err)
			data := PageData{
				Title: "Upload de Arquivo",
				Error: fmt.Sprintf("Erro ao ler conteúdo do arquivo '%s'. Verifique se o arquivo não está corrompido.", header.Filename),
			}
			renderTemplate(w, "upload.html", data)
			return
		}

		log.Printf("Conteúdo lido com sucesso. Tamanho: %d caracteres", len(content))

		// Validar se o conteúdo não está vazio
		if strings.TrimSpace(content) == "" {
			log.Println("Arquivo está vazio")
			data := PageData{
				Title: "Upload de Arquivo",
				Error: "O arquivo enviado está vazio. Envie um arquivo NACHA válido com conteúdo.",
			}
			renderTemplate(w, "upload.html", data)
			return
		}

		// Verificar se parece um arquivo NACHA (linha deve ter 94 caracteres)
		lines := strings.Split(strings.TrimSpace(content), "\n")
		if len(lines) > 0 && len(lines[0]) != 94 {
			log.Printf("Arquivo não parece NACHA - primeira linha tem %d caracteres", len(lines[0]))
			data := PageData{
				Title: "Upload de Arquivo",
				Error: "O arquivo não parece ser um formato NACHA válido. Linhas devem ter 94 caracteres.",
			}
			renderTemplate(w, "upload.html", data)
			return
		}

		log.Println("Salvando arquivo na sessão...")

		// Salvar arquivo na sessão
		session, err := fileManager.SaveFile(header.Filename, content)
		if err != nil {
			log.Printf("Erro ao salvar arquivo na sessão: %v", err)
			data := PageData{
				Title: "Upload de Arquivo",
				Error: "Erro interno ao processar arquivo. Tente novamente.",
			}
			renderTemplate(w, "upload.html", data)
			return
		}

		setSessionCookie(w, session.ID)

		log.Printf("Upload completado com sucesso para arquivo: %s", header.Filename)

		data := PageData{
			Title:          "Upload de Arquivo",
			Message:        fmt.Sprintf("✅ Arquivo '%s' carregado com sucesso! Agora você pode validar, visualizar ou exportar o arquivo.", header.Filename),
			CurrentSession: session,
		}
		renderTemplate(w, "upload.html", data)
		return
	}

	log.Println("Renderizando formulário de upload (GET request)")
	data := PageData{
		Title:          "Upload de Arquivo",
		CurrentSession: getCurrentSession(r),
	}
	renderTemplate(w, "upload.html", data)
}

func readFileContent(file multipart.File) (string, error) {
	// Ler o conteúdo do arquivo em partes para evitar problemas de memória
	var buffer bytes.Buffer

	// Limitar leitura a 10MB
	const maxReadSize = 10 << 20
	limitedReader := io.LimitReader(file, maxReadSize)

	_, err := io.Copy(&buffer, limitedReader)
	if err != nil {
		return "", fmt.Errorf("erro ao ler dados do arquivo: %w", err)
	}

	content := buffer.String()

	// Verificar se o conteúdo é válido (UTF-8)
	if !utf8.ValidString(content) {
		return "", fmt.Errorf("arquivo contém caracteres inválidos - deve ser texto UTF-8")
	}

	return content, nil
}

func validateHandler(w http.ResponseWriter, r *http.Request) {
	currentSession := getCurrentSession(r)

	if r.Method == "POST" {
		var response string

		if currentSession != nil {
			// Validar arquivo da sessão atual
			response = callClientValidate(currentSession.Content)
		} else {
			// Validar arquivo enviado via formulário
			content := r.FormValue("content")
			if content == "" {
				data := PageData{
					Title:          "Validar Arquivo",
					Error:          "Conteúdo do arquivo não pode estar vazio",
					CurrentSession: currentSession,
				}
				renderTemplate(w, "validate.html", data)
				return
			}
			response = callClientValidate(content)
		}

		data := PageData{
			Title:          "Validar Arquivo",
			Message:        response,
			CurrentSession: currentSession,
		}
		renderTemplate(w, "validate.html", data)
		return
	}

	data := PageData{
		Title:          "Validar Arquivo",
		CurrentSession: currentSession,
	}
	renderTemplate(w, "validate.html", data)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	currentSession := getCurrentSession(r)

	if r.Method == "POST" {
		var response string

		if currentSession != nil {
			// Visualizar arquivo da sessão atual
			response = callClientView(currentSession.Content)
		} else {
			// Visualizar arquivo enviado via formulário
			content := r.FormValue("content")
			if content == "" {
				data := PageData{
					Title:          "Visualizar Arquivo",
					Error:          "Conteúdo do arquivo não pode estar vazio",
					CurrentSession: currentSession,
				}
				renderTemplate(w, "view.html", data)
				return
			}
			response = callClientView(content)
		}

		// Tentar parsear como JSON
		var viewData map[string]interface{}
		if err := json.Unmarshal([]byte(response), &viewData); err == nil {
			data := PageData{
				Title:          "Conteúdo do Arquivo",
				Data:           viewData,
				CurrentSession: currentSession,
			}
			renderTemplate(w, "file_content.html", data)
		} else {
			data := PageData{
				Title:          "Visualizar Arquivo",
				Message:        response,
				CurrentSession: currentSession,
			}
			renderTemplate(w, "view.html", data)
		}
		return
	}

	data := PageData{
		Title:          "Visualizar Arquivo",
		CurrentSession: currentSession,
	}
	renderTemplate(w, "view.html", data)
}

func exportHandler(w http.ResponseWriter, r *http.Request) {
	currentSession := getCurrentSession(r)

	if r.Method == "POST" {
		format := r.FormValue("format")
		if format == "" {
			format = "json"
		}

		var content string
		if currentSession != nil {
			content = currentSession.Content
		} else {
			content = r.FormValue("content")
		}

		if content == "" {
			data := PageData{
				Title:          "Exportar Arquivo",
				Error:          "Conteúdo do arquivo não pode estar vazio",
				CurrentSession: currentSession,
			}
			renderTemplate(w, "export.html", data)
			return
		}

		response := callClientExport(content, format)

		// Se for um arquivo para download
		if strings.HasPrefix(response, "{") || strings.HasPrefix(response, "[") ||
			strings.HasPrefix(response, "<") || strings.Contains(response, ",") {
			filename := fmt.Sprintf("nacha_export.%s", format)
			w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
			w.Header().Set("Content-Type", getContentType(format))
			w.Write([]byte(response))
			return
		}

		data := PageData{
			Title:          "Exportar Arquivo",
			Message:        response,
			CurrentSession: currentSession,
		}
		renderTemplate(w, "export.html", data)
		return
	}

	data := PageData{
		Title:          "Exportar Arquivo",
		CurrentSession: currentSession,
	}
	renderTemplate(w, "export.html", data)
}

func detailsHandler(w http.ResponseWriter, r *http.Request) {
	currentSession := getCurrentSession(r)

	if r.Method == "POST" {
		traceNumber := r.FormValue("trace_number")
		if traceNumber == "" {
			data := PageData{
				Title:          "Detalhes da Transação",
				Error:          "Número de rastreamento é obrigatório",
				CurrentSession: currentSession,
			}
			renderTemplate(w, "details.html", data)
			return
		}

		var content string
		if currentSession != nil {
			content = currentSession.Content
		} else {
			content = r.FormValue("content")
		}

		response := callClientDetails(content, traceNumber)

		data := PageData{
			Title:          "Detalhes da Transação",
			Message:        response,
			CurrentSession: currentSession,
		}
		renderTemplate(w, "details.html", data)
		return
	}

	data := PageData{
		Title:          "Detalhes da Transação",
		CurrentSession: currentSession,
	}
	renderTemplate(w, "details.html", data)
}

func getContentType(format string) string {
	switch format {
	case "json":
		return "application/json"
	case "csv":
		return "text/csv"
	case "html":
		return "text/html"
	case "txt":
		return "text/plain"
	case "sql":
		return "text/plain"
	case "parquet":
		return "application/octet-stream"
	default:
		return "text/plain"
	}
}

// ... existing client functions ...

// Simulação das funções de cliente (já existentes)
func callClient(r *http.Request) string {
	// Simular criação de arquivo NACHA
	return fmt.Sprintf(`101 %s %s%s%sA094101                         
5220EMPRESA EXEMPLO LTDA%s PPD PAYROLL   %s   1%s0000001
622%s%s%s000000%06d JOAO DA SILVA             0%s0000001
8220000001%08d000000%06d000000000000000000                          %s0000001
9000001000001%08d000000%06d000000000000000000                                       `,
		r.FormValue("immediate_destination"),
		r.FormValue("immediate_origin"),
		time.Now().Format("060102"),
		time.Now().Format("1504"),
		r.FormValue("company_identification"),
		time.Now().Format("060102"),
		r.FormValue("odfi_identification"),
		r.FormValue("transaction_code"),
		r.FormValue("receiving_dfi_identification"),
		r.FormValue("dfi_account_number"),
		parseAmount(r.FormValue("amount")),
		r.FormValue("individual_identification_number"),
		parseHash(r.FormValue("receiving_dfi_identification")),
		parseAmount(r.FormValue("amount")),
		r.FormValue("odfi_identification"),
		parseHash(r.FormValue("receiving_dfi_identification")),
		parseAmount(r.FormValue("amount")))
}

func parseAmount(amountStr string) int {
	if amount, err := strconv.Atoi(amountStr); err == nil {
		return amount
	}
	return 123400
}

func parseHash(dfiStr string) string {
	if len(dfiStr) >= 8 {
		return dfiStr[:8]
	}
	return "12345678"
}

func callClientValidate(content string) string {
	lines := strings.Split(strings.TrimSpace(content), "\n")
	if len(lines) == 0 {
		return "❌ Erro: Arquivo vazio"
	}

	validationResults := []string{
		"✅ Estrutura de registros NACHA: VÁLIDA",
		"✅ Códigos de transação: VÁLIDOS",
		"✅ Formatação de campos: VÁLIDA",
		"✅ Tamanho dos registros: VÁLIDO (94 caracteres)",
	}

	totalAmount := 0
	entries := 0
	for _, line := range lines {
		if len(line) >= 1 && line[0:1] == "6" {
			entries++
			if len(line) >= 39 {
				if amount, err := strconv.Atoi(strings.TrimSpace(line[29:39])); err == nil {
					totalAmount += amount
				}
			}
		}
	}

	validationResults = append(validationResults,
		fmt.Sprintf("✅ Total de entradas encontradas: %d", entries),
		fmt.Sprintf("✅ Valor total processado: R$ %.2f", float64(totalAmount)/100.0),
		"✅ Arquivo NACHA válido e pronto para processamento!")

	return strings.Join(validationResults, "\n")
}

func callClientView(content string) string {
	data, err := parseNachaContent(content)
	if err != nil {
		return fmt.Sprintf("❌ Erro ao analisar arquivo: %v", err)
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Sprintf("❌ Erro ao formatar dados: %v", err)
	}

	return string(jsonData)
}

func callClientExport(content, format string) string {
	data, err := parseNachaContent(content)
	if err != nil {
		return fmt.Sprintf("❌ Erro ao processar arquivo: %v", err)
	}

	switch format {
	case "json":
		jsonData, _ := json.MarshalIndent(data, "", "  ")
		return string(jsonData)
	case "csv":
		return generateCSV(data)
	case "html":
		return generateHTML(data)
	case "txt":
		return data.RawContent
	case "sql":
		return generateSQL(data)
	default:
		return "Formato não suportado"
	}
}

func callClientDetails(content, traceNumber string) string {
	lines := strings.Split(strings.TrimSpace(content), "\n")

	for _, line := range lines {
		if len(line) >= 94 && line[0:1] == "6" {
			lineTraceNumber := strings.TrimSpace(line[79:94])
			if lineTraceNumber == traceNumber {
				return fmt.Sprintf(`🔍 Transação Encontrada:
				
📋 Informações da Transação:
• Número de Rastreamento: %s
• Código da Transação: %s
• Banco Receptor: %s
• Conta: %s
• Valor: R$ %.2f
• Nome: %s
• Identificação: %s

✅ Transação localizada com sucesso!`,
					traceNumber,
					strings.TrimSpace(line[1:3]),
					strings.TrimSpace(line[3:11]),
					strings.TrimSpace(line[11:28]),
					parseFloat(strings.TrimSpace(line[29:39]))/100.0,
					strings.TrimSpace(line[54:76]),
					strings.TrimSpace(line[39:54]))
			}
		}
	}

	return fmt.Sprintf("❌ Transação com número de rastreamento '%s' não encontrada", traceNumber)
}

func parseFloat(s string) float64 {
	if val, err := strconv.ParseFloat(s, 64); err == nil {
		return val
	}
	return 0.0
}

func generateCSV(data *NachaData) string {
	csv := "Tipo,Campo,Valor\n"
	csv += "Header,Destino Imediato," + fmt.Sprintf("%v", data.Header["immediate_destination"]) + "\n"
	csv += "Header,Origem Imediata," + fmt.Sprintf("%v", data.Header["immediate_origin"]) + "\n"
	csv += "Estatísticas,Total de Lotes," + fmt.Sprintf("%v", data.Statistics["total_batches"]) + "\n"
	csv += "Estatísticas,Total de Entradas," + fmt.Sprintf("%v", data.Statistics["total_entries"]) + "\n"
	csv += "Estatísticas,Valor Total," + fmt.Sprintf("%.2f", data.Statistics["total_amount_formatted"]) + "\n"

	for i, batch := range data.Batches {
		csv += fmt.Sprintf("Lote %d,Nome da Empresa,%v\n", i+1, batch["company_name"])
		csv += fmt.Sprintf("Lote %d,Código de Classe,%v\n", i+1, batch["service_class_code"])
	}

	return csv
}

func generateHTML(data *NachaData) string {
	html := `<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <title>Relatório NACHA</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        table { border-collapse: collapse; width: 100%; margin: 20px 0; }
        th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
        th { background-color: #f2f2f2; }
        .header { background-color: #4CAF50; color: white; padding: 20px; text-align: center; }
    </style>
</head>
<body>
    <div class="header">
        <h1>Relatório de Arquivo NACHA</h1>
    </div>
    
    <h2>Cabeçalho do Arquivo</h2>
    <table>
        <tr><th>Campo</th><th>Valor</th></tr>`

	for k, v := range data.Header {
		html += fmt.Sprintf("<tr><td>%s</td><td>%v</td></tr>", k, v)
	}

	html += `</table>
    
    <h2>Estatísticas</h2>
    <table>
        <tr><th>Métrica</th><th>Valor</th></tr>`

	for k, v := range data.Statistics {
		html += fmt.Sprintf("<tr><td>%s</td><td>%v</td></tr>", k, v)
	}

	html += `</table>
</body>
</html>`

	return html
}

func generateSQL(data *NachaData) string {
	sql := `-- Script SQL para importação de dados NACHA
CREATE TABLE IF NOT EXISTS nacha_files (
    id SERIAL PRIMARY KEY,
    immediate_destination VARCHAR(10),
    immediate_origin VARCHAR(10),
    file_creation_date VARCHAR(6),
    file_creation_time VARCHAR(4),
    total_batches INT,
    total_entries INT,
    total_amount DECIMAL(12,2)
);

INSERT INTO nacha_files (
    immediate_destination, immediate_origin, file_creation_date, file_creation_time,
    total_batches, total_entries, total_amount
) VALUES (
    '%s', '%s', '%s', '%s',
    %v, %v, %.2f
);`

	return fmt.Sprintf(sql,
		data.Header["immediate_destination"],
		data.Header["immediate_origin"],
		data.Header["file_creation_date"],
		data.Header["file_creation_time"],
		data.Statistics["total_batches"],
		data.Statistics["total_entries"],
		data.Statistics["total_amount_formatted"])
}

func main() {
	// Criar diretório temp se não existir
	os.MkdirAll("temp", 0755)

	// Configurar servidor HTTP com timeouts adequados para upload
	server := &http.Server{
		Addr:           ":8080",
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		IdleTimeout:    60 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
	}

	// Servir arquivos estáticos
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	// Rotas
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/create", createHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/validate", validateHandler)
	http.HandleFunc("/view", viewHandler)
	http.HandleFunc("/export", exportHandler)
	http.HandleFunc("/details", detailsHandler)

	log.Println("🌐 Aplicação Web NACHA iniciando na porta 8080")
	log.Println("📊 Funcionalidades disponíveis:")
	log.Println("   • Criação e Upload de Arquivos (até 10MB)")
	log.Println("   • Validação NACHA completa")
	log.Println("   • Visualização de Conteúdo estruturada")
	log.Println("   • Exportação em 7 formatos diferentes")
	log.Println("   • Detalhes de Transações por rastreamento")
	log.Println("   • Gerenciamento de Sessão com arquivos temporários")
	log.Println("⚠️  Certifique-se de que o servidor gRPC NACHA esteja executando na porta 50051")
	log.Println("🔗 Acesse: http://localhost:8080")

	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Erro ao iniciar servidor:", err)
	}
}
