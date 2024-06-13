package pack

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type DbMock struct {
	data []int
	err  error
}

func NewDbMock(dt []int, err error) *DbMock {
	return &DbMock{data: dt, err: err}
}

func (db *DbMock) AddPacks(size []int) error { return db.err }

func (db *DbMock) RemovePack(size int) error { return db.err }

func (db *DbMock) GetPacks() []int { return nil }

func (db *DbMock) RemovePacks() {}

func TestHandler_handleAddPacks(t *testing.T) {
	type testCaseInput struct {
		dbMock         *DbMock
		requestPayload SizePayload
	}
	type testCaseOutput struct {
		status int
		err    error
	}
	type testCase struct {
		name     string
		input    testCaseInput
		expected testCaseOutput
	}

	tests := []testCase{
		{
			name: "test happy flow for adding new pack, no error returned",
			input: testCaseInput{
				dbMock:         NewDbMock([]int{}, nil),
				requestPayload: SizePayload{Sizes: []int{200}},
			},
			expected: testCaseOutput{
				status: http.StatusCreated,
				err:    nil,
			},
		},
		{
			name: "test adding pack with value 0, error returned",
			input: testCaseInput{
				dbMock:         NewDbMock([]int{}, fmt.Errorf("an error has occurred")),
				requestPayload: SizePayload{Sizes: []int{0}},
			},
			expected: testCaseOutput{
				status: http.StatusInternalServerError,
				err:    fmt.Errorf("an error has occurred"),
			},
		},
		{
			name: "test adding pack with value -1, error returned",
			input: testCaseInput{
				dbMock:         NewDbMock([]int{}, fmt.Errorf("an error has occurred")),
				requestPayload: SizePayload{Sizes: []int{-1}},
			},
			expected: testCaseOutput{
				status: http.StatusInternalServerError,
				err:    fmt.Errorf("an error has occurred"),
			},
		},
		{
			name: "test error adding pack to DB, error returned",
			input: testCaseInput{
				dbMock:         NewDbMock([]int{}, errors.New("an error has occurred")),
				requestPayload: SizePayload{Sizes: []int{1}},
			},
			expected: testCaseOutput{
				status: http.StatusInternalServerError,
				err:    errors.New("an error has occurred"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				db: tt.input.dbMock,
			}

			payload, _ := json.Marshal(&tt.input.requestPayload)
			req := httptest.NewRequest(http.MethodPost, "/pack", io.NopCloser(bytes.NewBuffer(payload)))

			w := httptest.NewRecorder()
			h.handleAddPacks(w, req)

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
			}

			assert.Equal(t, tt.expected.status, resp.StatusCode)
		})
	}
}

func TestHandler_handleRemovePack(t *testing.T) {
	type testCaseInput struct {
		dbMock *DbMock
		size   string
	}
	type testCaseOutput struct {
		status int
		err    error
	}
	type testCase struct {
		name     string
		input    testCaseInput
		expected testCaseOutput
	}

	tests := []testCase{
		{
			name: "test happy flow for removing 200 pack, no error returned",
			input: testCaseInput{
				dbMock: NewDbMock([]int{200}, nil),
				size:   "200",
			},
			expected: testCaseOutput{
				status: http.StatusOK,
				err:    nil,
			},
		},
		{
			name: "test removing pack with value 0, error returned",
			input: testCaseInput{
				dbMock: NewDbMock([]int{200}, nil),
				size:   "0",
			},
			expected: testCaseOutput{
				status: http.StatusBadRequest,
				err:    fmt.Errorf("you must provide a positive value"),
			},
		},
		{
			name: "test removing pack with value -1, error returned",
			input: testCaseInput{
				dbMock: NewDbMock([]int{200}, nil),
				size:   "-1",
			},
			expected: testCaseOutput{
				status: http.StatusBadRequest,
				err:    fmt.Errorf("you must provide a positive value"),
			},
		},
		{
			name: "test removing pack with invalid, error returned",
			input: testCaseInput{
				dbMock: NewDbMock([]int{200}, nil),
				size:   "test",
			},
			expected: testCaseOutput{
				status: http.StatusBadRequest,
				err:    fmt.Errorf("please provide a numeric value"),
			},
		},
		{
			name: "test error removing value from DB, error returned",
			input: testCaseInput{
				dbMock: NewDbMock([]int{400}, errors.New("an error has occurred")),
				size:   "200",
			},
			expected: testCaseOutput{
				status: http.StatusInternalServerError,
				err:    errors.New("an error has occurred"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				db: tt.input.dbMock,
			}

			req := httptest.NewRequest(http.MethodDelete, "/pack", nil)
			req.SetPathValue("size", tt.input.size)

			w := httptest.NewRecorder()
			h.handleRemovePack(w, req)

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
			}

			assert.Equal(t, tt.expected.status, resp.StatusCode)
		})
	}
}
