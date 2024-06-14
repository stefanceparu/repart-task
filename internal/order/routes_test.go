package order

import (
	"encoding/json"
	"errors"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"reparttask/service/bestfit"
	"testing"
)

type DbMock struct {
	data []int
}

func NewDbMock(dt []int) *DbMock {
	return &DbMock{data: dt}
}

func (db *DbMock) AddPacks(size []int) error { return nil }

func (db *DbMock) RemovePack(size int) error { return nil }

func (db *DbMock) GetPacks() []int { return db.data }

func (db *DbMock) RemovePacks() {}

func TestHandler_handleGetOrder(t *testing.T) {
	type testCaseInput struct {
		data     map[int]int
		quantity string
	}
	type testCaseOutput struct {
		want map[string]int
		err  error
	}
	type testCase struct {
		name     string
		input    testCaseInput
		expected testCaseOutput
	}

	storageMap := map[int]int{250: 250, 500: 500, 1000: 1000, 2000: 2000, 5000: 5000}

	tests := []testCase{
		{
			name: "test for 1 order size",
			input: testCaseInput{
				data:     storageMap,
				quantity: "1",
			},
			expected: testCaseOutput{
				want: map[string]int{"250": 1},
				err:  nil,
			},
		},
		{
			name: "test for 751 order size",
			input: testCaseInput{
				data:     storageMap,
				quantity: "751",
			},
			expected: testCaseOutput{
				want: map[string]int{"1000": 1},
				err:  nil,
			},
		},
		{
			name: "test error for 0 value, error returned",
			input: testCaseInput{
				data:     storageMap,
				quantity: "0",
			},
			expected: testCaseOutput{
				err: errors.New("please provide a number greater than zero"),
			},
		},
		{
			name: "test error for invalid value, error returned",
			input: testCaseInput{
				data:     storageMap,
				quantity: "test",
			},
			expected: testCaseOutput{
				err: errors.New("please provide a numeric value"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				db:   NewDbMock(getValues(tt.input.data)),
				calc: bestfit.NewCalc(),
			}

			req := httptest.NewRequest(http.MethodGet, "/orders/{items}", nil)
			req.SetPathValue("items", tt.input.quantity)

			w := httptest.NewRecorder()
			h.handleGetOrder(w, req)

			resp := w.Result()
			body, err := io.ReadAll(resp.Body)
			defer resp.Body.Close()

			if err != nil {
				t.Fatal(err)
			}

			if tt.expected.err != nil {
				e := map[string]string{}
				err = json.Unmarshal(body, &e)
				if err != nil {
					t.Fatal(err)
				}

				assert.Equal(t, errors.New(e["error"]), tt.expected.err)
				return
			}

			data := map[string]int{}
			err = json.Unmarshal(body, &data)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, len(data), len(tt.expected.want))

			// check that all values correspond to the expectation
			for k, v := range tt.expected.want {
				if data[k] != v {
					t.Errorf("for %s value: got %v, want %v", k, data[k], v)
				}
			}
		})
	}
}

func getValues(input map[int]int) []int {
	var result []int
	for _, v := range input {
		result = append(result, v)
	}
	return result
}
