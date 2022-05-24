package twitter

import (
	"fmt"
	"log"

	"github.com/mrjones/oauth"
)

func NewDesktopClient(consumerKey, consumerSecret string) *DesktopClient {
	newDesktop := new(DesktopClient)
	newDesktop.OAuthConsumer = oauth.NewConsumer(
		consumerKey,
		consumerSecret,
		oauth.ServiceProvider{
			RequestTokenUrl:   OAUTH_REQUES_TOKEN,
			AuthorizeTokenUrl: OAUTH_AUTH_TOKEN,
			AccessTokenUrl:    OAUTH_ACCESS_TOKEN,
		},
	)
	//Enable debug info
	newDesktop.OAuthConsumer.Debug(false)

	return newDesktop
}

type DesktopClient struct {
	Client
	OAuthConsumer *oauth.Consumer
}

func (d *DesktopClient) DoAuth() error {
	requestToken, u, err := d.OAuthConsumer.GetRequestTokenAndUrl("oob")
	fmt.Println("rest token=", requestToken, " err=", err)
	if err != nil {
		log.Println(err)
		return err
	}

	fmt.Println("(1) Go to: " + u)
	fmt.Println("(2) Grant access, you should get back a verification code.")
	fmt.Println("(3) Enter that verification code here: ")

	verificationCode := ""
	fmt.Scanln(&verificationCode)

	accessToken, err := d.OAuthConsumer.AuthorizeToken(requestToken, verificationCode)
	if err != nil {
		log.Println(err)
		return err
	}

	d.HttpConn, err = d.OAuthConsumer.MakeHttpClient(accessToken)
	if err != nil {
		log.Println(err)
		return err
	}

	return err
}
