package tools

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
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
	fullPath := filepath.Join(dir, name)
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
	fullPath := filepath.Join(dir, name)
	geometry = strings.TrimSpace(geometry)
	err = script.Exec(fmt.Sprintf(`grim -g "%s" "%s"`, geometry, fullPath)).Error()
	return
}

// nmFieldCount is the number of colon-separated fields expected in the tabular
// nmcli output parsed by the Wi-Fi profile resolver below.
const nmFieldCount = 4

// nmActiveFieldCount is the number of colon-separated fields expected in the
// tabular nmcli device-status output parsed by parseActiveWiFiConnection.
const nmActiveFieldCount = 3

// nmCommandTimeout bounds direct nmcli and ping command execution.
const nmCommandTimeout = 10 * time.Second

// nmRestartTimeout bounds the NetworkManager restart command execution.
const nmRestartTimeout = 30 * time.Second

// ResolveWiFiConnection returns the name of the active Wi-Fi connection when one
// is connected, otherwise the highest-priority autoconnect Wi-Fi profile.
//
// It first inspects device status for a connected Wi-Fi device and returns its
// active connection name. Otherwise it falls back to connection profiles and
// selects the autoconnect Wi-Fi profile with the highest AUTOCONNECT-PRIORITY.
// A clear error is returned if command execution fails or no suitable Wi-Fi
// profile exists.
func ResolveWiFiConnection() (connectionName string, err error) {
	// Step 1: Prefer an already-connected Wi-Fi connection.
	out, err := runNMCLI([]string{"-t", "-e", "no", "-f", "TYPE,STATE,CONNECTION", "dev", "status"}, nmCommandTimeout)
	if err != nil {
		return "", fmt.Errorf("failed to query device status: %w", err)
	}

	if connectionName = parseActiveWiFiConnection(out); connectionName != "" {
		return connectionName, nil
	}

	// Step 2: Fall back to the highest-priority autoconnect Wi-Fi profile.
	out, err = runNMCLI([]string{"-t", "-e", "no", "-f", "NAME,TYPE,AUTOCONNECT,AUTOCONNECT-PRIORITY", "connection", "show"}, nmCommandTimeout)
	if err != nil {
		return "", fmt.Errorf("failed to query connection profiles: %w", err)
	}

	connectionName = parseBestAutoConnectProfile(out)
	if connectionName == "" {
		return "", errors.New("no suitable autoconnect Wi-Fi profile found")
	}
	return connectionName, nil
}

// parseActiveWiFiConnection parses `nmcli dev status` tabular output restricted
// to TYPE,STATE,CONNECTION and returns the connection name of the first
// connected Wi-Fi device. An empty result indicates no active Wi-Fi connection.
func parseActiveWiFiConnection(out string) (connectionName string) {
	for line := range strings.SplitSeq(out, "\n") {
		fields := strings.SplitN(strings.TrimSpace(line), ":", nmActiveFieldCount)
		if len(fields) < nmActiveFieldCount {
			continue
		}
		connType, state, connection := fields[0], fields[1], fields[2]
		if connType == "wifi" && state == "connected" && connection != "" {
			return connection
		}
	}
	return ""
}

// parseBestAutoConnectProfile parses `nmcli connection show` tabular output and
// returns the autoconnect 802-11-wireless profile with the highest
// AUTOCONNECT-PRIORITY. Empty names are excluded during parsing, so an empty
// result indicates that no suitable profile exists.
func parseBestAutoConnectProfile(out string) (name string) {
	bestName := ""
	bestPriority := 0
	for line := range strings.SplitSeq(out, "\n") {
		fields := strings.SplitN(strings.TrimSpace(line), ":", nmFieldCount)
		if len(fields) < nmFieldCount {
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
		if bestName == "" || priority > bestPriority {
			bestName = name
			bestPriority = priority
		}
	}
	return bestName
}

// ConnectionGatewayReachable reports whether the IPv4 gateway of the given
// connection is reachable. It resolves the gateway via nmcli and pings it once.
// It returns false on command failure, an empty gateway, or a failed ping.
func ConnectionGatewayReachable(connectionName string) bool {
	gateway, err := runNMCLI([]string{"-g", "IP4.GATEWAY", "connection", "show", connectionName}, nmCommandTimeout)
	if err != nil {
		return false
	}
	gateway = strings.TrimSpace(gateway)
	if gateway == "" {
		return false
	}
	_, err = runCommandWithTimeout(nmCommandTimeout, "ping", "-c", "1", "-W", "2", gateway)
	return err == nil
}

// RestartNetworkManager restarts NetworkManager non-interactively via sudo with
// a bounded timeout. Any error is returned to the caller without logging.
func RestartNetworkManager() error {
	if _, err := runCommandWithTimeout(nmRestartTimeout, "sudo", "-n", "/usr/bin/systemctl", "restart", "NetworkManager"); err != nil {
		return fmt.Errorf("failed to restart NetworkManager: %w", err)
	}
	return nil
}

// runCommand executes an external command (never a shell) under a context
// timeout and returns the trimmed stdout output. Errors are wrapped with the
// command name for easier diagnosis.
func runCommandWithTimeout(timeout time.Duration, name string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	out, err := exec.CommandContext(ctx, name, args...).Output()
	if err != nil {
		return "", fmt.Errorf("%s command failed: %w", name, err)
	}
	return strings.TrimSpace(string(out)), nil
}

// runNMCLI executes nmcli with the supplied argument array under a context
// timeout and returns the trimmed stdout output.
func runNMCLI(args []string, timeout time.Duration) (string, error) {
	return runCommandWithTimeout(timeout, "nmcli", args...)
}

func PromptText(text string) (result string, err error) {
	result, err = script.Echo(text).Exec("zenity --editable --text-info").String()
	return
}
