package gpfs

import (
        "bufio"
        "os/exec"
        "strings"
        "strconv"
        "github.com/influxdata/telegraf"
        "github.com/influxdata/telegraf/plugins/inputs"
)

type CommandRunner interface {
        RunCommand(command string) (*bufio.Scanner, error)
}

type RealCommandRunner struct{}

func (r *RealCommandRunner) RunCommand(command string) (*bufio.Scanner, error) {
        cmd := exec.Command("bash", "-c", command)
        stdout, err := cmd.StdoutPipe()
        if err != nil {
                return nil, err
        }
        if err := cmd.Start(); err != nil {
                return nil, err
        }
        return bufio.NewScanner(stdout), nil
}

type GPFSStats struct {
        Command string `toml:"command"`
        Runner  CommandRunner
}

func (g *GPFSStats) Description() string {
        return "Gather GPFS Performance Statistics using mmpmon"
}

func (g *GPFSStats) SampleConfig() string {
        return `command = "mmpmon -p -c 'fs_io_s'" `
}

func (g *GPFSStats) Gather(acc telegraf.Accumulator) error {
        if g.Runner == nil {
                g.Runner = &RealCommandRunner{}
        }
        scanner, err := g.Runner.RunCommand(g.Command)
        if err != nil {
                return err
        }

        var priorT, priorTu, priorBr, priorBw float64
        count := 0

        for scanner.Scan() {
                line := scanner.Text()
                fields := strings.Fields(line)
                if len(fields) < 22 {
                        continue // Ungültige Zeile, überspringen
                }

                t, _ := strconv.ParseFloat(fields[8], 64)
                tu, _ := strconv.ParseFloat(fields[10], 64)
                br, _ := strconv.ParseFloat(fields[18], 64)
                bw, _ := strconv.ParseFloat(fields[20], 64)

                if count > 0 {
                        deltaT := t - priorT
                        deltaTu := tu - priorTu
                        deltaBr := br - priorBr
                        deltaBw := bw - priorBw
                        dt := deltaT + (deltaTu / 1000000.0)
                        if dt > 0 {
                                rrate := (deltaBr / dt) / 1000000.0
                                wrate := (deltaBw / dt) / 1000000.0

                                // Füge die Messwerte zum Akkumulator hinzu
                                fields := map[string]interface{}{
                                        "read_rate":  rrate,
                                        "write_rate": wrate,
                                }
                                acc.AddFields("gpfs_io", fields, nil)
                        }
                }

                priorT = t
                priorTu = tu
                priorBr = br
                priorBw = bw
                count++
        }

        if err := scanner.Err(); err != nil {
                return err
        }

        return nil
}


func init() {
        inputs.Add("gpfs_stats", func() telegraf.Input {
                return &GPFSStats{
                        Command: "mmpmon -p -c 'fs_io_s'",
                }
        })
}

