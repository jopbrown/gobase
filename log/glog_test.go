package log_test

import (
	"bytes"
	"fmt"
	"os"

	"log/slog"

	"github.com/jopbrown/gobase/log"
)

func ExamplePrint() {
	log.SetGlobalLogger(log.NewLogger(os.Stdout, log.LevelInfo, log.LevelFatal))
	log.SetGlobalVerbose(1)

	log.Print("Print\n")
	log.Printf("%s\n", "Printf")
	log.Println("Println")
	log.Printlnf("%s", "Printlnf")

	v1 := log.V(1)
	v1.Print("v1: Print\n")
	v1.Printf("v1: %s\n", "Printf")
	v1.Println("v1: Println")
	v1.Printlnf("v1: %s", "Printlnf")

	v2 := log.V(2)
	v2.Print("v2: Print\n")
	v2.Printf("v2: %s\n", "Printf")
	v2.Println("v2: Println")
	v2.Printlnf("v2: %s", "Printlnf")

	// Output:
	// 	Print
	// Printf
	// Println
	// Printlnf
	// v1: Print
	// v1: Printf
	// v1: Println
	// v1: Printlnf
}

func ExampleLogByLevel() {
	log.SetGlobalLogger(log.NewLoggerWithFormat(os.Stdout, log.LevelInfo, log.LevelFatal, log.TestLoggerFormat()))
	log.SetGlobalVerbose(1)

	log.Debug("Debug\n")
	log.Debugf("%s", "Debugf")
	log.Info("Info\n")
	log.Infof("%s", "Infof")
	log.Warn("Warn\n")
	log.Warnf("%s", "Warnf")
	log.Error("Error\n")
	log.Errorf("%s", "Errorf")

	v1 := log.V(1).With("V1")
	v1.Debug("Debug\n")
	v1.Debugf("%s", "Debugf")
	v1.Info("Info\n")
	v1.Infof("%s", "Infof")
	v1.Warn("Warn\n")
	v1.Warnf("%s", "Warnf")
	v1.Error("Error\n")
	v1.Errorf("%s", "Errorf")

	v2 := log.V(2).With("V2")
	v2.Debug("Debug\n")
	v2.Debugf("%s", "Debugf")
	v2.Info("Info\n")
	v2.Infof("%s", "Infof")
	v2.Warn("Warn\n")
	v2.Warnf("%s", "Warnf")
	v2.Error("Error\n")
	v2.Errorf("%s", "Errorf")

	// Output:
	// INFO  V0 log_test.ExampleLogByLevel Info
	// INFO  V0 log_test.ExampleLogByLevel Infof
	// WARN  V0 log_test.ExampleLogByLevel Warn
	// WARN  V0 log_test.ExampleLogByLevel Warnf
	// ERROR V0 log_test.ExampleLogByLevel Error
	// ERROR V0 log_test.ExampleLogByLevel Errorf
	// INFO  V1 log_test.ExampleLogByLevel V1 Info
	// INFO  V1 log_test.ExampleLogByLevel V1 Infof
	// WARN  V1 log_test.ExampleLogByLevel V1 Warn
	// WARN  V1 log_test.ExampleLogByLevel V1 Warnf
	// ERROR V1 log_test.ExampleLogByLevel V1 Error
	// ERROR V1 log_test.ExampleLogByLevel V1 Errorf
}

