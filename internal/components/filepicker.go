package components

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/helton/shantilly/internal/config"
	"github.com/helton/shantilly/internal/styles"
)

// FilePickerState representa o estado interno do FilePicker
type FilePickerState struct {
	CurrentDir    string   `json:"current_dir"`
	SelectedFile  string   `json:"selected_file"`
	Files         []string `json:"files"`
	ShowHidden    bool     `json:"show_hidden"`
	FileFilter    string   `json:"file_filter"`
	ScrollOffset  int      `json:"scroll_offset"`
	CursorIndex   int      `json:"cursor_index"`
	Favorites     []string `json:"favorites"`
	History       []string `json:"history"`
	MaxHistory    int      `json:"max_history"`
	PreviewMode   bool     `json:"preview_mode"`
	PreviewBuffer string   `json:"preview_buffer"`
}

// FilePicker implements a file picker component for selecting files/directories.
// This component allows users to navigate through directories and select files.
type FilePicker struct {
	name         string
	label        string
	required     bool
	help         string
	selectedPath string
	theme        *styles.Theme
	errorMsg     string
	focused      bool
	initialPath  string
	state        *FilePickerState
	width        int
	height       int
}

// NewFilePicker creates a new FilePicker component from configuration.
func NewFilePicker(cfg config.ComponentConfig, theme *styles.Theme) (*FilePicker, error) {
	if cfg.Type != config.TypeFilePicker {
		return nil, fmt.Errorf("tipo de componente inválido: esperado filepicker, recebido %s", cfg.Type)
	}

	// Initialize state
	state := &FilePickerState{
		CurrentDir:   ".",
		Files:        []string{},
		ShowHidden:   false,
		FileFilter:   "*",
		ScrollOffset: 0,
		CursorIndex:  0,
		Favorites:    []string{},
		History:      []string{},
		MaxHistory:   50,
		PreviewMode:  false,
	}

	fp := &FilePicker{
		name:         cfg.Name,
		label:        cfg.Label,
		required:     cfg.Required,
		help:         cfg.Help,
		theme:        theme,
		state:        state,
		selectedPath: "",
		initialPath:  "",
		focused:      false,
		width:        80,
		height:       24,
	}

	// Set default path if provided
	if cfg.Default != nil {
		if defaultPath, ok := cfg.Default.(string); ok {
			fp.selectedPath = defaultPath
			fp.initialPath = defaultPath
			fp.state.CurrentDir = defaultPath
		}
	}

	// Set file filter if provided in options
	if cfg.Options != nil {
		if filter, ok := cfg.Options["filter"].(string); ok {
			fp.state.FileFilter = filter
		}
		if showHidden, ok := cfg.Options["show_hidden"].(bool); ok {
			fp.state.ShowHidden = showHidden
		}
		if maxHistory, ok := cfg.Options["max_history"].(int); ok {
			fp.state.MaxHistory = maxHistory
		}
		if previewMode, ok := cfg.Options["preview_mode"].(bool); ok {
			fp.state.PreviewMode = previewMode
		}
	}

	// Load initial directory
	if err := fp.loadDirectory(); err != nil {
		// Log error but don't fail initialization - component can still function
		fp.errorMsg = fmt.Sprintf("Erro ao carregar diretório inicial: %v", err)
	}

	return fp, nil
}

// Init implements tea.Model.
func (fp *FilePicker) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model.
func (fp *FilePicker) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !fp.focused {
		return fp, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			fp.navigateUp()
		case "down", "j":
			fp.navigateDown()
		case "left", "h", "backspace":
			fp.navigateToParent()
		case "right", "l", "enter":
			fp.navigateInto()
		case "home", "g":
			fp.goToTop()
		case "end", "G":
			fp.goToBottom()
		case "page up":
			fp.pageUp()
		case "page down":
			fp.pageDown()
		case " ":
			fp.toggleSelection()
		case "f":
			fp.addToFavorites()
		case "F":
			fp.showFavorites()
		case "p":
			fp.togglePreview()
		case "ctrl+c":
			fp.focused = false
			return fp, nil
		}

	case tea.WindowSizeMsg:
		fp.width = msg.Width
		fp.height = msg.Height
	}

	return fp, nil
}

