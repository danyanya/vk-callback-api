package main

import (
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

const (
	TYPE_CONFIRM      = "confirmation"
	TYPE_LEAD_NEW     = "lead_forms_new"
	DEFAULT_CONF_CODE = "confirm"

	KEY_EMAIL = "email"
	KEY_PHONE = "phone_number"

	MINDBOX_POINT = "Email"
)

type LeadFormAnswer struct {
	Key      string `json:"key"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type LeadFormObject struct {
	LeadID   int              `json:"lead_id"`
	GroupID  int              `json:"group_id"`
	UserID   int              `json:"user_id"`
	FormID   int              `json:"form_id"`
	FormName string           `json:"form_name"`
	AdId     int              `json:"ad_id"`
	Answers  []LeadFormAnswer `json:"answers"`
}

type Request struct {
	Type    string         `json:"type"`
	GroupID int            `json:"group_id"`
	Object  LeadFormObject `json:"object"`
}

type Response struct {
	Status string
}

type MindBoxSubscription struct {
	Brand          string `xml:"brand"`
	PointOfContact string `xml:"pointOfContact"`
}

type MindBoxCustomer struct {
	MobilePhone   string                `xml:"mobilePhone"`
	Email         string                `xml:"email"`
	Subscriptions []MindBoxSubscription `xml:"subscriptions"`
}

type MindBoxOperation struct {
	XMLName        xml.Name        `xml:"operation"`
	PointOfContact string          `xml:"pointOfContact"`
	Customer       MindBoxCustomer `xml:"customer"`
}

type MindBoxConfig struct {
	URL            string
	Key            string
	Brand          string
	PointOfContact string
}

var mbConfig = MindBoxConfig{}

func sendPostRequest(apiUrl string, apiKey string, payload []byte) string {

	client := &http.Client{}
	body := bytes.NewBuffer(payload)
	r, _ := http.NewRequest(http.MethodPost, apiUrl, body) // URL-encoded payload
	r.Header.Add("Authorization", fmt.Sprintf("Mindbox secretKey=\"%s\"", apiKey))
	r.Header.Add("Content-Type", "application/xml")
	r.Header.Add("Content-Length", fmt.Sprintf("%d", body.Len()))

	resp, err := client.Do(r)

	if err != nil {
		return err.Error()
	}

	return resp.Status
}

func dataSendMindBox(l echo.Logger, request *Request) {
	if request.Type != TYPE_LEAD_NEW {
		return
	}

	var phone string
	var email string

	for _, v := range request.Object.Answers {
		switch v.Key {
		case KEY_EMAIL:
			email = v.Answer
		case KEY_PHONE:
			phone = v.Answer
		}
	}

	l.Info("Phone and email ", phone, " ", email)

	if len(phone)+len(email) == 0 {
		return
	}

	op := MindBoxOperation{
		PointOfContact: mbConfig.PointOfContact,
		Customer: MindBoxCustomer{
			MobilePhone: phone,
			Email:       email,
			Subscriptions: []MindBoxSubscription{
				MindBoxSubscription{
					Brand:          mbConfig.Brand,
					PointOfContact: MINDBOX_POINT,
				},
			},
		},
	}

	payload, _ := xml.MarshalIndent(op, "", "  ")

	fmt.Printf("Payload:\n%s\n", string(payload))
	l.Debug("Payload to send: ", string(payload))
	l.Info("Response code: ",
		sendPostRequest(mbConfig.URL, mbConfig.Key, payload))
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Logger.SetLevel(log.DEBUG)
	e.HideBanner = true

	var confCode string
	var servPort int
	flag.StringVar(&confCode, "c", "confirm", "Confirmation code")
	flag.IntVar(&servPort, "p", 9911, "Port to serve")
	flag.Parse()

	if confCode == DEFAULT_CONF_CODE {
		conf := os.Getenv("CONF_CODE")
		if len(conf) != 0 {
			confCode = conf
		}
	}

	mbConfig.URL = os.Getenv("MINDBOX_URL")
	mbConfig.Key = os.Getenv("MINDBOX_KEY")
	mbConfig.Brand = os.Getenv("MINDBOX_BRAND")
	mbConfig.PointOfContact = os.Getenv("MINDBOX_POINT_OF_CONTACT")

	e.Logger.Debug("Confirmation code to use: ", confCode)

	e.Logger.Info(mbConfig)

	e.Any("/", func(c echo.Context) error {

		req := &Request{}
		if err := c.Bind(req); err != nil {
			e.Logger.Error("Error in bind: ", err.Error())
			return c.String(200, "ok")
		}

		if req.Type == TYPE_CONFIRM {
			e.Logger.Debug("Confirmation code response to VK")
			return c.String(200, confCode)
		}

		e.Logger.Info("Parsed request: ", req)

		// goroutine for read-and-send user data to MindBox
		go dataSendMindBox(e.Logger, req)

		return c.String(200, "ok")
	})

	e.Logger.Info("Start server at port ", servPort)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", servPort)))
}
