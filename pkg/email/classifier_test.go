package email_test

import (
	"github.com/StudioAquatan/hacku2020/pkg/email"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Classifier", func() {
	Context("email.ClassifyMail()", func() {
		When("the email subject is screening", func() {
			It("returns true", func() {
				testCases := []string{
					"選考結果",
				}
				for _, s := range testCases {
					res := email.ClassifyScreeningMailBySubj(s)
					Expect(res).To(Equal(true))
				}
			})
			It("returns false", func() {
				testCases := []string{
					"先日はお忙しい中、xxxの面接にご参加いただきまして誠にありがとうございました。",
				}
				for _, s := range testCases {
					res := email.ClassifyScreeningMailBySubj(s)
					Expect(res).To(Equal(false))
				}
			})
		})
		When("the email body is oinori", func() {
			It("returns true", func() {
				testCases := []string{
					"この度は弊社の面接にご参加頂き、誠にありがとうございました。誠に残念ながら、今回は貴意に添いかねる結果となりました。",
					"この度はご期待に添えない結果となりました。",
					"貴殿の今後のご活躍を心よりお祈り申し上げます。",
					"誠に申し訳ございませんか",
				}
				for _, s := range testCases {
					res := email.ClassifyOinoriMailByBody(s)
					Expect(res).To(Equal(true))
				}
			})
			It("returns false", func() {
				testCases := []string{
					"先日はお忙しい中、xxxの面接にご参加いただきまして誠にありがとうございました。",
				}
				for _, s := range testCases {
					res := email.ClassifyOinoriMailByBody(s)
					Expect(res).To(Equal(false))
				}
			})
		})
		When("the email body is acceptance", func() {
			It("returns true", func() {
				testCases := []string{
					"是非参加していただきたいと思います",
					"ぜひ参加していただきたいと思います",
					"ご参加いただきたく思います",
					"お会いできますことを",
					"心よりお待ちしております",
				}
				for _, s := range testCases {
					res := email.ClassifyAcceptanceMailByBody(s)
					Expect(res).To(Equal(true))
				}
			})
			It("returns false", func() {
				testCases := []string{
					"先日はお忙しい中、xxxの面接にご参加いただきまして誠にありがとうございました。",
				}
				for _, s := range testCases {
					res := email.ClassifyOinoriMailByBody(s)
					Expect(res).To(Equal(false))
				}
			})
		})
	})
})
