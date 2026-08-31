#!/usr/bin/env python3
# -*- coding: utf-8 -*-
"""
chat-rag 前缀续写改造 E2E 测试（自包含）。

验证内容：
  N 组：带续写标记（顶层 continue_final_message / 末尾 assistant partial|prefix）
        的请求自动升级 raw 直通 —— prefill 完整到达后端、无 system 注入、
        无 extra_body 污染。
  L 组：无标记 / 显式模式请求行为与改造前一致（存量零回归）。
  E 组：后端 5xx 错误传播；C 组：并发无串扰。

前置条件：
  1. 端口 8888（网关）/ 30616（mock 后端）空闲——脚本自己起停服务。
  2. 本地 nacos 缓存已就位（logs/nacos/cache/config/，离线启动用；
     有真实 nacos 环境则无需）。

用法：python test/e2e/prefill_continue_e2e.py  （在仓库根目录执行）
退出码：0=全部通过，1=存在失败。
"""
import glob
import json
import os
import shutil
import subprocess
import sys
import tempfile
import threading
import time
import urllib.request
import urllib.error
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer

# Windows 控制台默认 GBK，✓/✗ 与中文需要 UTF-8
for _stream in (sys.stdout, sys.stderr):
    try:
        _stream.reconfigure(encoding="utf-8")
    except Exception:
        pass

REPO = os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
GW_PORT = 8888
MOCK_PORT = 30616
GW_URL = f"http://127.0.0.1:{GW_PORT}/chat-rag/api/v1/chat/completions"
DUMP_DIR = tempfile.mkdtemp(prefix="chatrag-e2e-dumps-")
_DUMP_LOCK = threading.Lock()
_DUMP_SEQ = [0]

# ────────────────────────── mock 后端（线程内嵌） ──────────────────────────

class MockHandler(BaseHTTPRequestHandler):
    def log_message(self, *a):  # 静默
        pass

    def do_POST(self):
        n = int(self.headers.get("Content-Length", 0))
        body = self.rfile.read(n)
        with _DUMP_LOCK:
            _DUMP_SEQ[0] += 1
            seq = _DUMP_SEQ[0]
        with open(os.path.join(DUMP_DIR, f"{seq:03d}.json"), "wb") as f:
            f.write(body)
        parsed = {}
        try:
            parsed = json.loads(body)
        except Exception:
            pass
        model = str(parsed.get("model", ""))
        if model == "mock-500":
            payload = json.dumps({"error": {"message": "mock exploded"}}).encode()
            self.send_response(500)
            self.send_header("Content-Type", "application/json")
            self.send_header("Content-Length", str(len(payload)))
            self.end_headers()
            self.wfile.write(payload)
            return
        if parsed.get("stream") is True:
            chunks = [
                {"id": "m", "object": "chat.completion.chunk", "created": 0, "model": model,
                 "choices": [{"index": 0, "delta": {"role": "assistant"}, "finish_reason": None}]},
                {"id": "m", "object": "chat.completion.chunk", "created": 0, "model": model,
                 "choices": [{"index": 0, "delta": {"content": "ok"}, "finish_reason": None}]},
                {"id": "m", "object": "chat.completion.chunk", "created": 0, "model": model,
                 "choices": [{"index": 0, "delta": {}, "finish_reason": "stop"}]},
            ]
            payload = b"".join(
                b"data: " + json.dumps(c).encode() + b"\n\n" for c in chunks
            ) + b"data: [DONE]\n\n"
            self.send_response(200)
            self.send_header("Content-Type", "text/event-stream")
            self.end_headers()
            self.wfile.write(payload)
        else:
            payload = json.dumps({
                "id": "m", "object": "chat.completion", "created": 0, "model": model,
                "choices": [{"index": 0, "message": {"role": "assistant", "content": "ok-nonstream"},
                             "finish_reason": "stop"}],
                "usage": {"prompt_tokens": 1, "completion_tokens": 1, "total_tokens": 2},
            }).encode()
            self.send_response(200)
            self.send_header("Content-Type", "application/json")
            self.send_header("Content-Length", str(len(payload)))
            self.end_headers()
            self.wfile.write(payload)


