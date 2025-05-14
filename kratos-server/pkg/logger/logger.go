package log

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"time"
)

type SqlLogger struct {
	SlowThreshold time.Duration
	log           *log.Helper
	traceErrStr   string
	traceWarnStr  string
	traceStr      string
}

func NewSqlLogger(log *log.Helper, slow time.Duration) *SqlLogger {
	return &SqlLogger{
		SlowThreshold: slow,
		log:           log,
		traceWarnStr:  "%s %s [%.3fms] [rows:%v] %s",
	}
}

func (s *SqlLogger) LogMode(level logger.LogLevel) logger.Interface {
	return NewSqlLogger(s.log, s.SlowThreshold)
}

func (s *SqlLogger) Info(ctx context.Context, s2 string, i ...interface{}) {
	s.log.Infof(s2, i...)
}

func (s *SqlLogger) Warn(ctx context.Context, s2 string, i ...interface{}) {
	s.log.Warnf(s2, i...)
}

func (s *SqlLogger) Error(ctx context.Context, s2 string, i ...interface{}) {
	s.log.Errorf(s2, i...)
}

func (s *SqlLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {

	elapsed := time.Since(begin)

	if elapsed > s.SlowThreshold && s.SlowThreshold != 0 {
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", s.SlowThreshold)
		if rows == -1 {
			s.Warn(ctx, s.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			s.Warn(ctx, s.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}
