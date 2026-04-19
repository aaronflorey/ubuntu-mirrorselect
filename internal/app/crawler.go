package app

import (
	"strings"

	"github.com/haukened/mirrorselect/internal/llog"

	"github.com/biter777/countries"
	"github.com/gocolly/colly"
)

var launchpadCountryCodeAliases = map[string]string{
	"IRAN, ISLAMIC REPUBLIC OF":    "IR",
	"KOREA, REPUBLIC OF":           "KR",
	"MACEDONIA, REPUBLIC OF":       "MK",
	"MOLDOVA, REPUBLIC OF":         "MD",
	"TANZANIA, UNITED REPUBLIC OF": "TZ",
	"VIET NAM":                     "VN",
}

func crawlLaunchpad(desiredCC string) (mirrors []Mirror, err error) {
	desiredCC = strings.ToUpper(strings.TrimSpace(desiredCC))
	currentCC := ""

	c := colly.NewCollector(
		colly.AllowedDomains("launchpad.net"), // only visit launchpad.net
		colly.MaxDepth(1),                     // only scrape the first page, dont recurse
	)
	c.OnError(func(r *colly.Response, e error) {
		err = e
	})
	c.OnHTML("table#mirrors_list > tbody", func(h *colly.HTMLElement) {
		h.ForEach("tr", func(_ int, row *colly.HTMLElement) {
			row.ForEach("*", func(_ int, cell *colly.HTMLElement) {
				if cell.Attr("colspan") == "2" {
					cName := strings.TrimSpace(cell.Text)
					switch cName {
					case "":
						return
					case "Total":
						return
					default:
						countryCode, ok := parseCountryCode(cName)
						if !ok {
							currentCC = ""
							llog.Warnf("Unable to map launchpad country heading %q to ISO-3166 country code", cName)
							return
						}
						currentCC = countryCode
						llog.Debugf("Updated country to %s", currentCC)
					}
				} else if cell.Attr("href") != "" {
					link := cell.Attr("href")
					if strings.HasPrefix(link, "http") && currentCC == desiredCC {
						mirror, ok := NewMirror(link)
						if ok {
							mirrors = append(mirrors, mirror)
						}
					}
				}
			})
		})
	})
	c.OnScraped(func(r *colly.Response) {
		llog.Debug("Finished scraping launchpad.net")
	})
	err = c.Visit("https://launchpad.net/ubuntu/+archivemirrors")
	return
}

func parseCountryCode(country string) (string, bool) {
	name := strings.TrimSpace(country)
	if name == "" {
		return "", false
	}

	code := countries.ByName(name)
	if code != countries.Unknown {
		return code.Alpha2(), true
	}

	parts := strings.SplitN(name, ",", 2)
	baseName := strings.TrimSpace(parts[0])
	if baseName != "" {
		code = countries.ByName(baseName)
		if code != countries.Unknown {
			return code.Alpha2(), true
		}
	}

	aliasCode, ok := launchpadCountryCodeAliases[strings.ToUpper(name)]
	if ok {
		return aliasCode, true
	}

	return "", false
}
