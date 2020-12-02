// Package timelimit is a CoreDNS plugin that restricts certain dns resolutions during certain timeframes
//
// It serves as a crude way of implementing parental controls for network devices
package timelimit

import (
	"context"
	"time"

	"github.com/coredns/coredns/plugin"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
)

// Define log to be a logger with the plugin name in it. This way we can just use log.Info and
// friends to log.
var log = clog.NewWithPlugin("timelimit")

// TimeLimit is an example plugin to show how to write a plugin.
type TimeLimit struct {
	Next plugin.Handler
}

// ServeDNS implements the plugin.Handler interface. This method gets called when example is used
// in a Server.
func (tl TimeLimit) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	// Here we check if its a restricted dns name and if so, only allow during certain times.
	// If not, forward the reqeust.

	// Debug log that we've have seen the query. This will only be shown when the debug plugin is loaded.
	log.Debug("Received response")

	state := request.Request{W: w, Req: r}

	if state.QName() == "youtube.com" {
		hour := time.Now().Hour()
		if hour >= 13 && hour < 14 {
			return dns.RcodeServerFailure, nil
		} else if hour >= 15 && hour < 17 {
			return dns.RcodeServerFailure, nil
		}
	}

	// Call next plugin (if any).
	return plugin.NextOrFailure(tl.Name(), tl.Next, ctx, w, r)
}

// Name implements the Handler interface.
func (tl TimeLimit) Name() string { return "timelimit" }
