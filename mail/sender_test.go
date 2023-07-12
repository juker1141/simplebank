package mail

import (
	"testing"

	"github.com/juker1141/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	config, err := util.LoadConfig("..")
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "A test email"
	content := `
	<h1>Hello World</h1>
	<p>This is a test message from <a href="http://google.com">Google</a></p>
	`

	to := []string{"ryu@realtime.tw"}
	attachFiles := []string{"../README.md"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	require.NoError(t, err)
}