package channels

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/mymmrac/telego"
	"github.com/mymmrac/telego/telegohandler"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"

	"github.com/sipeed/picoclaw/pkg/bus"
	"github.com/sipeed/picoclaw/pkg/config"
	"github.com/sipeed/picoclaw/pkg/logger"
	"github.com/sipeed/picoclaw/pkg/utils"
	"github.com/sipeed/picoclaw/pkg/voice"
)

const defaultCoolifyContext = "mycoolify"

type TelegramChannel struct {
	*BaseChannel
	bot          *telego.Bot
	commands     TelegramCommander
	config       *config.Config
	chatIDs      map[string]int64
	transcriber  *voice.GroqTranscriber
	placeholders sync.Map // chatID -> messageID
	stopThinking sync.Map // chatID -> thinkingCancel
}

type thinkingCancel struct {
	fn context.CancelFunc
}

func (c *thinkingCancel) Cancel() {
	if c != nil && c.fn != nil {
		c.fn()
	}
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func collectParentDirs(start string, maxDepth int) []string {
	if start == "" || maxDepth <= 0 {
		return nil
	}

	dirs := make([]string, 0, maxDepth)
	current := filepath.Clean(start)
	for i := 0; i < maxDepth; i++ {
		dirs = append(dirs, current)
		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}
	return dirs
}

func uniqueNonEmpty(values []string) []string {
	out := make([]string, 0, len(values))
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			continue
		}
		if _, ok := seen[trimmed]; ok {
			continue
		}
		seen[trimmed] = struct{}{}
		out = append(out, trimmed)
	}
	return out
}

func resolvePicoclawRoot() (string, error) {
	candidates := make([]string, 0, 12)
	if envRoot := strings.TrimSpace(os.Getenv("PICOCLAW_ROOT")); envRoot != "" {
		candidates = append(candidates, envRoot)
	}

	if cwd, err := os.Getwd(); err == nil {
		candidates = append(candidates, collectParentDirs(cwd, 6)...)
	}

	if _, sourceFile, _, ok := runtime.Caller(0); ok {
		sourceRoot := filepath.Clean(filepath.Join(filepath.Dir(sourceFile), "..", ".."))
		candidates = append(candidates, collectParentDirs(sourceRoot, 3)...)
	}

	for _, candidate := range uniqueNonEmpty(candidates) {
		if fileExists(filepath.Join(candidate, "workspace", "bin", "me.sh")) {
			return candidate, nil
		}
	}

	return "", fmt.Errorf("unable to locate picoclaw workspace root; set PICOCLAW_ROOT")
}

func resolveWorkspaceRunner() (string, error) {
	root, err := resolvePicoclawRoot()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, "workspace", "bin", "me.sh"), nil
}

func shellQuote(value string) string {
	if value == "" {
		return "''"
	}
	return "'" + strings.ReplaceAll(value, "'", `'"'"'`) + "'"
}

func buildMaxExtractScriptCommand(script string, args ...string) (string, error) {
	wrapper, err := resolveWorkspaceRunner()
	if err != nil {
		return "", err
	}

	parts := []string{
		"MAXEXTRACT_USE_SSH=1",
		"MAXEXTRACT_OUTPUT_FORMAT=telegram",
		shellQuote(wrapper),
		shellQuote(script),
	}
	for _, arg := range args {
		parts = append(parts, shellQuote(arg))
	}
	return strings.Join(parts, " "), nil
}

func isValidMode(value string, allowAll bool) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "paper", "live":
		return true
	case "all":
		return allowAll
	default:
		return false
	}
}

func parseFleetModeArgs(cmdName, args string) (string, error) {
	fields := strings.Fields(args)
	if len(fields) == 0 {
		return "all", nil
	}
	if len(fields) > 1 {
		return "", fmt.Errorf("Uso: /%s [all|paper|live]", cmdName)
	}
	mode := strings.ToLower(fields[0])
	if !isValidMode(mode, true) {
		return "", fmt.Errorf("Mode non valido: %s. Usa all, paper o live.", fields[0])
	}
	return mode, nil
}

