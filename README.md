# AdStdEndpoints

Standard client ad-request endpoints providing comprehensive ad serving capabilities for different formats: direct redirects, dynamic/native content, and proxy-based HTML rendering.

## Table of Contents

- [Overview](#overview)
- [Endpoint Types](#endpoint-types)
  - [Direct Endpoint](#direct-endpoint)
  - [Dynamic Endpoint](#dynamic-endpoint)
  - [Proxy Endpoint](#proxy-endpoint)
- [Protocol Documentation](#protocol-documentation)
  - [Request Parameters](#request-parameters)
  - [Response Formats](#response-formats)
- [API Examples](#api-examples)
- [Debug Mode](#debug-mode)
- [Event Tracking](#event-tracking)
- [License](#license)

## Overview

AdStdEndpoints provides three distinct endpoint types for different ad serving scenarios:

- **Direct**: Immediate redirects for popunders, click-through ads, and direct action campaigns
- **Dynamic**: JSON-based responses for native advertising, rich media, and programmatic content
- **Proxy**: HTML template-based rendering for embedded banner ads and iframe content

Each endpoint type supports comprehensive event tracking, debug capabilities, and flexible response formatting.

## Endpoint Types

### Direct Endpoint

The direct endpoint (`/direct`) handles immediate redirects and is optimized for:

- Popunder campaigns
- Direct action advertising
- Click-through campaigns
- Simple redirect-based monetization

**Key Features:**

- Immediate HTTP redirects (302 Found)
- Superfailover URL support
- Alternative link handling
- Debug mode with detailed response information

### Dynamic Endpoint

The dynamic endpoint (`/dynamic`) provides structured JSON responses for:

- Native advertising integration
- Rich media content delivery
- Programmatic advertising
- Mobile app monetization
- Content recommendation systems

**Key Features:**

- Flexible JSON/JSONP responses
- Asset management with thumbnails
- Advanced tracking capabilities
- Grouped content organization
- Meta information for compliance

### Proxy Endpoint

The proxy endpoint (`/proxy`) renders HTML templates for:

- Banner ad placement
- Iframe-based advertising
- Legacy ad integration
- Direct HTML embedding

**Key Features:**

- Server-side HTML rendering
- Template-based customization
- Direct HTML embedding
- Legacy system compatibility

## Protocol Documentation

### Request Parameters

All endpoints support the following common GET parameters:

| Parameter  | Type     | Description | Aliases |
|------------|----------|-------------|---------|
| `debug`    | `bool`   | Enables debug mode with detailed response information | - |
| `x`        | `int`    | X coordinate for positioning/targeting | - |
| `y`        | `int`    | Y coordinate for positioning/targeting | - |
| `w`        | `int`    | Maximum desired width | - |
| `h`        | `int`    | Maximum desired height | - |
| `mw`       | `int`    | Minimum width constraint | `width` |
| `mh`       | `int`    | Minimum height constraint | `height` |
| `fmt`      | `string` | Size format (e.g., `300x250`) | - |
| `format`   | `string` | Comma-separated format codes (`auto`, `all`, or specific) | - |
| `type`     | `string` | Comma-separated format types (`auto`, `all`, or specific) | - |
| `keywords` | `string` | Comma-separated targeting keywords | `keyword`, `kw` |
| `count`    | `int`    | Desired number of items | - |
| `subid1`   | `string` | Primary tracking identifier | `subid`, `s1` |
| `subid2`   | `string` | Secondary tracking identifier | `s2` |
| `subid3`   | `string` | Tertiary tracking identifier | `s3` |
| `subid4`   | `string` | Quaternary tracking identifier | `s4` |
| `subid5`   | `string` | Quinary tracking identifier | `s5` |

#### Dynamic Endpoint Specific Parameters

| Parameter  | Type     | Description |
|------------|----------|-------------|
| `format`   | `string` | Response format: `json` (default) or `jsonp` |
| `callback` | `string` | JSONP callback function name (default: `callback`) |

#### Direct Endpoint Specific Parameters

| Parameter    | Type   | Description |
|--------------|--------|-------------|
| `noredirect` | `bool` | Return JSON debug response instead of redirect |

### Response Formats

#### Direct Endpoint Response

**Normal Operation:**

- **Success with ad**: `HTTP 302 Found` redirect to ad URL
- **Alternative link**: `HTTP 302 Found` redirect with `X-Status-Alternative: 1` header
- **No ad/failover**: `HTTP 302 Found` redirect to superfailover URL with `X-Status-Failover: 1` header

**Debug Mode (`debug=true&noredirect=true`):**

```json
{
  "id": "ad-unit-id",
  "zone_id": 12345,
  "auction_id": "auction-uuid",
  "impression_id": "impression-uuid",
  "is_alternative_link": false,
  "link": "https://example.com/click-url",
  "superfailover": "https://fallback.com",
  "error": null,
  "is_empty": false
}
```

#### Dynamic Endpoint Response

**JSON Structure:**

```json
{
  "version": "1",
  "custom_tracker": {
    "impressions": ["https://track.example.com/imp/1"],
    "views": ["https://track.example.com/view/1"],
    "clicks": ["https://track.example.com/click/1"]
  },
  "groups": [
    {
      "id": "impression-id-1",
      "custom_tracker": {
        "impressions": ["https://track.example.com/custom/imp/1"],
        "views": ["https://track.example.com/custom/view/1"],
        "clicks": ["https://track.example.com/custom/click/1"]
      },
      "items": [
        {
          "id": "ad-item-id",
          "type": "banner",
          "url": "https://example.com/click",
          "content": "<div>Ad Content</div>",
          "content_url": "https://example.com/iframe.html",
          "fields": {
            "title": "Ad Title",
            "description": "Ad Description",
            "brandname": "Brand Name",
            "call_to_action": "Click Here"
          },
          "assets": [
            {
              "name": "main-image",
              "path": "https://cdn.example.com/image.jpg",
              "type": "image",
              "width": 300,
              "height": 250,
              "thumbs": [
                {
                  "path": "https://cdn.example.com/thumb.jpg",
                  "type": "image",
                  "width": 100,
                  "height": 83
                }
              ]
            }
          ],
          "tracker": {
            "impressions": [
              "https://track.example.com/imp/item/1",
              "https://third-party.com/imp/123"
            ],
            "views": [
              "https://track.example.com/view/item/1",
              "https://third-party.com/view/123"
            ],
            "clicks": [
              "https://third-party.com/click/123"
            ]
          },
          "meta": {
            "advertiser": {
              "id": 123,
              "name": "Advertiser Name",
              "about_url": "https://advertiser.com/about",
              "contact_url": "https://advertiser.com/contact",
              "privacy_url": "https://advertiser.com/privacy",
              "terms_url": "https://advertiser.com/terms"
            },
            "ad": {
              "id": 456,
              "campaign_id": 789,
              "description": "Advertisement Description",
              "min_age": 18,
              "about_url": "https://advertiser.com/ad/about",
              "contact_url": "https://advertiser.com/ad/contact",
              "privacy_url": "https://advertiser.com/ad/privacy",
              "terms_url": "https://advertiser.com/ad/terms"
            },
            "complaint_url": "https://example.com/complaint",
            "hide": {
              "type": "cookie",
              "hide_ad_url_type": "post",
              "hide_ad_url": "https://example.com/hide",
              "hide_ad_url_params": {
                "ad_id": "456"
              },
              "name": "hide_ad_456"
            }
          },
          "debug": {
            "adUnit": {
              "internal_id": "internal-123",
              "auction_details": "..."
            }
          }
        }
      ]
    }
  ],
  "debug": {
    "http": {
      "uri": "/dynamic?w=300&h=250",
      "ip": "192.168.1.1",
      "method": "GET",
      "query": "w=300&h=250&debug=true",
      "headers": {
        "user-agent": "Mozilla/5.0...",
        "accept": "application/json"
      }
    }
  }
}
```

**Response Structure Explanation:**

- **`version`**: API version identifier
- **`custom_tracker`**: Global tracking URLs applied to all items
- **`groups`**: Array of impression groups, each containing related items
  - **`id`**: Impression/group identifier
  - **`custom_tracker`**: Group-specific tracking for empty responses
  - **`items`**: Array of ad items in this group
    - **`id`**: Unique item identifier
    - **`type`**: Ad format type (banner, native, video, etc.)
    - **`url`**: Click-through URL
    - **`content`**: Raw HTML/text content
    - **`content_url`**: URL for iframe-based content
    - **`fields`**: Dynamic key-value pairs for ad content
    - **`assets`**: Media files (images, videos) with metadata
    - **`tracker`**: Item-specific tracking URLs
    - **`meta`**: Compliance and advertiser information
    - **`debug`**: Debug information (only in debug mode)

#### Proxy Endpoint Response

Returns rendered HTML content directly:

```html
<!DOCTYPE html>
<html>
<head>
    <title>Ad Content</title>
</head>
<body>
    <div class="ad-container">
        <!-- Rendered ad content -->
    </div>
</body>
</html>
```

## API Examples

### Direct Endpoint Examples

**Basic Request:**

```bash
curl -X GET 'https://api.example.com/direct?zone=123&w=300&h=250'
# Response: HTTP 302 redirect to ad URL or superfailover
```

**Debug Request:**

```bash
curl -X GET 'https://api.example.com/direct?zone=123&debug=true&noredirect=true' \
     -H "Accept: application/json"
```

### Dynamic Endpoint Examples

**Basic JSON Request:**

```bash
curl -X GET 'https://api.example.com/dynamic?zone=123&w=300&h=250&count=3' \
     -H "Accept: application/json"
```

**JSONP Request:**

```bash
curl -X GET 'https://api.example.com/dynamic?zone=123&format=jsonp&callback=handleAds' \
     -H "Accept: application/javascript"
```

**Native Ad Request:**

```bash
curl -X GET 'https://api.example.com/dynamic?zone=123&type=native&keywords=technology,mobile' \
     -H "Accept: application/json"
```

### Proxy Endpoint Examples

**HTML Banner Request:**

```bash
curl -X GET 'https://api.example.com/proxy?zone=123&w=728&h=90' \
     -H "Accept: text/html"
```

## Debug Mode

Enable debug mode by adding `debug=true` to any request. Debug mode provides:

- Detailed request information (headers, IP, query parameters)
- Internal ad unit details
- Auction information
- Error details
- Performance metrics

**Debug mode should only be used in development/testing environments.**

## Event Tracking

The system provides comprehensive event tracking through pixel URLs:

- **Impressions**: Fired when ad content is served
- **Views**: Fired when ad becomes viewable to user
- **Clicks**: Fired when user interacts with ad
- **Custom**: Application-specific events

Tracking URLs support both first-party (system-generated) and third-party (advertiser-provided) pixels.

## Integration Examples

### JavaScript Integration (Dynamic)

```javascript
// Fetch ads dynamically
fetch('/dynamic?zone=123&w=300&h=250&count=2')
  .then(response => response.json())
  .then(data => {
    data.groups.forEach(group => {
      group.items.forEach(item => {
        // Render ad item
        renderAd(item);
        
        // Fire impression tracking
        item.tracker.impressions.forEach(url => {
          new Image().src = url;
        });
      });
    });
  });

function renderAd(item) {
  const container = document.createElement('div');
  container.innerHTML = item.content;
  container.onclick = () => {
    // Fire click tracking
    item.tracker.clicks.forEach(url => {
      new Image().src = url;
    });
    // Navigate to ad URL
    window.open(item.url, '_blank');
  };
  document.body.appendChild(container);
}
```

### Server-Side Integration (Direct)

```go
// Redirect user to ad or superfailover
http.Redirect(w, r, "/direct?zone=123&subid1=user123", http.StatusFound)
```

## Error Handling

The system provides graceful error handling:

- **No ads available**: Redirects to superfailover (direct) or returns empty groups (dynamic)
- **Invalid parameters**: Returns appropriate HTTP error codes
- **Network issues**: Implements timeout and retry mechanisms
- **Debug information**: Available in debug mode for troubleshooting

## License

[LICENSE](LICENSE)

Copyright 2024 Dmitry Ponomarev & Geniusrabbit

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at

<http://www.apache.org/licenses/LICENSE-2.0>

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
