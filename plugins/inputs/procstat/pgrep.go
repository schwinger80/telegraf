package procstat

import (
        "fmt"
        "os"
        "os/exec"
        "strconv"
        "strings"

        "github.com/influxdata/telegraf/internal"
)

// Implementation of PIDGatherer that execs pgrep to find processes
type Pgrep struct {
        path string
}

func NewPgrep() (PIDFinder, error) {
        path, err := exec.LookPath("pgrep")
        if err != nil {
                return nil, fmt.Errorf("could not find pgrep binary: %w", err)
        }
        return &Pgrep{path}, nil
}

func (pg *Pgrep) PidFile(path string) ([]PID, error) {
        var pids []PID
        pidString, err := os.ReadFile(path)
        if err != nil {
                return pids, fmt.Errorf("failed to read pidfile %q: %w",
                        path, err)
        }
        pid, err := strconv.ParseInt(strings.TrimSpace(string(pidString)), 10, 32)
        if err != nil {
                return pids, err
        }
        pids = append(pids, PID(pid))

   childPIDs, err := pg.getChildPIDs(pids)
   if err != nil {
      return nil, err
   }

   allPIDs := append(pids, childPIDs...)
   return allPIDs, nil
}

func (pg *Pgrep) Pattern(pattern string) ([]PID, error) {
        args := []string{pattern}
        return find(pg.path, args)
}

func (pg *Pgrep) UID(user string) ([]PID, error) {
        args := []string{"-u", user}
        return find(pg.path, args)
}

//func (pg *Pgrep) FullPattern(pattern string) ([]PID, error) {
//      args := []string{"-f", pattern}
//      return find(pg.path, args)
//}

func (pg *Pgrep) FullPattern(pattern string) ([]PID, error) {
        args := []string{"-f", pattern}
        matchingPIDs, err := find(pg.path, args)
        if err != nil {
                return nil, err
        }

        childPIDs, err := pg.getChildPIDs(matchingPIDs)
        if err != nil {
                return nil, err
        }

        allPIDs := append(matchingPIDs, childPIDs...)
        return allPIDs, nil
}

//func (pg *Pgrep) getChildPIDs(parentPIDs []PID) ([]PID, error) {
//      var childPIDs []PID
//      for _, parentPID := range parentPIDs {
//              args := []string{"-P", fmt.Sprint(parentPID)}
//              childPIDList, err := find(pg.path, args)
//              if err != nil {
//                      return nil, err
//              }
//              childPIDs = append(childPIDs, childPIDList...)
//      }
//      return childPIDs, nil
//}


func (pg *Pgrep) getChildPIDs(parentPIDs []PID) ([]PID, error) {
        var childPIDs []PID
        for _, parentPID := range parentPIDs {
                args := []string{"-P", fmt.Sprint(parentPID)}
                childPIDList, err := find(pg.path, args)
                if err != nil {
                        return nil, err
                }

                // Rekursion: Finde die Child-Prozesse der aktuellen Child-Prozesse
                grandChildPIDList, err := pg.getChildPIDs(childPIDList)
                if err != nil {
                        return nil, err
                }

                childPIDs = append(childPIDs, childPIDList...)
                childPIDs = append(childPIDs, grandChildPIDList...)
        }
        return childPIDs, nil
}

func find(path string, args []string) ([]PID, error) {
        out, err := run(path, args)
        if err != nil {
                return nil, err
        }

        return parseOutput(out)
}

func run(path string, args []string) (string, error) {
        out, err := exec.Command(path, args...).Output()

        //if exit code 1, ie no processes found, do not return error
        if i, _ := internal.ExitStatus(err); i == 1 {
                return "", nil
        }

        if err != nil {
                return "", fmt.Errorf("error running %q: %w", path, err)
        }
        return string(out), err
}

func parseOutput(out string) ([]PID, error) {
        pids := []PID{}
        fields := strings.Fields(out)
        for _, field := range fields {
                pid, err := strconv.ParseInt(field, 10, 32)
                if err != nil {
                        return nil, err
                }
                pids = append(pids, PID(pid))
        }
        return pids, nil
}