func parseFleetModeAndDaysArgs(cmdName, args string) (string, string, error) {
	fields := strings.Fields(args)
	mode := "all"
	days := "auto"

	switch len(fields) {
	case 0:
		return mode, days, nil
	case 1:
		if isValidMode(fields[0], true) {
			return strings.ToLower(fields[0]), days, nil
		}
		return mode, fields[0], nil
	case 2:
		if !isValidMode(fields[0], true) {
			return "", "", fmt.Errorf("Mode non valido: %s. Usa all, paper o live.", fields[0])
		}
		return strings.ToLower(fields[0]), fields[1], nil
	default:
		return "", "", fmt.Errorf("Uso: /%s [all|paper|live] [days]", cmdName)
	}
}

func parseBotSelectorArgs(cmdName, args string, allowDays bool) (string, string, string, string, error) {
	fields := strings.Fields(args)
	if len(fields) < 3 {
		usage := fmt.Sprintf("Uso: /%s <mode> <strategy> <market>", cmdName)
		if allowDays {
			usage += " [days]"
		}
		return "", "", "", "", fmt.Errorf("%s\nEsempio: /%s paper ema-until-expiry btc-5m", usage, cmdName)
	}
	if len(fields) > 3 && !allowDays {
		return "", "", "", "", fmt.Errorf("Uso: /%s <mode> <strategy> <market>", cmdName)
	}
	if len(fields) > 4 {
		return "", "", "", "", fmt.Errorf("Uso: /%s <mode> <strategy> <market> [days]", cmdName)
	}

	mode := strings.ToLower(fields[0])
	if !isValidMode(mode, false) {
		return "", "", "", "", fmt.Errorf("Mode non valido: %s. Usa paper o live.", fields[0])
	}

	days := "auto"
	if allowDays && len(fields) == 4 {
		days = fields[3]
	}
	return mode, fields[1], fields[2], days, nil
}

func isDirectShellSlashCommand(cmd string) bool {
	switch cmd {
	case "bots", "inventory", "digest", "pnl", "truth_slo", "bot", "health", "roi", "run":
		return true
	default:
		return false
	}
}

func NewTelegramChannel(cfg *config.Config, bus *bus.MessageBus) (*TelegramChannel, error) {
	var opts []telego.BotOption
	telegramCfg := cfg.Channels.Telegram

	if telegramCfg.Proxy != "" {
		proxyURL, parseErr := url.Parse(telegramCfg.Proxy)
		if parseErr != nil {
			return nil, fmt.Errorf("invalid proxy URL %q: %w", telegramCfg.Proxy, parseErr)
		}
		opts = append(opts, telego.WithHTTPClient(&http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			},
		}))
	} else if os.Getenv("HTTP_PROXY") != "" || os.Getenv("HTTPS_PROXY") != "" {
		// Use environment proxy if configured
		opts = append(opts, telego.WithHTTPClient(&http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
			},
		}))
	}

	bot, err := telego.NewBot(telegramCfg.Token, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create telegram bot: %w", err)
	}

	base := NewBaseChannel("telegram", telegramCfg, bus, telegramCfg.AllowFrom)

	return &TelegramChannel{
		BaseChannel:  base,
		commands:     NewTelegramCommands(bot, cfg),
		bot:          bot,
		config:       cfg,
		chatIDs:      make(map[string]int64),
		transcriber:  nil,
		placeholders: sync.Map{},
		stopThinking: sync.Map{},
	}, nil
}

func (c *TelegramChannel) SetTranscriber(transcriber *voice.GroqTranscriber) {
	c.transcriber = transcriber
}

