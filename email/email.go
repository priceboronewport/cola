package email

import (
	"../tempfile"
	"bytes"
	"errors"
	"fmt"
	"mime/quotedprintable"
	"os"
	"os/exec"
	"strings"
)

type Email struct {
	Headers map[string]string
	Content string
}

func listContains(list []string, target string) bool {
	for _, v := range list {
		if strings.ToLower(v) == strings.ToLower(target) {
			return true
		}
	}
	return false
}

func New() *Email {
	em := Email{Headers: make(map[string]string), Content: ""}
	em.Headers["MIME-Version"] = "1.0"
	em.Headers["Content-Transfer-Encoding"] = "quoted-printable"
	em.Headers["Content-Disposition"] = "inline"
	return &em
}

func QuotedPrintable(s string) (string, error) {
	var ac bytes.Buffer
	w := quotedprintable.NewWriter(&ac)
	_, err := w.Write([]byte(s))
	if err != nil {
		return "", err
	}
	err = w.Close()
	if err != nil {
		return "", err
	}
	return ac.String(), nil
}

func (em *Email) SaveToFile(filename string) (err error) {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(em.AsString())
	if err == nil {
		f.Sync()
	}
	return
}

func (em *Email) Send(address_list string) (err error) {
	tf := tempfile.New("email")
	defer tf.Close()
	em.SaveToFile(tf.Filename)
	if address_list == "" {
		address_list = em.Addresses()
	}
	if address_list != "" {
		command := fmt.Sprintf("cat \"%s\"|/usr/sbin/sendmail -O DeliveryMode=b \"%s\"", tf.Filename, address_list)
		cmd := exec.Command("sh", "-c", command)
		_, err = cmd.CombinedOutput()
		fmt.Printf("eamil.Send('%s')\n", address_list)
	} else {
		err = errors.New("No Addresses")
	}
	return
}

func (em *Email) SendWhitelist(white_list []string) (err error) {
	address_list := strings.Split(em.Addresses(), ",")
	var send_list string
	for _, v := range address_list {
		if listContains(white_list, v) {
			if send_list != "" {
				send_list += ","
			}
			send_list += v
		} else {
			fmt.Printf("email.SendWhitelist(%v): Excluded %s\n", white_list, v)
		}
	}
	if send_list != "" {
		err = em.Send(send_list)
	}
	return
}

func (em *Email) AsString() (str string) {
	for k, v := range em.Headers {
		str += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	str += "\r\n"
	str += em.Content
	return
}

func HeaderAddress(name string, email_address string) string {
	return fmt.Sprintf("\"%s\" <%s>", name, email_address)
}

func (em *Email) Addresses() (addresses string) {
	tos := strings.Split(em.Headers["To"], ";")
	ccs := strings.Split(em.Headers["Cc"], ";")
	list := append(tos, ccs...)
	for _, v := range list {
		start := strings.Index(v, "<")
		if start != -1 {
			end := strings.Index(v, ">")
			if end != -1 {
				if len(addresses) > 0 {
					addresses += ","
				}
				addresses += v[start+1 : end]
			}
		}
	}
	return
}
