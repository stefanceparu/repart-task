package bestfit

import (
	"log"
	"testing"
)

func Test_CalculatePacks(t *testing.T) {
	type testCaseInput struct {
		input         []int
		orderQuantity int
	}
	type testCaseOutput struct {
		want map[int]int
	}
	type testCase struct {
		name     string
		input    testCaseInput
		expected testCaseOutput
	}

	tests := []testCase{
		{
			name: "test for 1 order size",
			input: testCaseInput{
				input:         []int{250, 2000, 500, 1000, 5000},
				orderQuantity: 1,
			},
			expected: testCaseOutput{
				want: map[int]int{250: 1},
			},
		},
		{
			name: "test for 250 order size",
			input: testCaseInput{
				input:         []int{250, 2000, 500, 1000, 5000},
				orderQuantity: 250,
			},
			expected: testCaseOutput{
				want: map[int]int{250: 1},
			},
		},
		{
			name: "test for 251 order size",
			input: testCaseInput{
				input:         []int{250, 2000, 500, 1000, 5000},
				orderQuantity: 251,
			},
			// qty = 251
			// 1) (1 * 250) + (1 * 250) = 500 ===> 500 - 251 = 249 items left
			// 2) 1 * 500 = 500 			  ===> 500 - 251 = 249 items left
			// both are equal but the second option has fewer packages
			expected: testCaseOutput{
				want: map[int]int{500: 1},
			},
		},
		{
			name: "test for 501 order size",
			input: testCaseInput{
				input:         []int{250, 2000, 500, 1000, 5000},
				orderQuantity: 501,
			},
			// qty = 501
			// 1) (1 * 250) + (1 * 500) = 750 ===> 750 - 501 = 249 items left
			// 2) 1 * 1000 = 1000 			  ===> 1000 - 501 = 499 items left
			// first option has fewer items left
			expected: testCaseOutput{
				want: map[int]int{250: 1, 500: 1},
			},
		},
		{
			name: "test for 12001 order size",
			input: testCaseInput{
				input:         []int{250, 2000, 500, 1000, 5000},
				orderQuantity: 12001,
			},
			expected: testCaseOutput{
				want: map[int]int{5000: 2, 2000: 1, 250: 1},
			},
		},
		{
			name: "test for 751 order size",
			input: testCaseInput{
				input:         []int{250, 2000, 500, 1000, 5000},
				orderQuantity: 751,
			},
			// qty = 751
			// 1) (2 * 250) + (1 * 500) = 1000 ===> 1000 - 751 = 249 items left
			// 2) 2 * 500  = 1000 			   ===> 1000 - 751 = 249 items left
			// 3) 1 * 1000 = 1000 			   ===> 1000 - 751 = 249 items left
			// all return the same value but third option has fewer packs used
			expected: testCaseOutput{
				want: map[int]int{1000: 1},
			},
		},
		{
			name: "test for 1251 order size",
			input: testCaseInput{
				input:         []int{250, 2000, 500, 1000, 5000},
				orderQuantity: 1251,
			},
			expected: testCaseOutput{
				want: map[int]int{1000: 1, 500: 1},
			},
		},
		{
			name: "test for negative order size",
			input: testCaseInput{
				input:         []int{250, 2000, 500, 1000, 5000},
				orderQuantity: -1,
			},
			expected: testCaseOutput{
				want: map[int]int{},
			},
		},
		{
			name: "test for different input pack sizes",
			input: testCaseInput{
				input:         []int{23, 37, 45, 100, 500},
				orderQuantity: 46,
			},
			expected: testCaseOutput{
				// qty = 46
				// 1) (2 * 23) = 46 ==> 46 - 46 = 0 items left
				// 2) 1 * 100 = 100 ==> 100 - 46 = 54 items left
				// the first option has fewer items left.
				want: map[int]int{23: 2},
			},
		},
		{
			name: "test 2 for different input pack sizes",
			input: testCaseInput{
				input:         []int{23, 70, 73, 100, 500},
				orderQuantity: 69,
			},
			expected: testCaseOutput{
				// qty = 69
				// 1) 3 * 23 = 69 				==> 69 - 69 = 0 items left
				// 2) (1 * 70) + (1 * 23) = 93  ==> 93 - 69 = 24 items left
				// the first option has fewer items left.
				want: map[int]int{23: 3},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			c := NewCalc()
			got := c.CalculatePacks(tt.input.input, tt.input.orderQuantity)

			log.Println(got)
			// check expected values one by one
			for k, v := range tt.expected.want {
				if got[k] != v {
					t.Errorf("for required pack: %d: got %v, want %v", k, got[k], v)
					t.FailNow()
				}
			}
		})
	}
}

func Test_combinePacks(t *testing.T) {
	type testCaseInput struct {
		packs []int
		m     map[int]int
	}

	tests := []struct {
		name  string
		input testCaseInput
		want  map[int]int
	}{
		{
			name: "test combine packs",
			input: testCaseInput{
				packs: []int{250, 500, 1000, 5000},
				m:     map[int]int{250: 2, 1000: 1},
			},
			want: map[int]int{500: 1, 1000: 1},
		},
		{
			name: "test recursion, 250 x 2 can be combined into 500 x 1 this results 2 x 500 that can be combined into 1000",
			input: testCaseInput{
				packs: []int{250, 500, 1000, 5000},
				m:     map[int]int{250: 2, 500: 1},
			},
			want: map[int]int{1000: 1},
		},
		{
			name: "test with different packs for which we also have larger ones",
			input: testCaseInput{
				packs: []int{33, 46, 66, 92},
				m:     map[int]int{33: 3, 46: 2},
			},
			want: map[int]int{33: 1, 66: 1, 92: 1},
		},
		{
			name: "test with no larger packs available",
			input: testCaseInput{
				packs: []int{33, 46, 73, 100},
				m:     map[int]int{33: 3, 46: 2},
			},
			want: map[int]int{33: 3, 46: 2},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := combinePacks(tt.input.m, tt.input.packs)
			for k, v := range tt.want {
				if got[k] != v {
					t.Errorf("for required pack: %d: got %v, want %v", k, got[k], v)
					t.FailNow()
				}
			}
		})
	}
}
