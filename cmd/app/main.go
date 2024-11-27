package main

import (
	"fmt"
	"password-manager/pkg/policy"
)

func main() {
	//console.Init()

	var policies = policy.Policies{}

	err := policies.Load()

	newPolicy := policy.NewPolicy(
		"Example",
		[]rune{'a', 'b', 'c'},
		0,
		0,
		0,
		0,
		0,
		0,
		0,
		0,
	)

	policies.New(newPolicy)

	changeTo := make(policy.Fvm)
	//
	changeTo["MinimumNumbers"] = 10
	//
	err = policies.UpdateByName("Example", changeTo)

	policies.Save()

	fmt.Println(err)

}
