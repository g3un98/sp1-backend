package main

import (
	"context"
	"time"

	"github.com/chromedp/chromedp"
)

func containAdminMembers(s []member, v string) bool {
	for _, vv := range s {
		if v == vv.AppId && vv.IsAdmin == 1 {
			return true
		}
	}
	return false
}

func containMembers(s []member, v string) bool {
	for _, vv := range s {
		if v == vv.AppId {
			return true
		}
	}
	return false
}

func newChromedp() (*context.Context, func()) {
	ctx, cancel1 := chromedp.NewContext(
		context.Background(),
	)
	ctx, cancel2 := context.WithTimeout(ctx, 1*time.Minute)

	return &ctx, func() {
		defer cancel1()
		defer cancel2()
	}
}
