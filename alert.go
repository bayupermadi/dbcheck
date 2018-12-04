package dbcheck

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
	"github.com/wjaoss/aws-wrapper/apps"
)

func AlertCheck(value int, threshold int, monitoringType string) {
	sesService := viper.GetBool("app.aws.service.ses.enabled")
	if value > threshold && sesService {
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
	}
}
