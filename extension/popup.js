document.addEventListener('DOMContentLoaded', async function() {
    const loginSection = document.getElementById('loginSection');
    const userSection = document.getElementById('userSection');
    const loginButton = document.getElementById('loginButton');
    const addProxyButton = document.getElementById('addProxyButton');
    const removeProxyButton = document.getElementById('removeProxyButton');
    const logoutButton = document.getElementById('logoutButton');
    const statusMessage = document.getElementById('statusMessage');
    const userName = document.getElementById('userName');
    const userEmail = document.getElementById('userEmail');

    // Check if user is already logged in
    const backendUrl = "http://88.198.127.219:3000"
    await checkAuthState();

    // Event listeners
    loginButton.addEventListener('click', handleLogin);
    addProxyButton.addEventListener('click', handleAddProxy);
    removeProxyButton.addEventListener('click', handleRemoveProxy);
    logoutButton.addEventListener('click', handleLogout);

    async function checkAuthState() {
        try {
            const result = await chrome.storage.local.get(['userToken', 'userInfo']);

            if (result.userToken) {
                // Call backend to verify token and update IPs
                const response = await fetch(`${backendUrl}/auth/verify`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({ token: result.userToken })
                });

                const data = await response.json();

                if (!data.success) {
                    // Token invalid or user not found → clear and show login
                    await clearAuthData();
                    showLoginSection();
                    return false; // Authentication failed
                }

                // Update user info locally (email + IPs)
                await chrome.storage.local.set({
                    userInfo: {
                        email: data.email,
                        name: data.name || data.email.split('@')[0], // Use name if available, otherwise extract from email
                        ips: data.ips
                    }
                });

                // Show user section with complete user info
                showUserSection({
                    email: data.email,
                    name: data.name || data.email.split('@')[0],
                    ips: data.ips
                });
                
                return true; // Authentication successful
            } else {
                showLoginSection();
                return false; // No token found
            }
        } catch (error) {
            console.error('Error checking auth state:', error);
            showLoginSection();
            return false; // Error occurred
        }
    }

    async function handleLogin() {
        try {
            showLoading(loginButton, 'Signing in...');

            // 1. Get Google OAuth token from Chrome Identity API
            const googleToken = await new Promise((resolve, reject) => {
                chrome.identity.getAuthToken({ interactive: true }, (token) => {
                    if (chrome.runtime.lastError) {
                        reject(chrome.runtime.lastError);
                    } else {
                        resolve(token);
                    }
                });
            });

            // 2. Send token to backend
            const response = await fetch(`${backendUrl}/login`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    google_token: googleToken
                })
            });

            const data = await response.json();
            console.log("response from backend:", data);

            if (!data.success) {
                showStatus('Login failed. Invalid user.', 'error');
                return;
            }

            // 3. Store backend JWT + user info locally
            const userInfo = {
                email: data.email,
                name: data.name || data.email.split('@')[0], // Extract name from email if not provided
                ips: data.ips || []
            };

            await chrome.storage.local.set({
                userToken: data.token,   // backend JWT
                userInfo: userInfo
            });

            // 4. Update UI
            showUserSection(userInfo);
            showStatus('Successfully signed in!', 'success');

        } catch (error) {
            console.error('Login error:', error);
            showStatus('Login failed. Please try again.', 'error');
        } finally {
            hideLoading(loginButton, 'Sign in with Google');
        }
    }

    async function handleAddProxy() {
        try {
            showLoading(addProxyButton, 'Adding...');
            
            const proxyConfig = {
                mode: "fixed_servers",
                rules: {
                    singleProxy: {
                        scheme: "http",
                        host: "88.198.127.219",
                        port: 8080
                    },
                    bypassList: ["localhost", "127.0.0.1"]
                }
            };

            await new Promise((resolve, reject) => {
                chrome.proxy.settings.set({
                    value: proxyConfig,
                    scope: 'regular'
                }, () => {
                    if (chrome.runtime.lastError) {
                        reject(chrome.runtime.lastError);
                    } else {
                        resolve();
                    }
                });
            });

            showStatus('Proxy set to localhost:8080', 'success');
            
        } catch (error) {
            console.error('Error setting proxy:', error);
            showStatus('Failed to set proxy. Please try again.', 'error');
        } finally {
            hideLoading(addProxyButton, 'Add Proxy');
        }
    }

    async function handleRemoveProxy() {
        try {
            showLoading(removeProxyButton, 'Removing...');
            
            await new Promise((resolve, reject) => {
                chrome.proxy.settings.clear({
                    scope: 'regular'
                }, () => {
                    if (chrome.runtime.lastError) {
                        reject(chrome.runtime.lastError);
                    } else {
                        resolve();
                    }
                });
            });

            showStatus('Proxy removed successfully', 'success');
            
        } catch (error) {
            console.error('Error removing proxy:', error);
            showStatus('Failed to remove proxy. Please try again.', 'error');
        } finally {
            hideLoading(removeProxyButton, 'Remove Proxy');
        }
    }

    async function handleLogout() {
        try {
            const result = await chrome.storage.local.get(['userToken']);
            
            if (result.userToken) {
                // Revoke the token
                chrome.identity.removeCachedAuthToken({ token: result.userToken });
            }

            await clearAuthData();
            showLoginSection();
            showStatus('Successfully signed out', 'success');
            
        } catch (error) {
            console.error('Logout error:', error);
            showStatus('Logout failed', 'error');
        }
    }

    async function clearAuthData() {
        await chrome.storage.local.remove(['userToken', 'userInfo']);
    }

    function showLoginSection() {
        loginSection.style.display = 'block';
        userSection.style.display = 'none';
    }

    function showUserSection(userInfo) {
        loginSection.style.display = 'none';
        userSection.style.display = 'block';
        
        // Properly set the username and email
        userName.textContent = userInfo.name || 'User';
        userEmail.textContent = userInfo.email || '';
    }

    function showLoading(button, text) {
        button.disabled = true;
        button.innerHTML = `<span class="loading"></span> ${text}`;
    }

    function hideLoading(button, originalText) {
        button.disabled = false;
        button.textContent = originalText;
    }

    function showStatus(message, type) {
        statusMessage.textContent = message;
        statusMessage.className = `status ${type}`;
        statusMessage.style.display = 'block';
        
        setTimeout(() => {
            statusMessage.style.display = 'none';
        }, 3000);
    }
});