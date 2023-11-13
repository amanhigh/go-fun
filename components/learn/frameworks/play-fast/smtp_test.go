package play_fast_test

import (
	"encoding/base64"
	"fmt"
	"net/smtp"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Smtp", func() {

	//sendemail -f from@gmail.com -t to@gmail.com -u "Bash Subject" -s smtp.mailtrap.io:2525 -m "I am Body" -v -o message-charset=$CHARSET -o username=bc705c85d0f7dc -o password=4dd5c28282e88a
	//https://mailtrap.io/
	var (
		//HACK:Move to Github Secrets and Rotate
		auth    = smtp.CRAMMD5Auth("bc705c85d0f7dc", "4dd5c28282e88a")
		server  = "smtp.mailtrap.io:2525"
		from    = "from@gmail.com"
		to      = "to@gmail.com"
		subject = "I am Subject"
		body    = "I am Body"
	)

	It("should mail", func() {
		err := smtp.SendMail(server, auth, from, []string{to}, composeMimeMail(to, from, subject, body))
		Expect(err).To(BeNil())
	})

})

func composeMimeMail(to string, from string, subject string, body string) []byte {
	header := make(map[string]string)
	header["From"] = from
	header["To"] = to
	header["Subject"] = subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/plain; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	return []byte(message)
}
