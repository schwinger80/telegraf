package gpfs_io

import (
	"testing"
	"strconv"
)

func TestGPFSIO_Gather(t *testing.T) {
	// Erstellen Sie eine Instanz Ihres Plugins mit den gewünschten Konfigurationsoptionen
	g := &GPFSIO{
		PipePath: "/path/to/your/namedpipe",
	}

	// Erstellen Sie eine einfache Map, um die gesammelten Daten zu speichern
	collectedData := map[string]interface{}{}

	// Rufen Sie die Gather-Methode Ihres Plugins auf und übergeben die Map
	err := g.Gather(collectedData)
	if err != nil {
		t.Errorf("Fehler beim Ausführen von Gather: %v", err)
	}

	// Überprüfen Sie, ob die gesammelten Daten den erwarteten Werten entsprechen
	expectedMetrics := map[string]interface{}{
		// Hier können Sie erwartete Felder und Werte hinzufügen, die Sie erwarten.
		// Beispiel: "status_code": 0,
		// "timestamp": int64(1697982519),
		// "microseconds": int(24222),
		// "disk_count": int(84),
		// "bytes_read": uint64(22794336358085),
		// "bytes_written": uint64(0),
		// "open_calls": int(127818),
		// "close_calls": int(127817),
		// "read_calls": int(34356262),
		// "write_calls": int(0),
		// "readdir_calls": int(46906),
		// "inode_updates": int(76710),
		// Weitere Felder hier
	}

	// Vergleichen Sie die gesammelten Daten mit den erwarteten Werten
	for field, expectedValue := range expectedMetrics {
		actualValue, found := collectedData[field]
		if !found {
			t.Errorf("Feld %s nicht in den gesammelten Daten gefunden", field)
		}

		// Konvertieren Sie den erwarteten Wert in den richtigen Typ und vergleichen Sie ihn
		switch expectedValue.(type) {
		case int:
			actualInt, err := strconv.Atoi(actualValue.(string))
			if err != nil {
				t.Errorf("Feld %s sollte einen int-Wert haben, aber der Wert konnte nicht konvertiert werden", field)
			}
			if actualInt != expectedValue.(int) {
				t.Errorf("Feld %s hat einen unerwarteten Wert. Erwartet: %d, erhalten: %d", field, expectedValue.(int), actualInt)
			}
		case int64:
			actualInt64, err := strconv.ParseInt(actualValue.(string), 10, 64)
			if err != nil {
				t.Errorf("Feld %s sollte einen int64-Wert haben, aber der Wert konnte nicht konvertiert werden", field)
			}
			if actualInt64 != expectedValue.(int64) {
				t.Errorf("Feld %s hat einen unerwarteten Wert. Erwartet: %d, erhalten: %d", field, expectedValue.(int64), actualInt64)
			}
		case uint64:
			actualUint64, err := strconv.ParseUint(actualValue.(string), 10, 64)
			if err != nil {
				t.Errorf("Feld %s sollte einen uint64-Wert haben, aber der Wert konnte nicht konvertiert werden", field)
			}
			if actualUint64 != expectedValue.(uint64) {
				t.Errorf("Feld %s hat einen unerwarteten Wert. Erwartet: %d, erhalten: %d", field, expectedValue.(uint64), actualUint64)
			}
		case string:
			if actualValue.(string) != expectedValue.(string) {
				t.Errorf("Feld %s hat einen unerwarteten Wert. Erwartet: %s, erhalten: %s", field, expectedValue.(string), actualValue.(string))
			}
		default:
			t.Errorf("Ungültiger Datentyp für das Feld %s", field)
		}
	}
}
