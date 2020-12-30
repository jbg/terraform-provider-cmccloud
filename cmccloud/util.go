package cmccloud

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// IfThenElse one line if else condition: IfThenElse(1 == 1, "Yes", false) => "Yes"
func IfThenElse(condition bool, a interface{}, b interface{}) interface{} {
	if condition {
		return a
	}
	return b
}

func stringArrayToSet(items []string) *schema.Set {
	set := schema.NewSet(schema.HashString, []interface{}{})
	for _, v := range items {
		set.Add(v)
	}
	return set
}

func setToStringArray(items *schema.Set) []string {
	flatten := make([]string, items.Len())

	for i, v := range items.List() {
		flatten[i] = v.(string)
	}
	return flatten
}

func interfaceToString(items []interface{}) []string {
	flatten := make([]string, len(items))

	for i, v := range items {
		flatten[i] = fmt.Sprint(v)
	}
	return flatten
}
