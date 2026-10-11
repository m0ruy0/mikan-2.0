package tgbot

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"fmt"
	"html"
	"strconv"
	"strings"
	"time"

	"mikan/internal/panel/domain"
	"mikan/internal/panel/store/db"
)

// adminCommand handles the restricted /admin surface of the subscriber bot. It is
// deliberately tied to the separately verified infrastructure-admin private chat:
// knowing the command name or a subscription link never grants administrative access.
func (b *Bot) adminCommand(ctx context.Context, out *Outbox, m *Message, raw string) bool {
	parts := strings.Fields(raw)
	if len(parts) == 0 {
		return false
	}
	command := strings.ToLower(strings.SplitN(parts[0], "@", 2)[0])
	if command != "/admin" {
		return false
	}

	chat := m.Chat.ID
	if m.Chat.Type != "private" || m.From == nil || m.From.ID != chat || !b.InfrastructureAdminChatIs(ctx, chat) {
		b.adminReply(ctx, out, chat, "⛔ Доступ запрещён.\n⛔ Access denied.")
		return true
	}

	args := parts[1:]
	if len(args) == 0 || args[0] == "help" {
		b.adminReply(ctx, out, chat, b.adminHelp(ctx))
		return true
	}

	switch strings.ToLower(args[0]) {
	case "stats":
		b.adminReply(ctx, out, chat, b.adminStats(ctx))
	case "users", "clients":
		b.adminReply(ctx, out, chat, b.adminUsers(ctx))
	case "user", "subscription":
		if len(args) != 2 {
			b.adminReply(ctx, out, chat, b.adminUsage(ctx))
			return true
		}
		id, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil || id < 1 {
			b.adminReply(ctx, out, chat, b.adminUsage(ctx))
			return true
		}
		b.adminReply(ctx, out, chat, b.adminUser(ctx, id))
	case "devices":
		if len(args) != 2 {
			b.adminReply(ctx, out, chat, b.adminUsage(ctx))
			return true
		}
		id, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil || id < 1 {
			b.adminReply(ctx, out, chat, b.adminUsage(ctx))
			return true
		}
		b.adminReply(ctx, out, chat, b.adminDevices(ctx, id))
	case "nodes":
		b.adminReply(ctx, out, chat, b.adminNodes(ctx))
	case "disable", "enable", "reset-traffic":
		if len(args) != 2 {
			b.adminReply(ctx, out, chat, b.adminUsage(ctx))
			return true
		}
		id, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil || id < 1 {
			b.adminReply(ctx, out, chat, b.adminUsage(ctx))
			return true
		}
		b.adminPrepareAction(ctx, out, chat, id, strings.ToLower(args[0]), 0)
	case "extend":
		if len(args) != 3 {
			b.adminReply(ctx, out, chat, b.adminUsage(ctx))
			return true
		}
		id, idErr := strconv.ParseInt(args[1], 10, 64)
		days, daysErr := strconv.ParseInt(args[2], 10, 64)
		if idErr != nil || id < 1 || daysErr != nil || days < 1 || days > 3650 {
			b.adminReply(ctx, out, chat, b.adminUsage(ctx))
			return true
		}
		b.adminPrepareAction(ctx, out, chat, id, "extend", days)
	case "confirm":
		if len(args) != 2 {
			b.adminReply(ctx, out, chat, b.adminUsage(ctx))
			return true
		}
		b.adminConfirmAction(ctx, out, chat, args[1])
	case "cancel":
		if len(args) != 2 {
			b.adminReply(ctx, out, chat, b.adminUsage(ctx))
			return true
		}
		b.adminCancelAction(ctx, out, chat, args[1])
	default:
		b.adminReply(ctx, out, chat, b.adminUsage(ctx))
	}
	return true
}

