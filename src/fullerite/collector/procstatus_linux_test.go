// +build linux

package collector

import (
	"fullerite/metric"

	"testing"
	"time"

	l "github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestProcStatusCollect(t *testing.T) {
	config := make(map[string]interface{})
	config["interval"] = 9999

	dims := map[string]string{
		"module": ".*",
	}

	config["generatedDimensions"] = dims

	channel := make(chan metric.Metric)

	testLog = l.WithFields(l.Fields{"testing": "procstatus_linux"})
	ps := NewProcStatus(channel, 12, testLog)
	ps.Configure(config)

	go ps.Collect()

	select {
	case <-ps.Channel():
		return
	case <-time.After(2 * time.Second):
		t.Fail()
	}
}

func TestProcStatusExtractDimensions(t *testing.T) {
	testLog = l.WithFields(l.Fields{"testing": "procstatus_linux"})

	config := make(map[string]interface{})

	dims := map[string]string{
		"module": "^python.*?test.*?\\.([^\\.]*)?\\-\\[\\d+\\]$",
		"order":  "^python.*?test.*?\\.[^\\.]*?\\-\\[(\\d+)\\]$",
	}
	config["generatedDimensions"] = dims

	ps := NewProcStatus(nil, 12, testLog)
	ps.Configure(config)

	dim := map[string]string{
		"module": "bond",
		"order":  "007",
	}

	extracted := ps.extractDimensions("python -m test.my.function.bond-[007]")
	assert.Equal(t, dim, extracted)
}

func TestProcStatusMetrics(t *testing.T) {
	testLog = l.WithFields(l.Fields{"testing": "procstatus_linux"})

	config := make(map[string]interface{})

	dims := map[string]string{
		"seven":  "(.......)",
		"eleven": "(...........)",
	}
	config["generatedDimensions"] = dims

	ps := NewProcStatus(nil, 12, testLog)
	ps.Configure(config)

	count := 0
	for _, m := range ps.procStatusMetrics() {
		mDims := m.Dimensions
		_, existsSeven := mDims["seven"]
		_, existsEleven := mDims["eleven"]
		if existsSeven == false || existsEleven == false {
			continue
		}
		count++
	}
	if count == 0 {
		t.Fail()
	}
}
