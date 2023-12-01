package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/demisang/codegen/internal"
	"github.com/demisang/codegen/internal/rest"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

type IntegrationTestSuite struct {
	suite.Suite
	ctx       context.Context
	host      string
	log       *logrus.Logger
	app       *internal.App
	server    *rest.Server
	client    *http.Client
	outputDir string
}

func TestIntegrationSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

const (
	host         = "localhost"
	port         = 4765
	rootDir      = "."
	templatesDir = "../templates/go-onion"
	outputDir    = "./_output"
)
const outputRelativeRootDir = "_output"

func (s *IntegrationTestSuite) SetupSuite() {
	var err error

	s.outputDir, err = filepath.Abs(outputDir)

	s.Require().NoError(err)

	serviceURL := url.URL{
		Scheme: "http",
		Host:   net.JoinHostPort(host, strconv.Itoa(port)),
	}
	s.host = serviceURL.String()
	s.ctx = context.Background()
	s.log = logrus.New()

	// New application
	s.app = internal.NewApp(s.log, rootDir, templatesDir)
	s.server = rest.NewServer(s.app, s.log, host, port)
	s.client = &http.Client{}

	go func() {
		err := s.server.Run(s.ctx, []func(){})
		s.Require().NoError(err)
	}()
	time.Sleep(100 * time.Millisecond)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	// Remove all generated files after all tests
	err := os.RemoveAll(s.outputDir)
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) sendRequest(ctx context.Context, method, endpoint string, data, dst any) int {
	s.T().Helper()

	var body []byte

	var err error

	if data != nil {
		body, err = json.Marshal(data)
		s.Require().NoError(err)
	}

	request, err := http.NewRequestWithContext(ctx, method, s.host+endpoint, bytes.NewReader(body))
	request.Header.Add("Content-Type", "application/json; charset=utf-8")
	request.Header.Add("Accept", "*/*")
	s.Require().NoError(err)

	response, err := s.client.Do(request)
	s.Require().NoError(err)

	if dst != nil {
		err = json.NewDecoder(response.Body).Decode(dst)
		s.Require().NoError(err)
	}

	return response.StatusCode
}
