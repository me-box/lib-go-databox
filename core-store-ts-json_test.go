package libDatabox

import (
	"strconv"
	s "strings"
	"testing"
	"time"
)

func TestAggregationFunctionMinOnEmptyDS(t *testing.T) {

	_, err := tsc.LastN(dsID+"TestAggregationFunctionOnEmptyDS", 4, JSONTimeSeriesQueryOptions{
		AggregationFunction: Min,
	})
	if err != nil {
		t.Errorf("Call to LastN failed with error %s", err.Error())
	}

}

func TestAggregationFunctionMaxOnEmptyDS(t *testing.T) {

	_, err := tsc.LastN(dsID+"TestAggregationFunctionOnEmptyDS", 4, JSONTimeSeriesQueryOptions{
		AggregationFunction: Max,
	})
	if err != nil {
		t.Errorf("Call to LastN failed with error %s", err.Error())
	}

}

func TestAggregationFunctionSumOnEmptyDS(t *testing.T) {

	_, err := tsc.LastN(dsID+"TestAggregationFunctionOnEmptyDS", 4, JSONTimeSeriesQueryOptions{
		AggregationFunction: Sum,
	})
	if err != nil {
		t.Errorf("Call to LastN failed with error %s", err.Error())
	}

}

func TestAggregationFunctionSDOnEmptyDS(t *testing.T) {

	_, err := tsc.LastN(dsID+"TestAggregationFunctionOnEmptyDS", 4, JSONTimeSeriesQueryOptions{
		AggregationFunction: StandardDeviation,
	})
	if err != nil {
		t.Errorf("Call to LastN failed with error %s", err.Error())
	}

}

func TestWrite(t *testing.T) {
	err := tsc.Write(dsID, []byte("{\"value\":3.1415}"))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}
}

func benchmarkWrite(num int, b *testing.B) {
	// write b.N times
	for n := 0; n < b.N; n++ {
		for i := 0; i < num; i++ {
			tsc.Write(dsID+"benchmarkWrite"+strconv.Itoa(num), []byte("{\"value\":"+strconv.Itoa(n)+"}"))
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
		tsc.LastN(dsID+"benchmarkLastN", num, JSONTimeSeriesQueryOptions{})
	}
}

//BenchmarkLastNWrite Not part of the benchmark just writes some data in to the store
func BenchmarkLastNWrite(b *testing.B) {
	for i := 0; i < 50000; i++ {
		tsc.Write(dsID+"benchmarkLastN", []byte("{\"value\":"+strconv.Itoa(i)+"}"))
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
		tsc.LastN(dsID+"benchmarkLastNSum", num, JSONTimeSeriesQueryOptions{AggregationFunction: Sum})
	}
}

//BenchmarkLastNSumWrite Not part of the benchmark just writes some data in to the store
func BenchmarkLastNSumWrite(b *testing.B) {
	for i := 0; i < 50000; i++ {
		tsc.Write(dsID+"benchmarkLastNSum", []byte("{\"value\":"+strconv.Itoa(i)+"}"))
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
		tsc.LastN(dsID+"benchmarkLastNMean", num, JSONTimeSeriesQueryOptions{AggregationFunction: Mean})
	}
}

//BenchmarkLastNMeanWrite Not part of the benchmark just writes some data in to the store
func BenchmarkLastNMeanWrite(b *testing.B) {
	for i := 0; i < 50000; i++ {
		tsc.Write(dsID+"benchmarkLastNMean", []byte("{\"value\":"+strconv.Itoa(i)+"}"))
	}
}

func BenchmarkLastNMean1(b *testing.B)     { benchmarkLastNMean(1, b) }
func BenchmarkLastNMean50(b *testing.B)    { benchmarkLastNMean(50, b) }
func BenchmarkLastNMean500(b *testing.B)   { benchmarkLastNMean(500, b) }
func BenchmarkLastNMean5000(b *testing.B)  { benchmarkLastNMean(5000, b) }
func BenchmarkLastNMean50000(b *testing.B) { benchmarkLastNMean(50000, b) }

func BenchmarkWriteReadMixed(b *testing.B) {
	for n := 0; n < b.N; n++ {
		tsc.Write(dsID, []byte("{\"value\":"+strconv.Itoa(n)+"}"))
		tsc.Latest(dsID)
	}
}

func TestLatest(t *testing.T) {

	err := tsc.Write(dsID, []byte("{\"value\":3.14}"))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}

	result, err := tsc.Latest(dsID)
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
		err := tsc.Write(dsID, []byte("{\"value\":"+strconv.Itoa(i)+"}"))
		if err != nil {
			t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
		}
	}

	result, err := tsc.Latest(dsID)
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

	for i := 1; i <= numRecToWrite; i++ {
		err := tsc.Write(dsID, []byte("{\"value\":"+strconv.Itoa(i)+"}"))
		if err != nil {
			t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
		}
	}

	result, err := tsc.Length(dsID)
	if err != nil {
		t.Errorf("Call to Latest failed with error %s", err.Error())
	}
	if numRecToWrite != result {
		t.Errorf("TestWriteLots failed expected %d but got %d", numRecToWrite, result)
	}
}

