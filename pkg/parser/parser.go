package parser

import (
	"regexp"
	"strconv"
	"strings"
)

var blockSizes = map[string]int{
	"B":  1,
	"KB": 1024,
	"MB": 1024 * 1024,
	"GB": 1024 * 1024 * 1024,
}

// ParseSize converts message string to bytes
func ParseSize(msgSizeStr string) (int, error) {
	msgSizeStr = strings.TrimSpace(strings.ToUpper(msgSizeStr))

	re := regexp.MustCompile(`([0-9]+)(\w+)`)
	res := re.FindAllStringSubmatch(msgSizeStr, -1)

	for k, v := range blockSizes {
		if !strings.HasSuffix(msgSizeStr, k) {
			continue
		}

		if res[0][2] == k {
			size, err := strconv.Atoi(res[0][1])
			if err != nil {
				return 0, err
			}
			return size * v, nil
		}
	}

	return 0, nil
}
