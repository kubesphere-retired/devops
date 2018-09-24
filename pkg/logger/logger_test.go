/*
Copyright 2018 The KubeSphere Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package logger

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func readBuf(buf *bytes.Buffer) string {
	str := buf.String()
	buf.Reset()
	return str
}

func TestLogger(t *testing.T) {
	buf := new(bytes.Buffer)
	SetOutput(buf)

	Debug("debug log, should ignore by default")
	assert.Empty(t, readBuf(buf))

	Info("info log, should visable")
	assert.Contains(t, readBuf(buf), "info log, should visable")

	Info("format [%d]", 111)
	log := readBuf(buf)
	assert.Contains(t, log, "format [111]")
	t.Log(log)

	SetLevelByString("debug")
	Debug("debug log, now it becomes visible")
	assert.Contains(t, readBuf(buf), "debug log, now it becomes visible")

	logger = NewLogger()
	logger.SetPrefix("(prefix)").SetSuffix("(suffix)").SetOutput(buf)

	logger.Warn("log_content")
	log = readBuf(buf)
	assert.Regexp(t, " -WARNING- \\(prefix\\)log_content \\(logger_test.go:\\d+\\)\\(suffix\\)", log)
	t.Log(log)

	logger.HideCallstack()
	logger.Warn("log_content")
	log = readBuf(buf)
	assert.Regexp(t, " -WARNING- \\(prefix\\)log_content\\(suffix\\)", log)
	t.Log(log)
}
