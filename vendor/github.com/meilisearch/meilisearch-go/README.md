<p align="center">
  <img src="https://raw.githubusercontent.com/meilisearch/integration-guides/main/assets/logos/meilisearch_go.svg" alt="Meilisearch-Go" width="200" height="200" />
</p>

<h1 align="center">Meilisearch Go</h1>

<h4 align="center">
  <a href="https://github.com/meilisearch/meilisearch">Meilisearch</a> |
  <a href="https://docs.meilisearch.com">Documentation</a> |
  <a href="https://slack.meilisearch.com">Slack</a> |
  <a href="https://roadmap.meilisearch.com/tabs/1-under-consideration">Roadmap</a> |
  <a href="https://www.meilisearch.com">Website</a> |
  <a href="https://docs.meilisearch.com/faq">FAQ</a>
</h4>

<p align="center">
  <a href="https://github.com/meilisearch/meilisearch-go/actions"><img src="https://github.com/meilisearch/meilisearch-go/workflows/Tests/badge.svg" alt="GitHub Workflow Status"></a>
  <a href="https://goreportcard.com/report/github.com/meilisearch/meilisearch-go"><img src="https://goreportcard.com/badge/github.com/meilisearch/meilisearch-go" alt="Test"></a>
  <a href="https://github.com/meilisearch/meilisearch-go/blob/main/LICENSE"><img src="https://img.shields.io/badge/license-MIT-informational" alt="License"></a>
  <a href="https://app.bors.tech/repositories/28783"><img src="https://bors.tech/images/badge_small.svg" alt="Bors enabled"></a>
</p>

<p align="center">⚡ The Meilisearch API client written for Golang</p>

**Meilisearch Go** is the Meilisearch API client for Go developers.

