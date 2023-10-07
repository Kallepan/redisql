package constant

import "testing"

type GetErrorTest struct {
	response Response
	expected error
}

var GetErrorTests = []GetErrorTest{
	{
		response: Success,
		expected: nil,
	},
	{
		response: InvalidRequest,
		expected: InvalidRequestError,
	},
	{
		response: Conflict,
		expected: ConflictError,
	},
	{
		response: NotFound,
		expected: NotFoundError,
	},
	{
		response: InternalError,

		expected: InternalErrorError,
	},
	{
		response: 0,
		expected: nil,
	},
}

func TestGetError(t *testing.T) {
	for _, test := range GetErrorTests {
		if test.response.GetError() != test.expected {
			t.Errorf("GetError(%v): expected %v, actual %v", test.response, test.expected, test.response.GetError())
		}
	}
}
