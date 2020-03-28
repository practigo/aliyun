package acm

import (
	"bufio"
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/practigo/aliyun"
	"github.com/practigo/aliyun/sts"
)

// AccessKey is an alias for sts.Credentials
// that used to sign the ACM APIs.
type AccessKey sts.Credentials

// Sign signs the content as described in
// https://help.aliyun.com/document_detail/64129.html.
func (k *AccessKey) Sign(content string) string {
	// signature=`echo -n $content | openssl dgst -hmac $secretKey -sha1 -binary | base64`
	h := hmac.New(sha1.New, []byte(k.AccessKeySecret))
	h.Write([]byte(content)) // sha1 Write() returns no error
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// Request sends an ACM API. Only data bytes of a 200 OK
// response will be returned.
func Request(cl *http.Client, req *http.Request, ak *AccessKey, opt ConfigOption, timestamp string) (data []byte, err error) {
	toSign := opt.Tenant + opt.Group + timestamp
	req.Header.Add("Spas-AccessKey", ak.AccessKeyID)
	req.Header.Add("Spas-Signature", ak.Sign(toSign))
	req.Header.Add("timeStamp", timestamp)
	if ak.SecurityToken != "" {
		req.Header.Add("Spas-SecurityToken", ak.SecurityToken)
	}

	resp, err := cl.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return data, fmt.Errorf("status %d", resp.StatusCode)
	}

	return ioutil.ReadAll(resp.Body)
}

// Timestamp returns a timestamp string in ms.
func Timestamp() string {
	return strconv.FormatInt(time.Now().Unix(), 10) + "000" // shortcut to ms
}

// MD5 returns the md5 string.
func MD5(bs []byte) string {
	h := md5.New()
	h.Write(bs)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// ConfigOption specifies a certain config.
type ConfigOption struct {
	DataID string
	Group  string
	Tenant string // a.k.a. namespace
	MD5    string // for listen only
}

// Join returns a string form of the ConfigOption as
// `dataId^2group^2contentMD5^2tenant^1`.
func (o ConfigOption) Join() string {
	return strings.Join([]string{o.DataID, o.Group, o.MD5, o.Tenant}, "%02") + "%01"
}

// Service is the config service.
type Service struct {
	Host        string
	Port        string
	Path4IPs    string
	Path4Config string
	// Cl is used for all requests except for listenConfig.
	Cl *http.Client
	// Listener is for listenConfig which performs long polling.
	Listener *http.Client
}

// New returns an ACM Service with default port & path setting.
// Get the endpoint from https://help.aliyun.com/document_detail/64129.html.
func New(host string) Service {
	return Service{
		Host: host,
		// defaults
		Port:        ":8080",
		Path4IPs:    "/diamond-server/diamond",
		Path4Config: "/diamond-server/config.co",
		Cl:          aliyun.TimeoutClient(5 * time.Second),
		Listener:    aliyun.TimeoutClient(35 * time.Second),
	}
}

// GetServiceIPs gets the IPs for a certain service specified by the uri.
func GetServiceIPs(cl *http.Client, uri string) (ips []string, err error) {
	resp, err := cl.Get(uri)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ips, fmt.Errorf("status %d", resp.StatusCode)
	}

	ips = make([]string, 0)
	s := bufio.NewScanner(resp.Body)
	for s.Scan() {
		ips = append(ips, s.Text())
	}
	err = s.Err()
	return
}

// GetIPs gets the diamond service IPs.
func (s *Service) GetIPs() (ips []string, err error) {
	uri := "http://" + path.Join(s.Host+s.Port, s.Path4IPs)
	return GetServiceIPs(s.Cl, uri)
}

func (s *Service) getConfigRequest(opt ConfigOption, f func(string, ConfigOption) (*http.Request, error)) (*http.Request, error) {
	ips, err := s.GetIPs()
	if err != nil {
		return nil, fmt.Errorf("get service ips: %w", err)
	}
	if len(ips) < 1 {
		return nil, fmt.Errorf("no ip for diamond service")
	}
	uri := path.Join(ips[0]+s.Port, s.Path4Config) // just choose the first
	return f(uri, opt)
}

func configRequest(uri string, opt ConfigOption) (*http.Request, error) {
	v := url.Values{}
	v.Add("dataId", opt.DataID)
	v.Add("group", opt.Group)
	v.Add("Tenant", opt.Tenant)
	return http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s?%s", uri, v.Encode()), nil)
}

// GetConfig gets the specific config from ACM.
// https://help.aliyun.com/document_detail/64131.html
func (s *Service) GetConfig(ak *AccessKey, opt ConfigOption) ([]byte, error) {
	req, err := s.getConfigRequest(opt, configRequest)
	if err != nil {
		return nil, err
	}
	return Request(s.Cl, req, ak, opt, Timestamp())
}

func listenRequest(uri string, opt ConfigOption) (*http.Request, error) {
	data := "Probe-Modify-Request=" + opt.Join()
	// fmt.Println(data)
	req, err := http.NewRequest(http.MethodPost, "http://"+uri, bytes.NewBufferString(data))
	if err != nil {
		return req, err
	}
	req.Header.Add("longPullingTimeout", "30000") // TODO: make it an argument
	return req, err
}

// ListenConfig listens for the changed config(s).
// https://help.aliyun.com/document_detail/64132.html
func (s *Service) ListenConfig(ak *AccessKey, opt ConfigOption) (resp []byte, err error) {
	req, err := s.getConfigRequest(opt, listenRequest)
	if err != nil {
		return nil, err
	}
	return Request(s.Listener, req, ak, opt, Timestamp())
}

// ParseListenResponse gets the changed config(s).
func ParseListenResponse(resp []byte) (changed []ConfigOption) {
	changed = make([]ConfigOption, 0)
	cs := strings.Split(string(resp), "%01")
	for _, c := range cs {
		if len(c) <= 0 {
			continue
		}
		// dataId^2group^2tenant^1
		ps := strings.Split(c, "%02")
		if len(ps) < 2 {
			continue
		}
		opt := ConfigOption{
			DataID: ps[0],
			Group:  ps[1],
		}
		if len(ps) > 2 {
			opt.Tenant = ps[2]
		}
		changed = append(changed, opt)
	}
	return changed
}
