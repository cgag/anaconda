package anaconda

import (
	"net/url"
	"time"
)

// The thing you really want
type RateLimitStatus struct {
	RateLimitContext RateLimitContext
	Resources        Resources
}

type RateLimitContext struct {
	AccessToken string `json:"access_token"`
}

type EndpointStatus struct {
	Limit     int       `json:"limit"`
	Remaining int       `json:"remaining"`
	Reset     time.Time `json:"reset"` // need to parse from unix time ala RateLimitCheck()_
}

type Resources struct {
	Users          Resource `json:"users"`
	Friendships    Resource `json:"friendships"`
	Followers      Resource `json:"followers"`
	Account        Resource `json:"account"`
	Media          Resource `json:"media"`
	DirectMessages Resource `json:"direct_messages"`
	Mutes          Resource `json:"mutes"`
	Favorites      Resource `json:"favorites"`
	Statuses       Resource `json:"statuses"`
	Help           Resource `json:"help"`
	Search         Resource `json:"search"`
}

type Resource map[Endpoint]EndpointStatus

func (c *TwitterApi) GetRateLimitStatus() (RateLimitStatus, error) {
	raw, err := c.GetRateLimitStatusRaw()

	if err != nil {
		return RateLimitStatus{}, err
	}

	rateLimitStatus := parseTimes(raw)

	return rateLimitStatus, nil
}

func (c *TwitterApi) GetRateLimitStatusRaw() (RateLimitStatusRawTimes, error) {
	var raw RateLimitStatusRawTimes
	respChan := make(chan response)
	rawQuery := query{
		BaseUrl + "/application/rate_limit_status.json",
		url.Values{},
		&raw,
		_GET,
		respChan,
	}
	c.queryQueue <- rawQuery
	rawResp := <-respChan
	if rawResp.err != nil {
		return RateLimitStatusRawTimes{}, rawResp.err
	}

	return raw, nil
}

// RateLimitStatus with unparsed times
type RateLimitStatusRawTimes struct {
	RateLimitContext RateLimitContextRaw `json:"rate_limit_context"`
	Resources        ResourcesRaw        `json:"resources"`
}

type EndpointStatusRaw struct {
	Limit     int `json:"limit"`
	Remaining int `json:"remaining"`
	Reset     int `json:"reset"` // need to parse from unix time ala RateLimitCheck()_
}

type RateLimitContextRaw struct {
	AccessToken string `json:"access_token"`
}

type Endpoint string
type ResourceRaw map[Endpoint]EndpointStatusRaw

type ResourcesRaw struct {
	Users          ResourceRaw `json:"users"`
	Friendships    ResourceRaw `json:"friendships"`
	Followers      ResourceRaw `json:"followers"`
	Account        ResourceRaw `json:"account"`
	DirectMessages ResourceRaw `json:"direct_messages"`
	Mutes          ResourceRaw `json:"mutes"`
	Favorites      ResourceRaw `json:"favorites"`
	Statuses       ResourceRaw `json:"statuses"`
	Help           ResourceRaw `json:"help"`
	Search         ResourceRaw `json:"search"`
}

func parseTimes(raw RateLimitStatusRawTimes) RateLimitStatus {
	var parsedStatus RateLimitStatus

	rr := raw.Resources

	parsedStatus.RateLimitContext = RateLimitContext{
		raw.RateLimitContext.AccessToken,
	}

	parsedStatus.Resources = Resources{
		Users:          parseRawResource(rr.Users),
		Friendships:    parseRawResource(rr.Friendships),
		Followers:      parseRawResource(rr.Followers),
		Account:        parseRawResource(rr.Account),
		DirectMessages: parseRawResource(rr.DirectMessages),
		Mutes:          parseRawResource(rr.Mutes),
		Favorites:      parseRawResource(rr.Favorites),
		Statuses:       parseRawResource(rr.Statuses),
		Help:           parseRawResource(rr.Help),
		Search:         parseRawResource(rr.Search),
	}

	return parsedStatus
}

func parseRawResource(resourceRaw ResourceRaw) Resource {
	resource := make(Resource)

	for endpoint, endpointStatusRaw := range resourceRaw {
		resource[endpoint] = EndpointStatus{
			Limit:     endpointStatusRaw.Limit,
			Remaining: endpointStatusRaw.Remaining,
			Reset:     time.Unix(int64(endpointStatusRaw.Reset), 0),
		}
	}

	return resource
}
