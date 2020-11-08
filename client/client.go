package client

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"strconv"

	"github.com/gorilla/websocket"
)

// Version WebShell Client current version
const Version = "1.1"

// UserAgent Request header[User-Agent]
var UserAgent = fmt.Sprintf("web-shell-client/%s (%s; %s; %s)", Version, runtime.GOOS, runtime.GOARCH, runtime.Version())

// WebShellClient connect to WebShellServer
type WebShellClient struct {
	Client *http.Client
	Dialer *websocket.Dialer
}

// Init http client
func (c *WebShellClient) Init(https bool, crt string) {
	if https {
		certBytes, err := ioutil.ReadFile(crt)
		if err != nil {
			panic("Unable to read cert.pem")
		}
		clientCertPool := x509.NewCertPool()
		if ok := clientCertPool.AppendCertsFromPEM(certBytes); !ok {
			panic("failed to parse root certificate")
		}
		tlsConfig := &tls.Config{RootCAs: clientCertPool}
		c.Client = &http.Client{Transport: &http.Transport{TLSClientConfig: tlsConfig}}
		c.Dialer = &websocket.Dialer{TLSClientConfig: tlsConfig}
	} else {
		c.Client = &http.Client{}
		c.Dialer = &websocket.Dialer{}
	}
}

// Run WebShellClient
func (c *WebShellClient) Run(https bool, host, post, contentpath string) {
	path, err := LoginServer(https, host, post, contentpath, c.GetJSON)
	if err != nil {
		log.Println("Login to Server failed:", err.Error())
		return
	}
	ConnectSocket(https, host, post, contentpath, path, UserAgent, c.GetWebsocket)
}

// GetRes http get request
func (c *WebShellClient) GetRes(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", UserAgent)
	if err != nil {
		return nil, err
	}
	return c.Client.Do(req)
}

// GetJSON http get request and parse JSON
func (c *WebShellClient) GetJSON(url string) (map[string]interface{}, error) {
	res, err := c.GetRes(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, errors.New("response status is " + strconv.Itoa(res.StatusCode))
	}
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	data := make(map[string]interface{})
	err = json.Unmarshal(bytes, &data)
	return data, err
}

// GetWebsocket get websocket connection
func (c *WebShellClient) GetWebsocket(url string) (*websocket.Conn, error) {
	h := make(http.Header)
	h["User-Agent"] = []string{UserAgent}
	skt, _, err := c.Dialer.Dial(url, h)
	return skt, err
}
