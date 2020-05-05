package handlers

import (
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"rickmaster2008/dm/models"
)

var wg sync.WaitGroup

// Head used to check file size and if url accepts range download
func Head(url string) (bool, int) {
	res, err := http.Head(url)
	if err != nil {
		log.Fatal(err)
	}
	r := res.Header["Accept-Ranges"][0] != "none"
	fsize, _ := strconv.Atoi(res.Header["Content-Length"][0])

	return r, fsize
}

// DownloadFile will download a url to a local file
func DownloadFile(filename string, url string, fsize int) error {

	limit := 8

	// Get the data
	client := &http.Client{}

	psize := fsize / limit
	diff := fsize % limit

	counter := models.NewWriteCounter(fsize)
	counter.Start()
	for i := 0; i < limit; i++ {
		wg.Add(1)
		s := psize * i
		e := psize * (i + 1)
		if i == limit-1 {
			e += diff
		}
		go DownloadPart(url, filename, s, e, i, client, &wg, counter)

	}

	wg.Wait()
	counter.Finish()
	f, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer f.Close()
	for i := 0; i < limit; i++ {
		tmp, err := os.Open(filename + ".tmp" + strconv.Itoa(i))
		if err != nil {
			return err
		}
		_, err = io.Copy(f, tmp)

		err = os.Remove(filename + ".tmp" + strconv.Itoa(i))
		if err != nil {
			return err
		}
	}

	return nil
}

// DownloadPart downlaods a range of the file
func DownloadPart(url string, filename string, s int, e int, i int, client *http.Client, wg *sync.WaitGroup, counter *models.WriteCounter) {
	defer wg.Done()
	//Create the file, but give it a tmp file extension, this means we won't overwrite a
	// file until it's downloaded, but we'll remove the tmp extension once downloaded
	out, err := os.Create(filename + ".tmp" + strconv.Itoa(i))
	if err != nil {
		panic(err)
	}
	defer out.Close()

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Range", "bytes="+strconv.Itoa(s)+"-"+strconv.Itoa(e-1))

	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	tr := io.TeeReader(res.Body, counter)
	_, err = io.Copy(out, tr)
	if err != nil {
		panic(err)
	}
}
