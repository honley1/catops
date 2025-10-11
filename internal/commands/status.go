package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	constants "catops/config"
	"catops/internal/config"
	"catops/internal/metrics"
	"catops/internal/process"
	"catops/internal/ui"
	"catops/pkg/utils"
)

// NewStatusCmd creates the status command
func NewStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Display current system metrics and alert thresholds",
		Long: `Display real-time system information including:
  • System Information (Hostname, OS, IP, Uptime)
  • Current Metrics (CPU, Memory, Disk, HTTPS Connections)
  • Alert Thresholds (configured limits for alerts)

Examples:
  catops status          # Show all system information`,
		Run: func(cmd *cobra.Command, args []string) {
			ui.PrintHeader()

			// Load configuration
			cfg, err := config.LoadConfig()
			if err != nil {
				ui.PrintStatus("error", "Failed to load configuration")
				ui.PrintStatus("info", "Using default thresholds")
				cfg = &config.Config{
					CPUThreshold:  constants.DEFAULT_CPU_THRESHOLD,
					MemThreshold:  constants.DEFAULT_MEMORY_THRESHOLD,
					DiskThreshold: constants.DEFAULT_DISK_THRESHOLD,
				}
			}

			// get system information
			hostname, _ := os.Hostname()
			metrics, err := metrics.GetMetrics()
			if err != nil {
				ui.PrintStatus("error", fmt.Sprintf("Error getting metrics: %v", err))
				return
			}

			// system information section
			ui.PrintSection("System Information")
			systemData := map[string]string{
				"Hostname": hostname,
				"OS":       metrics.OSName,
				"IP":       metrics.IPAddress,
				"Uptime":   metrics.Uptime,
			}
			fmt.Print(ui.CreateBeautifulList(systemData))
			ui.PrintSectionEnd()

			// timestamp section
			ui.PrintSection("Timestamp")
			timestampData := map[string]string{
				"Current Time": metrics.Timestamp,
			}
			fmt.Print(ui.CreateBeautifulList(timestampData))
			ui.PrintSectionEnd()

			// metrics section
			ui.PrintSection("Current Metrics")
			metricsData := map[string]string{
				"CPU Usage":         fmt.Sprintf("%s (%d cores, %d active)", utils.FormatPercentage(metrics.CPUUsage), metrics.CPUDetails.Total, metrics.CPUDetails.Used),
				"Memory Usage":      fmt.Sprintf("%s (%s / %s)", utils.FormatPercentage(metrics.MemoryUsage), utils.FormatBytes(metrics.MemoryDetails.Used*1024), utils.FormatBytes(metrics.MemoryDetails.Total*1024)),
				"Disk Usage":        fmt.Sprintf("%s (%s / %s)", utils.FormatPercentage(metrics.DiskUsage), utils.FormatBytes(metrics.DiskDetails.Used*1024), utils.FormatBytes(metrics.DiskDetails.Total*1024)),
				"HTTPS Connections": utils.FormatNumber(metrics.HTTPSRequests),
				"IOPS":              utils.FormatNumber(metrics.IOPS),
				"I/O Wait":          utils.FormatPercentage(metrics.IOWait),
			}
			fmt.Print(ui.CreateBeautifulList(metricsData))
			ui.PrintSectionEnd()

			// thresholds section
			ui.PrintSection("Alert Thresholds")
			thresholdData := map[string]string{
				"CPU Threshold":    utils.FormatPercentage(cfg.CPUThreshold),
				"Memory Threshold": utils.FormatPercentage(cfg.MemThreshold),
				"Disk Threshold":   utils.FormatPercentage(cfg.DiskThreshold),
			}
			fmt.Print(ui.CreateBeautifulList(thresholdData))
			ui.PrintSectionEnd()

			// daemon status
			ui.PrintSection("Daemon Status")
			if process.IsRunning() {
				ui.PrintStatus("success", "Monitoring daemon is running")
			} else {
				ui.PrintStatus("warning", "Monitoring daemon is not running")
			}
			ui.PrintSectionEnd()
		},
	}
}