func (b *Bot) adminReply(ctx context.Context, out *Outbox, chat int64, text string) {
	out.Reply(chat, "admin:"+strconv.FormatInt(chat, 10), 1, func(ctx context.Context, c *Client) error {
		_, err := c.Send(ctx, chat, text, nil, false)
		return err
	})
}

func (b *Bot) adminHelp(ctx context.Context) string {
	if b.Config(ctx).Lang == "en" {
		return "<b>Mikan admin</b>\n/admin stats — overview\n/admin users — latest 10 clients\n/admin user ID — subscription details\n/admin devices ID — devices and recent IPs\n/admin nodes — panel nodes\n/admin help — this help\n\nWrite commands require a one-time confirmation token valid for five minutes."
	}
	return "<b>Администрирование Mikan</b>\n/admin stats — общая статистика\n/admin users — последние 10 клиентов\n/admin user ID — данные подписки\n/admin devices ID — устройства и последние IP\n/admin nodes — ноды панели\n/admin help — эта справка\n\nКоманды доступны только подтверждённому чату администратора. Изменения требуют одноразового токена подтверждения на 5 минут."
}

func (b *Bot) adminUsage(ctx context.Context) string {
	if b.Config(ctx).Lang == "en" {
		return "Usage: /admin stats, /admin users, /admin user ID, /admin devices ID, /admin nodes, /admin disable ID, /admin enable ID, /admin extend ID DAYS, /admin reset-traffic ID, /admin confirm TOKEN, /admin cancel TOKEN."
	}
	return "Использование: /admin stats, /admin users, /admin user ID, /admin devices ID, /admin nodes, /admin disable ID, /admin enable ID, /admin extend ID DAYS, /admin reset-traffic ID, /admin confirm TOKEN, /admin cancel TOKEN."
}

func (b *Bot) adminStats(ctx context.Context) string {
	users, err := b.d.Store.Q.ListUsers(ctx)
	if err != nil {
		return b.adminReadError(ctx)
	}
	nodes, err := b.d.Store.Q.ListNodes(ctx)
	if err != nil {
		return b.adminReadError(ctx)
	}
	active, expiring := 0, 0
	now := b.d.Now().Unix()
	for _, u := range users {
		if u.Status == "active" && (!u.ExpiresAt.Valid || u.ExpiresAt.Int64 == 0 || u.ExpiresAt.Int64 > now) {
			active++
		}
		if u.ExpiresAt.Valid && u.ExpiresAt.Int64 > now && u.ExpiresAt.Int64-now <= int64(7*24*time.Hour/time.Second) {
			expiring++
		}
	}
	enabledNodes := 0
	for _, n := range nodes {
		if n.Enabled != 0 {
			enabledNodes++
		}
	}
	linked, err := b.d.Store.Q.CountTgLinks(ctx)
	if err != nil {
		return b.adminReadError(ctx)
	}
	if b.Config(ctx).Lang == "en" {
		return fmt.Sprintf("<b>Panel overview</b>\nClients: <b>%d</b>\nActive subscriptions: <b>%d</b>\nExpiring in 7 days: <b>%d</b>\nTelegram-linked subscriptions: <b>%d</b>\nNodes online-enabled: <b>%d/%d</b>", len(users), active, expiring, linked, enabledNodes, len(nodes))
	}
	return fmt.Sprintf("<b>Обзор панели</b>\nКлиенты: <b>%d</b>\nАктивные подписки: <b>%d</b>\nИстекают за 7 дней: <b>%d</b>\nПодписок, привязанных к Telegram: <b>%d</b>\nВключённые ноды: <b>%d из %d</b>", len(users), active, expiring, linked, enabledNodes, len(nodes))
}

