// Background script for the Chrome extension

chrome.runtime.onInstalled.addListener(() => {
    console.log('Proxy OAuth Extension installed');
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
});

// Clear badge when extension is opened
chrome.action.onClicked.addListener(() => {
    chrome.action.setBadgeText({
        text: ''
    });
});

// Listen for storage changes to handle token expiration
chrome.storage.onChanged.addListener((changes, namespace) => {
    if (namespace === 'local' && changes.userToken) {
        if (!changes.userToken.newValue && changes.userToken.oldValue) {
            // Token was removed, clear proxy settings
            chrome.proxy.settings.clear({
                scope: 'regular'
            }, () => {
                console.log('Proxy settings cleared due to logout');
            });
        }
    }
});

// Handle extension startup - check if user is logged in and proxy should be active
chrome.runtime.onStartup.addListener(async () => {
    try {
        const result = await chrome.storage.local.get(['userToken', 'proxyEnabled']);
        
        if (result.userToken && result.proxyEnabled) {
            // TODO: Verify token with backend
            // const isValid = await verifyTokenWithBackend(result.userToken);
            // if (!isValid) {
            //     await chrome.storage.local.remove(['userToken', 'userInfo']);
            //     return;
            // }
            
            console.log('User is logged in, proxy settings may be active');
        }
    } catch (error) {
        console.error('Error checking auth state on startup:', error);
    }
});

// Function to verify token with backend (placeholder)
async function verifyTokenWithBackend(token) {
    // TODO: Implement actual backend verification
    // try {
    //     const response = await fetch('YOUR_BACKEND_URL/verify-token', {
    //         method: 'POST',
    //         headers: {
    //             'Content-Type': 'application/json',
    //             'Authorization': `Bearer ${token}`
    //         }
    //     });
    //     
    //     return response.ok;
    // } catch (error) {
    //     console.error('Backend verification failed:', error);
    //     return false;
    // }
    
    return true; // For now, assume token is valid
}

// Handle messages from popup or content scripts
chrome.runtime.onMessage.addListener((request, sender, sendResponse) => {
    switch (request.action) {
        case 'verifyToken':
            verifyTokenWithBackend(request.token)
                .then(isValid => sendResponse({ isValid }))
                .catch(error => sendResponse({ isValid: false, error: error.message }));
            return true; // Will respond asynchronously
            
        case 'getProxyStatus':
            chrome.proxy.settings.get({}, (config) => {
                const isProxySet = config.value.mode === 'fixed_servers' &&
                                 config.value.rules &&
                                 config.value.rules.singleProxy &&
                                 config.value.rules.singleProxy.host === 'localhost' &&
                                 config.value.rules.singleProxy.port === 8080;
                
                sendResponse({ isProxySet, config: config.value });
            });
            return true;
            
        default:
            sendResponse({ error: 'Unknown action' });
    }
});

// Monitor network requests for debugging (optional)
chrome.webRequest.onBeforeRequest.addListener(
    (details) => {
        // Log requests when proxy is active (for debugging)
        // console.log('Request:', details.url);
    },
    { urls: ["<all_urls>"] },
    []
);