package main

/*
	This code is written by a programmer and ethical hacker ..
    What this code is doing? Well, If you'll read this code, then you will understand that.
*/

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"time"
)

type Config struct {
	URL 	 string
	UserName string
	Password string
}

type Recipients struct {
	Id int `json:"id"`
	Name string `json:"name"`
}

type Message struct {
	Subject string `json:"subject"`
	Body string `json:"body"`
	DraftID int `json:"draftId"`
	PreviousMessageId int `json:"previousMessageId"`
	Attachments []string `json:"attachments"`
	DraftType string `json:"draftType"`
	Recipients []Recipients `json:"recipients"`
}

type Mail struct {
	Message Message `json:"message"`	
}

func main(){
	// Setting config for post request
	config := Config{
		URL: 		"https://my.e-klase.lv/?v=15", // Login panel
		UserName: 	"{YOUR_USERNAME}", 			   // Your username
		Password: 	"{YOUR_PASSWORD}",			   // Your password
	}

	// Oh, wait, I need also set data for the request
	data := url.Values{}
	data.Set("UserName", config.UserName)
	data.Set("Password", config.Password)
	data.Set("fake_pass", "")
	data.Set("cmdLogIn", "")

	// Creating POST reqeust to DOS the system
	req, err := http.NewRequest(http.MethodPost, config.URL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		log.Fatal(err)
	}

	// Setting headers for the request
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.140 Safari/537.36 OPR/69.2.2399.217")

	// Setting proxies to be anonymous, hahaha.
	proxy, err := url.Parse("socks5://localhost:9050")
	if err != nil {
		log.Fatal(err)
	}
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxy),
	}

	// Creating client who will send this request.. and cookies?
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{
		Jar: jar,
		Timeout: time.Second * 10,
		Transport: transport,
	}

	// Sending this request to the server.. yo, #HACKTHEWORLD
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// response from the server
	if resp.StatusCode == 200 {
		// Reading body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		// Let's check if we are in the right place where we want to be in.
		if !(strings.Contains(string(body), "/Family/Home")){
			err := fmt.Errorf("log in into %s failed, please check your config and try again", config.URL)
			fmt.Println(err)
			os.Exit(1)	
		}

		// shall we continue this game? p.s. there we are already in the system
		// so, let's send some messages? ... I need to find what we need for that.
		mail := Mail {
			Message: Message {
				Subject: "Penetration testing against eklase's server",
				Body: "Hello, this is not a ddos attack, this is penetration testing to test if servers can stand against multiple reqeusts at once",
				DraftType: "mdt_new",
				Recipients: []Recipients {
					{
						Name: "Anonymous ethical hacker (administrator), 3.kt",
						Id: 2157814,
					},
				},
			},
		}

		// Making this message to json, why? because we need json and all.
		json, err := json.Marshal(mail)
		if err != nil {
			log.Fatal(err)
		}

		// Creating new request to test this shit
		req, err := http.NewRequest(http.MethodPost, "https://my.e-klase.lv/api/family/mail/send", bytes.NewReader(json))
		if err != nil {
			log.Fatal(err)
		}

		// Okay, It's time for our sweet headers
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_3; en-MT) AppleWebKit/536.14.27 (KHTML, like Gecko) Version/13.0.2 Safari/536.14.27")
		req.Header.Set("Content-Type", "application/json;charset=utf-8")

		// Maybe now we can run send this request and see if that works? (of course, that not works)
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		// response from mail API
		if resp.StatusCode == 200 {
			// Reading body
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%s\n", body)
		}
	}
}
