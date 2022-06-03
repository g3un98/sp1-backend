package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"github.com/bytedance/sonic"
	"github.com/chromedp/chromedp"
	"github.com/gofiber/fiber/v2"
)

func getNetflixAccount(id, pw string) (*account, error) {
	ctx, cancel := newChromedp()
	defer cancel()

	if err := netflixLogin(ctx, id, pw); err != nil {
		return nil, err
	}
	defer netflixLogout(ctx)

	var account account
	account.Id = id
	account.Pw = pw

	var rawPayment, rawMembership string
	if err := chromedp.Run(
		*ctx,
		chromedp.Navigate(`https://www.netflix.com/kr/youraccount`),
		chromedp.Text(`div[class="account-section-group payment-details -wide"]`, &rawPayment, chromedp.NodeVisible),
		chromedp.Text(`div[data-uia="plan-section"] > section`, &rawMembership, chromedp.NodeVisible),
	); err != nil {
		return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
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
			return nil, fiber.NewError(fiber.StatusInternalServerError, err.Error())
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
		return nil, fiber.NewError(fiber.StatusInternalServerError, strings.Split(rawMembership, "\n")[0])
	}

	return &account, nil
}

func netflixLogin(c *context.Context, id, pw string) error {
	if len(id) < 5 || len(id) > 50 || len(pw) < 4 || len(pw) > 60 {
		return fiber.ErrBadRequest
	}

	var url, msg string

	if err := chromedp.Run(
		*c,
		chromedp.Navigate(`https://www.netflix.com/kr/login`),
		chromedp.Click(`input[data-uia="login-field"]`, chromedp.NodeVisible),
		chromedp.SendKeys(`input[data-uia="login-field"]`, id),
		chromedp.Click(`input[data-uia="password-field"]`, chromedp.NodeVisible),
		chromedp.SendKeys(`input[data-uia="password-field"]`, pw),
		chromedp.Click(`button[data-uia="login-submit-button"]`, chromedp.NodeVisible),
		chromedp.Sleep(1*time.Second),
		chromedp.Location(&url),
	); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if url == "https://www.netflix.com/kr/login" {
		if err := chromedp.Run(
			*c,
			chromedp.Text(`div[data-uia="error-message-container"]`, &msg, chromedp.NodeVisible),
		); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return fiber.NewError(fiber.StatusUnauthorized, msg)
	}

	return nil
}

func netflixLogout(c *context.Context) error {
	return chromedp.Run(
		*c,
		chromedp.Navigate(`https://www.netflix.com/kr/signout`),
	)
}

func postNetflixAccount(c *fiber.Ctx) error {
	_, cancel := newChromedp()
	defer cancel()

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

	account, err := getNetflixAccount(parser.OttId, parser.OttPw)
	if err != nil {
		return err
	}

	bodyByte, err := sonic.Marshal(&account)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.Send(bodyByte)
}

func putNetflixAccount(c *fiber.Ctx) error {
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
    filter := bson.M{"ott": "Netflix", "account.id": parser.OttId, "account.pw": parser.OttPw}
    if err := getCollection(client, "group").FindOne(ctx, filter).Decode(&group); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
    }
    
    _, cancel = newChromedp()
	defer cancel()

	account, err := getNetflixAccount(parser.OttId, parser.OttPw)
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

func deleteNetflixMembership(c *fiber.Ctx) error {
	ctx, cancel := newChromedp()
	defer cancel()

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

	if err := netflixLogin(ctx, parser.OttId, parser.OttPw); err != nil {
		return err
	}
	defer netflixLogout(ctx)

	var url string
	if err := chromedp.Run(
		*ctx,
		chromedp.Navigate(`https://www.netflix.com/kr/cancel`),
		chromedp.Sleep(1*time.Second),
		chromedp.Location(&url),
	); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	switch url {
	case "https://www.netflix.com/kr/":
		return fiber.ErrBadRequest
	case "https://www.netflix.com/cancelplan/confirm":
		return c.SendStatus(fiber.StatusOK)
	case "https://www.netflix.com/CancelPlan?locale=ko-KR":
		break
	default:
		return fiber.NewError(fiber.StatusInternalServerError, url)
	}

	if err := chromedp.Run(
		*ctx,
		chromedp.Click(`button[data-uia="action-finish-cancellation"]`, chromedp.NodeVisible),
		chromedp.Sleep(1*time.Second),
		chromedp.Location(&url),
	); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if url == "https://www.netflix.com/cancelplan/confirm" {
		return c.SendStatus(fiber.StatusOK)
	}

	return fiber.NewError(fiber.StatusInternalServerError, url)
}
