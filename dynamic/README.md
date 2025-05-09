# Dynamic Package

The `dynamic` package provides a structured way to manage assets, track user interactions (such as clicks, impressions, and views), and organize items into groups. It is designed to handle responses that consist of grouped content with associated metadata and tracking information, making it suitable for applications like native advertising, content recommendation engines, and more.

## Table of Contents

- [Struct Descriptions](#struct-descriptions)
  - [`tracker`](#tracker)
  - [`assetThumb`](#assetthumb)
  - [`asset`](#asset)
  - [`item`](#item)
  - [`group`](#group)
  - [`Response`](#response)
- [Methods](#methods)
  - [`group.addItem`](#groupadditemi-item-group)
  - [`Response.getGroupOrCreate`](#responsegetgrouporcreategroupid-string-group)
- [Example Usage](#example-usage)
- [API Interaction Example](#api-interaction-example)
  - [Sample Curl Request](#sample-curl-request)
  - [Sample JSON Response](#sample-json-response)

## Struct Descriptions

### `tracker`

The `tracker` struct tracks user interactions for clicks, impressions, and views.

- **Fields:**
  - `Clicks` (`[]string`, optional): List of click tracking URLs.
  - `Impressions` (`[]string`): List of impression tracking URLs.
  - `Views` (`[]string`): List of view tracking URLs.

### `assetThumb`

Represents a thumbnail for an asset, including its path, type, width, and height.

- **Fields:**
  - `Path` (`string`): Path to the thumbnail.
  - `Type` (`string`, optional): Type of thumbnail.
  - `Width` (`int`, optional): Width of the thumbnail.
  - `Height` (`int`, optional): Height of the thumbnail.

### `asset`

Represents an asset with its main details and thumbnails.

- **Fields:**
  - `Name` (`string`, optional): Name of the asset.
  - `Path` (`string`): Path to the asset.
  - `Type` (`string`, optional): Type of asset (e.g., image, video).
  - `Width` (`int`, optional): Width of the asset.
  - `Height` (`int`, optional): Height of the asset.
  - `Thumbs` (`[]assetThumb`, optional): List of thumbnails for the asset.

### `item`

Represents a content item within a group. It can include assets, custom fields, and a tracker.

- **Fields:**
  - `ID` (`any`): Unique identifier for the item.
  - `Type` (`string`): Type of the item.
  - `URL` (`string`, optional): URL associated with the item.
  - `Content` (`string`, optional): Raw content of the item.
  - `ContentURL` (`string`, optional): URL for the content.
  - `Fields` (`map[string]any`, optional): Custom fields for additional metadata.
  - `Assets` (`[]asset`, optional): List of assets associated with the item.
  - `Tracker` (`tracker`): Tracking information for the item.
  - `Debug` (`any`, optional): Debug information related to the item.

### `group`

Represents a group containing multiple items.

- **Fields:**
  - `ID` (`string`): Unique identifier for the group.
  - `Items` (`[]*item`): List of items in the group.

### `Response`

The main response object, which contains groups of items.

- **Fields:**
  - `Version` (`string`): Version of the response format.
  - `Groups` (`[]*group`, optional): List of groups in the response.

## Methods

### `group.addItem(i *item) *group`

Adds an item to the group and returns the updated group.

**Parameters:**

- `i` (`*item`): The item to be added to the group.

**Returns:**

- `*group`: The updated group with the new item added.

### `Response.getGroupOrCreate(groupID string) *group`

Retrieves an existing group by `groupID` or creates a new group if it doesn't exist.

**Parameters:**

- `groupID` (`string`): The unique identifier for the group.

**Returns:**

- `*group`: The existing or newly created group.

## Example Usage

```go
package pkg

import (
    "fmt"
    "dynamic"
)

func do() {
    // Initialize a Response with version "1.0"
    response := &dynamic.Response{Version: "1.0"}

    // Retrieve an existing group or create a new one with ID "group1"
    group := response.getGroupOrCreate("group1")

    // Create a new item with tracking information
    item := &dynamic.item{
        ID:    "item1",
        Type:  "image",
        URL:   "https://example.com/image1.jpg",
        Tracker: dynamic.tracker{
            Clicks: []string{"https://example.com/click1"},
        },
    }

    // Add the item to the group
    group.addItem(item)

    // Print the response structure
    fmt.Printf("%+v\n", response)
}
```

This example initializes a Response, creates a new group or retrieves an existing one, and adds an item with tracking information.

## Request Parameters

The API accepts the following GET parameters to customize the response:

| Parameter  | Type     | Description |
|------------|----------|-------------|
| `debug`    | `bool`   | Enables debug mode if set to `true`. |
| `x`        | `int`    | Optional coordinate X for positioning or targeting logic. |
| `y`        | `int`    | Optional coordinate Y for positioning or targeting logic. |
| `w`        | `int`    | Desired maximum width of the asset. |
| `h`        | `int`    | Desired maximum height of the asset. |
| `mw`       | `int`    | Minimum width constraint. |
| `mh`       | `int`    | Minimum height constraint. |
| `fmt`      | `string` | Alternative shorthand for specifying size, e.g. `300x250`. Can override `w` and `h`. |
| `width`    | `int`    | Minimum required width (alias for `mw`). |
| `height`   | `int`    | Minimum required height (alias for `mh`). |
| `format`   | `string` | Comma-separated list of format codes. Accepts `auto`, `all`, or specific formats. |
| `type`     | `string` | Comma-separated list of format types. Accepts `auto`, `all`, or specific types. |
| `keywords` | `string` | Comma-separated keywords to match content. Aliases: `keyword`, `kw`. |
| `count`    | `int`    | Desired number of items in the response. |
| `subid1`   | `string` | Custom tracking ID. Aliases: `subid`, `s1`. |
| `subid2`   | `string` | Additional tracking ID. Alias: `s2`. |
| `subid3`   | `string` | Additional tracking ID. Alias: `s3`. |
| `subid4`   | `string` | Additional tracking ID. Alias: `s4`. |
| `subid5`   | `string` | Additional tracking ID. Alias: `s5`. |

## API Interaction Example

This section demonstrates how to interact with the system using a curl request and the corresponding JSON response format.

### Sample Curl Request

To retrieve the response data, use the following curl command. Note that since it’s a GET request, no request body is needed.

```sh
curl -X GET 'http://localhost:8080/api/response?format=auto&type=banner&w=100&h=50' \
     -H "Accept: application/json"
```

### Sample JSON Response

```json
{
  "version": "1",
  "groups": [
    {
      "id": "unique-group-id",
      "items": [
        {
          "id": "unique-item-id",
          "type": "item-type",
          "url": "item-click-url",
          "fields": {
            "brandname": "BrandName",
            "description": "Item Description",
            "title": "Item Title"
          },
          "assets": [
            {
              "name": "asset-name",
              "path": "asset-url",
              "type": "asset-type",
              "width": 200,
              "height": 200
            }
          ],
          "tracker": {
            "impressions": [
              "impression-tracking-url-1",
              "impression-tracking-url-2"
            ],
            "views": [
              "view-tracking-url-1",
              "view-tracking-url-2"
            ],
            "clicks": [
              "click-tracking-url-1",
              "click-tracking-url-2"
            ]
          }
        }
      ]
    }
  ]
}
```

This JSON response showcases the structure of the Response object, including a group containing an item with associated assets and tracking URLs for impressions and views. The identifiers and URLs are anonymized for security and privacy.

## Summary

The dynamic package offers a comprehensive way to structure and manage dynamic content, making it easier to handle grouped items with rich metadata and tracking capabilities. Whether you’re building a native advertising platform, a content recommendation system, or any application that requires organized and trackable content groups, this package provides the necessary tools to streamline your workflow.

For further information or assistance, feel free to open an issue or contribute to the repository.
