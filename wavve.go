package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"github.com/chromedp/chromedp"
	"github.com/gofiber/fiber/v2"
)

func getWavveAccount(id, pw string) (*account, error) {
	ctx, cancel := chromedp.NewContext(
		context.Background(),
	)
	defer cancel()
	ctx, cancel = context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	var account account
    account.Id = id
    account.Pw = pw

	if len(account.Id) < 1 || len(account.Pw) < 1 {
		return nil, fiber.ErrBadRequest
	}

	msg, err := wavveLogin(&ctx, account)
    checkError(err)
	if msg != "" {
		return nil, fiber.NewError(fiber.StatusUnauthorized, msg)
	}
	defer wavveLogout(&ctx)

	var contents string
	err = chromedp.Run(
		ctx,
		chromedp.Navigate(`https://www.wavve.com/my/subscription_ticket`),
		chromedp.Text(`#contents`, &contents, chromedp.NodeVisible),
	)
    checkError(err)

	if contents == "이용권 결제 내역이 없어요." {
		account.Payment = payment{}
		account.Membership = membership{MEMBERSHIP_TYPE_NO, MEMBERSHIP_COST_NO}

		return &account, nil
	}

	var rawPaymentType, rawPaymentNext, rawMembershipType, rawMembershipCost string
	err = chromedp.Run(
		ctx,
		chromedp.Text(`#contents > div.mypooq-inner-wrap > section > div > div > div > table > tbody > tr > td:nth-child(6) > span > span`, &rawPaymentType, chromedp.NodeVisible),
		chromedp.Text(`#contents > div.mypooq-inner-wrap > section > div > div > div > table > tbody > tr > td:nth-child(5)`, &rawPaymentNext, chromedp.NodeVisible),
		chromedp.Text(`#contents > div.mypooq-inner-wrap > section > div > div > div > table > tbody > tr > td:nth-child(2) > div > p.my-pay-tit > span:nth-child(3)`, &rawMembershipType, chromedp.NodeVisible),
		chromedp.Text(`#contents > div.mypooq-inner-wrap > section > div > div > div > table > tbody > tr > td:nth-child(4)`, &rawMembershipCost, chromedp.NodeVisible),
	)
    checkError(err)

	var year, month, day int
	fmt.Sscanf(strings.Split(rawPaymentNext, " ")[0], "%d-%d-%d", &year, &month, &day)
	account.Payment = payment{
		Type: rawPaymentType,
		Next: time.Date(year, time.Month(month), day+1, 0, 0, 0, 0, time.FixedZone("KST", 9*60*60)).Unix(),
	}

	var dummy string
	_, err = fmt.Sscanf(rawMembershipCost, "%d%s", &account.Membership.Cost, &dummy)
    checkError(err)

	switch rawMembershipType {
	case "Basic":
		account.Membership.Type = MEMBERSHIP_WAVVE_TYPE_BASIC
	case "Standard":
		account.Membership.Type = MEMBERSHIP_WAVVE_TYPE_STANDARD
	case "Premium":
		account.Membership.Type = MEMBERSHIP_WAVVE_TYPE_PREMIUM
	case "Basic X FLO 무제한":
		account.Membership.Type = MEMBERSHIP_WAVVE_TYPE_FLO
	case "Basic X Bugs 듣기":
		account.Membership.Type = MEMBERSHIP_WAVVE_TYPE_BUGS
	case "Basic X KB 나라사랑카드":
		account.Membership.Type = MEMBERSHIP_WAVVE_TYPE_KB
	}

	return &account, nil
}

func wavveLogin(c *context.Context, a account) (string, error) {
	var url, msg string
	if err := chromedp.Run(
		*c,
		chromedp.Navigate(`https://www.wavve.com/login`),
		chromedp.Click(`input[title="아이디"]`, chromedp.NodeVisible),
		chromedp.SendKeys(`input[title="아이디"]`, a.Id, chromedp.NodeVisible),
		chromedp.Click(`input[title="비밀번호"]`, chromedp.NodeVisible),
		chromedp.SendKeys(`input[title="비밀번호"]`, a.Pw, chromedp.NodeVisible),
		chromedp.Click(`a[title="로그인"]`, chromedp.NodeVisible),
		chromedp.Sleep(1*time.Second),
		chromedp.Location(&url),
	); err != nil {
		return "", err
	}

	if url == "https://www.wavve.com/login" {
		if err := chromedp.Run(
			*c,
			chromedp.Text(`p[class="login-error-red"]`, &msg, chromedp.NodeVisible),
		); err != nil {
			return "", err
		}
		return msg, nil
	}

	return "", nil
}

