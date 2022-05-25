package controllers

import (
	"context"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	hnewsv1 "github.com/vadasambar/hnews/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("HNews Controller", func() {
	Context("When creating HNews", func() {
		It("It should create HNews from an empty `spec`", func() {
			By("By filling in the defaults for empty `spec`")
			ctx := context.Background()
			hnews := &hnewsv1.HNews{
				TypeMeta: metav1.TypeMeta{
					APIVersion: "apps.vadasambar.com/v1",
					Kind:       "HNews",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "hnews-sample",
					Namespace: "default",
				},
				Spec: hnewsv1.HNewsSpec{},
			}

			Expect(k8sClient.Create(ctx, hnews)).Should(Succeed())

			Eventually(func() bool {
				var hnewsCreated hnewsv1.HNews
				err := k8sClient.Get(ctx, types.NamespacedName{Name: "hnews-sample", Namespace: "default"}, &hnewsCreated)
				Expect(err).NotTo(HaveOccurred())

				filter := hnewsCreated.Spec.Filter
				return filter.Descendants == defaultDescendents &&
					filter.Limit == defaultLimit &&
					filter.Score == defaultScore &&
					filter.Type == defaultType &&
					len(hnewsCreated.Status.Links) <= defaultLimit

			}, time.Second*30, time.Second*2).Should(BeTrue())
		})

	})
})
