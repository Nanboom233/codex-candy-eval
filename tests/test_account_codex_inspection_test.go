// sub2api-package: backend/internal/repository
package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	dbent "github.com/Wei-Shaw/sub2api/ent"
	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/stretchr/testify/require"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
)

func TestCodexInspectionStateAtomicWrite(t *testing.T) {
	for _, outcome := range []string{"commit", "conflict", "outbox error"} {
		t.Run(outcome, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			t.Cleanup(func() { _ = db.Close() })
			client := dbent.NewClient(dbent.Driver(entsql.OpenDB(dialect.Postgres, db)))
			t.Cleanup(func() { _ = client.Close() })
			repo := newAccountRepositoryWithSQL(client, nil, nil)
			expected := &service.Account{ID: 20, Platform: service.PlatformOpenAI, Type: service.AccountTypeAPIKey,
				Status: service.StatusActive, Schedulable: true, Priority: 10,
				Credentials: map[string]any{"api_key": "old-key", "model_mapping": map[string]any{"m": "m"}},
				Extra:       map[string]any{"codex_eval_inspection_enabled": true, "codex_eval_inspection_verdict": "normal", "apikey_quota_exhausted": false}}
			mapping := map[string]any{"blacklist": "blacklist"}
			extra := map[string]any{"codex_eval_inspection_verdict": "degraded", "codex_eval_inspection_removed_models": []string{"m"}, "apikey_quota_exhausted": false}
			mock.ExpectBegin()
			query := mock.ExpectExec(`(?s)`+regexp.QuoteMeta("UPDATE accounts SET extra = COALESCE(extra, '{}'::jsonb) || $1::jsonb, credentials = jsonb_set(COALESCE(credentials, '{}'::jsonb), '{model_mapping}', $7::jsonb), updated_at = NOW()")+`.*`+
				regexp.QuoteMeta("AND COALESCE(credentials, '{}'::jsonb) = $5::jsonb AND proxy_id IS NOT DISTINCT FROM $6")+`.*`+
				regexp.QuoteMeta("AND COALESCE(extra -> 'codex_eval_inspection_enabled', 'true'::jsonb) <> 'false'::jsonb")).
				WithArgs(`{"codex_eval_inspection_removed_models":["m"],"codex_eval_inspection_verdict":"degraded"}`, int64(20), service.PlatformOpenAI, service.AccountTypeAPIKey,
					`{"api_key":"old-key","model_mapping":{"m":"m"}}`, nil, `{"blacklist":"blacklist"}`)
			if outcome == "conflict" {
				query.WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectRollback()
			} else {
				query.WillReturnResult(sqlmock.NewResult(0, 1))
				outbox := mock.ExpectExec(regexp.QuoteMeta("INSERT INTO scheduler_outbox")).
					WithArgs(service.SchedulerOutboxEventAccountChanged, int64(20), nil, nil, sqlmock.AnyArg())
				if outcome == "outbox error" {
					outbox.WillReturnError(errors.New("fixture outbox failure"))
					mock.ExpectRollback()
				} else {
					outbox.WillReturnResult(sqlmock.NewResult(1, 1))
					mock.ExpectCommit()
				}
			}
			applied, err := repo.UpdateCodexInspectionState(context.Background(), expected, mapping, extra)
			if outcome == "outbox error" {
				require.ErrorContains(t, err, "fixture outbox failure")
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, outcome == "commit", applied)
			require.Equal(t, "old-key", expected.Credentials["api_key"])
			require.Equal(t, false, expected.Extra["apikey_quota_exhausted"])
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
