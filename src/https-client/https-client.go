package main

import (
	"bytes"
	"crypto/md5"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"time"

	"gopkg.in/yaml.v2"
)

var (
	caCertPath = "C:/Users/LX/Desktop/ca.pem"
	certFile   = "C:/Users/LX/Desktop/client.pem"
	keyFile    = "C:/Users/LX/Desktop/client-key.pem"
)

const (
	ID = "hpe"
)

type Conf struct {
	VppConfig VppConfig `yaml:"vppConfig"`
}

type VppConfig struct {
	Acls []Acl `yaml:"acls"`
}

type Acl struct {
	Name  string `yaml:"name"`
	Rules []Rule `yaml:"rules"`
}

type Rule struct {
	Action string `yaml:"action"`
	Iprule IpRule `yaml:"ipRule"`
}
type IpRule struct {
	Ip IP `yaml:"ip"`
}
type IP struct {
	DestinationNetwork string `yaml:"destinationNetwork"`
	SourceNetwork      string `yaml:"sourceNetwork"`
}

func main() {
	//1. 获取token
	url := "https://10.20.141.243:9191/auth/login"

	var caCrt []byte
	caCrt, err := ioutil.ReadFile(caCertPath)
	pool := x509.NewCertPool()
	if err != nil {
		return
	}
	pool.AppendCertsFromPEM(caCrt)

	var cliCrt tls.Certificate
	cliCrt, err = tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return
	}

	requestBody := `{
		"username":"hpe",
		"password": "$2a$07$sQxSGrZfYA/6InuBUAIrg.vLKl1mNAYJUbfYdj2PrFLXG49gsPVfa"
		}`
	jsonStr := []byte(requestBody)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))

	if err != nil {
		fmt.Println("GetHttp Request Error:", err)
		return
	}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:      pool,
				Certificates: []tls.Certificate{cliCrt},
			},
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	token, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(string(token))

	//2. 发起请求
	url = "https://10.20.141.243:9191/configuration?replace=true"

	var acl []Acl
	acl = append(acl, Acl{Name: "acl_in", Rules: []Rule{Rule{Action: "PERMIT", Iprule: IpRule{Ip: IP{DestinationNetwork: "192.168.50.0/24", SourceNetwork: "0.0.0.0/0"}}}}})
	newConfig := Conf{VppConfig: VppConfig{Acls: acl}}

	reqBody := new(bytes.Buffer)

	yaml.NewEncoder(reqBody).Encode(newConfig)

	req, err = http.NewRequest("PUT", url, reqBody)
	if err != nil {
		panic(err)
	}

	req.Header.Add("Content-Type", "application/yaml")
	//添加token
	req.Header.Add("Authorization", "Bearer "+string(token))

	//添加时间戳
	q := req.URL.Query()
	timestamp := fmt.Sprintf("%d", time.Now().Unix())
	q.Add("ts", timestamp)

	//添加签名
	q.Add("sn", createSign(q, string(token), reqBody.String()))
	req.URL.RawQuery = q.Encode()

	fmt.Println(req.URL.String())
	fmt.Println(reqBody.String())

	resp, err = client.Do(req)

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(string(body))
}

func createSign(params url.Values, token string, body string) string {
	// 自定义 MD5 组合  replace=true&ts=123456 + token + hpe
	str := createEncryptStr(params) + token + ID + body

	s := md5.New()
	s.Write([]byte(str))
	return hex.EncodeToString(s.Sum(nil))
}

func createEncryptStr(params url.Values) string {
	var key []string
	var str = ""
	for k := range params {
		if k != "sn" {
			key = append(key, k)
		}
	}
	sort.Strings(key)
	for i := 0; i < len(key); i++ {
		if i == 0 {
			str = fmt.Sprintf("%v=%v", key[i], params.Get(key[i]))
		} else {
			str = str + fmt.Sprintf("&%v=%v", key[i], params.Get(key[i]))
		}
	}
	return str
}
