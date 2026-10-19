//go:build unit

package service

import (
	"bufio"
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/config"
	"github.com/Wei-Shaw/sub2api/internal/pkg/apicompat"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestCodexEvalNativeSender(t *testing.T) {
	for _, chat := range []bool{false, true} {
		name := "responses"
		if chat {
			name = "chat-completions"
		}
		t.Run(name, func(t *testing.T) {
			account := &Account{ID: 17, Name: "selected", Platform: PlatformOpenAI, Type: AccountTypeAPIKey, Concurrency: 1, Credentials: map[string]any{"api_key": "fixture-upstream-key", "base_url": "https://api.openai.com", "model_mapping": map[string]any{"requested": "mapped"}}, Extra: map[string]any{"openai_responses_supported": !chat}}
			repo := &openAIAccountTestRepo{mockAccountRepoForGemini: mockAccountRepoForGemini{accountsByID: map[int64]*Account{17: account}}}
			events := `data: {"type":"response.completed","response":{"status":"completed","output":[{"type":"message","content":[{"type":"output_text","text":"21"}]}],"usage":{"input_tokens":9,"output_tokens":2}}}` + "\n\n"
			if chat {
				events = "data: {\"id\":\"fixture\",\"model\":\"requested\",\"choices\":[{\"index\":0,\"delta\":{\"role\":\"assistant\",\"content\":\"21\"},\"finish_reason\":null}]}\n\ndata: {\"id\":\"fixture\",\"choices\":[{\"index\":0,\"delta\":{},\"finish_reason\":\"stop\"}],\"usage\":{\"prompt_tokens\":9,\"completion_tokens\":2,\"total_tokens\":11}}\n\ndata: [DONE]\n\n"
			}
			upstream := &queuedHTTPUpstream{responses: []*http.Response{newJSONResponse(200, events)}}
			cfg := &config.Config{}
			svc := &AccountTestService{accountRepo: repo, httpUpstream: upstream, cfg: cfg}
			body := []byte(`{"model":"requested","stream":true,"input":[{"role":"user","content":[{"type":"input_text","text":"complete prompt"}]}],"tools":[{"type":"function","name":"fixture_tool","parameters":{"type":"object"}}],"reasoning":{"effort":"none"},"client_metadata":{"turn_id":"fixture"}}`)
			resp, err := svc.SendCodexEvaluation(context.Background(), 17, body, http.Header{"Authorization": {"Bearer administrator-session"}, "Session-Id": {"fixture-session"}})
			require.NoError(t, err)
			defer resp.Body.Close()
			raw, err := io.ReadAll(resp.Body)
			require.NoError(t, err)
			require.Contains(t, string(raw), "response.completed")
			require.Contains(t, string(raw), "21")
			require.Contains(t, string(raw), "input_tokens\":9")
			require.Len(t, upstream.requests, 1)
			req := upstream.requests[0]
			sent, err := io.ReadAll(req.Body)
			require.NoError(t, err)
			require.Equal(t, "Bearer fixture-upstream-key", req.Header.Get("Authorization"))
			require.Equal(t, "fixture-session", req.Header.Get("Session-Id"))
			require.Equal(t, "requested", gjson.GetBytes(sent, "model").String())
			require.Contains(t, string(sent), "complete prompt")
			require.Contains(t, string(sent), "fixture_tool")
			if chat {
				require.True(t, strings.HasSuffix(req.URL.Path, "/chat/completions"))
				require.Equal(t, "none", gjson.GetBytes(sent, "reasoning_effort").String())
			} else {
				require.True(t, strings.HasSuffix(req.URL.Path, "/responses"))
				require.Equal(t, "fixture", gjson.GetBytes(sent, "client_metadata.turn_id").String())
			}
			require.Zero(t, repo.clearedErrorID)
			_, err = svc.SendCodexEvaluation(context.Background(), 18, body, nil)
			require.ErrorContains(t, err, "account not found")
			require.Len(t, upstream.requests, 1)
			upstream.responses = []*http.Response{newJSONResponse(401, `{"error":{"message":"fixture authentication failure"}}`)}
			resp, err = svc.SendCodexEvaluation(context.Background(), 17, body, nil)
			require.NoError(t, err)
			raw, err = io.ReadAll(resp.Body)
			resp.Body.Close()
			require.NoError(t, err)
			require.Equal(t, 401, resp.StatusCode)
			require.Contains(t, string(raw), "fixture authentication failure")
			require.Equal(t, int64(17), repo.setErrorID)
			for _, status := range []int{401, 429} {
				upstream.responses = []*http.Response{{StatusCode: status, Header: http.Header{}, Body: codexErrorBody{}}}
				resp, err = svc.SendCodexEvaluation(context.Background(), 17, body, nil)
				require.ErrorContains(t, err, "fixture body read")
				require.NotNil(t, resp)
				require.Equal(t, status, resp.StatusCode)
				resp.Body.Close()
			}
		})
	}
}

type codexErrorBody struct{}

func (codexErrorBody) Read([]byte) (int, error) { return 0, errors.New("fixture body read") }
func (codexErrorBody) Close() error             { return nil }

func TestCodexEvalModelsBypassesAccountProjection(t *testing.T) {
	gateway := newCodexModelsAPIKeyTestService(&codexModelsHTTPUpstreamStub{do: func(_ *http.Request, _ string, _ int64, _ int) (*http.Response, error) {
		return ordinaryModelsUpstreamResponse(`{"data":[{"id":"upstream-target"},{"id":"unconfigured-model"}]}`), nil
	}})
	account := newCodexModelsAPIKeyTestAccount("https://models.example/v1")
	account.ID = 17
	account.Credentials["model_mapping"] = map[string]any{"configured-alias": "upstream-target"}
	repo := &openAIAccountTestRepo{mockAccountRepoForGemini: mockAccountRepoForGemini{accountsByID: map[int64]*Account{17: account}}}
	svc := &AccountTestService{accountRepo: repo, openaiGatewayService: gateway}
	models, err := svc.CodexEvalModels(context.Background(), 17)
	require.NoError(t, err)
	require.Len(t, models, 2)
	require.Equal(t, "upstream-target", models[0].ID)
	require.Equal(t, "unconfigured-model", models[1].ID)
	ordinary, err := svc.FetchOpenAIAccountModels(context.Background(), account)
	require.NoError(t, err)
	require.Len(t, ordinary, 1)
	require.Equal(t, "configured-alias", ordinary[0].ID)
}

func TestCodexEvalChatStreamFailures(t *testing.T) {
	for _, failure := range []struct {
		wire, message string
	}{
		{`data: {"error":{"message":"actual upstream failure"}}` + "\n\n", "actual upstream failure"},
		{"", "stream ended before finish_reason"},
	} {
		t.Run(failure.message, func(t *testing.T) {
			wire := "data: {\"choices\":[],\"usage\":{\"prompt_tokens\":9,\"completion_tokens\":2,\"total_tokens\":11}}\n\n" + failure.wire
			body := io.NopCloser(strings.NewReader(wire))
			stream := &codexChatStream{body: body, reader: bufio.NewReader(body), state: apicompat.NewChatCompletionsToResponsesStreamState("fixture")}
			raw, err := io.ReadAll(stream)
			require.NoError(t, err)
			require.Contains(t, string(raw), `"type":"error"`)
			require.Contains(t, string(raw), failure.message)
			require.Contains(t, string(raw), `"input_tokens":9`)
			require.NotContains(t, string(raw), "response.completed")
		})
	}
}
