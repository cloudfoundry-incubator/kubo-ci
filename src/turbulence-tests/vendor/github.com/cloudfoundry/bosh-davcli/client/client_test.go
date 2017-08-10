package client_test

import (
	"errors"
	"io/ioutil"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/cloudfoundry/bosh-davcli/client"
	davconf "github.com/cloudfoundry/bosh-davcli/config"
	fakehttp "github.com/cloudfoundry/bosh-utils/http/fakes"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

var _ = Describe("Client", func() {
	var (
		fakeHTTPClient *fakehttp.FakeClient
		config         davconf.Config
		client         Client
		logger         boshlog.Logger
	)

	BeforeEach(func() {
		config.Endpoint = "http://example.com/"
		config.User = "some_user"
		config.Password = "some password"
		fakeHTTPClient = fakehttp.NewFakeClient()
		logger = boshlog.NewLogger(boshlog.LevelNone)
		client = NewClient(config, fakeHTTPClient, logger)
	})

	Describe("Exists", func() {
		BeforeEach(func() {
			fakeHTTPClient.StatusCode = 200
		})

		It("does not return an error if file exists", func() {
			err := client.Exists("/somefile")
			Expect(err).NotTo(HaveOccurred())
		})

		Context("the file does not exist", func() {
			BeforeEach(func() {
				fakeHTTPClient.StatusCode = 404
			})

			It("returns an error saying blob was not found", func() {
				err := client.Exists("/somefile")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Checking if dav blob /somefile exists: /somefile not found"))
			})
		})

		Context("unexpected http status code returned", func() {
			BeforeEach(func() {
				fakeHTTPClient.StatusCode = 601
			})

			It("returns an error saying an unexpected error occurred", func() {
				err := client.Exists("/somefile")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Checking if dav blob /somefile exists:"))
			})
		})
	})

	Describe("Delete", func() {
		Context("when the file does not exist", func() {
			BeforeEach(func() {
				fakeHTTPClient.StatusCode = 404
			})

			It("does not return an error if file does not exists", func() {
				fakeHTTPClient.StatusCode = 404

				err := client.Delete("/somefile")
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("when the file exists", func() {
			BeforeEach(func() {
				fakeHTTPClient.StatusCode = 204
			})

			It("does not return an error", func() {
				err := client.Delete("/somefile")
				Expect(err).ToNot(HaveOccurred())
				Expect(fakeHTTPClient.Requests[0].URL.Path).To(Equal("/19/somefile"))
				Expect(fakeHTTPClient.Requests[0].Method).To(Equal("DELETE"))
				Expect(fakeHTTPClient.Requests[0].Header["Authorization"]).To(Equal([]string{"Basic c29tZV91c2VyOnNvbWUgcGFzc3dvcmQ="}))
				Expect(fakeHTTPClient.Requests[0].Host).To(Equal("example.com"))
			})
		})

		Context("unexpected http status code returned", func() {
			BeforeEach(func() {
				fakeHTTPClient.StatusCode = 601
			})

			It("returns an error saying an unexpected error occurred", func() {
				err := client.Delete("/somefile")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Deleting blob '/somefile': Request failed, response: Response{ StatusCode: 601, Status: '' }"))
			})
		})
	})

	Describe("Get", func() {
		It("returns the response body from the given path", func() {
			fakeHTTPClient.StatusCode = 200
			fakeHTTPClient.SetMessage("response")

			responseBody, err := client.Get("/")
			Expect(err).NotTo(HaveOccurred())
			buf := make([]byte, 1024)
			n, _ := responseBody.Read(buf)
			Expect(string(buf[0:n])).To(Equal("response"))
		})

		Context("when the http request fails", func() {
			BeforeEach(func() {
				fakeHTTPClient.Error = errors.New("")
			})

			It("returns err", func() {
				responseBody, err := client.Get("/")
				Expect(responseBody).To(BeNil())
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Getting dav blob /"))
			})
		})

		Context("when the http response code is not 200", func() {
			BeforeEach(func() {
				fakeHTTPClient.StatusCode = 300
				fakeHTTPClient.SetMessage("response")
			})

			It("returns err", func() {
				responseBody, err := client.Get("/")
				Expect(responseBody).To(BeNil())
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Getting dav blob /: Request failed, response: Response{ StatusCode: 300, Status: '' }"))
				Expect(len(fakeHTTPClient.Requests)).To(Equal(3))
			})
		})
	})

	Describe("Put", func() {
		Context("When the put request succeeds", func() {
			itUploadsABlob := func() {
				body := ioutil.NopCloser(strings.NewReader("content"))
				err := client.Put("/", body, int64(7))
				Expect(err).NotTo(HaveOccurred())
				Expect(len(fakeHTTPClient.Requests)).To(Equal(1))
				req := fakeHTTPClient.Requests[0]
				Expect(req.ContentLength).To(Equal(int64(7)))
				Expect(fakeHTTPClient.RequestBodies).To(Equal([]string{"content"}))
			}

			It("uploads the given content if the blob does not exist", func() {
				fakeHTTPClient.StatusCode = 201
				itUploadsABlob()
			})

			It("uploads the given content if the blob exists", func() {
				fakeHTTPClient.StatusCode = 204
				itUploadsABlob()
			})
		})

		Context("when the http request fails", func() {
			BeforeEach(func() {
				fakeHTTPClient.Error = errors.New("EOF")
			})

			It("returns err", func() {
				body := ioutil.NopCloser(strings.NewReader("content"))
				err := client.Put("/", body, int64(7))
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Putting dav blob /: EOF"))
				Expect(len(fakeHTTPClient.Requests)).To(Equal(3))
			})
		})

		Context("when the http response code is not 201 or 204", func() {
			BeforeEach(func() {
				fakeHTTPClient.StatusCode = 300
				fakeHTTPClient.SetMessage("response")
			})

			It("returns err", func() {
				body := ioutil.NopCloser(strings.NewReader("content"))
				err := client.Put("/", body, int64(7))
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Putting dav blob /: Request failed, response: Response{ StatusCode: 300, Status: '' }"))
			})
		})
	})

	Describe("retryable count is configurable", func() {
		BeforeEach(func() {
			fakeHTTPClient.Error = errors.New("EOF")
			config = davconf.Config{RetryAttempts: 7}
			client = NewClient(config, fakeHTTPClient, logger)
		})

		It("tries the specified number of times", func() {
			body := ioutil.NopCloser(strings.NewReader("content"))
			err := client.Put("/", body, int64(7))
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Putting dav blob /: EOF"))
			Expect(len(fakeHTTPClient.Requests)).To(Equal(7))
		})

	})
})
