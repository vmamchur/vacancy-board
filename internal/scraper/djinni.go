package scraper

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/vmamchur/vacancy-board/internal/model"
)

type DjinniScraper struct{}

func (d DjinniScraper) Scrape() ([]model.Vacancy, error) {
	var vacancies []model.Vacancy

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background())
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
		return nil, err
	}

	page := 1
	for {
		var jobNodes []*cdp.Node
		err = chromedp.Run(ctx,
			chromedp.Navigate(fmt.Sprintf("https://djinni.co/jobs/?primary_keyword=fullstack&page=%d", page)),
			chromedp.WaitVisible("li[id^=job-item-]", chromedp.ByQuery),
			chromedp.Nodes("li[id^=job-item-]", &jobNodes, chromedp.ByQueryAll),
		)
		if err != nil {
			return nil, err
		}

		for _, node := range jobNodes {
			var title, url, companyName string

			err := chromedp.Run(ctx,
				chromedp.Text(".job-item__title-link", &title, chromedp.ByQuery, chromedp.FromNode(node)),
				chromedp.AttributeValue(".job-item__title-link", "href", &url, nil, chromedp.ByQuery, chromedp.FromNode(node)),
				chromedp.Text(`[data-analytics="company_page"]`, &companyName, chromedp.ByQuery),
			)
			if err != nil {
				continue
			}

			vacancies = append(vacancies, model.Vacancy{
				Title: title,
				CompanyName: sql.NullString{
					String: companyName,
					Valid:  strings.TrimSpace(companyName) != "",
				},
				Url: "https://djinni.co" + url,
			})
		}

		var isNextBtnVisible bool
		err = chromedp.Run(ctx,
			chromedp.Evaluate(`document.querySelector('li.page-item:not(.disabled) a.page-link span.bi-chevron-right') !== null`, &isNextBtnVisible),
		)
		if err != nil {
			return nil, err
		}
		if !isNextBtnVisible {
			break
		}

		page++
	}

	return vacancies, nil
}
