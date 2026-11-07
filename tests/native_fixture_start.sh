#!/bin/sh
set -eu
if [ "${1:-}" = "-version" ]; then
    exec /app/sub2api-bin "$@"
fi
grep -q "chatgpt.com" /etc/hosts 2>/dev/null || printf "127.0.0.1 chatgpt.com api.openai.com\n" >> /etc/hosts 2>/dev/null || true
if ! command -v postgres >/dev/null 2>&1 || [ ! -f /usr/share/postgresql18/extension/pg_trgm.control ]; then
    apk add --no-cache postgresql18 postgresql18-contrib redis python3 su-exec
fi
mkdir -p /var/lib/postgresql/data /run/postgresql
chown postgres:postgres /var/lib/postgresql/data /run/postgresql
if [ ! -f /var/lib/postgresql/data/PG_VERSION ]; then
    su-exec postgres initdb -D /var/lib/postgresql/data --auth-local=trust --auth-host=trust
fi
su-exec postgres postgres -D /var/lib/postgresql/data -k /run/postgresql &
postgres_pid=$!
redis-server --bind 127.0.0.1 --save '' --appendonly no &
redis_pid=$!
python3 /app/fixture-supplier.py &
supplier_pid=$!
app_pid=
stop() {
    trap - TERM INT
    [ -z "$app_pid" ] || kill -TERM "$app_pid" 2>/dev/null || true
    [ -z "$app_pid" ] || wait "$app_pid" || true
    kill -TERM "$supplier_pid" "$redis_pid" "$postgres_pid" 2>/dev/null || true
    wait || true
}
trap 'stop; exit 0' TERM INT
until pg_isready -h 127.0.0.1 -U postgres -d postgres -q; do sleep 1; done
if ! su-exec postgres psql -tAc "SELECT 1 FROM pg_roles WHERE rolname='sub2api'" | grep -q 1; then
    su-exec postgres createuser --superuser sub2api
fi
if ! su-exec postgres psql -tAc "SELECT 1 FROM pg_database WHERE datname='sub2api'" | grep -q 1; then
    su-exec postgres createdb -O sub2api sub2api
fi
while [ ! -f /app/fixture.ready ]; do sleep 1; done
export DATABASE_HOST=127.0.0.1 DATABASE_PORT=5432 REDIS_HOST=127.0.0.1 REDIS_PORT=6379 TZ=UTC
/app/sub2api-bin "$@" &
app_pid=$!
wait "$app_pid" || result=$?
stop
exit "${result:-0}"
