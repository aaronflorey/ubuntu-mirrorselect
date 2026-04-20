#!/usr/bin/env bash

set -euo pipefail

repo_root="$(git rev-parse --show-toplevel)"
tmp_dir="$(mktemp -d)"
server_log="$tmp_dir/server.log"
port_file="$tmp_dir/port"
tag="v9.9.9"
version="${tag#v}"
asset_name="mirrorselect_${version}_linux_amd64.tar.gz"
release_root="$tmp_dir/repo/releases"
download_dir="$release_root/download/$tag"
tag_dir="$release_root/tag/$tag"

cleanup() {
  if [[ -n "${server_pid:-}" ]]; then
    kill "$server_pid" >/dev/null 2>&1 || true
    wait "$server_pid" >/dev/null 2>&1 || true
  fi
  rm -rf "$tmp_dir"
}

trap cleanup EXIT

mkdir -p "$download_dir" "$tag_dir" "$tmp_dir/bin"

printf '#!/usr/bin/env sh\necho mirrorselect smoke test\n' >"$tmp_dir/bin/mirrorselect"
chmod +x "$tmp_dir/bin/mirrorselect"
tar -C "$tmp_dir/bin" -czf "$download_dir/$asset_name" mirrorselect
printf 'ok\n' >"$tag_dir/index.html"

python3 -u - <<'PY' "$tmp_dir" "$port_file" >"$server_log" 2>&1 &
import http.server
import socketserver
import sys

root = sys.argv[1]
port_file = sys.argv[2]


class Handler(http.server.SimpleHTTPRequestHandler):
    def __init__(self, *args, **kwargs):
        super().__init__(*args, directory=root, **kwargs)

    def do_GET(self):
        if self.path == "/repo/releases/latest":
            self.send_response(302)
            self.send_header("Location", "/repo/releases/tag/v9.9.9")
            self.end_headers()
            return
        return super().do_GET()


with socketserver.TCPServer(("127.0.0.1", 0), Handler) as httpd:
    with open(port_file, "w", encoding="utf-8") as handle:
        handle.write(str(httpd.server_address[1]))
    httpd.serve_forever()
PY
server_pid="$!"

for _ in $(seq 1 50); do
  if [[ -s "$port_file" ]]; then
    break
  fi
  sleep 0.1
done

if [[ ! -s "$port_file" ]]; then
  printf 'smoke test server did not start\n' >&2
  exit 1
fi

port="$(<"$port_file")"
test_home="$tmp_dir/home"
mkdir -p "$test_home"

HOME="$test_home" \
MIRRORSELECT_INSTALLER_BASE_URL="http://127.0.0.1:$port/repo" \
bash "$repo_root/install.sh"

test -x "$test_home/.local/bin/mirrorselect"
"$test_home/.local/bin/mirrorselect" >/dev/null
