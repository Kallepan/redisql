package app

import (
	"context"
	"testing"

	"github.com/kallepan/redisql/constant"
	"github.com/kallepan/redisql/mock"
	"github.com/redis/go-redis/v9"
)

func setUp() *redis.Client {
	c := mock.InitMockDB()

	c.FlushAll(context.Background())

	return c
}

type AddRowsTest struct {
	table Table
	id    string
	data  map[string]interface{}
}

var AddRowsTests = []AddRowsTest{
	{
		table: Table{Name: "test_table"},
		id:    "test_id",
		data: map[string]interface{}{
			"test_key":  "test_data",
			"test_key2": "test_data2",
		},
	},
	{
		table: Table{Name: "test_table"},
		id:    "test_id",
		data: map[string]interface{}{
			"test_key":  "test_data",
			"test_key2": "test_data2",
		},
	},
	{
		table: Table{Name: "test_table"},
		id:    "test_new",
		data: map[string]interface{}{
			"test_key":  "test_data",
			"test_key2": "test_data2",
		},
	},
}

func TestAddRows(t *testing.T) {
	ctx := context.Background()
	client := setUp()

	for _, test := range AddRowsTests {
		// Create Row
		err := test.table.AddRows(ctx, client, test.id, test.data)
		if err != constant.Success.GetError() {
			t.Errorf("AddRows(%s, %s): expected %v, actual %v", test.id, test.data, constant.Success.GetError(), err)
		}

		// Check if Row was created
		data, err := test.table.GetRow(ctx, client, test.id)
		if err != constant.Success.GetError() {
			t.Errorf("GetRow(%s): expected %v, actual %v", test.id, constant.Success, err)
		}

		for key, value := range test.data {
			if data[key] != value {
				t.Errorf("GetRow(%s): expected %v, actual %v", test.id, value, data[key])
			}
		}
	}
}

type AddRowTest struct {
	table Table
	id    string
	key   string
	data  interface{}

	expected constant.Response
}

var AddRowTests = []AddRowTest{
	{
		table:    Table{Name: "test_table"},
		id:       "test_id",
		key:      "test_key",
		data:     "test_data",
		expected: constant.Success,
	},
	{
		table:    Table{Name: "test_table"},
		id:       "test_id",
		key:      "test_key2",
		data:     "test_data2",
		expected: constant.Success,
	},
	{
		table:    Table{Name: "test_table"},
		id:       "test_new",
		key:      "test_key",
		data:     "test_data",
		expected: constant.Success,
	},
	{
		table:    Table{Name: "test_table"},
		id:       "test_id",
		key:      "test_key",
		data:     "test_data",
		expected: constant.Conflict,
	},
	{
		table:    Table{Name: "new_table"},
		id:       "test_id",
		key:      "test_key",
		data:     "test_data",
		expected: constant.Success,
	},
}

func TestAddRow(t *testing.T) {
	ctx := context.Background()
	client := setUp()

	for _, test := range AddRowTests {
		// Create Row
		err := test.table.AddRow(ctx, client, test.id, test.key, test.data)
		if err != test.expected.GetError() {
			t.Errorf("AddRow(%s, %s, %s): expected %v, actual %v", test.id, test.key, test.data, test.expected.GetError(), err)
		}

		// Check if Row was created
		data, err := test.table.GetRow(ctx, client, test.id)
		if err != constant.Success.GetError() {
			t.Errorf("GetRow(%s): expected %v, actual %v", test.id, constant.Success, err)
		}

		if data[test.key] != test.data {
			t.Errorf("GetRow(%s): expected %v, actual %v", test.id, test.data, data[test.key])
		}
	}
}

type DeleteRowTest struct {
	table    Table
	id       string
	key      string
	data     interface{} // create row with this data
	expected constant.Response
}

var DeleteRowTests = []DeleteRowTest{
	{
		table:    Table{Name: "test"},
		id:       "test_id",
		key:      "test_key",
		data:     "test_data",
		expected: constant.Success,
	},
	{
		table:    Table{Name: "test"},
		id:       "test_id",
		key:      "test_key2",
		data:     "test_data2",
		expected: constant.Success,
	},
	{
		table:    Table{Name: "test"},
		id:       "test_new",
		key:      "test_key",
		data:     "test_data",
		expected: constant.Success,
	},
	{
		table:    Table{Name: "test"},
		id:       "test_id",
		key:      "test_not_exist",
		data:     "test_data",
		expected: constant.NotFound,
	},
}

func TestDeleteRow(t *testing.T) {
	ctx := context.Background()
	client := setUp()

	// setup
	table := Table{Name: "test"}
	if err := table.AddRow(ctx, client, "test_id", "test_key", "test_data"); err != constant.Success.GetError() {
		t.Errorf("AddRow(%s, %s, %s): expected %v, actual %v", "test_id", "test_key", "test_data", constant.Success.GetError(), err)
	}
	if err := table.AddRow(ctx, client, "test_id", "test_key2", "test_data2"); err != constant.Success.GetError() {
		t.Errorf("AddRow(%s, %s, %s): expected %v, actual %v", "test_id", "test_key2", "test_data2", constant.Success.GetError(), err)
	}
	if err := table.AddRow(ctx, client, "test_new", "test_key", "test_data"); err != constant.Success.GetError() {
		t.Errorf("AddRow(%s, %s, %s): expected %v, actual %v", "test_new", "test_key", "test_data", constant.Success.GetError(), err)
	}

	for i, test := range DeleteRowTests {
		t.Logf("Test %d", i)

		// Delete Row
		err := test.table.DeleteRow(ctx, client, test.id, test.key)
		if err != test.expected.GetError() {
			t.Errorf("DeleteRow(%s, %s): expected %v, actual %v", test.id, test.key, test.expected.GetError(), err)
		}

		// Check if Row was deleted
		data, _ := test.table.GetRow(ctx, client, test.id)
		if data[test.key] != nil {
			t.Errorf("GetRow(%s): expected %v, actual %v", test.id, nil, data[test.key])
		}
	}
}
