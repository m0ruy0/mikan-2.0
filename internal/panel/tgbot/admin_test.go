package tgbot

import (
	"strconv"
	"strings"
	"testing"

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
