package main

import (
	"net/url"
	"regexp"
	"strings"
	"sync"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
)

// CleanURLPlugin strips well-known tracking query parameters (utm_*, fbclid,
// gclid, ...) from URLs contained in posted messages before they are stored.
type CleanURLPlugin struct {
	plugin.MattermostPlugin

	configLock    sync.RWMutex
	configuration *configuration
}

type configuration struct {
	ExtraTrackingParams string
}

// urlRegexp matches bare http(s) URLs while stopping at characters that
// commonly delimit a URL in Markdown (e.g. the closing ")" of a link, or
// wrapping quotes/angle brackets), so it doesn't swallow trailing syntax.
var urlRegexp = regexp.MustCompile(`https?://[^\s<>()\[\]"']+`)

// defaultTrackingParams lists query parameter names (lower-case) that are
// removed regardless of any admin-configured extra list.
var defaultTrackingParams = map[string]bool{
	"gclid":          true,
	"gclsrc":         true,
	"dclid":          true,
	"fbclid":         true,
	"msclkid":        true,
	"mc_cid":         true,
	"mc_eid":         true,
	"igshid":         true,
	"igsh":           true,
	"yclid":          true,
	"twclid":         true,
	"ttclid":         true,
	"vero_id":        true,
	"mkt_tok":        true,
	"_hsenc":         true,
	"_hsmi":          true,
	"ref_src":        true,
	"ref_url":        true,
	"spm":            true,
	"scid":           true,
	"si":             true,
	"s_kwcid":        true,
	"wt.mc_id":       true,
	"gbraid":         true,
	"wbraid":         true,
	"gad_source":     true,
	"gad_campaignid": true,
	"dm_cam":         true,
	"dm_grp":         true,
	"dm_ad":          true,
	"dm_kw":          true,
	"dm_net":         true,
	"att_gcid":       true,
	"att_gbid":       true,
	"att_wbid":       true,
}

func (p *CleanURLPlugin) OnConfigurationChange() error {
	var configuration configuration
	if err := p.API.LoadPluginConfiguration(&configuration); err != nil {
		return err
	}

	p.configLock.Lock()
	defer p.configLock.Unlock()
	p.configuration = &configuration
	return nil
}

func (p *CleanURLPlugin) extraTrackingParams() map[string]bool {
	p.configLock.RLock()
	config := p.configuration
	p.configLock.RUnlock()

	if config == nil || strings.TrimSpace(config.ExtraTrackingParams) == "" {
		return nil
	}

	extra := make(map[string]bool)
	for _, raw := range strings.Split(config.ExtraTrackingParams, ",") {
		name := strings.ToLower(strings.TrimSpace(raw))
		if name != "" {
			extra[name] = true
		}
	}
	return extra
}

func (p *CleanURLPlugin) MessageWillBePosted(c *plugin.Context, post *model.Post) (*model.Post, string) {
	return p.cleanPost(post), ""
}

func (p *CleanURLPlugin) MessageWillBeUpdated(c *plugin.Context, newPost, oldPost *model.Post) (*model.Post, string) {
	return p.cleanPost(newPost), ""
}

func (p *CleanURLPlugin) cleanPost(post *model.Post) *model.Post {
	newMessage, changed := cleanMessage(post.Message, p.extraTrackingParams())
	if !changed {
		return post
	}

	clone := post.Clone()
	clone.Message = newMessage
	return clone
}

// cleanMessage rewrites every URL found in message, stripping tracking
// query parameters, and reports whether anything was changed.
func cleanMessage(message string, extra map[string]bool) (string, bool) {
	changedAny := false
	result := urlRegexp.ReplaceAllStringFunc(message, func(match string) string {
		cleaned, changed := cleanURL(match, extra)
		if changed {
			changedAny = true
			return cleaned
		}
		return match
	})
	return result, changedAny
}

// cleanURL removes tracking query parameters from a single URL while
// preserving the order and encoding of the parameters that are kept. It
// only returns changed=true (and a rebuilt URL) when something was
// actually removed, so untouched URLs are never needlessly re-encoded.
func cleanURL(raw string, extra map[string]bool) (string, bool) {
	u, err := url.Parse(raw)
	if err != nil || u.RawQuery == "" {
		return raw, false
	}

	pairs := strings.Split(u.RawQuery, "&")
	kept := make([]string, 0, len(pairs))
	changed := false

	for _, pair := range pairs {
		if pair == "" {
			continue
		}

		key := pair
		if idx := strings.IndexByte(pair, '='); idx >= 0 {
			key = pair[:idx]
		}
		if decoded, err := url.QueryUnescape(key); err == nil {
			key = decoded
		}

		if isTrackingParam(strings.ToLower(key), extra) {
			changed = true
			continue
		}
		kept = append(kept, pair)
	}

	if !changed {
		return raw, false
	}

	u.RawQuery = strings.Join(kept, "&")
	u.ForceQuery = false
	return u.String(), true
}

func isTrackingParam(key string, extra map[string]bool) bool {
	if strings.HasPrefix(key, "utm_") {
		return true
	}
	if defaultTrackingParams[key] {
		return true
	}
	return extra != nil && extra[key]
}

func main() {
	plugin.ClientMain(&CleanURLPlugin{})
}
