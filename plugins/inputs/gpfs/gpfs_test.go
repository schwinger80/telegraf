package gpfs_io

import (
	"testing"
	"github.com/influxdata/telegraf/testutil"
)

func TestGPFSIO_Gather(t *testing.T) {
	// Erstellen Sie eine Instanz Ihres Plugins mit den gewünschten Konfigurationsoptionen
	g := &GPFSIO{
		PipePath: "/path/to/your/namedpipe",
	}

	// Erstellen Sie einen Akkumulator für Testzwecke
	acc := testutil.Accumulator{}

	// Rufen Sie die Gather-Methode Ihres Plugins auf
	err := g.Gather(&acc)
	if err != nil {
		t.Errorf("Fehler beim Ausführen von Gather: %v", err)
	}

	// Erstellen Sie eine Slice von telegraf.Metric für erwartete Metriken
	expectedMetrics := []telegraf.Metric{
		// Hier können Sie erwartete Metriken hinzufügen, die Sie erwarten.
		// Beispiel: testutil.MustMetric(
		//     "gpfs_io_mmpmon",
		//     map[string]string{
		//         "ip_address":   "10.156.153.84",
		//         "hostname":     "hpdar03c04s08",
		//         "cluster_name": "LRZ_DSS03.dss.lrz.de",
		//         // Weitere Tags hier
		//     },
		//     map[string]interface{}{
		//         "status_code":    0,
		//         "timestamp":      int64(1697982519),
		//         "microseconds":   int(24222),
		//         "disk_count":     int(84),
		//         "bytes_read":     uint64(22794336358085),
		//         "bytes_written":  uint64(0),
		//         "open_calls":     int(127818),
		//         "close_calls":    int(127817),
		//         "read_calls":     int(34356262),
		//         "write_calls":    int(0),
		//         "readdir_calls":  int(46906),
		//         "inode_updates":  int(76710),
		//         // Weitere Felder hier
		//     },
		//     int64(1697982519),
		// )
	}

	// Überprüfen Sie, ob die tatsächlichen Metriken den erwarteten Metriken entsprechen
	testutil.RequireMetricsEqual(t, expectedMetrics, acc.GetTelegrafMetrics())
}
