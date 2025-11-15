# Direct Endpoint

The `direct` package provides immediate HTTP redirect functionality for ad serving, optimized for popunder campaigns, direct action advertising, and click-through monetization.

## Table of Contents

- [Overview](#overview)
- [Use Cases](#use-cases)
- [Request Parameters](#request-parameters)
- [Response Types](#response-types)
- [API Examples](#api-examples)
- [Debug Mode](#debug-mode)
- [Error Handling](#error-handling)
- [Integration Examples](#integration-examples)

## Overview

The direct endpoint (`/direct`) handles immediate redirects with minimal latency, making it ideal for:

- **Popunder Campaigns**: Direct browser window/tab redirection
- **Click-Through Advertising**: Immediate navigation to advertiser landing pages
- **Direct Action Campaigns**: One-click conversions and app downloads
- **Alternative Link Handling**: Fallback content when no ads are available

**Key Features:**

- Immediate HTTP 302 redirects
- Superfailover URL support for monetization continuity
- Alternative link handling with custom headers
- Comprehensive debug mode with JSON responses
- Event tracking for clicks and impressions
- Support for multiple tracking identifiers

## Use Cases

### Popunder Advertising

Perfect for popunder campaigns where users are redirected to advertiser content in a new window or tab.

### Direct Response Marketing

Ideal for campaigns requiring immediate user action, such as app downloads, sign-ups, or purchases.

### Content Monetization

Provides seamless monetization for content sites with minimal integration complexity.

### Legacy System Integration

Simple HTTP redirect mechanism compatible with any system that can handle HTTP redirects.

## Request Parameters

### Core Parameters

| Parameter | Type | Description | Example |
|-----------|------|-------------|---------|
| `zone` | `int` | **Required.** Zone/placement identifier | `zone=123` |
| `debug` | `bool` | Enable debug mode | `debug=true` |
| `noredirect` | `bool` | Return JSON response instead of redirect (debug mode only) | `noredirect=true` |

### Targeting Parameters

| Parameter | Type | Description | Example |
|-----------|------|-------------|---------|
| `keywords` | `string` | Comma-separated targeting keywords | `keywords=tech,mobile` |
| `x` | `int` | X coordinate for positioning | `x=100` |
| `y` | `int` | Y coordinate for positioning | `y=200` |

### Tracking Parameters

| Parameter | Type | Description | Aliases |
|-----------|------|-------------|---------|
| `subid1` | `string` | Primary tracking identifier | `subid`, `s1` |
| `subid2` | `string` | Secondary tracking identifier | `s2` |
| `subid3` | `string` | Tertiary tracking identifier | `s3` |
| `subid4` | `string` | Quaternary tracking identifier | `s4` |
| `subid5` | `string` | Quinary tracking identifier | `s5` |

## Response Types

### Normal Operation Responses

#### Successful Ad Delivery

```http
HTTP/1.1 302 Found
Location: https://advertiser.com/landing?campaign=abc123
```

#### Alternative Link Response

```http
HTTP/1.1 302 Found
Location: https://alternative.com/content
X-Status-Alternative: 1
```

#### No Ad Available (Superfailover)

```http
HTTP/1.1 302 Found
Location: https://superfailover.com/monetization
X-Status-Failover: 1
```

### Debug Mode Response

When using `debug=true&noredirect=true`, returns JSON instead of redirect:

```json
{
  "id": "ad-unit-12345",
  "zone_id": 123,
  "auction_id": "auction-uuid-678",
  "impression_id": "impression-uuid-901",
  "is_alternative_link": false,
  "link": "https://advertiser.com/click?id=abc123",
  "superfailover": "https://superfailover.com/default",
  "error": null,
  "is_empty": false
}
```

### Response Fields Explanation

| Field | Type | Description |
|-------|------|-------------|
| `id` | `string` | Unique ad unit identifier |
| `zone_id` | `int` | Zone/placement ID from request |
| `auction_id` | `string` | Internal auction identifier |
| `impression_id` | `string` | Unique impression identifier |
| `is_alternative_link` | `bool` | Whether response is alternative content |
| `link` | `string` | Target URL for redirect |
| `superfailover` | `string` | Fallback URL when no ads available |
| `error` | `string` | Error message if applicable |
| `is_empty` | `bool` | Whether no ads were returned |

## API Examples

### Basic Direct Request

```bash
curl -I 'https://api.example.com/direct?zone=123'
# Response: HTTP 302 redirect to ad URL
```

### Request with Tracking

```bash
curl -I 'https://api.example.com/direct?zone=123&subid1=user456&subid2=campaign789'
# Response: HTTP 302 redirect with tracking parameters
```

### Debug Request

```bash
curl -X GET 'https://api.example.com/direct?zone=123&debug=true&noredirect=true' \
     -H "Accept: application/json"
```

**Debug Response:**

```json
{
  "id": "direct-ad-123",
  "zone_id": 123,
  "auction_id": "auction-abc-123",
  "impression_id": "imp-def-456",
  "is_alternative_link": false,
  "link": "https://advertiser.com/landing?utm_source=direct&utm_campaign=mobile",
  "superfailover": "https://fallback.com/monetize",
  "error": null,
  "is_empty": false
}
```

### Alternative Link Example

```bash
curl -I 'https://api.example.com/direct?zone=456'
# Response: HTTP 302 with X-Status-Alternative: 1 header
```

### No Ads Available Example

```bash
curl -I 'https://api.example.com/direct?zone=999'
# Response: HTTP 302 to superfailover with X-Status-Failover: 1 header
```

## Debug Mode

Debug mode provides comprehensive information about the request processing:

### Enabling Debug Mode

Add `debug=true&noredirect=true` to prevent redirects and receive JSON responses:

```bash
curl 'https://api.example.com/direct?zone=123&debug=true&noredirect=true'
```

### Debug Information Includes

- **Request Details**: Zone ID, auction ID, impression ID
- **Response Status**: Whether ad was found, alternative link used, or fallback triggered
- **URLs**: Target link and superfailover configuration
- **Error Information**: Any processing errors encountered
- **Tracking Data**: Internal identifiers for debugging

**Note:** Debug mode should only be used in development/testing environments.

## Error Handling

### Common Error Scenarios

#### Multiple Direct Responses

```json
{
  "error": "direct: multiple direct responses not supported",
  "is_empty": true
}
```

#### Invalid Response Type

```json
{
  "error": "direct: invalid response type", 
  "is_empty": true
}
```

#### No Superfailover Configured

```http
HTTP/1.1 200 OK
Content-Type: text/plain

Please add superfailover link
```

### Error Response Headers

- **`X-Status-Alternative: 1`**: Indicates alternative content was served
- **`X-Status-Failover: 1`**: Indicates superfailover was triggered

## Integration Examples

### JavaScript Popunder Integration

```javascript
// Create popunder window with direct endpoint
function createPopunder(zoneId, trackingId) {
  const directUrl = `/direct?zone=${zoneId}&subid1=${trackingId}`;
  
  // Open in new window/tab
  const popup = window.open(directUrl, '_blank', 
    'width=800,height=600,scrollbars=yes,resizable=yes');
  
  // Optional: Focus back to original window
  setTimeout(() => {
    window.focus();
  }, 100);
  
  return popup;
}

// Usage
createPopunder(123, 'user-session-456');
```

### PHP Integration

```php
<?php
// Server-side redirect handling
function redirectToAd($zoneId, $subid1 = null) {
    $params = ['zone' => $zoneId];
    if ($subid1) {
        $params['subid1'] = $subid1;
    }
    
    $directUrl = 'https://api.example.com/direct?' . http_build_query($params);
    
    // Redirect user
    header("Location: $directUrl", true, 302);
    exit;
}

// Usage
redirectToAd(123, 'user123');
?>
```

### Python Integration

```python
import requests
from flask import redirect

def get_ad_redirect(zone_id, tracking_params=None):
    """Get direct ad redirect URL"""
    params = {'zone': zone_id}
    if tracking_params:
        params.update(tracking_params)
    
    # For debugging, you can check the target URL first
    params.update({'debug': 'true', 'noredirect': 'true'})
    
    response = requests.get('https://api.example.com/direct', params=params)
    data = response.json()
    
    return data.get('link') or data.get('superfailover')

# Flask route example
@app.route('/ad/<int:zone_id>')
def serve_ad(zone_id):
    return redirect(f'https://api.example.com/direct?zone={zone_id}')
```

### Mobile App Integration (React Native)

```javascript
import { Linking } from 'react-native';

class DirectAdService {
  static async openAd(zoneId, trackingParams = {}) {
    const params = new URLSearchParams({
      zone: zoneId,
      ...trackingParams
    });
    
    const directUrl = `https://api.example.com/direct?${params}`;
    
    try {
      // Open in device browser
      await Linking.openURL(directUrl);
    } catch (error) {
      console.error('Failed to open ad:', error);
    }
  }
  
  static async getAdUrl(zoneId, trackingParams = {}) {
    const params = new URLSearchParams({
      zone: zoneId,
      debug: 'true',
      noredirect: 'true',
      ...trackingParams
    });
    
    try {
      const response = await fetch(`https://api.example.com/direct?${params}`);
      const data = await response.json();
      return data.link || data.superfailover;
    } catch (error) {
      console.error('Failed to get ad URL:', error);
      return null;
    }
  }
}

// Usage
DirectAdService.openAd(123, { subid1: 'mobile-user-456' });
```

### WordPress Plugin Integration

```php
<?php
// WordPress plugin integration
function sspserver_direct_ad_shortcode($atts) {
    $atts = shortcode_atts([
        'zone' => '',
        'text' => 'Click Here',
        'target' => '_blank',
        'class' => 'direct-ad-link'
    ], $atts);
    
    if (empty($atts['zone'])) {
        return '';
    }
    
    $direct_url = add_query_arg([
        'zone' => $atts['zone'],
        'subid1' => get_current_user_id(),
        'subid2' => get_the_ID()
    ], 'https://api.example.com/direct');
    
    return sprintf(
        '<a href="%s" target="%s" class="%s">%s</a>',
        esc_url($direct_url),
        esc_attr($atts['target']),
        esc_attr($atts['class']),
        esc_html($atts['text'])
    );
}

add_shortcode('sspserver_direct', 'sspserver_direct_ad_shortcode');

// Usage in WordPress: [sspserver_direct zone="123" text="Special Offer"]
?>
```

## Performance Considerations

### Response Time Optimization

- **Minimal Processing**: Direct endpoint processes requests with minimal overhead
- **Connection Pooling**: Reuse HTTP connections for better performance
- **CDN Integration**: Use CDN for superfailover URLs to reduce latency

### Caching Strategy

```javascript
// Client-side caching for repeated requests
class DirectAdCache {
  constructor(ttl = 300000) { // 5 minutes
    this.cache = new Map();
    this.ttl = ttl;
  }
  
  async getAdUrl(zoneId, params = {}) {
    const cacheKey = `${zoneId}-${JSON.stringify(params)}`;
    const cached = this.cache.get(cacheKey);
    
    if (cached && Date.now() - cached.timestamp < this.ttl) {
      return cached.url;
    }
    
    const url = await this.fetchAdUrl(zoneId, params);
    this.cache.set(cacheKey, {
      url,
      timestamp: Date.now()
    });
    
    return url;
  }
  
  async fetchAdUrl(zoneId, params) {
    const queryParams = new URLSearchParams({
      zone: zoneId,
      debug: 'true',
      noredirect: 'true',
      ...params
    });
    
    const response = await fetch(`/direct?${queryParams}`);
    const data = await response.json();
    return data.link || data.superfailover;
  }
}
```

## Best Practices

### 1. Always Configure Superfailover

Ensure your system has a superfailover URL configured to maintain monetization when no ads are available.

### 2. Use Tracking Parameters

Implement tracking parameters (`subid1-5`) for analytics and optimization:

```javascript
const trackingParams = {
  subid1: userId,
  subid2: sessionId,
  subid3: campaignId,
  subid4: sourceId,
  subid5: deviceType
};
```

### 3. Handle Errors Gracefully

Always implement fallback mechanisms for network errors or service unavailability.

### 4. Test with Debug Mode

Use debug mode during development to understand response behavior:

```bash
curl 'https://api.example.com/direct?zone=123&debug=true&noredirect=true'
```

### 5. Monitor Response Headers

Check for alternative link and failover headers to understand ad serving performance:

```javascript
fetch('/direct?zone=123', { method: 'HEAD' })
  .then(response => {
    const isAlternative = response.headers.get('X-Status-Alternative');
    const isFailover = response.headers.get('X-Status-Failover');
    
    // Handle different response types
    if (isFailover) {
      console.log('No ads available, showing fallback content');
    } else if (isAlternative) {
      console.log('Alternative content served');
    } else {
      console.log('Regular ad served');
    }
  });
```

## Summary

The direct endpoint provides a simple, fast, and reliable way to serve redirect-based advertising. Its minimal overhead and comprehensive error handling make it ideal for popunder campaigns, direct response marketing, and legacy system integration. The debug mode and tracking capabilities ensure you can monitor and optimize performance effectively.
