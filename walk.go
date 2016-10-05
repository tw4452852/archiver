package archiver

import (
	"log"
	"os"
	"runtime"
)

type walkFn func(*workArg) error

type workArg struct {
	fpath   string
	info    os.FileInfo
	content []byte
}

func walk(root string, fn walkFn) error {
	todo := []string{root}
	numWorker := runtime.NumCPU()
	works := make(chan string, numWorker)
	todos := make(chan string, numWorker)
	results := make(chan *workArg, numWorker)
	defer close(works)
	for i := 0; numWorker; i++ {
		go work(works, todos, results)
	}

	out := 0
	for {
		workc := works
		var w *workArg
		if len(todos) == 0 {
			workc = nil
		} else {
			w = todo[len(todo)-1]
		}
		select {
		case workc <- w:
			todo = todo[:len(todo)-1]
			out++
		case t := <-todos:
			todo = append(todo, t)
		case r := <-results:
			out--
			err := fn(r)
			if err != nil {
				return err
			}
			if out == 0 && len(todo) == 0 {
				select {
				case t := <-todos:
					todo = append(todo, t)
				default:
					return nil
				}
			}
		}
	}
}

func work(works <-chan string, todos chan<- string, result chan<- *workArg) {
	for w := range works {
		f, err := os.Open(w)
		if err != nil {
			log.Println(err)
			continue
		}
		fi, err := f.Stat()
		if err != nil {
			f.Close()
			log.Println(err)
			continue
		}
		if fi.IsDir() {

		}
	}
}