func (c *TelegramChannel) Start(ctx context.Context) error {
	logger.InfoC("telegram", "Starting Telegram bot (polling mode)...")

	updates, err := c.bot.UpdatesViaLongPolling(ctx, &telego.GetUpdatesParams{
		Timeout: 30,
	})
	if err != nil {
		return fmt.Errorf("failed to start long polling: %w", err)
	}

	bh, err := telegohandler.NewBotHandler(c.bot, updates)
	if err != nil {
		return fmt.Errorf("failed to create bot handler: %w", err)
	}

	bh.HandleMessage(func(ctx *th.Context, message telego.Message) error {
		c.commands.Help(ctx, message)
		return nil
	}, th.CommandEqual("help"))
	bh.HandleMessage(func(ctx *th.Context, message telego.Message) error {
		return c.commands.Start(ctx, message)
	}, th.CommandEqual("start"))

	bh.HandleMessage(func(ctx *th.Context, message telego.Message) error {
		return c.commands.Show(ctx, message)
	}, th.CommandEqual("show"))

	bh.HandleMessage(func(ctx *th.Context, message telego.Message) error {
		return c.commands.List(ctx, message)
	}, th.CommandEqual("list"))

	bh.HandleMessage(func(ctx *th.Context, message telego.Message) error {
		return c.handleOpsSlashCommand(ctx, &message, "bots")
	}, th.CommandEqual("bots"))
	bh.HandleMessage(func(ctx *th.Context, message telego.Message) error {
		return c.handleOpsSlashCommand(ctx, &message, "bots")
	}, th.CommandEqual("vedibots"))

	bh.HandleMessage(func(ctx *th.Context, message telego.Message) error {
		return c.handleOpsSlashCommand(ctx, &message, "inventory")
	}, th.CommandEqual("inventory"))

	bh.HandleMessage(func(ctx *th.Context, message telego.Message) error {
		return c.handleOpsSlashCommand(ctx, &message, "digest")
	}, th.CommandEqual("digest"))

	bh.HandleMessage(func(ctx *th.Context, message telego.Message) error {
		return c.handleOpsSlashCommand(ctx, &message, "pnl")
	}, th.CommandEqual("pnl"))

	bh.HandleMessage(func(ctx *th.Context, message telego.Message) error {
		return c.handleOpsSlashCommand(ctx, &message, "truth_slo")
	}, th.CommandEqual("truth_slo"))
	bh.HandleMessage(func(ctx *th.Context, message telego.Message) error {
		return c.handleOpsSlashCommand(ctx, &message, "truth_slo")
	}, th.CommandEqual("truthslo"))

	bh.HandleMessage(func(ctx *th.Context, message telego.Message) error {
		return c.handleOpsSlashCommand(ctx, &message, "bot")
	}, th.CommandEqual("bot"))
	bh.HandleMessage(func(ctx *th.Context, message telego.Message) error {
		return c.handleOpsSlashCommand(ctx, &message, "bot")
	}, th.CommandEqual("vedibot"))

	bh.HandleMessage(func(ctx *th.Context, message telego.Message) error {
		return c.handleOpsSlashCommand(ctx, &message, "health")
	}, th.CommandEqual("health"))

	bh.HandleMessage(func(ctx *th.Context, message telego.Message) error {
		return c.handleOpsSlashCommand(ctx, &message, "roi")
	}, th.CommandEqual("roi"))

	bh.HandleMessage(func(ctx *th.Context, message telego.Message) error {
		return c.handleOpsSlashCommand(ctx, &message, "memclear")
	}, th.CommandEqual("memclear"))
	bh.HandleMessage(func(ctx *th.Context, message telego.Message) error {
		return c.handleOpsSlashCommand(ctx, &message, "memclear")
	}, th.CommandEqual("memoria"))
	bh.HandleMessage(func(ctx *th.Context, message telego.Message) error {
		return c.handleOpsSlashCommand(ctx, &message, "memclear")
	}, th.CommandEqual("cancellamemoria"))

	bh.HandleMessage(func(ctx *th.Context, message telego.Message) error {
		return c.handleOpsSlashCommand(ctx, &message, "run")
	}, th.CommandEqual("run"))

	bh.HandleMessage(func(ctx *th.Context, message telego.Message) error {
		return c.handleMessage(ctx, &message)
	}, th.AnyMessage())

	c.setRunning(true)
	logger.InfoCF("telegram", "Telegram bot connected", map[string]any{
		"username": c.bot.Username(),
	})

	go bh.Start()

	go func() {
		<-ctx.Done()
		bh.Stop()
	}()

	return nil
}

func (c *TelegramChannel) Stop(ctx context.Context) error {
	logger.InfoC("telegram", "Stopping Telegram bot...")
	c.setRunning(false)
	return nil
}

