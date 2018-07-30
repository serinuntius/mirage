package main

import (
	"fmt"
	"math/rand"
	"net/http"

	"gopkg.in/acidlemon/rocket.v2"
)

type WebApi struct {
	rocket.WebApp
	cfg *Config
}

func NewWebApi(cfg *Config) *WebApi {
	app := &WebApi{}
	app.Init()
	app.cfg = cfg

	view := &rocket.View{
		BasicTemplates: []string{cfg.Storage.HtmlDir + "/layout.html"},
	}

	app.AddRoute("/", app.List, view)
	app.AddRoute("/launcher", app.Launcher, view)
	app.AddRoute("/launch", app.Launch, view)
	app.AddRoute("/terminate", app.Terminate, view)
	app.AddRoute("/api/list", app.ApiList, view)
	app.AddRoute("/api/launch", app.ApiLaunch, view)
	app.AddRoute("/api/terminate", app.ApiTerminate, view)

	app.BuildRouter()

	return app
}

func (api *WebApi) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	api.Handler(w, req)
}

func (api *WebApi) List(c rocket.CtxData) {
	info, err := app.Docker.List()
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}
	value := rocket.RenderVars{
		"info":  info,
		"error": errStr,
	}

	c.Render(api.cfg.Storage.HtmlDir+"/list.html", value)
}

func (api *WebApi) Launcher(c rocket.CtxData) {
	c.Render(api.cfg.Storage.HtmlDir+"/launcher.html", rocket.RenderVars{
		"DefaultImage": api.cfg.Docker.DefaultImage,
		"Parameters":   api.cfg.Parameter,
	})
}

func (api *WebApi) Launch(c rocket.CtxData) {
	result := api.launch(c)
	if result["result"] == "ok" {
		c.Redirect("/")
	} else {
		c.RenderJSON(result)
	}
}

func (api *WebApi) Terminate(c rocket.CtxData) {
	result := api.terminate(c)
	if result["result"] == "ok" {
		c.Redirect("/")
	} else {
		c.RenderJSON(result)
	}
}

func (api *WebApi) ApiList(c rocket.CtxData) {
	info, err := app.Docker.List()
	var status interface{}
	if err != nil {
		status = err.Error()
	} else {
		status = info
	}

	result := rocket.RenderVars{
		"result": status,
	}

	c.RenderJSON(result)
}

func (api *WebApi) ApiLaunch(c rocket.CtxData) {
	result := api.launch(c)

	c.RenderJSON(result)
}

func (api *WebApi) ApiTerminate(c rocket.CtxData) {
	result := api.terminate(c)

	c.RenderJSON(result)
}

func (api *WebApi) launch(c rocket.CtxData) rocket.RenderVars {
	if c.Req().Method != "POST" {
		c.Res().StatusCode = http.StatusMethodNotAllowed
		c.RenderText("you must use POST")
		return rocket.RenderVars{}
	}

	subdomain, _ := c.ParamSingle("subdomain")
	image, _ := c.ParamSingle("image")
	name, _ := c.ParamSingle("name")

	if name == "" {
		name = subdomain + "-" + randomString(5)
	}

	parameter, err := api.loadParameter(c)
	if err != nil {
		result := rocket.RenderVars{
			"result": err.Error(),
		}

		return result
	}

	status := "ok"

	if subdomain == "" || image == "" {
		status = fmt.Sprintf("parameter required: subdomain=%s, image=%s",
			subdomain, image)
	} else {
		err := app.Docker.Launch(subdomain, image, name, parameter)
		if err != nil {
			status = err.Error()
		}
	}

	result := rocket.RenderVars{
		"result": status,
	}

	return result
}

func (api *WebApi) terminate(c rocket.CtxData) rocket.RenderVars {
	if c.Req().Method != "POST" {
		c.Res().StatusCode = http.StatusMethodNotAllowed
		c.RenderText("you must use POST")
		return rocket.RenderVars{}
	}

	status := "ok"

	subdomain, _ := c.ParamSingle("subdomain")
	if subdomain == "" {
		status = fmt.Sprintf("parameter required: subdomain")
	} else {
		err := app.Docker.Terminate(subdomain)

		if err != nil {
			status = err.Error()
		}
	}

	result := rocket.RenderVars{
		"result": status,
	}

	return result
}

func (api *WebApi) loadParameter(c rocket.CtxData) (map[string]string, error) {
	var parameter map[string]string = make(map[string]string)

	for _, v := range api.cfg.Parameter {
		param, _ := c.ParamSingle(v.Name)
		if param == "" && v.Required == true {
			return nil, fmt.Errorf("lack require parameter: %s", v.Name)
		} else if param == "" {
			continue
		}

		if v.Rule != "" {
			if !v.Regexp.MatchString(param) {
				return nil, fmt.Errorf("parameter %s value is rule error", v.Name)
			}
		}

		parameter[v.Name] = param
	}

	return parameter, nil
}

const rsLetters = "0123456789abcdef"

func randomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = rsLetters[rand.Intn(len(rsLetters))]
	}
	return string(b)
}
