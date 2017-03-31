package main

import (
	"flag"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func main() {
	name := flag.String("name", "", "the name to search")
	flag.Parse()

	linearStart := time.Now()
	log.Println("Start crawling.")

	var babies []Baby

	// http://www.babygalerie24.de/de_DE/fulda/?page=96
	rawurl := "http://www.babygalerie24.de/de_DE/fulda/?page=%d"
	start := 1
	end := 96

	for i := start; i < end; i++ {
		url := fmt.Sprintf(rawurl, i)
		src, err := downloadPage(url)
		if err != nil {
			log.Printf("Error while downloading %s: %v\n", url, err)
		}

		tmp := extractBabyKlinikum(src)

		for _, b := range tmp {
			if strings.Contains(b.Name, *name) {
				log.Printf("Found: %s\n", b.String())
			}
			// log.Printf("Scanned: %s\n", b.String())
		}

		babies = append(babies, tmp...)
	}
	linearElapsed := time.Since(linearStart)

	log.Println("Linear Took: ", linearElapsed.String())
}

func extractBabyKlinikum(src string) []Baby {
	var result []Baby
	rr := `<li class="baby_(?P<gender>(female|male|twins))"><a href="/de_DE/fulda/babygalerie/baby/(.*).html" title="Geboren am (?P<birthdate>[0-9.]{10})(| um (?P<when>[0-9:]{5}) Uhr , Gewicht: (?P<weight>[0-9]{3,5}) g, Größe: (?P<size>[0-9]{2}) cm)"><img
        src="/resources/gallery/thumb/(.*)/small/(.*)/(.*).jpg"
        data-midsrc="/resources/gallery/thumb/(.*)/mid/(.*)/(.*).jpg"
        data-fullsrc="/resources/gallery/thumb/(.*)/full/(.*)/(.*).jpg"
        alt="(?P<name>(.*))" width="201" height="141" /><span class="name">(?P<namea>(.*))</span></a></li>`

	r := regexp.MustCompile(rr)

	names := r.SubexpNames() // holt alle Namen raus
	matches := r.FindAllStringSubmatch(src, -1)

	for _, m := range matches {
		// range over all found babies on a page

		md := map[string]string{}
		for i, n := range m {
			md[names[i]] = n
		}
		var birthdate time.Time
		var size, weight int
		var err error

		// create baby
		birthdate, err = time.Parse("02.01.2006 15:04", md["birthdate"]+" "+md["when"])
		if err != nil {
			//log.Printf("could not parse datetime %s - %s: %v\n", md["birthdate"], md["when"], err)
			birthdate, err = time.Parse("02.01.2006", md["birthdate"])
			if err != nil {
				log.Println("still not possible to parse time -.-")
				birthdate = time.Time{}
			}
		}
		// only when is available
		if md["size"] != "" {
			size, err = strconv.Atoi(md["size"])
			if err != nil {
				log.Printf("could not parse size integer %s: %v\n", md["size"], err)
				size = 0
			}
		} else {
			size = 0
		}
		// only when is available
		if md["weight"] != "" {
			weight, err = strconv.Atoi(md["weight"])
			if err != nil {
				log.Printf("could not parse weight integer %s: %v\n", md["weight"], err)
				weight = 0
			}
		} else {
			weight = 0
		}
		b := Baby{
			ID:        0,
			Name:      md["name"],
			Type:      md["gender"],
			Birthdate: birthdate,
			Size:      size,
			Weight:    weight,
		}
		result = append(result, b)
	}

	return result
}

/*
CREATE TABLE `babys` (
	`id`	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
	`hospital`	INTEGER NOT NULL,
	`gender`	INTEGER NOT NULL,
	`birthdate`	INTEGER
);
*/
