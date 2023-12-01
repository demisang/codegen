package codegen

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"net/url"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"syscall"

	"github.com/demisang/codegen/internal"
	"github.com/demisang/codegen/internal/rest"
	"github.com/sirupsen/logrus"
)

func Run() error {
	log := logrus.New()

	// Parse flags
	rootDirFlag := flag.String("root", "", "where project root files are placed")
	templatesDirFlag := flag.String("templates", "", "where templates are placed")
	serviceHostFlag := flag.String("host", "0.0.0.0", "service host")
	servicePortFlag := flag.Int("port", 4765, "service port")
	openBrowserFlag := flag.Bool("browser", true, "open browser GUI after service running")

	flag.Parse()

	// Directories
	rootDir, err := getAbsoluteDirPathFromFlagValue("root dir", rootDirFlag)
	if err != nil {
		return fmt.Errorf("root dir absolute path: %w", err)
	}

	templatesDir, err := getAbsoluteDirPathFromFlagValue("templates dir", templatesDirFlag)
	if err != nil {
		return fmt.Errorf("templates dir absolute path: %w", err)
	}

	// Service addr
	if serviceHostFlag == nil {
		return errors.New("service host is nil")
	}

	if servicePortFlag == nil {
		return errors.New("service port is nil")
	}

	// New context
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// New application
	app := internal.NewApp(log, rootDir, templatesDir)

	// Run http-server for service
	if err := runGuiServer(ctx, app, log, *serviceHostFlag, *servicePortFlag, *openBrowserFlag); err != nil {
		return fmt.Errorf("server run: %w", err)
	}

	return nil
}

func runGuiServer(ctx context.Context, app *internal.App, log *logrus.Logger, host string, port int, browser bool) error {
	server := rest.NewServer(app, log, host, port)

	var onStarted []func()

	onStarted = append(onStarted, func() {
		serviceURL := url.URL{
			Scheme: "http",
			Host:   net.JoinHostPort(host, strconv.Itoa(port)),
		}
		if host == "0.0.0.0" {
			serviceURL.Host = net.JoinHostPort("localhost", strconv.Itoa(port))
		}

		if !browser {
			log.Infof("\nService available on %s\n", serviceURL.String())

			return
		}

		//goland:noinspection HttpUrlsUsage
		err := openBrowserTab(serviceURL.String())
		if err != nil {
			log.Errorf("open browser: %v", err)
		}
	})

	return server.Run(ctx, onStarted)
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

func openBrowserTab(url string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	return err
}
