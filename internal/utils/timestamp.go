package utils

import (
	"time"
)

// timestamps are formatted according to RFC 3339
type Timestamp time.Time

func TimeNow() Timestamp {
	return Timestamp(time.Now().UTC())
}

func (t Timestamp) String() string {
	return time.Time(t).Format(time.RFC3339Nano)
}

func StringToTimestamp(input string) (Timestamp, error) {

	parsedTime, err := time.Parse(time.RFC3339Nano, input)
	if err != nil {
		return Timestamp{}, err
	}

	return Timestamp(parsedTime), nil 

} 
