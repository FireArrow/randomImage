package main

import (
	"fmt"
)

func rootPage(basePath string) []byte {
	return []byte(fmt.Sprintf(rootTemplate))
}

const rootTemplate = `
`
