package tests

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/demisang/codegen/internal"
)

type generateRequest struct {
	TemplateID   string                 `json:"template_id"`
	TargetDir    string                 `json:"target_dir"`
	Placeholders []internal.Placeholder `json:"placeholders"`
}

var generateModelRequest = generateRequest{
	TemplateID: "model",
	TargetDir:  outputRelativeRootDir,
	Placeholders: internal.Placeholders{
		internal.Placeholder{Value: "__PascalCase__", Replace: "ProductOrder"},
		internal.Placeholder{Value: "__camelCase__", Replace: "productOrder"},
		internal.Placeholder{Value: "__snake_case__", Replace: "product_order"},
	},
}

func (s *IntegrationTestSuite) TestGetTemplatesList() {
	s.Run("get templates ok", func() {
		var templates []internal.Template

		code := s.sendRequest(s.ctx, http.MethodGet, "/templates", nil, &templates)

		s.Require().Equal(http.StatusOK, code)
		s.Require().Len(templates, 2, "returned 2 templates")
		s.Require().Equal(templates[0].Name, "Model")
		s.Require().Equal(templates[1].Name, "Store")
	})
}

func (s *IntegrationTestSuite) TestPreview() {
	s.Run("preview model template ok", func() {
		var previewItems []internal.PreviewListItem

		code := s.sendRequest(s.ctx, http.MethodPost, "/raw-list", generateModelRequest, &previewItems)

		s.Require().Equal(http.StatusOK, code)
		s.Require().Len(previewItems, 1, "returned 1 file")
		s.Require().False(previewItems[0].IsDir)
		s.Require().Equal(previewItems[0].Path, "/__snake_case__.go")
		s.Require().Contains(previewItems[0].Content, generateModelRequest.Placeholders[0].Value)
		s.Require().Contains(previewItems[0].Content, generateModelRequest.Placeholders[1].Value)
	})
}

func (s *IntegrationTestSuite) TestGenerate() {
	s.Run("generate model template ok", func() {
		var helpMessage string

		code := s.sendRequest(s.ctx, http.MethodPost, "/generate", generateModelRequest, &helpMessage)

		s.Require().Equal(http.StatusOK, code)
		s.Require().NotEmpty(helpMessage, "help message returned")
		s.Require().Contains(helpMessage, generateModelRequest.Placeholders[0].Replace)

		// Check files are created
		filename := filepath.Join(s.outputDir, generateModelRequest.Placeholders[2].Replace+".go")
		fileContentBytes, err := os.ReadFile(filename)
		s.Require().NoError(err, "file created")
		fileContent := string(fileContentBytes)

		s.Require().Contains(fileContent, generateModelRequest.Placeholders[0].Replace)
		s.Require().Contains(fileContent, generateModelRequest.Placeholders[1].Replace)
	})
}
