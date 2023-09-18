package gpfs

import (
        "errors"
        "testing"
        "bufio"
        "strings"
        "github.com/influxdata/telegraf/testutil"
        "github.com/stretchr/testify/require"
)

type MockCommandRunner struct {
        Command string
        Output  string
        Err     error
}

func (m *MockCommandRunner) RunCommand(command string) (*bufio.Scanner, error) {
   if command != m.Command {
      return nil, errors.New("unexpected command")
   }
   return bufio.NewScanner(strings.NewReader(m.Output)), m.Err
}


func TestGPFSStats(t *testing.T) {
        tests := []struct {
                name    string
                command string
                output  string
                err     error
                metrics map[string]interface{}
        }{
                // Beispieltestfälle hier hinzufügen
                {
                        name:    "valid case",
                        command: "mmpmon -p -c 'fs_io_s'",
                        output:  "92.1 MB/sec read      0.1 MB/sec write",
                        err:     nil,
                        metrics: map[string]interface{}{
                                "read_rate":  92.1,
                                "write_rate": 0.1,
                        },
                },
        }

        for _, tt := range tests {
                t.Run(tt.name, func(t *testing.T) {
                        mcr := &MockCommandRunner{
                                Command: tt.command,
                                Output:  tt.output,
                                Err:     tt.err,
                        }
                        gpfsStats := &GPFSStats{
                                Command: tt.command,
                                Runner:  mcr,
                        }

                        var acc testutil.Accumulator

                        err := gpfsStats.Gather(&acc)
                        require.NoError(t, err)

                        for k, v := range tt.metrics {
                                require.True(t, acc.HasMeasurement(k),
                                        "missing measurement: %q: %v", k, v)

                        }
                })
        }

}
