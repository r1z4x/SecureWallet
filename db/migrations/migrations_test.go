package migrations

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestMigrationFilesExist verifies that V002 forward and rollback files exist
func TestMigrationFilesExist(t *testing.T) {
	migrationsDir := findMigrationsDir(t)

	forwardFile := filepath.Join(migrationsDir, "V002_wallets_and_transactions.sql")
	rollbackFile := filepath.Join(migrationsDir, "V002_wallets_and_transactions_rollback.sql")

	if _, err := os.Stat(forwardFile); os.IsNotExist(err) {
		t.Fatalf("forward migration file missing: %s", forwardFile)
	}
	if _, err := os.Stat(rollbackFile); os.IsNotExist(err) {
		t.Fatalf("rollback migration file missing: %s", rollbackFile)
	}
}

// TestV004MigrationFilesExist verifies that V004 forward and rollback files exist
func TestV004MigrationFilesExist(t *testing.T) {
	migrationsDir := findMigrationsDir(t)

	forwardFile := filepath.Join(migrationsDir, "V004_uuid_fk_enforcement.sql")
	rollbackFile := filepath.Join(migrationsDir, "V004_uuid_fk_enforcement_rollback.sql")

	if _, err := os.Stat(forwardFile); os.IsNotExist(err) {
		t.Fatalf("V004 forward migration file missing: %s", forwardFile)
	}
	if _, err := os.Stat(rollbackFile); os.IsNotExist(err) {
		t.Fatalf("V004 rollback migration file missing: %s", rollbackFile)
	}
}

// TestForwardMigrationContainsRequiredTables verifies the forward migration defines all required tables
func TestForwardMigrationContainsRequiredTables(t *testing.T) {
	content := readForwardMigration(t)

	requiredTables := []string{"wallets", "transactions", "idempotency_records", "schema_migrations"}
	for _, table := range requiredTables {
		if !strings.Contains(content, fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s", table)) {
			t.Errorf("forward migration missing table: %s", table)
		}
	}
}

// TestForwardMigrationContainsFinancialConstraints verifies critical financial integrity constraints
func TestForwardMigrationContainsFinancialConstraints(t *testing.T) {
	content := readForwardMigration(t)

	constraints := map[string]string{
		"non-negative balance":     "chk_wallet_balance_nonnegative",
		"positive amount":          "chk_transaction_amount_positive",
		"valid transaction types":  "chk_transaction_type",
		"valid transaction status": "chk_transaction_status",
		"idempotency uniqueness":   "uq_transaction_idempotency",
		"wallet user FK":           "fk_wallet_user",
		"transaction wallet FK":    "fk_transaction_wallet",
		"counterparty wallet FK":   "fk_transaction_counterparty_wallet",
	}

	for desc, constraint := range constraints {
		if !strings.Contains(content, constraint) {
			t.Errorf("forward migration missing constraint for %s: %s", desc, constraint)
		}
	}
}

// TestForwardMigrationContainsCounterpartyField verifies transfer tracking support
func TestForwardMigrationContainsCounterpartyField(t *testing.T) {
	content := readForwardMigration(t)

	if !strings.Contains(content, "counterparty_wallet_id") {
		t.Error("forward migration missing counterparty_wallet_id column for transfer tracking")
	}
}

// TestForwardMigrationContainsIdempotencyKey verifies idempotency support
func TestForwardMigrationContainsIdempotencyKey(t *testing.T) {
	content := readForwardMigration(t)

	if !strings.Contains(content, "idempotency_key") {
		t.Error("forward migration missing idempotency_key column")
	}
}

// TestForwardMigrationContainsVersionField verifies optimistic locking support
func TestForwardMigrationContainsVersionField(t *testing.T) {
	content := readForwardMigration(t)

	if !strings.Contains(content, "version INT") {
		t.Error("forward migration missing version column for optimistic locking")
	}
}

// TestRollbackMigrationDropsTablesInCorrectOrder verifies rollback drops dependents first
func TestRollbackMigrationDropsTablesInCorrectOrder(t *testing.T) {
	rollbackContent := readRollbackMigration(t)

	// Find positions of DROP statements
	idempotencyPos := strings.Index(rollbackContent, "DROP TABLE IF EXISTS idempotency_records")
	transactionsPos := strings.Index(rollbackContent, "DROP TABLE IF EXISTS transactions")
	walletsPos := strings.Index(rollbackContent, "DROP TABLE IF EXISTS wallets")

	if idempotencyPos == -1 {
		t.Error("rollback missing DROP for idempotency_records")
	}
	if transactionsPos == -1 {
		t.Error("rollback missing DROP for transactions")
	}
	if walletsPos == -1 {
		t.Error("rollback missing DROP for wallets")
	}

	// Verify order: transactions before wallets (transactions depends on wallets)
	if transactionsPos != -1 && walletsPos != -1 && transactionsPos > walletsPos {
		t.Error("rollback must DROP transactions before wallets (foreign key dependency)")
	}
}

