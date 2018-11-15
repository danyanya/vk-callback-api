package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

const (
	TYPE_CONFIRM      = "confirmation"
	DEFAULT_CONF_CODE = "confirm"
)

type LeadFormAnswer struct {
	Key      string `json:"key"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type LeadFormObject struct {
	LeadID   string           `json:"lead_id"`
	GroupID  string           `json:"group_id"`
	UserID   string           `json:"user_id"`
	FormID   string           `json:"form_id"`
	FormName string           `json:"form_name"`
	Answers  []LeadFormAnswer `json:"answers"`
}

type Request struct {
	Type    string         `json:"type"`
	GroupID string         `json:"group_id"`
	Object  LeadFormObject `json:"object"`
}

type Response struct {
	Status string
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

	e.Logger.Debug("Confirmation code to use: ", confCode)

	e.Any("/", func(c echo.Context) error {

		req := &Request{}
		if err := c.Bind(req); err != nil {
			e.Logger.Error("Error in bind: ", err.Error())
			return c.JSON(200, Response{Status: "success"})
		}

		if req.Type == TYPE_CONFIRM {
			e.Logger.Debug("Confirmation code response to VK")
			return c.String(200, confCode)
		}

		// print as object
		e.Logger.Info("Parsed request: ", req)

		return c.JSON(200, Response{Status: "success"})
	})

	e.Logger.Info("Start server at port ", servPort)
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", servPort)))
}
