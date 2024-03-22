package snowflake

import (
	"github.com/bwmarrin/snowflake"
)

// BUCKET_SIZE is the number of seconds in a bucket.
const BUCKET_SIZE = 1000 * 60 * 60 * 24 * 10

// MakeBucket returns the bucket number based on the input ID.
// It takes an integer ID as a parameter and returns an integer value.
func MakeBucket(id int64) int64 {
	snowflakeID := snowflake.ParseInt64(id)

	return snowflakeID.Time() / BUCKET_SIZE
}

// MakeBuckets generates a slice of integers representing the buckets between the startID and endID.
func MakeBuckets(startID, endID int64) []int64 {
	start := MakeBucket(startID)
	end := MakeBucket(endID)
	result := make([]int64, 0, end-start+1)
	for i := start; i <= end; i++ {
		result = append(result, i)
	}

	return result
}