func (c *TelegramChannel) Send(ctx context.Context, msg bus.OutboundMessage) error {
	if !c.IsRunning() {
		return fmt.Errorf("telegram bot not running")
	}

	chatID, err := parseChatID(msg.ChatID)
	if err != nil {
		return fmt.Errorf("invalid chat ID: %w", err)
	}

	// Stop thinking animation
	if stop, ok := c.stopThinking.Load(msg.ChatID); ok {
		if cf, ok := stop.(*thinkingCancel); ok && cf != nil {
			cf.Cancel()
		}
		c.stopThinking.Delete(msg.ChatID)
	}

	htmlContent := markdownToTelegramHTML(msg.Content)

	// Try to edit placeholder
	if pID, ok := c.placeholders.Load(msg.ChatID); ok {
		c.placeholders.Delete(msg.ChatID)
		editMsg := tu.EditMessageText(tu.ID(chatID), pID.(int), htmlContent)
		editMsg.ParseMode = telego.ModeHTML

		if _, err = c.bot.EditMessageText(ctx, editMsg); err == nil {
			return nil
		}
		// Fallback to new message if edit fails
	}

	tgMsg := tu.Message(tu.ID(chatID), htmlContent)
	tgMsg.ParseMode = telego.ModeHTML

	if _, err = c.bot.SendMessage(ctx, tgMsg); err != nil {
		logger.ErrorCF("telegram", "HTML parse failed, falling back to plain text", map[string]any{
			"error": err.Error(),
		})
		tgMsg.ParseMode = ""
		tgMsg.Text = msg.Content
		_, err = c.bot.SendMessage(ctx, tgMsg)
		return err
	}

	return nil
}

func (c *TelegramChannel) handleMessage(ctx context.Context, message *telego.Message) error {
	if message == nil {
		return fmt.Errorf("message is nil")
	}

	user := message.From
	if user == nil {
		return fmt.Errorf("message sender (user) is nil")
	}

	senderID := fmt.Sprintf("%d", user.ID)
	if user.Username != "" {
		senderID = fmt.Sprintf("%d|%s", user.ID, user.Username)
	}

	// check allowlist to avoid downloading attachments for rejected users
	if !c.IsAllowed(senderID) {
		logger.DebugCF("telegram", "Message rejected by allowlist", map[string]any{
			"user_id": senderID,
		})
		return nil
	}

	chatID := message.Chat.ID
	c.chatIDs[senderID] = chatID

	content := ""
	mediaPaths := []string{}
	localFiles := []string{} // track local files that need cleanup

	// ensure temp files are cleaned up when function returns
	defer func() {
		for _, file := range localFiles {
			if err := os.Remove(file); err != nil {
				logger.DebugCF("telegram", "Failed to cleanup temp file", map[string]any{
					"file":  file,
					"error": err.Error(),
				})
			}
		}
	}()

	if message.Text != "" {
		content += message.Text
	}

	if message.Caption != "" {
		if content != "" {
			content += "\n"
		}
		content += message.Caption
	}

	if len(message.Photo) > 0 {
		photo := message.Photo[len(message.Photo)-1]
		photoPath := c.downloadPhoto(ctx, photo.FileID)
		if photoPath != "" {
			localFiles = append(localFiles, photoPath)
			mediaPaths = append(mediaPaths, photoPath)
			if content != "" {
				content += "\n"
			}
			content += "[image: photo]"
		}
	}

	if message.Voice != nil {
		voicePath := c.downloadFile(ctx, message.Voice.FileID, ".ogg")
		if voicePath != "" {
			localFiles = append(localFiles, voicePath)
			mediaPaths = append(mediaPaths, voicePath)

			transcribedText := ""
			if c.transcriber != nil && c.transcriber.IsAvailable() {
				transcriberCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
				defer cancel()

				result, err := c.transcriber.Transcribe(transcriberCtx, voicePath)
				if err != nil {
					logger.ErrorCF("telegram", "Voice transcription failed", map[string]any{
						"error": err.Error(),
						"path":  voicePath,
					})
					transcribedText = "[voice (transcription failed)]"
				} else {
					transcribedText = fmt.Sprintf("[voice transcription: %s]", result.Text)
					logger.InfoCF("telegram", "Voice transcribed successfully", map[string]any{
						"text": result.Text,
					})
				}
			} else {
				transcribedText = "[voice]"
			}

			if content != "" {
				content += "\n"
			}
			content += transcribedText
		}
	}

	if message.Audio != nil {
		audioPath := c.downloadFile(ctx, message.Audio.FileID, ".mp3")
		if audioPath != "" {
			localFiles = append(localFiles, audioPath)
			mediaPaths = append(mediaPaths, audioPath)
			if content != "" {
				content += "\n"
			}
			content += "[audio]"
		}
	}

	if message.Document != nil {
		docPath := c.downloadFile(ctx, message.Document.FileID, "")
		if docPath != "" {
			localFiles = append(localFiles, docPath)
			mediaPaths = append(mediaPaths, docPath)
			if content != "" {
				content += "\n"
			}
			content += "[file]"
		}
	}

	if content == "" {
		content = "[empty message]"
	}

	logger.DebugCF("telegram", "Received message", map[string]any{
		"sender_id": senderID,
		"chat_id":   fmt.Sprintf("%d", chatID),
		"preview":   utils.Truncate(content, 50),
	})

	// Thinking indicator
	err := c.bot.SendChatAction(ctx, tu.ChatAction(tu.ID(chatID), telego.ChatActionTyping))
	if err != nil {
		logger.ErrorCF("telegram", "Failed to send chat action", map[string]any{
			"error": err.Error(),
		})
	}

	// Stop any previous thinking animation
	chatIDStr := fmt.Sprintf("%d", chatID)
	if prevStop, ok := c.stopThinking.Load(chatIDStr); ok {
		if cf, ok := prevStop.(*thinkingCancel); ok && cf != nil {
			cf.Cancel()
		}
	}

	// Create cancel function for thinking state
	_, thinkCancel := context.WithTimeout(ctx, 5*time.Minute)
	c.stopThinking.Store(chatIDStr, &thinkingCancel{fn: thinkCancel})

	pMsg, err := c.bot.SendMessage(ctx, tu.Message(tu.ID(chatID), "Thinking... 💭"))
	if err == nil {
		pID := pMsg.MessageID
		c.placeholders.Store(chatIDStr, pID)
	}

	peerKind := "direct"
	peerID := fmt.Sprintf("%d", user.ID)
	if message.Chat.Type != "private" {
		peerKind = "group"
		peerID = fmt.Sprintf("%d", chatID)
	}

	metadata := map[string]string{
		"message_id": fmt.Sprintf("%d", message.MessageID),
		"user_id":    fmt.Sprintf("%d", user.ID),
		"username":   user.Username,
		"first_name": user.FirstName,
		"is_group":   fmt.Sprintf("%t", message.Chat.Type != "private"),
		"peer_kind":  peerKind,
		"peer_id":    peerID,
	}

	c.HandleMessage(fmt.Sprintf("%d", user.ID), fmt.Sprintf("%d", chatID), content, mediaPaths, metadata)
	return nil
}