func (b *Bot) adminUsers(ctx context.Context) string {
	users, err := b.d.Store.Q.ListUsers(ctx)
	if err != nil {
		return b.adminReadError(ctx)
	}
	if len(users) == 0 {
		if b.Config(ctx).Lang == "en" {
			return "No clients found."
		}
		return "Клиентов пока нет."
	}
	const limit = 10
	if len(users) > limit {
		users = users[:limit]
	}
	var out strings.Builder
	if b.Config(ctx).Lang == "en" {
		out.WriteString("<b>Latest clients</b>\n")
	} else {
		out.WriteString("<b>Последние клиенты</b>\n")
	}
	for _, u := range users {
		name := html.EscapeString(u.Name)
		status := html.EscapeString(u.Status)
		fmt.Fprintf(&out, "• <code>#%d</code> %s — %s", u.ID, name, status)
		if u.ExpiresAt.Valid && u.ExpiresAt.Int64 > 0 {
			fmt.Fprintf(&out, " · до %s", time.Unix(u.ExpiresAt.Int64, 0).UTC().Format("02.01.2006"))
		}
		out.WriteByte('\n')
	}
	if len(users) == limit {
		if b.Config(ctx).Lang == "en" {
			out.WriteString("\nUse /admin user ID for details.")
		} else {
			out.WriteString("\nДля подробностей: /admin user ID.")
		}
	}
	return out.String()
}

func (b *Bot) adminUser(ctx context.Context, id int64) string {
	u, err := b.d.Store.Q.GetUser(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			if b.Config(ctx).Lang == "en" {
				return "Client not found."
			}
			return "Клиент не найден."
		}
		return b.adminReadError(ctx)
	}
	name := html.EscapeString(u.Name)
	status := html.EscapeString(u.Status)
	traffic := formatAdminBytes(u.UsedUp + u.UsedDown)
	limit := "без лимита"
	if b.Config(ctx).Lang == "en" {
		limit = "unlimited"
	}
	if u.TrafficLimit.Valid && u.TrafficLimit.Int64 > 0 {
		limit = formatAdminBytes(u.TrafficLimit.Int64)
	}
	expiry := "без срока"
	if b.Config(ctx).Lang == "en" {
		expiry = "no expiry"
	}
	if u.ExpiresAt.Valid && u.ExpiresAt.Int64 > 0 {
		expiry = time.Unix(u.ExpiresAt.Int64, 0).UTC().Format("02 Jan 2006 15:04 UTC")
	}
	plan := "—"
	if u.TariffID.Valid {
		if tariff, err := b.d.Store.Q.GetTariff(ctx, u.TariffID.Int64); err == nil {
			plan = html.EscapeString(tariff.Name)
		}
	}
	bound, boundErr := b.d.Store.Q.ListBoundDevices(ctx, id)
	devices, devicesErr := b.d.Store.Q.ListUserDevices(ctx, id)
	if boundErr != nil || devicesErr != nil {
		return b.adminReadError(ctx)
	}
	if b.Config(ctx).Lang == "en" {
		return fmt.Sprintf("<b>Client #%d</b>\nName: %s\nStatus: %s\nPlan: %s\nExpires: %s\nTraffic: %s / %s\nDevices: %d bound, %d recent connections\nDetails: /admin devices %d", id, name, status, plan, expiry, traffic, limit, len(bound), len(devices), id)
	}
	return fmt.Sprintf("<b>Клиент #%d</b>\nИмя: %s\nСтатус: %s\nТариф: %s\nСрок: %s\nТрафик: %s / %s\nУстройства: %d привязанных, %d недавних подключений\nПодробнее: /admin devices %d", id, name, status, plan, expiry, traffic, limit, len(bound), len(devices), id)
}

