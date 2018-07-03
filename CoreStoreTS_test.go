package libDatabox

import (
	"strconv"
	s "strings"
	"testing"
	"time"
)

func TestAggregationFunctionMinOnEmptyDS(t *testing.T) {

	_, err := StoreClient.TSJSON.LastN(dsID+"TestAggregationFunctionOnEmptyDS", 4, TimeSeriesQueryOptions{
		AggregationFunction: Min,
	})
	if err != nil {
		t.Errorf("Call to LastN failed with error %s", err.Error())
	}

}

func TestAggregationFunctionMaxOnEmptyDS(t *testing.T) {

	_, err := StoreClient.TSJSON.LastN(dsID+"TestAggregationFunctionOnEmptyDS", 4, TimeSeriesQueryOptions{
		AggregationFunction: Max,
	})
	if err != nil {
		t.Errorf("Call to LastN failed with error %s", err.Error())
	}

}

func TestAggregationFunctionSumOnEmptyDS(t *testing.T) {

	_, err := StoreClient.TSJSON.LastN(dsID+"TestAggregationFunctionOnEmptyDS", 4, TimeSeriesQueryOptions{
		AggregationFunction: Sum,
	})
	if err != nil {
		t.Errorf("Call to LastN failed with error %s", err.Error())
	}

}

func TestAggregationFunctionSDOnEmptyDS(t *testing.T) {

	_, err := StoreClient.TSJSON.LastN(dsID+"TestAggregationFunctionOnEmptyDS", 4, TimeSeriesQueryOptions{
		AggregationFunction: StandardDeviation,
	})
	if err != nil {
		t.Errorf("Call to LastN failed with error %s", err.Error())
	}

}

func TestWrite(t *testing.T) {
	err := StoreClient.TSJSON.Write(dsID, []byte("{\"value\":3.1415}"))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}
}

func benchmarkWrite(num int, b *testing.B) {
	// write b.N times
	for n := 0; n < b.N; n++ {
		for i := 0; i < num; i++ {
			StoreClient.TSJSON.Write(dsID+"benchmarkWrite"+strconv.Itoa(num), []byte("{\"value\":"+strconv.Itoa(n)+"}"))
		}
	}
}

func BenchmarkWrite1(b *testing.B)     { benchmarkWrite(1, b) }
func BenchmarkWrite10(b *testing.B)    { benchmarkWrite(10, b) }
func BenchmarkWrite100(b *testing.B)   { benchmarkWrite(100, b) }
func BenchmarkWrite1000(b *testing.B)  { benchmarkWrite(1000, b) }
func BenchmarkWrite10000(b *testing.B) { benchmarkWrite(10000, b) }

func benchmarkLastN(num int, b *testing.B) {

	// write then read b.N times
	for n := 0; n < b.N; n++ {
		StoreClient.TSJSON.LastN(dsID+"benchmarkLastN", num, TimeSeriesQueryOptions{})
	}
}

//BenchmarkLastNWrite Not part of the benchmark just writes some data in to the store
func BenchmarkLastNWrite(b *testing.B) {
	for i := 0; i < 50000; i++ {
		StoreClient.TSJSON.Write(dsID+"benchmarkLastN", []byte("{\"value\":"+strconv.Itoa(i)+"}"))
	}
}
func BenchmarkLastN1(b *testing.B)     { benchmarkLastN(1, b) }
func BenchmarkLastN50(b *testing.B)    { benchmarkLastN(50, b) }
func BenchmarkLastN500(b *testing.B)   { benchmarkLastN(500, b) }
func BenchmarkLastN5000(b *testing.B)  { benchmarkLastN(5000, b) }
func BenchmarkLastN50000(b *testing.B) { benchmarkLastN(50000, b) }

func benchmarkLastNSum(num int, b *testing.B) {

	// write then read b.N times
	for n := 0; n < b.N; n++ {
		StoreClient.TSJSON.LastN(dsID+"benchmarkLastNSum", num, TimeSeriesQueryOptions{AggregationFunction: Sum})
	}
}

//BenchmarkLastNSumWrite Not part of the benchmark just writes some data in to the store
func BenchmarkLastNSumWrite(b *testing.B) {
	for i := 0; i < 50000; i++ {
		StoreClient.TSJSON.Write(dsID+"benchmarkLastNSum", []byte("{\"value\":"+strconv.Itoa(i)+"}"))
	}
}

