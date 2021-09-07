package fcmSdk

import (
	"context"
	"encoding/json"
	"fmt"
	"go-example/fcmSdk/internal"
	"net/http"
	"strings"
)

const (
	iidEndpoint    = "https://iid.googleapis.com/iid/v1"
	iidSubscribe   = "batchAdd"
	iidUnsubscribe = "batchRemove"
)

// TopicManagementResponse is the result produced by topic management operations.
//
// TopicManagementResponse provides an overview of how many input tokens were successfully handled,
// and how many failed. In case of failures, the Errors list provides specific details concerning
// each error.
type TopicManagementResponse struct {
	SuccessCount int
	FailureCount int
	Errors       []*ErrorInfo
}

func NewTopicManagementResponse(resp *IidResponse) *TopicManagementResponse {
	tmr := &TopicManagementResponse{}
	for idx, res := range resp.Results {
		if len(res) == 0 {
			tmr.SuccessCount++
		} else {
			tmr.FailureCount++
			reason := res["error"].(string)
			tmr.Errors = append(tmr.Errors, &ErrorInfo{
				Index:  idx,
				Reason: reason,
			})
		}
	}
	return tmr
}

type IidClient struct {
	IidEndpoint string
	HttpClient  *internal.HTTPClient
}

func NewIIDClient(hc *http.Client) *IidClient {
	client := internal.WithDefaultRetryConfig(hc)
	client.CreateErrFn = handleIIDError
	client.Opts = []internal.HTTPOption{internal.WithHeader("access_token_auth", "true")}
	return &IidClient{
		IidEndpoint: iidEndpoint,
		HttpClient:  client,
	}
}

// SubscribeToTopic subscribes a list of registration tokens to a topic.
//
// The tokens list must not be empty, and have at most 1000 tokens.
func (c *IidClient) SubscribeToTopic(ctx context.Context, tokens []string, topic string) (*TopicManagementResponse, error) {
	req := &IidRequest{
		Topic:  topic,
		Tokens: tokens,
		op:     iidSubscribe,
	}
	return c.makeTopicManagementRequest(ctx, req)
}

// UnsubscribeFromTopic unsubscribes a list of registration tokens from a topic.
//
// The tokens list must not be empty, and have at most 1000 tokens.
func (c *IidClient) UnsubscribeFromTopic(ctx context.Context, tokens []string, topic string) (*TopicManagementResponse, error) {
	req := &IidRequest{
		Topic:  topic,
		Tokens: tokens,
		op:     iidUnsubscribe,
	}
	return c.makeTopicManagementRequest(ctx, req)
}

type IidRequest struct {
	Topic  string   `json:"to"`
	Tokens []string `json:"registration_tokens"`
	op     string
}

type IidResponse struct {
	Results []map[string]interface{} `json:"results"`
}

type iidErrorResponse struct {
	Error string `json:"error"`
}

func (c *IidClient) makeTopicManagementRequest(ctx context.Context, req *IidRequest) (*TopicManagementResponse, error) {
	if len(req.Tokens) == 0 {
		return nil, fmt.Errorf("no tokens specified")
	}
	if len(req.Tokens) > 1000 {
		return nil, fmt.Errorf("tokens list must not contain more than 1000 items")
	}
	for _, token := range req.Tokens {
		if token == "" {
			return nil, fmt.Errorf("tokens list must not contain empty strings")
		}
	}

	if req.Topic == "" {
		return nil, fmt.Errorf("topic name not specified")
	}
	if !topicNamePattern.MatchString(req.Topic) {
		return nil, fmt.Errorf("invalid topic name: %q", req.Topic)
	}

	if !strings.HasPrefix(req.Topic, "/topics/") {
		req.Topic = "/topics/" + req.Topic
	}

	request := &internal.Request{
		Method: http.MethodPost,
		URL:    fmt.Sprintf("%s:%s", c.IidEndpoint, req.op),
		Body:   internal.NewJSONEntity(req),
	}
	var result IidResponse
	if _, err := c.HttpClient.DoAndUnmarshal(ctx, request, &result); err != nil {
		return nil, err
	}

	return NewTopicManagementResponse(&result), nil
}

func handleIIDError(resp *internal.Response) error {
	base := internal.NewFirebaseError(resp)
	var ie iidErrorResponse
	json.Unmarshal(resp.Body, &ie) // ignore any json parse errors at this level
	if ie.Error != "" {
		base.String = fmt.Sprintf("error while calling the iid service: %s", ie.Error)
	}

	return base
}