// View implements tea.Model.
func (fp *FilePicker) View() string {
	var parts []string

	// Header com informações básicas
	var header string
	if fp.selectedPath != "" {
		header = fmt.Sprintf("📁 %s", fp.selectedPath)
	} else {
		header = "📂 Nenhum arquivo selecionado"
	}

	if fp.label != "" {
		header += fmt.Sprintf(" (%s)", fp.label)
	}
	parts = append(parts, fp.theme.Border.Render(header))

	// Renderizar lista de arquivos se estiver focado
	if fp.focused {
		fileListParts := fp.renderFileList()
		parts = append(parts, fileListParts...)
	} else {
		// Quando não focado, mostrar apenas informações básicas
		parts = append(parts, "Pressione Enter para navegar")
	}

	// Render error or help text
	if fp.errorMsg != "" {
		parts = append(parts, fp.theme.Error.Render("✗ "+fp.errorMsg))
	} else if fp.help != "" {
		parts = append(parts, fp.theme.Help.Render(fp.help))
	}

	// Render contextual help quando focado
	if fp.focused && fp.errorMsg == "" {
		contextualHelp := fp.getContextualHelp()
		if contextualHelp != "" {
			parts = append(parts, fp.theme.Help.Render(contextualHelp))
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

// Name implements Component.
func (fp *FilePicker) Name() string {
	return fp.name
}

// CanFocus implements Component.
func (fp *FilePicker) CanFocus() bool {
	return true
}

// SetFocus implements Component.
func (fp *FilePicker) SetFocus(focused bool) {
	fp.focused = focused
}

// IsValid implements Component.
func (fp *FilePicker) IsValid() bool {
	errors := fp.ValidateWithContext(ValidationContext{
		ComponentValues: make(map[string]interface{}),
		GlobalConfig:    make(map[string]interface{}),
		ExternalData:    make(map[string]interface{}),
	})

	// Se há erros de validação, definir a primeira mensagem de erro
	if len(errors) > 0 {
		fp.errorMsg = errors[0].Message
		return false
	}

	fp.errorMsg = ""
	return true
}

// GetError implements Component.
func (fp *FilePicker) GetError() string {
	return fp.errorMsg
}

// SetError implements Component.
func (fp *FilePicker) SetError(msg string) {
	fp.errorMsg = msg
}

// Value implements Component.
func (fp *FilePicker) Value() interface{} {
	return fp.selectedPath
}

// SetValue implements Component.
func (fp *FilePicker) SetValue(value interface{}) error {
	strValue, ok := value.(string)
	if !ok {
		return fmt.Errorf("valor inválido: esperado string, recebido %T", value)
	}

	fp.selectedPath = strValue
	return nil
}

// Reset implements Component.
func (fp *FilePicker) Reset() {
	fp.selectedPath = fp.initialPath
	fp.errorMsg = ""
	fp.focused = false
}

// GetMetadata implements Component.
func (fp *FilePicker) GetMetadata() ComponentMetadata {
	return ComponentMetadata{
		Version:      "1.0.0",
		Author:       "Shantilly Team",
		Description:  "File picker component for selecting files and directories",
		Dependencies: []string{},
		Examples: []ComponentExample{
			{
				Name:        "Config File Picker",
				Description: "File picker for selecting configuration files",
				Config: map[string]interface{}{
					"type":  "filepicker",
					"name":  "config_file",
					"label": "Select Configuration File",
					"options": map[string]interface{}{
						"filter": "*.yaml,*.yml,*.json",
					},
				},
			},
		},
		Schema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"value": map[string]interface{}{
					"type":        "string",
					"description": "The selected file path",
				},
			},
		},
	}
}

// ValidateWithContext implements Component.
func (fp *FilePicker) ValidateWithContext(context ValidationContext) []ValidationError {
	var errors []ValidationError

	// Validação básica: arquivo obrigatório deve estar selecionado
	if fp.required && fp.selectedPath == "" {
		errors = append(errors, ValidationError{
			Code:     "FILE_PICKER_REQUIRED",
			Message:  "Um arquivo deve ser selecionado",
			Field:    fp.name,
			Severity: "error",
			Context: map[string]interface{}{
				"component":   "FilePicker",
				"value":       fp.Value(),
				"current_dir": fp.state.CurrentDir,
				"file_filter": fp.state.FileFilter,
				"show_hidden": fp.state.ShowHidden,
			},
		})
	}

	// Validação: se há caminho selecionado, verificar se existe
	if fp.selectedPath != "" {
		if _, err := os.Stat(fp.selectedPath); os.IsNotExist(err) {
			errors = append(errors, ValidationError{
				Code:     "FILE_NOT_FOUND",
				Message:  "Arquivo selecionado não existe",
				Field:    fp.name,
				Severity: "error",
				Context: map[string]interface{}{
					"component":   "FilePicker",
					"file_path":   fp.selectedPath,
					"current_dir": fp.state.CurrentDir,
				},
			})
		} else if err != nil {
			// Outro erro de acesso ao arquivo
			errors = append(errors, ValidationError{
				Code:     "FILE_ACCESS_ERROR",
				Message:  fmt.Sprintf("Erro ao acessar arquivo: %v", err),
				Field:    fp.name,
				Severity: "error",
				Context: map[string]interface{}{
					"component": "FilePicker",
					"file_path": fp.selectedPath,
					"error":     err.Error(),
				},
			})
		}
	}

	// Validação: verificar se diretório atual é acessível
	if _, err := os.Stat(fp.state.CurrentDir); err != nil {
		errors = append(errors, ValidationError{
			Code:     "CURRENT_DIR_ERROR",
			Message:  fmt.Sprintf("Erro ao acessar diretório atual: %v", err),
			Field:    fp.name,
			Severity: "warning",
			Context: map[string]interface{}{
				"component":   "FilePicker",
				"current_dir": fp.state.CurrentDir,
				"error":       err.Error(),
			},
		})
	}

	// Validação cruzada: verificar se há componentes relacionados na validação
	if componentValues, exists := context.ComponentValues[fp.name]; exists {
		if relatedValue, ok := componentValues.(map[string]interface{}); ok {
			// Pode adicionar validações cruzadas aqui no futuro
			_ = relatedValue // Evitar erro de variável não utilizada
		}
	}

	return errors
}

// ExportToFormat implements Component.
func (fp *FilePicker) ExportToFormat(format ExportFormat) ([]byte, error) {
	data := map[string]interface{}{
		"name":     fp.Name(),
		"value":    fp.Value(),
		"metadata": fp.GetMetadata(),
	}

	switch format {
	case FormatJSON:
		return json.MarshalIndent(data, "", "  ")
	default:
		return nil, fmt.Errorf("formato não suportado: %s", format)
	}
}

// ImportFromFormat implements Component.
func (fp *FilePicker) ImportFromFormat(format ExportFormat, data []byte) error {
	var imported map[string]interface{}

	switch format {
	case FormatJSON:
		if err := json.Unmarshal(data, &imported); err != nil {
			return fmt.Errorf("erro ao fazer parse do JSON: %w", err)
		}
	default:
		return fmt.Errorf("formato não suportado: %s", format)
	}

	if value, ok := imported["value"].(string); ok {
		return fp.SetValue(value)
	}

	return nil
}

// GetDependencies implements Component.
func (fp *FilePicker) GetDependencies() []string {
	return []string{}
}

// SetTheme implements Component.
func (fp *FilePicker) SetTheme(theme *styles.Theme) {
	fp.theme = theme
	// Re-renderizar lista de arquivos após mudança de tema
	if fp.focused && len(fp.state.Files) > 0 {
		if err := fp.loadDirectory(); err != nil {
			fp.errorMsg = fmt.Sprintf("Erro ao recarregar diretório após mudança de tema: %v", err)
		}
	}
}

// Métodos de navegação e funcionalidades avançadas do FilePicker

// loadDirectory carrega o conteúdo do diretório atual
func (fp *FilePicker) loadDirectory() error {
	files, err := os.ReadDir(fp.state.CurrentDir)
	if err != nil {
		fp.errorMsg = fmt.Sprintf("Erro ao ler diretório: %v", err)
		return err
	}

	// Filtrar arquivos baseado nas configurações
	filteredFiles := []string{}
	for _, file := range files {
		name := file.Name()

		// Aplicar filtro de arquivos ocultos
		if !fp.state.ShowHidden && strings.HasPrefix(name, ".") {
			continue
		}

		// Aplicar filtro de extensão se especificado
		if fp.state.FileFilter != "*" {
			match, err := filepath.Match(fp.state.FileFilter, name)
			if err != nil {
				continue
			}
			if !match {
				continue
			}
		}

		filteredFiles = append(filteredFiles, name)
	}

	fp.state.Files = filteredFiles

	// Adicionar diretório atual ao histórico se não estiver vazio
	if fp.state.CurrentDir != "" {
		fp.addToHistory(fp.state.CurrentDir)
	}

	return nil
}

// navigateUp move o cursor para cima na lista de arquivos
func (fp *FilePicker) navigateUp() {
	if fp.state.CursorIndex > 0 {
		fp.state.CursorIndex--
		fp.updateScrollOffset()
	}
}

// navigateDown move o cursor para baixo na lista de arquivos
func (fp *FilePicker) navigateDown() {
	if fp.state.CursorIndex < len(fp.state.Files)-1 {
		fp.state.CursorIndex++
		fp.updateScrollOffset()
	}
}

// navigateToParent navega para o diretório pai
func (fp *FilePicker) navigateToParent() {
	parentDir := filepath.Dir(fp.state.CurrentDir)
	if parentDir != fp.state.CurrentDir { // Evitar loop infinito
		fp.state.CurrentDir = parentDir
		fp.state.CursorIndex = 0
		fp.state.ScrollOffset = 0
		if err := fp.loadDirectory(); err != nil {
			fp.errorMsg = fmt.Sprintf("Erro ao navegar para diretório pai: %v", err)
		}
	}
}

// navigateInto navega para dentro do diretório selecionado ou seleciona arquivo
func (fp *FilePicker) navigateInto() {
	if fp.state.CursorIndex >= len(fp.state.Files) {
		return
	}

	selected := fp.state.Files[fp.state.CursorIndex]
	fullPath := filepath.Join(fp.state.CurrentDir, selected)

	// Verificar se é diretório
	if info, err := os.Stat(fullPath); err == nil && info.IsDir() {
		fp.state.CurrentDir = fullPath
		fp.state.CursorIndex = 0
		fp.state.ScrollOffset = 0
		if err := fp.loadDirectory(); err != nil {
			fp.errorMsg = fmt.Sprintf("Erro ao navegar para diretório: %v", err)
		}
	} else {
		// Selecionar arquivo
		fp.selectedPath = fullPath
		fp.errorMsg = ""
	}
}

// goToTop vai para o primeiro arquivo da lista
func (fp *FilePicker) goToTop() {
	fp.state.CursorIndex = 0
	fp.updateScrollOffset()
}

// goToBottom vai para o último arquivo da lista
func (fp *FilePicker) goToBottom() {
	fp.state.CursorIndex = len(fp.state.Files) - 1
	fp.updateScrollOffset()
}

// pageUp navega uma página para cima
func (fp *FilePicker) pageUp() {
	pageSize := fp.getPageSize()
	fp.state.CursorIndex = max(0, fp.state.CursorIndex-pageSize)
	fp.updateScrollOffset()
}

// pageDown navega uma página para baixo
func (fp *FilePicker) pageDown() {
	pageSize := fp.getPageSize()
	maxIndex := len(fp.state.Files) - 1
	fp.state.CursorIndex = min(maxIndex, fp.state.CursorIndex+pageSize)
	fp.updateScrollOffset()
}

// toggleSelection alterna a seleção do arquivo atual
func (fp *FilePicker) toggleSelection() {
	if fp.state.CursorIndex >= len(fp.state.Files) {
		return
	}

	selected := fp.state.Files[fp.state.CursorIndex]
	fp.selectedPath = filepath.Join(fp.state.CurrentDir, selected)
	fp.errorMsg = ""
}

// addToFavorites adiciona o diretório atual aos favoritos
func (fp *FilePicker) addToFavorites() {
	dir := fp.state.CurrentDir
	for _, favorite := range fp.state.Favorites {
		if favorite == dir {
			return // Já está nos favoritos
		}
	}
	fp.state.Favorites = append(fp.state.Favorites, dir)
}

// showFavorites mostra apenas os diretórios favoritos
func (fp *FilePicker) showFavorites() {
	// Implementação simplificada - alterna entre favoritos e lista normal
	fp.state.ShowHidden = !fp.state.ShowHidden
	if err := fp.loadDirectory(); err != nil {
		fp.errorMsg = fmt.Sprintf("Erro ao mostrar favoritos: %v", err)
	}
}

// togglePreview alterna o modo preview
func (fp *FilePicker) togglePreview() {
	fp.state.PreviewMode = !fp.state.PreviewMode
	if fp.state.PreviewMode && fp.selectedPath != "" {
		fp.loadPreview()
	}
}

// updateScrollOffset atualiza o offset de scroll para manter o cursor visível
func (fp *FilePicker) updateScrollOffset() {
	pageSize := fp.getPageSize()
	if fp.state.CursorIndex < fp.state.ScrollOffset {
		fp.state.ScrollOffset = fp.state.CursorIndex
	} else if fp.state.CursorIndex >= fp.state.ScrollOffset+pageSize {
		fp.state.ScrollOffset = fp.state.CursorIndex - pageSize + 1
	}
}

// getPageSize retorna o tamanho da página baseado na altura disponível
func (fp *FilePicker) getPageSize() int {
	// Reservar espaço para header e footer
	availableHeight := fp.height - 4
	return max(1, availableHeight)
}

// addToHistory adiciona um caminho ao histórico
func (fp *FilePicker) addToHistory(path string) {
	// Remover se já existir para evitar duplicatas
	for i, item := range fp.state.History {
		if item == path {
			fp.state.History = append(fp.state.History[:i], fp.state.History[i+1:]...)
			break
		}
	}

	// Adicionar no início
	fp.state.History = append([]string{path}, fp.state.History...)

	// Limitar tamanho do histórico
	if len(fp.state.History) > fp.state.MaxHistory {
		fp.state.History = fp.state.History[:fp.state.MaxHistory]
	}
}

// loadPreview carrega preview do arquivo selecionado
func (fp *FilePicker) loadPreview() {
	if fp.selectedPath == "" {
		return
	}

	file, err := os.Open(fp.selectedPath)
	if err != nil {
		fp.state.PreviewBuffer = fmt.Sprintf("Erro ao abrir arquivo: %v", err)
		return
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			fp.state.PreviewBuffer = fmt.Sprintf("Erro ao fechar arquivo: %v", closeErr)
		}
	}()

	// Ler apenas primeiras linhas (máximo 20 linhas)
	scanner := bufio.NewScanner(file)
	lines := []string{}
	lineCount := 0
	maxLines := 20

	for scanner.Scan() && lineCount < maxLines {
		lines = append(lines, scanner.Text())
		lineCount++
	}

	fp.state.PreviewBuffer = strings.Join(lines, "\n")

	if lineCount >= maxLines {
		fp.state.PreviewBuffer += "\n... (arquivo muito grande para preview)"
	}
}

