package collectors

import (
	"strconv"
	"strings"

	"github.com/StackExchange/scollector/opentsdb"
)

func init() {
	collectors = append(collectors, &IntervalCollector{F: c_dfstat_darwin})
}

func c_dfstat_darwin() opentsdb.MultiDataPoint {
	var md opentsdb.MultiDataPoint
	readCommand(func(line string) {
		fields := strings.Fields(line)
		if line == "" || len(fields) < 9 || !IsDigit(fields[2]) {
			return
		}
		mount := fields[8]
		if strings.HasPrefix(mount, "/Volumes/Time Machine Backups") {
			return
		}
		f5, _ := strconv.Atoi(fields[5])
		f6, _ := strconv.Atoi(fields[6])
		tags := opentsdb.TagSet{"mount": mount}
		Add(&md, "darwin.disk.fs.total", fields[1], tags)
		Add(&md, "darwin.disk.fs.used", fields[2], tags)
		Add(&md, "darwin.disk.fs.free", fields[3], tags)
		Add(&md, "darwin.disk.fs.inodes.total", f5+f6, tags)
		Add(&md, "darwin.disk.fs.inodes.used", fields[5], tags)
		Add(&md, "darwin.disk.fs.inodes.free", fields[6], tags)
	}, "df", "-lki")
	return md
}