func BenchmarkLastNSum1(b *testing.B)     { benchmarkLastNSum(1, b) }
func BenchmarkLastNSum50(b *testing.B)    { benchmarkLastNSum(50, b) }
func BenchmarkLastNSum500(b *testing.B)   { benchmarkLastNSum(500, b) }
func BenchmarkLastNSum5000(b *testing.B)  { benchmarkLastNSum(5000, b) }
func BenchmarkLastNSum50000(b *testing.B) { benchmarkLastNSum(50000, b) }

func benchmarkLastNMean(num int, b *testing.B) {

	// write then read b.N times
	for n := 0; n < b.N; n++ {
		StoreClient.TSJSON.LastN(dsID+"benchmarkLastNMean", num, TimeSeriesQueryOptions{AggregationFunction: Mean})
	}
}

//BenchmarkLastNMeanWrite Not part of the benchmark just writes some data in to the store
func BenchmarkLastNMeanWrite(b *testing.B) {
	for i := 0; i < 50000; i++ {
		StoreClient.TSJSON.Write(dsID+"benchmarkLastNMean", []byte("{\"value\":"+strconv.Itoa(i)+"}"))
	}
}

func BenchmarkLastNMean1(b *testing.B)     { benchmarkLastNMean(1, b) }
func BenchmarkLastNMean50(b *testing.B)    { benchmarkLastNMean(50, b) }
func BenchmarkLastNMean500(b *testing.B)   { benchmarkLastNMean(500, b) }
func BenchmarkLastNMean5000(b *testing.B)  { benchmarkLastNMean(5000, b) }
func BenchmarkLastNMean50000(b *testing.B) { benchmarkLastNMean(50000, b) }

func BenchmarkWriteReadMixed(b *testing.B) {
	for n := 0; n < b.N; n++ {
		StoreClient.TSJSON.Write(dsID, []byte("{\"value\":"+strconv.Itoa(n)+"}"))
		StoreClient.TSJSON.Latest(dsID)
	}
}

func TestLatest(t *testing.T) {

	err := StoreClient.TSJSON.Write(dsID, []byte("{\"value\":3.14}"))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}

	result, err := StoreClient.TSJSON.Latest(dsID)
	if err != nil {
		t.Errorf("Call to Latest failed with error %s", err.Error())
	}
	expected := []byte(`{"value":3.14}`)
	cont := s.Contains(string(result), string(expected))
	if cont != true {
		t.Errorf("Write to %s failed expected %s but got %s", dsID, expected, result)
	}
}

func TestWriteLots(t *testing.T) {

	for i := 1; i <= 10; i++ {
		err := StoreClient.TSJSON.Write(dsID, []byte("{\"value\":"+strconv.Itoa(i)+"}"))
		if err != nil {
			t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
		}
	}

	result, err := StoreClient.TSJSON.Latest(dsID)
	if err != nil {
		t.Errorf("Call to Latest failed with error %s", err.Error())
	}
	expected := []byte(`{"value":10}`)
	cont := s.Contains(string(result), string(expected))
	if cont != true {
		t.Errorf("TestWriteLots failed expected %s but got %s", expected, result)
	}
}

func TestWriteLength(t *testing.T) {

	numRecToWrite := 50
	_dsID := dsID + "TestWriteLength"

	for i := 1; i <= numRecToWrite; i++ {
		err := StoreClient.TSJSON.Write(_dsID, []byte("{\"value\":"+strconv.Itoa(i)+"}"))
		if err != nil {
			t.Errorf("Write to %s failed expected err to be nil got %s", _dsID, err.Error())
		}
	}

	result, err := StoreClient.TSJSON.Length(_dsID)
	if err != nil {
		t.Errorf("Call to Latest failed with error %s", err.Error())
	}
	if numRecToWrite != result {
		t.Errorf("TestWriteLots failed expected %d but got %d", numRecToWrite, result)
	}
}

