package gpfs_io

import (
	"testing"
	"strconv"
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/testutil"
)

func TestGPFSIO_Gather(t *testing.T) {
	// Erstellen Sie eine Instanz Ihres Plugins mit den gewünschten Konfigurationsoptionen
	g := &GPFSIO{
		PipePath: "/path/to/your/namedpipe",
	}

	// Erstellen Sie eine leere Instanz von telegraf.Accumulator
	acc := &testutil.Accumulator{}

	// Rufen Sie die Gather-Methode Ihres Plugins auf und übergeben die Instanz von telegraf.Accumulator
	err := g.Gather(acc)
	if err != nil {
		t.Errorf("Fehler beim Ausführen von Gather: %v", err)
	}

	// Überprüfen Sie, ob die gesammelten Daten den erwarteten Werten entsprechen
	expectedMetrics := []testutil.Metric{
		testutil.MustMetric(
			"gpfs_io_mmpmon",
			map[string]string{
				"_n_":  "10.156.153.84",
				"_nn_": "hpdar03c04s08",
				"_cl_": "LRZ_DSS03.dss.lrz.de",
				"_fs_": "dsstbyfs01",
			},
			map[string]interface{}{
				"_rc_":           0,
				"_t_":            int64(1697982519),
				"_tu_":           24222,
				"_d_":            84,
				"_br_":           uint64(22794336358085),
				"_bw_":           uint64(0),
				"_oc_":           127818,
				"_cc_":           127817,
				"_rdc_":          34356262,
				"_wc_":           0,
				"_dir_":          46906,
				"_iu_":           76710,
			},
		),
		// Weitere erwartete Metriken hier hinzufügen, falls erforderlich
	}

	// Vergleichen Sie die gesammelten Daten mit den erwarteten Metriken
	testutil.RequireMetricsEqual(t, expectedMetrics, acc.GetTelegrafMetrics())
}
