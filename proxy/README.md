# Proxy Endpoint

The `proxy` package provides server-side HTML template rendering for embedded ad delivery, optimized for banner placement, iframe integration, and legacy system compatibility.

## Table of Contents

- [Overview](#overview)
- [Use Cases](#use-cases)
- [Template System](#template-system)
- [Request Parameters](#request-parameters)
- [Response Format](#response-format)
- [API Examples](#api-examples)
- [Integration Examples](#integration-examples)
- [Build Configuration](#build-configuration)
- [JavaScript Integration](#javascript-integration)

## Overview

The proxy endpoint (`/proxy`) renders HTML templates server-side and delivers complete HTML documents for ad display. This approach is ideal for:

- **Banner Ad Placement**: Traditional display advertising in HTML format
- **IFrame Integration**: Content delivered via iframe embedding
- **Legacy System Compatibility**: HTML-based ad delivery for older systems
- **Direct HTML Embedding**: Server-rendered content without client-side processing

**Key Features:**

- Server-side HTML template rendering
- Embedded JavaScript SDK integration
- Dynamic content loading via JSONP
- Responsive ad rendering
- Legacy browser compatibility
- Build-time template compilation

## Use Cases

### Traditional Banner Advertising

Perfect for standard display banner placements where HTML content is directly embedded in web pages.

### IFrame-Based Delivery

Ideal for serving ads in sandboxed iframe environments, providing security isolation between ad content and parent page.

### Legacy System Integration

Seamless integration with older content management systems that require direct HTML insertion without JavaScript complexity.

### Content Management System Plugins

Server-rendered HTML that can be easily integrated into CMS templates, WordPress plugins, or static site generators.

## Template System

The proxy endpoint uses a template-based rendering system with the following components:

### Core Templates

- **`ad_base.qtpl`**: Base HTML structure and meta tags
- **`ad_dinamic_proxy.qtpl`**: Dynamic proxy banner rendering
- **`ad_native.qtpl`**: Native ad styling and layout

### Template Features

- **Responsive Design**: CSS media queries for different screen sizes
- **Loading States**: Preloader animations during content fetch
- **Error Handling**: Graceful fallback for failed ad requests
- **SEO Friendly**: Proper HTML structure and meta tags

## Request Parameters

### Core Parameters

| Parameter | Type | Description | Example |
|-----------|------|-------------|---------|
| `zone` | `int` | **Required.** Zone/placement identifier | `zone=123` |
| `w` | `int` | Banner width in pixels | `w=728` |
| `h` | `int` | Banner height in pixels | `h=90` |

### Positioning Parameters

| Parameter | Type | Description | Example |
|-----------|------|-------------|---------|
| `x` | `int` | X coordinate for ad positioning | `x=100` |
| `y` | `int` | Y coordinate for ad positioning | `y=200` |

### Targeting Parameters

| Parameter | Type | Description | Example |
|-----------|------|-------------|---------|
| `keywords` | `string` | Comma-separated targeting keywords | `keywords=tech,mobile` |
| `type` | `string` | Ad format type | `type=banner` |
| `format` | `string` | Format specification | `format=display` |

### Tracking Parameters

| Parameter | Type | Description | Aliases |
|-----------|------|-------------|---------|
| `subid1` | `string` | Primary tracking identifier | `subid`, `s1` |
| `subid2` | `string` | Secondary tracking identifier | `s2` |
| `subid3` | `string` | Tertiary tracking identifier | `s3` |
| `subid4` | `string` | Quaternary tracking identifier | `s4` |
| `subid5` | `string` | Quinary tracking identifier | `s5` |

## Response Format

The proxy endpoint returns complete HTML documents with embedded ads:

```html
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Ad Content</title>
    <style>
        /* Responsive ad styling */
        .ad-container {
            position: relative;
            max-width: 100%;
            overflow: hidden;
        }
        
        /* Loading animation */
        .loading {
            display: flex;
            justify-content: center;
            align-items: center;
            min-height: 100px;
        }
        
        /* Responsive breakpoints */
        @media (max-width: 768px) {
            .ad-container {
                padding: 10px;
            }
        }
    </style>
</head>
<body>
    <div id="loadingBlock" class="loading">
        <div class="spinner">Loading...</div>
    </div>
    
    <ins id="element_123"></ins>
    
    <script src="https://cdn.sspserver.com/embedded.js"></script>
    <script>
        (function() {
            new EmbeddedAd({
                element: "element_123",
                zone_id: 123,
                JSONPLink: '//api.sspserver.com/dynamic/123?format=jsonp&'
            }).on('render', function() {
                var loader = document.getElementById('loadingBlock');
                if (loader) {
                    loader.parentElement.removeChild(loader);
                }
            }).on('error', function(err) {
                console.log('Ad loading error:', err);
            }).render();
        })();
    </script>
</body>
</html>
```

### Response Structure

- **HTML Document**: Complete HTML5 document structure
- **Responsive CSS**: Media queries for different screen sizes
- **Loading States**: Visual feedback during ad loading
- **JavaScript Integration**: Embedded SDK for dynamic content
- **Error Handling**: Graceful fallback for loading failures

## API Examples

### Basic Banner Request

```bash
curl 'https://api.example.com/proxy?zone=123&w=728&h=90' \
     -H "Accept: text/html"
```

### Mobile Banner Request

```bash
curl 'https://api.example.com/proxy?zone=456&w=320&h=50&type=mobile' \
     -H "Accept: text/html"
```

### Leaderboard Banner Request

```bash
curl 'https://api.example.com/proxy?zone=789&w=728&h=90&keywords=technology' \
     -H "Accept: text/html"
```

### Targeted Banner with Tracking

```bash
curl 'https://api.example.com/proxy?zone=123&w=300&h=250&subid1=user456&keywords=finance' \
     -H "Accept: text/html"
```

## Integration Examples

### Direct HTML Embedding

```html
<!-- Direct iframe integration -->
<iframe src="https://api.example.com/proxy?zone=123&w=728&h=90" 
        width="728" 
        height="90" 
        frameborder="0" 
        scrolling="no">
</iframe>
```

### WordPress Integration

```php
<?php
// WordPress shortcode for proxy ads
function sspserver_proxy_ad_shortcode($atts) {
    $atts = shortcode_atts([
        'zone' => '',
        'width' => '300',
        'height' => '250',
        'keywords' => ''
    ], $atts);
    
    if (empty($atts['zone'])) {
        return '';
    }
    
    $proxy_url = add_query_arg([
        'zone' => $atts['zone'],
        'w' => $atts['width'],
        'h' => $atts['height'],
        'keywords' => $atts['keywords'],
        'subid1' => get_current_user_id()
    ], 'https://api.example.com/proxy');
    
    return sprintf(
        '<iframe src="%s" width="%s" height="%s" frameborder="0" scrolling="no" style="border:none;"></iframe>',
        esc_url($proxy_url),
        esc_attr($atts['width']),
        esc_attr($atts['height'])
    );
}

add_shortcode('sspserver_proxy', 'sspserver_proxy_ad_shortcode');

// Usage: [sspserver_proxy zone="123" width="728" height="90" keywords="technology"]
?>
```

### JavaScript Dynamic Loading

```javascript
// Dynamically load proxy ads
class ProxyAdLoader {
    constructor(containerId) {
        this.container = document.getElementById(containerId);
    }
    
    loadAd(zoneId, width, height, options = {}) {
        const iframe = document.createElement('iframe');
        
        const params = new URLSearchParams({
            zone: zoneId,
            w: width,
            h: height,
            ...options
        });
        
        iframe.src = `https://api.example.com/proxy?${params}`;
        iframe.width = width;
        iframe.height = height;
        iframe.frameBorder = '0';
        iframe.scrolling = 'no';
        iframe.style.border = 'none';
        
        // Handle load events
        iframe.onload = () => {
            console.log('Proxy ad loaded successfully');
        };
        
        iframe.onerror = () => {
            console.error('Failed to load proxy ad');
            this.showFallback();
        };
        
        this.container.appendChild(iframe);
        return iframe;
    }
    
    showFallback() {
        this.container.innerHTML = '<div style="background:#f0f0f0;padding:20px;text-align:center;">Advertisement</div>';
    }
}

// Usage
const adLoader = new ProxyAdLoader('ad-container');
adLoader.loadAd(123, 728, 90, { keywords: 'technology', subid1: 'user123' });
```

### React Component Integration

```jsx
import React, { useEffect, useRef } from 'react';

const ProxyAd = ({ zoneId, width, height, keywords, trackingId }) => {
    const iframeRef = useRef(null);
    
    useEffect(() => {
        const params = new URLSearchParams({
            zone: zoneId,
            w: width,
            h: height,
            ...(keywords && { keywords }),
            ...(trackingId && { subid1: trackingId })
        });
        
        const proxyUrl = `https://api.example.com/proxy?${params}`;
        
        if (iframeRef.current) {
            iframeRef.current.src = proxyUrl;
        }
    }, [zoneId, width, height, keywords, trackingId]);
    
    return (
        <iframe
            ref={iframeRef}
            width={width}
            height={height}
            frameBorder="0"
            scrolling="no"
            style={{ border: 'none', display: 'block' }}
            title={`Ad Zone ${zoneId}`}
        />
    );
};

// Usage
<ProxyAd 
    zoneId={123} 
    width={300} 
    height={250} 
    keywords="technology,mobile" 
    trackingId="user456" 
/>
```

### Server-Side Integration (Node.js)

```javascript
const express = require('express');
const app = express();

// Middleware to serve proxy ads
app.get('/ad/:zoneId', (req, res) => {
    const { zoneId } = req.params;
    const { w = 300, h = 250, keywords } = req.query;
    
    const params = new URLSearchParams({
        zone: zoneId,
        w,
        h,
        ...(keywords && { keywords }),
        subid1: req.sessionID
    });
    
    const proxyUrl = `https://api.example.com/proxy?${params}`;
    
    // Return iframe HTML
    res.send(`
        <iframe src="${proxyUrl}" 
                width="${w}" 
                height="${h}" 
                frameborder="0" 
                scrolling="no" 
                style="border:none;">
        </iframe>
    `);
});

// Usage: GET /ad/123?w=728&h=90&keywords=technology
```

## Build Configuration

The proxy endpoint requires build-time configuration to enable HTML template support:

### Go Build Tags

```bash
# Build with HTML template support
go build -tags htmltemplates

# Build without HTML template support (proxy endpoint disabled)
go build
```

### Build Tags

The proxy endpoint is conditionally compiled based on build tags:

- **With `htmltemplates` tag**: Full proxy functionality enabled
- **Without `htmltemplates` tag**: Proxy endpoint returns `nil` (disabled)

### Template Compilation

Templates are compiled at build time using the qtpl template system:

```bash
# Generate template Go code
qtpl -f templates/

# Build with templates
go build -tags htmltemplates
```

## JavaScript Integration

The proxy endpoint generates HTML that integrates with the JavaScript SDK:

### EmbeddedAd Class

The rendered HTML includes JavaScript that creates an `EmbeddedAd` instance:

```javascript
new EmbeddedAd({
    element: "element_123",          // Target DOM element ID
    zone_id: 123,                   // Zone identifier
    JSONPLink: '/dynamic/123?format=jsonp&'  // JSONP endpoint for dynamic content
})
.on('render', function() {
    // Ad successfully rendered
    hideLoader();
})
.on('error', function(err) {
    // Handle rendering errors
    console.error('Ad render error:', err);
    showFallback();
})
.render();
```

### Event Handling

The JavaScript integration provides event hooks:

- **`render`**: Fired when ad content is successfully displayed
- **`error`**: Fired when ad loading or rendering fails
- **Custom events**: Additional events based on ad interaction

### Responsive Behavior

The generated HTML includes responsive CSS and JavaScript to adapt to different screen sizes:

```javascript
// Responsive ad handling
function adaptAdSize() {
    const container = document.getElementById('element_123');
    const containerWidth = container.offsetWidth;
    
    if (containerWidth < 728) {
        // Switch to mobile ad format
        loadMobileAd();
    } else {
        // Load desktop ad format
        loadDesktopAd();
    }
}

window.addEventListener('resize', adaptAdSize);
```

## Performance Considerations

### Template Caching

Templates are compiled at build time, providing optimal runtime performance:

- **No runtime compilation**: Templates pre-compiled to Go code
- **Memory efficiency**: Compiled templates loaded once at startup
- **Fast rendering**: Minimal CPU overhead for HTML generation

### Delivery Optimization

```javascript
// Optimize iframe loading
function optimizeProxyAd(iframe) {
    // Lazy loading
    iframe.loading = 'lazy';
    
    // Intersection Observer for viewport detection
    const observer = new IntersectionObserver((entries) => {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                // Load ad only when visible
                iframe.src = iframe.dataset.src;
                observer.unobserve(iframe);
            }
        });
    });
    
    observer.observe(iframe);
}
```

### Content Security Policy

Configure CSP headers for secure proxy ad delivery:

```html
<meta http-equiv="Content-Security-Policy" 
      content="default-src 'self'; 
               script-src 'self' https://cdn.sspserver.com; 
               img-src 'self' https:; 
               connect-src https://api.sspserver.com;">
```

## Error Handling and Fallbacks

### JavaScript Error Handling

```javascript
// Robust error handling for proxy ads
class RobustProxyAd {
    constructor(containerId, zoneId, options = {}) {
        this.container = document.getElementById(containerId);
        this.zoneId = zoneId;
        this.options = options;
        this.retryCount = 0;
        this.maxRetries = 3;
    }
    
    load() {
        const iframe = this.createIframe();
        
        iframe.onload = () => {
            console.log('Proxy ad loaded successfully');
        };
        
        iframe.onerror = () => {
            this.handleError();
        };
        
        // Timeout handling
        setTimeout(() => {
            if (!iframe.contentDocument || !iframe.contentDocument.body.innerHTML) {
                this.handleError();
            }
        }, 5000);
        
        this.container.appendChild(iframe);
    }
    
    createIframe() {
        const iframe = document.createElement('iframe');
        const params = new URLSearchParams({
            zone: this.zoneId,
            ...this.options
        });
        
        iframe.src = `https://api.example.com/proxy?${params}`;
        iframe.style.cssText = 'border:none;width:100%;height:100%;';
        
        return iframe;
    }
    
    handleError() {
        if (this.retryCount < this.maxRetries) {
            this.retryCount++;
            console.log(`Retrying proxy ad load (${this.retryCount}/${this.maxRetries})`);
            setTimeout(() => this.load(), 1000 * this.retryCount);
        } else {
            this.showFallback();
        }
    }
    
    showFallback() {
        this.container.innerHTML = `
            <div style="
                background: #f5f5f5;
                border: 1px solid #ddd;
                padding: 20px;
                text-align: center;
                color: #666;
            ">
                Advertisement space
            </div>
        `;
    }
}
```

## Best Practices

### 1. Iframe Sandboxing

Always use appropriate iframe sandbox attributes for security:

```html
<iframe src="proxy-ad-url"
        sandbox="allow-scripts allow-same-origin allow-popups allow-forms"
        width="300" 
        height="250">
</iframe>
```

### 2. Responsive Design

Implement responsive ad containers:

```css
.ad-container {
    position: relative;
    width: 100%;
    max-width: 728px;
    margin: 0 auto;
}

.ad-container iframe {
    width: 100%;
    height: auto;
    min-height: 90px;
}

@media (max-width: 768px) {
    .ad-container {
        max-width: 320px;
    }
}
```

### 3. Loading States

Provide visual feedback during ad loading:

```javascript
function showAdLoader(containerId) {
    const container = document.getElementById(containerId);
    container.innerHTML = `
        <div class="ad-loader">
            <div class="spinner"></div>
            <p>Loading advertisement...</p>
        </div>
    `;
}
```

### 4. Accessibility

Ensure ads are accessible:

```html
<iframe src="proxy-ad-url" 
        title="Advertisement" 
        role="img" 
        aria-label="Sponsored content">
</iframe>
```

## Summary

The proxy endpoint provides a robust, server-side rendered solution for HTML-based ad delivery. Its template system, JavaScript integration, and iframe compatibility make it ideal for traditional banner advertising, legacy system integration, and environments requiring server-rendered content. The build-time compilation and responsive design ensure optimal performance and user experience across different devices and platforms.
