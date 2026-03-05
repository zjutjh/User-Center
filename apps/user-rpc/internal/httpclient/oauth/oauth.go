package oauth

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"github.com/zjutjh/User-Center/common/constants"
	"github.com/zjutjh/User-Center/common/errorsx"
	"golang.org/x/sync/errgroup"
)

type Client struct {
	ctx context.Context
}

var (
	httpClient *resty.Client
	once       sync.Once
)

const maxRedirect = 10

func createClient() *resty.Client {
	client := resty.New()
	client.
		SetRedirectPolicy(resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		})).
		SetRetryCount(3).
		SetTimeout(5 * time.Second)
	client.SetCookieJar(nil)
	return client
}

func New(ctx context.Context) *Client {
	once.Do(func() {
		httpClient = createClient()
	})
	return &Client{ctx: ctx}
}

func (c *Client) R() *resty.Request {
	return httpClient.R().SetContext(c.ctx)
}

func (c *Client) Login(username string, password string) ([]*http.Cookie, error) {
	var (
		initCookies      = make([]*http.Cookie, 0)
		publicKeyCookies = make([]*http.Cookie, 0)
		execution        string
		encPwd           string
		g                errgroup.Group
	)

	g.Go(func() error {
		resp, err := c.R().Get(constants.OAuthLoginURL)
		if err != nil {
			return err
		}

		initCookies = resp.Cookies()
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body()))
		if err != nil {
			return err
		}
		execution = doc.Find("input[type=hidden][name=execution]").AttrOr("value", "")
		return nil
	})

	g.Go(func() error {
		var publicKey publicKeyData
		resp, err := c.R().
			SetResult(&publicKey).
			Get(constants.OAuthLoginPublicKeyURL)
		if err != nil {
			return err
		}

		publicKeyCookies = resp.Cookies()
		encPwd, err = getEncryptedPassword(&publicKey, password)
		return err
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	resp, err := c.R().
		SetCookies(append(initCookies, publicKeyCookies...)).
		SetFormData(map[string]string{
			"username":  username,
			"password":  encPwd,
			"execution": execution,
			"_eventId":  "submit",
		}).
		Post(constants.OAuthLoginURL)
	if err != nil && !errors.Is(err, resty.ErrAutoRedirectDisabled) {
		return nil, err
	}
	if err := checkLogin(resp); err != nil {
		return nil, err
	}

	for _, cookie := range resp.Cookies() {
		if cookie.Name == "iPlanetDirectoryPro" {
			return []*http.Cookie{cookie}, nil
		}
	}

	return nil, errorsx.ErrUnknown
}

func checkLogin(resp *resty.Response) error {
	if resp.StatusCode() == http.StatusFound && resp.Header().Get("Location") == constants.PersonalCenterURL {
		return nil
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body()))
	if err != nil {
		return err
	}
	if doc.Find("title").Text() == "修改密码" {
		return errorsx.ErrPasswordNeedEdited
	}

	switch strings.TrimSpace(doc.Find("#msg").Text()) {
	case constants.WrongPasswordMsg, constants.WrongAccountMsg:
		return errorsx.ErrWrongAccountOrPassword
	case constants.NotActivatedMsg:
		return errorsx.ErrNotActivated
	default:
		return errorsx.ErrUnknown
	}
}

func (c *Client) SSOLogin(cookies []*http.Cookie, serviceURL string) ([]*http.Cookie, error) {
	location := constants.OAuthLoginURL + "?service=" + url.QueryEscape(serviceURL)
	finalCookies := append([]*http.Cookie{}, cookies...)

	for redirectCount := 0; redirectCount < maxRedirect; redirectCount++ {
		resp, err := c.R().
			SetCookies(finalCookies).
			Get(location)
		if err != nil && !errors.Is(err, resty.ErrAutoRedirectDisabled) {
			return nil, fmt.Errorf("%w: %v", errorsx.ErrInvalidCookie, err)
		}
		next := resp.Header().Get("Location")
		if resp.StatusCode() != http.StatusFound || next == "" {
			break
		}
		finalCookies = append(finalCookies, resp.Cookies()...)
		location = resolveRedirect(location, next)
	}

	return finalCookies, nil
}

func resolveRedirect(current string, next string) string {
	baseURL, err := url.Parse(current)
	if err != nil {
		return next
	}
	nextURL, err := url.Parse(next)
	if err != nil {
		return next
	}
	return baseURL.ResolveReference(nextURL).String()
}
