package email_test

import (
	"github.com/StudioAquatan/hacku2020/pkg/email"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Classifier", func() {
	Context("email.ClassifyMail()", func() {
		When("the email has come", func() {
			It("returns true", func() {
				testCases := []string{
					"この度は弊社の面接にご参加頂き、誠にありがとうございました。誠に残念ながら、今回は貴意に添いかねる結果となりました。",
					"この度はご期待に添えない結果となりました。",
					"貴殿の今後のご活躍を心よりお祈り申し上げます。",
					"慎重に選考を進めた結果，",
				}
				for _, s := range testCases {
					res := email.ClassifyMail(s)
					Expect(res).To(Equal(true))
				}
			})
			It("returns false", func() {
				testCases := []string{
					"先日はお忙しい中、xxxの面接にご参加いただきまして誠にありがとうございました。",
				}
				for _, s := range testCases {
					res := email.ClassifyMail(s)
					Expect(res).To(Equal(false))
				}
			})
		})
	})
})
