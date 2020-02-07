package main

import (
	"compress/gzip"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"

	corev2 "github.com/sensu/sensu-go/api/core/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExecuteHandler(t *testing.T) {
	assert := assert.New(t)
	event := corev2.FixtureEvent("entity1", "check1")
	event.Check = nil
	event.Metrics = corev2.FixtureMetrics()
	event.Metrics.Points[0].Timestamp = 1580922166749062000

	var test = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gr, err := gzip.NewReader(r.Body)
		assert.NoError(err)
		body, err := ioutil.ReadAll(gr)
		assert.NoError(err)
		expectedBody := `"answer" 42 158092216674906200 source="entity1" "foo"="bar"`
		assert.Equal(expectedBody, strings.Trim(string(body), "\n"))
		w.WriteHeader(http.StatusOK)
	}))

	url, err := url.ParseRequestURI(test.URL)
	assert.NoError(err)
	config.Host = url.Hostname()
	port, err := strconv.Atoi(url.Port())
	require.NoError(t, err)
	config.Port = uint64(port)
	assert.NoError(sendMetrics(event))
}

func TestMain(t *testing.T) {
	assert := assert.New(t)
	file, _ := ioutil.TempFile(os.TempDir(), "sensu-go-graphite-handler")
	defer func() {
		_ = os.Remove(file.Name())
	}()

	event := corev2.FixtureEvent("entity1", "check1")
	event.Check = nil
	event.Metrics = corev2.FixtureMetrics()
	eventJSON, _ := json.Marshal(event)
	_, err := file.WriteString(string(eventJSON))
	require.NoError(t, err)
	require.NoError(t, file.Sync())
	_, err = file.Seek(0, 0)
	require.NoError(t, err)
	os.Stdin = file
	requestReceived := false

	var test = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestReceived = true
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"ok": true}`))
		require.NoError(t, err)
	}))

	url, err := url.ParseRequestURI(test.URL)
	assert.NoError(err)
	config.Host = url.Hostname()
	port, err := strconv.Atoi(url.Port())
	require.NoError(t, err)
	config.Port = uint64(port)
	oldArgs := os.Args
	os.Args = []string{"sensu-go-graphite-handler", "--host", url.Hostname(), "--port", url.Port()}
	defer func() { os.Args = oldArgs }()

	main()
	assert.True(requestReceived)
}