// TestRollbackMigrationRemovesTrackingEntry verifies cleanup of schema_migrations
func TestRollbackMigrationRemovesTrackingEntry(t *testing.T) {
	content := readRollbackMigration(t)

	if !strings.Contains(content, "DELETE FROM schema_migrations WHERE version = 2") {
		t.Error("rollback must remove version 2 entry from schema_migrations")
	}
}

// TestForwardMigrationIsIdempotent verifies all CREATE statements use IF NOT EXISTS
func TestForwardMigrationIsIdempotent(t *testing.T) {
	content := readForwardMigration(t)

	lines := strings.Split(content, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "CREATE TABLE") && !strings.Contains(trimmed, "IF NOT EXISTS") {
			t.Errorf("line %d: CREATE TABLE without IF NOT EXISTS (breaks idempotency): %s", i+1, trimmed)
		}
	}
}

// TestForwardMigrationUsesInnoDB verifies engine choice for transaction support
func TestForwardMigrationUsesInnoDB(t *testing.T) {
	content := readForwardMigration(t)

	tables := []string{"wallets", "transactions", "idempotency_records", "schema_migrations"}
	for _, table := range tables {
		// Find the table block and verify ENGINE=InnoDB
		tableStart := strings.Index(content, fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s", table))
		if tableStart == -1 {
			continue
		}

		// Find the next CREATE TABLE or end of file
		nextCreate := strings.Index(content[tableStart+1:], "CREATE TABLE")
		var tableBlock string
		if nextCreate == -1 {
			tableBlock = content[tableStart:]
		} else {
			tableBlock = content[tableStart : tableStart+1+nextCreate]
		}

		if !strings.Contains(tableBlock, "ENGINE=InnoDB") {
			t.Errorf("table %s does not use InnoDB engine (required for foreign keys and transactions)", table)
		}
	}
}

// TestForwardMigrationRecordsMigration verifies INSERT into schema_migrations
func TestForwardMigrationRecordsMigration(t *testing.T) {
	content := readForwardMigration(t)

	if !strings.Contains(content, "INSERT IGNORE INTO schema_migrations") {
		t.Error("forward migration must record itself in schema_migrations using INSERT IGNORE")
	}
	if !strings.Contains(content, "VALUES (2,") {
		t.Error("forward migration must record version 2 in schema_migrations")
	}
}

// TestMigrationScriptExists verifies the shell migration runner exists
func TestMigrationScriptExists(t *testing.T) {
	migrationsDir := findMigrationsDir(t)
	scriptPath := filepath.Join(filepath.Dir(migrationsDir), "migrate.sh")

	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		t.Fatalf("migration runner script missing: %s", scriptPath)
	}

	info, err := os.Stat(scriptPath)
	if err != nil {
		t.Fatalf("cannot stat migration script: %v", err)
	}

	// Check executable bit (on Unix systems)
	if info.Mode()&0111 == 0 {
		t.Error("migration runner script is not executable")
	}
}

// --- helpers ---

func findMigrationsDir(t *testing.T) string {
	// Try common locations
	candidates := []string{
		"db/migrations",
		"../db/migrations",
		"../../db/migrations",
	}

	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}

	t.Fatal("cannot find db/migrations directory")
	return ""
}

func readForwardMigration(t *testing.T) string {
	migrationsDir := findMigrationsDir(t)
	content, err := os.ReadFile(filepath.Join(migrationsDir, "V002_wallets_and_transactions.sql"))
	if err != nil {
		t.Fatalf("cannot read forward migration: %v", err)
	}
	return string(content)
}

func readRollbackMigration(t *testing.T) string {
	migrationsDir := findMigrationsDir(t)
	content, err := os.ReadFile(filepath.Join(migrationsDir, "V002_wallets_and_transactions_rollback.sql"))
	if err != nil {
		t.Fatalf("cannot read rollback migration: %v", err)
	}
	return string(content)
}

func readV004ForwardMigration(t *testing.T) string {
	migrationsDir := findMigrationsDir(t)
	content, err := os.ReadFile(filepath.Join(migrationsDir, "V004_uuid_fk_enforcement.sql"))
	if err != nil {
		t.Fatalf("cannot read V004 forward migration: %v", err)
	}
	return string(content)
}