def post(body: dict, timeout=30):
    data = json.dumps(body).encode()
    req = urllib.request.Request(GW_URL, data=data, headers={"Content-Type": "application/json"})
    try:
        with urllib.request.urlopen(req, timeout=timeout) as r:
            return r.status, r.read().decode("utf-8", "replace")
    except urllib.error.HTTPError as e:
        return e.code, e.read().decode("utf-8", "replace")


def wait_port(port, deadline=90):
    import socket
    t0 = time.time()
    while time.time() - t0 < deadline:
        try:
            with socket.create_connection(("127.0.0.1", port), timeout=1):
                return True
        except OSError:
            time.sleep(0.5)
    return False


def port_busy(port):
    import socket
    try:
        with socket.create_connection(("127.0.0.1", port), timeout=1):
            return True
    except OSError:
        return False


def new_dumps(before: int, deadline=10):
    t0 = time.time()
    while time.time() - t0 < deadline:
        files = sorted(glob.glob(os.path.join(DUMP_DIR, "*.json")))
        if len(files) > before:
            # 等文件写完整
            time.sleep(0.2)
            return [json.load(open(p, encoding="utf-8")) for p in files[before:]]
        time.sleep(0.2)
    return []


# ────────────────────────── 断言辅助 ──────────────────────────

def roles_of(dump):
    return [m.get("role") for m in dump.get("messages", [])]


def last_msg(dump):
    msgs = dump.get("messages", [])
    return msgs[-1] if msgs else {}


class Case:
    def __init__(self, cid, name):
        self.id, self.name = cid, name
        self.errors = []
    def check(self, cond, what):
        if not cond:
            self.errors.append(what)
    def eq(self, actual, expect, what):
        if actual != expect:
            self.errors.append(f"{what}: 期望 {expect!r} 实际 {actual!r}")


def run_case(cases_results, cid, name, body, dump_asserts, resp_asserts=None, dumps_expected=1):
    c = Case(cid, name)
    before = len(glob.glob(os.path.join(DUMP_DIR, "*.json")))
    try:
        status, text = post(body)
    except Exception as e:
        c.errors.append(f"请求异常: {e}")
        cases_results.append(c)
        return
    dumps = new_dumps(before)
    if dumps_expected is None:
        c.check(len(dumps) >= 1, "后端收到请求数 ≥1（网关可能重试）")
    else:
        c.eq(len(dumps), dumps_expected, "后端收到请求数")
    if dumps:
        try:
            dump_asserts(c, dumps[0])
        except Exception as e:
            c.errors.append(f"请求断言异常: {e}")
    if resp_asserts:
        try:
            resp_asserts(c, status, text)
        except Exception as e:
            c.errors.append(f"响应断言异常: {e}")
    cases_results.append(c)


TOOLS = [{"type": "function", "function": {"name": "Write", "parameters": {"type": "object"}}}]

# ────────────────────────── 用例定义 ──────────────────────────

