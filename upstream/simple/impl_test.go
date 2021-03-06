package simple_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/ironsmile/nedomi/config"
	"github.com/ironsmile/nedomi/upstream/simple"
)

func TestGetRequest(t *testing.T) {
	var b = "Hello world is an awesome test string"
	var start, end uint64 = 3, 500000
	file, err := ioutil.TempFile(os.TempDir(), "test")
	if err != nil {
		t.Fatal(err)
	}
	_, err = io.WriteString(file, b)
	if err != nil {
		t.Fatal(err)
	}
	file.Sync()
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatal(err)
	}
	go func() {
		http.Serve(listener, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.ServeContent(w, r, "", time.Time{}, file)
		}))
	}()

	cz := &config.CacheZoneSection{
		ID:             "42",
		Path:           "/var/tmp",
		StorageObjects: 1024,
		PartSize:       2,
	}
	var m = make(map[string]*config.CacheZoneSection)
	m[cz.ID] = cz

	vh := new(config.VirtualHost)
	vh.Name = "localhost"
	vh.HandlerType = "proxy"
	vh.CacheKey = "test"
	vh.UpstreamType = "simple"
	vh.UpstreamAddress, _ = url.Parse("http://" + listener.Addr().String())
	vh.CacheZone = cz
	if err := vh.Validate(); err != nil {
		t.Fatal(err)
	}

	upstreamURL, err := url.Parse("http://" + listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}

	u := simple.New(upstreamURL)

	resp, err := u.GetRequestPartial(file.Name(), start, end)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	expected := []byte(b)[start:]
	if !bytes.Equal(expected, result) {
		t.Fatalf("Expected [%s] got [%s]", expected, result)
	}
}
