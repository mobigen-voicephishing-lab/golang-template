package logger

import (
	"bytes"
	"fmt"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// Formatter - logrus formatter, implements logrus.Formatter
type Formatter struct {
	// TimestampFormat - default: time.StampMilli = "Jan _2 15:04:05.000"
	TimestampFormat string

	// ShowFullLevel - show a full level [WARNING] instead of [WARN]
	ShowFullLevel bool
	// NoUppercaseLevel - no upper case for level value
	NoUppercaseLevel bool

	// ShowFields
	ShowFields bool
	// SortField sort enable/disable
	SortFields bool
}

// Format an log entry
func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = time.StampMilli
	}

	// output buffer
	b := &bytes.Buffer{}

	// write time
	b.WriteString(entry.Time.Format(timestampFormat))

	// write level
	var level string
	if f.NoUppercaseLevel {
		level = entry.Level.String()
	} else {
		level = strings.ToUpper(entry.Level.String())
	}
	b.WriteString(" [")
	if f.ShowFullLevel {
		b.WriteString(level)
	} else {
		b.WriteString(level[:4])
	}
	b.WriteString("] ")

	// write caller
	f.writeCaller(b, entry)

	// write fields
	if f.ShowFields {
		f.writeFields(b, entry)
	}

	// write message
	b.WriteString(entry.Message)

	b.WriteByte('\n')
	return b.Bytes(), nil
}

func (f *Formatter) writeCaller(b *bytes.Buffer, entry *logrus.Entry) {
	if entry.HasCaller() {
		_, filename := path.Split(entry.Caller.File)
		if len(filename) < 16 {
			fmt.Fprintf(b, "[%-16s : %3d] ",
				filename, entry.Caller.Line)
		} else {
			fmt.Fprintf(b, "[%.16s : %3d] ",
				filename, entry.Caller.Line)
		}
	}
}

func (f *Formatter) writeFields(b *bytes.Buffer, entry *logrus.Entry) {
	if len(entry.Data) != 0 {
		fields := make([]string, 0, len(entry.Data))
		for field := range entry.Data {
			fields = append(fields, field)
		}
		if f.SortFields {
			sort.Strings(fields)
		}
		fmt.Fprintf(b, "[ ")
		for idx, field := range fields {
			f.writeField(b, entry, field)
			if idx+1 < len(fields) {
				fmt.Fprintf(b, ", ")
			}
		}
		fmt.Fprintf(b, " ] ")
	}
}

func (f *Formatter) writeField(b *bytes.Buffer, entry *logrus.Entry, field string) {
	fmt.Fprintf(b, "%s:%v", field, entry.Data[field])
}
