package admin

import (
	"context"
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/service"
)

type inspectionShadowAdminService struct {
	*stubAdminService
	parent  *service.Account
	lookups int
}

func (s *inspectionShadowAdminService) GetAccountsByIDs(_ context.Context, ids []int64) ([]*service.Account, error) {
	s.lookups++
	if len(ids) != 1 || ids[0] != s.parent.ID {
		panic("unexpected parent lookup")
	}
	return []*service.Account{s.parent}, nil
}

func TestCodexInspectionShadowParentPlan(t *testing.T) {
	parentID := int64(10)
	svc := &inspectionShadowAdminService{stubAdminService: newStubAdminService(), parent: &service.Account{
		ID: parentID, Credentials: map[string]any{"plan_type": "pro"},
	}}
	h := &AccountHandler{adminService: svc}
	shadow := service.Account{ID: 20, ParentAccountID: &parentID, Type: service.AccountTypeOAuth,
		GroupIDs: []int64{101}, Status: service.StatusActive, Schedulable: true, Priority: 10}
	owned := shadow
	owned.ID, owned.Credentials = 21, map[string]any{"plan_type": "plus"}
	rows := h.inspectionAccounts(context.Background(), []service.Account{shadow, owned})
	if svc.lookups != 1 || rows[0].PlanType != "pro" || rows[1].PlanType != "plus" {
		t.Fatalf("inspection lost native plan fallback or credential precedence: %+v", rows)
	}
	if !rows[0].Available || rows[0].GroupIDs[0] != 101 || rows[0].Priority != 10 {
		t.Fatalf("plan enrichment changed native account state: %+v", rows[0])
	}
	// The single-account recheck must resolve the parent's current plan again.
	svc.parent.Credentials["plan_type"] = "pro_lite"
	current := h.inspectionAccounts(context.Background(), []service.Account{shadow})[0]
	if svc.lookups != 2 || current.PlanType != "pro_lite" {
		t.Fatalf("priority recheck used a stale or missing parent plan: %+v", current)
	}
	if shadow.GetCredential("plan_type") != "" {
		t.Fatal("inspection copied parent credentials into the shadow")
	}
}