func readV004RollbackMigration(t *testing.T) string {
	migrationsDir := findMigrationsDir(t)
	content, err := os.ReadFile(filepath.Join(migrationsDir, "V004_uuid_fk_enforcement_rollback.sql"))
	if err != nil {
		t.Fatalf("cannot read V004 rollback migration: %v", err)
	}
	return string(content)
}

// TestV004EnforcesForeignKeys verifies V004 adds FK constraints for all entity relationships
func TestV004EnforcesForeignKeys(t *testing.T) {
	content := readV004ForwardMigration(t)

	requiredFKs := []string{
		"fk_wallets_user",
		"fk_transactions_wallet",
		"fk_transactions_counterparty_wallet",
		"fk_sessions_user",
		"fk_audit_logs_user",
		"fk_idempotency_records_user",
		"fk_support_tickets_user",
		"fk_login_history_user",
		"fk_blog_posts_author",
		"fk_blog_comments_post",
		"fk_security_alerts_user",
	}

	for _, fk := range requiredFKs {
		if !strings.Contains(content, fk) {
			t.Errorf("V004 missing foreign key constraint: %s", fk)
		}
	}
}

// TestV004RecordsV03 verifies V004 records V003 in schema_migrations
func TestV004RecordsV03(t *testing.T) {
	content := readV004ForwardMigration(t)

	if !strings.Contains(content, "VALUES (3,") {
		t.Error("V004 must record V003 in schema_migrations")
	}
	if !strings.Contains(content, "V003_audit_log_enhancements") {
		t.Error("V004 must reference V003_audit_log_enhancements migration name")
	}
}

// TestV004RecordsItself verifies V004 records itself in schema_migrations
func TestV004RecordsItself(t *testing.T) {
	content := readV004ForwardMigration(t)

	if !strings.Contains(content, "VALUES (4,") {
		t.Error("V004 must record itself in schema_migrations")
	}
	if !strings.Contains(content, "V004_uuid_fk_enforcement") {
		t.Error("V004 must reference its own migration name")
	}
}

// TestV004RollbackDropsAllForeignKeys verifies V004 rollback removes all FK constraints
func TestV004RollbackDropsAllForeignKeys(t *testing.T) {
	content := readV004RollbackMigration(t)

	requiredFKs := []string{
		"fk_wallets_user",
		"fk_transactions_wallet",
		"fk_transactions_counterparty_wallet",
		"fk_sessions_user",
		"fk_audit_logs_user",
		"fk_idempotency_records_user",
		"fk_support_tickets_user",
		"fk_login_history_user",
		"fk_blog_posts_author",
		"fk_blog_comments_post",
		"fk_security_alerts_user",
	}

	for _, fk := range requiredFKs {
		if !strings.Contains(content, fk) {
			t.Errorf("V004 rollback missing DROP for foreign key: %s", fk)
		}
	}
}

// TestV004UsesMySQLConstraintSyntax verifies V004 does not use non-MySQL DROP CONSTRAINT syntax.
func TestV004UsesMySQLConstraintSyntax(t *testing.T) {
	files := map[string]string{
		"forward":  readV004ForwardMigration(t),
		"rollback": readV004RollbackMigration(t),
	}

	for name, content := range files {
		if strings.Contains(content, "DROP CONSTRAINT") {
			t.Errorf("V004 %s migration uses DROP CONSTRAINT; MySQL requires DROP FOREIGN KEY or DROP CHECK", name)
		}
	}
}

// TestV004RollbackRemovesTrackingEntry verifies V004 rollback removes its version entry
func TestV004RollbackRemovesTrackingEntry(t *testing.T) {
	content := readV004RollbackMigration(t)

	if !strings.Contains(content, "DELETE FROM schema_migrations WHERE version = 4") {
		t.Error("V004 rollback must remove version 4 entry from schema_migrations")
	}
}

