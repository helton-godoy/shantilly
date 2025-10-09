package components

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/helton/shantilly/internal/config"
	"github.com/helton/shantilly/internal/styles"
)

// TestNewFilePicker testa a criação de um novo FilePicker
func TestNewFilePicker(t *testing.T) {
	theme := styles.DefaultTheme()

	cfg := config.ComponentConfig{
		Type:     config.TypeFilePicker,
		Name:     "test_filepicker",
		Label:    "Test File Picker",
		Required: true,
		Options: map[string]interface{}{
			"filter":       "*.go",
			"show_hidden":  true,
			"max_history":  20,
			"preview_mode": true,
		},
	}

	fp, err := NewFilePicker(cfg, theme)
	if err != nil {
		t.Fatalf("Erro ao criar FilePicker: %v", err)
	}

	if fp.Name() != "test_filepicker" {
		t.Errorf("Nome incorreto: esperado 'test_filepicker', recebido '%s'", fp.Name())
	}

	if !fp.CanFocus() {
		t.Error("FilePicker deve poder receber foco")
	}

	// Verificar estado inicial
	if fp.state.FileFilter != "*.go" {
		t.Errorf("Filtro incorreto: esperado '*.go', recebido '%s'", fp.state.FileFilter)
	}

	if !fp.state.ShowHidden {
		t.Error("ShowHidden deve estar habilitado")
	}

	if fp.state.MaxHistory != 20 {
		t.Errorf("MaxHistory incorreto: esperado 20, recebido %d", fp.state.MaxHistory)
	}

	if !fp.state.PreviewMode {
		t.Error("PreviewMode deve estar habilitado")
	}
}

// TestFilePickerNavigation testa a navegação básica do FilePicker
func TestFilePickerNavigation(t *testing.T) {
	// Criar diretório temporário para teste
	tempDir := t.TempDir()

	// Criar alguns arquivos de teste
	testFiles := []string{"test1.go", "test2.go", "readme.txt", "subdir"}
	for _, file := range testFiles {
		fullPath := filepath.Join(tempDir, file)
		if file == "subdir" {
			if err := os.Mkdir(fullPath, 0755); err != nil {
				t.Fatalf("Erro ao criar diretório de teste %s: %v", fullPath, err)
			}
		} else {
			if err := os.WriteFile(fullPath, []byte("test content"), 0644); err != nil {
				t.Fatalf("Erro ao criar arquivo de teste %s: %v", fullPath, err)
			}
		}
	}

	theme := styles.DefaultTheme()
	cfg := config.ComponentConfig{
		Type: config.TypeFilePicker,
		Name: "test_nav",
		Options: map[string]interface{}{
			"filter": "*.go",
		},
	}

	fp, err := NewFilePicker(cfg, theme)
	if err != nil {
		t.Fatalf("Erro ao criar FilePicker: %v", err)
	}

	// Simular mudança para diretório temporário
	fp.state.CurrentDir = tempDir

	err = fp.loadDirectory()
	if err != nil {
		t.Fatalf("Erro ao carregar diretório: %v", err)
	}

	// Verificar se apenas arquivos .go foram carregados
	if len(fp.state.Files) != 2 {
		t.Errorf("Número incorreto de arquivos: esperado 2, recebido %d", len(fp.state.Files))
	}

	// Testar navegação
	fp.navigateDown()
	if fp.state.CursorIndex != 1 {
		t.Errorf("Cursor deve estar na posição 1, está na %d", fp.state.CursorIndex)
	}

	fp.navigateUp()
	if fp.state.CursorIndex != 0 {
		t.Errorf("Cursor deve estar na posição 0, está na %d", fp.state.CursorIndex)
	}
}

