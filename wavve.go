package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	"github.com/chromedp/chromedp"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

func getWavveAccount(id, pw string) (*account, error) {
	ctx, cancel := newChromedp()
	defer cancel()

	var account account
	account.Id = id
	account.Pw = pw

	if err := wavveLogin(ctx, id, pw); err != nil {
		return nil, err
	}
	defer wavveLogout(ctx)

	var contents string
	if err := chromedp.Run(
		*ctx,
		chromedp.Navigate(`https://www.wavve.com/my/subscription_ticket`),
		chromedp.Text(`#contents`, &contents, chromedp.NodeVisible),
	); err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if contents == "이용권 결제 내역이 없어요." {
		account.Payment = payment{}
		account.Membership = membership{MEMBERSHIP_TYPE_NO, MEMBERSHIP_COST_NO}

		return &account, nil
	}

	var rawPaymentType, rawPaymentNext, rawMembershipType, rawMembershipCost string
	if err := chromedp.Run(
		*ctx,
		chromedp.Text(`#contents > div.mypooq-inner-wrap > section > div > div > div > table > tbody > tr > td:nth-child(6) > span > span`, &rawPaymentType, chromedp.NodeVisible),
		chromedp.Text(`#contents > div.mypooq-inner-wrap > section > div > div > div > table > tbody > tr > td:nth-child(5)`, &rawPaymentNext, chromedp.NodeVisible),
		chromedp.Text(`#contents > div.mypooq-inner-wrap > section > div > div > div > table > tbody > tr > td:nth-child(2) > div > p.my-pay-tit > span:nth-child(3)`, &rawMembershipType, chromedp.NodeVisible),
		chromedp.Text(`#contents > div.mypooq-inner-wrap > section > div > div > div > table > tbody > tr > td:nth-child(4)`, &rawMembershipCost, chromedp.NodeVisible),
	); err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var (
		dummy            string
		year, month, day int
	)
	if _, err := fmt.Sscanf(strings.Split(rawPaymentNext, " ")[0], "%d-%d-%d", &year, &month, &day); err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	account.Payment = payment{
		Type: rawPaymentType,
		Next: time.Date(year, time.Month(month), day+1, 0, 0, 0, 0, time.FixedZone("KST", 9*60*60)).Unix(),
	}

    rawMembershipCost = strings.ReplaceAll(rawMembershipCost, ",", "")
	if _, err := fmt.Sscanf(rawMembershipCost, "%d%s", &account.Membership.Cost, &dummy); err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
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
	default:
		return nil, fiber.NewError(fiber.StatusInternalServerError, rawMembershipType)
	}

	return &account, nil
}

func wavveLogin(c *context.Context, id, pw string) error {
	if len(id) < 1 || len(pw) < 1 {
		return fiber.ErrBadRequest
	}

	var url, msg string
	if err := chromedp.Run(
		*c,
		chromedp.Navigate(`https://www.wavve.com/login`),
		chromedp.Click(`input[title="아이디"]`, chromedp.NodeVisible),
		chromedp.SendKeys(`input[title="아이디"]`, id, chromedp.NodeVisible),
		chromedp.Click(`input[title="비밀번호"]`, chromedp.NodeVisible),
		chromedp.SendKeys(`input[title="비밀번호"]`, pw, chromedp.NodeVisible),
		chromedp.Click(`a[title="로그인"]`, chromedp.NodeVisible),
		chromedp.Sleep(1*time.Second),
		chromedp.Location(&url),
	); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if url == "https://www.wavve.com/login" {
		if err := chromedp.Run(
			*c,
			chromedp.Text(`p[class="login-error-red"]`, &msg, chromedp.NodeVisible),
		); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return fiber.NewError(fiber.StatusBadRequest, msg)
	}

	return nil
}

func wavveLogout(c *context.Context) error {
	return chromedp.Run(
		*c,
		chromedp.Navigate(`https://www.wavve.com`),
		chromedp.Click(`#app > div.body > div:nth-child(2) > header > div:nth-child(1) > div.header-nav > div > ul > li.over-parent-1depth > div > ul > li:nth-child(4) > button`, chromedp.NodeVisible),
	)
}

func postWavveAccount(c *fiber.Ctx) error {
	var parser struct {
		OttId string `json:"ott_id"`
		OttPw string `json:"ott_pw"`
	}
	if err := c.BodyParser(&parser); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if parser.OttId == "" || parser.OttPw == "" {
		return fiber.ErrBadRequest
	}

	account, err := getWavveAccount(parser.OttId, parser.OttPw)
	if err != nil {
		return err
	}

	bodyByte, err := sonic.Marshal(&account)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Send(bodyByte)
}

func putWavveAccount(c *fiber.Ctx) error {
	client, ctx, cancel, err := newClient()
	if err != nil {
		return err
	}
	defer cancel()

	var parser struct {
		OttId string `json:"ott_id" bson:"ott_id"`
		OttPw string `json:"ott_pw" bson:"ott_pw"`
	}
	if err := c.BodyParser(&parser); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if parser.OttId == "" || parser.OttPw == "" {
		return fiber.ErrBadRequest
	}

	var group group
	filter := bson.M{"ott": "wavve", "account.id": parser.OttId, "account.pw": parser.OttPw}
	if err := getCollection(client, "group").FindOne(ctx, filter).Decode(&group); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	account, err := getWavveAccount(parser.OttId, parser.OttPw)
	if err != nil {
		return err
	}

	if group.Account != *account {
		newAccountBsonByte, err := bson.Marshal(account)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		var newAccountBson bson.M
		if err = bson.Unmarshal(newAccountBsonByte, &newAccountBson); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		update := bson.M{"$set": bson.M{"account": newAccountBson}}
		getCollection(client, "group").FindOneAndUpdate(ctx, filter, update)
	}

	accountByte, err := sonic.Marshal(account)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.Send(accountByte)
}