// TestInitSQLUsesUUIDForAllPrimaryKeys verifies init.sql uses CHAR(36) for all primary keys
func TestInitSQLUsesUUIDForAllPrimaryKeys(t *testing.T) {
	migrationsDir := findMigrationsDir(t)
	initPath := filepath.Join(filepath.Dir(migrationsDir), "init.sql")

	content, err := os.ReadFile(initPath)
	if err != nil {
		t.Fatalf("cannot read init.sql: %v", err)
	}

	tables := []string{"users", "wallets", "transactions", "sessions", "audit_logs",
		"idempotency_records", "support_tickets", "login_history",
		"blog_posts", "blog_comments", "blog_categories", "blog_tags", "security_alerts"}

	for _, table := range tables {
		// Find the table definition
		tableStart := strings.Index(string(content), fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s", table))
		if tableStart == -1 {
			t.Errorf("init.sql missing table: %s", table)
			continue
		}

		// Find the end of the table definition (next CREATE or end of file)
		nextCreate := strings.Index(string(content)[tableStart+1:], "CREATE TABLE")
		var tableBlock string
		if nextCreate == -1 {
			tableBlock = string(content)[tableStart:]
		} else {
			tableBlock = string(content)[tableStart : tableStart+1+nextCreate]
		}

		// Check that id column uses CHAR(36)
		if !strings.Contains(tableBlock, "id CHAR(36)") && !strings.Contains(tableBlock, "id VARCHAR(36)") {
			t.Errorf("table %s primary key does not use CHAR(36) UUID type", table)
		}
	}
}

// TestInitSQLHasForeignKeysForUserReferences verifies all tables referencing users have FK constraints
func TestInitSQLHasForeignKeysForUserReferences(t *testing.T) {
	migrationsDir := findMigrationsDir(t)
	initPath := filepath.Join(filepath.Dir(migrationsDir), "init.sql")

	content, err := os.ReadFile(initPath)
	if err != nil {
		t.Fatalf("cannot read init.sql: %v", err)
	}

	// Tables that should have FK to users
	userFKTables := []string{"wallets", "sessions", "audit_logs", "idempotency_records",
		"support_tickets", "login_history", "blog_posts", "security_alerts"}

	for _, table := range userFKTables {
		tableStart := strings.Index(string(content), fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s", table))
		if tableStart == -1 {
			continue
		}

		nextCreate := strings.Index(string(content)[tableStart+1:], "CREATE TABLE")
		var tableBlock string
		if nextCreate == -1 {
			tableBlock = string(content)[tableStart:]
		} else {
			tableBlock = string(content)[tableStart : tableStart+1+nextCreate]
		}

		if !strings.Contains(tableBlock, "REFERENCES users(id)") {
			t.Errorf("table %s missing foreign key reference to users(id)", table)
		}
	}
}

// TestAllModelsUseUUIDType verifies all Go models use uuid.UUID for ID fields
func TestAllModelsUseUUIDType(t *testing.T) {
	modelFiles := []string{
		"user.go", "wallet.go", "transaction.go", "session.go",
		"audit_log.go", "idempotency_record.go", "support_ticket.go",
		"login_history.go", "blog.go",
	}

	modelsDir := findModelsDir(t)

	for _, file := range modelFiles {
		path := filepath.Join(modelsDir, file)
		content, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("cannot read model file %s: %v", file, err)
		}

		// Check that the file imports github.com/google/uuid
		if !strings.Contains(string(content), `"github.com/google/uuid"`) {
			t.Errorf("model file %s does not import github.com/google/uuid", file)
		}

		// Check that ID field uses uuid.UUID
		if !strings.Contains(string(content), "uuid.UUID") {
			t.Errorf("model file %s does not use uuid.UUID for ID field", file)
		}
	}
}

// TestTransactionTypesMatchDBConstraint verifies Go code uses lowercase transaction types matching DB CHECK
func TestTransactionTypesMatchDBConstraint(t *testing.T) {
	// Read transfer.go
	transferPath := findTransferGoPath(t)
	content, err := os.ReadFile(transferPath)
	if err != nil {
		t.Fatalf("cannot read transfer.go: %v", err)
	}

	validTypes := []string{"`deposit`", "`withdrawal`", "`transfer_in`", "`transfer_out`"}
	invalidTypes := []string{"`DEPOSIT`", "`WITHDRAWAL`", "`TRANSFER_IN`", "`TRANSFER_OUT`",
		"\"DEPOSIT\"", "\"WITHDRAWAL\"", "\"TRANSFER_IN\"", "\"TRANSFER_OUT\""}

	for _, invalidType := range invalidTypes {
		// Remove backticks for string search
		searchType := strings.Trim(invalidType, "`")
		if strings.Contains(string(content), fmt.Sprintf("Type:        %s", searchType)) {
			t.Errorf("transfer.go uses invalid transaction type %s (should be lowercase to match DB CHECK constraint)", searchType)
		}
	}

	_ = validTypes
}

// --- additional helpers ---

func findModelsDir(t *testing.T) string {
	candidates := []string{
		"internal/models",
		"../internal/models",
		"../../internal/models",
	}

	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}

	t.Fatal("cannot find internal/models directory")
	return ""
}

func findTransferGoPath(t *testing.T) string {
	candidates := []string{
		"internal/services/transfer.go",
		"../internal/services/transfer.go",
		"../../internal/services/transfer.go",
	}

	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}

	t.Fatal("cannot find internal/services/transfer.go")
	return ""
}