//TODO this fails looks like a timing thing
/*func TestWriteThenWriteAT(t *testing.T) {

	now := time.Now().UnixNano() / int64(time.Millisecond)

	for i := 1; i <= 10; i++ {
		err := tsbc.Write(dsID, []byte("{\"TestWriteThenWriteAT\":\"data"+strconv.Itoa(i)+"\"}"))
		if err != nil {
			t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
		}
	}

	//time.Sleep(time.Second * 2)

	now = time.Now().UnixNano() / int64(time.Millisecond)
	fmt.Println(now + 1000)
	err := tsbc.WriteAt(dsID, now+1000, []byte("{\"TestWriteThenWriteAT\":\"data11\"}"))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}

	err = tsbc.WriteAt(dsID, now+1001, []byte("{\"TestWriteThenWriteAT\":\"data12\"}"))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}

	//time.Sleep(time.Second * 5)

	result, err := tsbc.LastN(dsID, 2)
	if err != nil {
		t.Errorf("Call to LastN failed with error %s", err.Error())
	}

	expected := []byte(`{"TestWriteThenWriteAT":"data11"}`)
	cont := s.Contains(string(result), string(expected))
	if cont != true {
		t.Errorf("Call to LastN failed expected %s but got %s", expected, result)
	}

	expected = []byte(`{"TestWriteThenWriteAT":"data12"}`)
	cont = s.Contains(string(result), string(expected))
	if cont != true {
		t.Errorf("Call to LastN failed expected %s but got %s", expected, result)
	}

}*/

func TestLastNWithTag(t *testing.T) {

	//now := time.Now().UnixNano() / int64(time.Millisecond)

	err := tsc.Write(dsID, []byte("{\"value\":11, \"lastNTag\":\"one\"}"))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}
	err = tsc.Write(dsID, []byte("{\"value\":12, \"lastNTag\":\"one\"}"))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}
	err = tsc.Write(dsID, []byte("{\"value\":13, \"lastNTag\":\"two\"}"))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}
	err = tsc.Write(dsID, []byte("{\"value\":14, \"lastNTag\":\"two\"}"))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}

	result, err := tsc.LastN(dsID, 99, JSONTimeSeriesQueryOptions{
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

	err := tsc.Write(dsID, []byte("{\"value\":11, \"lastNTag\":\"one\"}"))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}
	err = tsc.Write(dsID, []byte("{\"value\":12}, \"lastNTag\":\"one\""))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}
	err = tsc.Write(dsID, []byte("{\"value\":13}, \"lastNTag\":\"two\""))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}
	err = tsc.Write(dsID, []byte("{\"value\":14}, \"lastNTag\":\"two\""))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}

	result, err := tsc.LastN(dsID, 4, JSONTimeSeriesQueryOptions{
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

	err := tsc.Write(dsID, []byte("{\"value\":11.0}"))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}
	err = tsc.Write(dsID, []byte("{\"value\":12.0}"))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}
	err = tsc.Write(dsID, []byte("{\"value\":13.0}"))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}
	err = tsc.Write(dsID, []byte("{\"value\":14.0}"))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", dsID, err.Error())
	}

	result, err := tsc.LastN(dsID, 4, JSONTimeSeriesQueryOptions{
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

	err := tsc.WriteAt(thisDsID, now+20, []byte("{\"value\":11}"))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", thisDsID, err.Error())
	}
	err = tsc.WriteAt(thisDsID, now+40, []byte("{\"value\":12}"))
	if err != nil {
		t.Errorf("Write to %s failed expected err to be nil got %s", thisDsID, err.Error())
	}

	result, err := tsc.LastN(thisDsID, 2, JSONTimeSeriesQueryOptions{})
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
		err := tsc.Write(dsID+"TestEarliest", []byte("{\"value\":"+strconv.Itoa(i)+"}"))
		if err != nil {
			t.Errorf("Write to %s failed expected err to be nil got %s", dsID+"TestEarliest", err.Error())
		}
		//time.Sleep(time.Millisecond * 10)
	}

	result, err := tsc.Earliest(dsID + "TestEarliest")
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
		err := tsc.Write(dsID+"TestFirstN", []byte("{\"value\":"+strconv.Itoa(i)+"}"))
		if err != nil {
			t.Errorf("Write to %s failed expected err to be nil got %s", dsID+"TestFirstN", err.Error())
		}
		time.Sleep(time.Millisecond * 10)
	}

	result, err := tsc.FirstN(dsID+"TestFirstN", 20, JSONTimeSeriesQueryOptions{})
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
		err := tsc.Write(dsID+"TestFirstN", []byte("{\"value\":"+strconv.Itoa(i)+"}"))
		if err != nil {
			t.Errorf("Write to %s failed expected err to be nil got %s", dsID+"TestFirstN", err.Error())
		}
		time.Sleep(time.Millisecond * 10)
	}

	startTime := time.Now().Unix()
	result, err := tsc.FirstN(dsID+"TestFirstN", 20, JSONTimeSeriesQueryOptions{})
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

		err := tsc.WriteAt(dsID, now+int64(timeStepMs*i), []byte("{\"value\":"+strconv.Itoa(i)+"}"))
		if err != nil {
			t.Errorf("WriteAt to %s failed expected err to be nil got %s", dsID, err.Error())
		}
	}

	result, err := tsc.Range(dsID, now, now+int64(numRecords*timeStepMs), JSONTimeSeriesQueryOptions{})
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