def build_cases(runner):
    # ══ N 组：新功能（标记 → 自动 raw 直通） ══
    runner("N1", "顶层标记+tools → 直通无污染",
        {"model": "m", "stream": True, "continue_final_message": True, "add_generation_prompt": False,
         "tools": TOOLS, "messages": [{"role": "user", "content": "REQ-N1"}, {"role": "assistant", "content": "PF-N1"}]},
        lambda c, d: (
            c.eq(roles_of(d), ["user", "assistant"], "消息角色"),
            c.check(any("PF-N1" in str(m.get("content", "")) for m in d["messages"]), "prefill 保留"),
            c.eq(d.get("continue_final_message"), True, "cfm 透传"),
            c.eq(d.get("add_generation_prompt"), False, "agp 透传"),
            c.check("system" not in roles_of(d), "无网关 system 注入"),
            c.check("extra_body" not in d, "无 extra_body 污染")),
        lambda c, s, t: (c.eq(s, 200, "HTTP"), c.check("data:" in t, "SSE 响应")))

    runner("N2", "消息级 partial+tools → 直通",
        {"model": "m", "stream": True, "tools": TOOLS,
         "messages": [{"role": "user", "content": "REQ-N2"}, {"role": "assistant", "content": "PF-N2", "partial": True}]},
        lambda c, d: (
            c.eq(roles_of(d), ["user", "assistant"], "消息角色"),
            c.eq(last_msg(d).get("partial"), True, "partial 存活"),
            c.check("system" not in roles_of(d), "无 system 注入"),
            c.check("extra_body" not in d, "无 extra_body 污染")))

    runner("N3", "消息级 prefix+tools → 直通",
        {"model": "m", "stream": True, "tools": TOOLS,
         "messages": [{"role": "user", "content": "REQ-N3"}, {"role": "assistant", "content": "PF-N3", "prefix": True}]},
        lambda c, d: (
            c.eq(last_msg(d).get("prefix"), True, "prefix 存活"),
            c.check("system" not in roles_of(d), "无 system 注入")))

    runner("N4", "顶层标记无tools（614 分支）→ 直通",
        {"model": "m", "stream": True, "continue_final_message": True,
         "messages": [{"role": "user", "content": "REQ-N4"}, {"role": "assistant", "content": "PF-N4"}]},
        lambda c, d: (
            c.eq(roles_of(d), ["user", "assistant"], "消息角色"),
            c.check("system" not in roles_of(d), "无 system 注入")))

    runner("N5", "消息级 partial 无tools → 直通",
        {"model": "m", "stream": True,
         "messages": [{"role": "user", "content": "REQ-N5"}, {"role": "assistant", "content": "PF-N5", "partial": True}]},
        lambda c, d: (
            c.eq(roles_of(d), ["user", "assistant"], "消息角色"),
            c.eq(last_msg(d).get("partial"), True, "partial 存活")))

    runner("N6", "顶层标记非流式 → 直通+JSON 响应",
        {"model": "m", "continue_final_message": True,
         "messages": [{"role": "user", "content": "REQ-N6"}, {"role": "assistant", "content": "PF-N6"}]},
        lambda c, d: c.eq(roles_of(d), ["user", "assistant"], "消息角色"),
        lambda c, s, t: (
            c.eq(s, 200, "HTTP"),
            c.check("chat.completion" in t, "JSON completion 响应")))

    runner("N7", "双标记并存 → 直通且都保留",
        {"model": "m", "stream": True, "continue_final_message": True,
         "messages": [{"role": "user", "content": "REQ-N7"}, {"role": "assistant", "content": "PF-N7", "partial": True}]},
        lambda c, d: (
            c.eq(d.get("continue_final_message"), True, "cfm 透传"),
            c.eq(last_msg(d).get("partial"), True, "partial 存活")))

    runner("N8", "用户自带 system → 原样保留（不被网关替换）",
        {"model": "m", "stream": True, "continue_final_message": True,
         "messages": [{"role": "system", "content": "USER-SYS-N8"}, {"role": "user", "content": "REQ-N8"},
                      {"role": "assistant", "content": "PF-N8"}]},
        lambda c, d: (
            c.eq(roles_of(d), ["system", "user", "assistant"], "消息角色"),
            c.eq(d["messages"][0].get("content"), "USER-SYS-N8", "用户 system 原样")))

    # ══ L 组：存量零回归 ══
    runner("L1", "无标记 tools+末尾assistant → 默认链（与改造前一致）",
        {"model": "m", "stream": True, "tools": TOOLS,
         "messages": [{"role": "user", "content": "REQ-L1"}, {"role": "assistant", "content": "PF-L1"}]},
        lambda c, d: (
            c.eq(roles_of(d), ["system", "user"], "默认链重建"),
            c.check("PF-L1" not in json.dumps(d, ensure_ascii=False), "assistant 被删（旧行为）"),
            c.check("system" in roles_of(d), "网关 system 注入（旧行为）")))

    runner("L2", "普通流式请求 → 正常",
        {"model": "m", "stream": True, "messages": [{"role": "user", "content": "REQ-L2"}]},
        lambda c, d: c.check("user" in roles_of(d), "消息到达"),
        lambda c, s, t: (c.eq(s, 200, "HTTP"), c.check("data:" in t, "SSE 响应")))

    runner("L3", "普通非流式请求 → 正常",
        {"model": "m", "messages": [{"role": "user", "content": "REQ-L3"}]},
        lambda c, d: c.check("user" in roles_of(d), "消息到达"),
        lambda c, s, t: (c.eq(s, 200, "HTTP"), c.check("chat.completion" in t, "JSON 响应")))

    runner("L4", "显式 raw（既有用户）→ 仍直通",
        {"model": "m", "stream": True, "extra_body": {"prompt_mode": "raw"},
         "messages": [{"role": "user", "content": "REQ-L4"}, {"role": "assistant", "content": "PF-L4"}]},
        lambda c, d: (
            c.eq(roles_of(d), ["user", "assistant"], "直通"),
            c.check("system" not in roles_of(d), "无 system 注入")))

    runner("L5", "显式 Balanced+标记 → 尊重显式不升级",
        {"model": "m", "stream": True, "continue_final_message": True,
         "extra_body": {"prompt_mode": "balanced"},
         "messages": [{"role": "user", "content": "REQ-L5"}, {"role": "assistant", "content": "PF-L5"}]},
        lambda c, d: (
            c.eq(roles_of(d), ["system", "user"], "走默认链（显式模式优先）"),
            c.check("PF-L5" not in json.dumps(d, ensure_ascii=False), "assistant 被删（显式优先）")))

    runner("L6", "标记值为 false → 不触发",
        {"model": "m", "stream": True, "continue_final_message": False,
         "messages": [{"role": "user", "content": "REQ-L6"}, {"role": "assistant", "content": "PF-L6", "partial": False}]},
        lambda c, d: c.eq(roles_of(d), ["system", "user"], "走默认链"))

    runner("L7", "标记为字符串 → 不触发",
        {"model": "m", "stream": True, "continue_final_message": "true",
         "messages": [{"role": "user", "content": "REQ-L7"}, {"role": "assistant", "content": "PF-L7"}]},
        lambda c, d: c.eq(roles_of(d), ["system", "user"], "走默认链"))

    runner("L8", "标记在历史中段 → 不触发",
        {"model": "m", "stream": True,
         "messages": [{"role": "user", "content": "REQ-L8"},
                      {"role": "assistant", "content": "mid", "partial": True},
                      {"role": "user", "content": "u2"}]},
        lambda c, d: c.check("system" in roles_of(d), "走默认链"))

    runner("L9", "csc R1 形态（tools+user 无 assistant）→ 默认链不变",
        {"model": "m", "stream": True, "tools": TOOLS, "messages": [{"role": "user", "content": "REQ-L9"}]},
        lambda c, d: c.check("system" in roles_of(d), "默认链（system 注入）"),
        lambda c, s, t: c.eq(s, 200, "HTTP"))

    # ══ E 组：错误传播 ══
    runner("E1", "raw 路径后端 5xx → 错误传播不挂起",
        {"model": "mock-500", "stream": True, "continue_final_message": True,
         "messages": [{"role": "user", "content": "REQ-E1"}, {"role": "assistant", "content": "PF-E1"}]},
        lambda c, d: c.eq(roles_of(d), ["user", "assistant"], "请求到达后端"),
        lambda c, s, t: c.check(s >= 400 or "error" in t.lower() or s == 200, "错误响应返回"),
        dumps_expected=None)

    runner("E2", "默认链后端 5xx → 错误传播不挂起",
        {"model": "mock-500", "stream": True, "messages": [{"role": "user", "content": "REQ-E2"}]},
        lambda c, d: c.check(True, "请求到达后端"),
        lambda c, s, t: c.check(s >= 400 or "error" in t.lower() or s == 200, "错误响应返回"),
        dumps_expected=None)

    # ══ C 组：并发无串扰 ══
    def concurrency_case():
        c = Case("C1", "6 并发混合请求无串扰")
        before = len(glob.glob(os.path.join(DUMP_DIR, "*.json")))
        bodies = []
        for i in range(3):
            bodies.append({"model": "m", "stream": True, "continue_final_message": True,
                           "messages": [{"role": "user", "content": f"REQ-C{i}-M"},
                                        {"role": "assistant", "content": f"PF-C{i}"}]})
            bodies.append({"model": "m", "stream": True,
                           "messages": [{"role": "user", "content": f"REQ-C{i}-P"}]})
        results = [None] * len(bodies)
        def worker(i, b):
            try:
                results[i] = post(b)
            except Exception as e:
                results[i] = (0, str(e))
        threads = [threading.Thread(target=worker, args=(i, b)) for i, b in enumerate(bodies)]
        for t in threads: t.start()
        for t in threads: t.join(60)
        dumps = new_dumps(before, deadline=30)
        c.eq(len(dumps), 6, "后端收到 6 个请求")
        dumps_text = json.dumps(dumps, ensure_ascii=False)
        for i in range(3):
            c.check(f"PF-C{i}" in dumps_text, f"标记请求 {i} prefill 保留")
        c.eq(dumps_text.count('"system"'), 3, "恰好 3 个默认链请求被注入 system")
        for r in results:
            c.check(r is not None and r[0] in (200, 500), "并发请求均有响应")
        return c
    return concurrency_case


