package test

import (
	"testing"
)

var relationInputs = []test{
	{
		Input: `["table1"]`,
	},
	{
		Input: `["table1", "table2"]`,
	},
}

func TestRelation(t *testing.T) {
	// for i, test := range relationInputs {
	// 	f := query.ParseRelation(test.Input)

	// 	fmt.Println("-------------------------------------------------------")
	// 	fmt.Println("Relation input:", i)
	// 	fmt.Println("-------------------------------------------------------")
	// 	fmt.Println("f", utils.PrettyPrint(f))

	// }

}
