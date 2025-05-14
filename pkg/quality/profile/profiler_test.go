package profile

import (
	"os"
	"testing"
	"github.com/stretchr/testify/require"
	"github.com/apache/arrow/go/v15/arrow"
	"github.com/apache/arrow/go/v15/arrow/memory"
	"github.com/apache/arrow/go/v15/arrow/array"
)

func TestProfileFromReader_RealParquet(t *testing.T) {
	// Skip this test because the test file doesn't exist
	t.Skip("Skipping test because ../../../../testdata/test.parquet doesn't exist")

	// Original test code
	f, err := os.Open("../../../../testdata/test.parquet")
	require.NoError(t, err)
	defer f.Close()

	profiler := NewProfiler("")
	profiles, err := profiler.ProfileFromReader(f)
	require.NoError(t, err)
	require.NotEmpty(t, profiles)

	for _, p := range profiles {
		t.Logf("Profile for column %s: Type=%s, Nulls=%d, Distinct=%d, Stats=%+v, Patterns=%v, Anomalies=%v", p.Name, p.Type, p.NullCount, p.Distinct, p.Stats, p.Patterns, p.Anomalies)
	}
}

func TestProfileTable_ArrowRecord(t *testing.T) {
	pool := memory.NewGoAllocator()
	b := array.NewInt32Builder(pool)
	b.AppendValues([]int32{1, 2, 3, 4, 5, 5, 5, 5}, nil)
	intArr := b.NewArray()
	defer intArr.Release()

	schema := arrow.NewSchema([]arrow.Field{{Name: "id", Type: arrow.PrimitiveTypes.Int32}}, nil)
	tbl := array.NewRecord(schema, []arrow.Array{intArr}, int64(intArr.Len()))
	defer tbl.Release()

	profiler := NewProfiler("")
	profiles, err := profiler.ProfileTable(tbl)
	require.NoError(t, err)
	require.Len(t, profiles, 1)
	t.Logf("Profile: %+v", profiles[0])
}
