# Dynamic Package

The `dynamic` package provides a comprehensive ad serving solution with structured JSON responses for different ad formats. It manages assets, tracks user interactions (clicks, impressions, views), and organizes content into groups, making it ideal for native advertising, programmatic campaigns, rich media, and modern ad serving scenarios.

## Table of Contents

- [Overview](#overview)
- [Ad Format Types](#ad-format-types)
- [Struct Descriptions](#struct-descriptions)
- [Ad Format Examples](#ad-format-examples)
  - [Native Ads](#native-ads)
  - [Banner Ads](#banner-ads)
  - [Slider Banner Ads](#slider-banner-ads)
  - [Slider Video Ads](#slider-video-ads)
  - [Proxy Ads](#proxy-ads)
- [Integration Examples](#integration-examples)
- [Request Parameters](#request-parameters)

## Overview

The dynamic endpoint serves structured JSON responses optimized for:

- **Native Advertising**: Content-style ads that blend with editorial content
- **Rich Media Campaigns**: Interactive banners with multiple assets
- **Video Advertising**: Pre-roll, mid-roll, and native video content
- **Programmatic Buying**: Real-time bidding integration
- **Mobile Applications**: App monetization and in-app advertising
- **Content Recommendation**: Sponsored content and related articles

## Ad Format Types

| Format Type | Description | Use Cases |
|-------------|-------------|-----------|
| `native` | Content-style ads with title, description, images, and videos | Editorial integration, sponsored content, video content |
| `banner` | Traditional display banners with images, videos, and HTML | Website monetization, display campaigns, rich media |
| `slider_banner` | Multi-asset banner carousels | Product showcases, brand campaigns |
| `slider_video` | Video carousel presentations | Entertainment, media content |
| `proxy` | Server-rendered HTML content with `content` or `content_url` | Legacy systems, iframe integration |

## Struct Descriptions

### `tracker`

The `tracker` struct manages comprehensive event tracking for ad interactions.

**Fields:**

- **`Clicks`** (`[]string`, optional): Click tracking URLs fired when user interacts with ad
- **`Impressions`** (`[]string`): Impression tracking URLs fired when ad is served
- **`Views`** (`[]string`): View tracking URLs fired when ad becomes viewable

**Usage:** Supports both first-party system tracking and third-party advertiser pixels.

### `assetThumb`

Represents thumbnails and preview images for assets with size information.

**Fields:**

- **`Path`** (`string`): CDN URL to the thumbnail image
- **`Type`** (`string`, optional): Thumbnail format (`image`, `webp`, `jpeg`, etc.)
- **`Width`** (`int`, optional): Thumbnail width in pixels
- **`Height`** (`int`, optional): Thumbnail height in pixels

### `asset`

Represents media assets (images, videos, documents) with metadata and thumbnails.

**Fields:**

- **`Name`** (`string`, optional): Asset identifier (`main_image`, `logo`, `video`, etc.)
- **`Path`** (`string`): CDN URL to the full asset
- **`Type`** (`string`, optional): Asset type (`image`, `video`, `audio`, `document`)
- **`Width`** (`int`, optional): Asset width in pixels (for images/videos)
- **`Height`** (`int`, optional): Asset height in pixels (for images/videos)
- **`Thumbs`** (`[]assetThumb`, optional): Array of thumbnail variations

### `item`

Represents an individual ad unit with content, assets, and tracking.

**Fields:**

- **`ID`** (`any`): Unique identifier for the ad item
- **`Type`** (`string`): Ad format type (`native`, `banner`, `video_banner`, etc.)
- **`URL`** (`string`, optional): Click-through destination URL
- **`Content`** (`string`, optional): Raw HTML/text content for direct rendering (proxy ads only, mutually exclusive with `ContentURL`)
- **`ContentURL`** (`string`, optional): IFrame URL for proxy content delivery (proxy ads only, mutually exclusive with `Content`)
- **`Fields`** (`map[string]any`, optional): Dynamic key-value pairs for ad content
- **`Assets`** (`[]asset`, optional): Media files associated with the ad (images, videos, etc.)

**Note:** `Content` and `ContentURL` are supported only by proxy ads and are mutually exclusive:

- **Proxy ads with HTML injection**: Use `Content` with raw HTML
- **Proxy ads with iframe**: Use `ContentURL` with iframe source URL  
- **Banner/Native ads**: Use `URL` for click destination and `Assets` for media content (images, videos, etc.)
- **`Tracker`** (`tracker`): Event tracking configuration
- **`Meta`** (`*itemMetaInfo`, optional): Advertiser and compliance information
- **`Debug`** (`any`, optional): Debug information (development mode only)

### `group`

Container for related ad items, typically organized by impression or placement.

**Fields:**

- **`ID`** (`string`): Group identifier (usually impression ID)
- **`CustomTracker`** (`tracker`, optional): Group-level tracking for empty responses
- **`Items`** (`[]*item`): Array of ad items in this group

### `Response`

Root response object containing all ad groups and global metadata.

**Fields:**

- **`Version`** (`string`): API version identifier (currently "1")
- **`CustomTracker`** (`tracker`, optional): Global tracking applied to all items
- **`Groups`** (`[]*group`, optional): Array of ad groups
- **`Debug`** (`any`, optional): Request debug information

## Ad Format Examples

### Native Ads

Native ads blend seamlessly with editorial content, providing non-intrusive advertising experiences. They support both image and video content through the assets array.

**Image Native Ad Request:**

```bash
curl -X GET 'https://api.example.com/dynamic?zone=123&type=native&count=3&keywords=technology' \
     -H "Accept: application/json"
```

**Image Native Ad Response:**

```json
{
  "version": "1",
  "groups": [
    {
      "id": "impression-native-001",
      "items": [
        {
          "id": "native-ad-123",
          "type": "native",
          "url": "https://advertiser.com/landing?utm_source=sspserver",
          "fields": {
            "title": "Revolutionary Tech Breakthrough Changes Everything",
            "description": "Scientists discover new technology that could transform how we live and work. Learn more about this groundbreaking innovation.",
            "brandname": "TechCorp",
            "call_to_action": "Learn More",
            "sponsored_label": "Sponsored",
            "category": "Technology",
            "author": "TechCorp Editorial Team",
            "reading_time": "3 min read"
          },
          "assets": [
            {
              "name": "main_image",
              "path": "https://cdn.sspserver.com/assets/native/main-1920x1080.jpg",
              "type": "image",
              "width": 1920,
              "height": 1080,
              "thumbs": [
                {
                  "path": "https://cdn.sspserver.com/assets/native/thumb-400x300.jpg",
                  "type": "image",
                  "width": 400,
                  "height": 300
                },
                {
                  "path": "https://cdn.sspserver.com/assets/native/thumb-200x150.jpg",
                  "type": "image", 
                  "width": 200,
                  "height": 150
                }
              ]
            },
            {
              "name": "brand_logo",
              "path": "https://cdn.sspserver.com/assets/logos/techcorp-100x100.png",
              "type": "image",
              "width": 100,
              "height": 100
            }
          ],
          "tracker": {
            "impressions": [
              "https://track.sspserver.com/imp?id=native-ad-123",
              "https://advertiser.com/track/imp/xyz123"
            ],
            "views": [
              "https://track.sspserver.com/view?id=native-ad-123",
              "https://advertiser.com/track/view/xyz123"
            ],
            "clicks": [
              "https://advertiser.com/track/click/xyz123"
            ]
          },
          "meta": {
            "advertiser": {
              "id": 456,
              "name": "TechCorp Inc.",
              "about_url": "https://techcorp.com/about",
              "privacy_url": "https://techcorp.com/privacy"
            },
            "ad": {
              "id": 789,
              "campaign_id": 101112,
              "description": "Technology innovation campaign"
            }
          }
        }
      ]
    }
  ]
}
```

**Video Native Ad Request:**

```bash
curl -X GET 'https://api.example.com/dynamic?zone=321&type=native&keywords=gaming' \
     -H "Accept: application/json"
```

**Video Native Ad Response:**

```json
{
  "version": "1",
  "groups": [
    {
      "id": "impression-native-video-002",
      "items": [
        {
          "id": "native-video-321",
          "type": "native",
          "url": "https://gaming.example.com/download?ref=native_video",
          "fields": {
            "title": "Epic Adventure Awaits in New Mobile Game",
            "description": "Join millions of players in the most exciting mobile RPG of 2024. Download now and get exclusive starter rewards!",
            "brandname": "GameStudio",
            "call_to_action": "Play Now",
            "sponsored_label": "Sponsored",
            "duration": "30s",
            "category": "Gaming",
            "rating": "4.8★",
            "downloads": "10M+"
          },
          "assets": [
            {
              "name": "video_mp4",
              "path": "https://cdn.sspserver.com/videos/game-trailer-16x9.mp4",
              "type": "video",
              "width": 1280,
              "height": 720
            },
            {
              "name": "video_poster",
              "path": "https://cdn.sspserver.com/videos/game-poster-1280x720.jpg",
              "type": "image",
              "width": 1280,
              "height": 720,
              "thumbs": [
                {
                  "path": "https://cdn.sspserver.com/videos/game-thumb-640x360.jpg",
                  "type": "image",
                  "width": 640,
                  "height": 360
                }
              ]
            },
            {
              "name": "app_icon",
              "path": "https://cdn.sspserver.com/icons/game-icon-512x512.png",
              "type": "image",
              "width": 512,
              "height": 512
            }
          ],
          "tracker": {
            "impressions": [
              "https://track.sspserver.com/imp?id=native-video-321"
            ],
            "views": [
              "https://track.sspserver.com/view?id=native-video-321",
              "https://track.sspserver.com/video-quartile?id=native-video-321&q=25",
              "https://track.sspserver.com/video-quartile?id=native-video-321&q=50",
              "https://track.sspserver.com/video-quartile?id=native-video-321&q=75",
              "https://track.sspserver.com/video-complete?id=native-video-321"
            ]
          }
        }
      ]
    }
  ]
}
```

### Banner Ads

Traditional display banners optimized for website placements and programmatic buying. They support both image and video content through the assets array.

**Image Banner Request:**

```bash
curl -X GET 'https://api.example.com/dynamic?zone=456&type=banner&w=728&h=90&format=display' \
     -H "Accept: application/json"
```

**Image Banner Response:**

```json
{
  "version": "1",
  "groups": [
    {
      "id": "impression-banner-002",
      "items": [
        {
          "id": "banner-ad-456",
          "type": "banner",
          "url": "https://shop.example.com/sale?utm_campaign=summer2024",
          "fields": {
            "headline": "Summer Sale - Up to 50% Off",
            "brandname": "ShopMax",
            "offer": "Limited Time Offer",
            "call_to_action": "Shop Now"
          },
          "assets": [
            {
              "name": "banner_image",
              "path": "https://cdn.sspserver.com/banners/summer-sale-728x90.jpg",
              "type": "image",
              "width": 728,
              "height": 90
            },
            {
              "name": "retina_banner",
              "path": "https://cdn.sspserver.com/banners/summer-sale-1456x180.jpg",
              "type": "image",
              "width": 1456,
              "height": 180
            }
          ],
          "tracker": {
            "impressions": [
              "https://track.sspserver.com/imp?id=banner-ad-456"
            ],
            "views": [
              "https://track.sspserver.com/view?id=banner-ad-456"
            ]
          }
        }
      ]
    }
  ]
}
```

**Video Banner Request:**

```bash
curl -X GET 'https://api.example.com/dynamic?zone=789&type=banner&w=300&h=250' \
     -H "Accept: application/json"
```

**Video Banner Response:**

```json
{
  "version": "1",
  "groups": [
    {
      "id": "impression-video-banner-003",
      "items": [
        {
          "id": "video-banner-789",
          "type": "banner",
          "url": "https://streaming.example.com/signup?promo=video2024",
          "fields": {
            "title": "Stream Your Favorite Shows",
            "brandname": "StreamFlix",
            "duration": "15s",
            "autoplay": true,
            "muted": true,
            "controls": true
          },
          "assets": [
            {
              "name": "video_mp4",
              "path": "https://cdn.sspserver.com/videos/streamflix-promo-300x250.mp4",
              "type": "video",
              "width": 300,
              "height": 250
            },
            {
              "name": "video_webm",
              "path": "https://cdn.sspserver.com/videos/streamflix-promo-300x250.webm",
              "type": "video",
              "width": 300,
              "height": 250
            },
            {
              "name": "poster_image",
              "path": "https://cdn.sspserver.com/videos/streamflix-poster-300x250.jpg",
              "type": "image",
              "width": 300,
              "height": 250
            }
          ],
          "tracker": {
            "impressions": [
              "https://track.sspserver.com/imp?id=video-banner-789"
            ],
            "views": [
              "https://track.sspserver.com/view?id=video-banner-789",
              "https://track.sspserver.com/video-start?id=video-banner-789"
            ],
            "clicks": [
              "https://streaming.example.com/track/click/abc789"
            ]
          }
        }
      ]
    }
  ]
}
```

### Slider Banner Ads

Multi-asset carousel banners for showcasing multiple products or messages.

**Request:**

```bash
curl -X GET 'https://api.example.com/dynamic?zone=654&type=slider_banner&w=320&h=250&count=1' \
     -H "Accept: application/json"
```

**Response:**

```json
{
  "version": "1",
  "groups": [
    {
      "id": "impression-slider-005",
      "items": [
        {
          "id": "slider-banner-654",
          "type": "slider_banner",
          "url": "https://fashion.example.com/collection/spring2024",
          "fields": {
            "title": "Spring Collection 2024",
            "brandname": "FashionForward",
            "slide_count": 4,
            "auto_advance": true,
            "slide_duration": "3s",
            "show_dots": true,
            "show_arrows": true
          },
          "assets": [
            {
              "name": "slide_1",
              "path": "https://cdn.sspserver.com/fashion/spring-slide-1-320x250.jpg",
              "type": "image",
              "width": 320,
              "height": 250
            },
            {
              "name": "slide_2", 
              "path": "https://cdn.sspserver.com/fashion/spring-slide-2-320x250.jpg",
              "type": "image",
              "width": 320,
              "height": 250
            },
            {
              "name": "slide_3",
              "path": "https://cdn.sspserver.com/fashion/spring-slide-3-320x250.jpg",
              "type": "image",
              "width": 320,
              "height": 250
            },
            {
              "name": "slide_4",
              "path": "https://cdn.sspserver.com/fashion/spring-slide-4-320x250.jpg",
              "type": "image",
              "width": 320,
              "height": 250
            }
          ],
          "tracker": {
            "impressions": [
              "https://track.sspserver.com/imp?id=slider-banner-654"
            ],
            "views": [
              "https://track.sspserver.com/view?id=slider-banner-654",
              "https://track.sspserver.com/slide-view?id=slider-banner-654&slide=1",
              "https://track.sspserver.com/slide-view?id=slider-banner-654&slide=2",
              "https://track.sspserver.com/slide-view?id=slider-banner-654&slide=3",
              "https://track.sspserver.com/slide-view?id=slider-banner-654&slide=4"
            ]
          }
        }
      ]
    }
  ]
}
```

### Slider Video Ads

Video carousel presentations with multiple video assets.

**Request:**

```bash
curl -X GET 'https://api.example.com/dynamic?zone=987&type=slider_video&w=400&h=300' \
     -H "Accept: application/json"
```

**Response:**

```json
{
  "version": "1",
  "groups": [
    {
      "id": "impression-slider-video-006",
      "items": [
        {
          "id": "slider-video-987",
          "type": "slider_video",
          "url": "https://travel.example.com/destinations",
          "fields": {
            "title": "Discover Amazing Destinations",
            "brandname": "TravelCorp",
            "slide_count": 3,
            "auto_advance": false,
            "show_controls": true,
            "muted": true
          },
          "assets": [
            {
              "name": "video_slide_1",
              "path": "https://cdn.sspserver.com/travel/destination-1-400x300.mp4",
              "type": "video",
              "width": 400,
              "height": 300
            },
            {
              "name": "video_slide_2",
              "path": "https://cdn.sspserver.com/travel/destination-2-400x300.mp4", 
              "type": "video",
              "width": 400,
              "height": 300
            },
            {
              "name": "video_slide_3",
              "path": "https://cdn.sspserver.com/travel/destination-3-400x300.mp4",
              "type": "video", 
              "width": 400,
              "height": 300
            },
            {
              "name": "poster_slide_1",
              "path": "https://cdn.sspserver.com/travel/poster-1-400x300.jpg",
              "type": "image",
              "width": 400,
              "height": 300
            },
            {
              "name": "poster_slide_2", 
              "path": "https://cdn.sspserver.com/travel/poster-2-400x300.jpg",
              "type": "image",
              "width": 400,
              "height": 300
            },
            {
              "name": "poster_slide_3",
              "path": "https://cdn.sspserver.com/travel/poster-3-400x300.jpg",
              "type": "image",
              "width": 400,
              "height": 300
            }
          ],
          "tracker": {
            "impressions": [
              "https://track.sspserver.com/imp?id=slider-video-987"
            ],
            "views": [
              "https://track.sspserver.com/view?id=slider-video-987",
              "https://track.sspserver.com/video-slide-start?id=slider-video-987&slide=1",
              "https://track.sspserver.com/video-slide-start?id=slider-video-987&slide=2", 
              "https://track.sspserver.com/video-slide-start?id=slider-video-987&slide=3"
            ]
          }
        }
      ]
    }
  ]
}
```

### Proxy Ads

Server-rendered HTML content for iframe delivery and legacy system integration. Proxy ads can use either:

- **HTML Injection** (`content`): Direct HTML content rendered in-page
- **IFrame Delivery** (`content_url`): URL loaded in an iframe

**HTML Injection Example:**

```bash
curl -X GET 'https://api.example.com/dynamic?zone=147&type=proxy&w=300&h=600' \
     -H "Accept: application/json"
```

**Response (HTML Injection):**

```json
{
  "version": "1",
  "groups": [
    {
      "id": "impression-proxy-007",
      "items": [
        {
          "id": "proxy-ad-147",
          "type": "proxy",
          "url": "https://ecommerce.example.com/products/electronics",
          "content_url": "https://proxy.sspserver.com/render/proxy-ad-147",
          "content": "<div class=\"proxy-ad\" style=\"width:300px;height:600px;background:#f0f0f0;\"><div class=\"ad-header\"><h3>Latest Electronics</h3></div><div class=\"product-grid\"><!-- Server-rendered product listings --></div><div class=\"cta-button\"><a href=\"#\">Shop Now</a></div></div>",
          "fields": {
            "title": "Latest Electronics - Up to 40% Off",
            "brandname": "ElectroMax",
            "template": "product_grid_vertical",
            "background_color": "#f0f0f0",
            "text_color": "#333333"
          },
          "tracker": {
            "impressions": [
              "https://track.sspserver.com/imp?id=proxy-ad-147"
            ],
            "views": [
              "https://track.sspserver.com/view?id=proxy-ad-147"
            ]
          }
        }
      ]
    }
  ]
}
```

**IFrame Delivery Example:**

```bash
curl -X GET 'https://api.example.com/dynamic?zone=258&type=proxy&w=320&h=480' \
     -H "Accept: application/json"
```

**Response (IFrame Delivery):**

```json
{
  "version": "1",
  "groups": [
    {
      "id": "impression-proxy-iframe-008",
      "items": [
        {
          "id": "proxy-iframe-258",
          "type": "proxy",
          "url": "https://marketplace.example.com/mobile-deals",
          "content_url": "https://proxy.sspserver.com/iframe/mobile-deals-frame",
          "fields": {
            "title": "Mobile Deals Marketplace",
            "brandname": "MobileMart",
            "template": "mobile_responsive_frame",
            "frame_border": "0",
            "allow_fullscreen": true
          },
          "tracker": {
            "impressions": [
              "https://track.sspserver.com/imp?id=proxy-iframe-258"
            ],
            "views": [
              "https://track.sspserver.com/view?id=proxy-iframe-258"
            ]
          }
        }
      ]
    }
  ]
}
```

## Integration Examples

### JavaScript Native Ad Integration

```javascript
// Fetch native ads and integrate with content feed
async function loadNativeAds() {
  try {
    const response = await fetch('/dynamic?zone=123&type=native&count=3&keywords=technology');
    const data = await response.json();
    
    data.groups.forEach(group => {
      group.items.forEach(item => {
        renderNativeAd(item);
        
        // Fire impression tracking
        item.tracker.impressions.forEach(url => {
          new Image().src = url;
        });
      });
    });
  } catch (error) {
    console.error('Failed to load ads:', error);
  }
}

function renderNativeAd(item) {
  const container = document.createElement('article');
  container.className = 'native-ad';
  container.innerHTML = `
    <div class="sponsored-label">${item.fields.sponsored_label}</div>
    <img src="${item.assets.find(a => a.name === 'main_image').path}" 
         alt="${item.fields.title}" class="ad-image">
    <h3 class="ad-title">${item.fields.title}</h3>
    <p class="ad-description">${item.fields.description}</p>
    <div class="ad-brand">${item.fields.brandname}</div>
    <button class="ad-cta">${item.fields.call_to_action}</button>
  `;
  
  // Add click handler
  container.addEventListener('click', () => {
    // Fire click tracking
    item.tracker.clicks.forEach(url => {
      new Image().src = url;
    });
    window.open(item.url, '_blank');
  });
  
  // Add to content feed
  document.querySelector('#content-feed').appendChild(container);
}
```

### Banner Ad Integration (Assets-based)

```javascript
// Banner ad with assets (no HTML content)
function loadAssetBanner(containerId, zone) {
  fetch(`/dynamic?zone=${zone}&type=banner&w=300&h=250`)
    .then(response => response.json())
    .then(data => {
      const item = data.groups[0].items[0];
      const container = document.getElementById(containerId);
      
      // Banner ads use assets for media content
      const bannerAsset = item.assets.find(a => a.name === 'main_banner') || item.assets[0];
      
      const bannerDiv = document.createElement('div');
      bannerDiv.className = 'asset-banner';
      bannerDiv.innerHTML = `
        <img src="${bannerAsset.path}" 
             alt="${item.fields.title || 'Advertisement'}"
             style="width: 100%; height: auto; cursor: pointer;">
      `;
      
      // Add click handler
      bannerDiv.addEventListener('click', () => {
        item.tracker.clicks.forEach(url => {
          new Image().src = url;
        });
        window.open(item.url, '_blank');
      });
      
      container.appendChild(bannerDiv);
      
      // Fire impression tracking
      item.tracker.impressions.forEach(url => {
        new Image().src = url;
      });
    });
}
```

### Video Banner Integration

```javascript
// Video banner with autoplay and tracking (using banner type with video assets)
function loadVideoBanner(containerId, zone) {
  fetch(`/dynamic?zone=${zone}&type=banner&w=300&h=250`)
    .then(response => response.json())
    .then(data => {
      const item = data.groups[0].items[0];
      const container = document.getElementById(containerId);
      
      // Check if this banner has video assets
      const videoAsset = item.assets.find(a => a.type === 'video');
      if (videoAsset) {
        const video = document.createElement('video');
        video.src = videoAsset.path;
        video.poster = item.assets.find(a => a.name === 'poster_image')?.path;
        video.autoplay = item.fields.autoplay;
        video.muted = item.fields.muted;
        video.controls = item.fields.controls;
        video.style.width = '100%';
        video.style.height = '100%';
        
        // Video event tracking
        video.addEventListener('play', () => {
          item.tracker.views.forEach(url => {
            if (url.includes('video-start')) {
              new Image().src = url;
            }
          });
        });
        
        video.addEventListener('click', () => {
          item.tracker.clicks.forEach(url => {
            new Image().src = url;
          });
          window.open(item.url, '_blank');
        });
        
        container.appendChild(video);
      }
      
      // Fire impression tracking
      item.tracker.impressions.forEach(url => {
        new Image().src = url;
      });
    });
}
```

### Slider Banner Implementation

```javascript
// Multi-slide banner carousel
class SliderBanner {
  constructor(containerId, zone) {
    this.container = document.getElementById(containerId);
    this.currentSlide = 0;
    this.loadAd(zone);
  }
  
  async loadAd(zone) {
    try {
      const response = await fetch(`/dynamic?zone=${zone}&type=slider_banner`);
      const data = await response.json();
      this.item = data.groups[0].items[0];
      this.render();
      
      // Fire impression tracking
      this.item.tracker.impressions.forEach(url => {
        new Image().src = url;
      });
    } catch (error) {
      console.error('Failed to load slider banner:', error);
    }
  }
  
  render() {
    const slides = this.item.assets.filter(asset => asset.name.startsWith('slide_'));
    
    this.container.innerHTML = `
      <div class="slider-container">
        <div class="slides">
          ${slides.map((slide, index) => `
            <img src="${slide.path}" 
                 class="slide ${index === 0 ? 'active' : ''}"
                 alt="Slide ${index + 1}">
          `).join('')}
        </div>
        <div class="slider-controls">
          <button class="prev" onclick="this.previousSlide()">‹</button>
          <div class="dots">
            ${slides.map((_, index) => `
              <span class="dot ${index === 0 ? 'active' : ''}" 
                    onclick="this.goToSlide(${index})"></span>
            `).join('')}
          </div>
          <button class="next" onclick="this.nextSlide()">›</button>
        </div>
      </div>
    `;
    
    // Auto-advance if enabled
    if (this.item.fields.auto_advance) {
      this.startAutoAdvance();
    }
    
    // Add click handler
    this.container.addEventListener('click', () => {
      this.item.tracker.clicks.forEach(url => {
        new Image().src = url;
      });
      window.open(this.item.url, '_blank');
    });
  }
  
  nextSlide() {
    const slides = this.container.querySelectorAll('.slide');
    const dots = this.container.querySelectorAll('.dot');
    
    slides[this.currentSlide].classList.remove('active');
    dots[this.currentSlide].classList.remove('active');
    
    this.currentSlide = (this.currentSlide + 1) % slides.length;
    
    slides[this.currentSlide].classList.add('active');
    dots[this.currentSlide].classList.add('active');
    
    // Fire slide view tracking
    this.fireSlideTracking(this.currentSlide + 1);
  }
  
  fireSlideTracking(slideNumber) {
    this.item.tracker.views.forEach(url => {
      if (url.includes(`slide=${slideNumber}`)) {
        new Image().src = url;
      }
    });
  }
}
```

### Mobile App Integration (React Native)

```jsx
// React Native component for dynamic ads
import React, { useState, useEffect } from 'react';
import { View, Text, Image, TouchableOpacity, Video } from 'react-native';

const DynamicAd = ({ zone, adType, width, height }) => {
  const [adData, setAdData] = useState(null);
  
  useEffect(() => {
    loadAd();
  }, [zone, adType]);
  
  const loadAd = async () => {
    try {
      const response = await fetch(
        `/dynamic?zone=${zone}&type=${adType}&w=${width}&h=${height}`
      );
      const data = await response.json();
      setAdData(data.groups[0].items[0]);
      
      // Fire impression tracking
      if (data.groups[0].items[0]) {
        fireTracking(data.groups[0].items[0].tracker.impressions);
      }
    } catch (error) {
      console.error('Ad loading failed:', error);
    }
  };
  
  const fireTracking = (urls) => {
    urls.forEach(url => {
      fetch(url, { method: 'GET' }).catch(() => {});
    });
  };
  
  const handleAdClick = () => {
    if (adData) {
      fireTracking(adData.tracker.clicks);
      // Open ad URL in browser
      Linking.openURL(adData.url);
    }
  };
  
  if (!adData) return <View style={{width, height}} />;
  
  switch (adType) {
    case 'native':
      return (
        <TouchableOpacity onPress={handleAdClick} style={styles.nativeAd}>
          <Text style={styles.sponsoredLabel}>{adData.fields.sponsored_label}</Text>
          <Image 
            source={{uri: adData.assets.find(a => a.name === 'main_image').path}}
            style={styles.adImage}
          />
          <Text style={styles.adTitle}>{adData.fields.title}</Text>
          <Text style={styles.adDescription}>{adData.fields.description}</Text>
          <Text style={styles.brandName}>{adData.fields.brandname}</Text>
        </TouchableOpacity>
      );
      
    case 'banner':
      // Check if banner has video assets
      const videoAsset = adData.assets.find(a => a.type === 'video');
      if (videoAsset) {
        return (
          <TouchableOpacity onPress={handleAdClick}>
            <Video
              source={{uri: videoAsset.path}}
              poster={adData.assets.find(a => a.name === 'poster_image')?.path}
              shouldPlay={adData.fields.autoplay}
              isMuted={adData.fields.muted}
              resizeMode="cover"
              style={{width, height}}
            />
          </TouchableOpacity>
        );
      } else {
        // Image banner
        return (
          <TouchableOpacity onPress={handleAdClick}>
            <Image 
              source={{uri: adData.assets[0].path}}
              style={{width, height}}
            />
          </TouchableOpacity>
        );
      }
      
    default:
      return (
        <TouchableOpacity onPress={handleAdClick}>
          <Image 
            source={{uri: adData.assets[0].path}}
            style={{width, height}}
          />
        </TouchableOpacity>
      );
  }
};

const styles = {
  nativeAd: {
    padding: 12,
    backgroundColor: '#f9f9f9',
    borderRadius: 8,
    margin: 8
  },
  sponsoredLabel: {
    fontSize: 10,
    color: '#666',
    marginBottom: 8
  },
  adImage: {
    width: '100%',
    height: 200,
    borderRadius: 4,
    marginBottom: 8
  },
  adTitle: {
    fontSize: 16,
    fontWeight: 'bold',
    marginBottom: 4
  },
  adDescription: {
    fontSize: 14,
    color: '#666',
    marginBottom: 8
  },
  brandName: {
    fontSize: 12,
    color: '#888'
  }
};
```

## Request Parameters

### Core Parameters

| Parameter  | Type     | Description | Example |
|------------|----------|-------------|---------|
| `zone`     | `int`    | **Required.** Zone/placement identifier | `zone=123` |
| `type`     | `string` | Ad format type | `type=native,banner` |
| `count`    | `int`    | Number of ads requested (1-10) | `count=3` |

### Sizing Parameters

| Parameter  | Type     | Description | Example |
|------------|----------|-------------|---------|
| `w`        | `int`    | Maximum desired width in pixels | `w=300` |
| `h`        | `int`    | Maximum desired height in pixels | `h=250` |
| `mw`       | `int`    | Minimum width constraint | `mw=250` |
| `mh`       | `int`    | Minimum height constraint | `mh=200` |
| `fmt`      | `string` | Size format shorthand | `fmt=300x250` |
| `width`    | `int`    | Alias for `mw` | `width=250` |
| `height`   | `int`    | Alias for `mh` | `height=200` |

### Targeting Parameters

| Parameter  | Type     | Description | Example |
|------------|----------|-------------|---------|
| `keywords` | `string` | Comma-separated targeting keywords | `keywords=tech,mobile,apps` |
| `category` | `string` | Content category for targeting | `category=technology` |
| `x`        | `int`    | X coordinate for geo/position targeting | `x=100` |
| `y`        | `int`    | Y coordinate for geo/position targeting | `y=200` |

### Format Control Parameters

| Parameter  | Type     | Description | Values |
|------------|----------|-------------|--------|
| `format`   | `string` | Response format | `json` (default), `jsonp` |
| `callback` | `string` | JSONP callback function name | `callback=handleAds` |
| `debug`    | `bool`   | Enable debug information | `debug=true` |

### Tracking Parameters

| Parameter  | Type     | Description | Example |
|------------|----------|-------------|---------|
| `subid1`   | `string` | Primary custom tracking ID | `subid1=user123` |
| `subid2`   | `string` | Secondary tracking ID | `subid2=campaign456` |
| `subid3`   | `string` | Tertiary tracking ID | `subid3=source789` |
| `subid4`   | `string` | Additional tracking ID | `subid4=medium_cpc` |
| `subid5`   | `string` | Additional tracking ID | `subid5=creative_a` |

### Ad Type Values

| Type | Description | Best For |
|------|-------------|----------|
| `native` | Content-style ads with title, description, images, and videos | Editorial integration, sponsored content, video content |
| `banner` | Traditional display rectangles with images, videos, and HTML | Website monetization, programmatic, rich media |
| `slider_banner` | Multi-image carousel banners | Product showcases, brand stories |
| `slider_video` | Video carousel presentations | Entertainment, travel, lifestyle |
| `proxy` | Server-rendered HTML content with `content` or `content_url` | Legacy integration, custom templates |

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
