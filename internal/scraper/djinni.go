package scraper

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/vmamchur/vacancy-board/internal/model"
	"github.com/vmamchur/vacancy-board/internal/repository"
)

type DjinniScraper struct {
	vacancyRepository repository.VacancyRepository
}

func (d DjinniScraper) Scrape() error {
	allocCtx, cancel := chromedp.NewRemoteAllocator(context.Background(), "http://chrome:9222/json/version")
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	err := chromedp.Run(ctx,
		chromedp.Navigate("https://djinni.co/login"),
		chromedp.WaitVisible("form#signup", chromedp.ByQuery),
		chromedp.SendKeys(`form#signup input[name="email"]`, "nerddcity@gmail.com", chromedp.ByQuery),
		chromedp.SendKeys(`form#signup input[name="password"]`, "11121314", chromedp.ByQuery),
		chromedp.Click(`form#signup button[type="submit"]`, chromedp.ByQuery),
		chromedp.Sleep(2*time.Second),
	)
	if err != nil {
		return err
	}

	page := 1
outer:
	for {
		log.Printf("Navigating to page %d...", page)

		var jobNodes []*cdp.Node
		err = chromedp.Run(ctx,
			chromedp.Navigate(fmt.Sprintf("https://djinni.co/jobs/?primary_keyword=fullstack&page=%d", page)),
			chromedp.WaitVisible("li[id^=job-item-]", chromedp.ByQuery),
			chromedp.Nodes("li[id^=job-item-]", &jobNodes, chromedp.ByQueryAll),
		)
		if err != nil {
			return err
		}

		log.Printf("Found %d vacancies on page %d", len(jobNodes), page)

		for _, node := range jobNodes {
			var title, url, companyName string

			err := chromedp.Run(ctx,
				chromedp.Text(".job-item__title-link", &title, chromedp.ByQuery, chromedp.FromNode(node)),
				chromedp.AttributeValue(".job-item__title-link", "href", &url, nil, chromedp.ByQuery, chromedp.FromNode(node)),
				chromedp.Text(`[data-analytics="company_page"]`, &companyName, chromedp.ByQuery),
			)
			if err != nil {
				log.Printf("Skipping vacancy due to extraction error: %v", err)
				continue
			}

			fullUrl := "https://djinni.co" + url
			log.Printf("Processing: \"%s\" at \"%s\" (%s)", title, companyName, fullUrl)

			_, err = d.vacancyRepository.Create(ctx, model.CreateVacancyDTO{
				Title: title,
				CompanyName: sql.NullString{
					String: companyName,
					Valid:  strings.TrimSpace(companyName) != "",
				},
				Url: fullUrl,
			})
			if err != nil {
				log.Printf("Stopping: vacancy already exists — \"%s\" (%s)", title, fullUrl)
				break outer
			}

			log.Printf("Saved: \"%s\" (%s)", title, fullUrl)
		}

		var isNextBtnVisible bool
		err = chromedp.Run(ctx,
			chromedp.Evaluate(`document.querySelector('li.page-item:not(.disabled) a.page-link span.bi-chevron-right') !== null`, &isNextBtnVisible),
		)
		if err != nil {
			return err
		}
		if !isNextBtnVisible {
			log.Println("No more pages. Done scraping.")
			break
		}

		page++
	}

	return nil
}
