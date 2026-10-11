package tgbot

import (
	"strconv"
	"strings"
	"testing"

	"mikan/internal/panel/domain"
	"mikan/internal/panel/settings"
)

func TestAdminCommandsUseVerifiedInfrastructureChat(t *testing.T) {
	e := setup(t)
	const admin int64 = 777
	if err := settings.Set(e.ctx, e.set, KeyInfraAdminChat, admin); err != nil {
		t.Fatal(err)
	}

	n := e.tg.count()
	e.say(admin, "/admin stats")
	calls := e.tg.wait(t, n, "sendMessage")
	msg, ok := find(calls, "sendMessage")
	if !ok || !strings.Contains(text(msg), "Клиенты: <b>1</b>") || !strings.Contains(text(msg), "Ноды панели") && !strings.Contains(text(msg), "Обзор панели") {
		t.Fatalf("admin stats response: %q", text(msg))
	}

	// An ordinary subscriber who knows the command must not get any panel data.
	e.later()
	n = e.tg.count()
	e.say(555, "/admin users")
	calls = e.tg.wait(t, n, "sendMessage")
	msg, ok = find(calls, "sendMessage")
	if !ok || !strings.Contains(text(msg), "Доступ запрещён") {
		t.Fatalf("unauthorized admin command: %q", text(msg))
	}

	e.later()
	n = e.tg.count()
	e.say(admin, "/admin user "+strconv.FormatInt(e.user.ID, 10))
	calls = e.tg.wait(t, n, "sendMessage")
	msg, ok = find(calls, "sendMessage")
	if !ok || !strings.Contains(text(msg), "<b>Клиент #") || !strings.Contains(text(msg), "Анна") || !strings.Contains(text(msg), "/admin devices") {
		t.Fatalf("admin user response: %q", text(msg))
	}

	e.later()
	n = e.tg.count()
	e.say(admin, "/admin users")
	calls = e.tg.wait(t, n, "sendMessage")
	msg, ok = find(calls, "sendMessage")
	if !ok || !strings.Contains(text(msg), "Последние клиенты") || !strings.Contains(text(msg), "Анна") {
		t.Fatalf("admin users response: %q", text(msg))
	}

	e.later()
	n = e.tg.count()
	e.say(admin, "/admin devices "+strconv.FormatInt(e.user.ID, 10))
	calls = e.tg.wait(t, n, "sendMessage")
	msg, ok = find(calls, "sendMessage")
	if !ok || !strings.Contains(text(msg), "Устройства") {
		t.Fatalf("admin devices response: %q", text(msg))
	}

	e.later()
	n = e.tg.count()
	e.say(admin, "/admin nodes")
	calls = e.tg.wait(t, n, "sendMessage")
	msg, ok = find(calls, "sendMessage")
	if !ok || !strings.Contains(text(msg), "Ноды панели") {
		t.Fatalf("admin nodes response: %q", text(msg))
	}
}

func TestAdminMutationsRequireConfirmationAndAreAudited(t *testing.T) {
	e := setup(t, func(e *env, d *Deps) {
		d.Users = domain.NewUsers(e.st, domain.NewPool(e.st, e.clock), noChanges{}, e.clock)
	})
	const admin int64 = 777
	if err := settings.Set(e.ctx, e.set, KeyInfraAdminChat, admin); err != nil {
		t.Fatal(err)
	}

	e.say(admin, "/admin disable "+strconv.FormatInt(e.user.ID, 10))
	calls := e.tg.wait(t, 0, "sendMessage")
	msg, ok := find(calls, "sendMessage")
	if !ok || !strings.Contains(text(msg), "/admin confirm ") {
		t.Fatalf("disable must request confirmation: %q", text(msg))
	}
	const marker = "/admin confirm "
	start := strings.Index(text(msg), marker)
	token := strings.Fields(text(msg)[start+len(marker):])[0]
	token = strings.TrimSuffix(token, "</code>")
	if token == "" {
		t.Fatal("confirmation token is empty")
	}

	e.later()
	n := e.tg.count()
	e.say(admin, "/admin confirm "+token)
	calls = e.tg.wait(t, n, "sendMessage")
	msg, ok = find(calls, "sendMessage")
	if !ok || !strings.Contains(text(msg), "выключен") {
		t.Fatalf("disable confirmation: %q", text(msg))
	}
	u, err := e.st.Q.GetUser(e.ctx, e.user.ID)
	if err != nil || u.Status != "disabled" {
		t.Fatalf("the domain action did not disable the user: status=%q err=%v", u.Status, err)
	}
	var action string
	if err := e.st.DB.QueryRowContext(e.ctx, "SELECT action FROM audit_log ORDER BY id DESC LIMIT 1").Scan(&action); err != nil || action != "telegram.admin.disable" {
		t.Fatalf("admin action was not audited: action=%q err=%v", action, err)
	}

	// One-time confirmation tokens cannot be replayed.
	e.later()
	n = e.tg.count()
	e.say(admin, "/admin confirm "+token)
	calls = e.tg.wait(t, n, "sendMessage")
	msg, ok = find(calls, "sendMessage")
	if !ok || !strings.Contains(text(msg), "Токен подтверждения неверен") {
		t.Fatalf("a replayed token should be rejected: %q", text(msg))
	}
}
