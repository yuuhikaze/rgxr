package storage

import (
    "fmt"
    "os"
    "path/filepath"
)

func SaveSVG(uuid, svg string) error {
    path := fmt.Sprintf("/data/images/%s.svg", uuid)
    dir := filepath.Dir(path)
    
    // Ensure directory exists
    if err := os.MkdirAll(dir, 0755); err != nil {
        return err
    }
    
    return os.WriteFile(path, []byte(svg), 0644)
}

func SaveTeX(uuid, tex string) error {
    path := fmt.Sprintf("/data/tex/%s.tex", uuid)
    dir := filepath.Dir(path)
    
    // Ensure directory exists
    if err := os.MkdirAll(dir, 0755); err != nil {
        return err
    }
    
    return os.WriteFile(path, []byte(tex), 0644)
}

func GetSVG(uuid string) (string, error) {
    path := fmt.Sprintf("/data/images/%s.svg", uuid)
    data, err := os.ReadFile(path)
    return string(data), err
}

func GetTeX(uuid string) (string, error) {
    path := fmt.Sprintf("/data/tex/%s.tex", uuid)
    data, err := os.ReadFile(path)
    return string(data), err
}
