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
	inputs.Add("gpfs", func() telegraf.Input { return &Gpfs{} })
}

