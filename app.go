package codegen

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

// Имена файлов, расположенных в корневой директории каждого шаблона.
// Они будут исключены из результатов генерации.
const (
	// metaFileName хранит имя, описание и замены выбранного шаблона
	metaFileName = ".meta.json"
	// helpFileName показывает пользователю справку для пост-генерации выбранного шаблона
	helpFileName = ".help.md"
)

const (
	// dirPermissions drwxr-xr-x
	dirPermissions = 0755
	// filePermissions -rw-r--r--
	filePermissions = 0644
)

type App struct {
	log          *logrus.Logger
	RootDir      string
	TemplatesDir string
}

func NewApp(log *logrus.Logger, rootDir, templatesDir string) *App {
	return &App{
		log:          log,
		RootDir:      rootDir,
		TemplatesDir: templatesDir,
	}
}

func (a *App) GetTemplatesList(ctx context.Context) ([]Template, error) {
	var results []Template

	files, err := os.ReadDir(a.TemplatesDir)
	if err != nil {
		return results, fmt.Errorf("get templates: %w", err)
	}

	// Parse ".meta.json" from each template dir
	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		template := Template{
			Id: file.Name(),
		}

		jsonString, err := os.ReadFile(a.TemplatesDir + "/" + file.Name() + "/.meta.json")
		if err != nil {
			return results, fmt.Errorf("read .meta.json template(%s) file: %w", file.Name(), err)
		}

		err = json.Unmarshal(jsonString, &template)
		if err != nil {
			return results, fmt.Errorf("json template(%s): %w", file.Name(), err)
		}

		results = append(results, template)
	}

	return results, nil
}

func (a *App) RawList(ctx context.Context, templateID, targetRelativeRootDir string, placeholders []Placeholder) ([]PreviewListItem, error) {
	list, err := a.PreviewList(ctx, templateID, targetRelativeRootDir, nil)
	if err != nil {
		return list, fmt.Errorf("raw preview: %w", err)
	}

	// Check if replaced dir/file names already exists
	for k, item := range list {
		path := filepath.Join(a.RootDir, targetRelativeRootDir, item.Path)
		for _, placeholder := range placeholders {
			path = strings.ReplaceAll(path, placeholder.Value, placeholder.Replace)
		}
		list[k].IsNew = a.checkFileExits(path)
	}

	return list, nil
}

func (a *App) PreviewList(ctx context.Context, templateID, targetRelativeRootDir string, placeholders []Placeholder) ([]PreviewListItem, error) {
	var items []PreviewListItem
	templateDir := filepath.Join(a.TemplatesDir, templateID)
	targetDir := filepath.Join(a.RootDir, targetRelativeRootDir)

	err := filepath.WalkDir(templateDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if path == templateDir {
			return nil // ignore root dir
		}

		if a.isIgnoreFile(d.Name()) {
			return nil // ignore template meta-files
		}

		// extract "/models/product.go" from "/var/www/project/codegen/templates/example1/models/product.go"
		relativeTargetPath := strings.Replace(path, templateDir, "", 1)
		for _, placeholder := range placeholders {
			relativeTargetPath = strings.ReplaceAll(relativeTargetPath, placeholder.Value, placeholder.Replace)
		}

		var content string

		if !d.IsDir() {
			content, err = a.replacePlaceholdersInFile(path, placeholders)
			if err != nil {
				return err
			}
		}

		destinationFile := filepath.Join(targetDir, relativeTargetPath)

		item := PreviewListItem{
			IsDir:   d.IsDir(),
			IsNew:   a.checkFileExits(destinationFile),
			Path:    strings.ReplaceAll(relativeTargetPath, "\\", "/"),
			Content: content,
		}

		items = append(items, item)

		return nil
	})

	if err != nil {
		return items, err
	}

	return items, nil
}

func (a *App) Generate(ctx context.Context, templateID, targetRelativeRootDir string, placeholders []Placeholder) (string, error) {
	previewList, err := a.PreviewList(ctx, templateID, targetRelativeRootDir, placeholders)
	if err != nil {
		return "", fmt.Errorf("preview: %w", err)
	}

	for _, previewListItem := range previewList {
		destinationDir := filepath.Join(a.RootDir, targetRelativeRootDir)
		destinationFile := filepath.Join(destinationDir, previewListItem.Path)

		if previewListItem.IsDir {
			if previewListItem.IsNew {
				err = os.Mkdir(destinationFile, dirPermissions)
				if err != nil {
					return "", fmt.Errorf("create new dir: %w", err)
				}
			}

			continue // directory handled
		}

		// Ensure dir path exist
		err = os.MkdirAll(destinationDir, dirPermissions)
		if err != nil {
			return "", fmt.Errorf("mkdir all: %w", err)
		}

		err = os.WriteFile(destinationFile, []byte(previewListItem.Content), filePermissions)
		if err != nil {
			return "", fmt.Errorf("write file: %w", err)
		}
	}

	// Generate help message
	templateHelpFilepath := filepath.Join(a.TemplatesDir, templateID, helpFileName)

	helpString, err := a.replacePlaceholdersInFile(templateHelpFilepath, placeholders)
	if err != nil {
		helpString = "Files successfully generated! But caused error while generating help message: " + err.Error()
	}

	return helpString, nil
}

func (a *App) GetDirectories(ctx context.Context, selectedDir string) ([]string, error) {
	var result []string

	dirs, err := os.ReadDir(filepath.Join(a.RootDir, selectedDir))
	if err != nil && !os.IsNotExist(err) {
		return result, fmt.Errorf("read dir: %w", err)
	}

	for _, dirEntry := range dirs {
		if dirEntry.IsDir() && dirEntry.Name()[:1] != "." {
			result = append(result, dirEntry.Name())
		}
	}

	return result, nil
}

func (a *App) replacePlaceholdersInFile(path string, placeholders []Placeholder) (string, error) {
	input, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	for _, placeholder := range placeholders {
		fromBytes := []byte(placeholder.Value)
		toBytes := []byte(placeholder.Replace)

		input = bytes.ReplaceAll(input, fromBytes, toBytes)
	}

	return string(input), nil
}

func (a *App) checkFileExits(path string) bool {
	_, err := os.Stat(path)

	return errors.Is(err, os.ErrNotExist)
}

func (a *App) isIgnoreFile(name string) bool {
	return name == metaFileName || name == helpFileName
}
