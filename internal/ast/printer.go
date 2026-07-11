package ast

import (
	"encoding/json"
	"fmt"
)

func Print(node any) {
	data, err := json.MarshalIndent(node, "", "  ")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(data))
}
