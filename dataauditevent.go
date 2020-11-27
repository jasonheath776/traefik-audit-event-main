// Package traefik_audit_event can trigger Audit events when some pattern matches
package traefik_audit_event

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
)


// Config the plugin configuration.
type Config struct {
	Url          string `yaml:"Url"`
	IncludeRequest bool    `yaml:"IncludeRequest"`
	IncludeResponse bool `yaml:"IncludeResponse"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		Url: "::6000",
		IncludeRequest: true,
		IncludeResponse: true,
	}
}

// AuditEvent plugin
type AuditEvent struct {
	next        http.Handler
	Req string
	Res string
	Config *Config
}

// New created a new plugin
func New(ctx context.Context, next http.Handler, config *Config, request string, response string) (http.Handler, error) {
	if config.Url == "" {
		return nil, fmt.Errorf("traefik-audit-event: you need to specity your url endpoint")
	}
	return &AuditEvent{
		Req:   request,
		Res:  response,
		Config: config,
		next:     next,

	}, nil
}

func (a *AuditEvent) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	resr := httptest.NewRecorder()
	reqr := httptest.NewRecorder()

	if a.Config.IncludeResponse {
		a.next.ServeHTTP(resr, req)
		_, _ = rw.Write(resr.Body.Bytes())
	}

	if a.Config.IncludeRequest {
		_ = req.Write(reqr.Body)
	}

	if resr.Code == http.StatusOK	{
		SendEvent(a)
	}
}

// GenerateEventPayload generates the JSON payload required by the audit API
func GenerateEventPayload(a *AuditEvent) *bytes.Buffer {
	return bytes.NewBuffer([]byte(`{"req":"` + a.Req + `","res":"` + a.Res + `"}"`))
}

// SendEvent rpc call fire and forget
func SendEvent(a *AuditEvent){
	fmt.Printf("%s",a.Req)
	fmt.Printf("%s", a.Res)
	//client, err := rpc.DialHTTP("tcp", a.Config.Url)
	//
	//if err != nil {
	//	log.Print(`dialing: ` + a.Config.Url, err)
	//}
	//
	//_= client.Go("WriteAudit", GenerateEventPayload(a), nil, nil)
}

