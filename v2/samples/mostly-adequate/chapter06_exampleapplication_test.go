// Copyright (c) 2023 - 2025 IBM Corp.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mostlyadequate

import (
	"context"
	"fmt"
	"net/http"
	"regexp"

	A "github.com/IBM/fp-go/v2/array"
	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	J "github.com/IBM/fp-go/v2/json"
	S "github.com/IBM/fp-go/v2/string"

	R "github.com/IBM/fp-go/v2/context/readerioresult"
	H "github.com/IBM/fp-go/v2/context/readerioresult/http"
)

type (
	FlickrMedia struct {
		Link string `json:"m"`
	}

	FlickrItem struct {
		Media FlickrMedia `json:"media"`
	}

	FlickrFeed struct {
		Items []FlickrItem `json:"items"`
	}
)

func (f FlickrMedia) getLink() string {
	return f.Link
}

func (f FlickrItem) getMedia() FlickrMedia {
	return f.Media
}

func (f FlickrFeed) getItems() []FlickrItem {
	return f.Items
}

func Example_application() {
	// pure
	host := "api.flickr.com"
	path := "/services/feeds/photos_public.gne"
	query := S.Format[string]("?tags=%s&format=json&jsoncallback=?")
	url := F.Flow2(
		query,
		S.Format[string](fmt.Sprintf("https://%s%s%%s", host, path)),
	)
	// flick returns jsonP, we extract the JSON body, this is handled by jquery in the original code
	sanitizeJSONP := Replace(regexp.MustCompile(`(?s)^\s*\((.*)\)\s*$`))("$1")
	// parse jsonP
	parseJSONP := F.Flow3(
		sanitizeJSONP,
		S.ToBytes,
		J.Unmarshal[FlickrFeed],
	)
	// markup
	img := S.Format[string]("<img src='%s'/>")
	// lenses
	mediaURL := F.Flow2(
		FlickrItem.getMedia,
		FlickrMedia.getLink,
	)
	mediaURLs := F.Flow2(
		FlickrFeed.getItems,
		A.Map(mediaURL),
	)
	images := F.Flow2(
		mediaURLs,
		A.Map(img),
	)

	client := H.MakeClient(http.DefaultClient)

	// func(string) R.ReaderIOEither[[]string]
	app := F.Flow5(
		url,
		H.MakeGetRequest,
		H.ReadText(client),
		R.ChainEitherK(parseJSONP),
		R.Map(images),
	)

	// R.ReaderIOEither[[]string]
	// this is the managed effect that can be called to download and render the images
	catImageEffect := app("cats")

	// impure, actually executes the effect
	catImages := catImageEffect(context.TODO())()
	fmt.Println(E.IsRight(catImages))

	// Output:
	// true

}
