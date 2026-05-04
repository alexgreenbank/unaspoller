package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
)

type UNAS struct {
	username  string
	password  string
	target    string
	baseURL   string
	client    *http.Client
	csrf      string
	cookies   []*http.Cookie // FIXME cookie jar
	m         *metrics
	probeFile *os.File
	c         *config
	apidefs   map[string]*UNASDriveAPIDef
}

func NewUNAS(ip string, c *config) *UNAS {
	ret := UNAS{}
	ret.target = ip
	ret.baseURL = fmt.Sprintf("https://%s", ip)
	ret.cookies = make([]*http.Cookie, 0)
	ret.apidefs = make(map[string]*UNASDriveAPIDef)
	ret.c = c
	ret.client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: !c.optVerifySSL},
		},
	}
	return &ret
}

func (u *UNAS) SetUsername(un string) {
	u.username = un
}

func (u *UNAS) SetPassword(pw string) {
	u.password = pw
}
