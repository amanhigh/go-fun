package tools

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/bitfield/script"
)

// ErrScreenshotAborted is returned when the user cancels a region screenshot
// (e.g., pressing Escape during slurp region selection).
var ErrScreenshotAborted = errors.New("screenshot aborted")

func Screenshot(dir, name string) (err error) {
	var monitor string
	if monitor, err = GetActiveMonitor(); err != nil {
		return
	}
	fullPath := dir + "/" + name
	err = script.Exec(fmt.Sprintf("grim -o %s %s", monitor, fullPath)).Error()
	return
}

func NamedRegionScreenshot(dir, name string) (err error) {
	// Step 1: Run slurp to get the selected region geometry.
	// slurp exits with code 1 when the user cancels (Escape).
	geometry, err := script.Exec("slurp").String()
	if err != nil {
		return ErrScreenshotAborted
	}

	// Step 2: Use the geometry string for grim capture.
	fullPath := dir + "/" + name
	geometry = strings.TrimSpace(geometry)
	err = script.Exec(fmt.Sprintf(`grim -g "%s" "%s"`, geometry, fullPath)).Error()
	return
}

// ResolveWiFiConnection returns the name of the active Wi-Fi connection, or the
// name of the best autoconnect Wi-Fi profile when no Wi-Fi is currently active.
//
// It first inspects device status for a connected Wi-Fi device. If none is
// found, it falls back to connection profiles and selects the autoconnect
// Wi-Fi profile with the highest AUTOCONNECT-PRIORITY. A clear error is
// returned if command execution fails or no suitable Wi-Fi profile exists.
func ResolveWiFiConnection() (string, error) {
	// Step 1: Prefer an already-connected Wi-Fi device.
	out, err := runNMCLI([]string{"-t", "-e", "no", "-f", "TYPE,STATE,CONNECTION", "dev", "status"}, 10*time.Second)
	if err != nil {
		return "", fmt.Errorf("failed to query device status: %w", err)
	}
	for line := range strings.SplitSeq(out, "\n") {
		fields := strings.SplitN(strings.TrimSpace(line), ":", 3)
		if len(fields) < 3 {
			continue
		}
		devType, state, connection := fields[0], fields[1], fields[2]
		if devType == "wifi" && state == "connected" && connection != "" {
			return connection, nil
		}
	}

	// Step 2: Fall back to the highest-priority autoconnect Wi-Fi profile.
	out, err = runNMCLI([]string{"-t", "-e", "no", "-f", "NAME,TYPE,AUTOCONNECT,AUTOCONNECT-PRIORITY", "connection", "show"}, 10*time.Second)
	if err != nil {
		return "", fmt.Errorf("failed to query connection profiles: %w", err)
	}

	bestName := ""
	bestPriority := 0
	found := false
	for line := range strings.SplitSeq(out, "\n") {
		fields := strings.SplitN(strings.TrimSpace(line), ":", 4)
		if len(fields) < 4 {
			continue
		}
		name, connType, autoconnect, priorityStr := fields[0], fields[1], fields[2], fields[3]
		if connType != "802-11-wireless" || autoconnect != "yes" || name == "" {
			continue
		}
		priority, perr := strconv.Atoi(strings.TrimSpace(priorityStr))
		if perr != nil {
			priority = 0
		}
		if !found || priority > bestPriority {
			found = true
			bestName = name
			bestPriority = priority
		}
	}

	if !found {
		return "", errors.New("no suitable autoconnect Wi-Fi profile found")
	}
	return bestName, nil
}

// ConnectionGatewayReachable reports whether the IPv4 gateway of the given
// connection is reachable. It resolves the gateway via nmcli and pings it once.
// It returns false on command failure, an empty gateway, or a failed ping.
func ConnectionGatewayReachable(connectionName string) bool {
	gateway, err := runNMCLI([]string{"-g", "IP4.GATEWAY", "connection", "show", connectionName}, 10*time.Second)
	if err != nil {
		return false
	}
	gateway = strings.TrimSpace(gateway)
	if gateway == "" {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "ping", "-c", "1", "-W", "2", gateway)
	return cmd.Run() == nil
}

// ReconnectConnection performs a targeted up-only recovery of the given
// connection by activating it directly. It does not bring the connection down
// first; instead it relies on nmcli's activation to re-establish the link.
//
// Recovery is attempted up to four times on an approximately one-minute cadence:
// attempt 1 starts immediately, then attempts 2, 3, and 4 begin at +1m, +2m, and
// +3m relative to the first attempt. Each activation is bounded by a per-attempt
// budget sized to absorb the observed ~39s "ssid-not-found" failure: nmcli runs
// with --wait 45 under a slightly larger 50s command context. After every
// successful activation the gateway is probed once via
// ConnectionGatewayReachable; if the gateway is reachable the recovery returns
// success immediately. If activation or gateway verification fails, the next
// scheduled one-minute slot is awaited and the attempt is retried.
//
// After the fourth attempt fails, a descriptive error is returned containing the
// total attempt count and the last error, so callers (e.g., the network monitor)
// do not treat the recovery as successful and can fall back to a heavier restart.
// All command execution is bounded by context timeouts and uses direct argument
// arrays (never a shell). The tools layer does not log; errors are returned for
// callers to handle.
func ReconnectConnection(connectionName string) error {
	const attempts = 4
	const activationWaitSeconds = 45
	const activationBudget = 50 * time.Second

	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	var lastErr error

	for attempt := 1; attempt <= attempts; attempt++ {
		// Wait until the scheduled one-minute slot; attempt 1 is immediate.
		if attempt > 1 {
			<-ticker.C
		}

		ctxUp, cancelUp := context.WithTimeout(context.Background(), activationBudget)
		up := exec.CommandContext(ctxUp, "nmcli", "--wait",
			strconv.Itoa(activationWaitSeconds),
			"connection", "up", "id", connectionName)
		out, err := up.CombinedOutput()

		if err != nil {
			// Preserve the context deadline/cancellation before cancelling the
			// context, so a per-attempt timeout is not masked by the command
			// error. Otherwise keep the nmcli command error itself, and append
			// the trimmed combined nmcli output directly to it.
			if ctxUp.Err() != nil {
				lastErr = ctxUp.Err()
			} else {
				lastErr = err
			}
			if diag := strings.TrimSpace(string(out)); diag != "" {
				lastErr = fmt.Errorf("%w: %s", lastErr, diag)
			}
			cancelUp()
			continue
		}
		cancelUp()

		// The connection reports "up", but verify the gateway actually returned
		// before reporting success. A connection can associate yet leave the
		// gateway unreachable (e.g., no DHCP lease or missing route), so a bare
		// "up" must not be treated as recovered.
		if ConnectionGatewayReachable(connectionName) {
			return nil
		}
		lastErr = errors.New("connection is up but its gateway is unreachable")
	}

	return fmt.Errorf("failed to recover connection %q after %d attempts: %w",
		connectionName, attempts, lastErr)
}

// RestartNetworkManager restarts NetworkManager non-interactively via sudo with
// a bounded timeout. Any error is returned to the caller without logging.
func RestartNetworkManager() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "sudo", "-n", "/usr/bin/systemctl", "restart", "NetworkManager")
	return cmd.Run()
}

// runNMCLI executes nmcli with the supplied argument array (never a shell) under
// a context timeout and returns the trimmed stdout output.
func runNMCLI(args []string, timeout time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	out, err := exec.CommandContext(ctx, "nmcli", args...).Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func PromptText(text string) (result string, err error) {
	result, err = script.Echo(text).Exec("zenity --editable --text-info").String()
	return
}