func TestRegisterDatasource(t *testing.T) {

	dsmd := DataSourceMetadata{
		DataSourceID:   dsID,
		Vendor:         "testing",
		ContentType:    "application/json",
		StoreType:      "ts",
		Description:    "A test DS",
		DataSourceType: "test",
	}

	err := tsc.RegisterDatasource(dsmd)
	if err != nil {
		t.Errorf("RegisterDatasource failed expected err to be nil got %s", err.Error())
	}

	catByteArray, getErr := tsc.GetDatasourceCatalogue()
	if getErr != nil {
		t.Errorf("GetDatasourceCatalogue failed expected err to be nil got %s", getErr.Error())
	}

	dsmdByteArray, _ := dataSourceMetadataToHypercat(dsmd, "tcp://127.0.0.1:5555/ts/")
	cont := s.Contains(string(catByteArray), string(dsmdByteArray))
	if cont != true {
		t.Errorf("GetDatasourceCatalogue Error '%s' does not contain  %s", string(catByteArray), string(dsmdByteArray))
	}
}

func TestConcurrentWriteAndRead(t *testing.T) {

	doneChanWrite := make(chan int)
	doneChanRead := make(chan int)
	//now := time.Now().UnixNano() / int64(time.Millisecond)
	startAt := 1000
	numRecords := 1100

	go func() {
		for i := startAt; i <= numRecords; i++ {
			err := tsc.Write(dsID, []byte("{\"value\":"+strconv.Itoa(i)+"}"))
			if err != nil {
				t.Errorf("WriteAt to %s failed expected err to be nil got %s", dsID, err.Error())
			}
			//fmt.Println(string("written:: " + strconv.Itoa(i)))
		}
		doneChanWrite <- 1
	}()

	go func() {
		for i := 1; i <= numRecords; i++ {
			_, err := tsc.Latest(dsID)
			if err != nil {
				t.Errorf("Latest failed expected err to be nil got %s", err.Error())
			}
			//fmt.Println("Got:: ", string(data))
		}
		doneChanRead <- 1
	}()

	<-doneChanWrite
	<-doneChanRead

	result, err := tsc.LastN(dsID, numRecords, JSONTimeSeriesQueryOptions{})
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

	receivedData := []string{}

	go func() {
		dataChan, err := tsc.Observe(dsID)
		if err != nil {
			t.Errorf("Observing %s failed expected err to be nil got %s", dsID, err.Error())
		}

		for data := range dataChan {
			receivedData = append(receivedData, string(data))
			t.Log("received:: " + string(data))
		}

	}()

	//Observe take a bit of time to register we miss some values if we dont wait before writing
	time.Sleep(time.Second)

	go func() {
		for i := startAt; i <= numRecords; i++ {
			err := tsc.Write(dsID, []byte("{\"value\":"+strconv.Itoa(i)+"}"))
			if err != nil {
				t.Errorf("WriteAt to %s failed expected err to be nil got %s", dsID, err.Error())
			}
			t.Log(string("written:: " + strconv.Itoa(i)))
		}
		doneChanWrite <- 1
	}()

	<-doneChanWrite

	for i := startAt; i <= numRecords; i++ {
		expected := []byte("{\"value\":" + strconv.Itoa(i) + "}")
		cont := s.Contains(receivedData[i], string(expected))
		t.Log(receivedData[i])
		if cont != true {
			t.Errorf("receivedData Error '%s' does not contain  %s", string(receivedData[i]), string(expected))
			break
		}
	}

}