func (b *Bot) adminDevices(ctx context.Context, id int64) string {
	u, err := b.d.Store.Q.GetUser(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			if b.Config(ctx).Lang == "en" {
				return "Client not found."
			}
			return "Клиент не найден."
		}
		return b.adminReadError(ctx)
	}
	bound, err := b.d.Store.Q.ListBoundDevices(ctx, id)
	if err != nil {
		return b.adminReadError(ctx)
	}
	recent, err := b.d.Store.Q.ListUserDevices(ctx, id)
	if err != nil {
		return b.adminReadError(ctx)
	}
	var out strings.Builder
	if b.Config(ctx).Lang == "en" {
		fmt.Fprintf(&out, "<b>Devices · client #%d (%s)</b>\n", id, html.EscapeString(u.Name))
	} else {
		fmt.Fprintf(&out, "<b>Устройства · клиент #%d (%s)</b>\n", id, html.EscapeString(u.Name))
	}
	if len(bound) == 0 && len(recent) == 0 {
		if b.Config(ctx).Lang == "en" {
			out.WriteString("No device records.\n")
		} else {
			out.WriteString("Записей об устройствах нет.\n")
		}
	}
	const limit = 8
	for i, d := range bound {
		if i >= limit {
			break
		}
		label := strings.TrimSpace(strings.Join([]string{d.Model, d.OS, d.App}, " · "))
		if label == "" {
			label = "device"
		}
		fmt.Fprintf(&out, "• %s — %s\n", html.EscapeString(label), time.Unix(d.LastSeen, 0).UTC().Format("02.01.2006"))
	}
	for i, d := range recent {
		if i >= limit {
			break
		}
		fmt.Fprintf(&out, "• IP <code>%s</code> — %s\n", html.EscapeString(d.Ip), time.Unix(d.LastSeen, 0).UTC().Format("02.01.2006"))
	}
	return out.String()
}

func (b *Bot) adminNodes(ctx context.Context) string {
	nodes, err := b.d.Store.Q.ListNodes(ctx)
	if err != nil {
		return b.adminReadError(ctx)
	}
	var out strings.Builder
	if b.Config(ctx).Lang == "en" {
		out.WriteString("<b>Panel nodes</b>\n")
	} else {
		out.WriteString("<b>Ноды панели</b>\n")
	}
	if len(nodes) == 0 {
		if b.Config(ctx).Lang == "en" {
			out.WriteString("No nodes.\n")
		} else {
			out.WriteString("Нод нет.\n")
		}
		return out.String()
	}
	for _, n := range nodes {
		state := "выключена"
		if n.Enabled != 0 {
			state = "включена"
		}
		if b.Config(ctx).Lang == "en" {
			state = "disabled"
			if n.Enabled != 0 {
				state = "enabled"
			}
		}
		host := strings.TrimSpace(n.PublicHost)
		if host == "" {
			host = strings.TrimSpace(n.Domain)
		}
		if host == "" {
			host = "—"
		}
		inbounds, inboundErr := b.d.Store.Q.ListNodeInbounds(ctx, n.ID)
		if inboundErr != nil {
			return b.adminReadError(ctx)
		}
		fmt.Fprintf(&out, "• <b>%s</b> · %s · %d inbounds · <code>%s</code>\n", html.EscapeString(n.Name), state, len(inbounds), html.EscapeString(host))
	}
	return out.String()
}

func (b *Bot) adminReadError(ctx context.Context) string {
	if b.Config(ctx).Lang == "en" {
		return "Could not load panel data. Check the panel logs."
	}
	return "Не удалось получить данные панели. Проверь журналы."
}

func formatAdminBytes(n int64) string {
	if n < 0 {
		n = 0
	}
	const unit = int64(1024)
	if n < unit {
		return fmt.Sprintf("%d B", n)
	}
	div, exp := unit, 0
	for v := n / unit; v >= unit && exp < 5; v /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB", float64(n)/float64(div), "KMGTPE"[exp])
}

type pendingAdminAction struct {
	ChatID int64
	UserID int64
	Action string
	Days   int64
	Until  time.Time
}