func commandArgsFromText(text string) string {
	parts := strings.SplitN(strings.TrimSpace(text), " ", 2)
	if len(parts) < 2 {
		return ""
	}
	return strings.TrimSpace(parts[1])
}

func (c *TelegramChannel) sendSlashText(ctx context.Context, message *telego.Message, text string) error {
	params := &telego.SendMessageParams{
		ChatID: telego.ChatID{ID: message.Chat.ID},
		Text:   markdownToTelegramHTML(text),
		ParseMode: telego.ModeHTML,
		ReplyParameters: &telego.ReplyParameters{
			MessageID: message.MessageID,
		},
	}
	if _, err := c.bot.SendMessage(ctx, params); err == nil {
		return nil
	}

	params.Text = text
	params.ParseMode = ""
	_, err := c.bot.SendMessage(ctx, params)
	return err
}

func isSafeRunCommand(raw string) bool {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return false
	}
	if strings.ContainsAny(trimmed, "`") || strings.Contains(trimmed, "$(") {
		return false
	}
	// Keep slash-command runner simple and safe.
	disallowed := []string{"&&", "||", ";", "|", ">", "<"}
	for _, token := range disallowed {
		if strings.Contains(trimmed, token) {
			return false
		}
	}
	// Allow only Coolify and known MaxExtract ops scripts.
	if strings.HasPrefix(trimmed, "coolify ") {
		return true
	}
	if strings.HasPrefix(trimmed, "me_bot") || strings.HasPrefix(trimmed, "me_bots") {
		return true
	}
	if strings.HasPrefix(trimmed, "./workspace/bin/me.sh ") ||
		strings.HasPrefix(trimmed, "workspace/bin/me.sh ") ||
		strings.Contains(trimmed, "/workspace/bin/me.sh ") {
		return true
	}
	if strings.Contains(trimmed, "/scripts/me_bot") || strings.Contains(trimmed, "/scripts/me_bots") {
		return true
	}
	return false
}

