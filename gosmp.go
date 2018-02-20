package gosmp

import (
	"fmt"
)

const (
	UNIT_SECOND 		= 1
	UNIT_MILLISECOND 	= 2
	UNIT_MICROSECOND 	= 3
	UNIT_PERCENTAGE		= 4
	UNIT_BYTES			= 5
	UNIT_KILOBYTES		= 6
	UNIT_MEGABYTES		= 7
	UNIT_GIGABYTES		= 8
	UNIT_TERABYTES		= 9
	UNIT_COUNTER		= 0

	EXIT_OK				= 0
	EXIT_WARNING		= 1
	EXIT_CRITICAL		= 2
	EXIT_UNKNOWN		= 3
)

var outputFormat string

type CheckResult struct {
	ExitStatus int
	CheckFunction func() (interface {}, string)
	Output string
	PerformanceData PerfdataValue
}

type PerfdataValue struct {
	label string
	counter bool
	unit string
	value CheckValue
	Crit CheckValue
	Warn CheckValue
	Min interface {}
	Max interface {}
}

type CheckValue struct {
	Value interface {}
}

func (checkResult *CheckResult) Run() {
	checkResult.PerformanceData.value.Value, checkResult.Output = checkResult.CheckFunction()

	checkResult.PerformanceData.value.convertToType()
	checkResult.PerformanceData.Warn.convertToType()
	checkResult.PerformanceData.Crit.convertToType()

	checkResult.validateCheckFunctionResult()
}

func (checkResult *CheckResult) Get() string {
	return checkResult.FormatOutput()
}

func (checkResult *CheckResult) uintCheck(warn uint64, crit uint64, checkValue uint64) int {
	if warn > crit {
		switch {
		case warn <= checkValue:
			return EXIT_WARNING
		case crit <= checkValue:
			return EXIT_CRITICAL
		default:
			return EXIT_OK
		}
	} else {
		switch {
		case crit >= checkValue:
			return EXIT_CRITICAL
		case warn >= checkValue:
			return EXIT_WARNING
		default:
			return EXIT_OK
		}
	}
}

func (checkResult *CheckResult) intCheck(warn int64, crit int64, checkValue int64) int {
	if warn > crit {
		switch {
		case warn <= checkValue:
			return EXIT_WARNING
		case crit <= checkValue:
			return EXIT_CRITICAL
		default:
			return EXIT_OK
		}
	} else {
		switch {
		case crit >= checkValue:
			return EXIT_CRITICAL
		case warn >= checkValue:
			return EXIT_WARNING
		default:
			return EXIT_OK
		}
	}
}

func (checkResult *CheckResult) floatCheck(warn float64, crit float64, checkValue float64) int {
	if warn > crit {
		switch {
		case checkValue < crit:
			return EXIT_CRITICAL
		case checkValue < warn:
			return EXIT_WARNING
		default:
			return EXIT_OK
		}
	} else {
		fmt.Println(warn, crit, checkValue)
		switch {
		case checkValue > crit:
			return EXIT_CRITICAL
		case checkValue > warn:
			return EXIT_WARNING
		default:
			return EXIT_OK
		}
	}
}

func (checkValue *CheckValue) convertToType() {
	switch checkValue.Value.(type) {
	case uint:
		checkValue.Value = uint64(checkValue.Value.(uint))
	case uint8:
		checkValue.Value = uint64(checkValue.Value.(uint8))
	case uint16:
		checkValue.Value = uint64(checkValue.Value.(uint16))
	case uint32:
		checkValue.Value = uint64(checkValue.Value.(uint32))
	case uint64:
		checkValue.Value = uint64(checkValue.Value.(uint64))

	case int:
		checkValue.Value = int64(checkValue.Value.(int))
	case int8:
		checkValue.Value = int64(checkValue.Value.(int8))
	case int16:
		checkValue.Value = int64(checkValue.Value.(int16))
	case int32:
		checkValue.Value = int64(checkValue.Value.(int32))
	case int64:
		checkValue.Value = int64(checkValue.Value.(int64))

	case float32:
		checkValue.Value = float64(checkValue.Value.(float32))
	case float64:
		checkValue.Value = float64(checkValue.Value.(float64))
	}
}

func (checkResult *CheckResult) validateCheckFunctionResult() {
	switch checkResult.PerformanceData.value.Value.(type) {
	case uint64:
		outputFormat = "%d"
		checkResult.ExitStatus = checkResult.uintCheck(checkResult.PerformanceData.Warn.Value.(uint64), checkResult.PerformanceData.Crit.Value.(uint64), checkResult.PerformanceData.value.Value.(uint64))
	case int64:
		outputFormat = "%d"
		checkResult.ExitStatus = checkResult.intCheck(checkResult.PerformanceData.Warn.Value.(int64), checkResult.PerformanceData.Crit.Value.(int64), checkResult.PerformanceData.value.Value.(int64))
	case float64:
		outputFormat = "%g"
		checkResult.ExitStatus = checkResult.floatCheck(checkResult.PerformanceData.Warn.Value.(float64), checkResult.PerformanceData.Crit.Value.(float64), checkResult.PerformanceData.value.Value.(float64))

	//case complex64:
	//case complex128:

	//case string:
	//	outputFormat = "%s"
	//	checkStatus = EXIT_OK
	default:
		panic("CheckFunction returned an unsupported data type")
	}
}


func (checkResult *CheckResult) FormatOutput() string {
	switch checkResult.ExitStatus {
	case EXIT_OK:
		checkResult.Output = "OK - " + checkResult.Output
	case EXIT_WARNING:
		checkResult.Output = "WARNING - " + checkResult.Output
	case EXIT_CRITICAL:
		checkResult.Output = "CRITICAL - " + checkResult.Output
	case EXIT_UNKNOWN:
		panic("Unknown state")
	}

	var output string;

	output = checkResult.Output + "|"

	output += checkResult.PerformanceData.label + "=" + fmt.Sprintf(outputFormat, checkResult.PerformanceData.value) + checkResult.PerformanceData.unit + ";;;;"

	return output
}