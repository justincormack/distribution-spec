package conformance

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/bloodorangeio/reggie"
	g "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var test03Discovery = func() {
	g.Context("Discovery", func() {

		var numTags = 4
		var tagList []string

		g.Context("Setup", func() {
			g.Specify("Populate registry with test tags", func() {
				RunOnlyIf(runDiscoverySetup)
				SkipIfDisabled(discovery)
				for i := 0; i < numTags; i++ {
					tag := fmt.Sprintf("test%d", i)
					tagList = append(tagList, tag)
					req := client.NewRequest(reggie.PUT, "/v2/<name>/manifests/<reference>",
						reggie.WithReference(tag)).
						SetHeader("Content-Type", "application/vnd.oci.image.manifest.v1+json").
						SetBody(manifestContent)
					_, err := client.Do(req)
					_ = err
					req = client.NewRequest(reggie.GET, "/v2/<name>/tags/list")
					resp, err := client.Do(req)
					tagList = getTagList(resp)
					_ = err
				}
			})

			g.Specify("Populate registry with test tags (no push)", func() {
				RunOnlyIfNot(runDiscoverySetup)
				SkipIfDisabled(discovery)
				tagList = strings.Split(os.Getenv(envVarTagList), ",")
			})
		})

		g.Context("Test discovery endpoints", func() {
			g.Specify("GET request to list tags should yield 200 response", func() {
				SkipIfDisabled(discovery)
				req := client.NewRequest(reggie.GET, "/v2/<name>/tags/list")
				resp, err := client.Do(req)
				Expect(err).To(BeNil())
				Expect(resp.StatusCode()).To(Equal(http.StatusOK))
				Expect(err).To(BeNil())
				numTags = len(tagList)
			})

			g.Specify("GET number of tags should be limitable by `n` query parameter", func() {
				SkipIfDisabled(discovery)
				numResults := numTags / 2
				req := client.NewRequest(reggie.GET, "/v2/<name>/tags/list").
					SetQueryParam("n", strconv.Itoa(numResults))
				resp, err := client.Do(req)
				Expect(err).To(BeNil())
				Expect(resp.StatusCode()).To(Equal(http.StatusOK))
				tagList = getTagList(resp)
				Expect(err).To(BeNil())
				Expect(len(tagList)).To(Equal(numResults))
			})

			g.Specify("GET start of tag is set by `last` query parameter", func() {
				SkipIfDisabled(discovery)
				numResults := numTags / 2
				req := client.NewRequest(reggie.GET, "/v2/<name>/tags/list").
					SetQueryParam("n", strconv.Itoa(numResults))
				resp, err := client.Do(req)
				Expect(err).To(BeNil())
				tagList = getTagList(resp)
				req = client.NewRequest(reggie.GET, "/v2/<name>/tags/list").
					SetQueryParam("n", strconv.Itoa(numResults)).
					SetQueryParam("last", tagList[numResults-1])
				resp, err = client.Do(req)
				Expect(err).To(BeNil())
				Expect(resp.StatusCode()).To(Equal(http.StatusOK))
				Expect(err).To(BeNil())
				Expect(len(tagList)).To(BeNumerically("<=", numResults))
				Expect(tagList).To(ContainElement(tagList[numResults-1]))
			})
		})

		g.Context("Teardown", func() {
			g.Specify("Delete created manifest & associated tags", func() {
				RunOnlyIf(runDiscoverySetup)
				SkipIfDisabled(discovery)
				req := client.NewRequest(reggie.DELETE, "/v2/<name>/manifests/<digest>", reggie.WithDigest(manifestDigest))
				_, err := client.Do(req)
				_ = err
			})
		})
	})
}