**Meilisearch** is an open-source search engine. [Learn more about Meilisearch.](https://github.com/meilisearch/Meilisearch)

## Table of Contents <!-- omit in toc -->

- [🎃 Hacktoberfest](#-hacktoberfest)
- [📖 Documentation](#-documentation)
- [🔧 Installation](#-installation)
- [🚀 Getting started](#-getting-started)
- [🤖 Compatibility with Meilisearch](#-compatibility-with-meilisearch)
- [💡 Learn more](#-learn-more)
- [⚙️ Contributing](#️-contributing)

## 🎃 Hacktoberfest

It’s Hacktoberfest 2022 @Meilisearch

[Hacktoberfest](https://hacktoberfest.com/) is a celebration of the open-source community. This year, and for the third time in a row, Meilisearch is participating in this fantastic event.

You’d like to contribute? Don’t hesitate to check out our [contributing guidelines](./CONTRIBUTING.md).

## 📖 Documentation

This readme contains all the documentation you need to start using this Meilisearch SDK.

For general information on how to use Meilisearch—such as our API reference, tutorials, guides, and in-depth articles—refer to our [main documentation website](https://docs.meilisearch.com/).


## 🔧 Installation

With `go get` in command line:
```bash
go get github.com/meilisearch/meilisearch-go
```

### Run Meilisearch <!-- omit in toc -->

There are many easy ways to [download and run a Meilisearch instance](https://docs.meilisearch.com/reference/features/installation.html#download-and-launch).

For example, using the `curl` command in [your Terminal](https://itconnect.uw.edu/learn/workshops/online-tutorials/web-publishing/what-is-a-terminal/):

```bash
# Install Meilisearch
curl -L https://install.meilisearch.com | sh

# Launch Meilisearch
./meilisearch --master-key=masterKey
```

NB: you can also download Meilisearch from **Homebrew** or **APT** or even run it using **Docker**.

## 🚀 Getting started

#### Add documents <!-- omit in toc -->

```go
package main

import (
	"fmt"
	"os"

	"github.com/meilisearch/meilisearch-go"
)

func main() {
	client := meilisearch.NewClient(meilisearch.ClientConfig{
                Host: "http://127.0.0.1:7700",
                APIKey: "masterKey",
        })
	// An index is where the documents are stored.
	index := client.Index("movies")

	// If the index 'movies' does not exist, Meilisearch creates it when you first add the documents.
	documents := []map[string]interface{}{
        { "id": 1, "title": "Carol", "genres": []string{"Romance", "Drama"} },
        { "id": 2, "title": "Wonder Woman", "genres": []string{"Action", "Adventure"} },
        { "id": 3, "title": "Life of Pi", "genres": []string{"Adventure", "Drama"} },
        { "id": 4, "title": "Mad Max: Fury Road", "genres": []string{"Adventure", "Science Fiction"} },
        { "id": 5, "title": "Moana", "genres": []string{"Fantasy", "Action"} },
        { "id": 6, "title": "Philadelphia", "genres": []string{"Drama"} },
	}
	task, err := index.AddDocuments(documents)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(task.TaskUID)
}
```

With the `taskUID`, you can check the status (`enqueued`, `processing`, `succeeded` or `failed`) of your documents addition using the [task endpoint](https://docs.meilisearch.com/learn/advanced/asynchronous_operations.html#task-status).

#### Basic Search <!-- omit in toc -->

```go
package main

import (
    "fmt"
    "os"

    "github.com/meilisearch/meilisearch-go"
)

func main() {
    // Meilisearch is typo-tolerant:
    searchRes, err := client.Index("movies").Search("philoudelphia",
        &meilisearch.SearchRequest{
            Limit: 10,
        })
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    fmt.Println(searchRes.Hits)
}
```

JSON output:
```json
{
  "hits": [{
    "id": 6,
    "title": "Philadelphia",
    "genres": ["Drama"]
  }],
  "offset": 0,
  "limit": 10,
  "processingTimeMs": 1,
  "query": "philoudelphia"
}
```

#### Custom Search <!-- omit in toc -->

All the supported options are described in the [search parameters](https://docs.meilisearch.com/reference/features/search_parameters.html) section of the documentation.

```go
func main() {
    searchRes, err := client.Index("movies").Search("wonder",
        &meilisearch.SearchRequest{
            AttributesToHighlight: []string{"*"},
        })
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    fmt.Println(searchRes.Hits)
}
```

JSON output:
```json
{
    "hits": [
        {
            "id": 2,
            "title": "Wonder Woman",
            "genres": ["Action", "Adventure"],
            "_formatted": {
                "id": 2,
                "title": "<em>Wonder</em> Woman"
            }
        }
    ],
    "offset": 0,
    "limit": 20,
    "processingTimeMs": 0,
    "query": "wonder"
}
```

#### Custom Search With Filters <!-- omit in toc -->

If you want to enable filtering, you must add your attributes to the `filterableAttributes` index setting.

```go
task, err := index.UpdateFilterableAttributes(&[]string{"id", "genres"})
```

You only need to perform this operation once.

Note that Meilisearch will rebuild your index whenever you update `filterableAttributes`. Depending on the size of your dataset, this might take time. You can track the process using the [task status](https://docs.meilisearch.com/learn/advanced/asynchronous_operations.html#task-status).

Then, you can perform the search:

```go
searchRes, err := index.Search("wonder",
    &meilisearch.SearchRequest{
        Filter: "id > 1 AND genres = Action",
    })
```

```json
{
  "hits": [
    {
      "id": 2,
      "title": "Wonder Woman",
      "genres": ["Action","Adventure"]
    }
  ],
  "offset": 0,
  "limit": 20,
  "estimatedTotalHits": 1,
  "processingTimeMs": 0,
  "query": "wonder"
}
```

## 🤖 Compatibility with Meilisearch

This package only guarantees the compatibility with the [version v0.29.0 of Meilisearch](https://github.com/meilisearch/meilisearch/releases/tag/v0.29.0).

## 💡 Learn more

The following sections in our main documentation website may interest you:

- **Manipulate documents**: see the [API references](https://docs.meilisearch.com/reference/api/documents.html) or read more about [documents](https://docs.meilisearch.com/learn/core_concepts/documents.html).
- **Search**: see the [API references](https://docs.meilisearch.com/reference/api/search.html) or follow our guide on [search parameters](https://docs.meilisearch.com/reference/features/search_parameters.html).
- **Manage the indexes**: see the [API references](https://docs.meilisearch.com/reference/api/indexes.html) or read more about [indexes](https://docs.meilisearch.com/learn/core_concepts/indexes.html).
- **ClientConfigure the index settings**: see the [API references](https://docs.meilisearch.com/reference/api/settings.html) or follow our guide on [settings parameters](https://docs.meilisearch.com/reference/features/settings.html).

## ⚙️ Contributing

Any new contribution is more than welcome in this project!

If you want to know more about the development workflow or want to contribute, please visit our [contributing guidelines](/CONTRIBUTING.md) for detailed instructions!

<hr>

**Meilisearch** provides and maintains many **SDKs and Integration tools** like this one. We want to provide everyone with an **amazing search experience for any kind of project**. If you want to contribute, make suggestions, or just know what's going on right now, visit us in the [integration-guides](https://github.com/meilisearch/integration-guides) repository.
