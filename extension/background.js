// Background script for the Chrome extension
chrome.runtime.onInstalled.addListener(() => {
    console.log('Free Proxy Extension installed');
});

// Listen for proxy errors
chrome.proxy.onProxyError.addListener((details) => {
    console.error('Proxy error:', details);
    
    // Notify the user about proxy errors
    chrome.action.setBadgeText({
        text: '!'
    });
    chrome.action.setBadgeBackgroundColor({
        color: '#ff0000'
    });
    
    // Auto-clear badge after 10 seconds
    setTimeout(() => {
        chrome.action.setBadgeText({
            text: ''
        });
    }, 10000);
});

// Clear badge when extension is opened
chrome.action.onClicked.addListener(() => {
    chrome.action.setBadgeText({
        text: ''
    });
});

// Update badge based on proxy status
async function updateProxyBadge() {
    try {
        chrome.proxy.settings.get({}, (config) => {
            if (config.value.mode === 'fixed_servers') {
                // Proxy is active
                chrome.action.setBadgeText({
                    text: 'ON'
                });
                chrome.action.setBadgeBackgroundColor({
                    color: '#4CAF50'
                });
            } else {
                // Proxy is inactive
                chrome.action.setBadgeText({
                    text: 'OFF'
                });
                chrome.action.setBadgeBackgroundColor({
                    color: '#757575'
                });
            }
        });
    } catch (error) {
        console.error('Error updating badge:', error);
    }
}

// Handle extension startup - update badge
chrome.runtime.onStartup.addListener(() => {
    updateProxyBadge();
});

// Initial badge update
updateProxyBadge();

// Handle messages from popup
chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
    switch (request.action) {
        case 'getProxyStatus':
            chrome.proxy.settings.get({}, (config) => {
                const isProxySet = config.value.mode === 'fixed_servers' &&
                    config.value.rules &&
                    config.value.rules.singleProxy &&
                    config.value.rules.singleProxy.host === 'proxybackend.bhar.xyz' &&
                    config.value.rules.singleProxy.port === 8080;
                
                sendResponse({ isProxySet, config: config.value });
            });
            return true;
        
        case 'updateBadge':
            updateProxyBadge();
            sendResponse({ success: true });
            return true;
        
        default:
            sendResponse({ error: 'Unknown action' });
    }
});

// Listen for proxy settings changes to update badge
chrome.proxy.settings.onChange.addListener((details) => {
    console.log('Proxy settings changed:', details);
    updateProxyBadge();
});