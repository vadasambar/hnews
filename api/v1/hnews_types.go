/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

type Type string

// Based on the possible types specified in https://github.com/HackerNews/API#items
const (
	Story   Type = "story"
	Comment Type = "comment"
	Job     Type = "job"
	Poll    Type = "poll"
	PollOpt Type = "pollopt"
)

// Generated using https://mholt.github.io/json-to-go/
// by converting get Id http json response to go struct
type GetIdResponse struct {
	By          string `json:"by"`
	Descendants int    `json:"descendants"`
	ID          int    `json:"id"`
	Kids        []int  `json:"kids"`
	Score       int    `json:"score"`
	Time        int    `json:"time"`
	Title       string `json:"title"`
	Type        Type   `json:"type"`
	URL         string `json:"url"`
}

// Filter allows you to filter and get the
// Hacker News articles you want
type Filter struct {
	// Number of Hacker News articles you want.
	// +kubebuilder:validation:Maximum:=20
	Limit int `json:"limit"`
	// Type of Hacker News articles you are looking for.
	// Has to be either of: job,story,comment,poll,pollopt
	// +kubebuilder:validation:Enum:=job;story;comment;poll;pollopt
	Type string `json:"type,omitempty"`
	// Score of Hacker News articles you are looking for.
	// Specify it like:
	// score: ">=10", score: "<10", score: "=10", score: "!=10"
	Score Comparison `json:"score"`
	// Number of direct (first level) comments in the article.
	Descendants Comparison `json:"descendents"`
}

type Comparison string

// HNewsSpec defines the desired state of HNews
type HNewsSpec struct {
	Filter Filter `json:"filter,omitempty"`
}

// HNewsStatus defines the observed state of HNews
type HNewsStatus struct {
	// Important: Run "make" to regenerate code after modifying this file
	Links []Link `json:"link"`
}

// Link holds the information about
// Hacker News article for which satisfies the filter
type Link struct {
	// HNewsUrl refers to the URL of the HNews page
	// e.g., https://news.ycombinator.com/item?id=31316372
	HNewsUrl string `json:"hnews_url"`
	// ArticleUrl refers to the URL which is shared on the HNews page above
	// e.g., https://swelltype.com/yep-i-created-the-new-avatar-font/
	ArticleUrl  string `json:"article_url"`
	Descendents int    `json:"descendents"`
	Score       int    `json:"score"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// HNews is the Schema for the hnews API
type HNews struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   HNewsSpec   `json:"spec"`
	Status HNewsStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// HNewsList contains a list of HNews
type HNewsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HNews `json:"items"`
}

func init() {
	SchemeBuilder.Register(&HNews{}, &HNewsList{})
}