func TestLastNWithTag(t *testing.T) {

	//now := time.Now().UnixNano() / int64(time.Millisecond)

	err := StoreClient.TSJSON.Write(dsID, []byte("{\"value\":11, \"lastNTag\":\"one\"}"))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}
	err = StoreClient.TSJSON.Write(dsID, []byte("{\"value\":12, \"lastNTag\":\"one\"}"))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}
	err = StoreClient.TSJSON.Write(dsID, []byte("{\"value\":13, \"lastNTag\":\"two\"}"))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}
	err = StoreClient.TSJSON.Write(dsID, []byte("{\"value\":14, \"lastNTag\":\"two\"}"))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}

	result, err := StoreClient.TSJSON.LastN(dsID, 99, TimeSeriesQueryOptions{
		Filter: &Filter{
			TagName:    "lastNTag",
			FilterType: "equals",
			Value:      "one",
		},
	})
	if err != nil {
		t.Errorf("Call to LastN failed with error %s", err.Error())
	}

	expected := []byte(`{"value":11`)
	cont := s.Contains(string(result), string(expected))
	if cont != true {
		t.Errorf("Call to LastN failed expected %s but got %s", expected, result)
	}

	expected = []byte(`{"value":12`)
	cont = s.Contains(string(result), string(expected))
	if cont != true {
		t.Errorf("Call to LastN failed expected %s but got %s", expected, result)
	}
}

func TestLastNWithSum(t *testing.T) {

	//now := time.Now().UnixNano() / int64(time.Millisecond)

	err := StoreClient.TSJSON.Write(dsID, []byte("{\"value\":11, \"lastNTag\":\"one\"}"))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}
	err = StoreClient.TSJSON.Write(dsID, []byte("{\"value\":12}, \"lastNTag\":\"one\""))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}
	err = StoreClient.TSJSON.Write(dsID, []byte("{\"value\":13}, \"lastNTag\":\"two\""))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}
	err = StoreClient.TSJSON.Write(dsID, []byte("{\"value\":14}, \"lastNTag\":\"two\""))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}

	result, err := StoreClient.TSJSON.LastN(dsID, 4, TimeSeriesQueryOptions{
		AggregationFunction: Sum,
	})
	if err != nil {
		t.Errorf("Call to LastN failed with error %s", err.Error())
	}

	expected := []byte(`{"result":50}`)
	cont := s.Contains(string(result), string(expected))
	if cont != true {
		t.Errorf("Call to LastN failed expected %s but got %s", expected, result)
	}
}

func TestLastNWithMean(t *testing.T) {

	//now := time.Now().UnixNano() / int64(time.Millisecond)

	err := StoreClient.TSJSON.Write(dsID, []byte("{\"value\":11.0}"))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}
	err = StoreClient.TSJSON.Write(dsID, []byte("{\"value\":12.0}"))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}
	err = StoreClient.TSJSON.Write(dsID, []byte("{\"value\":13.0}"))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}
	err = StoreClient.TSJSON.Write(dsID, []byte("{\"value\":14.0}"))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}

	result, err := StoreClient.TSJSON.LastN(dsID, 4, TimeSeriesQueryOptions{
		AggregationFunction: Mean,
	})
	if err != nil {
		t.Errorf("Call to LastN failed with error %s", err.Error())
	}

	expected := []byte(`{"result":12.5}`)
	cont := s.Contains(string(result), string(expected))
	if cont != true {
		t.Errorf("Call to LastN failed expected %s but got %s", expected, result)
	}
}

func TestLastN(t *testing.T) {

	now := time.Now().UnixNano() / int64(time.Millisecond)
	thisDsID := dsID + "TestLastN"

	err := StoreClient.TSJSON.WriteAt(thisDsID, now+20, []byte("{\"value\":11}"))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", thisDsID, err.Error())
	}
	err = StoreClient.TSJSON.WriteAt(thisDsID, now+40, []byte("{\"value\":12}"))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", thisDsID, err.Error())
	}

	result, err := StoreClient.TSJSON.LastN(thisDsID, 2, TimeSeriesQueryOptions{})
	if err != nil {
		t.Errorf("Call to LastN failed with error %s", err.Error())
	}

	expected := []byte(`{"value":11}`)
	cont := s.Contains(string(result), string(expected))
	if cont != true {
		t.Errorf("Call to LastN failed expected %s but got %s", expected, result)
	}

	expected = []byte(`{"value":12}`)
	cont = s.Contains(string(result), string(expected))
	if cont != true {
		t.Errorf("Call to LastN failed expected %s but got %s", expected, result)
	}
}