func ExampleLogTee() {
	l1buf := bytes.NewBuffer(nil)
	l2buf := bytes.NewBuffer(nil)
	l1 := log.NewLoggerWithFormat(l1buf, log.LevelDebug, log.LevelInfo, log.TestLoggerFormat())
	l2 := log.NewLoggerWithFormat(l2buf, log.LevelWarn, log.LevelFatal, log.TestLoggerFormat())
	log.SetGlobalLogger(log.NewTeeLogger(l1, l2))
	log.SetGlobalVerbose(1)

	log.Print("Print\n")
	log.Debug("Debug\n")
	log.Debugf("%s", "Debugf")
	log.Info("Info\n")
	log.Infof("%s", "Infof")
	log.Warn("Warn\n")
	log.Warnf("%s", "Warnf")
	log.Error("Error\n")
	log.Errorf("%s", "Errorf")

	v1 := log.V(1).With("V1")
	v1.Print("Print\n")
	v1.Debug("Debug\n")
	v1.Debugf("%s", "Debugf")
	v1.Info("Info\n")
	v1.Infof("%s", "Infof")
	v1.Warn("Warn\n")
	v1.Warnf("%s", "Warnf")
	v1.Error("Error\n")
	v1.Errorf("%s", "Errorf")

	v2 := log.V(2).With("V2")
	v2.Print("Print\n")
	v2.Debug("Debug\n")
	v2.Debugf("%s", "Debugf")
	v2.Info("Info\n")
	v2.Infof("%s", "Infof")
	v2.Warn("Warn\n")
	v2.Warnf("%s", "Warnf")
	v2.Error("Error\n")
	v2.Errorf("%s", "Errorf")

	fmt.Print("logger1:\n", l1buf.String())
	fmt.Print("logger2:\n", l2buf.String())

	// Output:
	// logger1:
	// Print
	// DEBUG V0 log_test.ExampleLogTee Debug
	// DEBUG V0 log_test.ExampleLogTee Debugf
	// INFO  V0 log_test.ExampleLogTee Info
	// INFO  V0 log_test.ExampleLogTee Infof
	// Print
	// DEBUG V1 log_test.ExampleLogTee V1 Debug
	// DEBUG V1 log_test.ExampleLogTee V1 Debugf
	// INFO  V1 log_test.ExampleLogTee V1 Info
	// INFO  V1 log_test.ExampleLogTee V1 Infof
	// logger2:
	// Print
	// WARN  V0 log_test.ExampleLogTee Warn
	// WARN  V0 log_test.ExampleLogTee Warnf
	// ERROR V0 log_test.ExampleLogTee Error
	// ERROR V0 log_test.ExampleLogTee Errorf
	// Print
	// WARN  V1 log_test.ExampleLogTee V1 Warn
	// WARN  V1 log_test.ExampleLogTee V1 Warnf
	// ERROR V1 log_test.ExampleLogTee V1 Error
	// ERROR V1 log_test.ExampleLogTee V1 Errorf
}

func ExampleTeeGetWriter() {
	l1buf := bytes.NewBuffer(nil)
	l2buf := bytes.NewBuffer(nil)
	l1 := log.NewLoggerWithFormat(l1buf, log.LevelDebug, log.LevelInfo, log.TestLoggerFormat())
	l2 := log.NewLoggerWithFormat(l2buf, log.LevelWarn, log.LevelFatal, log.TestLoggerFormat())
	log.SetGlobalLogger(log.NewTeeLogger(l1, l2))
	log.SetGlobalVerbose(1)

	fmt.Fprintln(log.GetWriter(log.LevelNone), "None")
	fmt.Fprintln(log.GetWriter(log.LevelDebug), "Debug")
	fmt.Fprintln(log.GetWriter(log.LevelInfo), "Info")
	fmt.Fprintln(log.GetWriter(log.LevelWarn), "Warn")
	fmt.Fprintln(log.GetWriter(log.LevelError), "Error")
	fmt.Fprintln(log.GetWriter(log.LevelFatal), "Fatal")
	fmt.Fprintln(log.GetWriter(log.LevelAll), "all")

	fmt.Print("logger1:\n", l1buf.String())
	fmt.Print("logger2:\n", l2buf.String())
	// Output:
	// logger1:
	// Debug
	// Info
	// all
	// logger2:
	// Warn
	// Error
	// Fatal
	// all
}

