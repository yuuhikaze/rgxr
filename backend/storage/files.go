package storage

import (
	"fmt"
	"os"
)

func SaveSVG(uuid, svg string) error {
	path := fmt.Sprintf("/data/images/%s.svg", uuid)
	return os.WriteFile(path, []byte(svg), 0644)
}
