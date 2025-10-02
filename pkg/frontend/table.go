package frontend

import (
	"fmt"
	"slices"
	"strings"

	"github.com/aknea001/goDoList/pkg"
)

func line(lgstT int, lgstD int) {
	fmt.Printf("+%s+%s+\n",
		strings.Repeat("-", lgstT+2),
		strings.Repeat("-", lgstD+2),
	)
}

func DrawTable(tasks []pkg.Task) {
	titleLens := make([]int, 0)
	descLens := make([]int, 0)

	for i := range tasks {
		titleLens = append(titleLens, len(tasks[i].Title))
		descLens = append(descLens, len(tasks[i].Description))
	}

	// add the len of "title" and "description"
	titleLens = append(titleLens, 5)
	descLens = append(descLens, 11)

	longestTitle := slices.Max(titleLens)
	longestDesc := slices.Max(descLens)

	line(longestTitle, longestDesc)

	fmt.Printf("| Title%s | Description%s |\n",
		strings.Repeat(" ", longestTitle-5),
		strings.Repeat(" ", longestDesc-11),
	)

	line(longestTitle, longestDesc)

	for i := range tasks {
		currentTitle := tasks[i].Title
		currentDesc := tasks[i].Description

		fmt.Printf("| %s%s | %s%s |\n",
			currentTitle, strings.Repeat(" ", longestTitle-len(currentTitle)),
			currentDesc, strings.Repeat(" ", longestDesc-len(currentDesc)),
		)
	}

	line(longestTitle, longestDesc)
}
