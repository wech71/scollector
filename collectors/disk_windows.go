package collectors

import (
	"github.com/StackExchange/scollector/opentsdb"
	"github.com/StackExchange/slog"
	"github.com/StackExchange/wmi"
)

func init() {
	collectors = append(collectors, &IntervalCollector{F: c_physical_disk_windows})
	collectors = append(collectors, &IntervalCollector{F: c_diskspace_windows})
}

func c_diskspace_windows() opentsdb.MultiDataPoint {
	var dst []Win32_PerfFormattedData_PerfDisk_LogicalDisk
	var q = wmi.CreateQuery(&dst, `WHERE Name <> '_Total'`)
	err := queryWmi(q, &dst)
	if err != nil {
		slog.Infoln("diskspace:", err)
		return nil
	}
	var md opentsdb.MultiDataPoint
	for _, v := range dst {
		Add(&md, "win.disk.fs.space_free", v.FreeMegabytes*1048576, opentsdb.TagSet{"partition": v.Name})
		Add(&md, osDiskFree, v.FreeMegabytes*1048576, opentsdb.TagSet{"disk": v.Name})
		if v.PercentFreeSpace != 0 {
			space_total := v.FreeMegabytes * 1048576 * 100 / v.PercentFreeSpace
			space_used := space_total - v.FreeMegabytes*1048576
			Add(&md, "win.disk.fs.space_total", space_total, opentsdb.TagSet{"partition": v.Name})
			Add(&md, "win.disk.fs.space_used", space_used, opentsdb.TagSet{"partition": v.Name})
			Add(&md, osDiskTotal, space_total, opentsdb.TagSet{"disk": v.Name})
			Add(&md, osDiskUsed, space_used, opentsdb.TagSet{"disk": v.Name})
		}

		Add(&md, "win.disk.fs.percent_free", v.PercentFreeSpace, opentsdb.TagSet{"partition": v.Name})
		Add(&md, osDiskPctFree, v.PercentFreeSpace, opentsdb.TagSet{"disk": v.Name})
	}
	return md
}

type Win32_PerfFormattedData_PerfDisk_LogicalDisk struct {
	FreeMegabytes    uint64
	Name             string
	PercentFreeSpace uint64
}

func c_physical_disk_windows() opentsdb.MultiDataPoint {
	var dst []Win32_PerfRawData_PerfDisk_PhysicalDisk
	var q = wmi.CreateQuery(&dst, `WHERE Name <> '_Total'`)
	err := queryWmi(q, &dst)
	if err != nil {
		slog.Infoln("disk_physical:", err)
		return nil
	}
	var md opentsdb.MultiDataPoint
	for _, v := range dst {
		Add(&md, "win.disk.duration", v.AvgDiskSecPerRead, opentsdb.TagSet{"disk": v.Name, "type": "read"})
		Add(&md, "win.disk.duration", v.AvgDiskSecPerWrite, opentsdb.TagSet{"disk": v.Name, "type": "write"})
		Add(&md, "win.disk.queue", v.AvgDiskReadQueueLength, opentsdb.TagSet{"disk": v.Name, "type": "read"})
		Add(&md, "win.disk.queue", v.AvgDiskWriteQueueLength, opentsdb.TagSet{"disk": v.Name, "type": "write"})
		Add(&md, "win.disk.ops", v.DiskReadsPerSec, opentsdb.TagSet{"disk": v.Name, "type": "read"})
		Add(&md, "win.disk.ops", v.DiskWritesPerSec, opentsdb.TagSet{"disk": v.Name, "type": "write"})
		Add(&md, "win.disk.bytes", v.DiskReadBytesPerSec, opentsdb.TagSet{"disk": v.Name, "type": "read"})
		Add(&md, "win.disk.bytes", v.DiskWriteBytesPerSec, opentsdb.TagSet{"disk": v.Name, "type": "write"})
		Add(&md, "win.disk.percent_time", v.PercentDiskReadTime, opentsdb.TagSet{"disk": v.Name, "type": "read"})
		Add(&md, "win.disk.percent_time", v.PercentDiskWriteTime, opentsdb.TagSet{"disk": v.Name, "type": "write"})
		Add(&md, "win.disk.spltio", v.SplitIOPerSec, opentsdb.TagSet{"disk": v.Name})
	}
	return md
}

type Win32_PerfRawData_PerfDisk_PhysicalDisk struct {
	AvgDiskReadQueueLength  uint64
	AvgDiskSecPerRead       uint32
	AvgDiskSecPerWrite      uint32
	AvgDiskWriteQueueLength uint64
	DiskReadBytesPerSec     uint64
	DiskReadsPerSec         uint32
	DiskWriteBytesPerSec    uint64
	DiskWritesPerSec        uint32
	Name                    string
	PercentDiskReadTime     uint64
	PercentDiskWriteTime    uint64
	SplitIOPerSec           uint32
}
