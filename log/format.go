package log

type LoggerFormat struct {
	AddLevel       bool
	AddVerbose     bool
	AddPrefix      bool
	AddCaller      bool
	AddSource      bool
	AddDateTime    bool
	LongCaller     bool
	SourceDepth    int
	DateTimeFormat LoggerFormatDateTime
}

type LoggerFormatDateTime struct {
	AddDate         bool
	AddTime         bool
	AddMicroseconds bool
}

func DefaultLoggerFormat() LoggerFormat {
	return LoggerFormat{
		AddLevel:    true,
		AddVerbose:  false,
		AddPrefix:   true,
		AddCaller:   false,
		AddSource:   true,
		AddDateTime: true,
		LongCaller:  false,
		SourceDepth: 2,
		DateTimeFormat: LoggerFormatDateTime{
			AddDate:         false,
			AddTime:         true,
			AddMicroseconds: false,
		},
	}
}

func SimpleLoggerFormat() LoggerFormat {
	return LoggerFormat{
		AddLevel:    false,
		AddVerbose:  false,
		AddPrefix:   true,
		AddCaller:   false,
		AddSource:   false,
		AddDateTime: false,
		LongCaller:  false,
		SourceDepth: 0,
		DateTimeFormat: LoggerFormatDateTime{
			AddDate:         false,
			AddTime:         false,
			AddMicroseconds: false,
		},
	}
}

func FileLoggerFormat() LoggerFormat {
	return LoggerFormat{
		AddLevel:    true,
		AddVerbose:  true,
		AddPrefix:   true,
		AddCaller:   true,
		AddSource:   true,
		AddDateTime: true,
		LongCaller:  false,
		SourceDepth: 3,
		DateTimeFormat: LoggerFormatDateTime{
			AddDate:         true,
			AddTime:         true,
			AddMicroseconds: true,
		},
	}
}

func TestLoggerFormat() LoggerFormat {
	return LoggerFormat{
		AddLevel:    true,
		AddVerbose:  true,
		AddPrefix:   true,
		AddCaller:   true,
		AddSource:   false,
		AddDateTime: false,
		LongCaller:  false,
		SourceDepth: 0,
		DateTimeFormat: LoggerFormatDateTime{
			AddDate:         false,
			AddTime:         false,
			AddMicroseconds: false,
		},
	}
}

func FullLoggerFormat() LoggerFormat {
	return LoggerFormat{
		AddLevel:    true,
		AddVerbose:  true,
		AddPrefix:   true,
		AddCaller:   true,
		AddSource:   true,
		AddDateTime: true,
		LongCaller:  true,
		SourceDepth: 0,
		DateTimeFormat: LoggerFormatDateTime{
			AddDate:         true,
			AddTime:         true,
			AddMicroseconds: true,
		},
	}
}
