package v2b

import (
	"encoding/json"
	"time"

	"github.com/avast/retry-go/v4"
	"github.com/go-resty/resty/v2"
	"github.com/wyx2685/UniProxy/common/balance"
)

var (
	clients *balance.List[*resty.Client]
	etag    string
)

func Init(b string, url []string, auth string) {
	cs := make([]*resty.Client, len(url))
	for i, u := range url {
		cs[i] = resty.New().
			SetTimeout(time.Second*40).
			SetQueryParam("auth_data", auth).
			SetBaseURL(u).
			SetRetryCount(3).
			SetRetryWaitTime(3 * time.Second)
	}
	clients = balance.New[*resty.Client](b, cs)
}

type ServerFetchRsp struct {
	Data []ServerInfo `json:"data"`
}

type ServerInfo struct {
	Id             int         `json:"id"`
	Name           string      `json:"name"`
	Host           string      `json:"host"`
	Port           int         `json:"port"`
	Network        string      `json:"network"`
	Type           string      `json:"type"`
	Cipher         string      `json:"cipher"`
	Tls            int         `json:"tls"`
	Flow           string      `json:"flow"`
	TlsSettings    struct {
		AllowInsecure string `json:"allow_insecure"`
		Fingerprint   string `json:"fingerprint"`
		PublicKey     string `json:"public_key"`
		ServerName    string `json:"server_name"`
		ShortId       string `json:"short_id"`
	} `json:"tls_settings"`
	NetworkSettings struct {
		Path       string      `json:"path"`
		Headers    interface{} `json:"headers"`
		ServerName string      `json:"serverName"`
	} `json:"networkSettings"`
	CreatedAt      int	   `json:"created_at"`
	AllowInsecure  int         `json:"insecure"`
	Allow_Insecure int         `json:"allow_insecure"`
	LastCheckAt    interface{} `json:"last_check_at"`
	Tags           interface{} `json:"tags"`
	UpMbps        int         `json:"up_mbps"`
	ServerName    string      `json:"server_name"`
	ServerKey     string      `json:"server_key"`
	DownMbps      int         `json:"down_mbps"`
	HysteriaVersion int        `json:"version"`
	Hy2Obfs       string      `json:"obfs"`
	Hy2ObfsPassword string     `json:"obfs_password"`
}

func GetServers() ([]ServerInfo, error) {
	var r *resty.Response
	err := retry.Do(func() error {
		c := clients.Next()
		rsp, err := c.R().
			SetHeader("If-None-Match", etag).
			Get("api/v1/user/server/fetch")
		if err != nil {
			return err
		}
		if rsp.StatusCode() == 304 {
			return nil
		}
		etag = rsp.Header().Get("ETag")
		if rsp.StatusCode() != 200 {
			return nil
		}
		r = rsp
		return nil
	}, retry.Attempts(3))
	if err != nil {
		return nil, err
	}
	if r.StatusCode() == 304 {
		return nil, nil
	}
	rsp := &ServerFetchRsp{}
	err = json.Unmarshal(r.Body(), rsp)
	if err != nil {
		return nil, err
	}
	if len(rsp.Data) == 0 {
		// 返回一个默认的节点
		defaultServer := ServerInfo{
			Id:          0,
			Name:        "默认节点",
			Host:        "default.host",
			Port:        80,
			Network:     "default",
			Type:        "default",
			Cipher:      "none",
			Tls:         0,
			Flow:        "default",
			// 填充其他字段为默认值
		}
		return []ServerInfo{defaultServer}, nil
	}
	return rsp.Data, nil
}