func ExampleSlog() {
	l1buf := bytes.NewBuffer(nil)
	l2buf := bytes.NewBuffer(nil)
	l1 := log.NewLoggerWithFormat(l1buf, log.LevelDebug, log.LevelInfo, log.TestLoggerFormat())
	l2 := log.NewLoggerWithFormat(l2buf, log.LevelWarn, log.LevelFatal, log.TestLoggerFormat())
	log.SetGlobalLogger(log.NewTeeLogger(l1, l2))
	log.SetGlobalVerbose(1)

	s := log.S(true)
	s.Info("Debug", slog.String("attrString", "DebugString"))
	s.Info("Info", slog.String("attrString", "InfoString"))
	s.Warn("Warn", slog.String("attrString", "WarnString"))
	s.Error("Error", slog.String("attrString", "ErrorString"))

	s1 := log.V(1).With("V1").S(true)
	s1.Info("Debug", slog.String("attrString", "DebugString"))
	s1.Info("Info", slog.String("attrString", "InfoString"))
	s1.Warn("Warn", slog.String("attrString", "WarnString"))
	s1.Error("Error", slog.String("attrString", "ErrorString"))

	s2 := log.V(2).With("V2").S(true)
	s2.Info("Debug", slog.String("attrString", "DebugString"))
	s2.Info("Info", slog.String("attrString", "InfoString"))
	s2.Warn("Warn", slog.String("attrString", "WarnString"))
	s2.Error("Error", slog.String("attrString", "ErrorString"))

	fmt.Print("logger1:\n", l1buf.String())
	fmt.Print("logger2:\n", l2buf.String())

	// Output:
	// logger1:
	// {"level":"INFO","msg":"Debug","attrString":"DebugString"}
	// {"level":"INFO","msg":"Info","attrString":"InfoString"}
	// {"level":"INFO","msg":"Debug","verbose":1,"prefix":"V1","attrString":"DebugString"}
	// {"level":"INFO","msg":"Info","verbose":1,"prefix":"V1","attrString":"InfoString"}
	// logger2:
	// {"level":"WARN","msg":"Warn","attrString":"WarnString"}
	// {"level":"ERROR","msg":"Error","attrString":"ErrorString"}
	// {"level":"WARN","msg":"Warn","verbose":1,"prefix":"V1","attrString":"WarnString"}
	// {"level":"ERROR","msg":"Error","verbose":1,"prefix":"V1","attrString":"ErrorString"}
}

func ExampleSlogWithGroup() {
	l1buf := bytes.NewBuffer(nil)
	l2buf := bytes.NewBuffer(nil)
	l1 := log.NewLoggerWithFormat(l1buf, log.LevelDebug, log.LevelInfo, log.TestLoggerFormat())
	l2 := log.NewLoggerWithFormat(l2buf, log.LevelWarn, log.LevelFatal, log.TestLoggerFormat())
	log.SetGlobalLogger(log.NewTeeLogger(l1, l2))
	log.SetGlobalVerbose(1)

	s := log.S(true)
	s = s.WithGroup("GROUP").With(slog.String("with", "something"))
	s.Info("Debug", slog.String("attrString", "DebugString"))
	s.Info("Info", slog.String("attrString", "InfoString"))
	s.Warn("Warn", slog.String("attrString", "WarnString"))
	s.Error("Error", slog.String("attrString", "ErrorString"))

	fmt.Print("logger1:\n", l1buf.String())
	fmt.Print("logger2:\n", l2buf.String())

	// Output:
	// logger1:
	// {"level":"INFO","msg":"Debug","GROUP":{"with":"something","attrString":"DebugString"}}
	// {"level":"INFO","msg":"Info","GROUP":{"with":"something","attrString":"InfoString"}}
	// logger2:
	// {"level":"WARN","msg":"Warn","GROUP":{"with":"something","attrString":"WarnString"}}
	// {"level":"ERROR","msg":"Error","GROUP":{"with":"something","attrString":"ErrorString"}}
}
