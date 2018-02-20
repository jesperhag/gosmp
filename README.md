# gosmp
Library for creating simple monitoring plugins in Go

The aim with this project is to compile with the guidelines defined here: https://www.monitoring-plugins.org/doc/guidelines.html as well as streamline the process making a monitoring plugin to the point where you only need to define the check itself and not have to worry about handling CLI options, arguments etc as well as worrying about returning valid performance data.

WIP

Example usage, plugin checks diskspace used (in bytes):

```go
package main

import (
	"github.com/jesperhag/gosmp"
	"fmt"
	"os"
	"syscall"
)

func main() {
	checkResult := gosmp.CheckResult{
		PerformanceData: gosmp.PerfdataValue{},
	}

	checkResult.PerformanceData.Warn.Value = uint(180*1073741824)
	checkResult.PerformanceData.Crit.Value = uint(200*1073741824)

	checkResult.CheckFunction = checkDisk
	checkResult.PerformanceData.Label = "/"
	checkResult.PerformanceData.Unit = gosmp.UNIT_BYTES

	checkResult.Run()
	fmt.Println(checkResult.FormatOutput())

	os.Exit(checkResult.ExitStatus)
}

func checkDisk() (interface {}, string) {
	var stat syscall.Statfs_t
	wd, err := os.Getwd() ; if err != nil {
		fmt.Println(err)
	}

	wd = "/"

	syscall.Statfs(wd, &stat)

	available := stat.Bavail * uint64(stat.Bsize)
	total := stat.Blocks * uint64(stat.Bsize)
	used := total - available

	return used, wd
}
```
