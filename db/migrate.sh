#!/bin/bash
# migrate.sh - Apply or rollback database migrations for SecureWallet
# Usage:
#   ./migrate.sh up [version]     - Apply migrations up to version (default: latest)
#   ./migrate.sh down [version]   - Rollback to version (default: previous)
#   ./migrate.sh status           - Show applied migrations

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
MIGRATIONS_DIR="${SCRIPT_DIR}/migrations"

# Database connection from environment (with sensible defaults for dev)
DB_HOST="${DB_HOST:-127.0.0.1}"
DB_PORT="${DB_PORT:-3306}"
DB_USER="${DB_USER:-root}"
DB_PASSWORD="${DB_PASSWORD:-}"
DB_NAME="${DB_NAME:-securewallet_dev}"

MYSQL_CMD="mysql -h ${DB_HOST} -P ${DB_PORT} -u ${DB_USER}"
if [ -n "${DB_PASSWORD}" ]; then
    MYSQL_CMD="${MYSQL_CMD} -p${DB_PASSWORD}"
fi
MYSQL_CMD="${MYSQL_CMD} ${DB_NAME}"

log() {
    echo "[$(date '+%Y-%m-%d %H:%M:%S')] $*"
}

get_applied_versions() {
    ${MYSQL_CMD} -N -e "SELECT version FROM schema_migrations ORDER BY version ASC;" 2>/dev/null || echo ""
}

is_applied() {
    local version=$1
    local result
    result=$(${MYSQL_CMD} -N -e "SELECT COUNT(*) FROM schema_migrations WHERE version = ${version};" 2>/dev/null || echo "0")
    [ "${result}" -gt 0 ]
}

cmd_up() {
    local target_version="${1:-999}"
    log "Applying migrations up to version ${target_version}..."

    local applied=0
    for migration_file in "${MIGRATIONS_DIR}"/V[0-9]*.sql; do
        [ -f "${migration_file}" ] || continue
        [[ "${migration_file}" == *"_rollback.sql" ]] && continue

        local filename
        filename=$(basename "${migration_file}")
        local version
        version=$(echo "${filename}" | grep -oP 'V\K[0-9]+')

        if [ "${version}" -gt "${target_version}" ]; then
            break
        fi

        if is_applied "${version}"; then
            log "  SKIP: ${filename} (already applied)"
            continue
        fi

        log "  APPLY: ${filename}"
        ${MYSQL_CMD} < "${migration_file}"
        applied=$((applied + 1))
    done

    log "Applied ${applied} migration(s)."
}

cmd_down() {
    local target_version="${1:-0}"
    log "Rolling back migrations down to version ${target_version}..."

    local rolled_back=0
    local versions
    versions=$(get_applied_versions)

    for version in $(echo "${versions}" | tac); do
        if [ "${version}" -le "${target_version}" ]; then
            break
        fi

        local rollback_file
        rollback_file=$(find "${MIGRATIONS_DIR}" -name "V${version}_*_rollback.sql" -type f | head -1)

        if [ -z "${rollback_file}" ]; then
            log "  ERROR: No rollback file found for version ${version}"
            exit 1
        fi

        local filename
        filename=$(basename "${rollback_file}")
        log "  ROLLBACK: ${filename}"
        ${MYSQL_CMD} < "${rollback_file}"
        rolled_back=$((rolled_back + 1))
    done

    log "Rolled back ${rolled_back} migration(s)."
}

cmd_status() {
    log "Migration status for ${DB_NAME}:"
    local versions
    versions=$(get_applied_versions)

    if [ -z "${versions}" ]; then
        log "  No migrations applied."
        return
    fi

    for version in ${versions}; do
        local name
        name=$(${MYSQL_CMD} -N -e "SELECT name FROM schema_migrations WHERE version = ${version};" 2>/dev/null)
        local applied_at
        applied_at=$(${MYSQL_CMD} -N -e "SELECT applied_at FROM schema_migrations WHERE version = ${version};" 2>/dev/null)
        log "  [APPLIED] V${version}: ${name} (applied: ${applied_at})"
    done

    log ""
    log "Pending migrations:"
    local has_pending=false
    for migration_file in "${MIGRATIONS_DIR}"/V[0-9]*.sql; do
        [ -f "${migration_file}" ] || continue
        [[ "${migration_file}" == *"_rollback.sql" ]] && continue

        local filename
        filename=$(basename "${migration_file}")
        local version
        version=$(echo "${filename}" | grep -oP 'V\K[0-9]+')

        if ! is_applied "${version}"; then
            log "  [PENDING] V${version}: ${filename}"
            has_pending=true
        fi
    done

    if [ "${has_pending}" = false ]; then
        log "  (none)"
    fi
}

case "${1:-help}" in
    up)
        cmd_up "${2:-}"
        ;;
    down)
        cmd_down "${2:-}"
        ;;
    status)
        cmd_status
        ;;
    *)
        echo "Usage: $0 {up|down|status} [version]"
        exit 1
        ;;
esac