func (b *Bot) adminPrepareAction(ctx context.Context, out *Outbox, chat, userID int64, action string, days int64) {
	if b.d.Users == nil {
		b.adminReply(ctx, out, chat, b.adminMutationUnavailable(ctx))
		return
	}
	u, err := b.d.Users.Get(ctx, userID)
	if err != nil {
		if err == domain.ErrNotFound || err == sql.ErrNoRows {
			b.adminReply(ctx, out, chat, b.adminNotFound(ctx))
			return
		}
		b.adminReply(ctx, out, chat, b.adminReadError(ctx))
		return
	}

	raw := make([]byte, 18)
	if _, err := rand.Read(raw); err != nil {
		b.adminReply(ctx, out, chat, b.adminReadError(ctx))
		return
	}
	token := base64.RawURLEncoding.EncodeToString(raw)
	until := b.d.Now().Add(5 * time.Minute)
	pending := pendingAdminAction{ChatID: chat, UserID: userID, Action: action, Days: days, Until: until}
	b.adminActionMu.Lock()
	if b.adminActions == nil {
		b.adminActions = make(map[string]pendingAdminAction)
	}
	for key, item := range b.adminActions {
		if !b.d.Now().Before(item.Until) {
			delete(b.adminActions, key)
		}
	}
	b.adminActions[token] = pending
	b.adminActionMu.Unlock()

	name := html.EscapeString(u.Name)
	var prompt string
	switch action {
	case domain.BulkDisable:
		prompt = fmt.Sprintf("⚠️ Выключить доступ клиента #%d (%s)?\nПодтвердите: <code>/admin confirm %s</code>", userID, name, token)
		if b.Config(ctx).Lang == "en" {
			prompt = fmt.Sprintf("⚠️ Disable access for client #%d (%s)?\nConfirm: <code>/admin confirm %s</code>", userID, name, token)
		}
	case domain.BulkEnable:
		prompt = fmt.Sprintf("Включить доступ клиента #%d (%s)?\nПодтвердите: <code>/admin confirm %s</code>", userID, name, token)
		if b.Config(ctx).Lang == "en" {
			prompt = fmt.Sprintf("Enable access for client #%d (%s)?\nConfirm: <code>/admin confirm %s</code>", userID, name, token)
		}
	case domain.BulkExtend:
		prompt = fmt.Sprintf("Продлить подписку клиента #%d (%s) на %d дн.?\nПодтвердите: <code>/admin confirm %s</code>", userID, name, days, token)
		if b.Config(ctx).Lang == "en" {
			prompt = fmt.Sprintf("Extend client #%d (%s) by %d days?\nConfirm: <code>/admin confirm %s</code>", userID, name, days, token)
		}
	case domain.BulkReset:
		prompt = fmt.Sprintf("Сбросить счётчик трафика клиента #%d (%s)?\nПодтвердите: <code>/admin confirm %s</code>", userID, name, token)
		if b.Config(ctx).Lang == "en" {
			prompt = fmt.Sprintf("Reset traffic for client #%d (%s)?\nConfirm: <code>/admin confirm %s</code>", userID, name, token)
		}
	default:
		b.adminReply(ctx, out, chat, b.adminUsage(ctx))
		return
	}
	prompt += "\nТокен действует 5 минут. Отмена: <code>/admin cancel " + token + "</code>"
	if b.Config(ctx).Lang == "en" {
		prompt += "\nToken expires in 5 minutes. Cancel: <code>/admin cancel " + token + "</code>"
	}
	b.adminReply(ctx, out, chat, prompt)
}

