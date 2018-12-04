package dbcheck

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/spf13/viper"
	"github.com/wjaoss/aws-wrapper/apps"
)

func AlertCheck(value int, threshold int, monitoringType string) {
	sesService := viper.GetBool("app.aws.service.ses.enabled")
	sesRegion := viper.Get("app.aws.service.ses.region").(string)
	if sesRegion != "" {
		AwsSession(sesRegion)
	}
	contents, _ := ioutil.ReadFile(monitoringType)
	if value > threshold && sesService && string(contents) != "alert" {
		msg := "alert"
		sender := viper.Get("app.aws.service.ses.from").(string)
		recipient := viper.Get("app.aws.service.ses.to").(string)
		subject := viper.Get("app.aws.service.ses.subject").(string)
		body := fmt.Sprintf("Current Value of "+monitoringType+" is %d", value)

		dest := strings.Split(recipient, ", ")
		start := 0
		for i := 0; i < len(dest); i++ {
			start += i
			apps.SES(sender, dest[start], subject, body)
		}

		ioutil.WriteFile(monitoringType, []byte(msg), 0644)
	}
	if value <= threshold && string(contents) != "noalert" {
		msg := "noalert"
		ioutil.WriteFile(monitoringType, []byte(msg), 0644)
	}
}