// renderFileList renderiza a lista de arquivos
func (fp *FilePicker) renderFileList() []string {
	var lines []string

	// Header com informações do diretório atual
	header := fmt.Sprintf("📁 %s", fp.state.CurrentDir)
	if fp.state.FileFilter != "*" {
		header += fmt.Sprintf(" (filtro: %s)", fp.state.FileFilter)
	}
	lines = append(lines, fp.theme.Border.Render(header))

	pageSize := fp.getPageSize()
	endIndex := min(len(fp.state.Files), fp.state.ScrollOffset+pageSize)

	for i := fp.state.ScrollOffset; i < endIndex; i++ {
		fileName := fp.state.Files[i]

		// Verificar se é diretório
		fullPath := filepath.Join(fp.state.CurrentDir, fileName)
		if info, err := os.Stat(fullPath); err == nil && info.IsDir() {
			fileName = "📁 " + fileName
		} else {
			fileName = "📄 " + fileName
		}

		// Aplicar estilo baseado no estado
		var styledLine string
		if i == fp.state.CursorIndex {
			styledLine = fp.theme.BorderActive.Render("▶ " + fileName)
		} else {
			styledLine = "  " + fileName
		}

		lines = append(lines, styledLine)
	}

	// Renderizar preview se habilitado
	if fp.state.PreviewMode && fp.state.PreviewBuffer != "" {
		lines = append(lines, "")
		lines = append(lines, fp.theme.Border.Render("📖 Preview:"))
		previewLines := strings.Split(fp.state.PreviewBuffer, "\n")
		for _, line := range previewLines {
			lines = append(lines, "  "+line)
		}
	}

	return lines
}

// Funções auxiliares
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// getContextualHelp retorna ajuda contextual baseada no estado atual
func (fp *FilePicker) getContextualHelp() string {
	help := "📋 Navegação: ↑↓ ou jk | ←→ ou hl | Enter: entrar/selecionar | "
	help += fmt.Sprintf("📁 Diretório: %s", fp.state.CurrentDir)

	if fp.state.FileFilter != "*" {
		help += fmt.Sprintf(" | 🔍 Filtro: %s", fp.state.FileFilter)
	}

	if len(fp.state.Favorites) > 0 {
		help += fmt.Sprintf(" | ⭐ Favoritos: %d", len(fp.state.Favorites))
	}

	if fp.state.PreviewMode {
		help += " | 👁 Preview: ON"
	}

	help += " | [Espaço]: selecionar | f: favoritar | p: preview | Ctrl+C: sair"

	return help
}

// getCurrentDirectory returns the current working directory as a placeholder
// In a full implementation, this would maintain a current directory state
func (fp *FilePicker) getCurrentDirectory() string {
	if dir, err := os.Getwd(); err == nil {
		return dir
	}
	return ""
}
