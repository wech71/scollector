package collectors

import (
	"log"
	"strconv"
	"strings"

	"github.com/StackExchange/scollector/opentsdb"
)

func init() {
	collectors = append(collectors, &IntervalCollector{F: c_dfstat_blocks_linux})
	collectors = append(collectors, &IntervalCollector{F: c_dfstat_inodes_linux})
}

func c_dfstat_blocks_linux() opentsdb.MultiDataPoint {
	var md opentsdb.MultiDataPoint
	readCommand(func(line string) {
		fields := strings.Fields(line)
		if line == "" || len(fields) < 6 || !IsDigit(fields[2]) {
			return
		}
		mount := fields[5]
		tags := opentsdb.TagSet{"mount": mount}
		os_tags := opentsdb.TagSet{"disk": mount}
		//Meta Data will need to indicate that these are 1kblocks
		Add(&md, "linux.disk.fs.space_total", fields[1], tags)
		Add(&md, "linux.disk.fs.space_used", fields[2], tags)
		Add(&md, "linux.disk.fs.space_free", fields[3], tags)
		Add(&md, osDiskTotal, fields[1], os_tags)
		Add(&md, osDiskUsed, fields[2], os_tags)
		Add(&md, osDiskFree, fields[3], os_tags)
		st, err := strconv.ParseFloat(fields[1], 64)
		sf, err := strconv.ParseFloat(fields[3], 64)
		if err == nil {
			if st != 0 {
				Add(&md, osDiskPctFree, sf/st*100, os_tags)
			}
		} else {
			log.Println(err)
		}
	}, "df", "-lP", "--block-size", "1")
	return md
}

func c_dfstat_inodes_linux() opentsdb.MultiDataPoint {
	var md opentsdb.MultiDataPoint
	readCommand(func(line string) {
		fields := strings.Fields(line)
		if line == "" || len(fields) < 6 || !IsDigit(fields[2]) {
			return
		}
		mount := fields[5]
		tags := opentsdb.TagSet{"mount": mount}
		Add(&md, "linux.disk.fs.inodes_total", fields[1], tags)
		Add(&md, "linux.disk.fs.inodes_used", fields[2], tags)
		Add(&md, "linux.disk.fs.inodes_free", fields[3], tags)
	}, "df", "-liP")
	return md
}
