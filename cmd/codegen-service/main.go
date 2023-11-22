package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/demisang/codegen"
	"github.com/demisang/codegen/rest"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()

	// Parse flags
	rootDirFlag := flag.String("root", "", "where project root files are placed")
	templatesDirFlag := flag.String("templates", "", "where templates are placed")

	flag.Parse()

	// Directories
	rootDir, err := getAbsoluteDirPathFromFlagValue("root dir", rootDirFlag)
	if err != nil {
		log.Fatal(err)
	}

	templatesDir, err := getAbsoluteDirPathFromFlagValue("templates dir", templatesDirFlag)
	if err != nil {
		log.Fatal(err)
	}

	// New context
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// New application
	app := codegen.NewApp(log, rootDir, templatesDir)

	// Run http-server for GUI
	if err := runGuiServer(ctx, app, log, "0.0.0.0", 4765); err != nil {
		log.Fatal(err)
	}
}

func runGuiServer(ctx context.Context, app *codegen.App, log *logrus.Logger, host string, port int) error {
	server := rest.NewServer(app, log, host, port)

	return server.Run(ctx)
}

func getAbsoluteDirPathFromFlagValue(name string, value *string) (string, error) {
	if value == nil || *value == "" {
		return "", errors.New(name + " required")
	}

	absolute, err := filepath.Abs(*value)
	if err != nil {
		return "", fmt.Errorf(name+" path is wrong: %w", err)
	}

	return absolute, nil
}
