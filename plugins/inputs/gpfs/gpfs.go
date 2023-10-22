// go:generate ../../../tools/readme_config_includer/generator
package gpfs

import (
	"bufio"
	"os"
	"strings"

	"github.com/influxdata/telegraf"
	"github.com/influxdata/telegraf/plugins/inputs"
)

//go:embed sample.conf
var sampleConfig string

type GPFS struct {
	PipePath string         `toml:"pipe_path"`
	Log      telegraf.Logger `toml:"-"`
}

func (*GPFS) SampleConfig() string {
	return sampleConfig
}

func (g *GPFS) Init() error {
	// Hier können Sie eine Überprüfung durchführen, ob die named pipe existiert oder andere Initialisierungen.
	return nil
}
func (g *GPFS) Gather(acc telegraf.Accumulator) error {
	// Name der Named Pipe (FIFO) anpassen
	pipeName := "/path/to/your/namedpipe"

	// Öffnen Sie die Named Pipe zum Lesen
	pipe, err := os.Open(pipeName)
	if err != nil {
		log.Fatalf("Fehler beim Öffnen der Named Pipe: %v", err)
		return err
	}
	defer pipe.Close()

	// Erstellen Sie einen Scanner, um die Zeilen aus der Named Pipe zu lesen
	scanner := bufio.NewScanner(pipe)

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
		acc.AddFields("gpfs_mmpmon", fields, tags)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Fehler beim Lesen der Named Pipe: %v", err)
		return err
	}

	return nil
}


func (g *GPFS) Gather(acc telegraf.Accumulator) error {
	pipe, err := os.Open(g.PipePath)
	if err != nil {
		return err
	}
	defer pipe.Close()

	reader := bufio.NewReader(pipe)
	line, _, err := reader.ReadLine()
	if err != nil {
		return err
	}

	// Hier ist eine einfache Annahme: Der String ist "key=value"
	parts := strings.Split(string(line), "=")
	if len(parts) != 2 {
		return fmt.Errorf("Unexpected format: %s", line)
	}

	// Speichert den Wert im Telegraf-Akkumulator
	acc.AddFields("gpfs", map[string]interface{}{parts[0]: parts[1]}, nil)

	return nil
}

func init() {
	inputs.Add("gpfs", func() telegraf.Input { return &gpfs{} })
}

