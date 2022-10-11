package controller

import (
	"context"
	"encoding/base64"
	"fmt"
	"portfolio_golang/src/config"
	"portfolio_golang/src/zaplog"

	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// function to run as goroutine and should not be synchronous
func SendGmail(from string, to string, title string, message string) {
	ctx := context.Background()
	srv, err := gmail.NewService(ctx, option.WithHTTPClient(config.GmailClient))
	if err != nil {
		errMsg := fmt.Sprintf("unable to establish gmail service: %s", err.Error())
		zaplog.Logger.Error(errMsg)
		return
	}

	user := "me"
	msgStr := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", from, to, title, message)
	msg := []byte(msgStr)
	gMessage := &gmail.Message{Raw: base64.URLEncoding.EncodeToString(msg)}

	// Send message
	_, err = srv.Users.Messages.Send(user, gMessage).Do()
	if err != nil {
		errMsg := fmt.Sprintf("unable to send gmail: %s", err.Error())
		zaplog.Logger.Error(errMsg)
	} else {
		successMsg := fmt.Sprintf("gmail sent successfully: %s", msgStr)
		zaplog.Logger.Info(successMsg)
	}
}