func TestEarliest(t *testing.T) {

	for i := 1; i <= 10; i++ {
		err := StoreClient.TSJSON.Write(dsID+"TestEarliest", []byte("{\"value\":"+strconv.Itoa(i)+"}"))
		if err != nil {
			t.Errorf("Write to %s failed expected err to be nil got %s", dsID+"TestEarliest", err.Error())
		}

		// On a fast CPU its possible to write faster than the resolution of zestDBs timestamp!
		// when this happens it can return records out of order causing this test to fail ;-(
		// so lets just sleep for a while to make sure!!
		// https://github.com/jptmoore/zestdb/issues/25
		time.Sleep(time.Millisecond * 10)
	}

	time.Sleep(time.Second)

	result, err := StoreClient.TSJSON.Earliest(dsID + "TestEarliest")
	if err != nil {
		t.Errorf("Call to Earliest failed with error %s", err.Error())
	}

	expected := []byte(`{"value":1}`)
	cont := s.Contains(string(result), string(expected))
	if cont != true {
		t.Errorf("Call to Earliest failed expected %s but got %s", expected, result)
	}

}

func TestFirstN(t *testing.T) {

	for i := 1; i <= 100; i++ {
		err := StoreClient.TSJSON.Write(dsID+"TestFirstN", []byte("{\"value\":"+strconv.Itoa(i)+"}"))
		if err != nil {
			t.Errorf("Write to %s failed expected err to be nil got %s", dsID+"TestFirstN", err.Error())
		}
		// On a fast CPU its possible to write faster than the resolution of zestDBs timestamp!
		// when this happens it can return records out of order causing this test to fail ;-(
		// so lets just sleep for a while to make sure!!
		// https://github.com/jptmoore/zestdb/issues/25
		time.Sleep(time.Millisecond * 10)
	}

	result, err := StoreClient.TSJSON.FirstN(dsID+"TestFirstN", 20, TimeSeriesQueryOptions{})
	if err != nil {
		t.Errorf("Call to FirstN failed with error %s", err.Error())
	}

	expected := []byte(`{"value":1}`)
	cont := s.Contains(string(result), string(expected))
	if cont != true {
		t.Errorf("Call to FirstN failed expected %s but got %s", expected, result)
	}

	expected = []byte(`{"value":20}`)
	cont = s.Contains(string(result), string(expected))
	if cont != true {
		t.Errorf("Call to FirstN failed expected %s but got %s", expected, result)
	}
}

func TestFirstNPastInternalBuffer(t *testing.T) {

	for i := 1; i <= 1000; i++ {
		err := StoreClient.TSJSON.Write(dsID+"TestFirstN", []byte("{\"value\":"+strconv.Itoa(i)+"}"))
		if err != nil {
			t.Errorf("Write to %s failed expected err to be nil got %s", dsID+"TestFirstN", err.Error())
		}
		time.Sleep(time.Millisecond * 10)
	}

	startTime := time.Now().Unix()
	result, err := StoreClient.TSJSON.FirstN(dsID+"TestFirstN", 20, TimeSeriesQueryOptions{})
	if err != nil {
		t.Errorf("Call to FirstN failed with error %s", err.Error())
	}
	queryTime := time.Now().Unix() - startTime
	t.Log("query took :: ", queryTime)

	expected := []byte(`{"value":1}`)
	cont := s.Contains(string(result), string(expected))
	if cont != true {
		t.Errorf("Call to FirstN failed expected %s but got %s", expected, result)
	}

	expected = []byte(`{"value":20}`)
	cont = s.Contains(string(result), string(expected))
	if cont != true {
		t.Errorf("Call to FirstN failed expected %s but got %s", expected, result)
	}
}