func normalizeRunCommand(raw string) (string, error) {
	trimmed := strings.TrimSpace(raw)
	if strings.HasPrefix(trimmed, "me_bot") || strings.HasPrefix(trimmed, "me_bots") {
		wrapper, err := resolveWorkspaceRunner()
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%s %s", shellQuote(wrapper), trimmed), nil
	}
	if strings.HasPrefix(trimmed, "./workspace/bin/me.sh ") ||
		strings.HasPrefix(trimmed, "workspace/bin/me.sh ") ||
		strings.Contains(trimmed, "/workspace/bin/me.sh ") {
		wrapper, err := resolveWorkspaceRunner()
		if err != nil {
			return "", err
		}
		re := regexp.MustCompile(`^(\./)?workspace/bin/me\.sh\s+`)
		normalized := re.ReplaceAllString(trimmed, fmt.Sprintf("%s ", shellQuote(wrapper)))
		if normalized != trimmed {
			return normalized, nil
		}
	}
	return trimmed, nil
}

func (c *TelegramChannel) buildOpsShellCommand(cmd, args string) (string, error) {
	switch cmd {
	case "bots":
		mode, err := parseFleetModeArgs("bots", args)
		if err != nil {
			return "", err
		}
		return buildMaxExtractScriptCommand("me_bots_api_state.sh", "--context", defaultCoolifyContext, "--mode", mode)
	case "inventory":
		mode, err := parseFleetModeArgs("inventory", args)
		if err != nil {
			return "", err
		}
		return buildMaxExtractScriptCommand("me_bots_inventory.sh", "--context", defaultCoolifyContext, "--mode", mode)
	case "digest":
		mode, days, err := parseFleetModeAndDaysArgs("digest", args)
		if err != nil {
			return "", err
		}
		return buildMaxExtractScriptCommand("me_bots_digest.sh", "--context", defaultCoolifyContext, "--mode", mode, "--days", days)
	case "pnl":
		mode, err := parseFleetModeArgs("pnl", args)
		if err != nil {
			return "", err
		}
		return buildMaxExtractScriptCommand("me_bots_pnl.sh", "--context", defaultCoolifyContext, "--mode", mode)
	case "truth_slo":
		mode, err := parseFleetModeArgs("truth_slo", args)
		if err != nil {
			return "", err
		}
		return buildMaxExtractScriptCommand("me_bots_truth_slo.sh", "--context", defaultCoolifyContext, "--mode", mode)
	case "bot":
		mode, strategy, market, days, err := parseBotSelectorArgs("bot", args, true)
		if err != nil {
			return "", err
		}
		return buildMaxExtractScriptCommand(
			"me_bot_report.sh",
			"--context", defaultCoolifyContext,
			"--mode", mode,
			"--strategy", strategy,
			"--market", market,
			"--days", days,
		)
	case "health":
		mode, strategy, market, _, err := parseBotSelectorArgs("health", args, false)
		if err != nil {
			return "", err
		}
		return buildMaxExtractScriptCommand(
			"me_bot_health.sh",
			"--context", defaultCoolifyContext,
			"--mode", mode,
			"--strategy", strategy,
			"--market", market,
		)
	case "roi":
		mode, strategy, market, days, err := parseBotSelectorArgs("roi", args, true)
		if err != nil {
			return "", err
		}
		return buildMaxExtractScriptCommand(
			"me_bot_roi.sh",
			"--context", defaultCoolifyContext,
			"--mode", mode,
			"--strategy", strategy,
			"--market", market,
			"--days", days,
		)
	case "run":
		if args == "" {
			return "", fmt.Errorf("Uso: /run <command>\nConsentiti: coolify ..., me_bot..., me_bots..., workspace/bin/me.sh ...")
		}
		if !isSafeRunCommand(args) {
			return "", fmt.Errorf("Comando bloccato (unsafe).\nPrefissi consentiti: coolify, me_bot, me_bots, workspace/bin/me.sh")
		}
		return normalizeRunCommand(args)
	default:
		return "", fmt.Errorf("Unsupported command: %s", cmd)
	}
}