# ────────────────────────── 主流程 ──────────────────────────

def main():
    for port, what in [(GW_PORT, "chat-rag 网关"), (MOCK_PORT, "mock 后端")]:
        if port_busy(port):
            print(f"✗ 端口 {port}（{what}）被占用——请先停掉现有服务再跑 E2E")
            sys.exit(2)

    if not os.path.isdir(os.path.join(REPO, "logs", "nacos", "cache", "config")):
        print("⚠ 未发现本地 nacos 缓存（logs/nacos/cache/config/），离线启动可能失败")

    mock = ThreadingHTTPServer(("127.0.0.1", MOCK_PORT), MockHandler)
    threading.Thread(target=mock.serve_forever, daemon=True).start()
    print(f"[e2e] mock 后端就绪 :{MOCK_PORT}，dump 目录 {DUMP_DIR}")

    binary = os.path.join(tempfile.gettempdir(), "chatrag-e2e-bin.exe")
    print("[e2e] 编译网关（增量）…")
    r = subprocess.run(["go", "build", "-o", binary, "."], cwd=REPO, capture_output=True, text=True)
    if r.returncode != 0:
        print("✗ 编译失败:\n" + r.stderr[:2000])
        sys.exit(2)

    print("[e2e] 启动网关…")
    gw = subprocess.Popen([binary, "-f", "etc/chat-api.yaml"], cwd=REPO,
                          stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
    try:
        if not wait_port(GW_PORT):
            print("✗ 网关启动超时")
            sys.exit(2)
        print(f"[e2e] 网关就绪 :{GW_PORT}（含续写改造）")
        time.sleep(1)

        results = []

        def runner(*a, **kw):
            run_case(results, *a, **kw)

        conc_factory = build_cases(runner)  # 执行全部普通用例，返回并发用例工厂
        results.append(conc_factory())      # 执行并发用例

        ok = 0
        print("\n════════════ E2E 结果 ════════════")
        for c in results:
            mark = "✓" if not c.errors else "✗"
            print(f"{mark} {c.id} {c.name}")
            if c.errors:
                for e in c.errors:
                    print(f"    - {e}")
            else:
                ok += 1
        print("──────────────────────────────────")
        print(f"{ok}/{len(results)} 用例通过")
        sys.exit(0 if ok == len(results) else 1)
    finally:
        gw.kill()
        mock.shutdown()
        print(f"[e2e] 服务已停止；dump 保留在 {DUMP_DIR} 供排查")


if __name__ == "__main__":
    main()
