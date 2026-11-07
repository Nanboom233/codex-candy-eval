#!/usr/bin/env python3
"""Serve synthetic supplier responses inside the reusable test container."""
import collections
import http.server
import json
import time

counts = collections.Counter()
events = []


class Supplier(http.server.BaseHTTPRequestHandler):
    def log_message(self, *_args):
        pass

    def reply(self, payload, content_type="application/json", status=200):
        raw = payload.encode() if isinstance(payload, str) else json.dumps(payload).encode()
        self.send_response(status)
        self.send_header("Content-Type", content_type)
        self.send_header("Content-Length", str(len(raw)))
        self.end_headers()
        try:
            self.wfile.write(raw)
        except (BrokenPipeError, ConnectionResetError):
            pass

    def do_GET(self):
        if self.path == "/fixture/state":
            self.reply({"events": events})
            return
        self.reply({"object": "list", "data": [{"id": "inspection-literal", "object": "model"}]})

    def do_POST(self):
        raw = self.rfile.read(int(self.headers.get("Content-Length", "0")))
        if self.path.startswith(("/ntfy/", "/bark/")):
            path_parts = self.path.split("?", 1)[0].split("/")
            notification_status = 200 if len(path_parts) > 2 and path_parts[2] == "success" else 503
            events.append({"notification": path_parts[1], "status": notification_status,
                           "body": raw.decode(errors="replace")})
            payload = {"ok": True} if notification_status == 200 else {"error": "fixture notification failure"}
            self.reply(payload, status=notification_status)
            return
        body = json.loads(raw or b"{}")
        raw_text = raw.decode(errors="replace")
        account = (self.headers.get("Authorization", "").removeprefix("Bearer fixture-")
                   or self.headers.get("x-api-key", "").removeprefix("fixture-")
                   or self.headers.get("x-goog-api-key", "").removeprefix("fixture-"))
        model = body.get("model") or self.path.split("/models/")[-1].split(":")[0]
        events.append({"account": account, "model": model, "stream": body.get("stream"), "path": self.path})
        counts[account, model] += 1
        if model == "inspection-block":
            time.sleep(30)
        if "355" in raw_text:
            answer = ",".join(str((i * 37 + 13) % 355 + 1) for i in range(200))
        elif "Reply with OK." in raw_text:
            if account == "wrong" or "restricted" in model:
                self.reply({"error": {"code": "subscription_sharing_user_not_eligible",
                                      "message": "Your account is not eligible to use subscription sharing"}}, status=403)
                return
            answer = "OK"
        elif "鹈鹕" in raw_text and "```svg" in raw_text:
            answer = "```svg\n<svg xmlns=\"http://www.w3.org/2000/svg\" viewBox=\"0 0 100 100\"><circle cx=\"50\" cy=\"50\" r=\"40\" fill=\"#38bdf8\"/><text x=\"50\" y=\"55\" text-anchor=\"middle\">21</text></svg>\n```"
        elif "鹈鹕" in raw_text:
            answer = "```html\n<!DOCTYPE html><html><body><svg viewBox=\"0 0 100 100\"><circle cx=\"50\" cy=\"50\" r=\"40\"/></svg></body></html>\n```"
        else:
            is_wrong = (
                account == "wrong"
                or "wrong" in model
                or (model == "inspection-degrade-a" and account == "a-only")
                or (account == "mixed" and counts[account, model] != 2)
            )
            answer = "20" if is_wrong else "21"
        if "/v1/messages" in self.path:
            sse = (
                "event: message_start\n"
                'data: {"type":"message_start","message":{"id":"msg_fixture","type":"message","role":"assistant","usage":{"input_tokens":8,"output_tokens":1}}}\n\n'
                "event: content_block_delta\n"
                + "data: " + json.dumps({"type": "content_block_delta", "index": 0, "delta": {"type": "text_delta", "text": answer}}) + "\n\n"
                + "event: message_delta\n"
                'data: {"type":"message_delta","delta":{"stop_reason":"end_turn"},"usage":{"output_tokens":4}}\n\n'
                "event: message_stop\n"
                'data: {"type":"message_stop"}\n\n'
            )
            self.reply(sse, "text/event-stream")
            return
        if ":streamGenerateContent" in self.path:
            chunk = {
                "candidates": [{"content": {"role": "model", "parts": [{"text": answer}]}, "finishReason": "STOP"}],
                "usageMetadata": {"promptTokenCount": 9, "candidatesTokenCount": 3, "thoughtsTokenCount": 1},
            }
            self.reply("data: " + json.dumps(chunk) + "\n\n", "text/event-stream")
            return
        if self.path.endswith("/chat/completions"):
            if body.get("stream"):
                chunks = [{"id": "fixture", "object": "chat.completion.chunk", "model": model,
                           "choices": [{"index": 0, "delta": {"role": "assistant", "content": answer}, "finish_reason": None}]},
                          {"id": "fixture", "object": "chat.completion.chunk", "model": model,
                           "choices": [{"index": 0, "delta": {}, "finish_reason": "stop"}],
                           "usage": {"prompt_tokens": 5, "completion_tokens": 1, "total_tokens": 6}}]
                self.reply("".join("data: " + json.dumps(chunk) + "\n\n" for chunk in chunks) + "data: [DONE]\n\n", "text/event-stream")
            else:
                self.reply({"id": "fixture", "object": "chat.completion", "model": model,
                            "choices": [{"index": 0, "message": {"role": "assistant", "content": answer}, "finish_reason": "stop"}]})
            return
        response = {"id": "fixture-response", "object": "response", "status": "completed",
                    "output": [{"type": "message", "role": "assistant", "content": [{"type": "output_text", "text": answer}]}],
                    "usage": {"input_tokens": 5, "output_tokens": 1}}
        if body.get("stream"):
            self.reply("data: " + json.dumps({"type": "response.completed", "response": response}) + "\n\n", "text/event-stream")
        else:
            self.reply(response)


if __name__ == "__main__":
    http.server.ThreadingHTTPServer(("127.0.0.1", 18081), Supplier).serve_forever()