func wavveLogout(c *context.Context) error {
	return chromedp.Run(
		*c,
		chromedp.Navigate(`https://www.wavve.com`),
		chromedp.Click(`#app > div.body > div:nth-child(2) > header > div:nth-child(1) > div.header-nav > div > ul > li.over-parent-1depth > div > ul > li:nth-child(4) > button`, chromedp.NodeVisible),
	)
}

func wavveInfo(c *fiber.Ctx) error {
	ctx, cancel := chromedp.NewContext(
		context.Background(),
	)
	defer cancel()
	ctx, cancel = context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	var account account
	if err := c.BodyParser(&account); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if len(account.Id) < 1 || len(account.Pw) < 1 {
		return fiber.ErrBadRequest
	}

	msg, err := wavveLogin(&ctx, account)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if msg != "" {
		return fiber.NewError(fiber.StatusUnauthorized, msg)
	}
	defer wavveLogout(&ctx)

	var contents string
	if err := chromedp.Run(
		ctx,
		chromedp.Navigate(`https://www.wavve.com/my/subscription_ticket`),
		chromedp.Text(`#contents`, &contents, chromedp.NodeVisible),
	); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if contents == "이용권 결제 내역이 없어요." {
		account.Payment = payment{}
		account.Membership = membership{MEMBERSHIP_TYPE_NO, MEMBERSHIP_COST_NO}

		body, err := sonic.Marshal(&account)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		return c.Send(body)
	}

	var rawPaymentType, rawPaymentNext, rawMembershipType, rawMembershipCost string
	if err := chromedp.Run(
		ctx,
		chromedp.Text(`#contents > div.mypooq-inner-wrap > section > div > div > div > table > tbody > tr > td:nth-child(6) > span > span`, &rawPaymentType, chromedp.NodeVisible),
		chromedp.Text(`#contents > div.mypooq-inner-wrap > section > div > div > div > table > tbody > tr > td:nth-child(5)`, &rawPaymentNext, chromedp.NodeVisible),
		chromedp.Text(`#contents > div.mypooq-inner-wrap > section > div > div > div > table > tbody > tr > td:nth-child(2) > div > p.my-pay-tit > span:nth-child(3)`, &rawMembershipType, chromedp.NodeVisible),
		chromedp.Text(`#contents > div.mypooq-inner-wrap > section > div > div > div > table > tbody > tr > td:nth-child(4)`, &rawMembershipCost, chromedp.NodeVisible),
	); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var year, month, day int
	fmt.Sscanf(strings.Split(rawPaymentNext, " ")[0], "%d-%d-%d", &year, &month, &day)
	account.Payment = payment{
		Type: rawPaymentType,
		Next: time.Date(year, time.Month(month), day+1, 0, 0, 0, 0, time.FixedZone("KST", 9*60*60)).Unix(),
	}

	var dummy string
	if _, err = fmt.Sscanf(rawMembershipCost, "%d%s", &account.Membership.Cost, &dummy); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	switch rawMembershipType {
	case "Basic":
		account.Membership.Type = MEMBERSHIP_WAVVE_TYPE_BASIC
	case "Standard":
		account.Membership.Type = MEMBERSHIP_WAVVE_TYPE_STANDARD
	case "Premium":
		account.Membership.Type = MEMBERSHIP_WAVVE_TYPE_PREMIUM
	case "Basic X FLO 무제한":
		account.Membership.Type = MEMBERSHIP_WAVVE_TYPE_FLO
	case "Basic X Bugs 듣기":
		account.Membership.Type = MEMBERSHIP_WAVVE_TYPE_BUGS
	case "Basic X KB 나라사랑카드":
		account.Membership.Type = MEMBERSHIP_WAVVE_TYPE_KB
	}

	body, err := sonic.Marshal(&account)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Send(body)
}
