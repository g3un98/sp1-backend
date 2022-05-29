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

func netflixLogin(c context.Context, a account) (string, error) {
	var url, msg string

	if err := chromedp.Run(
		c,
		chromedp.Navigate(`https://www.netflix.com/kr/login`),
		chromedp.Click(`input[data-uia="login-field"]`, chromedp.NodeVisible),
		chromedp.SendKeys(`input[data-uia="login-field"]`, a.Id),
		chromedp.Click(`input[data-uia="password-field"]`, chromedp.NodeVisible),
		chromedp.SendKeys(`input[data-uia="password-field"]`, a.Pw),
		chromedp.Click(`button[data-uia="login-submit-button"]`, chromedp.NodeVisible),
		chromedp.Sleep(1*time.Second),
		chromedp.Location(&url),
	); err != nil {
		return "", err
	}

	if url == "https://www.netflix.com/kr/login" {
		if err := chromedp.Run(
			c,
			chromedp.Text(`div[data-uia="error-message-container"]`, &msg, chromedp.NodeVisible),
		); err != nil {
			return "", err
		}
		return msg, nil
	}

	return "", nil
}

func netflixLogout(c context.Context) error {
	return chromedp.Run(
		c,
		chromedp.Navigate(`https://www.netflix.com/kr/signout`),
	)
}

func netflixInfo(c *fiber.Ctx) error {
	ctx, cancel := chromedp.NewContext(
		context.Background(),
	)
	defer cancel()

	var account account
	if err := c.BodyParser(&account); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if len(account.Id) < 5 || len(account.Id) > 50 || len(account.Pw) < 4 || len(account.Pw) > 60 {
		return fiber.ErrBadRequest
	}

	msg, err := netflixLogin(ctx, account)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if msg != "" {
		return fiber.NewError(fiber.StatusUnauthorized, msg)
	}
	defer netflixLogout(ctx)

	var rawPayment, rawMembership string
	if err := chromedp.Run(
		ctx,
		chromedp.Navigate(`https://www.netflix.com/kr/youraccount`),
		chromedp.Text(`div[class="account-section-group payment-details -wide"]`, &rawPayment, chromedp.NodeVisible),
		chromedp.Text(`div[data-uia="plan-section"] > section`, &rawMembership, chromedp.NodeVisible),
	); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var (
		dummy            string
		year, month, day int
	)
	if rawPayment == "결제 정보가 없습니다" {
		account.Payment = payment{}
	} else {
		payments := strings.Split(rawPayment, "\n")
		if _, err := fmt.Sscanf(payments[2], "%s %s %d%s %d%s %d%s", &dummy, &dummy, &year, &dummy, &month, &dummy, &day, &dummy); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		account.Payment = payment{
			Type:   payments[0],
			Detail: payments[1],
			Next:   time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.FixedZone("KST", 9*60*60)).Unix(),
		}
	}

	switch strings.Split(rawMembership, "\n")[0] {
	case "스트리밍 멤버십에 가입하지 않으셨습니다.":
		account.Membership.Type = MEMBERSHIP_TYPE_NO
		account.Membership.Cost = MEMBERSHIP_COST_NO
	case "베이식":
		account.Membership.Type = MEMBERSHIP_NETLIFX_TYPE_BASIC
		account.Membership.Cost = MEMBERSHIP_NETLIFX_COST_BASIC
	case "스탠다드":
		account.Membership.Type = MEMBERSHIP_NETLIFX_TYPE_STANDARD
		account.Membership.Cost = MEMBERSHIP_NETLIFX_COST_STANDARD
	case "프리미엄":
		account.Membership.Type = MEMBERSHIP_NETLIFX_TYPE_PREMIUM
		account.Membership.Cost = MEMBERSHIP_NETLIFX_COST_PREMIUM
	default:
		return fiber.NewError(fiber.StatusInternalServerError, "Parse membership type")
	}

	body, err := sonic.Marshal(&account)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Send(body)
}

func netflixUnsubscribe(c *fiber.Ctx) error {
	ctx, cancel := chromedp.NewContext(
		context.Background(),
	)
	defer cancel()

	var account account
	if err := c.BodyParser(&account); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if len(account.Id) < 5 || len(account.Id) > 50 || len(account.Pw) < 4 || len(account.Pw) > 60 {
		return fiber.ErrBadRequest
	}

	msg, err := netflixLogin(ctx, account)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if msg != "" {
		return fiber.NewError(fiber.StatusUnauthorized, msg)
	}
	defer netflixLogout(ctx)

	var rawPayment, rawMembership string
	if err := chromedp.Run(
		ctx,
		chromedp.Navigate(`https://www.netflix.com/kr/youraccount`),
		chromedp.Text(`div[class="account-section-group payment-details -wide"]`, &rawPayment, chromedp.NodeVisible),
		chromedp.Text(`div[data-uia="plan-section"] > section`, &rawMembership, chromedp.NodeVisible),
	); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var (
		dummy            string
		year, month, day int
	)
	if rawPayment == "결제 정보가 없습니다" {
		account.Payment = payment{}
	} else {
		payments := strings.Split(rawPayment, "\n")
		if _, err := fmt.Sscanf(payments[2], "%s %s %d%s %d%s %d%s", &dummy, &dummy, &year, &dummy, &month, &dummy, &day, &dummy); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		account.Payment = payment{
			Type:   payments[0],
			Detail: payments[1],
			Next:   time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.FixedZone("KST", 9*60*60)).Unix(),
		}
	}

	switch strings.Split(rawMembership, "\n")[0] {
	case "스트리밍 멤버십에 가입하지 않으셨습니다.":
		account.Membership.Type = MEMBERSHIP_TYPE_NO
		account.Membership.Cost = MEMBERSHIP_COST_NO
	case "베이식":
		account.Membership.Type = MEMBERSHIP_NETLIFX_TYPE_BASIC
		account.Membership.Cost = MEMBERSHIP_NETLIFX_COST_BASIC
	case "스탠다드":
		account.Membership.Type = MEMBERSHIP_NETLIFX_TYPE_STANDARD
		account.Membership.Cost = MEMBERSHIP_NETLIFX_COST_STANDARD
	case "프리미엄":
		account.Membership.Type = MEMBERSHIP_NETLIFX_TYPE_PREMIUM
		account.Membership.Cost = MEMBERSHIP_NETLIFX_COST_PREMIUM
	default:
		return fiber.NewError(fiber.StatusInternalServerError, "Parse membership type")
	}

	body, err := sonic.Marshal(&account)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Send(body)
}