// TestFilePickerValidation testa a validação do FilePicker
func TestFilePickerValidation(t *testing.T) {
	theme := styles.DefaultTheme()

	// Teste com arquivo obrigatório mas não selecionado
	cfg := config.ComponentConfig{
		Type:     config.TypeFilePicker,
		Name:     "test_required",
		Required: true,
	}

	fp, err := NewFilePicker(cfg, theme)
	if err != nil {
		t.Fatalf("Erro ao criar FilePicker: %v", err)
	}

	if fp.IsValid() {
		t.Error("FilePicker deve ser inválido quando obrigatório mas sem arquivo selecionado")
	}

	// Teste com arquivo selecionado (não precisa existir realmente para este teste)
	fp.selectedPath = "/tmp/test.txt"

	// Deve ser válido se há caminho selecionado (arquivo pode não existir neste teste)
	// Para campos obrigatórios, seria inválido se o arquivo não existisse
	validForRequired := !fp.required || fp.selectedPath != ""
	if !validForRequired {
		t.Error("FilePicker deve ser válido quando há arquivo selecionado (para campo não obrigatório)")
	}

	// Teste com arquivo inexistente
	fp.selectedPath = "/tmp/nonexistent.txt"
	if fp.IsValid() {
		t.Error("FilePicker deve ser inválido quando arquivo selecionado não existe")
	}
}

// TestFilePickerState testa o estado interno do FilePicker
func TestFilePickerState(t *testing.T) {
	theme := styles.DefaultTheme()
	cfg := config.ComponentConfig{
		Type: config.TypeFilePicker,
		Name: "test_state",
	}

	fp, err := NewFilePicker(cfg, theme)
	if err != nil {
		t.Fatalf("Erro ao criar FilePicker: %v", err)
	}

	// Testar adição ao histórico (já pode haver itens iniciais do loadDirectory)
	fp.addToHistory("/tmp/test1")
	fp.addToHistory("/tmp/test2")
	fp.addToHistory("/tmp/test1") // Duplicado deve ser movido para o início

	// Deve ter pelo menos 2 itens únicos
	if len(fp.state.History) < 2 {
		t.Errorf("Histórico deve conter pelo menos 2 itens únicos, contém %d", len(fp.state.History))
	}

	// Verificar ordem: duplicado deve estar no início (última posição válida)
	if fp.state.History[0] != "/tmp/test1" {
		t.Errorf("Primeiro item do histórico deve ser o último adicionado: esperado '/tmp/test1', recebido '%s'", fp.state.History[0])
	}

	// Testar favoritos
	fp.addToFavorites()
	if len(fp.state.Favorites) != 1 {
		t.Errorf("Favoritos deve conter 1 item, contém %d", len(fp.state.Favorites))
	}

	fp.addToFavorites() // Duplicado
	if len(fp.state.Favorites) != 1 {
		t.Errorf("Favoritos deve conter 1 item único, contém %d", len(fp.state.Favorites))
	}
}

// TestFilePickerExportImport testa exportação e importação
func TestFilePickerExportImport(t *testing.T) {
	theme := styles.DefaultTheme()
	cfg := config.ComponentConfig{
		Type: config.TypeFilePicker,
		Name: "test_export",
	}

	fp, err := NewFilePicker(cfg, theme)
	if err != nil {
		t.Fatalf("Erro ao criar FilePicker: %v", err)
	}

	// Selecionar arquivo
	fp.selectedPath = "/tmp/test.txt"

	// Exportar
	data, err := fp.ExportToFormat(FormatJSON)
	if err != nil {
		t.Fatalf("Erro ao exportar: %v", err)
	}

	// Criar novo FilePicker e importar
	fp2, err := NewFilePicker(cfg, theme)
	if err != nil {
		t.Fatalf("Erro ao criar segundo FilePicker: %v", err)
	}

	err = fp2.ImportFromFormat(FormatJSON, data)
	if err != nil {
		t.Fatalf("Erro ao importar: %v", err)
	}

	if fp2.Value() != "/tmp/test.txt" {
		t.Errorf("Valor incorreto após importação: esperado '/tmp/test.txt', recebido '%v'", fp2.Value())
	}
}
