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
            events.append({"notification": self.path.split("/")[1]})
            self.reply({"error": "fixture notification failure"}, status=503)
            return
        body = json.loads(raw or b"{}")
        account = self.headers.get("Authorization", "").removeprefix("Bearer fixture-")
        model = body.get("model")
        events.append({"account": account, "model": model, "stream": body.get("stream")})
        counts[account, model] += 1
        if model == "inspection-block":
            time.sleep(30)
        answer = "20" if account == "wrong" or (account == "mixed" and counts[account, model] != 2) else "21"
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
