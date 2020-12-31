// (c) 2020 Gon Y Yi. <https://gonyyi.com/copyright.txt>

package aface

type LoggerDummy1 struct {}
func(d *LoggerDummy1) Debugf(a string, b ...interface{}){}
func(d *LoggerDummy1) Infof(a string, b ...interface{}){}
func(d *LoggerDummy1) Warnf(a string, b ...interface{}){}
func(d *LoggerDummy1) Errorf(a string, b ...interface{}){}
func(d *LoggerDummy1) Fatalf(a string, b ...interface{}){}

type LoggerDummy1a struct {}
func(d *LoggerDummy1a) Tracef(a string, b ...interface{}){}
func(d *LoggerDummy1a) Debugf(a string, b ...interface{}){}
func(d *LoggerDummy1a) Infof(a string, b ...interface{}){}
func(d *LoggerDummy1a) Warnf(a string, b ...interface{}){}
func(d *LoggerDummy1a) Errorf(a string, b ...interface{}){}
func(d *LoggerDummy1a) Fatalf(a string, b ...interface{}){}

type LoggerDummy2 struct {}
func(d *LoggerDummy2) Debug(b ...interface{}){}
func(d *LoggerDummy2) Info(b ...interface{}){}
func(d *LoggerDummy2) Warn(b ...interface{}){}
func(d *LoggerDummy2) Error(b ...interface{}){}
func(d *LoggerDummy2) Fatal(b ...interface{}){}

type LoggerDummy3 struct {}
func(d *LoggerDummy3) Debug(a string){}
func(d *LoggerDummy3) Info(a string){}
func(d *LoggerDummy3) Warn(a string){}
func(d *LoggerDummy3) Error(a string){}
func(d *LoggerDummy3) Fatal(a string){}

type LoggerDummy12 struct {}
func(d *LoggerDummy12) Debugf(a string, b ...interface{}){}
func(d *LoggerDummy12) Infof(a string, b ...interface{}){}
func(d *LoggerDummy12) Warnf(a string, b ...interface{}){}
func(d *LoggerDummy12) Errorf(a string, b ...interface{}){}
func(d *LoggerDummy12) Fatalf(a string, b ...interface{}){}
func(d *LoggerDummy12) Debug(b ...interface{}){}
func(d *LoggerDummy12) Info(b ...interface{}){}
func(d *LoggerDummy12) Warn(b ...interface{}){}
func(d *LoggerDummy12) Error(b ...interface{}){}
func(d *LoggerDummy12) Fatal(b ...interface{}){}

type LoggerDummy13 struct {}
func(d *LoggerDummy13) Debugf(a string, b ...interface{}){}
func(d *LoggerDummy13) Infof(a string, b ...interface{}){}
func(d *LoggerDummy13) Warnf(a string, b ...interface{}){}
func(d *LoggerDummy13) Errorf(a string, b ...interface{}){}
func(d *LoggerDummy13) Fatalf(a string, b ...interface{}){}
func(d *LoggerDummy13) Debug(a string){}
func(d *LoggerDummy13) Info(a string){}
func(d *LoggerDummy13) Warn(a string){}
func(d *LoggerDummy13) Error(a string){}
func(d *LoggerDummy13) Fatal(a string){}

type LoggerDummy1a3 struct {}
func(d *LoggerDummy1a3) Tracef(a string, b ...interface{}){}
func(d *LoggerDummy1a3) Debugf(a string, b ...interface{}){}
func(d *LoggerDummy1a3) Infof(a string, b ...interface{}){}
func(d *LoggerDummy1a3) Warnf(a string, b ...interface{}){}
func(d *LoggerDummy1a3) Errorf(a string, b ...interface{}){}
func(d *LoggerDummy1a3) Fatalf(a string, b ...interface{}){}
func(d *LoggerDummy1a3) Trace(a string){}
func(d *LoggerDummy1a3) Debug(a string){}
func(d *LoggerDummy1a3) Info(a string){}
func(d *LoggerDummy1a3) Warn(a string){}
func(d *LoggerDummy1a3) Error(a string){}
func(d *LoggerDummy1a3) Fatal(a string){}
