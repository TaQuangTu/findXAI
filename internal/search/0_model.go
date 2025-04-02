package search

import "time"

type (
	KeyBucket struct {
		NumberOfPartition int
		NumberOfRecords   int
		PartitionId       int
		PartitionAvg      float64
	}

	KeyBucketList []KeyBucket
)

func (k *KeyBucketList) Avg() int {
	var (
		total float64
	)
	for _, bucket := range *k {
		total += bucket.PartitionAvg
	}
	return int(int(total) / len(*k))
}

type (
	AvailableKey struct {
		ResetedAt time.Time
		ApiKey    string
		EngineId  string
	}

	Key struct {
		Id             int64
		Name           string
		ApiKey         string
		SearchEngineId string
		IsActive       bool
		DailyQueries   int32
		StatusCode     int32
		ErrorMsg       string
		CreatedAt      time.Time
		UpdatedAt      time.Time
	}
)
