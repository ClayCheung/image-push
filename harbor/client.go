package harbor

import (
	"crypto/tls"
	"fmt"
	"github.com/caicloud/cargo-admin/pkg/utils/matcher"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/caicloud/nirvana/log"

	"errors"
)

type Client struct {
	config   *Config
	baseURL  string
	client   *http.Client
	coockies []*http.Cookie
}

func NewClient(host, username, password string) (*Client, error) {
	return newClient(&Config{host, username, password})
}

func newClient(conf *Config) (*Client, error) {
	baseURL := strings.TrimRight(conf.Host, "/")
	if !strings.Contains(baseURL, "://") {
		baseURL = "http://" + baseURL
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	cookies, err := LoginAndGetCookies(client, conf)
	if err != nil {
		log.Errorf("login harbor: %s error: %v during background", conf.Host, err)
		return nil, err
	}
	log.Infof("harbor %s cookies has been refreshed", conf.Host)

	return &Client{
		config:   conf,
		baseURL:  baseURL,
		client:   client,
		coockies: cookies,
	}, nil
}

// do creates request and authorizes it if authorizer is not nil
func (c *Client) do(method, relativePath string, body io.Reader) (*http.Response, error) {
	url := c.baseURL + relativePath
	log.Infof("%s %s", method, url)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	if body != nil || method == http.MethodPost || method == http.MethodPut {
		req.Header.Set("Content-Type", "application/json")
	}
	for i := range c.coockies {
		req.AddCookie(c.coockies[i])
	}

	resp, err := c.client.Do(req)
	if err != nil {
		log.Errorf("unexpected error: %v", err)
		return nil, err
	}

	if resp.StatusCode/100 == 5 || resp.StatusCode == 401 {
		b, err := ioutil.ReadAll(resp.Body)
		defer func() {
			_ = resp.Body.Close()
		}()
		if err != nil {
			return nil, errors.ErrorUnknownInternal.Error(err)
		}

		log.Errorf("unexpected %d error from harbor: %s", resp.StatusCode, b)
		log.Errorf("need to refresh harbor: %s 's cookies now! refreshCookies error: %v", c.config.Host, c.RefreshCookies())
		return nil, errors.ErrorUnknownInternal.Error(fmt.Sprintf("harbor internal error: %s", b))
	}
	return resp, nil
}

func (c *Client) RefreshCookies() error {
	cookies, err := LoginAndGetCookies(c.client, c.config)
	if err != nil {
		log.Errorf("refresh harbor: %s 's cookies error: %v", c.config.Host, err)
		return err
	}
	c.coockies = cookies
	return nil
}

func (c *Client) GetConfig() *Config {
	return c.config
}


func LoginAndGetCookies(client *http.Client, conf *Config) ([]*http.Cookie, error) {
	url := LoginURL(conf.Host, conf.Username, conf.Password)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		msg := matcher.MaskPwd(err.Error())
		log.Error(msg)
		return nil, errors.ErrorUnknownInternal.Error(msg)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.ErrorUnknownInternal.Error(err)
	}

	if resp.StatusCode != 200 {
		log.Errorf("login harbor: %s error: %s", conf.Host, b)
		return nil, errors.ErrorUnknownInternal.Error(fmt.Sprintf("%s", b))
	}

	// If status code is 200 and no cookies set, it's not a valid harbor. For example, 1.1.1.1
	if string(b) != "" || len(resp.Cookies()) == 0 {
		return nil, errors.ErrorUnknownInternal.Error(fmt.Sprintf("%s is not a valid harbor", conf.Host))
	}

	return resp.Cookies(), nil
}