func truncateForTelegram(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen] + "\n\n[output truncated]"
}

func (c *TelegramChannel) runOpsShellCommand(ctx context.Context, cmdText string) string {
	runCtx, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	cmd := exec.CommandContext(runCtx, "/bin/sh", "-lc", cmdText)
	out, err := cmd.CombinedOutput()
	result := strings.TrimSpace(string(out))

	if runCtx.Err() == context.DeadlineExceeded {
		return "Comando timeout (120s). Riprova con scope piu piccolo."
	}

	if err != nil {
		if result == "" {
			result = err.Error()
		}
		return truncateForTelegram("Comando fallito:\n"+result, 3900)
	}

	if result == "" {
		result = "Comando eseguito, nessun output."
	}
	return truncateForTelegram(result, 3900)
}

func (c *TelegramChannel) buildOpsSlashPrompt(cmd, args string) (string, error) {
	switch cmd {
	case "memclear":
		return "Pulisci la memoria operativa del workspace corrente: svuota memory/MEMORY.md, aggiungi nota nel daily note di oggi con timestamp e motivo 'manual reset da /memclear', poi conferma in 3 bullet cosa hai cancellato.", nil
	case "run":
		if args == "" {
			return "", fmt.Errorf("Uso: /run <command>\nConsentiti: coolify ..., me_bot..., me_bots..., workspace/bin/me.sh ...")
		}
		if !isSafeRunCommand(args) {
			return "", fmt.Errorf("Comando bloccato (unsafe).\nPrefissi consentiti: coolify, me_bot, me_bots, workspace/bin/me.sh")
		}
		normalized, err := normalizeRunCommand(args)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("Esegui questo comando shell e restituisci output compatto in formato Telegram:\n%s", normalized), nil
	default:
		return "", fmt.Errorf("Unsupported command: %s", cmd)
	}
}

func (c *TelegramChannel) handleOpsSlashCommand(ctx context.Context, message *telego.Message, cmd string) error {
	args := commandArgsFromText(message.Text)

	if isDirectShellSlashCommand(cmd) {
		shellCmd, err := c.buildOpsShellCommand(cmd, args)
		if err != nil {
			return c.sendSlashText(ctx, message, err.Error())
		}
		return c.sendSlashText(ctx, message, c.runOpsShellCommand(ctx, shellCmd))
	}

	prompt, err := c.buildOpsSlashPrompt(cmd, args)
	if err != nil {
		return c.sendSlashText(ctx, message, err.Error())
	}
	message.Text = prompt
	return c.handleMessage(ctx, message)
}

func (c *TelegramChannel) downloadPhoto(ctx context.Context, fileID string) string {
	file, err := c.bot.GetFile(ctx, &telego.GetFileParams{FileID: fileID})
	if err != nil {
		logger.ErrorCF("telegram", "Failed to get photo file", map[string]any{
			"error": err.Error(),
		})
		return ""
	}

	return c.downloadFileWithInfo(file, ".jpg")
}

func (c *TelegramChannel) downloadFileWithInfo(file *telego.File, ext string) string {
	if file.FilePath == "" {
		return ""
	}

	url := c.bot.FileDownloadURL(file.FilePath)
	logger.DebugCF("telegram", "File URL", map[string]any{"url": url})

	// Use FilePath as filename for better identification
	filename := file.FilePath + ext
	return utils.DownloadFile(url, filename, utils.DownloadOptions{
		LoggerPrefix: "telegram",
	})
}

func (c *TelegramChannel) downloadFile(ctx context.Context, fileID, ext string) string {
	file, err := c.bot.GetFile(ctx, &telego.GetFileParams{FileID: fileID})
	if err != nil {
		logger.ErrorCF("telegram", "Failed to get file", map[string]any{
			"error": err.Error(),
		})
		return ""
	}

	return c.downloadFileWithInfo(file, ext)
}

func parseChatID(chatIDStr string) (int64, error) {
	var id int64
	_, err := fmt.Sscanf(chatIDStr, "%d", &id)
	return id, err
}