func (b *Bot) adminConfirmAction(ctx context.Context, out *Outbox, chat int64, token string) {
	b.adminActionMu.Lock()
	pending, ok := b.adminActions[token]
	if ok {
		delete(b.adminActions, token)
	}
	b.adminActionMu.Unlock()
	if !ok || pending.ChatID != chat || !b.d.Now().Before(pending.Until) {
		b.adminReply(ctx, out, chat, b.adminTokenInvalid(ctx))
		return
	}
	if b.d.Users == nil {
		b.adminReply(ctx, out, chat, b.adminMutationUnavailable(ctx))
		return
	}

	action := pending.Action
	opt := domain.BulkOpt{Days: pending.Days}
	changed, err := b.d.Users.Bulk(ctx, []int64{pending.UserID}, action, opt)
	if err != nil {
		b.adminReply(ctx, out, chat, b.adminReadError(ctx))
		return
	}
	if changed == 0 {
		b.adminReply(ctx, out, chat, b.adminNotFound(ctx))
		return
	}
	details := fmt.Sprintf("{\"via\":\"telegram\",\"chat_id\":%d,\"user_id\":%d,\"days\":%d}", chat, pending.UserID, pending.Days)
	auditErr := b.d.Store.Q.InsertAudit(ctx, db.InsertAuditParams{
		Ts: b.d.Now().Unix(), AdminID: sql.NullInt64{},
		Action: "telegram.admin." + action,
		TargetType: sql.NullString{String: "user", Valid: true},
		TargetID: sql.NullString{String: strconv.FormatInt(pending.UserID, 10), Valid: true},
		Ip: sql.NullString{}, Details: sql.NullString{String: details, Valid: true},
	})
	if auditErr != nil && b.d.Log != nil {
		b.d.Log.Warn("telegram admin action audit failed", "action", action, "user_id", pending.UserID, "err", auditErr)
	}
	b.adminReply(ctx, out, chat, b.adminActionDone(ctx, action, pending.UserID, pending.Days))
}

func (b *Bot) adminCancelAction(ctx context.Context, out *Outbox, chat int64, token string) {
	b.adminActionMu.Lock()
	pending, ok := b.adminActions[token]
	if ok && pending.ChatID == chat {
		delete(b.adminActions, token)
	}
	b.adminActionMu.Unlock()
	if !ok || pending.ChatID != chat {
		b.adminReply(ctx, out, chat, b.adminTokenInvalid(ctx))
		return
	}
	if b.Config(ctx).Lang == "en" {
		b.adminReply(ctx, out, chat, "✅ Pending action cancelled.")
	} else {
		b.adminReply(ctx, out, chat, "✅ Действие отменено.")
	}
}

func (b *Bot) adminActionDone(ctx context.Context, action string, id, days int64) string {
	if b.Config(ctx).Lang == "en" {
		switch action {
		case domain.BulkDisable:
			return fmt.Sprintf("✅ Access disabled for client #%d.", id)
		case domain.BulkEnable:
			return fmt.Sprintf("✅ Access enabled for client #%d.", id)
		case domain.BulkExtend:
			return fmt.Sprintf("✅ Subscription #%d extended by %d days.", id, days)
		case domain.BulkReset:
			return fmt.Sprintf("✅ Traffic reset for client #%d.", id)
		}
	}
	switch action {
	case domain.BulkDisable:
		return fmt.Sprintf("✅ Доступ клиента #%d выключен.", id)
	case domain.BulkEnable:
		return fmt.Sprintf("✅ Доступ клиента #%d включён.", id)
	case domain.BulkExtend:
		return fmt.Sprintf("✅ Подписка клиента #%d продлена на %d дн.", id, days)
	case domain.BulkReset:
		return fmt.Sprintf("✅ Счётчик трафика клиента #%d сброшен.", id)
	}
	return b.adminUsage(ctx)
}

func (b *Bot) adminMutationUnavailable(ctx context.Context) string {
	if b.Config(ctx).Lang == "en" {
		return "Admin actions are unavailable in this panel configuration."
	}
	return "Изменение подписок недоступно в этой конфигурации панели."
}

func (b *Bot) adminTokenInvalid(ctx context.Context) string {
	if b.Config(ctx).Lang == "en" {
		return "Confirmation token is invalid, expired or belongs to another chat. Request a new action."
	}
	return "Токен подтверждения неверен, истёк или выдан другому чату. Запросите действие заново."
}

func (b *Bot) adminNotFound(ctx context.Context) string {
	if b.Config(ctx).Lang == "en" {
		return "Client not found."
	}
	return "Клиент не найден."
}
