//go:build !custom || inputs || inputs.gpfs

package all

import _ "github.com/influxdata/telegraf/plugins/inputs/gpfs" // register plugin
