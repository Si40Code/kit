package config

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
)

func startWatcher(path string) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println("watcher error:", err)
		return
	}

	go func() {
		for {
			select {
			case event := <-w.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					mu.Lock()
					if err := k.Load(file.Provider(path), yaml.Parser()); err == nil {
						LogConfigDiff("file", lastSnapshot, k.Raw())
						lastSnapshot = cloneMap(k.Raw())
						mu.Unlock()
						notifyChange()
					} else {
						log.Println("reload failed:", err)
						mu.Unlock()
					}
				}
			case err := <-w.Errors:
				log.Println("watch error:", err)
			}
		}
	}()

	_ = w.Add(path)
}
