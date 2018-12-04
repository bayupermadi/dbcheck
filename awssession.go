package dbcheck

import (
	"github.com/spf13/viper"
	"github.com/wjaoss/aws-wrapper/session"
)

func AwsSession(awsRegion string) {
	awsKeyID := viper.Get("app.aws.credential.id-key").(string)
	awsSecretKey := viper.Get("app.aws.credential.secret-key").(string)

	session.SetConfiguration(awsKeyID, awsSecretKey, awsRegion)

}