func markdownToTelegramHTML(text string) string {
	if text == "" {
		return ""
	}

	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")

	codeBlocks := extractCodeBlocks(text)
	text = codeBlocks.text

	inlineCodes := extractInlineCodes(text)
	text = inlineCodes.text

	text = regexp.MustCompile(`(?m)^#{1,6}\s+(.+)$`).ReplaceAllString(text, "$1")

	text = regexp.MustCompile(`(?m)^>\s*(.*)$`).ReplaceAllString(text, "$1")

	text = escapeHTML(text)

	text = replaceMarkdownLinks(text)

	text = regexp.MustCompile(`\*\*(.+?)\*\*`).ReplaceAllString(text, "<b>$1</b>")

	text = regexp.MustCompile(`__(.+?)__`).ReplaceAllString(text, "<b>$1</b>")

	reItalic := regexp.MustCompile(`_([^_]+)_`)
	text = reItalic.ReplaceAllStringFunc(text, func(s string) string {
		match := reItalic.FindStringSubmatch(s)
		if len(match) < 2 {
			return s
		}
		return "<i>" + match[1] + "</i>"
	})

	text = regexp.MustCompile(`~~(.+?)~~`).ReplaceAllString(text, "<s>$1</s>")

	text = regexp.MustCompile(`(?m)^[-*]\s+(.*)$`).ReplaceAllString(text, "• $1")
	text = regexp.MustCompile(`(?m)^\d+\.\s+(.*)$`).ReplaceAllString(text, "• $1")

	for i, code := range inlineCodes.codes {
		escaped := escapeHTML(code)
		text = strings.ReplaceAll(text, fmt.Sprintf("\x00IC%d\x00", i), fmt.Sprintf("<code>%s</code>", escaped))
	}

	for i, code := range codeBlocks.codes {
		escaped := escapeHTML(code)
		text = strings.ReplaceAll(
			text,
			fmt.Sprintf("\x00CB%d\x00", i),
			fmt.Sprintf("<pre><code>%s</code></pre>", escaped),
		)
	}

	return text
}

type codeBlockMatch struct {
	text  string
	codes []string
}

func extractCodeBlocks(text string) codeBlockMatch {
	re := regexp.MustCompile("```[\\w]*\\n?([\\s\\S]*?)```")
	matches := re.FindAllStringSubmatch(text, -1)

	codes := make([]string, 0, len(matches))
	for _, match := range matches {
		codes = append(codes, match[1])
	}

	i := 0
	text = re.ReplaceAllStringFunc(text, func(m string) string {
		placeholder := fmt.Sprintf("\x00CB%d\x00", i)
		i++
		return placeholder
	})

	return codeBlockMatch{text: text, codes: codes}
}

type inlineCodeMatch struct {
	text  string
	codes []string
}

func extractInlineCodes(text string) inlineCodeMatch {
	re := regexp.MustCompile("`([^`]+)`")
	matches := re.FindAllStringSubmatch(text, -1)

	codes := make([]string, 0, len(matches))
	for _, match := range matches {
		codes = append(codes, match[1])
	}

	i := 0
	text = re.ReplaceAllStringFunc(text, func(m string) string {
		placeholder := fmt.Sprintf("\x00IC%d\x00", i)
		i++
		return placeholder
	})

	return inlineCodeMatch{text: text, codes: codes}
}

func escapeHTML(text string) string {
	text = strings.ReplaceAll(text, "&", "&amp;")
	text = strings.ReplaceAll(text, "<", "&lt;")
	text = strings.ReplaceAll(text, ">", "&gt;")
	return text
}

func replaceMarkdownLinks(text string) string {
	re := regexp.MustCompile(`\[([^\]]+)\]\(([^)\s]+)\)`)
	return re.ReplaceAllStringFunc(text, func(m string) string {
		parts := re.FindStringSubmatch(m)
		if len(parts) != 3 {
			return m
		}
		label := parts[1]
		href := escapeHTMLAttr(parts[2])
		return fmt.Sprintf(`<a href="%s">%s</a>`, href, label)
	})
}

func escapeHTMLAttr(text string) string {
	text = strings.ReplaceAll(text, "&", "&amp;")
	text = strings.ReplaceAll(text, `"`, "&quot;")
	text = strings.ReplaceAll(text, "<", "&lt;")
	text = strings.ReplaceAll(text, ">", "&gt;")
	return text
}