func TestWriteAtAndRange(t *testing.T) {

	now := time.Now().UnixNano() / int64(time.Millisecond)
	numRecords := 20
	timeStepMs := 50

	for i := 1; i <= numRecords; i++ {

		err := StoreClient.TSJSON.WriteAt(dsID, now+int64(timeStepMs*i), []byte("{\"value\":"+strconv.Itoa(i)+"}"))
		if err != nil {
			t.Errorf("WriteAt to %s failed expected err to be nil got %s", dsID, err.Error())
		}
		// On a fast CPU its possible to write faster than the resolution of zestDBs timestamp!
		// when this happens it can return records out of order causing this test to fail ;-(
		// so lets just sleep for a while to make sure!!
		// https://github.com/jptmoore/zestdb/issues/25
	}

	result, err := StoreClient.TSJSON.Range(dsID, now, now+int64(numRecords*timeStepMs), TimeSeriesQueryOptions{})
	if err != nil {
		t.Errorf("Call to Range failed with error %s", err.Error())
	}

	expected := []byte(`{"value":1}`)
	cont := s.Contains(string(result), string(expected))
	if cont != true {
		t.Errorf("TestWriteAtAndRange failed expected %s but got %s", expected, result)
	}
	expected = []byte(`{"value":5}`)
	cont = s.Contains(string(result), string(expected))
	if cont != true {
		t.Errorf("TestWriteAtAndRange failed expected %s but got %s", expected, result)
	}
	expected = []byte(`{"value":` + strconv.Itoa(numRecords) + `}`)
	cont = s.Contains(string(result), string(expected))
	if cont != true {
		t.Errorf("TestWriteAtAndRange failed expected %s but got %s", expected, result)
	}
}

func TestConcurrentWriteAndRead(t *testing.T) {

	doneChanWrite := make(chan int)
	doneChanRead := make(chan int)
	startAt := 1000
	numRecords := 1100

	go func() {
		for i := startAt; i <= numRecords; i++ {
			err := StoreClient.TSJSON.Write(dsID, []byte("{\"value\":"+strconv.Itoa(i)+"}"))
			if err != nil {
				t.Errorf("WriteAt to %s failed expected err to be nil got %s", dsID, err.Error())
			}
		}
		doneChanWrite <- 1
	}()

	go func() {
		for i := 1; i <= numRecords; i++ {
			_, err := StoreClient.TSJSON.Latest(dsID)
			if err != nil {
				t.Errorf("Latest failed expected err to be nil got %s", err.Error())
			}
		}
		doneChanRead <- 1
	}()

	<-doneChanWrite
	<-doneChanRead

	result, err := StoreClient.TSJSON.LastN(dsID, numRecords, TimeSeriesQueryOptions{})
	if err != nil {
		t.Errorf("Call to LastN failed with error %s", err.Error())
	}
	for i := startAt; i <= numRecords; i++ {
		expected := []byte("{\"value\":" + strconv.Itoa(i) + "}")
		cont := s.Contains(string(result), string(expected))
		if cont != true {
			t.Errorf("LastN Error '%s' does not contain  %s", string(result), string(expected))
			break
		}
	}
}

func TestObserve(t *testing.T) {

	t.Log("Hello !")
	doneChanWrite := make(chan int)
	//doneChanRead := make(chan int)
	//now := time.Now().UnixNano() / int64(time.Millisecond)
	startAt := 0
	numRecords := 100

	receivedData := []ObserveResponse{}

	go func() {
		dataChan, err := StoreClient.TSJSON.Observe(dsID)
		if err != nil {
			t.Errorf("Observing %s failed expected err to be nil got %s", dsID, err.Error())
		}

		for data := range dataChan {
			receivedData = append(receivedData, data)
			t.Log("received:: " + string(data.Data))
		}

	}()

	// Observe take a bit of time to register
	// we miss some values if we dont wait before writing
	time.Sleep(time.Second * 2)

	go func() {
		for i := startAt; i <= numRecords; i++ {
			err := StoreClient.TSJSON.Write(dsID, []byte("{\"value\":"+strconv.Itoa(i)+"}"))
			if err != nil {
				t.Errorf("WriteAt to %s failed expected err to be nil got %s", dsID, err.Error())
			}
			t.Log(string("written:: " + strconv.Itoa(i)))
		}

		// we miss some values if we dont wait before saying we are done!
		time.Sleep(time.Second * 3)
		doneChanWrite <- 1
	}()

	<-doneChanWrite

	if len(receivedData) < numRecords {
		t.Errorf("receivedData Error:  receivedData should contain '%d' items but contains  %d", numRecords, len(receivedData))
	}
	for i := startAt; i <= numRecords; i++ {
		expected := []byte("{\"value\":" + strconv.Itoa(i) + "}")
		cont := s.Contains(string(receivedData[i].Data), string(expected))
		//t.Log(receivedData[i])
		if cont != true {
			t.Errorf("receivedData Error '%s' does not contain  %s", string(receivedData[i].Data), string(expected))
			break
		}
	}

}
