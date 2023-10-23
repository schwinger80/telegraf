package gpfs_io

import (
	"bufio"
	"log"
	"os"
	"strings"
	"strconv"
	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
)

type GPFSIO struct {
	PipePath string `toml:"pipe_path"`
	Log      telegraf.Logger
}

func (g *GPFSIO) Description() string {
	return "GPFS I/O MMPMon Input Plugin"
}

func (g *GPFSIO) SampleConfig() string {
	return `
      ## Konfigurationsparameter hier, falls benötigt
    `
}

func (g *GPFSIO) Init() error {
	// Hier können Sie eine Überprüfung durchführen, ob die named pipe existiert oder andere Initialisierungen.
	return nil
}

func (g *GPFSIO) Gather(acc telegraf.Accumulator) error {
	// Öffnen Sie die Named Pipe zum Lesen
	pipe, err := os.Open(g.PipePath)
	if err != nil {
		if g.Log != nil {
			g.Log.Errorf("Fehler beim Öffnen der Named Pipe: %v", err)
		} else {
			log.Fatalf("Fehler beim Öffnen der Named Pipe: %v", err)
		}
		return err
	}
	defer pipe.Close()

	// Erstellen Sie einen Scanner, um die Zeilen aus der Named Pipe zu lesen
	scanner := bufio.NewScanner(pipe)

	// Initialisieren der Variablen zum sammeln der MEtriken
	fields := make(map[string]interface{})
	tags := make(map[string]string)
	
	for scanner.Scan() {
		line := scanner.Text()

		// Hier MMPMON-Ausgabe analysieren
		words := strings.Fields(line)

		fields := make(map[string]interface{})
		tags := make(map[string]string)

		for i := 0; i < len(words); i += 2 {
			keyword := words[i]
			value := words[i+1]

			switch keyword {
			case "_n_":
				tags["ip_address"] = value
			case "_nn_":
				tags["hostname"] = value
			case "_rc_":
				rc, err := strconv.Atoi(value)
				if err == nil {
					fields["status_code"] = rc
				}
			case "_t_":
				timestamp, err := strconv.ParseInt(value, 10, 64)
				if err == nil {
					fields["timestamp"] = timestamp
				}
			case "_tu_":
				microseconds, err := strconv.Atoi(value)
				if err == nil {
					fields["microseconds"] = microseconds
				}
			case "_cl_":
				tags["cluster_name"] = value
			case "_fs_":
				tags["file_system"] = value
			case "_d_":
				disks, err := strconv.Atoi(value)
				if err == nil {
					fields["disk_count"] = disks
				}
			case "_br_":
				bytesRead, err := strconv.ParseUint(value, 10, 64)
				if err == nil {
					fields["bytes_read"] = bytesRead
				}
			case "_bw_":
				bytesWritten, err := strconv.ParseUint(value, 10, 64)
				if err == nil {
					fields["bytes_written"] = bytesWritten
				}
			case "_oc_":
				openCalls, err := strconv.Atoi(value)
				if err == nil {
					fields["open_calls"] = openCalls
				}
			case "_cc_":
				closeCalls, err := strconv.Atoi(value)
				if err == nil {
					fields["close_calls"] = closeCalls
				}
			case "_rdc_":
				readCalls, err := strconv.Atoi(value)
				if err == nil {
					fields["read_calls"] = readCalls
				}
			case "_wc_":
				writeCalls, err := strconv.Atoi(value)
				if err == nil {
					fields["write_calls"] = writeCalls
				}
			case "_dir_":
				readdirCalls, err := strconv.Atoi(value)
				if err == nil {
					fields["readdir_calls"] = readdirCalls
				}
			case "_iu_":
				inodeUpdates, err := strconv.Atoi(value)
				if err == nil {
					fields["inode_updates"] = inodeUpdates
				}
			}
		}

		// Daten in das Telegraf-Datenmodell umwandeln und an den Akkumulator senden
		acc.AddFields("gpfs_io_mmpmon", fields, tags)
	}

	if err := scanner.Err(); err != nil {
		if g.Log != nil {
			g.Log.Errorf("Fehler beim Lesen der Named Pipe: %v", err)
		} else {
			log.Fatalf("Fehler beim Lesen der Named Pipe: %v", err)
		}
		return err
	}

	return nil
}

func init() {
	inputs.Add("gpfs_io", func() telegraf.Input {
		return &GPFSIO{}
	})
}
