//go:build windows

package riskmisc

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
	"regexp"
	"sync"
)

// hasApp checks Windows Registry to find software path.
func hasApp(appName string) []string {
	queryKey := func(w *sync.WaitGroup, startKey registry.Key, res *[]string) {
		defer w.Done()
		queryPath := "Software\\Microsoft\\Windows\\CurrentVersion\\App Paths\\"
		k, err := registry.OpenKey(startKey, queryPath, registry.READ)
		if err != nil {
			return
		}
		// 读取所有子项
		keyNames, err := k.ReadSubKeyNames(0)
		if err != nil {
			return
		}
		for _, v := range keyNames {
			matched, err := regexp.MatchString(appName, v)
			if err != nil {
				fmt.Println("regexp error:", err)
			} else {
				if matched {
					tmpRegPath := queryPath + "\\" + v
					appKey, _ := registry.OpenKey(startKey, tmpRegPath, registry.READ)
					s, _, err := appKey.GetStringValue("")
					if err != nil {
						fmt.Println(err)
					} else {
						*res = append(*res, s)
					}
				}
			}
		}
	}
	res := []string{}

	waitGroup := new(sync.WaitGroup)
	waitGroup.Add(2)

	go queryKey(waitGroup, registry.LOCAL_MACHINE, &res)
	go queryKey(waitGroup, registry.CURRENT_USER, &res)
	waitGroup.Wait()

	return res
}
